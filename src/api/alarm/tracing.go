package alarm

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) Push(ctx context.Context, title, content, metrics string, level Level, silencePeriod int) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Push", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Api.Alarm",
	})
	defer func() {
		span.LogKV(
			"title", title,
			"content", content,
			"metrics", metrics,
			"level", level,
			"silencePeriod", silencePeriod,
			"err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Push(ctx, title, content, metrics, level, silencePeriod)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
