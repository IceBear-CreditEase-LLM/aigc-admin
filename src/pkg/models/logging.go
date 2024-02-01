package models

import (
	"context"
	"fmt"
	"github.com/go-kit/log"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (l *logging) SyncDeployStatus(ctx context.Context, modelId string) error {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "SyncDeployStatus",
			"modelId", modelId,
			"took", time.Since(begin),
		)
	}(time.Now())
	return l.next.SyncDeployStatus(ctx, modelId)
}

func (l *logging) DeleteEval(ctx context.Context, id uint) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "DeleteEval",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.DeleteEval(ctx, id)
}

func (l *logging) CreateEval(ctx context.Context, request CreateEvalRequest) (res Eval, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CreateEval",
			"request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.CreateEval(ctx, request)
}

func (l *logging) ListEval(ctx context.Context, request ListEvalRequest) (res ListEvalResponse, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "ListEval",
			"request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.ListEval(ctx, request)
}

func (l *logging) CancelEval(ctx context.Context, id uint) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CancelEval",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.CancelEval(ctx, id)
}

func (l *logging) Undeploy(ctx context.Context, id uint) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Undeploy",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Undeploy(ctx, id)
}

func (l *logging) ListModels(ctx context.Context, request ListModelRequest) (res ListModelResponse, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "ListModels",
			"request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.ListModels(ctx, request)
}

func (l *logging) CreateModel(ctx context.Context, request CreateModelRequest) (res Model, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CreateModel",
			"request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.CreateModel(ctx, request)
}

func (l *logging) GetModel(ctx context.Context, id uint) (res Model, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "GetModel",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.GetModel(ctx, id)
}

func (l *logging) UpdateModel(ctx context.Context, request UpdateModelRequest) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "UpdateModel",
			"request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.UpdateModel(ctx, request)
}

func (l *logging) DeleteModel(ctx context.Context, id uint) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "DeleteModel",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.DeleteModel(ctx, id)
}

func (l *logging) Deploy(ctx context.Context, request ModelDeployRequest) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Deploy",
			"id", request,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Deploy(ctx, request)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "pkg.models", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  logger,
			next:    next,
			traceId: traceId,
		}
	}
}
