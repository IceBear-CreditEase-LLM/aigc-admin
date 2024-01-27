package channels

import (
	"context"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/channel"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/model"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/types"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/util"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"io"
	"time"
)

type Service interface {
	CreateChannel(ctx context.Context, request CreateChannelRequest) (resp Channel, err error)
	UpdateChannel(ctx context.Context, request UpdateChannelRequest) (err error)
	ListChannel(ctx context.Context, request ListChannelRequest) (resp ChannelList, err error)
	DeleteChannel(ctx context.Context, id uint) (err error)
	GetChannel(ctx context.Context, id uint) (resp Channel, err error)
	ListChannelModels(ctx context.Context, request ListChannelModelsRequest) (resp ChannelModelList, err error)
	ChatCompletionStream(ctx context.Context, request ChatCompletionRequest) (stream <-chan CompletionsStreamResult, err error)
}

type service struct {
	logger  log.Logger
	traceId string
	store   repository.Repository
	apiSvc  api.Service
}

func (s *service) ChatCompletionStream(ctx context.Context, request ChatCompletionRequest) (stream <-chan CompletionsStreamResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ChatCompletionStream")
	completionStream, err := s.apiSvc.PaasChat().ChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:       request.Model,
		Messages:    request.Messages,
		MaxTokens:   request.MaxTokens,
		Temperature: request.Temperature,
		TopP:        request.TopP,
	})
	if err != nil {
		_ = level.Error(logger).Log("apiSvc.PaasChat", "ChatCompletionStream", "err", err.Error())
		return stream, err
	}

	dot := make(chan CompletionsStreamResult)
	go func(completionStream *openai.ChatCompletionStream, dot chan CompletionsStreamResult) {
		var fullContent string
		defer func() {
			completionStream.Close()
			close(dot)
		}()
		for {
			completion, err := completionStream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				_ = level.Error(logger).Log("completionStream", "Recv", "err", err.Error())
				return
			}
			begin := time.Now()
			fullContent += completion.Choices[0].Delta.Content
			dot <- CompletionsStreamResult{
				FullContent: fullContent,
				Content:     completion.Choices[0].Delta.Content,
				CreatedAt:   begin,
				ContentType: "text",
				MessageId:   "",
				Model:       "",
				TopP:        0,
				Temperature: 0,
				MaxTokens:   0,
			}
		}
	}(completionStream, dot)
	return dot, nil
}

func (s *service) ListChannelModels(ctx context.Context, request ListChannelModelsRequest) (resp ChannelModelList, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ListChannelModels")
	if request.TenantId == types.SystemTenant {
		enabled := true
		req := model.ListModelRequest{
			Page:     -1,
			PageSize: -1,
			Enabled:  &enabled,
		}
		res, total, err := s.store.Model().ListModels(ctx, req)
		if err != nil {
			_ = level.Error(logger).Log("store.Model", "ListModels", "err", err.Error())
			return resp, err
		}
		resp.Total = total
		resp.Models = make([]Model, 0)
		for _, v := range res {
			resp.Models = append(resp.Models, convertModel(&v))
		}
		return resp, nil
	}
	res, err := s.store.Model().FindModelsByTenantId(ctx, request.TenantId)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "FindModelsByTenantId", "err", err.Error())
		return resp, err
	}
	resp.Models = make([]Model, 0)
	for _, v := range res {
		if v.Enabled {
			resp.Models = append(resp.Models, convertModel(&v))
		}
	}
	resp.Total = int64(len(resp.Models))
	return
}

func (s *service) GetChannel(ctx context.Context, id uint) (resp Channel, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "GetChannel")
	res, err := s.store.Channel().GetChannel(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "GetChannel", "err", err.Error(), "id", id)
		return
	}
	resp = convert(&res)
	return
}

func (s *service) CreateChannel(ctx context.Context, request CreateChannelRequest) (resp Channel, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "CreateChannel")
	data := &types.ChatChannels{
		Name:         uuid.New().String(),
		Alias:        request.Alias,
		Remark:       request.Remark,
		Quota:        request.Quota,
		ApiKey:       "sk-" + string(util.Krand(48, util.KC_RAND_KIND_ALL)),
		Email:        request.Email,
		LastOperator: request.LastOperator,
		TenantId:     request.TenantId,
		ModelId:      request.ModelId,
	}
	err = s.store.Channel().CreateChannel(ctx, data)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "CreateChannel", "err", err.Error())
		return
	}
	resp = convert(data)
	return
}

func (s *service) UpdateChannel(ctx context.Context, request UpdateChannelRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "UpdateChannel")
	res, err := s.store.Channel().GetChannel(ctx, request.Id)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "FindChannelById", "err", err.Error())
		return
	}
	if request.Name != nil {
		res.Name = *request.Name
	}
	if request.Alias != nil {
		res.Alias = *request.Alias
	}
	if request.Quota != nil {
		res.Quota = *request.Quota
	}
	if request.Email != nil {
		res.Email = *request.Email
	}
	if request.Remark != nil {
		res.Remark = *request.Remark
	}
	res.ModelId = request.ModelId
	res.UpdatedAt = time.Now()
	err = s.store.Channel().UpdateChannel(ctx, &res)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "UpdateChannel", "err", err.Error(), "id", request.Id)
		return
	}
	return
}

func (s *service) ListChannel(ctx context.Context, request ListChannelRequest) (resp ChannelList, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ListChannel")
	resp.Channels = make([]Channel, 0)
	listReq := channel.ListChannelRequest{
		Page:        request.Page,
		PageSize:    request.PageSize,
		Name:        request.Name,
		Email:       request.Email,
		ProjectName: request.ProjectName,
		ServiceName: request.ServiceName,
		TenantId:    request.TenantId,
	}
	res, total, err := s.store.Channel().ListChannels(ctx, listReq)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "ListChannels", "err", err.Error())
		return
	}
	for _, v := range res {
		resp.Channels = append(resp.Channels, convert(&v))
	}
	resp.Total = total
	return
}

func convert(data *types.ChatChannels) Channel {
	c := Channel{
		Id:           data.ID,
		Name:         data.Name,
		Alias:        data.Alias,
		Quota:        data.Quota,
		ApiKey:       data.ApiKey,
		Email:        data.Email,
		Remark:       data.Remark,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
		TenantId:     data.TenantId,
		LastOperator: data.LastOperator,
	}
	models := make([]Model, 0)
	for _, v := range data.ChannelModels {
		models = append(models, Model{
			Id:           v.ID,
			ProviderName: v.ProviderName,
			ModelType:    v.ModelType,
			ModelName:    v.ModelName,
			MaxTokens:    v.MaxTokens,
			IsPrivate:    v.IsPrivate,
			Remark:       v.Remark,
			Enabled:      v.Enabled,
			IsFineTuning: v.IsFineTuning,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
		})
	}
	c.Model.Num = len(models)
	c.Model.List = models
	return c
}

func convertModel(data *types.Models) Model {
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
	}
	return m
}

func (s *service) DeleteChannel(ctx context.Context, id uint) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "DeleteChannel")
	err = s.store.Channel().DeleteChannel(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "DeleteChannel", "err", err.Error(), "id", id)
		return
	}
	return
}

func NewService(logger log.Logger, traceId string, store repository.Repository, apiSvc api.Service) Service {
	return &service{
		logger:  log.With(logger, "service", "channels"),
		traceId: traceId,
		store:   store,
		apiSvc:  apiSvc,
	}
}
