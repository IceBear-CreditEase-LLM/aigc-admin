package finetuning

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (t *tracing) Estimate(ctx context.Context, tenantId uint, request CreateJobRequest) (response EstimateResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "Estimate", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "request", fmt.Sprintf("%+v", request), "response", response, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.Estimate(ctx, tenantId, request)
}

func (t *tracing) ListTemplate(ctx context.Context, tenantId uint, request ListTemplateRequest) (response ListTemplateResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListTemplate", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "request", fmt.Sprintf("%+v", request), "response", response, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.ListTemplate(ctx, tenantId, request)
}

func (t *tracing) GetJob(ctx context.Context, tenantId uint, jobId string) (response JobResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "GetJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "jobId", jobId, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.GetJob(ctx, tenantId, jobId)
}

func (t *tracing) DeleteJob(ctx context.Context, tenantId uint, jobId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "DeleteJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("jobId", jobId, "tenantId", tenantId, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.DeleteJob(ctx, tenantId, jobId)
}

func (t *tracing) DashBoard(ctx context.Context, tenantId uint) (res DashBoardResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "DashBoard", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.DashBoard(ctx, tenantId)
}

func (t *tracing) CreateJob(ctx context.Context, tenantId uint, request CreateJobRequest) (response JobResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CreateJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "tenantId", tenantId, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.CreateJob(ctx, tenantId, request)
}

func (t *tracing) ListJob(ctx context.Context, tenantId uint, request ListJobRequest) (response ListJobResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "tenantId", tenantId, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.ListJob(ctx, tenantId, request)
}

func (t *tracing) CancelJob(ctx context.Context, tenantId uint, jobId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CancelJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("jobId", jobId, "tenantId", tenantId, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.CancelJob(ctx, tenantId, jobId)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
