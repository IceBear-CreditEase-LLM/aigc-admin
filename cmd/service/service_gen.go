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
	err = initData()
	if err != nil {
		return err
	}
	return
}

// 初始化数据
func initData() (err error) {
	tenant := types.Tenants{
		Name:           "系统租户",
		PublicTenantID: uuid.New().String(),
		ContactEmail:   "admin@admin.com",
	}
	_ = logger.Log("init", "data", "SysDict", gormDB.Create(&tenant).Error)
	password, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_ = logger.Log("init", "data", "account", gormDB.Save(&types.Accounts{
		Email:        "admin@admin.com",
		Nickname:     "admin",
		Language:     "zh",
		IsLdap:       false,
		PasswordHash: string(password),
		Status:       true,
		Tenants:      []types.Tenants{tenant},
	}).Error)
	if aigcChannelKey == "" {
		aigcChannelKey = "sk-" + string(util.Krand(48, util.KC_RAND_KIND_ALL))
	}
	_ = logger.Log("init", "data", "ChatChannels", gormDB.Create(&types.ChatChannels{
		Name:       "default",
		Alias:      "默认渠道",
		Remark:     "默认渠道",
		Quota:      10000,
		Models:     "default",
		OnlyOpenAI: false,
		ApiKey:     aigcChannelKey,
		Email:      "admin@admin.com",
		TenantId:   tenant.ID,
	}).Error)
	return err
}
