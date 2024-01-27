/**
 * @Time : 2021/11/26 10:24 AM
 * @Author : solacowa@gmail.com
 * @File : service_job
 * @Software: GoLand
 */

package service

import (
	"github.com/spf13/cobra"
)

var (
	jobCmd = &cobra.Command{
		Use:               "job command <args> [flags]",
		Short:             "任务命令",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
可用的配置类型：
[clear]

aigc-admin job -h
`,
	}
)
