import json
import random
from transformers import AutoTokenizer, AutoModelForCausalLM, LlamaTokenizer
import torch
from tqdm import tqdm

import jieba
from nltk.translate.bleu_score import sentence_bleu, SmoothingFunction
from rouge import Rouge
from sklearn.metrics import accuracy_score, f1_score


from nltk.translate.bleu_score import sentence_bleu
from rouge import Rouge
from flask import Flask, request, jsonify
import torch

app = Flask(__name__)


def calculate_bleu_chinese(reference, candidate):
    ref_words = list(jieba.cut(reference))
    cand_words = list(jieba.cut(candidate))

    # 使用平滑函数
    smoothie = SmoothingFunction().method1
    return sentence_bleu([ref_words], cand_words, smoothing_function=smoothie)


def calculate_rouge_chinese(reference, candidate):
    # 将中文文本分割为字符列表
    reference_chars = list(jieba.cut(reference))
    candidate_chars = list(jieba.cut(candidate))

    # 将字符列表转换为字符串，每个字符后加空格以符合Rouge库的处理方式
    reference_joined = ' '.join(reference_chars)
    candidate_joined = ' '.join(candidate_chars)

    # 计算Rouge得分
    rouge = Rouge()
    scores = rouge.get_scores(candidate_joined, reference_joined)
    return scores[0]['rouge-l']['f']


def calculate_accuracy_and_f1(references, candidates):
    accuracy = accuracy_score(references, candidates)
    f1 = f1_score(references, candidates, average='weighted')
    return accuracy, f1


def load_model(model_name_or_path, device):
    """
    加载模型和分词器
    """
    # 检查是否为 Llama 模型
    if "llama" in model_name_or_path.lower():
        # 加载 Llama Tokenizer
        tokenizer = LlamaTokenizer.from_pretrained(model_name_or_path)
    else:
        # 为其他模型加载 AutoTokenizer
        tokenizer = AutoTokenizer.from_pretrained(model_name_or_path)

    # 确保tokenizer具有 pad_token
    if tokenizer.pad_token is None:
        tokenizer.add_special_tokens({'pad_token': '[PAD]'})

    # 加载模型
    model = AutoModelForCausalLM.from_pretrained(model_name_or_path)

    model.to(device)
    return model, tokenizer


def load_dataset(dataset_path):
    """
    加载测试集
    """
    with open(dataset_path, 'r', encoding='utf-8') as file:
        dataset = [json.loads(line) for line in file.readlines()]
    return dataset


def generate_answers_batch(model, tokenizer, batch_questions, max_length, device):
    """
    使用模型批量生成答案
    """
    inputs = tokenizer(batch_questions, padding=True,
                       return_tensors='pt').to(device)

    # 使用模型生成答案
    outputs = model.generate(
        **inputs,
        max_length=max_length,
        pad_token_id=tokenizer.pad_token_id)

    results = []
    for i in range(len(batch_questions)):
        output = tokenizer.decode(outputs[i], skip_special_tokens=True)
        if output.startswith(batch_questions[i]):
            output = output[len(batch_questions[i]):]
        output = output.strip()
        results.append(output)

    # 解码生成的答案
    return results


def write_evaluation_results_json(eval_output_path, hyperparams, evaluation_results, dataset, candidates, scores):
    """
    将评测结果以JSON格式写入文件
    """
    results = {
        "config": hyperparams,
        "overall_evaluation_metrics": evaluation_results,
        "detailed_results": []
    }

    for item, cand, score in zip(dataset, candidates, scores):
        question = item['messages'][0]['content']
        reference = item['messages'][1]['content']
        detailed_result = {
            "question": question,
            "reference": reference,
            "model_output": cand,
            "scores": score
        }
        results["detailed_results"].append(detailed_result)

    with open(eval_output_path, 'w', encoding='utf-8') as file:
        json.dump(results, file, indent=4, ensure_ascii=False)

    print(f"evaluation results written to {eval_output_path}")


def evaluate_model(model_name_or_path, dataset_path, evaluation_metrics, max_seq_len, per_device_batch_size, gpu_id, output_path):
    """
    真实的评估函数
    """
    # 移动模型到指定的GPU
    device = torch.device(
        f'cuda:{gpu_id}' if torch.cuda.is_available() else 'cpu')
    # 加载模型和分词器
    model, tokenizer = load_model(model_name_or_path, device)
    print("model and tokenizer loaded")
    # 加载测试集
    dataset = load_dataset(dataset_path)
    print("dataset loaded")

    # 初始化评估结果
    evaluation_results = {}

    # 准备问题列表
    questions = [item['messages'][0]['content'] for item in dataset]
    references = [item['messages'][1]['content'] for item in dataset]
    scores = []

    # 分批处理问题
    candidates = []
    for i in tqdm(range(0, len(questions), per_device_batch_size), desc="Evaluating", unit="batch"):
        batch_questions = questions[i:i+per_device_batch_size]
        batch_references = references[i:i+per_device_batch_size]
        batch_answers = generate_answers_batch(
            model, tokenizer, batch_questions, max_seq_len, device)

        candidates.extend(batch_answers)

        batch_scores = []
        for reference, answer in zip(batch_references, batch_answers):
            single_score = {}
            if "Acc" in evaluation_metrics:
                single_score["Acc"] = int(reference == answer)
            if "F1" in evaluation_metrics:
                single_score["F1"] = int(reference == answer)
            if "BLEU" in evaluation_metrics:
                single_score["BLEU"] = calculate_bleu_chinese(
                    reference, answer)
            if "Rouge" in evaluation_metrics:
                single_score["Rouge"] = calculate_rouge_chinese(
                    reference, answer)
            batch_scores.append(single_score)
        scores.extend(batch_scores)

    # 计算平均值
    for metric in evaluation_metrics:
        evaluation_results[metric] = sum(
            [score[metric] for score in scores]) / len(scores)

    print("answers generated")

    # 写入评测结果
    hyperparams = {
        "model_name_or_path": model_name_or_path,
        "dataset_path": dataset_path,
        "evaluation_metrics": evaluation_metrics,
        "max_seq_len": max_seq_len,
        "per_device_batch_size": per_device_batch_size,
        "gpu_id": gpu_id
    }
    eval_output_path = output_path + "/eval_results.json"

    write_evaluation_results_json(eval_output_path, hyperparams,
                                  evaluation_results, dataset, candidates, scores)

    # 清理模型以释放内存
    del model
    del tokenizer
    torch.cuda.empty_cache()

    # 构建返回的JSON格式结果
    result = {
        "code": 0,
        "msg": "Success",
        "data": {
            "evaluation_results": evaluation_results
        }
    }

    return result


@app.route('/evaluate', methods=['POST'])
def evaluate():
    data = request.json
    # 解析请求数据并进行相应处理
    model_name_or_path = data['model_name_or_path']
    dataset_path = data['dataset_path']
    evaluation_metrics = data['evaluation_metrics']
    max_seq_len = data['max_seq_len']
    per_device_batch_size = data['per_device_batch_size']
    gpu_id = data['gpu_id']
    output_path = data['output_path']

    # 评估逻辑
    result = evaluate_model(
        model_name_or_path=model_name_or_path,
        dataset_path=dataset_path,
        evaluation_metrics=evaluation_metrics,
        max_seq_len=max_seq_len,
        per_device_batch_size=per_device_batch_size,
        gpu_id=gpu_id,
        output_path=output_path
    )

    return jsonify(result)


if __name__ == "__main__":
    app.run(debug=True)

