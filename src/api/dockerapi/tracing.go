// Code generated . DO NOT EDIT.
package dockerapi

import (
	"context"
	"encoding/json"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) Create(ctx context.Context, name string, config Config) (id string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Create", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.dockerapi",
	})
	defer func() {

		configByte, _ := json.Marshal(config)
		configJson := string(configByte)

		span.LogKV(
			"name", name, "config", configJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.Create(ctx, name, config)

}

func (s *tracing) Logs(ctx context.Context, id string) (log string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Logs", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.dockerapi",
	})
	defer func() {

		span.LogKV(
			"id", id,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.Logs(ctx, id)

}

func (s *tracing) Remove(ctx context.Context, id string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Remove", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.dockerapi",
	})
	defer func() {

		span.LogKV(
			"id", id,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.Remove(ctx, id)

}

func (s *tracing) Restart(ctx context.Context, id string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Restart", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.dockerapi",
	})
	defer func() {

		span.LogKV(
			"id", id,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.Restart(ctx, id)

}

func (s *tracing) Status(ctx context.Context, id string) (status string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Status", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.dockerapi",
	})
	defer func() {

		span.LogKV(
			"id", id,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.Status(ctx, id)

}

func (s *tracing) Stop(ctx context.Context, id string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Stop", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.dockerapi",
	})
	defer func() {

		span.LogKV(
			"id", id,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.Stop(ctx, id)

}

func (s *tracing) Update(ctx context.Context, name string, id string, config Config) (newId string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Update", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.dockerapi",
	})
	defer func() {

		configByte, _ := json.Marshal(config)
		configJson := string(configByte)

		span.LogKV(
			"name", name, "id", id, "config", configJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.Update(ctx, name, id, config)

}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
