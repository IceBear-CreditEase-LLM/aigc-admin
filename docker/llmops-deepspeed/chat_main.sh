#!/bin/bash
clear

# 打印所有传入的参数（用于调试）
echo "Received arguments: $@"

MNT_PATH=$1

mkdir -p "$MNT_PATH"/formatted_datasets
mkdir -p "$MNT_PATH"/output_path/data_output
mkdir -p "$MNT_PATH"/output_path/tensorboard

export CUDA_VISIBLE_DEVICES=0,1,2,3

# 转换
python3 jsonl_to_arrow_format.py \
	--base_path "$MNT_PATH"

# DeepSpeed Team
ZERO_STAGE=2
deepspeed llmops_deepspeed_main.py \
	--data_path ./formatted_datasets \
  --data_output_path "$MNT_PATH"/output_path/data_output \
	--data_split 10,0,0 \
	--model_name_or_path "$MNT_PATH"/Llama-2-13b-chat-hf \
	--per_device_train_batch_size 1 \
	--per_device_eval_batch_size 1 \
	--max_seq_len 512 \
	--learning_rate 1e-5 \
	--weight_decay 0. \
	--num_train_epochs 1  \
	--gradient_accumulation_steps 1 \
	--lr_scheduler_type cosine \
	--num_warmup_steps 0 \
	--seed 42 \
	--gradient_checkpointing \
	--zero_stage $ZERO_STAGE \
	--deepspeed \
	--print_loss \
	--output_dir "$MNT_PATH"/output_path \
	--start_from_step -1 \
	--save_per_steps 100 \
	--tensorboard_path "$MNT_PATH"/output_path/tensorboard \
	--tensorboard_port 6007 \
	--enable_tensorboard| tee $MNT_PATH/output_path/log.txt \