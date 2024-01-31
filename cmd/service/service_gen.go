/**
 * @Time : 24/04/2020 11:49 AM
 * @Author : solacowa@gmail.com
 * @File : service_gen_table
 * @Software: GoLand
 */

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
				if rdb != nil {
					_ = level.Debug(logger).Log("rdb", "close", "err", rdb.Close())
				}
			}()

			if len(args) > 0 && args[0] == "all" {
				return generateTable()
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			logger = log.NewLogfmtLogger(os.Stdout)
			// 连接数据库
			if strings.EqualFold(dbDrive, "mysql") {
				dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
					mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
				var dbErr error
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
			}
			return nil
		},
	}
)

func generateTable() (err error) {
	_ = logger.Log("migrate", "table", "Chat", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.Chat{}))
	_ = logger.Log("migrate", "table", "ChatAllowUser", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.ChatAllowUser{}))
	_ = logger.Log("migrate", "table", "ChatRole", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.ChatRole{}))
	_ = logger.Log("migrate", "table", "ChatConversation", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.ChatConversation{}))
	_ = logger.Log("migrate", "table", "ChatSystemPrompt", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.ChatSystemPrompt{}))
	_ = logger.Log("migrate", "table", "ChatPromptTypes", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.ChatPromptTypes{}))
	_ = logger.Log("migrate", "table", "ChatChannels", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.ChatChannels{}))
	_ = logger.Log("migrate", "table", "ChatPrompts", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.ChatPrompts{}))
	_ = logger.Log("migrate", "table", "ChatChannelModels", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.ChatChannelModels{}))
	_ = logger.Log("migrate", "table", "ChatMessages", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.ChatMessages{}))
	_ = logger.Log("migrate", "table", "Dataset", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.Dataset{}))
	_ = logger.Log("migrate", "table", "DatasetSample", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.DatasetSample{}))
	_ = logger.Log("migrate", "table", "Assistants", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.Assistants{}))
	_ = logger.Log("migrate", "table", "Tools", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.Tools{}))
	_ = logger.Log("migrate", "table", "AssistantToolAssociations", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.AssistantToolAssociations{}))
	_ = logger.Log("migrate", "table", "Files", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.Files{}))
	_ = logger.Log("migrate", "table", "SysAudit", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.SysAudit{}))
	_ = logger.Log("migrate", "table", "FineTuningTrainJob", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.FineTuningTrainJob{}))
	_ = logger.Log("migrate", "table", "FineTuningTemplate", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.FineTuningTemplate{}))
	_ = logger.Log("migrate", "table", "Tenants", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.Tenants{}))
	_ = logger.Log("migrate", "table", "Models", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.Models{}))
	_ = logger.Log("migrate", "table", "SysDict", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.SysDict{}))
	_ = logger.Log("migrate", "table", "ModelDeploy", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.ModelDeploy{}))
	_ = logger.Log("migrate", "table", "LLMEvalResults", gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").AutoMigrate(types.LLMEvalResults{}))
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
	_ = logger.Log("init", "data", "models", gormDB.Create(&types.Models{
		ProviderName: types.ModelProviderLocalAI,
		ModelType:    types.ModelTypeTextGeneration,
		ModelName:    "qwen-14b-base",
		MaxTokens:    8192,
		IsPrivate:    true,
		IsFineTuning: false,
		Enabled:      true,
		Remark:       "通义千问14b",
		Tenants: []types.Tenants{
			tenant,
		},
		Parameters: 14.17,
	}))
	_ = logger.Log("init", "data", "models", gormDB.Create(&types.Models{
		ProviderName: types.ModelProviderOpenAI,
		ModelType:    types.ModelTypeTextGeneration,
		ModelName:    "gpt-3.5-turbo",
		MaxTokens:    4096,
		IsPrivate:    false,
		IsFineTuning: false,
		Enabled:      true,
		Remark:       "OpenAI GPT-3.5-turbo",
		Tenants: []types.Tenants{
			tenant,
		},
		Parameters: 20,
	}))

	_ = logger.Log("init", "data", "sys_dict", gormDB.Exec(initSysDictSql).Error)
	_ = logger.Log("init", "data", "sys_dict", gormDB.Exec(ftTemplateSql).Error)
	return err
}

var (
	ftTemplateSql = `INSERT INTO fine_tuning_template (created_at, updated_at, deleted_at, name, base_model, content, params, train_image, remark, base_model_path, template_type, script_file, output_dir, max_tokens, lora, enabled)
VALUES
	('2023-12-22 13:56:32.000', '2023-12-26 13:39:49.277', NULL, 'tpl-qwen-14b-base', 'qwen-14b-base', '#!/bin/bash\nexport AUTH=sk-001\nexport JOB_ID={{.JobId}}\n\nexport CUDA_DEVICE_MAX_CONNECTIONS=1\n\n\nGPUS_PER_NODE={{.ProcPerNode}}\nNNODES=1\nNODE_RANK=0\nMASTER_ADDR=localhost\nMASTER_PORT=6001\nUSE_LORA={{.Lora}}\n\n\nMODEL=\"{{.BaseModelPath}}\" # Set the path if you do not want to load from huggingface directly\n# ATTENTION: specify the path to your training data, which should be a json file consisting of a list of conversations.\n# See the section for finetuning in README for more information.\nDATA=\"{{.DataPath}}\"\n\nDISTRIBUTED_ARGS=\"\n    --nproc_per_node $GPUS_PER_NODE \\\n    --nnodes $NNODES \\\n    --node_rank $NODE_RANK \\\n    --master_addr $MASTER_ADDR \\\n    --master_port $MASTER_PORT\n\"\n\nmkdir -p /data/train-data/\nwget -O {{.DataPath}} {{.FileUrl}}\n\nif [ \"$USE_LORA\" -eq 1 ]; then\n    USE_LORA=True\n    DS_CONFIG_PATH=\"ds_config_zero2.json\"\nelse\n    USE_LORA=False\n    DS_CONFIG_PATH=\"ds_config_zero3.json\"\nfi\n\ntorchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \\\n    --model_name_or_path $MODEL \\\n    --data_path $DATA \\\n    --bf16 True \\\n    --output_dir {{.OutputDir}} \\\n    --num_train_epochs {{.TrainEpoch}} \\\n    --per_device_train_batch_size {{.TrainBatchSize}} \\\n    --per_device_eval_batch_size {{.EvalBatchSize}} \\\n    --gradient_accumulation_steps {{.AccumulationSteps}} \\\n    --evaluation_strategy \"no\" \\\n    --save_strategy \"steps\" \\\n    --save_steps 1000 \\\n    --save_total_limit 10 \\\n    --learning_rate 3e-4 \\\n    --weight_decay 0.1 \\\n    --adam_beta2 0.95 \\\n    --warmup_ratio 0.01 \\\n    --lr_scheduler_type \"cosine\" \\\n    --logging_steps 1 \\\n    --report_to \"none\" \\\n    --model_max_length {{.ModelMaxLength}} \\\n    --lazy_preprocess True \\\n    --use_lora $USE_LORA \\\n    --gradient_checkpointing \\\n    --deepspeed ${DS_CONFIG_PATH}\n', '', 'icowan/qwen-train:v0.2.34-1220', 'qwen-14b训练模版', '/data/base-model/qwen-14b-base', 'train', '/app/finetune.py', '/data/ft-model', 8192, 0, 1);
`
	initSysDictSql = `INSERT INTO sys_dict (id, parent_id, code, dict_value, dict_label, dict_type, sort, remark, created_at, updated_at, deleted_at)
VALUES
	(2, 0, 'speak_gender', 'gender', '性别', 'int', 1, '性别', '2023-11-22 16:19:52', '2024-01-29 10:32:18', NULL),
	(3, 2, 'speak_gender', '1', '男', 'int', 1, '性别:男', '2023-11-22 16:23:19', '2024-01-29 10:32:18', NULL),
	(4, 2, 'speak_gender', '2', '女', 'int', 0, '性别:女', '2023-11-22 16:24:27', '2024-01-29 10:32:18', NULL),
	(5, 0, 'speak_age_group', 'speak_age_group', '年龄段', 'int', 0, '', '2023-11-23 10:17:31', '2023-11-30 10:42:43', NULL),
	(6, 5, 'speak_age_group', '1', '少年', 'int', 5, '', '2023-11-23 10:18:31', '2023-11-23 10:20:51', NULL),
	(7, 5, 'speak_age_group', '2', '青年', 'int', 4, '', '2023-11-23 10:18:46', '2023-11-23 10:20:51', NULL),
	(8, 5, 'speak_age_group', '3', '中年', 'int', 4, '', '2023-11-23 10:18:56', '2023-11-23 10:20:51', NULL),
	(9, 5, 'speak_age_group', '4', '老年', 'int', 2, '', '2023-11-23 10:19:21', '2023-11-23 10:20:51', NULL),
	(10, 0, 'speak_style', 'speak_style', '风格', 'int', 0, '', '2023-11-23 10:25:07', '2023-11-30 10:42:43', NULL),
	(11, 10, 'speak_style', '1', '温柔', 'int', 5, '', '2023-11-23 10:25:53', '2023-11-23 10:25:53', NULL),
	(12, 10, 'speak_style', '2', '阳光', 'int', 4, '', '2023-11-23 10:26:04', '2023-11-23 10:26:04', NULL),
	(13, 0, 'speak_area', 'speak_area', '适应范围', 'int', 0, '', '2023-11-23 10:28:24', '2023-12-08 14:04:05', NULL),
	(14, 13, 'speak_area', '1', '客服', 'int', 5, '', '2023-11-23 10:28:57', '2023-11-23 10:28:57', NULL),
	(15, 13, 'speak_area', '2', '小说', 'int', 4, '', '2023-11-23 10:29:14', '2023-11-23 10:29:14', NULL),
	(16, 0, 'speak_lang', 'speak_lang', '语言', 'string', 0, '', '2023-11-23 10:32:39', '2023-11-23 10:32:39', NULL),
	(17, 16, 'speak_lang', 'zh-CN', '中文（普通话，简体）', 'string', 100, '', '2023-11-23 10:33:28', '2023-11-23 10:33:28', NULL),
	(18, 16, 'speak_lang', 'zh-HK', '中文（粤语，繁体）', 'string', 99, '', '2023-11-23 10:34:08', '2023-11-23 10:34:08', NULL),
	(19, 16, 'speak_lang', 'en-US', '英语（美国）', 'string', 98, '', '2023-11-23 10:34:30', '2023-11-23 10:34:30', NULL),
	(20, 16, 'speak_lang', 'en-GB', '英语（英国）', 'string', 97, '', '2023-11-23 10:35:07', '2023-11-23 10:35:07', NULL),
	(21, 0, 'speak_provider', 'speak_provider', '供应商', 'string', 0, '', '2023-11-23 10:44:23', '2023-11-23 10:44:23', NULL),
	(22, 21, 'speak_provider', 'azure', '微软', 'string', 0, '', '2023-11-23 10:44:50', '2023-11-23 10:44:50', NULL),
	(23, 10, 'speak_style', '3', '自然流畅', 'int', 0, '', '2023-11-23 10:26:04', '2023-11-23 10:26:04', NULL),
	(24, 10, 'speak_style', '4', '亲切温和', 'int', 0, '', '2023-11-23 10:26:04', '2023-11-23 10:26:04', NULL),
	(25, 10, 'speak_style', '5', '温柔甜美', 'int', 0, '', '2023-11-23 10:26:04', '2023-11-23 10:26:04', NULL),
	(26, 10, 'speak_style', '6', '成熟知性', 'int', 0, '', '2023-11-23 10:26:04', '2023-11-23 10:26:04', NULL),
	(27, 10, 'speak_style', '7', '大气浑厚', 'int', 0, '', '2023-11-23 10:26:04', '2023-11-23 10:26:04', NULL),
	(28, 10, 'speak_style', '8', '稳重磁性', 'int', 0, '', '2023-11-23 10:26:04', '2023-11-23 10:26:04', NULL),
	(29, 10, 'speak_style', '9', '年轻时尚', 'int', 0, '', '2023-11-23 10:26:04', '2023-11-23 10:26:04', NULL),
	(30, 10, 'speak_style', '10', '轻声耳语', 'int', 0, '', '2023-11-23 10:26:04', '2023-11-23 10:26:04', NULL),
	(31, 10, 'speak_style', '11', '可爱甜美', 'int', 0, '', '2023-11-23 10:26:04', '2023-11-23 10:26:04', NULL),
	(32, 10, 'speak_style', '12', '呆萌可爱', 'int', 0, '', '2023-11-23 10:26:04', '2023-11-23 10:26:04', NULL),
	(33, 10, 'speak_style', '13', '激情力度', 'int', 0, '', '2023-11-23 10:26:04', '2023-11-23 10:26:04', NULL),
	(34, 10, 'speak_style', '14', '饱满活泼', 'int', 0, '', '2023-11-23 10:26:04', '2023-11-23 10:26:04', NULL),
	(35, 10, 'speak_style', '15', '诙谐幽默', 'int', 0, '', '2023-11-23 10:26:04', '2023-11-23 10:26:04', NULL),
	(36, 10, 'speak_style', '16', '淳朴方言', 'int', 0, '', '2023-11-23 10:26:04', '2023-11-23 10:26:04', NULL),
	(37, 13, 'speak_area', '3', '新闻', 'int', 0, '', '2023-11-23 10:28:57', '2023-11-23 10:28:57', NULL),
	(38, 13, 'speak_area', '4', '纪录片', 'int', 0, '', '2023-11-23 10:28:57', '2023-11-23 10:28:57', NULL),
	(39, 13, 'speak_area', '5', '解说', 'int', 0, '', '2023-11-23 10:28:57', '2023-11-23 10:28:57', NULL),
	(40, 13, 'speak_area', '6', '教育', 'int', 0, '', '2023-11-23 10:28:57', '2023-11-23 10:28:57', NULL),
	(41, 13, 'speak_area', '7', '广告', 'int', 0, '', '2023-11-23 10:28:57', '2023-11-23 10:28:57', NULL),
	(42, 13, 'speak_area', '8', '直播', 'int', 0, '', '2023-11-23 10:28:57', '2023-11-23 10:28:57', NULL),
	(43, 13, 'speak_area', '9', '助理', 'int', 0, '', '2023-11-23 10:28:57', '2023-11-23 10:28:57', NULL),
	(44, 13, 'speak_area', '10', '特色', 'int', 0, '', '2023-11-23 10:28:57', '2023-11-23 10:28:57', NULL),
	(45, 16, 'speak_lang', 'zh-CN-henan', '中文（中原官话河南，简体）', 'string', 99, '', '2023-11-23 10:34:08', '2023-11-24 17:54:17', NULL),
	(46, 16, 'speak_lang', 'zh-CN-liaoning', '中文（东北官话，简体）', 'string', 99, '', '2023-11-23 10:34:08', '2023-11-24 17:54:19', NULL),
	(47, 16, 'speak_lang', 'zh-TW', '中文（台湾普通话，繁体）', 'string', 99, '', '2023-11-23 10:34:08', '2023-11-24 17:54:20', NULL),
	(48, 16, 'speak_lang', 'zh-CN-GUANGXI', '中文（广西口音普通话，简体）', 'string', 99, '', '2023-11-23 10:34:08', '2023-11-24 17:54:22', NULL),
	(49, 16, 'speak_lang', 'ko-KR', '韩语(韩国)', 'string', 97, '', '2023-11-23 10:35:07', '2023-11-23 10:35:07', NULL),
	(50, 16, 'speak_lang', 'ja-JP', '日语（日本）', 'string', 97, '', '2023-11-23 10:35:07', '2023-11-24 19:45:54', NULL),
	(51, 16, 'speak_lang', 'fil-PH', '菲律宾语（菲律宾）', 'string', 97, '', '2023-11-23 10:35:07', '2023-11-24 19:45:54', NULL),
	(52, 16, 'speak_lang', 'es-MX', '西班牙语(墨西哥)', 'string', 97, '', '2023-11-23 10:35:07', '2023-11-24 19:45:54', NULL),
	(53, 16, 'speak_lang', 'ru-RU', '俄语（俄罗斯）', 'string', 97, '', '2023-11-23 10:35:07', '2023-11-24 19:45:54', NULL),
	(54, 0, 'audio_tagged_lang', 'audio_tagged_lang', '标记语言', 'string', 0, '', '2023-11-23 10:32:39', '2023-11-23 10:32:39', NULL),
	(55, 54, 'audio_tagged_lang', 'zh', '中文', 'string', 1001, '', '2023-11-23 10:32:39', '2023-12-06 14:55:30', NULL),
	(56, 54, 'audio_tagged_lang', 'en', '英文', 'string', 99, '', '2023-11-23 10:32:39', '2023-11-23 10:32:39', NULL),
	(57, 54, 'audio_tagged_lang', 'es', '西班牙语', 'string', 98, '', '2023-11-23 10:32:39', '2023-11-23 10:32:39', NULL),
	(58, 54, 'audio_tagged_lang', 'de', '德语', 'string', 97, '', '2023-11-23 10:32:39', '2023-11-23 10:32:39', NULL),
	(59, 54, 'audio_tagged_lang', 'tl', '他加禄语', 'string', 96, '', '2023-11-23 10:32:39', '2023-11-28 17:32:38', NULL),
	(60, 54, 'audio_tagged_lang', 'fil', '菲律宾语', 'string', 95, '', '2023-11-23 10:32:39', '2023-11-23 10:32:39', NULL),
	(61, 0, 'sys_dict_type', 'sys_dict_type', '字典类型', 'string', 0, '', '2023-11-23 10:32:39', '2023-11-23 10:32:39', NULL),
	(62, 61, 'sys_dict_type', 'string', '字符串类型', 'string', 100, '', '2023-11-23 10:32:39', '2023-12-08 13:55:43', NULL),
	(63, 61, 'sys_dict_type', 'int', '数字类型', 'int', 99, '', '2023-11-23 10:32:39', '2023-12-08 13:55:43', NULL),
	(64, 61, 'sys_dict_type', 'bool', '布尔类型', 'bool', 98, '', '2023-11-23 10:32:39', '2023-12-08 13:55:43', NULL),
	(65, 0, 'test3333', '63', 'test', 'string', 1001, '备注', '2023-11-23 10:32:39', '2023-12-14 15:57:37', '2023-12-14 15:57:51'),
	(66, 65, 'test3333', '5', 'test', 'string', 1001, '备注', '2023-11-23 10:32:39', '2023-12-14 15:57:37', '2023-12-14 15:57:51'),
	(67, 66, 'test3333', '3', 'test2', 'bool', 0, '', '2023-11-23 10:32:39', '2023-11-30 17:33:16', NULL),
	(68, 67, 'test3333', 'true', 'test1-1', 'string', 0, '', '2023-11-23 10:32:39', '2023-11-30 17:33:16', NULL),
	(69, 0, 'test_1', 'dictValue', '语言2', 'string', 1, '兰格畏惧', '2023-11-30 10:25:11', '2023-11-30 10:33:45', '2023-11-30 10:33:58'),
	(70, 0, 'test_3', 'test_3', '飞桨', 'string', 0, '', '2023-11-30 11:13:46', '2023-12-14 16:03:36', '2023-12-14 16:03:50'),
	(71, 65, 'test3333', '4', '你是谁的谁心疼有为了了了', 'string', 0, '', '2023-11-30 18:01:57', '2023-12-14 15:57:37', '2023-12-14 15:57:51'),
	(72, 2, 'speak_gender', '3', '未知', 'int', 0, '未知未知嗜血粒', '2023-12-01 11:13:07', '2023-12-01 14:38:02', '2023-12-01 14:38:15'),
	(73, 72, 'speak_gender', '4', '未知子项', 'bool', 0, '', '2023-12-01 11:22:46', '2023-12-01 14:38:02', '2023-12-01 14:38:15'),
	(74, 72, 'speak_gender', 'asdfd', 'asd', 'string', 0, '', '2023-12-01 11:33:44', '2023-12-01 14:37:47', '2023-12-01 14:38:00'),
	(75, 0, 'test_language', 'test_language', '测试语言3阿斯顿发斯蒂芬', 'string', 1, '测试语言的备注还是很长的测试语言的备注还是很长的测试语言的备注还是很长的测试语言的备注还是很长的测试语言的备注还是很长的测试语言的备注还是很长的', '2023-12-01 14:55:21', '2023-12-01 15:27:49', '2023-12-01 15:28:02'),
	(76, 75, 'test_language', '1', '英语23343423232', 'int', 0, '', '2023-12-01 15:00:51', '2023-12-01 15:27:49', '2023-12-01 15:28:02'),
	(77, 75, 'test_language', '2', '中文2323', 'int', 0, '', '2023-12-01 15:01:09', '2023-12-01 15:27:49', '2023-12-01 15:28:02'),
	(78, 75, 'test_language', '3', '法语', 'string', 0, '', '2023-12-01 15:01:23', '2023-12-01 15:27:38', '2023-12-01 15:27:51'),
	(79, 75, 'test_language', '4', '德语', 'int', 0, '', '2023-12-01 15:01:42', '2023-12-01 15:02:47', '2023-12-01 15:03:00'),
	(80, 79, 'test_language', '1', '小德', 'bool', 0, '', '2023-12-01 15:02:28', '2023-12-01 15:02:47', '2023-12-01 15:03:00'),
	(81, 79, 'test_language', '2', '大德', 'int', 0, '', '2023-12-01 15:02:37', '2023-12-01 15:02:38', '2023-12-01 15:02:51'),
	(82, 0, 'test2', 'test2', 'test23', 'int', 0, '', '2023-12-04 14:47:48', '2023-12-14 15:57:52', '2023-12-14 15:58:06'),
	(83, 65, 'test3333', 'en', '英文', 'string', 0, '', '2023-12-07 17:11:41', '2023-12-14 15:57:37', '2023-12-14 15:57:51'),
	(84, 83, 'test3333', '1', 'male', 'string', 0, '', '2023-12-07 17:12:17', '2023-12-07 17:12:17', NULL),
	(85, 65, 'test3333', 'zh', '中文', 'string', 0, '', '2023-12-07 17:12:33', '2023-12-14 15:57:37', '2023-12-14 15:57:51'),
	(86, 85, 'test3333', '1', '男', 'string', 0, '', '2023-12-07 17:12:41', '2023-12-07 17:12:41', NULL),
	(87, 0, 'language', 'language', '国际化', 'string', 0, '', '2023-12-07 17:18:32', '2023-12-07 17:18:32', NULL),
	(88, 0, 'model_eval_dataset_type', 'model_eval_dataset_type', '模型评估数据集类型', 'string', 110, '', '2023-12-14 15:53:15', '2023-12-14 16:37:19', NULL),
	(89, 88, 'model_eval_dataset_type', 'train', '训练集', 'string', 99, '', '2023-12-14 15:54:28', '2023-12-14 16:37:19', NULL),
	(90, 88, 'model_eval_dataset_type', 'custom', '自定义', 'string', 98, '', '2023-12-14 15:54:57', '2023-12-14 16:37:19', NULL),
	(91, 0, 'model_eval_metric', 'model_eval_metric', '模型评估指标', 'string', 99, '', '2023-12-14 15:55:56', '2023-12-14 15:59:39', NULL),
	(92, 0, 'model_eval_status', 'model_eval_status', '模型评估状态', 'string', 99, '', '2023-12-14 15:58:48', '2023-12-14 15:58:48', NULL),
	(93, 92, 'model_eval_status', 'pending', '等待评估', 'string', 99, '', '2023-12-14 16:00:31', '2023-12-14 16:00:31', NULL),
	(94, 92, 'model_eval_status', 'running', '正在评估', 'string', 98, '', '2023-12-14 16:00:44', '2023-12-14 16:00:44', NULL),
	(95, 92, 'model_eval_status', 'success', '评估成功', 'string', 97, '', '2023-12-14 16:00:56', '2023-12-14 16:00:56', NULL),
	(96, 92, 'model_eval_status', 'failed', '评估失败', 'string', 97, '', '2023-12-14 16:01:09', '2023-12-14 16:01:09', NULL),
	(97, 92, 'model_eval_status', 'cancel', '评估取消', 'string', 96, '', '2023-12-14 16:01:23', '2023-12-14 16:01:23', NULL),
	(98, 91, 'model_eval_metric', 'equal', '完全匹配', 'string', 98, '', '2023-12-14 16:11:47', '2023-12-14 16:14:05', NULL),
	(99, 0, 'model_deploy_status', 'model_deploy_status', '模型部署状态', 'string', 0, '', '2023-12-14 17:20:04', '2023-12-14 17:20:04', NULL),
	(100, 99, 'model_deploy_status', 'pending', '部署中', 'string', 0, '', '2023-12-14 17:20:34', '2023-12-15 11:39:30', NULL),
	(101, 99, 'model_deploy_status', 'running', '运行中', 'string', 0, '', '2023-12-14 17:20:45', '2023-12-15 11:39:35', NULL),
	(102, 99, 'model_deploy_status', 'success', '完成', 'string', 0, '', '2023-12-14 17:21:06', '2023-12-14 17:22:34', NULL),
	(103, 99, 'model_deploy_status', 'failed', '失败', 'string', 0, '', '2023-12-14 17:24:54', '2023-12-14 17:24:54', NULL),
	(104, 0, 'model_provider_name', 'model_provider_name', '模型供应商', 'string', 0, '', '2023-12-14 20:44:04', '2024-01-30 11:19:36', NULL),
	(105, 104, 'model_provider_name', 'LocalAI', 'LocalAI', 'string', 0, '', '2023-12-14 20:45:32', '2024-01-30 11:19:36', NULL),
	(106, 104, 'model_provider_name', 'OpenAI', 'OpenAI', 'string', 0, '', '2023-12-14 20:45:44', '2024-01-30 11:19:36', NULL),
	(107, 0, 'digitalhuman_synthesis_status', 'digitalhuman_synthesis_status', '数字人合成状态', 'string', 0, '', '2023-12-15 15:09:42', '2023-12-15 15:09:42', NULL),
	(108, 107, 'digitalhuman_synthesis_status', 'running', '合成中', 'string', 0, '', '2023-12-15 15:10:48', '2023-12-15 15:10:48', NULL),
	(109, 107, 'digitalhuman_synthesis_status', 'success', '已完成', 'string', 0, '', '2023-12-15 15:11:12', '2023-12-15 15:11:12', NULL),
	(110, 107, 'digitalhuman_synthesis_status', 'failed', '失败', 'string', 0, '', '2023-12-15 15:11:30', '2023-12-15 15:11:30', NULL),
	(111, 107, 'digitalhuman_synthesis_status', 'waiting', '等待中', 'string', 0, '', '2023-12-15 15:11:48', '2023-12-15 15:11:48', NULL),
	(112, 107, 'digitalhuman_synthesis_status', 'cancel', '已取消', 'string', 0, '', '2023-12-15 15:12:02', '2023-12-15 15:12:02', NULL),
	(113, 0, 'digitalhuman_posture', 'digitalhuman_posture', '数字人姿势', 'int', 0, '', '2023-12-20 19:07:56', '2023-12-20 19:07:56', NULL),
	(114, 113, 'digitalhuman_posture', '1', '全身', 'int', 0, '', '2023-12-20 19:09:10', '2023-12-21 10:11:35', NULL),
	(115, 113, 'digitalhuman_posture', '2', '半身', 'int', 0, '', '2023-12-20 19:09:39', '2023-12-21 10:11:44', NULL),
	(116, 113, 'digitalhuman_posture', '3', '大半身', 'int', 0, '', '2023-12-20 19:10:22', '2023-12-21 10:11:53', NULL),
	(117, 113, 'digitalhuman_posture', '4', '坐姿', 'int', 0, '', '2023-12-20 19:10:34', '2023-12-21 10:11:58', NULL),
	(118, 0, 'digitalhuman_resolution', 'digitalhuman_resolution', '数字人分辨率', 'int', 0, '', '2023-12-20 19:16:05', '2023-12-20 19:16:05', NULL),
	(119, 118, 'digitalhuman_resolution', '1', '480P', 'int', 0, '', '2023-12-20 19:20:03', '2023-12-20 19:20:03', NULL),
	(120, 118, 'digitalhuman_resolution', '2', '720P', 'int', 0, '', '2023-12-20 19:20:22', '2023-12-20 19:20:22', NULL),
	(121, 118, 'digitalhuman_resolution', '3', '1080P', 'int', 0, '', '2023-12-20 19:20:43', '2023-12-20 19:20:43', NULL),
	(122, 118, 'digitalhuman_resolution', '4', '2K', 'int', 0, '', '2023-12-20 19:20:51', '2023-12-20 19:20:51', NULL),
	(123, 118, 'digitalhuman_resolution', '5', '4K', 'int', 0, '', '2023-12-20 19:21:13', '2023-12-20 19:21:13', NULL),
	(124, 118, 'digitalhuman_resolution', '6', '8K', 'int', 0, '', '2023-12-20 19:21:31', '2023-12-20 19:21:31', NULL),
	(125, 0, 'model_type', 'model_type', '模型类型', 'string', 0, '', '2023-12-22 11:13:26', '2023-12-22 11:22:53', '2023-12-22 11:23:07'),
	(126, 125, 'model_type', 'train', '微调训练', 'string', 0, '', '2023-12-22 11:17:02', '2023-12-22 11:22:53', '2023-12-22 11:23:07'),
	(127, 125, 'model_type', 'inference', '模型推理', 'string', 0, '', '2023-12-22 11:17:42', '2023-12-22 11:22:53', '2023-12-22 11:23:07'),
	(128, 0, 'template_type', 'template_type', '模板类型', 'string', 0, '', '2023-12-22 11:24:15', '2023-12-22 11:24:15', NULL),
	(129, 128, 'template_type', 'train', '微调训练', 'string', 0, '', '2023-12-22 11:25:41', '2023-12-22 11:30:11', NULL),
	(130, 128, 'template_type', 'inference', '模型推理', 'string', 0, '', '2023-12-22 11:26:50', '2023-12-22 11:29:04', NULL),
	(131, 0, 'model_quantify', 'model_quantify', '模型量化', 'string', 0, '模特部署量化', '2024-01-08 16:44:16', '2024-01-09 16:57:15', '2024-01-09 16:57:29'),
	(132, 131, 'model_quantify', 'bf16', '半精度', 'int', 0, '', '2024-01-08 16:45:30', '2024-01-09 16:57:15', '2024-01-09 16:57:29'),
	(133, 131, 'model_quantify', '8bit', '1/4精度', 'int', 1, '四分之一精度', '2024-01-08 16:47:23', '2024-01-09 16:57:15', '2024-01-09 16:57:29'),
	(134, 0, 'model_deploy_label', 'model_deploy_label', '模型部署标签', 'string', 0, '', '2024-01-09 10:50:12', '2024-01-09 10:50:12', NULL),
	(135, 134, 'model_deploy_label', 'a100-40x10', 'a100-40x10', 'string', 0, '', '2024-01-09 10:51:24', '2024-01-09 10:51:24', NULL),
	(136, 0, 'model_deploy_quantization', 'model_deploy_quantization', '模型部署量化', 'string', 0, '模型部署量化', '2024-01-09 10:52:19', '2024-01-09 10:52:19', NULL),
	(137, 136, 'model_deploy_quantization', 'float16', 'float16', 'string', 0, '', '2024-01-09 10:52:40', '2024-01-09 10:52:40', NULL),
	(138, 136, 'model_deploy_quantization', '8bit', '8bit', 'string', 0, '', '2024-01-09 10:52:46', '2024-01-09 10:52:46', NULL),
	(139, 0, 'vrp_model_type', 'vrp_model_type', '声纹比对模型类型', 'string', 0, '', '2024-01-16 16:11:09', '2024-01-16 16:11:09', NULL),
	(140, 139, 'vrp_model_type', 'CAMPPlus', 'CAMPPlus', 'string', 0, '', '2024-01-16 16:11:44', '2024-01-16 16:11:44', NULL),
	(141, 139, 'vrp_model_type', 'CAMPPlus-bbig', 'CAMPPlus-bbig', 'string', 0, '', '2024-01-16 16:11:56', '2024-01-16 16:11:56', NULL),
	(142, 139, 'vrp_model_type', 'CAMPPlus-big', 'CAMPPlus-big', 'string', 0, '', '2024-01-16 16:12:07', '2024-01-16 16:12:07', NULL),
	(143, 139, 'vrp_model_type', 'ERes2Net', 'ERes2Net', 'string', 0, '', '2024-01-16 16:12:16', '2024-01-16 16:12:16', NULL),
	(144, 139, 'vrp_model_type', 'ERes2Net-big', 'ERes2Net-big', 'string', 0, '', '2024-01-16 16:12:26', '2024-01-16 16:12:26', NULL),
	(145, 139, 'vrp_model_type', 'EcapaTdnn', 'EcapaTdnn', 'string', 0, '', '2024-01-16 16:12:36', '2024-01-16 16:12:36', NULL),
	(146, 139, 'vrp_model_type', 'Res2Net', 'Res2Net', 'string', 0, '', '2024-01-16 16:12:44', '2024-01-16 16:12:44', NULL),
	(147, 139, 'vrp_model_type', 'ResNetSE', 'ResNetSE', 'string', 0, '', '2024-01-16 16:12:54', '2024-01-16 16:12:54', NULL),
	(148, 139, 'vrp_model_type', 'TDNN', 'TDNN', 'string', 0, '', '2024-01-16 16:13:07', '2024-01-16 16:13:07', NULL),
	(149, 0, 'esrgan_model_type', 'esrgan_model_type', '图像超分模型类型', 'string', 0, '', '2024-01-17 18:45:20', '2024-01-17 18:45:20', NULL),
	(150, 149, 'esrgan_model_type', 'RealESRGAN_x4plus', 'RealESRGAN_x4plus', 'string', 0, '', '2024-01-17 18:45:43', '2024-01-17 18:45:43', NULL),
	(151, 149, 'esrgan_model_type', 'RealESRNet_x4plus', 'RealESRNet_x4plus', 'string', 0, '', '2024-01-17 18:45:56', '2024-01-17 18:45:56', NULL),
	(152, 149, 'esrgan_model_type', 'RealESRGAN_x4plus_anime_6B', 'RealESRGAN_x4plus_anime_6B', 'string', 0, '', '2024-01-17 18:46:09', '2024-01-17 18:46:09', NULL),
	(153, 149, 'esrgan_model_type', 'RealESRGAN_x2plus', 'RealESRGAN_x2plus', 'string', 0, '', '2024-01-17 18:46:21', '2024-01-17 18:46:21', NULL),
	(154, 149, 'esrgan_model_type', 'realesr-animevideov3', 'realesr-animevideov3', 'string', 0, '', '2024-01-17 18:46:31', '2024-01-17 18:46:31', NULL),
	(155, 149, 'esrgan_model_type', 'realesr-general-x4v3', 'realesr-general-x4v3', 'string', 0, '', '2024-01-17 18:46:49', '2024-01-17 18:46:49', NULL),
	(156, 0, 'assistant_tool_type', 'assistant_tool_type', 'AI助手工具类型', 'string', 0, 'AI助手工具类型', '2024-01-23 10:10:54', '2024-01-23 10:10:54', NULL),
	(157, 156, 'assistant_tool_type', 'function', 'API接口', 'string', 3, '', '2024-01-23 10:12:23', '2024-01-25 15:06:41', NULL),
	(158, 156, 'assistant_tool_type', 'retrieval', '知识库', 'string', 2, '', '2024-01-23 10:12:47', '2024-01-23 10:13:45', NULL),
	(159, 156, 'assistant_tool_type', 'code_interpreter', '代码执行', 'string', 1, '', '2024-01-23 10:13:01', '2024-01-23 10:13:25', NULL),
	(160, 0, 'http_method', 'http_method', '请求方法', 'string', 111, 'http请求方法', '2024-01-24 11:07:37', '2024-01-24 11:07:37', NULL),
	(161, 160, 'http_method', 'get', 'GET', 'string', 4, '', '2024-01-24 11:08:10', '2024-01-24 11:09:02', NULL),
	(162, 160, 'http_method', 'post', 'POST', 'string', 3, '', '2024-01-24 11:08:21', '2024-01-24 11:09:08', NULL),
	(163, 160, 'http_method', 'put', 'PUT', 'string', 2, '', '2024-01-24 11:08:36', '2024-01-24 11:09:12', NULL),
	(164, 160, 'http_method', 'delete', 'DEL', 'string', 1, '', '2024-01-24 11:09:54', '2024-01-24 11:09:54', NULL),
	(165, 0, 'programming_language', 'programming_language', '编程语言', 'string', 112, '', '2024-01-24 11:14:47', '2024-01-24 11:14:47', NULL),
	(166, 165, 'programming_language', 'python', 'Python', 'string', 1, '', '2024-01-24 11:15:12', '2024-01-24 11:15:12', NULL),
	(167, 0, 'denoise_sample_rate', 'denoise_sample_rate', '音频降噪采样率', 'int', 0, '', '2024-01-25 17:45:22', '2024-01-25 17:45:22', NULL),
	(168, 167, 'denoise_sample_rate', '0', '原始采样率', 'int', 0, '', '2024-01-25 17:45:51', '2024-01-25 17:45:51', NULL),
	(169, 167, 'denoise_sample_rate', '16', '16K', 'int', 0, '', '2024-01-25 17:46:04', '2024-01-25 17:46:04', NULL),
	(170, 104, 'model_provider_name', 'test', '测试', 'int', 0, '', '2024-01-30 10:43:10', '2024-01-30 10:52:02', '2024-01-30 10:52:18'),
	(171, 104, 'model_provider_name', 'test', 'test供应商', 'int', 0, '', '2024-01-31 10:38:32', '2024-01-31 10:38:32', NULL);`
)
