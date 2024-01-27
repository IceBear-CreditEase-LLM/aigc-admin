package langchain

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

func (s *logging) Summary(ctx context.Context, chainType ChainType, prompt, filePath, modelName string, maxTokens int, temperature, topP float64, streamingFunc func(ctx context.Context, chunk []byte) error) (res map[string]any, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Summary", "chainType", chainType, "prompt", prompt, "filePath", filePath, "modelName", modelName, "maxTokens", maxTokens, "temperature", temperature, "topP", topP, "streamingFunc", streamingFunc,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Summary(ctx, chainType, prompt, filePath, modelName, maxTokens, temperature, topP, streamingFunc)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "pkg.langchain", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
