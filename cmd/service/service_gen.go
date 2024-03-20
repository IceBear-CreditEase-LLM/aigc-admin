package service

import (
	"database/sql"
	"fmt"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/encode"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/types"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/util"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"strings"
)

var (
	generateCmd = &cobra.Command{
		Use:               "generate command <args> [flags]",
		Short:             "生成命令",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `## 生成命令
可用的配置类型：
[table, init-data]

aigc-admin generate -h
`,
	}

	genTableCmd = &cobra.Command{
		Use:               `table <args> [flags]`,
		Short:             "生成数据库表",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
aigc-admin generate table all
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 关闭资源连接
			defer func() {
				_ = level.Debug(logger).Log("db", "close", "err", db.Close())
			}()

			if len(args) > 0 && args[0] == "all" {
				return generateTable()
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			logger = log.NewLogfmtLogger(os.Stdout)
			// 连接数据库
			var dbErr error
			if strings.EqualFold(dbDrive, "mysql") {
				dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
					mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
				sqlDB, err := sql.Open("mysql", dbUrl)
				if err != nil {
					_ = level.Error(logger).Log("sql", "Open", "err", err.Error())
					return err
				}
				gormDB, err = gorm.Open(mysql.New(mysql.Config{
					Conn:              sqlDB,
					DefaultStringSize: 255,
				}), &gorm.Config{
					DisableForeignKeyConstraintWhenMigrating: true,
				})
				if dbErr != nil {
					_ = level.Error(logger).Log("db", "connect", "err", dbErr.Error())
					dbErr = encode.ErrServerStartDbConnect.Wrap(dbErr)
					return dbErr
				}
				//gormDB.Statement.Clauses["soft_delete_enabled"] = clause.Clause{}
				db, dbErr = gormDB.DB()
				if dbErr != nil {
					_ = level.Error(logger).Log("gormDB", "DB", "err", dbErr.Error())
					dbErr = encode.ErrServerStartDbConnect.Wrap(dbErr)
					return dbErr
				}
				_ = level.Debug(logger).Log("mysql", "connect", "success", true)
			} else if strings.EqualFold(dbDrive, "sqlite") {
				_ = os.MkdirAll(fmt.Sprintf("%s/database", serverStoragePath), 0755)
				sqliteDB, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s/database/aigc.db", serverStoragePath)), &gorm.Config{
					DisableForeignKeyConstraintWhenMigrating: true,
				})
				if err != nil {
					_ = level.Error(logger).Log("sqlite", "connect", "err", err.Error())
					return err
				}
				db, dbErr = sqliteDB.DB()
				if dbErr != nil {
					_ = level.Error(logger).Log("sqlite", "connect", "err", dbErr.Error())
					return dbErr
				}
				_ = level.Debug(logger).Log("sqlite", "connect", "success", true)
			} else {
				err = fmt.Errorf("db drive not support: %s", dbDrive)
				_ = level.Error(logger).Log("db", "drive", "err", err.Error())
				return err
			}
			return nil
		},
	}
)

func generateTable() (err error) {
	_ = logger.Log("migrate", "table", "Chat", gormDB.AutoMigrate(types.Chat{}))
	_ = logger.Log("migrate", "table", "ChatAllowUser", gormDB.AutoMigrate(types.ChatAllowUser{}))
	_ = logger.Log("migrate", "table", "ChatConversation", gormDB.AutoMigrate(types.ChatConversation{}))
	_ = logger.Log("migrate", "table", "ChatSystemPrompt", gormDB.AutoMigrate(types.ChatSystemPrompt{}))
	_ = logger.Log("migrate", "table", "ChatPromptTypes", gormDB.AutoMigrate(types.ChatPromptTypes{}))
	_ = logger.Log("migrate", "table", "ChatChannels", gormDB.AutoMigrate(types.ChatChannels{}))
	_ = logger.Log("migrate", "table", "ChatPrompts", gormDB.AutoMigrate(types.ChatPrompts{}))
	_ = logger.Log("migrate", "table", "ChatChannelModels", gormDB.AutoMigrate(types.ChatChannelModels{}))
	_ = logger.Log("migrate", "table", "ChatMessages", gormDB.AutoMigrate(types.ChatMessages{}))
	_ = logger.Log("migrate", "table", "Dataset", gormDB.AutoMigrate(types.Dataset{}))
	_ = logger.Log("migrate", "table", "DatasetSample", gormDB.AutoMigrate(types.DatasetSample{}))
	_ = logger.Log("migrate", "table", "Assistants", gormDB.AutoMigrate(types.Assistants{}))
	_ = logger.Log("migrate", "table", "Tools", gormDB.AutoMigrate(types.Tools{}))
	_ = logger.Log("migrate", "table", "AssistantToolAssociations", gormDB.AutoMigrate(types.AssistantToolAssociations{}))
	_ = logger.Log("migrate", "table", "Files", gormDB.AutoMigrate(types.Files{}))
	_ = logger.Log("migrate", "table", "SysAudit", gormDB.AutoMigrate(types.SysAudit{}))
	_ = logger.Log("migrate", "table", "FineTuningTrainJob", gormDB.AutoMigrate(types.FineTuningTrainJob{}))
	_ = logger.Log("migrate", "table", "FineTuningTemplate", gormDB.AutoMigrate(types.FineTuningTemplate{}))
	_ = logger.Log("migrate", "table", "Tenants", gormDB.AutoMigrate(types.Tenants{}))
	_ = logger.Log("migrate", "table", "Models", gormDB.AutoMigrate(types.Models{}))
	_ = logger.Log("migrate", "table", "SysDict", gormDB.AutoMigrate(types.SysDict{}))
	_ = logger.Log("migrate", "table", "ModelDeploy", gormDB.AutoMigrate(types.ModelDeploy{}))
	_ = logger.Log("migrate", "table", "LLMEvalResults", gormDB.AutoMigrate(types.LLMEvalResults{}))
	_ = logger.Log("migrate", "table", "DatasetDocument", gormDB.AutoMigrate(types.DatasetDocument{}))
	_ = logger.Log("migrate", "table", "DatasetDocumentSegment", gormDB.AutoMigrate(types.DatasetDocumentSegment{}))
	_ = logger.Log("migrate", "table", "DatasetAnnotationTask", gormDB.AutoMigrate(types.DatasetAnnotationTask{}))
	_ = logger.Log("migrate", "table", "DatasetAnnotationTaskSegment", gormDB.AutoMigrate(types.DatasetAnnotationTaskSegment{}))
	_ = logger.Log("migrate", "table", "ModelEvaluate", gormDB.AutoMigrate(types.ModelEvaluate{}))
	//err = initData()
	//if err != nil {
	//	return err
	//}
	return
}

// 初始化数据
func initData() (err error) {
	tenant := types.Tenants{
		Name:           "系统租户",
		PublicTenantID: uuid.New().String(),
		ContactEmail:   serverAdminUser,
	}
	_ = logger.Log("init", "data", "SysDict", gormDB.Create(&tenant).Error)
	password, err := bcrypt.GenerateFromPassword([]byte(serverAdminPass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_ = logger.Log("init", "data", "account", gormDB.Save(&types.Accounts{
		Email:        serverAdminUser,
		Nickname:     "系统管理员",
		Language:     "zh",
		IsLdap:       false,
		PasswordHash: string(password),
		Status:       true,
		Tenants:      []types.Tenants{tenant},
	}).Error)
	if serverChannelKey == "" {
		serverChannelKey = "sk-" + string(util.Krand(48, util.KC_RAND_KIND_ALL))
	}
	_ = logger.Log("init", "data", "ChatChannels", gormDB.Create(&types.ChatChannels{
		Name:       "default",
		Alias:      "默认渠道",
		Remark:     "默认渠道",
		Quota:      10000,
		Models:     "default",
		OnlyOpenAI: false,
		ApiKey:     serverChannelKey,
		Email:      serverAdminUser,
		TenantId:   tenant.ID,
	}).Error)

	_ = logger.Log("init", "data", "sys_dict", gormDB.Exec(initSysDictSql).Error)
	_ = logger.Log("init", "data", "finetuning_template", gormDB.Exec(ftTemplateSql).Error)
	_ = logger.Log("init", "data", "models", gormDB.Exec(modelSql).Error)
	return err
}

var (
	modelSql = `INSERT INTO models (id, created_at, updated_at, deleted_at, provider_name, model_type, model_name, max_tokens, is_private, is_fine_tuning, enabled, remark, parameters, last_operator, base_model_name, replicas, label, k8s_cluster, inferred_type, gpu, cpu, memory)
VALUES
	(2, '2024-02-04 13:02:48.112', '2024-03-19 14:08:35.667', NULL, 'OpenAI', 'text-generation', 'gpt-3.5-turbo', 4096, 0, 0, 1, 'OpenAI GPT-3.5-turbo', 20.00, '', NULL, 1, NULL, NULL, NULL, 0, 0, 1),
	(3, '2024-03-18 17:34:59.542', '2024-03-19 14:51:35.630', NULL, 'LocalAI', 'text-generation', 'qwen1.5-0.5b', 32768, 1, 0, 0, '', 0.50, 'admin', '', 1, '', '', '', 0, 0, 1),
	(4, '2024-03-19 14:03:11.073', '2024-03-19 14:41:08.194', NULL, 'LocalAI', 'text-generation', 'qwen1.5-1.8b', 32768, 0, 0, 0, '', 1.80, 'admin', '', 1, '', '', '', 0, 0, 1),
	(5, '2024-03-19 14:03:34.619', '2024-03-19 14:41:41.709', NULL, 'LocalAI', 'text-generation', 'qwen1.5-1.8b-chat', 32768, 0, 0, 0, '', 1.80, 'admin', '', 1, '', '', '', 0, 0, 1),
	(6, '2024-03-19 14:03:51.375', '2024-03-19 14:41:36.354', NULL, 'LocalAI', 'text-generation', 'qwen1.5-4b', 32768, 0, 0, 0, '', 3.98, 'admin', '', 1, '', '', '', 0, 0, 1),
	(7, '2024-03-19 14:04:11.425', '2024-03-19 14:41:11.423', NULL, 'LocalAI', 'text-generation', 'qwen1.5-4b-chat', 32768, 0, 0, 0, '', 3.98, 'admin', '', 1, '', '', '', 0, 0, 1),
	(8, '2024-03-19 14:04:29.257', '2024-03-19 14:41:18.790', NULL, 'LocalAI', 'text-generation', 'qwen1.5-7b', 32768, 0, 0, 0, '', 7.20, 'admin', '', 1, '', '', '', 0, 0, 1),
	(9, '2024-03-19 14:04:45.241', '2024-03-19 14:41:24.050', NULL, 'LocalAI', 'text-generation', 'qwen1.5-7b-chat', 32768, 0, 0, 0, '', 7.20, 'admin', '', 1, '', '', '', 0, 0, 1),
	(10, '2024-03-19 14:05:04.519', '2024-03-19 14:41:27.394', NULL, 'LocalAI', 'text-generation', 'qwen1.5-14b', 32768, 0, 0, 0, '', 14.20, 'admin', '', 1, '', '', '', 0, 0, 1),
	(11, '2024-03-19 14:05:27.624', '2024-03-19 14:41:47.741', NULL, 'LocalAI', 'text-generation', 'qwen1.5-14b-chat', 32768, 0, 0, 0, '', 14.20, 'admin', '', 1, '', '', '', 0, 0, 1),
	(12, '2024-03-19 14:06:26.666', '2024-03-19 14:41:33.633', NULL, 'LocalAI', 'text-generation', 'qwen1.5-72b', 32768, 0, 0, 0, '', 72.30, 'admin', '', 1, '', '', '', 0, 0, 1),
	(13, '2024-03-19 14:06:43.121', '2024-03-19 14:41:30.391', NULL, 'LocalAI', 'text-generation', 'qwen1.5-72b-chat', 32768, 0, 0, 0, '', 72.30, 'admin', '', 1, '', '', '', 0, 0, 1),
	(14, '2024-03-19 14:08:27.352', '2024-03-19 14:41:01.068', NULL, 'LocalAI', 'text-generation', 'qwen-plus', 32768, 0, 0, 0, '', 14.20, 'admin', '', 1, '', '', '', 0, 0, 1);
`

	ftTemplateSql = `INSERT INTO fine_tuning_template (id, created_at, updated_at, deleted_at, name, base_model, content, params, train_image, remark, base_model_path, script_file, output_dir, max_tokens, lora, enabled, template_type)
VALUES
	(2, '2024-03-18 17:36:17.927', '2024-03-19 14:34:05.091', NULL, 'qwen1.5-0.5b', 'qwen1.5-0.5b', '#!/bin/bash\n\nMODEL_WORKER=fastchat.serve.model_worker\nCONTROLLER_ADDRESS=http://fschat-controller:21001\nMODEL_NAME={{.modelName}}\nMODEL_PATH={{.modelPath}}\nHTTP_PORT={{.port}}\nQUANTIZATION={{.quantization}}\nNUM_GPUS={{.numGpus}}\nMAX_GPU_MEMORY={{.maxGpuMemory}}\nVLLM={{.vllm}}\nINFERRED_TYPE={{.inferredType}}\nOS_TYPE=$(uname)\n\n# awq 量化配置\n# awq_wbits\n# awq_groupsize\n\n# 并发限制\n# limit-worker-concurrency\n\n# gptq 量化配置\n# gptq_wbits\n# gptq_groupsize\n# gptq_act_order\n\n# MODEL_WORKER \nif [ \"$VLLM\" == \"true\" ]; then\n    MODEL_WORKER=\"fastchat.serve.vllm_worker\"\nfi\n\n# 量化配置\nif [ \"$QUANTIZATION\" == \"8bit\" ]; then\n    QUANTIZATION=\"--load-8bit\"\nelse\n    QUANTIZATION=\"\"\nfi\n\n# NUM_GPUS\nif [ \"$NUM_GPUS\" -gt 0 ]; then\n    NUM_GPUS=\"--num-gpus $NUM_GPUS\"\nelse\n    NUM_GPUS=\"\"\nfi\n\n# CPU推理CPU，mps\nif [ \"$INFERRED_TYPE\" == \"cpu\" ] && [ \"$OS_TYPE\" == \"Darwin\" ]; then\n    DEVICE_OPTION=\"--device mps\"\nelif [ \"$INFERRED_TYPE\" == \"cpu\" ]; then\n    DEVICE_OPTION=\"--device cpu\"\nelse\n    DEVICE_OPTION=\"\"\nfi\n\n# MAX_GPU_MEMORY\nif [ \"$MAX_GPU_MEMORY\" -gt 0 ]; then\n    MAX_GPU_MEMORY=\"--max-gpu-memory ${MAX_GPU_MEMORY}GiB\"\nelse\n    MAX_GPU_MEMORY=\"\"\nfi\n\npython3 -m $MODEL_WORKER --host 0.0.0.0 --port $HTTP_PORT \\\n    --controller-address $CONTROLLER_ADDRESS \\\n    --worker-address http://$MY_POD_IP:$HTTP_PORT \\\n    --model-name $MODEL_NAME \\\n    --model-path $MODEL_PATH \\\n    $QUANTIZATION $NUM_GPUS $MAX_GPU_MEMORY $DEVICE_OPTION\n', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-0-5b', '/app/start.sh', '/data/ft-model/', 32768, 0, 1, 'inference'),
	(3, '2024-03-19 14:13:04.852', '2024-03-19 14:38:20.662', NULL, 'qwen1.5-0.5b-train', 'qwen1.5-0.5b', '#!/bin/bash\nexport AUTH=sk-001\nexport JOB_ID={{.JobId}}\n\nexport CUDA_DEVICE_MAX_CONNECTIONS=1\n\n\nGPUS_PER_NODE={{.ProcPerNode}}\nNNODES=1\nNODE_RANK=0\nMASTER_ADDR=localhost\nMASTER_PORT={{.MasterPort}}\nUSE_LORA={{.Lora}}\nQ_LORA=False\n\nMODEL=\"{{.BaseModelPath}}\" # Set the path if you do not want to load from huggingface directly\n# ATTENTION: specify the path to your training data, which should be a json file consisting of a list of conversations.\n# See the section for finetuning in README for more information.\nDATA=\"{{.DataPath}}\"\n# 验证集\nEVAL_DATA=\"{{.ValidationFile}}\"\nDS_CONFIG_PATH=\"ds_config_zero3.json\"\n\nDISTRIBUTED_ARGS=\"\n    --nproc_per_node $GPUS_PER_NODE \\\n    --nnodes $NNODES \\\n    --node_rank $NODE_RANK \\\n    --master_addr $MASTER_ADDR \\\n    --master_port $MASTER_PORT\n\"\nif [ \"$USE_LORA\" == \"true\" ]; then\n    USE_LORA=True\n    DS_CONFIG_PATH=\"ds_config_zero2.json\"\nelse\n    USE_LORA=False\n    DS_CONFIG_PATH=\"ds_config_zero3.json\"\nfi\n\nmkdir -p /data/train-data/\nwget -O {{.DataPath}} {{.FileUrl}}\n\ntorchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \\\n    --model_name_or_path $MODEL \\\n    --data_path $DATA \\\n    --bf16 True \\\n    --output_dir {{.OutputDir}} \\\n    --num_train_epochs {{.TrainEpoch}} \\\n    --per_device_train_batch_size {{.TrainBatchSize}} \\\n    --per_device_eval_batch_size {{.EvalBatchSize}} \\\n    --gradient_accumulation_steps {{.AccumulationSteps}} \\\n    --evaluation_strategy \"no\" \\\n    --save_strategy \"steps\" \\\n    --save_steps 1000 \\\n    --save_total_limit 10 \\\n    --learning_rate {{.LearningRate}} \\\n    --weight_decay 0.1 \\\n    --adam_beta2 0.95 \\\n    --warmup_ratio 0.01 \\\n    --lr_scheduler_type \"cosine\" \\\n    --logging_steps 1 \\\n    --report_to \"none\" \\\n    --model_max_length {{.ModelMaxLength}} \\\n    --gradient_checkpointing True \\\n    --lazy_preprocess True \\\n    --use_lora ${USE_LORA} \\\n    --q_lora ${Q_LORA} \\\n    --deepspeed $DS_CONFIG_PATH', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-0-5b', '/app/finetune.py', '/data/ft-model/', 32768, 0, 1, 'train'),
	(4, '2024-03-18 17:36:17.927', '2024-03-19 14:33:55.959', NULL, 'qwen1.5-1.8b', 'qwen1.5-1.8b', '#!/bin/bash\n\nMODEL_WORKER=fastchat.serve.model_worker\nCONTROLLER_ADDRESS=http://fschat-controller:21001\nMODEL_NAME={{.modelName}}\nMODEL_PATH={{.modelPath}}\nHTTP_PORT={{.port}}\nQUANTIZATION={{.quantization}}\nNUM_GPUS={{.numGpus}}\nMAX_GPU_MEMORY={{.maxGpuMemory}}\nVLLM={{.vllm}}\nINFERRED_TYPE={{.inferredType}}\nOS_TYPE=$(uname)\n\n# awq 量化配置\n# awq_wbits\n# awq_groupsize\n\n# 并发限制\n# limit-worker-concurrency\n\n# gptq 量化配置\n# gptq_wbits\n# gptq_groupsize\n# gptq_act_order\n\n# MODEL_WORKER \nif [ \"$VLLM\" == \"true\" ]; then\n    MODEL_WORKER=\"fastchat.serve.vllm_worker\"\nfi\n\n# 量化配置\nif [ \"$QUANTIZATION\" == \"8bit\" ]; then\n    QUANTIZATION=\"--load-8bit\"\nelse\n    QUANTIZATION=\"\"\nfi\n\n# NUM_GPUS\nif [ \"$NUM_GPUS\" -gt 0 ]; then\n    NUM_GPUS=\"--num-gpus $NUM_GPUS\"\nelse\n    NUM_GPUS=\"\"\nfi\n\n# CPU推理CPU，mps\nif [ \"$INFERRED_TYPE\" == \"cpu\" ] && [ \"$OS_TYPE\" == \"Darwin\" ]; then\n    DEVICE_OPTION=\"--device mps\"\nelif [ \"$INFERRED_TYPE\" == \"cpu\" ]; then\n    DEVICE_OPTION=\"--device cpu\"\nelse\n    DEVICE_OPTION=\"\"\nfi\n\n# MAX_GPU_MEMORY\nif [ \"$MAX_GPU_MEMORY\" -gt 0 ]; then\n    MAX_GPU_MEMORY=\"--max-gpu-memory ${MAX_GPU_MEMORY}GiB\"\nelse\n    MAX_GPU_MEMORY=\"\"\nfi\n\npython3 -m $MODEL_WORKER --host 0.0.0.0 --port $HTTP_PORT \\\n    --controller-address $CONTROLLER_ADDRESS \\\n    --worker-address http://$MY_POD_IP:$HTTP_PORT \\\n    --model-name $MODEL_NAME \\\n    --model-path $MODEL_PATH \\\n    $QUANTIZATION $NUM_GPUS $MAX_GPU_MEMORY $DEVICE_OPTION\n', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-1-8b', '/app/start.sh', '/data/ft-model/', 32768, 0, 1, 'inference'),
	(5, '2024-03-18 17:36:17.927', '2024-03-19 14:33:48.463', NULL, 'qwen1.5-1.8b-chat', 'qwen1.5-1.8b-chat', '#!/bin/bash\n\nMODEL_WORKER=fastchat.serve.model_worker\nCONTROLLER_ADDRESS=http://fschat-controller:21001\nMODEL_NAME={{.modelName}}\nMODEL_PATH={{.modelPath}}\nHTTP_PORT={{.port}}\nQUANTIZATION={{.quantization}}\nNUM_GPUS={{.numGpus}}\nMAX_GPU_MEMORY={{.maxGpuMemory}}\nVLLM={{.vllm}}\nINFERRED_TYPE={{.inferredType}}\nOS_TYPE=$(uname)\n\n# awq 量化配置\n# awq_wbits\n# awq_groupsize\n\n# 并发限制\n# limit-worker-concurrency\n\n# gptq 量化配置\n# gptq_wbits\n# gptq_groupsize\n# gptq_act_order\n\n# MODEL_WORKER \nif [ \"$VLLM\" == \"true\" ]; then\n    MODEL_WORKER=\"fastchat.serve.vllm_worker\"\nfi\n\n# 量化配置\nif [ \"$QUANTIZATION\" == \"8bit\" ]; then\n    QUANTIZATION=\"--load-8bit\"\nelse\n    QUANTIZATION=\"\"\nfi\n\n# NUM_GPUS\nif [ \"$NUM_GPUS\" -gt 0 ]; then\n    NUM_GPUS=\"--num-gpus $NUM_GPUS\"\nelse\n    NUM_GPUS=\"\"\nfi\n\n# CPU推理CPU，mps\nif [ \"$INFERRED_TYPE\" == \"cpu\" ] && [ \"$OS_TYPE\" == \"Darwin\" ]; then\n    DEVICE_OPTION=\"--device mps\"\nelif [ \"$INFERRED_TYPE\" == \"cpu\" ]; then\n    DEVICE_OPTION=\"--device cpu\"\nelse\n    DEVICE_OPTION=\"\"\nfi\n\n# MAX_GPU_MEMORY\nif [ \"$MAX_GPU_MEMORY\" -gt 0 ]; then\n    MAX_GPU_MEMORY=\"--max-gpu-memory ${MAX_GPU_MEMORY}GiB\"\nelse\n    MAX_GPU_MEMORY=\"\"\nfi\n\npython3 -m $MODEL_WORKER --host 0.0.0.0 --port $HTTP_PORT \\\n    --controller-address $CONTROLLER_ADDRESS \\\n    --worker-address http://$MY_POD_IP:$HTTP_PORT \\\n    --model-name $MODEL_NAME \\\n    --model-path $MODEL_PATH \\\n    $QUANTIZATION $NUM_GPUS $MAX_GPU_MEMORY $DEVICE_OPTION\n', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-1-8b-chat', '/app/start.sh', '/data/ft-model/', 32768, 0, 1, 'inference'),
	(6, '2024-03-18 17:36:17.927', '2024-03-19 14:33:38.874', NULL, 'qwen1.5-4b', 'qwen1.5-4b', '#!/bin/bash\n\nMODEL_WORKER=fastchat.serve.model_worker\nCONTROLLER_ADDRESS=http://fschat-controller:21001\nMODEL_NAME={{.modelName}}\nMODEL_PATH={{.modelPath}}\nHTTP_PORT={{.port}}\nQUANTIZATION={{.quantization}}\nNUM_GPUS={{.numGpus}}\nMAX_GPU_MEMORY={{.maxGpuMemory}}\nVLLM={{.vllm}}\nINFERRED_TYPE={{.inferredType}}\nOS_TYPE=$(uname)\n\n# awq 量化配置\n# awq_wbits\n# awq_groupsize\n\n# 并发限制\n# limit-worker-concurrency\n\n# gptq 量化配置\n# gptq_wbits\n# gptq_groupsize\n# gptq_act_order\n\n# MODEL_WORKER \nif [ \"$VLLM\" == \"true\" ]; then\n    MODEL_WORKER=\"fastchat.serve.vllm_worker\"\nfi\n\n# 量化配置\nif [ \"$QUANTIZATION\" == \"8bit\" ]; then\n    QUANTIZATION=\"--load-8bit\"\nelse\n    QUANTIZATION=\"\"\nfi\n\n# NUM_GPUS\nif [ \"$NUM_GPUS\" -gt 0 ]; then\n    NUM_GPUS=\"--num-gpus $NUM_GPUS\"\nelse\n    NUM_GPUS=\"\"\nfi\n\n# CPU推理CPU，mps\nif [ \"$INFERRED_TYPE\" == \"cpu\" ] && [ \"$OS_TYPE\" == \"Darwin\" ]; then\n    DEVICE_OPTION=\"--device mps\"\nelif [ \"$INFERRED_TYPE\" == \"cpu\" ]; then\n    DEVICE_OPTION=\"--device cpu\"\nelse\n    DEVICE_OPTION=\"\"\nfi\n\n# MAX_GPU_MEMORY\nif [ \"$MAX_GPU_MEMORY\" -gt 0 ]; then\n    MAX_GPU_MEMORY=\"--max-gpu-memory ${MAX_GPU_MEMORY}GiB\"\nelse\n    MAX_GPU_MEMORY=\"\"\nfi\n\npython3 -m $MODEL_WORKER --host 0.0.0.0 --port $HTTP_PORT \\\n    --controller-address $CONTROLLER_ADDRESS \\\n    --worker-address http://$MY_POD_IP:$HTTP_PORT \\\n    --model-name $MODEL_NAME \\\n    --model-path $MODEL_PATH \\\n    $QUANTIZATION $NUM_GPUS $MAX_GPU_MEMORY $DEVICE_OPTION\n', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-4b', '/app/start.sh', '/data/ft-model/', 32768, 0, 1, 'inference'),
	(7, '2024-03-18 17:36:17.927', '2024-03-19 14:34:18.717', NULL, 'qwen1.5-4b-chat', 'qwen1.5-4b-chat', '#!/bin/bash\n\nMODEL_WORKER=fastchat.serve.model_worker\nCONTROLLER_ADDRESS=http://fschat-controller:21001\nMODEL_NAME={{.modelName}}\nMODEL_PATH={{.modelPath}}\nHTTP_PORT={{.port}}\nQUANTIZATION={{.quantization}}\nNUM_GPUS={{.numGpus}}\nMAX_GPU_MEMORY={{.maxGpuMemory}}\nVLLM={{.vllm}}\nINFERRED_TYPE={{.inferredType}}\nOS_TYPE=$(uname)\n\n# awq 量化配置\n# awq_wbits\n# awq_groupsize\n\n# 并发限制\n# limit-worker-concurrency\n\n# gptq 量化配置\n# gptq_wbits\n# gptq_groupsize\n# gptq_act_order\n\n# MODEL_WORKER \nif [ \"$VLLM\" == \"true\" ]; then\n    MODEL_WORKER=\"fastchat.serve.vllm_worker\"\nfi\n\n# 量化配置\nif [ \"$QUANTIZATION\" == \"8bit\" ]; then\n    QUANTIZATION=\"--load-8bit\"\nelse\n    QUANTIZATION=\"\"\nfi\n\n# NUM_GPUS\nif [ \"$NUM_GPUS\" -gt 0 ]; then\n    NUM_GPUS=\"--num-gpus $NUM_GPUS\"\nelse\n    NUM_GPUS=\"\"\nfi\n\n# CPU推理CPU，mps\nif [ \"$INFERRED_TYPE\" == \"cpu\" ] && [ \"$OS_TYPE\" == \"Darwin\" ]; then\n    DEVICE_OPTION=\"--device mps\"\nelif [ \"$INFERRED_TYPE\" == \"cpu\" ]; then\n    DEVICE_OPTION=\"--device cpu\"\nelse\n    DEVICE_OPTION=\"\"\nfi\n\n# MAX_GPU_MEMORY\nif [ \"$MAX_GPU_MEMORY\" -gt 0 ]; then\n    MAX_GPU_MEMORY=\"--max-gpu-memory ${MAX_GPU_MEMORY}GiB\"\nelse\n    MAX_GPU_MEMORY=\"\"\nfi\n\npython3 -m $MODEL_WORKER --host 0.0.0.0 --port $HTTP_PORT \\\n    --controller-address $CONTROLLER_ADDRESS \\\n    --worker-address http://$MY_POD_IP:$HTTP_PORT \\\n    --model-name $MODEL_NAME \\\n    --model-path $MODEL_PATH \\\n    $QUANTIZATION $NUM_GPUS $MAX_GPU_MEMORY $DEVICE_OPTION\n', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-4b-chat', '/app/start.sh', '/data/ft-model/', 32768, 0, 1, 'inference'),
	(8, '2024-03-18 17:36:17.927', '2024-03-19 14:34:28.766', NULL, 'qwen1.5-7b', 'qwen1.5-7b', '#!/bin/bash\n\nMODEL_WORKER=fastchat.serve.model_worker\nCONTROLLER_ADDRESS=http://fschat-controller:21001\nMODEL_NAME={{.modelName}}\nMODEL_PATH={{.modelPath}}\nHTTP_PORT={{.port}}\nQUANTIZATION={{.quantization}}\nNUM_GPUS={{.numGpus}}\nMAX_GPU_MEMORY={{.maxGpuMemory}}\nVLLM={{.vllm}}\nINFERRED_TYPE={{.inferredType}}\nOS_TYPE=$(uname)\n\n# awq 量化配置\n# awq_wbits\n# awq_groupsize\n\n# 并发限制\n# limit-worker-concurrency\n\n# gptq 量化配置\n# gptq_wbits\n# gptq_groupsize\n# gptq_act_order\n\n# MODEL_WORKER \nif [ \"$VLLM\" == \"true\" ]; then\n    MODEL_WORKER=\"fastchat.serve.vllm_worker\"\nfi\n\n# 量化配置\nif [ \"$QUANTIZATION\" == \"8bit\" ]; then\n    QUANTIZATION=\"--load-8bit\"\nelse\n    QUANTIZATION=\"\"\nfi\n\n# NUM_GPUS\nif [ \"$NUM_GPUS\" -gt 0 ]; then\n    NUM_GPUS=\"--num-gpus $NUM_GPUS\"\nelse\n    NUM_GPUS=\"\"\nfi\n\n# CPU推理CPU，mps\nif [ \"$INFERRED_TYPE\" == \"cpu\" ] && [ \"$OS_TYPE\" == \"Darwin\" ]; then\n    DEVICE_OPTION=\"--device mps\"\nelif [ \"$INFERRED_TYPE\" == \"cpu\" ]; then\n    DEVICE_OPTION=\"--device cpu\"\nelse\n    DEVICE_OPTION=\"\"\nfi\n\n# MAX_GPU_MEMORY\nif [ \"$MAX_GPU_MEMORY\" -gt 0 ]; then\n    MAX_GPU_MEMORY=\"--max-gpu-memory ${MAX_GPU_MEMORY}GiB\"\nelse\n    MAX_GPU_MEMORY=\"\"\nfi\n\npython3 -m $MODEL_WORKER --host 0.0.0.0 --port $HTTP_PORT \\\n    --controller-address $CONTROLLER_ADDRESS \\\n    --worker-address http://$MY_POD_IP:$HTTP_PORT \\\n    --model-name $MODEL_NAME \\\n    --model-path $MODEL_PATH \\\n    $QUANTIZATION $NUM_GPUS $MAX_GPU_MEMORY $DEVICE_OPTION\n', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-7b', '/app/start.sh', '/data/ft-model/', 32768, 0, 1, 'inference'),
	(9, '2024-03-19 17:36:17.927', '2024-03-19 14:32:46.468', NULL, 'qwen1.5-7b-chat', 'qwen1.5-7b-chat', '#!/bin/bash\n\nMODEL_WORKER=fastchat.serve.model_worker\nCONTROLLER_ADDRESS=http://fschat-controller:21001\nMODEL_NAME={{.modelName}}\nMODEL_PATH={{.modelPath}}\nHTTP_PORT={{.port}}\nQUANTIZATION={{.quantization}}\nNUM_GPUS={{.numGpus}}\nMAX_GPU_MEMORY={{.maxGpuMemory}}\nVLLM={{.vllm}}\nINFERRED_TYPE={{.inferredType}}\nOS_TYPE=$(uname)\n\n# awq 量化配置\n# awq_wbits\n# awq_groupsize\n\n# 并发限制\n# limit-worker-concurrency\n\n# gptq 量化配置\n# gptq_wbits\n# gptq_groupsize\n# gptq_act_order\n\n# MODEL_WORKER \nif [ \"$VLLM\" == \"true\" ]; then\n    MODEL_WORKER=\"fastchat.serve.vllm_worker\"\nfi\n\n# 量化配置\nif [ \"$QUANTIZATION\" == \"8bit\" ]; then\n    QUANTIZATION=\"--load-8bit\"\nelse\n    QUANTIZATION=\"\"\nfi\n\n# NUM_GPUS\nif [ \"$NUM_GPUS\" -gt 0 ]; then\n    NUM_GPUS=\"--num-gpus $NUM_GPUS\"\nelse\n    NUM_GPUS=\"\"\nfi\n\n# CPU推理CPU，mps\nif [ \"$INFERRED_TYPE\" == \"cpu\" ] && [ \"$OS_TYPE\" == \"Darwin\" ]; then\n    DEVICE_OPTION=\"--device mps\"\nelif [ \"$INFERRED_TYPE\" == \"cpu\" ]; then\n    DEVICE_OPTION=\"--device cpu\"\nelse\n    DEVICE_OPTION=\"\"\nfi\n\n# MAX_GPU_MEMORY\nif [ \"$MAX_GPU_MEMORY\" -gt 0 ]; then\n    MAX_GPU_MEMORY=\"--max-gpu-memory ${MAX_GPU_MEMORY}GiB\"\nelse\n    MAX_GPU_MEMORY=\"\"\nfi\n\npython3 -m $MODEL_WORKER --host 0.0.0.0 --port $HTTP_PORT \\\n    --controller-address $CONTROLLER_ADDRESS \\\n    --worker-address http://$MY_POD_IP:$HTTP_PORT \\\n    --model-name $MODEL_NAME \\\n    --model-path $MODEL_PATH \\\n    $QUANTIZATION $NUM_GPUS $MAX_GPU_MEMORY $DEVICE_OPTION\n', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-0-5b', '/app/start.sh', '/data/ft-model/', 32768, 0, 1, 'inference'),
	(10, '2024-03-19 17:36:17.927', '2024-03-19 14:34:40.884', NULL, 'qwen1.5-14b', 'qwen1.5-14b', '#!/bin/bash\n\nMODEL_WORKER=fastchat.serve.model_worker\nCONTROLLER_ADDRESS=http://fschat-controller:21001\nMODEL_NAME={{.modelName}}\nMODEL_PATH={{.modelPath}}\nHTTP_PORT={{.port}}\nQUANTIZATION={{.quantization}}\nNUM_GPUS={{.numGpus}}\nMAX_GPU_MEMORY={{.maxGpuMemory}}\nVLLM={{.vllm}}\nINFERRED_TYPE={{.inferredType}}\nOS_TYPE=$(uname)\n\n# awq 量化配置\n# awq_wbits\n# awq_groupsize\n\n# 并发限制\n# limit-worker-concurrency\n\n# gptq 量化配置\n# gptq_wbits\n# gptq_groupsize\n# gptq_act_order\n\n# MODEL_WORKER \nif [ \"$VLLM\" == \"true\" ]; then\n    MODEL_WORKER=\"fastchat.serve.vllm_worker\"\nfi\n\n# 量化配置\nif [ \"$QUANTIZATION\" == \"8bit\" ]; then\n    QUANTIZATION=\"--load-8bit\"\nelse\n    QUANTIZATION=\"\"\nfi\n\n# NUM_GPUS\nif [ \"$NUM_GPUS\" -gt 0 ]; then\n    NUM_GPUS=\"--num-gpus $NUM_GPUS\"\nelse\n    NUM_GPUS=\"\"\nfi\n\n# CPU推理CPU，mps\nif [ \"$INFERRED_TYPE\" == \"cpu\" ] && [ \"$OS_TYPE\" == \"Darwin\" ]; then\n    DEVICE_OPTION=\"--device mps\"\nelif [ \"$INFERRED_TYPE\" == \"cpu\" ]; then\n    DEVICE_OPTION=\"--device cpu\"\nelse\n    DEVICE_OPTION=\"\"\nfi\n\n# MAX_GPU_MEMORY\nif [ \"$MAX_GPU_MEMORY\" -gt 0 ]; then\n    MAX_GPU_MEMORY=\"--max-gpu-memory ${MAX_GPU_MEMORY}GiB\"\nelse\n    MAX_GPU_MEMORY=\"\"\nfi\n\npython3 -m $MODEL_WORKER --host 0.0.0.0 --port $HTTP_PORT \\\n    --controller-address $CONTROLLER_ADDRESS \\\n    --worker-address http://$MY_POD_IP:$HTTP_PORT \\\n    --model-name $MODEL_NAME \\\n    --model-path $MODEL_PATH \\\n    $QUANTIZATION $NUM_GPUS $MAX_GPU_MEMORY $DEVICE_OPTION\n', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-14b', '/app/start.sh', '/data/ft-model/', 32768, 0, 1, 'inference'),
	(11, '2024-03-19 17:36:17.927', '2024-03-19 14:33:01.470', NULL, 'qwen1.5-14b-chat', 'qwen1.5-14b-chat', '#!/bin/bash\n\nMODEL_WORKER=fastchat.serve.model_worker\nCONTROLLER_ADDRESS=http://fschat-controller:21001\nMODEL_NAME={{.modelName}}\nMODEL_PATH={{.modelPath}}\nHTTP_PORT={{.port}}\nQUANTIZATION={{.quantization}}\nNUM_GPUS={{.numGpus}}\nMAX_GPU_MEMORY={{.maxGpuMemory}}\nVLLM={{.vllm}}\nINFERRED_TYPE={{.inferredType}}\nOS_TYPE=$(uname)\n\n# awq 量化配置\n# awq_wbits\n# awq_groupsize\n\n# 并发限制\n# limit-worker-concurrency\n\n# gptq 量化配置\n# gptq_wbits\n# gptq_groupsize\n# gptq_act_order\n\n# MODEL_WORKER \nif [ \"$VLLM\" == \"true\" ]; then\n    MODEL_WORKER=\"fastchat.serve.vllm_worker\"\nfi\n\n# 量化配置\nif [ \"$QUANTIZATION\" == \"8bit\" ]; then\n    QUANTIZATION=\"--load-8bit\"\nelse\n    QUANTIZATION=\"\"\nfi\n\n# NUM_GPUS\nif [ \"$NUM_GPUS\" -gt 0 ]; then\n    NUM_GPUS=\"--num-gpus $NUM_GPUS\"\nelse\n    NUM_GPUS=\"\"\nfi\n\n# CPU推理CPU，mps\nif [ \"$INFERRED_TYPE\" == \"cpu\" ] && [ \"$OS_TYPE\" == \"Darwin\" ]; then\n    DEVICE_OPTION=\"--device mps\"\nelif [ \"$INFERRED_TYPE\" == \"cpu\" ]; then\n    DEVICE_OPTION=\"--device cpu\"\nelse\n    DEVICE_OPTION=\"\"\nfi\n\n# MAX_GPU_MEMORY\nif [ \"$MAX_GPU_MEMORY\" -gt 0 ]; then\n    MAX_GPU_MEMORY=\"--max-gpu-memory ${MAX_GPU_MEMORY}GiB\"\nelse\n    MAX_GPU_MEMORY=\"\"\nfi\n\npython3 -m $MODEL_WORKER --host 0.0.0.0 --port $HTTP_PORT \\\n    --controller-address $CONTROLLER_ADDRESS \\\n    --worker-address http://$MY_POD_IP:$HTTP_PORT \\\n    --model-name $MODEL_NAME \\\n    --model-path $MODEL_PATH \\\n    $QUANTIZATION $NUM_GPUS $MAX_GPU_MEMORY $DEVICE_OPTION\n', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-0-5b', '/app/start.sh', '/data/ft-model/', 32768, 0, 1, 'inference'),
	(12, '2024-03-19 17:36:17.927', '2024-03-19 14:33:16.820', NULL, 'qwen1.5-72b', 'qwen1.5-72b', '#!/bin/bash\n\nMODEL_WORKER=fastchat.serve.model_worker\nCONTROLLER_ADDRESS=http://fschat-controller:21001\nMODEL_NAME={{.modelName}}\nMODEL_PATH={{.modelPath}}\nHTTP_PORT={{.port}}\nQUANTIZATION={{.quantization}}\nNUM_GPUS={{.numGpus}}\nMAX_GPU_MEMORY={{.maxGpuMemory}}\nVLLM={{.vllm}}\nINFERRED_TYPE={{.inferredType}}\nOS_TYPE=$(uname)\n\n# awq 量化配置\n# awq_wbits\n# awq_groupsize\n\n# 并发限制\n# limit-worker-concurrency\n\n# gptq 量化配置\n# gptq_wbits\n# gptq_groupsize\n# gptq_act_order\n\n# MODEL_WORKER \nif [ \"$VLLM\" == \"true\" ]; then\n    MODEL_WORKER=\"fastchat.serve.vllm_worker\"\nfi\n\n# 量化配置\nif [ \"$QUANTIZATION\" == \"8bit\" ]; then\n    QUANTIZATION=\"--load-8bit\"\nelse\n    QUANTIZATION=\"\"\nfi\n\n# NUM_GPUS\nif [ \"$NUM_GPUS\" -gt 0 ]; then\n    NUM_GPUS=\"--num-gpus $NUM_GPUS\"\nelse\n    NUM_GPUS=\"\"\nfi\n\n# CPU推理CPU，mps\nif [ \"$INFERRED_TYPE\" == \"cpu\" ] && [ \"$OS_TYPE\" == \"Darwin\" ]; then\n    DEVICE_OPTION=\"--device mps\"\nelif [ \"$INFERRED_TYPE\" == \"cpu\" ]; then\n    DEVICE_OPTION=\"--device cpu\"\nelse\n    DEVICE_OPTION=\"\"\nfi\n\n# MAX_GPU_MEMORY\nif [ \"$MAX_GPU_MEMORY\" -gt 0 ]; then\n    MAX_GPU_MEMORY=\"--max-gpu-memory ${MAX_GPU_MEMORY}GiB\"\nelse\n    MAX_GPU_MEMORY=\"\"\nfi\n\npython3 -m $MODEL_WORKER --host 0.0.0.0 --port $HTTP_PORT \\\n    --controller-address $CONTROLLER_ADDRESS \\\n    --worker-address http://$MY_POD_IP:$HTTP_PORT \\\n    --model-name $MODEL_NAME \\\n    --model-path $MODEL_PATH \\\n    $QUANTIZATION $NUM_GPUS $MAX_GPU_MEMORY $DEVICE_OPTION\n', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-72b', '/app/start.sh', '/data/ft-model/', 32768, 0, 1, 'inference'),
	(13, '2024-03-19 17:36:17.927', '2024-03-19 14:33:26.291', NULL, 'qwen1.5-72b-chat', 'qwen1.5-72b-chat', '#!/bin/bash\n\nMODEL_WORKER=fastchat.serve.model_worker\nCONTROLLER_ADDRESS=http://fschat-controller:21001\nMODEL_NAME={{.modelName}}\nMODEL_PATH={{.modelPath}}\nHTTP_PORT={{.port}}\nQUANTIZATION={{.quantization}}\nNUM_GPUS={{.numGpus}}\nMAX_GPU_MEMORY={{.maxGpuMemory}}\nVLLM={{.vllm}}\nINFERRED_TYPE={{.inferredType}}\nOS_TYPE=$(uname)\n\n# awq 量化配置\n# awq_wbits\n# awq_groupsize\n\n# 并发限制\n# limit-worker-concurrency\n\n# gptq 量化配置\n# gptq_wbits\n# gptq_groupsize\n# gptq_act_order\n\n# MODEL_WORKER \nif [ \"$VLLM\" == \"true\" ]; then\n    MODEL_WORKER=\"fastchat.serve.vllm_worker\"\nfi\n\n# 量化配置\nif [ \"$QUANTIZATION\" == \"8bit\" ]; then\n    QUANTIZATION=\"--load-8bit\"\nelse\n    QUANTIZATION=\"\"\nfi\n\n# NUM_GPUS\nif [ \"$NUM_GPUS\" -gt 0 ]; then\n    NUM_GPUS=\"--num-gpus $NUM_GPUS\"\nelse\n    NUM_GPUS=\"\"\nfi\n\n# CPU推理CPU，mps\nif [ \"$INFERRED_TYPE\" == \"cpu\" ] && [ \"$OS_TYPE\" == \"Darwin\" ]; then\n    DEVICE_OPTION=\"--device mps\"\nelif [ \"$INFERRED_TYPE\" == \"cpu\" ]; then\n    DEVICE_OPTION=\"--device cpu\"\nelse\n    DEVICE_OPTION=\"\"\nfi\n\n# MAX_GPU_MEMORY\nif [ \"$MAX_GPU_MEMORY\" -gt 0 ]; then\n    MAX_GPU_MEMORY=\"--max-gpu-memory ${MAX_GPU_MEMORY}GiB\"\nelse\n    MAX_GPU_MEMORY=\"\"\nfi\n\npython3 -m $MODEL_WORKER --host 0.0.0.0 --port $HTTP_PORT \\\n    --controller-address $CONTROLLER_ADDRESS \\\n    --worker-address http://$MY_POD_IP:$HTTP_PORT \\\n    --model-name $MODEL_NAME \\\n    --model-path $MODEL_PATH \\\n    $QUANTIZATION $NUM_GPUS $MAX_GPU_MEMORY $DEVICE_OPTION\n', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-72b-chat', '/app/start.sh', '/data/ft-model/', 32768, 0, 1, 'inference'),
	(14, '2024-03-19 14:13:04.852', '2024-03-19 14:38:34.152', NULL, 'qwen1.5-1.8b-train', 'qwen1.5-1.8b', '#!/bin/bash\nexport AUTH=sk-001\nexport JOB_ID={{.JobId}}\n\nexport CUDA_DEVICE_MAX_CONNECTIONS=1\n\n\nGPUS_PER_NODE={{.ProcPerNode}}\nNNODES=1\nNODE_RANK=0\nMASTER_ADDR=localhost\nMASTER_PORT={{.MasterPort}}\nUSE_LORA={{.Lora}}\nQ_LORA=False\n\nMODEL=\"{{.BaseModelPath}}\" # Set the path if you do not want to load from huggingface directly\n# ATTENTION: specify the path to your training data, which should be a json file consisting of a list of conversations.\n# See the section for finetuning in README for more information.\nDATA=\"{{.DataPath}}\"\n# 验证集\nEVAL_DATA=\"{{.ValidationFile}}\"\nDS_CONFIG_PATH=\"ds_config_zero3.json\"\n\nDISTRIBUTED_ARGS=\"\n    --nproc_per_node $GPUS_PER_NODE \\\n    --nnodes $NNODES \\\n    --node_rank $NODE_RANK \\\n    --master_addr $MASTER_ADDR \\\n    --master_port $MASTER_PORT\n\"\nif [ \"$USE_LORA\" == \"true\" ]; then\n    USE_LORA=True\n    DS_CONFIG_PATH=\"ds_config_zero2.json\"\nelse\n    USE_LORA=False\n    DS_CONFIG_PATH=\"ds_config_zero3.json\"\nfi\n\nmkdir -p /data/train-data/\nwget -O {{.DataPath}} {{.FileUrl}}\n\ntorchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \\\n    --model_name_or_path $MODEL \\\n    --data_path $DATA \\\n    --bf16 True \\\n    --output_dir {{.OutputDir}} \\\n    --num_train_epochs {{.TrainEpoch}} \\\n    --per_device_train_batch_size {{.TrainBatchSize}} \\\n    --per_device_eval_batch_size {{.EvalBatchSize}} \\\n    --gradient_accumulation_steps {{.AccumulationSteps}} \\\n    --evaluation_strategy \"no\" \\\n    --save_strategy \"steps\" \\\n    --save_steps 1000 \\\n    --save_total_limit 10 \\\n    --learning_rate {{.LearningRate}} \\\n    --weight_decay 0.1 \\\n    --adam_beta2 0.95 \\\n    --warmup_ratio 0.01 \\\n    --lr_scheduler_type \"cosine\" \\\n    --logging_steps 1 \\\n    --report_to \"none\" \\\n    --model_max_length {{.ModelMaxLength}} \\\n    --gradient_checkpointing True \\\n    --lazy_preprocess True \\\n    --use_lora ${USE_LORA} \\\n    --q_lora ${Q_LORA} \\\n    --deepspeed $DS_CONFIG_PATH', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-1-8b', '/app/finetune.py', '/data/ft-model/', 32768, 0, 1, 'train'),
	(15, '2024-03-19 14:13:04.852', '2024-03-19 14:38:40.229', NULL, 'qwen1.5-1.8b-chat-train', 'qwen1.5-1.8b-chat', '#!/bin/bash\nexport AUTH=sk-001\nexport JOB_ID={{.JobId}}\n\nexport CUDA_DEVICE_MAX_CONNECTIONS=1\n\n\nGPUS_PER_NODE={{.ProcPerNode}}\nNNODES=1\nNODE_RANK=0\nMASTER_ADDR=localhost\nMASTER_PORT={{.MasterPort}}\nUSE_LORA={{.Lora}}\nQ_LORA=False\n\nMODEL=\"{{.BaseModelPath}}\" # Set the path if you do not want to load from huggingface directly\n# ATTENTION: specify the path to your training data, which should be a json file consisting of a list of conversations.\n# See the section for finetuning in README for more information.\nDATA=\"{{.DataPath}}\"\n# 验证集\nEVAL_DATA=\"{{.ValidationFile}}\"\nDS_CONFIG_PATH=\"ds_config_zero3.json\"\n\nDISTRIBUTED_ARGS=\"\n    --nproc_per_node $GPUS_PER_NODE \\\n    --nnodes $NNODES \\\n    --node_rank $NODE_RANK \\\n    --master_addr $MASTER_ADDR \\\n    --master_port $MASTER_PORT\n\"\nif [ \"$USE_LORA\" == \"true\" ]; then\n    USE_LORA=True\n    DS_CONFIG_PATH=\"ds_config_zero2.json\"\nelse\n    USE_LORA=False\n    DS_CONFIG_PATH=\"ds_config_zero3.json\"\nfi\n\nmkdir -p /data/train-data/\nwget -O {{.DataPath}} {{.FileUrl}}\n\ntorchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \\\n    --model_name_or_path $MODEL \\\n    --data_path $DATA \\\n    --bf16 True \\\n    --output_dir {{.OutputDir}} \\\n    --num_train_epochs {{.TrainEpoch}} \\\n    --per_device_train_batch_size {{.TrainBatchSize}} \\\n    --per_device_eval_batch_size {{.EvalBatchSize}} \\\n    --gradient_accumulation_steps {{.AccumulationSteps}} \\\n    --evaluation_strategy \"no\" \\\n    --save_strategy \"steps\" \\\n    --save_steps 1000 \\\n    --save_total_limit 10 \\\n    --learning_rate {{.LearningRate}} \\\n    --weight_decay 0.1 \\\n    --adam_beta2 0.95 \\\n    --warmup_ratio 0.01 \\\n    --lr_scheduler_type \"cosine\" \\\n    --logging_steps 1 \\\n    --report_to \"none\" \\\n    --model_max_length {{.ModelMaxLength}} \\\n    --gradient_checkpointing True \\\n    --lazy_preprocess True \\\n    --use_lora ${USE_LORA} \\\n    --q_lora ${Q_LORA} \\\n    --deepspeed $DS_CONFIG_PATH', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-1-8b-chat', '/app/finetune.py', '/data/ft-model/', 32768, 0, 1, 'train'),
	(16, '2024-03-19 14:13:04.852', '2024-03-19 14:39:11.505', NULL, 'qwen1.5-4b-chat-train', 'qwen1.5-4b-chat', '#!/bin/bash\nexport AUTH=sk-001\nexport JOB_ID={{.JobId}}\n\nexport CUDA_DEVICE_MAX_CONNECTIONS=1\n\n\nGPUS_PER_NODE={{.ProcPerNode}}\nNNODES=1\nNODE_RANK=0\nMASTER_ADDR=localhost\nMASTER_PORT={{.MasterPort}}\nUSE_LORA={{.Lora}}\nQ_LORA=False\n\nMODEL=\"{{.BaseModelPath}}\" # Set the path if you do not want to load from huggingface directly\n# ATTENTION: specify the path to your training data, which should be a json file consisting of a list of conversations.\n# See the section for finetuning in README for more information.\nDATA=\"{{.DataPath}}\"\n# 验证集\nEVAL_DATA=\"{{.ValidationFile}}\"\nDS_CONFIG_PATH=\"ds_config_zero3.json\"\n\nDISTRIBUTED_ARGS=\"\n    --nproc_per_node $GPUS_PER_NODE \\\n    --nnodes $NNODES \\\n    --node_rank $NODE_RANK \\\n    --master_addr $MASTER_ADDR \\\n    --master_port $MASTER_PORT\n\"\nif [ \"$USE_LORA\" == \"true\" ]; then\n    USE_LORA=True\n    DS_CONFIG_PATH=\"ds_config_zero2.json\"\nelse\n    USE_LORA=False\n    DS_CONFIG_PATH=\"ds_config_zero3.json\"\nfi\n\nmkdir -p /data/train-data/\nwget -O {{.DataPath}} {{.FileUrl}}\n\ntorchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \\\n    --model_name_or_path $MODEL \\\n    --data_path $DATA \\\n    --bf16 True \\\n    --output_dir {{.OutputDir}} \\\n    --num_train_epochs {{.TrainEpoch}} \\\n    --per_device_train_batch_size {{.TrainBatchSize}} \\\n    --per_device_eval_batch_size {{.EvalBatchSize}} \\\n    --gradient_accumulation_steps {{.AccumulationSteps}} \\\n    --evaluation_strategy \"no\" \\\n    --save_strategy \"steps\" \\\n    --save_steps 1000 \\\n    --save_total_limit 10 \\\n    --learning_rate {{.LearningRate}} \\\n    --weight_decay 0.1 \\\n    --adam_beta2 0.95 \\\n    --warmup_ratio 0.01 \\\n    --lr_scheduler_type \"cosine\" \\\n    --logging_steps 1 \\\n    --report_to \"none\" \\\n    --model_max_length {{.ModelMaxLength}} \\\n    --gradient_checkpointing True \\\n    --lazy_preprocess True \\\n    --use_lora ${USE_LORA} \\\n    --q_lora ${Q_LORA} \\\n    --deepspeed $DS_CONFIG_PATH', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-4b-chat', '/app/finetune.py', '/data/ft-model/', 32768, 0, 1, 'train'),
	(17, '2024-03-19 14:13:04.852', '2024-03-19 14:39:03.600', NULL, 'qwen1.5-4b-train', 'qwen1.5-4b', '#!/bin/bash\nexport AUTH=sk-001\nexport JOB_ID={{.JobId}}\n\nexport CUDA_DEVICE_MAX_CONNECTIONS=1\n\n\nGPUS_PER_NODE={{.ProcPerNode}}\nNNODES=1\nNODE_RANK=0\nMASTER_ADDR=localhost\nMASTER_PORT={{.MasterPort}}\nUSE_LORA={{.Lora}}\nQ_LORA=False\n\nMODEL=\"{{.BaseModelPath}}\" # Set the path if you do not want to load from huggingface directly\n# ATTENTION: specify the path to your training data, which should be a json file consisting of a list of conversations.\n# See the section for finetuning in README for more information.\nDATA=\"{{.DataPath}}\"\n# 验证集\nEVAL_DATA=\"{{.ValidationFile}}\"\nDS_CONFIG_PATH=\"ds_config_zero3.json\"\n\nDISTRIBUTED_ARGS=\"\n    --nproc_per_node $GPUS_PER_NODE \\\n    --nnodes $NNODES \\\n    --node_rank $NODE_RANK \\\n    --master_addr $MASTER_ADDR \\\n    --master_port $MASTER_PORT\n\"\nif [ \"$USE_LORA\" == \"true\" ]; then\n    USE_LORA=True\n    DS_CONFIG_PATH=\"ds_config_zero2.json\"\nelse\n    USE_LORA=False\n    DS_CONFIG_PATH=\"ds_config_zero3.json\"\nfi\n\nmkdir -p /data/train-data/\nwget -O {{.DataPath}} {{.FileUrl}}\n\ntorchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \\\n    --model_name_or_path $MODEL \\\n    --data_path $DATA \\\n    --bf16 True \\\n    --output_dir {{.OutputDir}} \\\n    --num_train_epochs {{.TrainEpoch}} \\\n    --per_device_train_batch_size {{.TrainBatchSize}} \\\n    --per_device_eval_batch_size {{.EvalBatchSize}} \\\n    --gradient_accumulation_steps {{.AccumulationSteps}} \\\n    --evaluation_strategy \"no\" \\\n    --save_strategy \"steps\" \\\n    --save_steps 1000 \\\n    --save_total_limit 10 \\\n    --learning_rate {{.LearningRate}} \\\n    --weight_decay 0.1 \\\n    --adam_beta2 0.95 \\\n    --warmup_ratio 0.01 \\\n    --lr_scheduler_type \"cosine\" \\\n    --logging_steps 1 \\\n    --report_to \"none\" \\\n    --model_max_length {{.ModelMaxLength}} \\\n    --gradient_checkpointing True \\\n    --lazy_preprocess True \\\n    --use_lora ${USE_LORA} \\\n    --q_lora ${Q_LORA} \\\n    --deepspeed $DS_CONFIG_PATH', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-4b', '/app/finetune.py', '/data/ft-model/', 32768, 0, 1, 'train'),
	(18, '2024-03-19 14:13:04.852', '2024-03-19 14:39:24.104', NULL, 'qwen1.5-7b-train', 'qwen1.5-7b', '#!/bin/bash\nexport AUTH=sk-001\nexport JOB_ID={{.JobId}}\n\nexport CUDA_DEVICE_MAX_CONNECTIONS=1\n\n\nGPUS_PER_NODE={{.ProcPerNode}}\nNNODES=1\nNODE_RANK=0\nMASTER_ADDR=localhost\nMASTER_PORT={{.MasterPort}}\nUSE_LORA={{.Lora}}\nQ_LORA=False\n\nMODEL=\"{{.BaseModelPath}}\" # Set the path if you do not want to load from huggingface directly\n# ATTENTION: specify the path to your training data, which should be a json file consisting of a list of conversations.\n# See the section for finetuning in README for more information.\nDATA=\"{{.DataPath}}\"\n# 验证集\nEVAL_DATA=\"{{.ValidationFile}}\"\nDS_CONFIG_PATH=\"ds_config_zero3.json\"\n\nDISTRIBUTED_ARGS=\"\n    --nproc_per_node $GPUS_PER_NODE \\\n    --nnodes $NNODES \\\n    --node_rank $NODE_RANK \\\n    --master_addr $MASTER_ADDR \\\n    --master_port $MASTER_PORT\n\"\nif [ \"$USE_LORA\" == \"true\" ]; then\n    USE_LORA=True\n    DS_CONFIG_PATH=\"ds_config_zero2.json\"\nelse\n    USE_LORA=False\n    DS_CONFIG_PATH=\"ds_config_zero3.json\"\nfi\n\nmkdir -p /data/train-data/\nwget -O {{.DataPath}} {{.FileUrl}}\n\ntorchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \\\n    --model_name_or_path $MODEL \\\n    --data_path $DATA \\\n    --bf16 True \\\n    --output_dir {{.OutputDir}} \\\n    --num_train_epochs {{.TrainEpoch}} \\\n    --per_device_train_batch_size {{.TrainBatchSize}} \\\n    --per_device_eval_batch_size {{.EvalBatchSize}} \\\n    --gradient_accumulation_steps {{.AccumulationSteps}} \\\n    --evaluation_strategy \"no\" \\\n    --save_strategy \"steps\" \\\n    --save_steps 1000 \\\n    --save_total_limit 10 \\\n    --learning_rate {{.LearningRate}} \\\n    --weight_decay 0.1 \\\n    --adam_beta2 0.95 \\\n    --warmup_ratio 0.01 \\\n    --lr_scheduler_type \"cosine\" \\\n    --logging_steps 1 \\\n    --report_to \"none\" \\\n    --model_max_length {{.ModelMaxLength}} \\\n    --gradient_checkpointing True \\\n    --lazy_preprocess True \\\n    --use_lora ${USE_LORA} \\\n    --q_lora ${Q_LORA} \\\n    --deepspeed $DS_CONFIG_PATH', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model', '/app/finetune.py', '/data/base-model/qwen1-5-7b', 32768, 0, 1, 'train'),
	(19, '2024-03-19 14:13:04.852', '2024-03-19 14:39:34.841', NULL, 'qwen1.5-7b-chat-train', 'qwen1.5-7b-chat', '#!/bin/bash\nexport AUTH=sk-001\nexport JOB_ID={{.JobId}}\n\nexport CUDA_DEVICE_MAX_CONNECTIONS=1\n\n\nGPUS_PER_NODE={{.ProcPerNode}}\nNNODES=1\nNODE_RANK=0\nMASTER_ADDR=localhost\nMASTER_PORT={{.MasterPort}}\nUSE_LORA={{.Lora}}\nQ_LORA=False\n\nMODEL=\"{{.BaseModelPath}}\" # Set the path if you do not want to load from huggingface directly\n# ATTENTION: specify the path to your training data, which should be a json file consisting of a list of conversations.\n# See the section for finetuning in README for more information.\nDATA=\"{{.DataPath}}\"\n# 验证集\nEVAL_DATA=\"{{.ValidationFile}}\"\nDS_CONFIG_PATH=\"ds_config_zero3.json\"\n\nDISTRIBUTED_ARGS=\"\n    --nproc_per_node $GPUS_PER_NODE \\\n    --nnodes $NNODES \\\n    --node_rank $NODE_RANK \\\n    --master_addr $MASTER_ADDR \\\n    --master_port $MASTER_PORT\n\"\nif [ \"$USE_LORA\" == \"true\" ]; then\n    USE_LORA=True\n    DS_CONFIG_PATH=\"ds_config_zero2.json\"\nelse\n    USE_LORA=False\n    DS_CONFIG_PATH=\"ds_config_zero3.json\"\nfi\n\nmkdir -p /data/train-data/\nwget -O {{.DataPath}} {{.FileUrl}}\n\ntorchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \\\n    --model_name_or_path $MODEL \\\n    --data_path $DATA \\\n    --bf16 True \\\n    --output_dir {{.OutputDir}} \\\n    --num_train_epochs {{.TrainEpoch}} \\\n    --per_device_train_batch_size {{.TrainBatchSize}} \\\n    --per_device_eval_batch_size {{.EvalBatchSize}} \\\n    --gradient_accumulation_steps {{.AccumulationSteps}} \\\n    --evaluation_strategy \"no\" \\\n    --save_strategy \"steps\" \\\n    --save_steps 1000 \\\n    --save_total_limit 10 \\\n    --learning_rate {{.LearningRate}} \\\n    --weight_decay 0.1 \\\n    --adam_beta2 0.95 \\\n    --warmup_ratio 0.01 \\\n    --lr_scheduler_type \"cosine\" \\\n    --logging_steps 1 \\\n    --report_to \"none\" \\\n    --model_max_length {{.ModelMaxLength}} \\\n    --gradient_checkpointing True \\\n    --lazy_preprocess True \\\n    --use_lora ${USE_LORA} \\\n    --q_lora ${Q_LORA} \\\n    --deepspeed $DS_CONFIG_PATH', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-7b-chat', '/app/finetune.py', '/data/ft-model/', 32768, 0, 1, 'train'),
	(20, '2024-03-19 14:13:04.852', '2024-03-19 14:39:50.220', NULL, 'qwen1.5-14b-chat-train', 'qwen1.5-14b-chat', '#!/bin/bash\nexport AUTH=sk-001\nexport JOB_ID={{.JobId}}\n\nexport CUDA_DEVICE_MAX_CONNECTIONS=1\n\n\nGPUS_PER_NODE={{.ProcPerNode}}\nNNODES=1\nNODE_RANK=0\nMASTER_ADDR=localhost\nMASTER_PORT={{.MasterPort}}\nUSE_LORA={{.Lora}}\nQ_LORA=False\n\nMODEL=\"{{.BaseModelPath}}\" # Set the path if you do not want to load from huggingface directly\n# ATTENTION: specify the path to your training data, which should be a json file consisting of a list of conversations.\n# See the section for finetuning in README for more information.\nDATA=\"{{.DataPath}}\"\n# 验证集\nEVAL_DATA=\"{{.ValidationFile}}\"\nDS_CONFIG_PATH=\"ds_config_zero3.json\"\n\nDISTRIBUTED_ARGS=\"\n    --nproc_per_node $GPUS_PER_NODE \\\n    --nnodes $NNODES \\\n    --node_rank $NODE_RANK \\\n    --master_addr $MASTER_ADDR \\\n    --master_port $MASTER_PORT\n\"\nif [ \"$USE_LORA\" == \"true\" ]; then\n    USE_LORA=True\n    DS_CONFIG_PATH=\"ds_config_zero2.json\"\nelse\n    USE_LORA=False\n    DS_CONFIG_PATH=\"ds_config_zero3.json\"\nfi\n\nmkdir -p /data/train-data/\nwget -O {{.DataPath}} {{.FileUrl}}\n\ntorchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \\\n    --model_name_or_path $MODEL \\\n    --data_path $DATA \\\n    --bf16 True \\\n    --output_dir {{.OutputDir}} \\\n    --num_train_epochs {{.TrainEpoch}} \\\n    --per_device_train_batch_size {{.TrainBatchSize}} \\\n    --per_device_eval_batch_size {{.EvalBatchSize}} \\\n    --gradient_accumulation_steps {{.AccumulationSteps}} \\\n    --evaluation_strategy \"no\" \\\n    --save_strategy \"steps\" \\\n    --save_steps 1000 \\\n    --save_total_limit 10 \\\n    --learning_rate {{.LearningRate}} \\\n    --weight_decay 0.1 \\\n    --adam_beta2 0.95 \\\n    --warmup_ratio 0.01 \\\n    --lr_scheduler_type \"cosine\" \\\n    --logging_steps 1 \\\n    --report_to \"none\" \\\n    --model_max_length {{.ModelMaxLength}} \\\n    --gradient_checkpointing True \\\n    --lazy_preprocess True \\\n    --use_lora ${USE_LORA} \\\n    --q_lora ${Q_LORA} \\\n    --deepspeed $DS_CONFIG_PATH', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-14b-chat', '/app/finetune.py', '/data/ft-model/', 32768, 0, 1, 'train'),
	(21, '2024-03-19 14:13:04.852', '2024-03-19 14:39:55.548', NULL, 'qwen1.5-14b-train', 'qwen1.5-14b', '#!/bin/bash\nexport AUTH=sk-001\nexport JOB_ID={{.JobId}}\n\nexport CUDA_DEVICE_MAX_CONNECTIONS=1\n\n\nGPUS_PER_NODE={{.ProcPerNode}}\nNNODES=1\nNODE_RANK=0\nMASTER_ADDR=localhost\nMASTER_PORT={{.MasterPort}}\nUSE_LORA={{.Lora}}\nQ_LORA=False\n\nMODEL=\"{{.BaseModelPath}}\" # Set the path if you do not want to load from huggingface directly\n# ATTENTION: specify the path to your training data, which should be a json file consisting of a list of conversations.\n# See the section for finetuning in README for more information.\nDATA=\"{{.DataPath}}\"\n# 验证集\nEVAL_DATA=\"{{.ValidationFile}}\"\nDS_CONFIG_PATH=\"ds_config_zero3.json\"\n\nDISTRIBUTED_ARGS=\"\n    --nproc_per_node $GPUS_PER_NODE \\\n    --nnodes $NNODES \\\n    --node_rank $NODE_RANK \\\n    --master_addr $MASTER_ADDR \\\n    --master_port $MASTER_PORT\n\"\nif [ \"$USE_LORA\" == \"true\" ]; then\n    USE_LORA=True\n    DS_CONFIG_PATH=\"ds_config_zero2.json\"\nelse\n    USE_LORA=False\n    DS_CONFIG_PATH=\"ds_config_zero3.json\"\nfi\n\nmkdir -p /data/train-data/\nwget -O {{.DataPath}} {{.FileUrl}}\n\ntorchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \\\n    --model_name_or_path $MODEL \\\n    --data_path $DATA \\\n    --bf16 True \\\n    --output_dir {{.OutputDir}} \\\n    --num_train_epochs {{.TrainEpoch}} \\\n    --per_device_train_batch_size {{.TrainBatchSize}} \\\n    --per_device_eval_batch_size {{.EvalBatchSize}} \\\n    --gradient_accumulation_steps {{.AccumulationSteps}} \\\n    --evaluation_strategy \"no\" \\\n    --save_strategy \"steps\" \\\n    --save_steps 1000 \\\n    --save_total_limit 10 \\\n    --learning_rate {{.LearningRate}} \\\n    --weight_decay 0.1 \\\n    --adam_beta2 0.95 \\\n    --warmup_ratio 0.01 \\\n    --lr_scheduler_type \"cosine\" \\\n    --logging_steps 1 \\\n    --report_to \"none\" \\\n    --model_max_length {{.ModelMaxLength}} \\\n    --gradient_checkpointing True \\\n    --lazy_preprocess True \\\n    --use_lora ${USE_LORA} \\\n    --q_lora ${Q_LORA} \\\n    --deepspeed $DS_CONFIG_PATH', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-14b', '/app/finetune.py', '/data/ft-model/', 32768, 0, 1, 'train'),
	(22, '2024-03-19 14:13:04.852', '2024-03-19 14:40:05.661', NULL, 'qwen1.5-72b-train', 'qwen1.5-72b', '#!/bin/bash\nexport AUTH=sk-001\nexport JOB_ID={{.JobId}}\n\nexport CUDA_DEVICE_MAX_CONNECTIONS=1\n\n\nGPUS_PER_NODE={{.ProcPerNode}}\nNNODES=1\nNODE_RANK=0\nMASTER_ADDR=localhost\nMASTER_PORT={{.MasterPort}}\nUSE_LORA={{.Lora}}\nQ_LORA=False\n\nMODEL=\"{{.BaseModelPath}}\" # Set the path if you do not want to load from huggingface directly\n# ATTENTION: specify the path to your training data, which should be a json file consisting of a list of conversations.\n# See the section for finetuning in README for more information.\nDATA=\"{{.DataPath}}\"\n# 验证集\nEVAL_DATA=\"{{.ValidationFile}}\"\nDS_CONFIG_PATH=\"ds_config_zero3.json\"\n\nDISTRIBUTED_ARGS=\"\n    --nproc_per_node $GPUS_PER_NODE \\\n    --nnodes $NNODES \\\n    --node_rank $NODE_RANK \\\n    --master_addr $MASTER_ADDR \\\n    --master_port $MASTER_PORT\n\"\nif [ \"$USE_LORA\" == \"true\" ]; then\n    USE_LORA=True\n    DS_CONFIG_PATH=\"ds_config_zero2.json\"\nelse\n    USE_LORA=False\n    DS_CONFIG_PATH=\"ds_config_zero3.json\"\nfi\n\nmkdir -p /data/train-data/\nwget -O {{.DataPath}} {{.FileUrl}}\n\ntorchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \\\n    --model_name_or_path $MODEL \\\n    --data_path $DATA \\\n    --bf16 True \\\n    --output_dir {{.OutputDir}} \\\n    --num_train_epochs {{.TrainEpoch}} \\\n    --per_device_train_batch_size {{.TrainBatchSize}} \\\n    --per_device_eval_batch_size {{.EvalBatchSize}} \\\n    --gradient_accumulation_steps {{.AccumulationSteps}} \\\n    --evaluation_strategy \"no\" \\\n    --save_strategy \"steps\" \\\n    --save_steps 1000 \\\n    --save_total_limit 10 \\\n    --learning_rate {{.LearningRate}} \\\n    --weight_decay 0.1 \\\n    --adam_beta2 0.95 \\\n    --warmup_ratio 0.01 \\\n    --lr_scheduler_type \"cosine\" \\\n    --logging_steps 1 \\\n    --report_to \"none\" \\\n    --model_max_length {{.ModelMaxLength}} \\\n    --gradient_checkpointing True \\\n    --lazy_preprocess True \\\n    --use_lora ${USE_LORA} \\\n    --q_lora ${Q_LORA} \\\n    --deepspeed $DS_CONFIG_PATH', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-72b', '/app/finetune.py', '/data/ft-model/', 32768, 0, 1, 'train'),
	(23, '2024-03-19 14:13:04.852', '2024-03-19 14:40:15.699', NULL, 'qwen1.5-72b-chat-train', 'qwen1.5-72b-chat', '#!/bin/bash\nexport AUTH=sk-001\nexport JOB_ID={{.JobId}}\n\nexport CUDA_DEVICE_MAX_CONNECTIONS=1\n\n\nGPUS_PER_NODE={{.ProcPerNode}}\nNNODES=1\nNODE_RANK=0\nMASTER_ADDR=localhost\nMASTER_PORT={{.MasterPort}}\nUSE_LORA={{.Lora}}\nQ_LORA=False\n\nMODEL=\"{{.BaseModelPath}}\" # Set the path if you do not want to load from huggingface directly\n# ATTENTION: specify the path to your training data, which should be a json file consisting of a list of conversations.\n# See the section for finetuning in README for more information.\nDATA=\"{{.DataPath}}\"\n# 验证集\nEVAL_DATA=\"{{.ValidationFile}}\"\nDS_CONFIG_PATH=\"ds_config_zero3.json\"\n\nDISTRIBUTED_ARGS=\"\n    --nproc_per_node $GPUS_PER_NODE \\\n    --nnodes $NNODES \\\n    --node_rank $NODE_RANK \\\n    --master_addr $MASTER_ADDR \\\n    --master_port $MASTER_PORT\n\"\nif [ \"$USE_LORA\" == \"true\" ]; then\n    USE_LORA=True\n    DS_CONFIG_PATH=\"ds_config_zero2.json\"\nelse\n    USE_LORA=False\n    DS_CONFIG_PATH=\"ds_config_zero3.json\"\nfi\n\nmkdir -p /data/train-data/\nwget -O {{.DataPath}} {{.FileUrl}}\n\ntorchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \\\n    --model_name_or_path $MODEL \\\n    --data_path $DATA \\\n    --bf16 True \\\n    --output_dir {{.OutputDir}} \\\n    --num_train_epochs {{.TrainEpoch}} \\\n    --per_device_train_batch_size {{.TrainBatchSize}} \\\n    --per_device_eval_batch_size {{.EvalBatchSize}} \\\n    --gradient_accumulation_steps {{.AccumulationSteps}} \\\n    --evaluation_strategy \"no\" \\\n    --save_strategy \"steps\" \\\n    --save_steps 1000 \\\n    --save_total_limit 10 \\\n    --learning_rate {{.LearningRate}} \\\n    --weight_decay 0.1 \\\n    --adam_beta2 0.95 \\\n    --warmup_ratio 0.01 \\\n    --lr_scheduler_type \"cosine\" \\\n    --logging_steps 1 \\\n    --report_to \"none\" \\\n    --model_max_length {{.ModelMaxLength}} \\\n    --gradient_checkpointing True \\\n    --lazy_preprocess True \\\n    --use_lora ${USE_LORA} \\\n    --q_lora ${Q_LORA} \\\n    --deepspeed $DS_CONFIG_PATH', '', 'dudulu/qwen-train:v0.2.36-0319', '', '/data/base-model/qwen1-5-72b-chat', '/app/finetune.py', '/data/ft-model/', 32768, 0, 1, 'train');
`

	initSysDictSql = `INSERT INTO sys_dict (id, created_at, updated_at, deleted_at, parent_id, code, dict_value, dict_label, dict_type, sort, remark)
VALUES
	(2, '2023-11-22 16:19:52.000', '2024-01-29 10:32:18.000', NULL, 0, 'speak_gender', 'gender', '性别', 'int', 1, '性别'),
	(3, '2023-11-22 16:23:19.000', '2024-01-29 10:32:18.000', NULL, 2, 'speak_gender', '1', '男', 'int', 1, '性别:男'),
	(4, '2023-11-22 16:24:27.000', '2024-01-29 10:32:18.000', NULL, 2, 'speak_gender', '2', '女', 'int', 0, '性别:女'),
	(5, '2023-11-23 10:17:31.000', '2023-11-30 10:42:43.000', NULL, 0, 'speak_age_group', 'speak_age_group', '年龄段', 'int', 0, ''),
	(6, '2023-11-23 10:18:31.000', '2023-11-23 10:20:51.000', NULL, 5, 'speak_age_group', '1', '少年', 'int', 5, ''),
	(7, '2023-11-23 10:18:46.000', '2023-11-23 10:20:51.000', NULL, 5, 'speak_age_group', '2', '青年', 'int', 4, ''),
	(8, '2023-11-23 10:18:56.000', '2023-11-23 10:20:51.000', NULL, 5, 'speak_age_group', '3', '中年', 'int', 4, ''),
	(9, '2023-11-23 10:19:21.000', '2023-11-23 10:20:51.000', NULL, 5, 'speak_age_group', '4', '老年', 'int', 2, ''),
	(10, '2023-11-23 10:25:07.000', '2023-11-30 10:42:43.000', NULL, 0, 'speak_style', 'speak_style', '风格', 'int', 0, ''),
	(11, '2023-11-23 10:25:53.000', '2023-11-23 10:25:53.000', NULL, 10, 'speak_style', '1', '温柔', 'int', 5, ''),
	(12, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '2', '阳光', 'int', 4, ''),
	(13, '2023-11-23 10:28:24.000', '2023-12-08 14:04:05.000', NULL, 0, 'speak_area', 'speak_area', '适应范围', 'int', 0, ''),
	(14, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '1', '客服', 'int', 5, ''),
	(15, '2023-11-23 10:29:14.000', '2023-11-23 10:29:14.000', NULL, 13, 'speak_area', '2', '小说', 'int', 4, ''),
	(16, '2023-11-23 10:32:39.000', '2023-11-23 10:32:39.000', NULL, 0, 'speak_lang', 'speak_lang', '语言', 'string', 0, ''),
	(17, '2023-11-23 10:33:28.000', '2023-11-23 10:33:28.000', NULL, 16, 'speak_lang', 'zh-CN', '中文（普通话，简体）', 'string', 100, ''),
	(18, '2023-11-23 10:34:08.000', '2023-11-23 10:34:08.000', NULL, 16, 'speak_lang', 'zh-HK', '中文（粤语，繁体）', 'string', 99, ''),
	(19, '2023-11-23 10:34:30.000', '2023-11-23 10:34:30.000', NULL, 16, 'speak_lang', 'en-US', '英语（美国）', 'string', 98, ''),
	(20, '2023-11-23 10:35:07.000', '2023-11-23 10:35:07.000', NULL, 16, 'speak_lang', 'en-GB', '英语（英国）', 'string', 97, ''),
	(21, '2023-11-23 10:44:23.000', '2023-11-23 10:44:23.000', NULL, 0, 'speak_provider', 'speak_provider', '供应商', 'string', 0, ''),
	(22, '2023-11-23 10:44:50.000', '2023-11-23 10:44:50.000', NULL, 21, 'speak_provider', 'azure', '微软', 'string', 0, ''),
	(23, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '3', '自然流畅', 'int', 0, ''),
	(24, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '4', '亲切温和', 'int', 0, ''),
	(25, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '5', '温柔甜美', 'int', 0, ''),
	(26, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '6', '成熟知性', 'int', 0, ''),
	(27, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '7', '大气浑厚', 'int', 0, ''),
	(28, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '8', '稳重磁性', 'int', 0, ''),
	(29, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '9', '年轻时尚', 'int', 0, ''),
	(30, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '10', '轻声耳语', 'int', 0, ''),
	(31, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '11', '可爱甜美', 'int', 0, ''),
	(32, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '12', '呆萌可爱', 'int', 0, ''),
	(33, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '13', '激情力度', 'int', 0, ''),
	(34, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '14', '饱满活泼', 'int', 0, ''),
	(35, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '15', '诙谐幽默', 'int', 0, ''),
	(36, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '16', '淳朴方言', 'int', 0, ''),
	(37, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '3', '新闻', 'int', 0, ''),
	(38, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '4', '纪录片', 'int', 0, ''),
	(39, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '5', '解说', 'int', 0, ''),
	(40, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '6', '教育', 'int', 0, ''),
	(41, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '7', '广告', 'int', 0, ''),
	(42, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '8', '直播', 'int', 0, ''),
	(43, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '9', '助理', 'int', 0, ''),
	(44, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '10', '特色', 'int', 0, ''),
	(45, '2023-11-23 10:34:08.000', '2023-11-24 17:54:17.000', NULL, 16, 'speak_lang', 'zh-CN-henan', '中文（中原官话河南，简体）', 'string', 99, ''),
	(46, '2023-11-23 10:34:08.000', '2023-11-24 17:54:19.000', NULL, 16, 'speak_lang', 'zh-CN-liaoning', '中文（东北官话，简体）', 'string', 99, ''),
	(47, '2023-11-23 10:34:08.000', '2023-11-24 17:54:20.000', NULL, 16, 'speak_lang', 'zh-TW', '中文（台湾普通话，繁体）', 'string', 99, ''),
	(48, '2023-11-23 10:34:08.000', '2023-11-24 17:54:22.000', NULL, 16, 'speak_lang', 'zh-CN-GUANGXI', '中文（广西口音普通话，简体）', 'string', 99, ''),
	(49, '2023-11-23 10:35:07.000', '2023-11-23 10:35:07.000', NULL, 16, 'speak_lang', 'ko-KR', '韩语(韩国)', 'string', 97, ''),
	(50, '2023-11-23 10:35:07.000', '2023-11-24 19:45:54.000', NULL, 16, 'speak_lang', 'ja-JP', '日语（日本）', 'string', 97, ''),
	(51, '2023-11-23 10:35:07.000', '2023-11-24 19:45:54.000', NULL, 16, 'speak_lang', 'fil-PH', '菲律宾语（菲律宾）', 'string', 97, ''),
	(52, '2023-11-23 10:35:07.000', '2023-11-24 19:45:54.000', NULL, 16, 'speak_lang', 'es-MX', '西班牙语(墨西哥)', 'string', 97, ''),
	(53, '2023-11-23 10:35:07.000', '2023-11-24 19:45:54.000', NULL, 16, 'speak_lang', 'ru-RU', '俄语（俄罗斯）', 'string', 97, ''),
	(54, '2023-11-23 10:32:39.000', '2023-11-23 10:32:39.000', NULL, 0, 'audio_tagged_lang', 'audio_tagged_lang', '标记语言', 'string', 0, ''),
	(55, '2023-11-23 10:32:39.000', '2023-12-06 14:55:30.000', NULL, 54, 'audio_tagged_lang', 'zh', '中文', 'string', 1001, ''),
	(56, '2023-11-23 10:32:39.000', '2023-11-23 10:32:39.000', NULL, 54, 'audio_tagged_lang', 'en', '英文', 'string', 99, ''),
	(57, '2023-11-23 10:32:39.000', '2023-11-23 10:32:39.000', NULL, 54, 'audio_tagged_lang', 'es', '西班牙语', 'string', 98, ''),
	(58, '2023-11-23 10:32:39.000', '2023-11-23 10:32:39.000', NULL, 54, 'audio_tagged_lang', 'de', '德语', 'string', 97, ''),
	(59, '2023-11-23 10:32:39.000', '2023-11-28 17:32:38.000', NULL, 54, 'audio_tagged_lang', 'tl', '他加禄语', 'string', 96, ''),
	(60, '2023-11-23 10:32:39.000', '2023-11-23 10:32:39.000', NULL, 54, 'audio_tagged_lang', 'fil', '菲律宾语', 'string', 95, ''),
	(61, '2023-11-23 10:32:39.000', '2023-11-23 10:32:39.000', NULL, 0, 'sys_dict_type', 'sys_dict_type', '字典类型', 'string', 0, ''),
	(62, '2023-11-23 10:32:39.000', '2023-12-08 13:55:43.000', NULL, 61, 'sys_dict_type', 'string', '字符串类型', 'string', 100, ''),
	(63, '2023-11-23 10:32:39.000', '2023-12-08 13:55:43.000', NULL, 61, 'sys_dict_type', 'int', '数字类型', 'int', 99, ''),
	(64, '2023-11-23 10:32:39.000', '2023-12-08 13:55:43.000', NULL, 61, 'sys_dict_type', 'bool', '布尔类型', 'bool', 98, ''),
	(87, '2023-12-07 17:18:32.000', '2023-12-07 17:18:32.000', NULL, 0, 'language', 'language', '国际化', 'string', 0, ''),
	(88, '2023-12-14 15:53:15.000', '2023-12-14 16:37:19.000', NULL, 0, 'model_eval_dataset_type', 'model_eval_dataset_type', '模型评估数据集类型', 'string', 110, ''),
	(89, '2023-12-14 15:54:28.000', '2023-12-14 16:37:19.000', NULL, 88, 'model_eval_dataset_type', 'train', '训练集', 'string', 99, ''),
	(90, '2023-12-14 15:54:57.000', '2023-12-14 16:37:19.000', NULL, 88, 'model_eval_dataset_type', 'custom', '自定义', 'string', 98, ''),
	(91, '2023-12-14 15:55:56.000', '2023-12-14 15:59:39.000', NULL, 0, 'model_eval_metric', 'model_eval_metric', '模型评估指标', 'string', 99, ''),
	(92, '2023-12-14 15:58:48.000', '2023-12-14 15:58:48.000', NULL, 0, 'model_eval_status', 'model_eval_status', '模型评估状态', 'string', 99, ''),
	(93, '2023-12-14 16:00:31.000', '2023-12-14 16:00:31.000', NULL, 92, 'model_eval_status', 'pending', '等待评估', 'string', 99, ''),
	(94, '2023-12-14 16:00:44.000', '2023-12-14 16:00:44.000', NULL, 92, 'model_eval_status', 'running', '正在评估', 'string', 98, ''),
	(95, '2023-12-14 16:00:56.000', '2023-12-14 16:00:56.000', NULL, 92, 'model_eval_status', 'success', '评估成功', 'string', 97, ''),
	(96, '2023-12-14 16:01:09.000', '2023-12-14 16:01:09.000', NULL, 92, 'model_eval_status', 'failed', '评估失败', 'string', 97, ''),
	(97, '2023-12-14 16:01:23.000', '2023-12-14 16:01:23.000', NULL, 92, 'model_eval_status', 'cancel', '评估取消', 'string', 96, ''),
	(98, '2023-12-14 16:11:47.000', '2023-12-14 16:14:05.000', NULL, 91, 'model_eval_metric', 'equal', '完全匹配', 'string', 98, ''),
	(99, '2023-12-14 17:20:04.000', '2023-12-14 17:20:04.000', NULL, 0, 'model_deploy_status', 'model_deploy_status', '模型部署状态', 'string', 0, ''),
	(100, '2023-12-14 17:20:34.000', '2023-12-15 11:39:30.000', NULL, 99, 'model_deploy_status', 'pending', '部署中', 'string', 0, ''),
	(101, '2023-12-14 17:20:45.000', '2023-12-15 11:39:35.000', NULL, 99, 'model_deploy_status', 'running', '运行中', 'string', 0, ''),
	(102, '2023-12-14 17:21:06.000', '2023-12-14 17:22:34.000', NULL, 99, 'model_deploy_status', 'success', '完成', 'string', 0, ''),
	(103, '2023-12-14 17:24:54.000', '2023-12-14 17:24:54.000', NULL, 99, 'model_deploy_status', 'failed', '失败', 'string', 0, ''),
	(104, '2023-12-14 20:44:04.000', '2024-01-30 11:19:36.000', NULL, 0, 'model_provider_name', 'model_provider_name', '模型供应商', 'string', 0, ''),
	(105, '2023-12-14 20:45:32.000', '2024-01-30 11:19:36.000', NULL, 104, 'model_provider_name', 'LocalAI', 'LocalAI', 'string', 0, ''),
	(106, '2023-12-14 20:45:44.000', '2024-01-30 11:19:36.000', NULL, 104, 'model_provider_name', 'OpenAI', 'OpenAI', 'string', 0, ''),
	(107, '2023-12-15 15:09:42.000', '2023-12-15 15:09:42.000', NULL, 0, 'digitalhuman_synthesis_status', 'digitalhuman_synthesis_status', '数字人合成状态', 'string', 0, ''),
	(108, '2023-12-15 15:10:48.000', '2023-12-15 15:10:48.000', NULL, 107, 'digitalhuman_synthesis_status', 'running', '合成中', 'string', 0, ''),
	(109, '2023-12-15 15:11:12.000', '2023-12-15 15:11:12.000', NULL, 107, 'digitalhuman_synthesis_status', 'success', '已完成', 'string', 0, ''),
	(110, '2023-12-15 15:11:30.000', '2023-12-15 15:11:30.000', NULL, 107, 'digitalhuman_synthesis_status', 'failed', '失败', 'string', 0, ''),
	(111, '2023-12-15 15:11:48.000', '2023-12-15 15:11:48.000', NULL, 107, 'digitalhuman_synthesis_status', 'waiting', '等待中', 'string', 0, ''),
	(112, '2023-12-15 15:12:02.000', '2023-12-15 15:12:02.000', NULL, 107, 'digitalhuman_synthesis_status', 'cancel', '已取消', 'string', 0, ''),
	(113, '2023-12-20 19:07:56.000', '2023-12-20 19:07:56.000', NULL, 0, 'digitalhuman_posture', 'digitalhuman_posture', '数字人姿势', 'int', 0, ''),
	(114, '2023-12-20 19:09:10.000', '2023-12-21 10:11:35.000', NULL, 113, 'digitalhuman_posture', '1', '全身', 'int', 0, ''),
	(115, '2023-12-20 19:09:39.000', '2023-12-21 10:11:44.000', NULL, 113, 'digitalhuman_posture', '2', '半身', 'int', 0, ''),
	(116, '2023-12-20 19:10:22.000', '2023-12-21 10:11:53.000', NULL, 113, 'digitalhuman_posture', '3', '大半身', 'int', 0, ''),
	(117, '2023-12-20 19:10:34.000', '2023-12-21 10:11:58.000', NULL, 113, 'digitalhuman_posture', '4', '坐姿', 'int', 0, ''),
	(118, '2023-12-20 19:16:05.000', '2023-12-20 19:16:05.000', NULL, 0, 'digitalhuman_resolution', 'digitalhuman_resolution', '数字人分辨率', 'int', 0, ''),
	(119, '2023-12-20 19:20:03.000', '2023-12-20 19:20:03.000', NULL, 118, 'digitalhuman_resolution', '1', '480P', 'int', 0, ''),
	(120, '2023-12-20 19:20:22.000', '2023-12-20 19:20:22.000', NULL, 118, 'digitalhuman_resolution', '2', '720P', 'int', 0, ''),
	(121, '2023-12-20 19:20:43.000', '2023-12-20 19:20:43.000', NULL, 118, 'digitalhuman_resolution', '3', '1080P', 'int', 0, ''),
	(122, '2023-12-20 19:20:51.000', '2023-12-20 19:20:51.000', NULL, 118, 'digitalhuman_resolution', '4', '2K', 'int', 0, ''),
	(123, '2023-12-20 19:21:13.000', '2023-12-20 19:21:13.000', NULL, 118, 'digitalhuman_resolution', '5', '4K', 'int', 0, ''),
	(124, '2023-12-20 19:21:31.000', '2023-12-20 19:21:31.000', NULL, 118, 'digitalhuman_resolution', '6', '8K', 'int', 0, ''),
	(125, '2023-12-22 11:13:26.000', '2023-12-22 11:22:53.000', '2023-12-22 11:23:07.000', 0, 'model_type', 'model_type', '模型类型', 'string', 0, ''),
	(126, '2023-12-22 11:17:02.000', '2023-12-22 11:22:53.000', '2023-12-22 11:23:07.000', 125, 'model_type', 'train', '微调训练', 'string', 0, ''),
	(127, '2023-12-22 11:17:42.000', '2023-12-22 11:22:53.000', '2023-12-22 11:23:07.000', 125, 'model_type', 'inference', '模型推理', 'string', 0, ''),
	(128, '2023-12-22 11:24:15.000', '2023-12-22 11:24:15.000', NULL, 0, 'template_type', 'template_type', '模板类型', 'string', 0, ''),
	(129, '2023-12-22 11:25:41.000', '2023-12-22 11:30:11.000', NULL, 128, 'template_type', 'train', '微调训练', 'string', 0, ''),
	(130, '2023-12-22 11:26:50.000', '2023-12-22 11:29:04.000', NULL, 128, 'template_type', 'inference', '模型推理', 'string', 0, ''),
	(131, '2024-01-08 16:44:16.000', '2024-01-09 16:57:15.000', '2024-01-09 16:57:29.000', 0, 'model_quantify', 'model_quantify', '模型量化', 'string', 0, '模特部署量化'),
	(132, '2024-01-08 16:45:30.000', '2024-01-09 16:57:15.000', '2024-01-09 16:57:29.000', 131, 'model_quantify', 'bf16', '半精度', 'int', 0, ''),
	(133, '2024-01-08 16:47:23.000', '2024-01-09 16:57:15.000', '2024-01-09 16:57:29.000', 131, 'model_quantify', '8bit', '1/4精度', 'int', 1, '四分之一精度'),
	(134, '2024-01-09 10:50:12.000', '2024-01-09 10:50:12.000', NULL, 0, 'model_deploy_label', 'model_deploy_label', '模型部署标签', 'string', 0, ''),
	(135, '2024-01-09 10:51:24.000', '2024-03-19 17:17:07.017', NULL, 134, 'model_deploy_label', 'cpu-aigc-model', 'cpu-aigc-model', 'string', 0, ''),
	(136, '2024-01-09 10:52:19.000', '2024-01-09 10:52:19.000', NULL, 0, 'model_deploy_quantization', 'model_deploy_quantization', '模型部署量化', 'string', 0, '模型部署量化'),
	(137, '2024-01-09 10:52:40.000', '2024-01-09 10:52:40.000', NULL, 136, 'model_deploy_quantization', 'float16', 'float16', 'string', 0, ''),
	(138, '2024-01-09 10:52:46.000', '2024-01-09 10:52:46.000', NULL, 136, 'model_deploy_quantization', '8bit', '8bit', 'string', 0, ''),
	(156, '2024-01-23 10:10:54.000', '2024-01-23 10:10:54.000', NULL, 0, 'assistant_tool_type', 'assistant_tool_type', 'AI助手工具类型', 'string', 0, 'AI助手工具类型'),
	(157, '2024-01-23 10:12:23.000', '2024-01-25 15:06:41.000', NULL, 156, 'assistant_tool_type', 'function', 'API接口', 'string', 3, ''),
	(158, '2024-01-23 10:12:47.000', '2024-01-23 10:13:45.000', NULL, 156, 'assistant_tool_type', 'retrieval', '知识库', 'string', 2, ''),
	(159, '2024-01-23 10:13:01.000', '2024-01-23 10:13:25.000', NULL, 156, 'assistant_tool_type', 'code_interpreter', '代码执行', 'string', 1, ''),
	(160, '2024-01-24 11:07:37.000', '2024-01-24 11:07:37.000', NULL, 0, 'http_method', 'http_method', '请求方法', 'string', 111, 'http请求方法'),
	(161, '2024-01-24 11:08:10.000', '2024-01-24 11:09:02.000', NULL, 160, 'http_method', 'get', 'GET', 'string', 4, ''),
	(162, '2024-01-24 11:08:21.000', '2024-01-24 11:09:08.000', NULL, 160, 'http_method', 'post', 'POST', 'string', 3, ''),
	(163, '2024-01-24 11:08:36.000', '2024-01-24 11:09:12.000', NULL, 160, 'http_method', 'put', 'PUT', 'string', 2, ''),
	(164, '2024-01-24 11:09:54.000', '2024-01-24 11:09:54.000', NULL, 160, 'http_method', 'delete', 'DEL', 'string', 1, ''),
	(165, '2024-01-24 11:14:47.000', '2024-01-24 11:14:47.000', NULL, 0, 'programming_language', 'programming_language', '编程语言', 'string', 112, ''),
	(166, '2024-01-24 11:15:12.000', '2024-01-24 11:15:12.000', NULL, 165, 'programming_language', 'python', 'Python', 'string', 1, ''),
	(172, '2024-03-19 11:25:22.770', '2024-03-19 11:25:22.770', NULL, 0, 'textannotation_type', 'textannotation_type', '文本标注类型', 'string', 0, ''),
	(173, '2024-03-19 11:25:47.575', '2024-03-19 11:25:47.575', NULL, 172, 'textannotation_type', 'rag', '检索增强生成', 'string', 0, ''),
	(174, '2024-03-19 11:26:00.272', '2024-03-19 11:26:00.272', NULL, 172, 'textannotation_type', 'faq', '知识问答', 'string', 0, ''),
	(175, '2024-03-19 11:26:12.417', '2024-03-19 11:26:12.417', NULL, 172, 'textannotation_type', 'general', '通用', 'string', 0, ''),
	(176, '2024-03-19 14:01:56.036', '2024-03-19 14:01:56.036', NULL, 0, 'model_type', 'model_type', '模型类型', 'string', 0, '模型类型：文本模型，语音模型，数字人模型等'),
	(177, '2024-03-19 14:02:19.712', '2024-03-19 14:02:19.712', NULL, 176, 'model_type', 'text-generation', '文本', 'string', 0, ''),
	(178, '2024-03-19 14:02:28.164', '2024-03-19 14:02:28.164', NULL, 176, 'model_type', 'embedding', 'embedding', 'string', 0, ''),
	(179, '2024-03-19 18:04:29.830', '2024-03-19 18:04:29.830', NULL, 0, 'model_evaluate_target_type', 'model_evaluate_target_type', '模型评测指标', 'string', 0, ''),
	(180, '2024-03-19 18:05:04.555', '2024-03-19 18:05:04.555', NULL, 179, 'model_evaluate_target_type', 'Acc', 'ACC', 'string', 0, ''),
	(181, '2024-03-19 18:05:14.515', '2024-03-19 18:05:14.515', NULL, 179, 'model_evaluate_target_type', 'F1', 'F1', 'string', 0, ''),
	(182, '2024-03-19 18:05:22.487', '2024-03-19 18:05:22.487', NULL, 179, 'model_evaluate_target_type', 'BLEU', 'BLEU', 'string', 0, ''),
	(183, '2024-03-19 18:05:30.619', '2024-03-19 18:05:30.619', NULL, 179, 'model_evaluate_target_type', 'Rouge', 'Rouge', 'string', 0, ''),
	(184, '2024-03-19 18:05:38.596', '2024-03-19 18:05:38.596', NULL, 179, 'model_evaluate_target_type', 'five', '五维图', 'string', 0, '');`
)
