package finetuning

import (
	"context"
	"fmt"
	"github.com/go-kit/log"
	"time"
)

type Middleware func(Service) Service

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (l *logging) Estimate(ctx context.Context, tenantId uint, request CreateJobRequest) (response EstimateResponse, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Estimate",
			"request", fmt.Sprintf("%+v", request),
			"response", response,
			"tenantId", tenantId,
			"err", err,
		)
	}(time.Now())
	return l.next.Estimate(ctx, tenantId, request)
}

func (l *logging) ListTemplate(ctx context.Context, tenantId uint, request ListTemplateRequest) (response ListTemplateResponse, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "ListTemplate",
			"request", fmt.Sprintf("%+v", request),
			"tenantId", tenantId,
			"response", response,
			"err", err,
		)
	}(time.Now())
	return l.next.ListTemplate(ctx, tenantId, request)
}

func (l *logging) GetJob(ctx context.Context, tenantId uint, jobId string) (response JobResponse, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "GetJob",
			"jobId", jobId,
			"tenantId", tenantId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.GetJob(ctx, tenantId, jobId)
}

func (l *logging) DeleteJob(ctx context.Context, tenantId uint, jobId string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "DeleteJob",
			"jobId", jobId,
			"tenantId", tenantId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.DeleteJob(ctx, tenantId, jobId)
}

func (l *logging) DashBoard(ctx context.Context, tenantId uint) (res DashBoardResponse, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "DashBoard",
			"tenantId", tenantId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.DashBoard(ctx, tenantId)
}

func (l *logging) CreateJob(ctx context.Context, tenantId uint, request CreateJobRequest) (response JobResponse, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CreateJob",
			"tenantId", tenantId,
			"request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.CreateJob(ctx, tenantId, request)
}

func (l *logging) ListJob(ctx context.Context, tenantId uint, request ListJobRequest) (response ListJobResponse, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "ListJob",
			"tenantId", tenantId,
			"request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.ListJob(ctx, tenantId, request)
}

func (l *logging) CancelJob(ctx context.Context, tenantId uint, jobId string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CancelJob",
			"tenantId", tenantId,
			"jobId", jobId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.CancelJob(ctx, tenantId, jobId)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "pkg.finetuning", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  logger,
			next:    next,
			traceId: traceId,
		}
	}
}
