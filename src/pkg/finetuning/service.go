package finetuning

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/encode"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/finetuning"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/types"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/util"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"io"
	"k8s.io/apimachinery/pkg/util/rand"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Service
type Service interface {
	CreateJob(ctx context.Context, tenantId uint, request CreateJobRequest) (response JobResponse, err error)
	ListJob(ctx context.Context, tenantId uint, request ListJobRequest) (response ListJobResponse, err error)
	CancelJob(ctx context.Context, tenantId uint, jobId string) (err error)
	DashBoard(ctx context.Context, tenantId uint) (res DashBoardResponse, err error)
	DeleteJob(ctx context.Context, tenantId uint, jobId string) (err error)
	GetJob(ctx context.Context, tenantId uint, jobId string) (response JobResponse, err error)
	ListTemplate(ctx context.Context, tenantId uint, request ListTemplateRequest) (response ListTemplateResponse, err error)
	// Estimate 微调时间预估
	Estimate(ctx context.Context, tenantId uint, request CreateJobRequest) (response EstimateResponse, err error)
}

type service struct {
	traceId     string
	logger      log.Logger
	store       repository.Repository
	api         api.Service
	namespace   string
	bucketName  string
	s3AccessKey string
	s3SecretKey string
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
	err = s.api.PaasChat().CancelFineTuningJob(ctx, jobId)
	if err != nil {
		_ = level.Error(logger).Log("paasChat", "CancelFineTuningJob", "err", err.Error())
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

	paths := strings.Split(time.Now().Format(time.RFC3339), "-")
	targetPath := fmt.Sprintf("aigc/train-data/%s/%s/%s.json", paths[0], paths[1], util.Krand(12, util.KC_RAND_KIND_ALL))
	err = s.api.S3Client(ctx).Upload(ctx, s.bucketName, targetPath, file, "")
	if err != nil {
		_ = level.Error(logger).Log("s3Client", "Upload", "err", err.Error())
		return
	}
	shareUrl, err := s.api.S3Client(ctx).ShareGen(ctx, s.bucketName, targetPath, 60*24*31*12)
	if err != nil {
		_ = level.Error(logger).Log("pan", "ShareGen", "err", err.Error())
		return
	}

	return shareUrl, nil
}

func New(traceId string, logger log.Logger, store repository.Repository, bucketName, s3AccessKey, s3SecretKey string, apiSvc api.Service) Service {
	return &service{
		traceId:     traceId,
		logger:      logger,
		store:       store,
		bucketName:  bucketName,
		s3AccessKey: s3AccessKey,
		s3SecretKey: s3SecretKey,
		namespace:   "aigc",
		api:         apiSvc,
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
