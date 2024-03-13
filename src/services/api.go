package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/middleware"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/services/alarm"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/services/dockerapi"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/services/fastchat"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/services/ldapcli"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/services/runtime"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/services/s3"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
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

	DockerApi() dockerapi.Service
	// Runtime runtime服务
	Runtime() runtime.Service
}

type api struct {
	logger      log.Logger
	s3Client    s3.Service
	alarm       alarm.Service
	traceId     string
	fastChatSvc fastchat.Service
	ldapSvc     ldapcli.Service
	dockerapi   dockerapi.Service
	runtimeSvc  runtime.Service
}

func (s *api) Runtime() runtime.Service {
	return s.runtimeSvc
}

func (s *api) DockerApi() dockerapi.Service {
	return s.dockerapi
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

// NewApi 中间件有顺序,在后面的会最先执行
func NewApi(_ context.Context, logger log.Logger, traceId string, debug bool, tracer opentracing.Tracer, cfg *Config, opts []kithttp.ClientOption, workspace string) Service {
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
	fastChatSvc := fastchat.New(cfg.FastChat, fastChatSvcOpts)
	ldapSvc := ldapcli.New(cfg.Ldap)
	dockerapiSvc := dockerapi.New(workspace)
	runtimeSvc := runtime.NewDocker(runtime.WithWorkspace("tmp"))
	//runtimeSvc := runtime.NewDocker(workspace, dockerapiSvc)

	if logger != nil {
		ldapSvc = ldapcli.NewLogging(logger, traceId)(ldapSvc)
		alarmSvc = alarm.NewLogging(logger, traceId)(alarmSvc)
		s3Cli = s3.NewLogging(logger, traceId)(s3Cli)
		fastChatSvc = fastchat.NewLogging(logger, traceId)(fastChatSvc)
		runtimeSvc = runtime.NewLogging(logger, traceId)(runtimeSvc)
		dockerapiSvc = dockerapi.NewLogging(logger, traceId)(dockerapiSvc)

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
		runtimeSvc = runtime.NewTracing(tracer)(runtimeSvc)
		dockerapiSvc = dockerapi.NewTracing(tracer)(dockerapiSvc)
	}

	return &api{
		alarm:       alarmSvc,
		fastChatSvc: fastChatSvc,
		ldapSvc:     ldapSvc,
		s3Client:    s3Cli,
		dockerapi:   dockerapiSvc,
		runtimeSvc:  runtimeSvc,
	}
}
