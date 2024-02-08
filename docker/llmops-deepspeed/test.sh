clear
###
 # @Author: whateverisnottaken johnzhaozitian@gmail.com
 # @Date: 2024-01-09 07:36:53
 # @LastEditors: whateverisnottaken johnzhaozitian@gmail.com
 # @LastEditTime: 2024-01-11 08:47:39
 # @FilePath: /fantianming/xm/LLMOpsDeepSpeed后端_1218/src/test.sh
 # @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
### 
ps aux | grep llmops_deepspeed_main.py | awk '{print $2}' | xargs kill -9
MNT_PATH=/home/calf/ssd/workspace/fantianming/xm/LLMOpsDeepSpeed后端_1218/mnt

export CUDA_VISIBLE_DEVICES=0,1,2,3,4,5,6,7
# export CUDA_VISIBLE_DEVICES=0,6,7

# 转换
python jsonl_to_arrow_format.py \
	--base_path "$MNT_PATH" \

# DeepSpeed Team
ZERO_STAGE=2
deepspeed llmops_deepspeed_main.py \
	--data_path "$MNT_PATH"/formatted_datasets \
    --data_output_path "$MNT_PATH"/output_path/data_output \
	--data_split 10,0,0 \
	--model_name_or_path "$MNT_PATH"/Llama-2-13b-chat-hf \
	--per_device_train_batch_size 36 \
	--per_device_eval_batch_size 10 \
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
