package service

import (
	"context"
	"encoding/json"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/pkg/finetuning"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/types"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/util"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
	"time"
)

var (
	jobFineTuningCmd = &cobra.Command{
		Use:               "finetuning command <args> [flags]",
		Short:             "微调任务命令",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
可用的配置类型：
[run-waiting-train, running-log]

aigc-admin job -h
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err = prepare(cmd.Context()); err != nil {
				return errors.Wrap(err, "prepare")
			}
			fineTuningSvc = finetuning.New(traceId, logger, store, serviceS3Bucket, serviceS3AccessKey, serviceS3SecretKey, apiSvc, rdb, aigcDataCfsPath)
			return nil
		},
	}

	jobFineTuningJobRunWaitingTrainCmd = &cobra.Command{
		Use:               `run-waiting-train [flags]`,
		Short:             "微调任务等待训练",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
aigc-admin job finetuning run-waiting-train
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fineTuningRunWaitingTrain(cmd.Context())
		},
	}

	jobFineTuningJobRunningJobLogCmd = &cobra.Command{
		Use:               `running-log [flags]`,
		Short:             "同步正在训练脚本日志",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
aigc-admin job finetuning running-log
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			runningJobs, err := store.FineTuning().FindFineTuningJobRunning(cmd.Context())
			if err != nil {
				_ = level.Error(logger).Log("msg", "find running job failed", "err", err.Error())
				return err
			}
			return fineTuningRunningJobLog(cmd.Context(), runningJobs)
		},
	}
)

func fineTuningRunWaitingTrain(ctx context.Context) (err error) {
	return fineTuningSvc.RunWaitingTrain(ctx)
}

type logEntry struct {
	Timestamp    time.Time `json:"timestamp"`
	Loss         float64   `json:"loss"`
	LearningRate float64   `json:"learning_rate"`
	Epoch        float64   `json:"epoch"`
}

func fineTuningRunningJobLog(ctx context.Context, jobs []types.FineTuningTrainJob) (err error) {
	for _, job := range jobs {
		//jobLog, err := apiSvc.Paas().GetJobPodsLog(ctx, "aigc", job.PaasJobName)
		//if err != nil {
		//	_ = level.Warn(logger).Log("msg", "get job pods log failed", "err", err.Error())
		//	continue
		//}
		var jobLog string
		jobLog, err = apiSvc.DockerApi().Logs(ctx, job.PaasJobName)
		if err != nil {
			_ = level.Warn(logger).Log("msg", "get job pods log failed", "err", err.Error())
			continue
		}

		lineArr := strings.Split(jobLog, "\n")
		re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z) (\{.*?\})`)

		var logEntryList []logEntry

		for _, log := range lineArr {
			matches := re.FindStringSubmatch(log)
			if len(matches) == 3 {
				timestampStr, jsonStr := matches[1], matches[2]

				timestamp, err := time.Parse(time.RFC3339Nano, timestampStr)
				if err != nil {
					_ = level.Warn(logger).Log("msg", "parse timestamp failed", "err", err.Error())
					continue
				}

				jsonStr = strings.Replace(jsonStr, "'", "\"", -1)        // 将单引号替换为双引号
				jsonStr = strings.Replace(jsonStr, "False", "false", -1) // 将 False 替换为 false
				jsonStr = strings.Replace(jsonStr, "True", "true", -1)   // 将 True 替换为 true

				var entry logEntry
				err = json.Unmarshal([]byte(jsonStr), &entry)
				if err != nil {
					_ = level.Warn(logger).Log("msg", "unmarshal json failed", "err", err.Error())
					continue
				}

				entry.Timestamp = timestamp
				logEntryList = append(logEntryList, entry)
			}
		}
		if len(logEntryList) < 1 {
			continue
		}
		lastLine := logEntryList[len(logEntryList)-1]
		job.ProgressLoss = lastLine.Loss
		job.ProgressLearningRate = lastLine.LearningRate
		job.ProgressEpochs = lastLine.Epoch
		progress := util.RoundToFourDecimalPlaces(lastLine.Epoch / float64(job.TrainEpoch))
		job.Progress = progress
		job.TrainLog = jobLog
		if err = store.FineTuning().UpdateFineTuningJob(ctx, &job); err != nil {
			_ = level.Warn(logger).Log("msg", "update job log failed", "err", err)
		}
	}
	return
}
