package model

import (
	"context"
	"fmt"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/types"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) CreateDeploy(ctx context.Context, modelDeploy *types.ModelDeploy) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateDeploy", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("modelDeploy", fmt.Sprintf("%+v", modelDeploy), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateDeploy(ctx, modelDeploy)
}

func (s *tracing) DeleteDeploy(ctx context.Context, modelId uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteDeploy", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("modelId", modelId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.DeleteDeploy(ctx, modelId)
}

func (s *tracing) FindByModelId(ctx context.Context, modelId string, preloads ...string) (model types.Models, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByModelId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("modelId", modelId, "model", fmt.Sprintf("%+v", model), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindByModelId(ctx, modelId, preloads...)
}

func (s *tracing) FindDeployPendingModels(ctx context.Context) (models []types.Models, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindDeployPendingModels", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("models", fmt.Sprintf("%+v", models), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindDeployPendingModels(ctx)
}

func (s *tracing) UpdateDeployStatus(ctx context.Context, modelId uint, status types.ModelDeployStatus) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateDeployStatus", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("modelId", modelId, "status", status, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateDeployStatus(ctx, modelId, status)
}

func (s *tracing) SetModelEnabled(ctx context.Context, modelId string, enabled bool) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SetModelEnabled", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("modelId", modelId, "enabled", enabled, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.SetModelEnabled(ctx, modelId, enabled)
}

func (s *tracing) DeleteEval(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.DeleteEval(ctx, id)
}

func (s *tracing) CreateEval(ctx context.Context, data *types.LLMEvalResults) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("data", fmt.Sprintf("%+v", data), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateEval(ctx, data)
}

func (s *tracing) ListEval(ctx context.Context, request ListEvalRequest) (res []types.LLMEvalResults, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ListEval(ctx, request)
}

func (s *tracing) UpdateEval(ctx context.Context, data *types.LLMEvalResults) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("data", fmt.Sprintf("%+v", data), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateEval(ctx, data)
}

func (s *tracing) GetEval(ctx context.Context, id uint) (res types.LLMEvalResults, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetEval(ctx, id)
}

func (s *tracing) ListModels(ctx context.Context, request ListModelRequest) (res []types.Models, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListModels", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ListModels(ctx, request)
}

func (s *tracing) CreateModel(ctx context.Context, data *types.Models) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateModel(ctx, data)
}

func (s *tracing) GetModel(ctx context.Context, id uint, preload ...string) (res types.Models, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetModel(ctx, id, preload...)
}

func (s *tracing) UpdateModel(ctx context.Context, request UpdateModelRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateModel(ctx, request)
}

func (s *tracing) DeleteModel(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.DeleteModel(ctx, id)
}

func (s *tracing) FindModelsByTenantId(ctx context.Context, tenantId uint) (res []types.Models, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindModelsByTenantId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindModelsByTenantId(ctx, tenantId)
}

func (s *tracing) GetModelByModelName(ctx context.Context, modelName string) (res types.Models, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetModelByModelName", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("modelName", modelName, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetModelByModelName(ctx, modelName)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
