# LLMOps deepspeed训练后端

- 使用Docker启动

1. **构建Docker镜像**:
   在项目的根目录（包含Dockerfile的地方）运行以下命令来构建Docker镜像

```bash
docker build -t llmops-deepspeed-backend .
```

2. **运行Docker容器**:
   使用以下命令运行Docker容器：

```bash
docker run --gpus all --shm-size 64gb -p 6006:6006 5000:5000 \
-v /Your/Path/mnt:/app/mnt/ \
llmops-deepspeed-backend
```

- 如遇GPU访问的问题，可尝试安装 NVIDIA Container Toolkit
- 在运行Docker容器时，使用 -v 参数来挂载目录，将`/Your/Path/mnt`替换为你的目录
- 挂载的目录结构应如下：

```bash
mnt
   - datasets
      - train.jsonl
      - test.jsonl
   - Llama-2-13b-chat-hf
   - output_path
```

其中，datasets是数据集，Llama-2-13b-chat-hf是预训练模型，output_path是训练后的模型保存路径
数据集的格式为jsonl，例如：

```json
{"messages": [{"role": "user", "content": "ping"}, {"role": "assistant", "content": "Pong!"}]}
```


## 评估模块
- eval.py中启动
1. 请求格式
```bash
curl --location 'http://localhost:5000/evaluate' \
--header 'Content-Type: application/json' \
--data '{
    "model_name_or_path": "app/mnt/Llama-2-13b-chat-hf",
    "dataset_path": "app/mnt/datasets/test.jsonl",
    "evaluation_metrics": ["Acc", "F1", "BLEU", "Rouge"],
    "max_seq_len": 128,
    "per_device_batch_size": 1,
    "gpu_id": 7,
    "output_path": "app/mnt/output_path"
}
'
```

- 两个经验公式在test.ipynb
   - memory_usage: 评估过程中的显存使用情况和
   valuation_duration: 整个评估过程的持续时间