package paaschat

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sashabaranov/go-openai"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (t *tracing) Wav2lipSynthesisCheck(ctx context.Context) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.paaschat",
	})
	defer func() {
		span.LogKV("err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.Wav2lipSynthesisCheck(ctx)
}

func (t *tracing) Wav2lipSynthesisCancel(ctx context.Context, uuid string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.paaschat",
	})
	defer func() {
		span.LogKV("uuid", uuid, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return (t.next).Wav2lipSynthesisCancel(ctx, uuid)
}

func (t *tracing) CancelFineTuningJob(ctx context.Context, jobId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CancelFineTuningJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.paasChat",
	})
	defer func() {
		span.LogKV("jobId", jobId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.CancelFineTuningJob(ctx, jobId)
}

func (t *tracing) ChatCompletionStream(ctx context.Context, request openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ChatCompletionStream", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.paasChat",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.ChatCompletionStream(ctx, request)
}

func (t *tracing) DeployModel(ctx context.Context, request DeployModelRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "DeployModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.paasChat",
	})
	defer func() {
		span.LogKV("request", request, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.DeployModel(ctx, request)
}

func (t *tracing) UndeployModel(ctx context.Context, modelName string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "UndeployModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.paasChat",
	})
	defer func() {
		span.LogKV("modelName", modelName, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.UndeployModel(ctx, modelName)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
