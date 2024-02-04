package finetuning

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api/alarm"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api/dockerapi"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/encode"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/pkg/files"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/finetuning"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/types"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/util"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
	"io"
	"k8s.io/apimachinery/pkg/util/rand"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Service 模型微调模块
type Service interface {
	// CreateJob 创建微调任务
	CreateJob(ctx context.Context, tenantId uint, request CreateJobRequest) (response JobResponse, err error)
	// ListJob 获取微调任务列表
	ListJob(ctx context.Context, tenantId uint, request ListJobRequest) (response ListJobResponse, err error)
	// CancelJob 取消微调任务
	CancelJob(ctx context.Context, tenantId uint, jobId string) (err error)
	// DashBoard 微调任务面板
	DashBoard(ctx context.Context, tenantId uint) (res DashBoardResponse, err error)
	// DeleteJob 删除微调任务
	DeleteJob(ctx context.Context, tenantId uint, jobId string) (err error)
	// GetJob 获取微调任务详情
	GetJob(ctx context.Context, tenantId uint, jobId string) (response JobResponse, err error)
	// ListTemplate 获取微调模板列表
	ListTemplate(ctx context.Context, tenantId uint, request ListTemplateRequest) (response ListTemplateResponse, err error)
	// Estimate 微调时间预估
	Estimate(ctx context.Context, tenantId uint, request CreateJobRequest) (response EstimateResponse, err error)
	// _createJob 创建训练任务
	_createJob(ctx context.Context, tenantId, channelId uint, trainingFileId, model, suffix, validationFile string, epochs int) (res jobResult, err error)
	// _cancelJob 取消微调任务
	_cancelJob(ctx context.Context, channelId uint, fineTuningJob string) (res jobResult, err error)
	// UpdateJobFinishedStatus 更新微调完成的任务状态 (训练脚本调用) train.sh 执行完之后调用
	UpdateJobFinishedStatus(ctx context.Context, fineTuningJob string, status types.TrainStatus, message string) (err error)
	// RunWaitingTrain 执行等待中的脚本 (定时任务) 通常每两分钟或每五分钟执行一次
	RunWaitingTrain(ctx context.Context) (err error)
	// _createFineTuningJob 创建微调任务
	_createFineTuningJob(ctx context.Context, jobId string) (err error)
}

type service struct {
	traceId     string
	logger      log.Logger
	store       repository.Repository
	api         api.Service
	namespace   string
	dataCfsPath string
	mu          sync.Mutex
	rdb         redis.UniversalClient
	fileSvc     files.Service
}

func (s *service) _createJob(ctx context.Context, tenantId, channelId uint, trainingFileId, baseModel, suffix, validationFile string, epochs int) (res jobResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	// 获取文件 查询文件信息
	fileInfo, err := s.store.Files().FindFileByFileId(ctx, trainingFileId)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Files", "FindFileByFileId", "err", err.Error())
		return
	}

	// 判断是否是gpt模型
	if strings.Contains(baseModel, "gpt-3.5") || strings.Contains(baseModel, "gpt-4") {
		// 1. 获取文件上传到openai
		fileId, err := s.uploadFileToOpenAi(ctx, &fileInfo)
		if err != nil {
			_ = level.Error(logger).Log("service", "uploadFileToOpenAi", "err", err.Error())
			return res, err
		}
		// 2. 调用openai 创建微调任务
		openAiFtJob, err := s.api.FastChat().CreateFineTuningJob(ctx, openai.FineTuningJobRequest{
			TrainingFile:    fileId,
			ValidationFile:  "",
			Model:           baseModel,
			Hyperparameters: &openai.Hyperparameters{Epochs: epochs},
			Suffix:          suffix,
		})
		if err != nil {
			_ = level.Error(logger).Log("api.FastChat", "CreateFineTuningJob", "err", err.Error())
			return res, errors.Wrap(err, "api.FastChat.CreateFineTuningJob")
		}
		// 3. 入库
		// 创建job
		ftJob := &types.FineTuningTrainJob{
			JobId:          openAiFtJob.ID,
			ChannelId:      channelId,
			BaseModel:      baseModel,
			FineTunedModel: openAiFtJob.FineTunedModel,
			ValidationFile: openAiFtJob.ValidationFile,
			TrainEpoch:     epochs,
			FileUrl:        fileInfo.S3Url,
			TrainStatus:    types.TrainStatusWaiting,
			TenantID:       tenantId,
		}
		if err = s.store.FineTuning().CreateFineTuningJob(ctx, ftJob); err != nil {
			res.Error = err.Error()
			_ = level.Error(logger).Log("repository.FineTuningJob", "CreateFineTuningJob", "err", err.Error())
			return res, err
		}

		res.CreatedAt = openAiFtJob.CreatedAt
		res.Id = openAiFtJob.ID
		res.Model = openAiFtJob.Model
		res.TrainingFile = openAiFtJob.TrainingFile
		res.ValidationFile = openAiFtJob.ValidationFile
		res.HyperParameters.NEpochs = ftJob.TrainEpoch
		//res.Status = ftJob.TrainStatus
		res.FineTunedModel = openAiFtJob.FineTunedModel

		return res, nil
	}

	panUrl, err := s._fileConvertAlpaca(ctx, baseModel, fileInfo.S3Url)
	if err != nil {
		_ = level.Warn(logger).Log("service", "_fileConvertAlpaca", "err", err.Error())
		return
	}

	// 根据模型获取模版
	ftJobTpl, err := s.store.FineTuning().FindFineTuningTemplateByType(ctx, baseModel, types.TemplateTypeTrain)
	if err != nil {
		_ = level.Warn(logger).Log("repository.FineTuningJob", "FindFineTuningTemplateByModel", "err", err.Error())
		return
	}
	if !strings.EqualFold(suffix, "") {
		suffix = ":" + suffix
	}
	suffix = string(util.Krand(4, util.KC_RAND_KIND_LOWER)) + suffix

	fineTunedModel := fmt.Sprintf("ft::%s:%d-%s", baseModel, tenantId, suffix)

	// 创建job
	ftJob := &types.FineTuningTrainJob{
		JobId:          uuid.New().String(),
		ChannelId:      channelId,
		TemplateId:     ftJobTpl.ID,
		BaseModel:      baseModel,
		TrainEpoch:     epochs,
		BaseModelPath:  ftJobTpl.BaseModelPath,
		DataPath:       fmt.Sprintf("/data/train-data/%s", trainingFileId),
		OutputDir:      fmt.Sprintf("%s/ft-%s-%d-%s", ftJobTpl.OutputDir, baseModel, tenantId, strings.ReplaceAll(strings.ReplaceAll(suffix, ".", "-"), ":", "-")),
		ScriptFile:     ftJobTpl.ScriptFile,
		MasterPort:     rand.IntnRange(20000, 30000),
		FileUrl:        panUrl,
		TrainStatus:    types.TrainStatusWaiting,
		FineTunedModel: fineTunedModel,
		ProcPerNode:    2,  // 暂时写死 这个得与GPU数量对应
		TrainBatchSize: 8,  // 暂时写死
		EvalBatchSize:  32, // 暂时写死
		TenantID:       tenantId,
	}
	if err = s.store.FineTuning().CreateFineTuningJob(ctx, ftJob); err != nil {
		res.Error = err.Error()
		_ = level.Error(logger).Log("repository.FineTuningJob", "CreateFineTuningJob", "err", err.Error())
		return
	}

	defer func() {
		if err != nil {
			ftJob.TrainStatus = types.TrainStatusFailed
			ftJob.ErrorMessage = err.Error()
			if err = s.store.FineTuning().UpdateFineTuningJob(ctx, ftJob); err != nil {
				_ = level.Error(logger).Log("repository.FineTuningJob", "UpdateFineTuningJob", "err", err.Error())
				return
			}
		}
	}()

	res.CreatedAt = ftJob.CreatedAt.Unix()
	res.Id = ftJob.JobId
	res.Model = ftJob.BaseModel
	res.TrainingFile = ftJob.DataPath
	res.ValidationFile = ftJob.ValidationFile
	res.HyperParameters.NEpochs = ftJob.TrainEpoch
	//res.Status = ftJob.TrainStatus
	res.FineTunedModel = ftJob.FineTunedModel
	//res.TrainedTokens = ftJob.TrainedTokens

	// 如果没有waiting 状态的任务直接调用_createJob创建执行
	hasRunning, err := s.store.FineTuning().HasRunningJob(ctx)
	if err != nil {
		_ = level.Warn(logger).Log("repository.FineTuningJob", "HasRunningJob", "err", err.Error())
		return res, nil
	}
	if !hasRunning {
		// 如果没有等待中的任务，直接创建
		if err = s._createFineTuningJob(ctx, ftJob.JobId); err != nil {
			_ = level.Error(logger).Log("service", "_createFineTuningJob", "err", err.Error())
			return res, err
		}
		_ = level.Info(logger).Log("msg", "没有正在运行的任务，直接创建", "jobId", ftJob.JobId)
	}

	return
}

func (s *service) uploadFileToOpenAi(ctx context.Context, fileInfo *types.Files) (fileId string, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	body, err := getHttpFileBody(fileInfo.S3Url)
	if err != nil {
		_ = level.Error(logger).Log("getHttpFileBody", "getHttpFileBody", "err", err.Error())
		return "", err
	}
	// 创建临时文件
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		_ = level.Error(logger).Log("msg", "创建临时文件失败", "err", err.Error())
		return
	}
	defer tmpfile.Close()

	_ = level.Info(logger).Log("msg", "创建临时文件", "tmpfile", tmpfile.Name())

	// 写入去除序号后的文本到临时文件
	if _, err := tmpfile.Write(body); err != nil {
		_ = level.Error(logger).Log("msg", "写入临时文件失败", "err", err.Error())
		return "", err
	}
	defer func(tempFilePath string) {
		_ = os.Remove(tempFilePath)
	}(tmpfile.Name())

	// 上传到openai
	openAiRes, err := s.api.FastChat().UploadFile(ctx, openai.GPT3Dot5Turbo, fileInfo.Name, tmpfile.Name(), fileInfo.Purpose)
	if err != nil {
		_ = level.Error(logger).Log("api.FastChat", "UploadFile", "err", err.Error())
		return
	}
	return openAiRes.ID, nil
}

func (s *service) _cancelJob(ctx context.Context, channelId uint, fineTuningJob string) (res jobResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	// 查看相关数据
	jobInfo, err := s.store.FineTuning().FindFineTuningJobByJobId(ctx, fineTuningJob, "Template")
	if err != nil {
		_ = level.Error(logger).Log("repository.FineTuningJob", "FindFineTuningJobByJobId", "err", err.Error())
		return
	}

	err = s.api.DockerApi().Remove(ctx, jobInfo.PaasJobName)
	if err != nil {
		_ = level.Error(logger).Log("api.DockerApi", "Remove", "err", err.Error())
		err = errors.Wrap(err, "api.DockerApi.Remove")
		return
	}
	// 调用paas强制删除job
	//if err = s.api.Paas().DeleteJob(ctx, s.namespace, jobInfo.PaasJobName); err != nil {
	//	_ = level.Error(logger).Log("api.Paas", "DeleteJob", "err", err.Error())
	//	//return
	//}
	// 调用paas删除configmap
	//serviceName := strings.ReplaceAll(strings.ReplaceAll(jobInfo.FineTunedModel, "::", "-"), ":", "-")
	//serviceName = strings.ReplaceAll(serviceName, ".", "-")
	//if err = s.api.Paas().DeleteConfigMap(ctx, s.namespace, serviceName, fmt.Sprintf("%s-config", serviceName)); err != nil {
	//	_ = level.Error(logger).Log("api.Paas", "DeleteConfigMap", "err", err.Error())
	//	//return
	//}
	// 更新数据库状态
	jobInfo.TrainStatus = types.TrainStatusCancel
	if err = s.store.FineTuning().UpdateFineTuningJob(ctx, &jobInfo); err != nil {
		_ = level.Error(logger).Log("repository.FineTuningJob", "UpdateFineTuningJob", "err", err.Error())
		return
	}
	return
}

func (s *service) UpdateJobFinishedStatus(ctx context.Context, fineTuningJob string, status types.TrainStatus, message string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	s.mu.Lock()
	defer s.mu.Unlock()
	lockKey := fmt.Sprintf("fine-tune:%s", fineTuningJob)
	acquired, err := s.rdb.SetNX(ctx, lockKey, 1, 2*time.Minute).Result()
	if err != nil {
		_ = level.Warn(logger).Log("redis", "SetNX", "err", err.Error())
		return
	}
	_ = level.Info(logger).Log("lockKey", lockKey, "acquired", acquired)
	if acquired {
		defer func() {
			if rdbErr := s.rdb.Del(ctx, lockKey).Err(); rdbErr != nil {
				_ = level.Warn(logger).Log("redis", "Del", "err", rdbErr.Error())
			}
		}()
	}

	// 查看相关数据
	jobInfo, err := s.store.FineTuning().FindFineTuningJobByJobId(ctx, fineTuningJob, "Template", "BaseModelInfo")
	if err != nil {
		_ = level.Error(logger).Log("repository.FineTuningJob", "FindFineTuningJobByJobId", "err", err.Error())
		err = errors.Wrap(err, "repository.FineTuningJob.FindFineTuningJobByJobId")
		return
	}
	// 如果还没结束，更新数据 返回
	if jobInfo.TrainStatus != types.TrainStatusRunning {
		_ = level.Info(logger).Log("msg", "任务还没结束")
		err = errors.New("任务还没结束")
		return
	}

	t := time.Now()
	// 如果已结束，更新数据
	jobInfo.TrainStatus = status
	if status == types.TrainStatusSuccess {
		jobInfo.Progress = 1
		jobInfo.FinishedAt = &t
	}
	if jobInfo.StartTrainTime != nil {
		jobInfo.TrainDuration = int(time.Now().Unix() - jobInfo.StartTrainTime.Unix())
	}
	// message 截取后3000个字符
	if len(message) > 3000 {
		message = message[:3000]
	}
	jobInfo.ErrorMessage = message
	if err = s.store.FineTuning().UpdateFineTuningJob(ctx, &jobInfo); err != nil {
		_ = level.Error(logger).Log("repository.FineTuningJob", "UpdateFineTuningJob", "err", err.Error())
		return errors.Wrap(err, "repository.FineTuningJob.UpdateFineTuningJob")
	}
	// 删除job 和configmap
	// 调用paas强制删除job
	//if err = s.api.Paas().DeleteJob(ctx, s.namespace, jobInfo.PaasJobName); err != nil {
	//	_ = level.Error(logger).Log("api.Paas", "DeleteJob", "err", err.Error())
	//	//_ = s.api.Alarm().Push(ctx, "微调任务删除job失败", fmt.Sprintf("微调任务删除job失败, jobName: %s, err: %s", jobInfo.PaasJobName, err.Error()), "paas-chat-api", alarm.LevelWarning, 5)
	//}
	// 调用paas删除configmap
	//serviceName := strings.ReplaceAll(strings.ReplaceAll(jobInfo.FineTunedModel, "::", "-"), ":", "-")
	//serviceName = strings.ReplaceAll(serviceName, ".", "-")
	//if err = s.api.Paas().DeleteConfigMap(ctx, s.namespace, serviceName, fmt.Sprintf("%s-config", serviceName)); err != nil {
	//	_ = level.Error(logger).Log("api.Paas", "DeleteConfigMap", "err", err.Error())
	//	//_ = s.api.Alarm().Push(ctx, "微调任务删除configmap失败", fmt.Sprintf("微调任务删除configmap失败, serviceName: %s, err: %s", serviceName, err.Error()), "paas-chat-api", alarm.LevelWarning, 5)
	//}

	if status == types.TrainStatusFailed {
		_ = level.Warn(logger).Log("msg", "任务失败")
		return
	}

	model := &types.Models{
		ProviderName: types.ModelProviderLocalAI,
		ModelType:    types.ModelTypeTextGeneration,
		ModelName:    jobInfo.FineTunedModel,
		//BaseModelName: jobInfo.BaseModel,
		MaxTokens:    jobInfo.ModelMaxLength,
		IsPrivate:    true,
		IsFineTuning: true,
		Enabled:      false,
		Remark:       jobInfo.Remark,
		Parameters:   jobInfo.BaseModelInfo.Parameters,
	}

	if err = s.store.Model().CreateModel(ctx, model); err != nil {
		_ = level.Error(logger).Log("repository.Models", "Create", "err", err.Error())
		_ = s.api.Alarm().Push(ctx, "微调任务创建模型失败", fmt.Sprintf("微调任务创建模型失败, jobName: %s, err: %s", jobInfo.PaasJobName, err.Error()), "paas-chat-api", alarm.LevelInfo, 5)
		return errors.Wrap(err, "repository.Models.Create")
	}

	// 将模型授权给租户，如果有channelId的话 同时也授权给渠道
	if jobInfo.ChannelId != 0 {
		if err = s.store.Channel().AddChannelModels(ctx, jobInfo.ChannelId, model); err != nil {
			_ = level.Warn(logger).Log("msg", "update channel info failed", "err", err.Error())
			return errors.Wrap(err, "repository.Channels.UpdateChannel")
		}
	}

	tenantInfo, err := s.store.Tenants().FindTenant(ctx, jobInfo.TenantID)
	if err != nil {
		_ = level.Warn(logger).Log("msg", "find tenant info failed", "err", err.Error())
		return errors.Wrap(err, "repository.Tenants.FindTenant")
	}

	//tenantInfo.Models = append(tenantInfo.Models, *model)
	if err = s.store.Tenants().AddModel(ctx, tenantInfo.ID, model); err != nil {
		_ = level.Warn(logger).Log("msg", "update tenant info failed", "err", err.Error())
		return errors.Wrap(err, "repository.Tenants.UpdateTenant")
	}

	// 数据库获取等待中的微调任务
	nextJobInfo, err := s.store.FineTuning().FindFineTuningJobLastByStatus(ctx, types.TrainStatusWaiting, "id asc")
	if err != nil {
		_ = level.Info(logger).Log("repository.FineTuningJob", "FindFineTuningJobLastByStatus", "err", err.Error())
		return nil
	}
	// 开始启下一个微调任务
	return s._createFineTuningJob(ctx, nextJobInfo.JobId)
}

func (s *service) RunWaitingTrain(ctx context.Context) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	// 数据库获取等待中的微调任务
	jobs, err := s.store.FineTuning().FindFineTuningJobLastByStatus(ctx, types.TrainStatusWaiting, "id asc")
	if err != nil {
		_ = level.Error(logger).Log("repository.FineTuningJob", "FindFineTuningJobLastByStatus", "err", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return
	}
	if jobs.ID == 0 {
		_ = level.Info(logger).Log("msg", "没有等待中的微调任务")
		return
	}
	// 开始启下一个微调任务
	if err = s._createFineTuningJob(ctx, jobs.JobId); err != nil {
		_ = level.Error(logger).Log("service", "_createFineTuningJob", "jobs.JobId", jobs.JobId, "err", err.Error())
		return err
	}
	return
}

func (s *service) _createFineTuningJob(ctx context.Context, jobId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	jobInfo, err := s.store.FineTuning().FindFineTuningJobByJobId(ctx, jobId, "Template", "BaseModelInfo")
	if err != nil {
		_ = level.Error(logger).Log("repository.FineTuningJob", "FindFineTuningJobByJobId", "err", err.Error())
		return
	}
	// 生成模版
	tplContent, err := s.store.FineTuning().EncodeFineTuningJobTemplate(ctx, jobInfo.Template.Content, &jobInfo)
	if err != nil {
		_ = level.Error(logger).Log("repository.FineTuningJob", "EncodeFineTuningJobTemplate", "err", err.Error())
		return errors.Wrap(err, "repository.FineTuningJob.EncodeFineTuningJobTemplate")
	}
	jobInfo.TrainScript = tplContent

	//namespaceName := s.namespace
	serviceName := strings.ReplaceAll(strings.ReplaceAll(jobInfo.FineTunedModel, "::", "-"), ":", "-")
	serviceName = strings.ReplaceAll(serviceName, ".", "-")

	// 创建服务名
	//if err = s.api.Paas().CreateService(ctx, namespaceName, serviceName); err != nil {
	//	_ = level.Warn(logger).Log("api.Paas", "CreateService", "err", err.Error())
	//	//return errors.Wrap(err, "api.Paas.CreateService")
	//}
	// 创建configmap
	//if err = s.api.Paas().CreateConfigMap(ctx, namespaceName, serviceName, fmt.Sprintf("%s-config", serviceName), map[string]string{
	//	"train.sh": tplContent,
	//}); err != nil {
	//	_ = level.Error(logger).Log("api.Paas", "CreateConfigMap", "err", err.Error())
	//	return errors.Wrap(err, "api.Paas.CreateConfigMap")
	//}

	//var memory = 192
	//var cpu = jobInfo.ProcPerNode * 4
	//if jobInfo.BaseModelInfo.Parameters > 33 {
	//	memory = 348
	//} else if jobInfo.BaseModelInfo.Parameters > 20 {
	//	memory = 256
	//} else if jobInfo.BaseModelInfo.Parameters > 10 {
	//	memory = 192
	//} else if jobInfo.BaseModelInfo.Parameters > 5 {
	//	memory = 128
	//} else {
	//	memory = 96
	//}

	var jobName string
	jobName, err = s.api.DockerApi().Create(ctx, serviceName, dockerapi.Config{
		Image: jobInfo.Template.TrainImage,
		Command: []string{
			"/bin/sh",
			"-c",
			"/app/train.sh",
		},
		GPU: jobInfo.ProcPerNode,
		Volumes: []dockerapi.Volume{
			{
				Key:   "./train.sh",
				Value: "/app/",
			}, {
				Key:   s.dataCfsPath,
				Value: "/data/ft-model",
			}, {
				Key:   s.dataCfsPath,
				Value: "/data/base-model",
			},
		},
		ConfigData: map[string]string{
			"./train.sh": tplContent,
		},
	})
	if err != nil {
		err = errors.Wrap(err, "docker api create")
		return
	}

	//if jobName, err = s.api.Paas().CreateJob(ctx, namespaceName, serviceName, jobInfo.Template.TrainImage, jobInfo.ProcPerNode, 0, 0, []string{
	//	"/bin/sh",
	//	"-c",
	//	"/app/train.sh",
	//}, []paas.Volume{
	//	{
	//		ConfigMapKey: []string{"train.sh"}, // configMap
	//		MountPath:    "/app/",
	//		Type:         1,
	//		VolumeName:   fmt.Sprintf("%s-config", serviceName),
	//	},
	//	{
	//		VolumeName: "aigc-data-cfs", // pvc
	//		Type:       2,
	//		MountPath:  "/data/ft-model",
	//		SubPath:    "ft-model",
	//	},
	//	{
	//		VolumeName: "aigc-data-cfs", // pvc
	//		Type:       2,
	//		MountPath:  "/data/base-model",
	//		SubPath:    "base-model",
	//	},
	//}, s.gpuTolerationValue); err != nil {
	//	_ = level.Error(logger).Log("api.Paas", "CreateJob", "err", err.Error())
	//	return errors.Wrap(err, "api.Paas.CreateJob")
	//}
	t := time.Now()
	jobInfo.PaasJobName = jobName
	jobInfo.TrainStatus = types.TrainStatusRunning
	jobInfo.StartTrainTime = &t
	// 更新数据训
	if err = s.store.FineTuning().UpdateFineTuningJob(ctx, &jobInfo); err != nil {
		_ = level.Error(logger).Log("repository.FineTuningJob", "UpdateFineTuningJob", "err", err.Error())
		return errors.Wrap(err, "repository.FineTuningJob.UpdateFineTuningJob")
	}
	// 定时任务去获取各job 的进度 这块就不处理了，定时任务去处理吧

	return
}

func (s *service) Estimate(ctx context.Context, tenantId uint, request CreateJobRequest) (response EstimateResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Estimate")
	model, err := s.store.Model().GetModelByModelName(ctx, request.BaseModel)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, encode.InvalidParams.Wrap(errors.New("模型不存在"))
		}
		return response, encode.ErrSystem.Wrap(errors.New("查询模型失败"))
	}
	file, err := s.store.Files().FindFileByFileId(ctx, request.FileId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, encode.InvalidParams.Wrap(errors.New("文件不存在"))
		}
		return response, encode.ErrSystem.Wrap(errors.New("查询文件失败"))
	}
	if file.Purpose != types.FilePurposeFineTune.String() {
		return response, encode.InvalidParams.Wrap(errors.New("文件类型错误"))
	}
	tokens := float64(file.TokenCount)
	parameters := model.Parameters
	n := 6 * tokens * parameters * math.Pow(10, 9) * float64(request.TrainEpoch)
	d := float64(request.ProcPerNode*request.ProcPerNode) * 4.5 * math.Pow(10, 12)
	_ = level.Info(logger).Log("finetune estimate", model.ModelName, "tokens", tokens, "parameters", parameters, "procPerNode", request.ProcPerNode, "n", n, "d", d)
	seconds := n/d + 1800
	response.Time = util.FormatDuration(seconds, util.PrecisionMinutes)
	return response, nil
}

func (s *service) ListTemplate(ctx context.Context, tenantId uint, request ListTemplateRequest) (response ListTemplateResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ListTemplate")
	var templates []types.FineTuningTemplate
	var total int64
	templates, total, err = s.store.FineTuning().ListFineTuningTemplate(ctx, finetuning.ListFineTuningTemplateRequest{
		Page:     request.Page,
		PageSize: request.PageSize,
	})
	if err != nil {
		_ = level.Error(logger).Log("store.finetuning", "ListFineTuningTemplate", "err", err.Error())
		return
	}
	response.List = make([]Template, 0)
	for _, tpl := range templates {
		response.List = append(response.List, Template{
			Id:            tpl.ID,
			Name:          tpl.Name,
			BaseModel:     tpl.BaseModel,
			Content:       tpl.Content,
			MaxTokens:     tpl.MaxTokens,
			Params:        tpl.Params,
			ScriptFile:    tpl.ScriptFile,
			BaseModelPath: tpl.BaseModelPath,
			OutputDir:     tpl.OutputDir,
			Remark:        tpl.Remark,
			CreatedAt:     tpl.CreatedAt,
			UpdatedAt:     tpl.UpdatedAt,
			TemplateType:  string(tpl.TemplateType),
			TrainImage:    tpl.TrainImage,
		})
	}
	response.Total = total
	return
}

func (s *service) GetJob(ctx context.Context, tenantId uint, jobId string) (response JobResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "GetJob")
	job, err := s.store.FineTuning().FindFineTuningJobByJobId(ctx, jobId)
	if err != nil {
		_ = level.Error(logger).Log("store.finetuning", "FindFineTuningJob", "err", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, encode.InvalidParams.Wrap(errors.New("任务不存在"))
		}
		return response, encode.ErrSystem.Wrap(errors.New("查询任务失败"))
	}
	return convertJob(&job), nil
}

func (s *service) DeleteJob(ctx context.Context, tenantId uint, jobId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "DeleteJob")
	job, err := s.store.FineTuning().FindFineTuningJobByJobId(ctx, jobId)
	if err != nil {
		_ = level.Error(logger).Log("store.finetuning", "FindFineTuningJobByJobId", "err", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return encode.InvalidParams.Wrap(errors.New("任务不存在"))
		}
		return encode.ErrSystem.Wrap(errors.New("查询任务失败"))
	}
	// 判断任务是否可以删除
	if !job.CanDelete() {
		_ = level.Error(logger).Log("store.finetuning", "FindFineTuningJobByJobId", "err", "任务不可删除")
		return encode.Invalid.Wrap(errors.Errorf("任务不可删除, status:%s", job.TrainStatus))
	}
	err = s.store.FineTuning().DeleteFineTuningJob(ctx, job.ID)
	if err != nil {
		_ = level.Error(logger).Log("store.finetuning", "DeleteFineTuningJob", "err", err.Error())
		return
	}
	return
}

func (s *service) DashBoard(ctx context.Context, tenantId uint) (res DashBoardResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "DashBoard")
	duration, err := s.store.FineTuning().CountFineTuningJobDuration(ctx)
	if err != nil {
		_ = level.Error(logger).Log("store.finetuning", "CountFineTuningJobDuration", "err", err.Error())
		return
	}
	statusMap, err := s.store.FineTuning().CountFineTuningJobByStatus(ctx)
	if err != nil {
		_ = level.Error(logger).Log("store.finetuning", "CountFineTuningJobByStatus", "err", err.Error())
		return
	}
	res = DashBoardResponse{
		WaitingJobCount:    statusMap[types.TrainStatusWaiting.String()],
		SuccessJobCount:    statusMap[types.TrainStatusSuccess.String()],
		TotalDurationCount: util.FormatDuration(float64(duration), util.PrecisionMinutes),
	}
	return
}

func (s *service) CreateJob(ctx context.Context, tenantId uint, request CreateJobRequest) (response JobResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "CreateJob")
	fileInfo, err := s.store.Files().FindFileByFileId(ctx, request.FileId)
	if err != nil {
		_ = level.Error(logger).Log("repository.finetuning", "FindFileByFileId", "err", err.Error(), "fileId", request.FileId)
		return response, err
	}
	// 转换文件格式
	panUrl, err := s._fileConvertAlpaca(ctx, request.BaseModel, fileInfo.S3Url)
	if err != nil {
		_ = level.Error(logger).Log("service", "_fileConvertAlpaca", "err", err.Error(), "fileId", request.FileId, "s3Url", fileInfo.S3Url)
		return
	}
	ftJobTpl, err := s.store.FineTuning().FindFineTuningTemplateByModel(ctx, request.BaseModel)
	if err != nil {
		_ = level.Error(logger).Log("repository.finetuning", "FindFineTuningTemplateByModel", "err", err.Error(), "baseModel", request.BaseModel)
		return response, err
	}
	suffix := request.Suffix
	// 生成微调任务
	if !strings.EqualFold(request.Suffix, "") {
		suffix = ":" + suffix
	}
	suffix = string(util.Krand(4, util.KC_RAND_KIND_LOWER)) + suffix

	fineTunedModel := fmt.Sprintf("ft::%s:%d-%s", request.BaseModel, request.TenantId, suffix)
	ftJob := &types.FineTuningTrainJob{
		JobId:             uuid.New().String(),
		FileId:            request.FileId,
		ChannelId:         0,
		TemplateId:        ftJobTpl.ID,
		BaseModel:         request.BaseModel,
		TrainEpoch:        request.TrainEpoch,
		BaseModelPath:     ftJobTpl.BaseModelPath,
		DataPath:          fmt.Sprintf("/data/train-data/%s", request.FileId),
		OutputDir:         fmt.Sprintf("%s/ft-%s-%d-%s", ftJobTpl.OutputDir, request.BaseModel, request.TenantId, strings.ReplaceAll(suffix, ":", "-")),
		ScriptFile:        ftJobTpl.ScriptFile,
		MasterPort:        rand.IntnRange(20000, 30000),
		FileUrl:           panUrl,
		TrainStatus:       types.TrainStatusWaiting,
		LearningRate:      request.LearningRate,
		FineTunedModel:    fineTunedModel,
		ProcPerNode:       request.ProcPerNode,
		AccumulationSteps: request.AccumulationSteps,
		TrainBatchSize:    request.TrainBatchSize,
		EvalBatchSize:     request.EvalBatchSize,
		TenantID:          request.TenantId,
		Remark:            request.Remark,
		TrainPublisher:    request.TrainPublisher,
		Lora:              request.Lora,
	}
	err = s.store.FineTuning().CreateFineTuningJob(ctx, ftJob)
	if err != nil {
		_ = level.Error(logger).Log("store.finetuning", "CreateFineTuningJob", "err", err.Error())
		return response, err
	}
	return convertJob(ftJob), nil
}

func (s *service) ListJob(ctx context.Context, tenantId uint, request ListJobRequest) (response ListJobResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ListJob")
	var jobs []types.FineTuningTrainJob
	var total int64
	jobs, total, err = s.store.FineTuning().ListFindTuningJob(ctx, finetuning.ListFindTuningJobRequest{
		Page:           request.Page,
		PageSize:       request.PageSize,
		FineTunedModel: request.FineTunedModel,
		TrainStatus:    request.TrainStatus,
	})
	if err != nil {
		_ = level.Error(logger).Log("store.finetuning", "FindFineTuningJobByStatus", "err", err.Error())
		return
	}
	response.List = make([]JobResponse, 0)
	for _, job := range jobs {
		response.List = append(response.List, convertJob(&job))
	}
	response.Total = total
	return
}

func (s *service) CancelJob(ctx context.Context, tenantId uint, jobId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "CancelJob")
	job, err := s.store.FineTuning().FindFineTuningJobByJobId(ctx, jobId)
	if err != nil {
		_ = level.Error(logger).Log("store.finetuning", "FindFineTuningJobByJobId", "err", err.Error())
		return
	}
	// 判断任务是否可以取消
	if !job.CanCancel() {
		_ = level.Error(logger).Log("store.finetuning", "FindFineTuningJobByJobId", "err", "任务不可取消")
		return encode.Invalid.Wrap(errors.Errorf("任务不可取消, status:%s", job.TrainStatus))
	}
	_, err = s._cancelJob(ctx, job.ChannelId, jobId)
	if err != nil {
		_ = level.Error(logger).Log("finetuning", "CancelFineTuningJob", "err", err.Error())
		return encode.ErrSystem.Wrap(errors.New("取消任务失败, 请稍后重试"))
	}
	return
}

func (s *service) _fileConvertAlpaca(ctx context.Context, modelName, sourceS3Url string) (newS3Url string, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	// 需要把 fileInfo.S3Url 内容格式转换成alpaca的那种
	alpacaDada, err := convertAlpaca(sourceS3Url, logger, modelName)
	if err != nil {
		_ = level.Error(logger).Log("convertAlpaca", "convertAlpaca", "err", err.Error())
		return
	}
	_ = level.Info(logger).Log("msg", "alpacaDada", "msg", "转换完成")

	// 将 *bytes.Reader 类型强制转换为 multipart.File 类型
	file := NewFile(alpacaDada) // 将 []byte 转换为 multipart.File

	fileUrl, err := s.fileSvc.UploadLocal(ctx, file, "json")
	if err != nil {
		_ = level.Error(logger).Log("fileSvc", "UploadLocal", "err", err.Error())
		return
	}

	return fileUrl, nil
}

func New(traceId string, logger log.Logger, store repository.Repository, bucketName, s3AccessKey, s3SecretKey string, apiSvc api.Service, rdb redis.UniversalClient, dataCfsPath string) Service {
	return &service{
		traceId:     traceId,
		logger:      logger,
		store:       store,
		namespace:   "aigc",
		api:         apiSvc,
		rdb:         rdb,
		dataCfsPath: dataCfsPath,
	}
}

type messageLine struct {
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type alpacaData struct {
	ID            string                `json:"id"`
	Conversations []alpacaConversations `json:"conversations"`
}

type alpacaConversations struct {
	From  string `json:"from"`
	Value string `json:"value"`
}

func convertAlpaca(httpUrl string, logger log.Logger, modelName string) (alpaca []byte, err error) {
	body, err := getHttpFileBody(httpUrl)
	if err != nil {
		err = errors.Wrap(err, "getHttpFileBody")
		return
	}
	var roleUser = "human"
	var roleAssistant = "gpt"
	if strings.Contains(modelName, "qwen") {
		roleUser = "user"
		roleAssistant = "assistant"
	}
	var alpacaDataList []alpacaData
	dataList := bytes.Split(body, []byte("\n"))
	for i, line := range dataList {
		var inputMsg messageLine
		if err := json.Unmarshal(line, &inputMsg); err != nil {
			_ = level.Error(logger).Log("json", "Unmarshal", "err", err.Error(), "line", string(line))
			continue
		}
		var conversations []alpacaConversations
		for _, msg := range inputMsg.Messages {
			if !util.StringInArray([]string{"user", "assistant"}, msg.Role) {
				continue
			}
			var role = roleUser
			if msg.Role == "assistant" {
				role = roleAssistant
			}
			conversations = append(conversations, alpacaConversations{
				From:  role,
				Value: msg.Content,
			})
		}
		alpacaDataList = append(alpacaDataList, alpacaData{
			ID:            fmt.Sprintf("ft_alpaca_%d", i),
			Conversations: conversations,
		})
	}
	return json.Marshal(alpacaDataList)
}

func getHttpFileBody(url string) (body []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		err = errors.Wrap(err, "http.Get")
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "io.ReadAll")
		return
	}
	return
}

func convertJob(data *types.FineTuningTrainJob) JobResponse {
	resp := JobResponse{
		Id:                data.ID,
		JobId:             data.JobId,
		BaseModel:         data.BaseModel,
		TrainEpoch:        data.TrainEpoch,
		TrainStatus:       string(data.TrainStatus),
		TrainDuration:     util.FormatDuration(float64(data.TrainDuration), util.PrecisionMinutes),
		Process:           data.Progress,
		FineTunedModel:    data.FineTunedModel,
		Remark:            data.Remark,
		CreatedAt:         data.CreatedAt,
		TrainPublisher:    data.TrainPublisher,
		TrainLog:          data.TrainLog,
		ErrorMessage:      data.ErrorMessage,
		Lora:              data.Lora,
		Suffix:            data.Suffix,
		ModelMaxLength:    data.ModelMaxLength,
		TrainBatchSize:    data.TrainBatchSize,
		FileId:            data.FileId,
		FileUrl:           data.FileUrl,
		LearningRate:      fmt.Sprintf("%.10f", data.LearningRate),
		EvalBatchSize:     data.EvalBatchSize,
		AccumulationSteps: data.AccumulationSteps,
		ProcPerNode:       data.ProcPerNode,
	}
	if data.FinishedAt != nil {
		resp.FinishedAt = data.FinishedAt.Format(time.RFC3339)
	}
	if data.StartTrainTime != nil {
		resp.StartTrainTime = data.StartTrainTime.Format(time.RFC3339)
	}

	if data.TrainStatus == types.TrainStatusRunning && data.StartTrainTime != nil {
		resp.TrainDuration = util.FormatDuration(float64(time.Now().Unix()-data.StartTrainTime.Unix()), util.PrecisionMinutes)
	}

	resp.TrainAnalysis = TrainAnalysis{
		Epoch:        TrainAnalysisObject{List: make([]TrainAnalysisDetail, 0)},
		Loss:         TrainAnalysisObject{List: make([]TrainAnalysisDetail, 0)},
		LearningRate: TrainAnalysisObject{List: make([]TrainAnalysisDetail, 0)},
	}
	if data.TrainLog != "" {
		ana, err := GetTrainInfoFromLog(data.TrainLog)
		if err == nil && len(ana) > 0 {
			for _, item := range ana {
				resp.TrainAnalysis.Epoch.List = append(resp.TrainAnalysis.Epoch.List, TrainAnalysisDetail{
					Timestamp: item.Timestamp,
					Value:     fmt.Sprintf("%.10f", item.Epoch),
				})
				resp.TrainAnalysis.Loss.List = append(resp.TrainAnalysis.Loss.List, TrainAnalysisDetail{
					Timestamp: item.Timestamp,
					Value:     fmt.Sprintf("%.10f", item.Loss),
				})
				resp.TrainAnalysis.LearningRate.List = append(resp.TrainAnalysis.LearningRate.List, TrainAnalysisDetail{
					Timestamp: item.Timestamp,
					Value:     fmt.Sprintf("%.10f", item.LearningRate),
				})
			}
		}
	}
	return resp
}

// File 实现 multipart.File 接口所需的方法
type File struct {
	*bytes.Reader
}

func (f *File) Close() error {
	return nil // bytes.Reader 不需要关闭资源，所以这里返回 nil 即可
}

// NewFile 创建一个新的 File 实例，该实例满足 multipart.File 接口
func NewFile(data []byte) *File {
	return &File{
		bytes.NewReader(data),
	}
}

// GetTrainInfoFromLog 从训练日志获取训练信息
func GetTrainInfoFromLog(jobLog string) (logEntryList []LogEntry, err error) {
	lineArr := strings.Split(jobLog, "\n")
	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z) (\{.*?\})`)

	for _, l := range lineArr {
		matches := re.FindStringSubmatch(l)
		if len(matches) == 3 {
			timestampStr, jsonStr := matches[1], matches[2]

			timestamp, err := time.Parse(time.RFC3339Nano, timestampStr)
			if err != nil {
				continue
			}

			jsonStr = strings.Replace(jsonStr, "'", "\"", -1)        // 将单引号替换为双引号
			jsonStr = strings.Replace(jsonStr, "False", "false", -1) // 将 False 替换为 false
			jsonStr = strings.Replace(jsonStr, "True", "true", -1)   // 将 True 替换为 true

			var entry LogEntry
			err = json.Unmarshal([]byte(jsonStr), &entry)
			if err != nil {
				continue
			}
			entry.Timestamp = timestamp
			logEntryList = append(logEntryList, entry)
		}
	}
	if len(logEntryList) < 1 {
		return
	}
	return
}
