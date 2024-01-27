package langchain

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) Summary(ctx context.Context, chainType ChainType, prompt, filePath, modelName string, maxTokens int, temperature, topP float64, streamingFunc func(ctx context.Context, chunk []byte) error) (res map[string]any, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Summary", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.langchain",
	})
	defer func() {
		span.LogKV("chainType", chainType, "prompt", prompt, "filePath", filePath, "modelName", modelName, "maxTokens", maxTokens, "temperature", temperature, "topP", topP, "streamingFunc", streamingFunc, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Summary(ctx, chainType, prompt, filePath, modelName, maxTokens, temperature, topP, streamingFunc)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
