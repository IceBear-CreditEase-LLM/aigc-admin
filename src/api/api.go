/**
 * @Time: 2020/12/27 22:06
 * @Author: solacowa@gmail.com
 * @File: api
 * @Software: GoLand
 */

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api/alarm"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api/azure"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api/fastchat"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api/ldapcli"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api/paaschat"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api/s3"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/middleware"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"net/http/httputil"
)

type Config struct {
	Namespace, ServiceName string
	FastChat               fastchat.Config
	Ldap                   ldapcli.Config
	Alarm                  alarm.Config
	S3                     S3
	PaasChat               paaschat.Config
	Azure                  azure.Config
	StorageType            string
}

type S3 struct {
	AccessKey, SecretKey string
	Region, Endpoint     string
	BucketName, Cluster  string
}

type ContextKey string

// Service 所有调用外部服务在API聚合
type Service interface {
	// S3Client s3 客户端
	S3Client(ctx context.Context) s3.Service
	// Alarm 统一告警中心客户端
	Alarm() alarm.Service
	// FastChat FastChat服务API
	FastChat() fastchat.Service
	// Ldap ldap客户端
	Ldap() ldapcli.Service
	// PaasChat paaschat服务
	PaasChat() paaschat.Service
	// Azure 微软服务
	Azure() azure.Service
}

type api struct {
	logger      log.Logger
	s3Client    s3.Service
	alarm       alarm.Service
	traceId     string
	fastChatSvc fastchat.Service
	ldapSvc     ldapcli.Service
	paasChatSvc paaschat.Service
	azure       azure.Service
}

func (s *api) Ldap() ldapcli.Service {
	return s.ldapSvc
}

func (s *api) FastChat() fastchat.Service {
	return s.fastChatSvc
}

func (s *api) Alarm() alarm.Service {
	return s.alarm
}

func (s *api) S3Client(ctx context.Context) s3.Service {
	return s.s3Client
}

func (s *api) PaasChat() paaschat.Service {
	return s.paasChatSvc
}

func (s *api) Azure() azure.Service {
	return s.azure
}

// NewApi 中间件有顺序,在后面的会最先执行
func NewApi(_ context.Context, logger log.Logger, traceId string, debug bool, tracer opentracing.Tracer, cfg *Config, opts []kithttp.ClientOption, rdb redis.UniversalClient) Service {
	logger = log.With(logger, "api", "Api")
	if debug {
		opts = append(opts, kithttp.ClientBefore(func(ctx context.Context, request *http.Request) context.Context {
			dump, _ := httputil.DumpRequest(request, true)
			fmt.Println(string(dump))
			return ctx
		}),
			kithttp.ClientAfter(func(ctx context.Context, response *http.Response) context.Context {
				dump, _ := httputil.DumpResponse(response, true)
				fmt.Println(string(dump))
				return ctx
			}),
		)
	}

	alarmSvc := alarm.New(traceId, cfg.Alarm, opts)
	s3Cli := s3.New(cfg.StorageType, cfg.S3.AccessKey, cfg.S3.SecretKey, cfg.S3.Endpoint, cfg.S3.Region)
	fastChatSvcOpts := opts
	if tracer != nil {
		fastChatSvcOpts = append(opts, kithttp.ClientBefore(middleware.RecordRequestAndBody(tracer, logger, "fastChat")))
	}
	fastChatSvc := fastchat.New(logger, cfg.FastChat, fastChatSvcOpts, rdb, alarmSvc)
	ldapSvc := ldapcli.New(cfg.Ldap)
	paasChatSvc := paaschat.New(logger, cfg.PaasChat, opts)
	azureSvc := azure.New(logger, cfg.Azure, opts)

	if logger != nil {
		ldapSvc = ldapcli.NewLogging(logger, traceId)(ldapSvc)
		alarmSvc = alarm.NewLogging(logger, traceId)(alarmSvc)
		s3Cli = s3.NewLogging(logger, traceId)(s3Cli)
		fastChatSvc = fastchat.NewLogging(logger, traceId)(fastChatSvc)
		paasChatSvc = paaschat.NewLogging(logger, traceId)(paasChatSvc)
		azureSvc = azure.NewLogging(logger, traceId)(azureSvc)

		if debug {
			b, _ := json.Marshal(cfg.Ldap)
			_ = level.Debug(logger).Log("ldap.config", string(b))
			b, _ = json.Marshal(cfg.S3)
			_ = level.Debug(logger).Log("s3.config", string(b))
		}
	}

	// 如果tracer有的话
	if tracer != nil {
		s3Cli = s3.NewTracing(tracer)(s3Cli)
		alarmSvc = alarm.NewTracing(tracer)(alarmSvc)
		fastChatSvc = fastchat.NewTracing(tracer)(fastChatSvc)
		ldapSvc = ldapcli.NewTracing(tracer)(ldapSvc)
		paasChatSvc = paaschat.NewTracing(tracer)(paasChatSvc)
		azureSvc = azure.NewTracing(tracer)(azureSvc)
	}

	return &api{
		alarm:       alarmSvc,
		fastChatSvc: fastChatSvc,
		ldapSvc:     ldapSvc,
		s3Client:    s3Cli,
		paasChatSvc: paasChatSvc,
		azure:       azureSvc,
	}
}
