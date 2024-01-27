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

func (t *tracing) DeleteEval(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "DeleteEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.DeleteEval(ctx, id)
}

func (t *tracing) CreateEval(ctx context.Context, data *types.LLMEvalResults) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CreateEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("data", fmt.Sprintf("%+v", data), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.CreateEval(ctx, data)
}

func (t *tracing) ListEval(ctx context.Context, request ListEvalRequest) (res []types.LLMEvalResults, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.ListEval(ctx, request)
}

func (t *tracing) UpdateEval(ctx context.Context, data *types.LLMEvalResults) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "UpdateEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("data", fmt.Sprintf("%+v", data), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.UpdateEval(ctx, data)
}

func (t *tracing) GetEval(ctx context.Context, id uint) (res types.LLMEvalResults, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "GetEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.GetEval(ctx, id)
}

func (t *tracing) ListModels(ctx context.Context, request ListModelRequest) (res []types.Models, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListModels", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.ListModels(ctx, request)
}

func (t *tracing) CreateModel(ctx context.Context, data *types.Models) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CreateModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.CreateModel(ctx, data)
}

func (t *tracing) GetModel(ctx context.Context, id uint, preload ...string) (res types.Models, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "GetModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.GetModel(ctx, id, preload...)
}

func (t *tracing) UpdateModel(ctx context.Context, request UpdateModelRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "UpdateModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.UpdateModel(ctx, request)
}

func (t *tracing) DeleteModel(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "DeleteModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.DeleteModel(ctx, id)
}

func (t *tracing) FindModelsByTenantId(ctx context.Context, tenantId uint) (res []types.Models, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "FindModelsByTenantId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.FindModelsByTenantId(ctx, tenantId)
}

func (t *tracing) GetModelByModelName(ctx context.Context, modelName string) (res types.Models, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "GetModelByModelName", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("modelName", modelName, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.GetModelByModelName(ctx, modelName)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
