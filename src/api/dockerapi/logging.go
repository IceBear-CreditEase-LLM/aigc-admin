// Code generated . DO NOT EDIT.
package dockerapi

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) Create(ctx context.Context, name string, config Config) (id string, err error) {
	defer func(begin time.Time) {

		configByte, _ := json.Marshal(config)
		configJson := string(configByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Create",

			"name", name,

			"config", configJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.Create(ctx, name, config)

}

func (s *logging) Logs(ctx context.Context, id string) (log string, err error) {
	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Logs",

			"id", id,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.Logs(ctx, id)

}

func (s *logging) Remove(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Remove",

			"id", id,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.Remove(ctx, id)

}

func (s *logging) Restart(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Restart",

			"id", id,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.Restart(ctx, id)

}

func (s *logging) Status(ctx context.Context, id string) (status string, err error) {
	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Status",

			"id", id,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.Status(ctx, id)

}

func (s *logging) Stop(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Stop",

			"id", id,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.Stop(ctx, id)

}

func (s *logging) Update(ctx context.Context, name string, id string, config Config) (newId string, err error) {
	defer func(begin time.Time) {

		configByte, _ := json.Marshal(config)
		configJson := string(configByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Update",

			"name", name,

			"id", id,

			"config", configJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.Update(ctx, name, id, config)

}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "api.dockerapi", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
