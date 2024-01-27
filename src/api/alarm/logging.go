/**
 * @Time : 2021/11/4 11:11 AM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package alarm

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

func (l *logging) Push(ctx context.Context, title, content, metrics string, level Level, silencePeriod int) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Push",
			"title", title,
			"content", content,
			"metrics", metrics,
			"level", level,
			"silencePeriod", silencePeriod,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Push(ctx, title, content, metrics, level, silencePeriod)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "api.alarm", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
