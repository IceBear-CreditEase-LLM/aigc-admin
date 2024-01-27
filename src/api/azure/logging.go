package azure

import (
	"context"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) TTS(ctx context.Context, request TTSRequest) (response string, filePath string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "TTS", "request", request,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.TTS(ctx, request)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "api.azure", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
