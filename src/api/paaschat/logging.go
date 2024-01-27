package paaschat

import (
	"context"
	"fmt"
	"github.com/go-kit/log"
	"github.com/sashabaranov/go-openai"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (l *logging) Wav2lipSynthesisCheck(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return (l.next).Wav2lipSynthesisCheck(ctx)
}

func (l *logging) Wav2lipSynthesisCancel(ctx context.Context, uuid string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "", "uuid", uuid,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Wav2lipSynthesisCancel(ctx, uuid)
}

func (l *logging) CancelFineTuningJob(ctx context.Context, jobId string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CancelFineTuningJob",
			"request", jobId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.CancelFineTuningJob(ctx, jobId)
}

func (l *logging) DeployModel(ctx context.Context, request DeployModelRequest) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "DeployModel",
			"request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.DeployModel(ctx, request)
}

func (l *logging) UndeployModel(ctx context.Context, modelName string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "UndeployModel",
			"request", modelName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.UndeployModel(ctx, modelName)
}

func (l *logging) ChatCompletionStream(ctx context.Context, request openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "ChatCompletionStream",
			"request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.ChatCompletionStream(ctx, request)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "api.paasChat", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  logger,
			next:    next,
			traceId: traceId,
		}
	}
}
