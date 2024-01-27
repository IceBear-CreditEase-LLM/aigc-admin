package models

import (
	"context"
	"fmt"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api/paaschat"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/encode"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/model"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/types"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/util"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"strings"
	"time"
)

type Middleware func(Service) Service

type Service interface {
	// ListModels 模型分页列表
	ListModels(ctx context.Context, request ListModelRequest) (res ListModelResponse, err error)
	// CreateModel 创建模型
	CreateModel(ctx context.Context, request CreateModelRequest) (res Model, err error)
	// GetModel 获取模型
	GetModel(ctx context.Context, id uint) (res Model, err error)
	// UpdateModel 更新模型
	UpdateModel(ctx context.Context, request UpdateModelRequest) (err error)
	// DeleteModel 删除模型
	DeleteModel(ctx context.Context, id uint) (err error)
	// Deploy 模型部署
	Deploy(ctx context.Context, request ModelDeployRequest) (err error)
	// Undeploy 模型取消部署
	Undeploy(ctx context.Context, id uint) (err error)
	// CreateEval 创建评估任务
	CreateEval(ctx context.Context, request CreateEvalRequest) (res Eval, err error)
	// ListEval 评估任务分页列表
	ListEval(ctx context.Context, request ListEvalRequest) (res ListEvalResponse, err error)
	// CancelEval 取消评估任务
	CancelEval(ctx context.Context, id uint) (err error)
	// DeleteEval 删除评估任务
	DeleteEval(ctx context.Context, id uint) (err error)
}

type service struct {
	logger  log.Logger
	traceId string
	store   repository.Repository
	apiSvc  api.Service
}

func (s *service) DeleteEval(ctx context.Context, id uint) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "DeleteEval")
	eval, err := s.store.Model().GetEval(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "GetEval", "err", err.Error(), "id", id)
		return encode.ErrSystem.Wrap(errors.New("查询评估任务失败"))
	}

	if utils.Contains([]string{types.EvalStatusPending.String(), types.EvalStatusRunning.String()}, eval.Status.String()) {
		_ = level.Error(logger).Log("msg", "等待中和运行中的任务不可删除", "status", eval.Status.String())
		return encode.InvalidParams.Wrap(errors.New("等待中和运行中的任务不可删除"))
	}
	err = s.store.Model().DeleteEval(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "DeleteEval", "err", err.Error(), "id", id)
		return encode.ErrSystem.Wrap(errors.New("删除评估任务失败"))
	}
	return
}

func (s *service) CancelEval(ctx context.Context, id uint) (err error) {
	eval, err := s.store.Model().GetEval(ctx, id)
	if err != nil {
		return encode.ErrSystem.Wrap(errors.New("查询模型失败"))
	}
	eval.Status = types.EvalStatusCancel
	err = s.store.Model().UpdateEval(ctx, &eval)
	if err != nil {
		return encode.ErrSystem.Wrap(errors.New("取消评估任务失败"))
	}
	return
}

func (s *service) ListEval(ctx context.Context, request ListEvalRequest) (res ListEvalResponse, err error) {
	eval, total, err := s.store.Model().ListEval(ctx, model.ListEvalRequest{
		Page:        request.Page,
		PageSize:    request.PageSize,
		ModelName:   request.ModelName,
		MetricName:  request.MetricName,
		Status:      request.Status,
		DatasetType: request.DatasetType,
	})
	if err != nil {
		return res, encode.ErrSystem.Wrap(errors.New("查询评估任务失败"))
	}
	list := make([]Eval, 0)
	for _, v := range eval {
		list = append(list, convertEval(&v))
	}
	res = ListEvalResponse{
		Total: total,
		List:  list,
	}
	return
}

func (s *service) CreateEval(ctx context.Context, request CreateEvalRequest) (res Eval, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "CreateEval")
	if request.DatasetType == types.EvalDataSetTypeTrain.String() {
		request.EvalPercent = request.EvalPercent / 100.0
		if request.EvalPercent <= 0 || request.EvalPercent > 1 {
			_ = level.Error(logger).Log("msg", "评估比例不正确", "evalPercent", request.EvalPercent)
			return res, encode.InvalidParams.Wrap(errors.New("评估比例不正确"))
		}
		job, err := s.store.FineTuning().GetFineTuningJobByModelName(ctx, request.ModelName)
		if err != nil {
			_ = level.Error(logger).Log("store.FineTuning", "GetFineTuningJobByModelName", "err", err.Error())
			return res, encode.ErrSystem.Wrap(errors.New("查询微调任务失败"))
		}
		request.DatasetFileId = job.FileId
	}

	file, err := s.store.Files().FindFileByFileId(ctx, request.DatasetFileId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, encode.InvalidParams.Wrap(errors.New("文件不存在"))
		}
		_ = level.Error(logger).Log("store.Files", "FindFileByFileId", "err", err.Error())
		return res, encode.ErrSystem.Wrap(errors.New("查询文件失败"))
	}

	if request.DatasetType == types.EvalDataSetTypeTrain.String() && file.Purpose != types.FilePurposeFineTune.String() {
		_ = level.Error(logger).Log("msg", "文件用途不正确", "purpose", file.Purpose)
		return res, encode.InvalidParams.Wrap(errors.New("文件用途不正确"))
	}

	if request.DatasetType == types.EvalDataSetTypeCustom.String() && file.Purpose != types.FilePurposeFineTuneEval.String() {
		_ = level.Error(logger).Log("msg", "文件用途不正确", "purpose", file.Purpose)
		return res, encode.InvalidParams.Wrap(errors.New("文件用途不正确"))
	}

	req := types.LLMEvalResults{
		Status:        types.EvalStatusPending,
		ModelName:     request.ModelName,
		MetricName:    request.MetricName,
		DatasetFileId: file.ID,
		DatasetType:   request.DatasetType,
		Remark:        request.Remark,
		UUid:          uuid.NewString(),
	}
	if request.DatasetType == types.EvalDataSetTypeTrain.String() {
		req.EvalTotal = int(float64(file.LineCount) * request.EvalPercent)
	}
	err = s.store.Model().CreateEval(ctx, &req)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "CreateEval", "err", err.Error())
		return
	}
	res = convertEval(&req)
	return
}

func (s *service) Undeploy(ctx context.Context, id uint) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Undeploy")
	m, err := s.store.Model().GetModel(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "GetModel", "err", err.Error(), "id", id)
		return
	}
	if !m.IsPrivate {
		// 公有模型不需要部署
		_ = level.Warn(logger).Log("msg", "public model not need undeploy", "modelName", m.ModelName)
		return encode.Invalid.Wrap(errors.Errorf("public model not need undeploy, model:%s", m.ModelName))
	}

	// 调用API取消部署模型
	err = s.apiSvc.PaasChat().UndeployModel(ctx, m.ModelName)
	if err != nil {
		_ = level.Error(logger).Log("api.PaasChat", "UndeployModel", "err", err.Error(), "modelName", m.ModelName)
		return
	}
	return
}

func (s *service) Deploy(ctx context.Context, request ModelDeployRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Deploy")
	m, err := s.store.Model().GetModel(ctx, request.Id)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "GetModel", "err", err.Error(), "request", fmt.Sprintf("%+v", request))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return encode.InvalidParams.Wrap(errors.New("模型不存在"))
		}
		return encode.ErrSystem.Wrap(errors.New("查询模型异常"))
	}
	if !m.IsPrivate {
		// 公有模型不需要部署
		_ = level.Warn(logger).Log("msg", "public model not need deploy", "modelName", m.ModelName)
		return encode.Invalid.Wrap(errors.Errorf("public model not need deploy, model:%s", m.ModelName))
	}
	if m.Enabled {
		// 已经被部署
		_ = level.Warn(logger).Log("msg", "model already deployed", "modelName", m.ModelName)
		return errors.Errorf("模型%s已部署", m.ModelName)
	}
	// 调用API部署模型
	req := paaschat.DeployModelRequest{
		ModelName:    m.ModelName,
		Replicas:     request.Replicas,
		Label:        request.Label,
		Gpu:          request.Gpu,
		Quantization: request.Quantization,
		Vllm:         request.Vllm,
	}
	err = s.apiSvc.PaasChat().DeployModel(ctx, req)
	if err != nil {
		_ = level.Error(logger).Log("api.PaasChat", "DeployModel", "err", err.Error(), "modelName", m.ModelName)
		return encode.ErrSystem.Wrap(errors.New("部署模型异常" + err.Error()))
	}
	return
}

func (s *service) ListModels(ctx context.Context, request ListModelRequest) (res ListModelResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ListModels")
	req := model.ListModelRequest{
		Page:         request.Page,
		PageSize:     request.PageSize,
		Enabled:      request.Enabled,
		IsPrivate:    request.IsPrivate,
		IsFineTuning: request.IsFineTuning,
		ModelName:    request.ModelName,
		ProviderName: request.ProviderName,
	}
	models, total, err := s.store.Model().ListModels(ctx, req)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "ListModels", "err", err.Error())
		return res, encode.ErrSystem.Wrap(errors.New("查询模型列表失败"))
	}
	list := make([]Model, 0)
	for _, v := range models {
		list = append(list, convert(&v))
	}
	res = ListModelResponse{
		Total:  total,
		Models: list,
	}
	return
}

func (s *service) CreateModel(ctx context.Context, request CreateModelRequest) (res Model, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "CreateModel")
	m, err := s.store.Model().GetModelByModelName(ctx, request.ModelName)
	if err == nil && m.ID > 0 {
		_ = level.Warn(logger).Log("store.Model", "GetModelByModelName", "err", "模型名称已存在", "modelName", request.ModelName)
		return res, encode.InvalidParams.Wrap(errors.Errorf("%s 模型已存在", request.ModelName))
	}
	provider := request.ProviderName
	if provider == "" {
		provider = providerName(request.ModelName).String()
	}
	req := types.Models{
		ProviderName: provider,
		ModelType:    types.ModelTypeTextGeneration.String(),
		ModelName:    request.ModelName,
		MaxTokens:    request.MaxTokens,
		IsPrivate:    request.IsPrivate,
		IsFineTuning: request.IsFineTuning,
		Enabled:      request.Enabled,
		Remark:       request.Remark,
		TenantId:     request.TenantId,
		Parameters:   request.Parameters,
		LastOperator: request.Email,
	}
	err = s.store.Model().CreateModel(ctx, &req)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "CreateModel", "err", err.Error())
		return
	}
	res = convert(&req)
	return
}

func (s *service) GetModel(ctx context.Context, id uint) (res Model, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "GetModel")
	m, err := s.store.Model().GetModel(ctx, id, "FineTuningTrainJob", "ModelDeploy")
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "GetModel", "err", err.Error(), "id", id)
		return
	}
	res = convert(&m)
	return
}

func (s *service) UpdateModel(ctx context.Context, request UpdateModelRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "UpdateModel")
	var req model.UpdateModelRequest
	req.Id = request.Id
	if request.TenantId != nil {
		req.TenantId = request.TenantId
	}
	if request.MaxTokens != nil {
		req.MaxTokens = request.MaxTokens
	}
	if request.Enabled != nil {
		req.Enabled = request.Enabled
	}
	if request.Remark != nil {
		req.Remark = request.Remark
	}
	err = s.store.Model().UpdateModel(ctx, req)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "UpdateModel", "err", err.Error())
		return
	}
	return
}

func (s *service) DeleteModel(ctx context.Context, id uint) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "DeleteModel")
	// todo 删除之前判断是否有绑定的渠道
	err = s.store.Model().DeleteModel(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "DeleteModel", "err", err.Error())
		return
	}
	return
}

func convert(data *types.Models) Model {
	m := Model{
		Id:           data.ID,
		ProviderName: data.ProviderName,
		ModelType:    data.ModelType,
		ModelName:    data.ModelName,
		MaxTokens:    data.MaxTokens,
		IsPrivate:    data.IsPrivate,
		IsFineTuning: data.IsFineTuning,
		Enabled:      data.Enabled,
		Remark:       data.Remark,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
		DeployStatus: data.ModelDeploy.Status,
		Parameters:   data.Parameters,
		LastOperator: data.LastOperator,
	}
	tenants := make([]Tenant, 0)
	for _, t := range data.Tenants {
		tenants = append(tenants, Tenant{
			Id:   t.ID,
			Name: t.Name,
		})
	}
	m.Tenants = tenants
	operation := make([]string, 0)
	operation = append(operation, "edit")
	if data.CanDelete() {
		operation = append(operation, "delete")
	}
	if data.CanDeploy() {
		operation = append(operation, "deploy")
	}
	if data.CanUndeploy() {
		operation = append(operation, "undeploy")
	}
	m.Operation = operation
	if data.FineTuningTrainJob.ID > 0 {
		m.JobId = data.FineTuningTrainJob.JobId
	}
	return m
}

func providerName(m string) types.ModelProvider {
	openAIModels := []string{
		"gpt-4-1106-preview",
		"gpt-4-vision-preview",
		"gpt-4",
		"gpt-4-32k",
		"gpt-4-0613",
		"gpt-4-32k-0613",
		"gpt-4-0314",
		"gpt-4-32k-0314",
		"gpt-3.5-turbo-1106",
		"gpt-3.5-turbo",
		"gpt-3.5-turbo-16k",
		"gpt-3.5-turbo-instruct",
		"gpt-3.5-turbo-0613",
		"gpt-3.5-turbo-16k-0613",
		"gpt-3.5-turbo-0301",
		"text-davinci-003",
		"text-davinci-002",
		"code-davinci-002",
	}
	for _, v := range openAIModels {
		if strings.EqualFold(v, m) {
			return types.ModelProviderOpenAI
		}
	}
	return types.ModelProviderLocalAI
}

func NewService(logger log.Logger, traceId string, store repository.Repository, apiSvc api.Service) Service {
	return &service{
		logger:  log.With(logger, "service", "models"),
		traceId: traceId,
		store:   store,
		apiSvc:  apiSvc,
	}
}

func convertEval(data *types.LLMEvalResults) Eval {
	e := Eval{
		Id:          data.ID,
		ModelName:   data.ModelName,
		DatasetType: data.DatasetType,
		Progress:    data.Progress,
		Score:       data.Score,
		CreatedAt:   data.CreatedAt,
		Status:      data.Status.String(),
		EvalTotal:   data.EvalTotal,
		Remark:      data.Remark,
		MetricName:  data.MetricName,
	}
	if data.Status == types.EvalStatusRunning && data.StartedAt != nil {
		start := data.StartedAt.Time
		e.Duration = util.FormatDuration(float64(time.Now().Sub(start)), util.PrecisionMinutes)
	}
	if data.Status == types.EvalStatusSuccess && data.CompletedAt != nil && data.StartedAt != nil {
		start := data.StartedAt.Time
		end := data.CompletedAt.Time
		e.Duration = util.FormatDuration(float64(end.Sub(start)), util.PrecisionMinutes)
	}
	if data.StartedAt != nil {
		e.StartedAt = data.StartedAt.Time.Format(time.RFC3339)
	}
	return e
}
