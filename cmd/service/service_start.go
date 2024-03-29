package service

import (
	"context"
	"embed"
	"fmt"
	tiktoken2 "github.com/IceBear-CreditEase-LLM/aigc-admin/src/helpers/tiktoken"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/pkg/assistants"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/pkg/auth"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/pkg/channels"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/pkg/datasets"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/pkg/files"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/pkg/finetuning"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/pkg/models"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/pkg/sys"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/pkg/tools"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/types"
	"github.com/pkoukk/tiktoken-go"
	"github.com/tmc/langchaingo/llms/openai"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-kit/kit/tracing/opentracing"

	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api/alarm"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/encode"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/logging"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/middleware"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	"github.com/oklog/oklog/pkg/group"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
)

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "启动http服务",
		Example: `## 启动命令
aigc-admin start -p :8080
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return start(cmd.Context())
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := prepare(cmd.Context()); err != nil {
				_ = level.Error(logger).Log("cmd", "start.PreRunE", "err", err.Error())
				return err
			}
			// 判断是否需要初始化数据，如果没有则初始化数据
			if !gormDB.Migrator().HasTable(types.Accounts{}) {
				if err = generateTable(); err != nil {
					_ = level.Error(logger).Log("cmd.start.PreRunE", "generateTable", "err", err.Error())
					return err
				}
				if err = initData(); err != nil {
					_ = level.Error(logger).Log("cmd.start.PreRunE", "initData", "err", err.Error())
					return err
				}
			}
			//aigc-admin channelID
			channelRes, err := store.Chat().FindChannelByApiKey(cmd.Context(), serverChannelKey)
			if err != nil {
				_ = level.Error(logger).Log("cmd.start.PreRunE", "FindChannelByApiKey", "err", err.Error())
				return err
			}
			channelId = int(channelRes.ID)
			return nil
		},
	}

	tracer stdopentracing.Tracer

	opts []kithttp.ServerOption

	WebFs  embed.FS
	DataFs embed.FS

	authSvc auth.Service

	fileSvc       files.Service
	channelSvc    channels.Service
	modelSvc      models.Service
	fineTuningSvc finetuning.Service
	sysSvc        sys.Service
	datasetSvc    datasets.Service
	toolsSvc      tools.Service
	assistantsSvc assistants.Service
)

func start(ctx context.Context) (err error) {

	tiktoken.SetBpeLoader(tiktoken2.NewBpeLoader(DataFs))

	authSvc = auth.New(logger, traceId, store, rdb, apiSvc)
	fileSvc = files.NewService(logger, traceId, store, apiSvc, files.Config{
		LocalDataPath: serverStoragePath,
		ServerUrl:     fmt.Sprintf("%s/storage", serverDomain),
	})
	channelSvc = channels.NewService(logger, traceId, store, apiSvc)
	modelSvc = models.NewService(logger, traceId, store, apiSvc, aigcDataCfsPath)
	fineTuningSvc = finetuning.New(traceId, logger, store, serviceS3Bucket, serviceS3AccessKey, serviceS3SecretKey, apiSvc, rdb, aigcDataCfsPath)
	sysSvc = sys.NewService(logger, traceId, store, apiSvc)
	datasetSvc = datasets.New(logger, traceId, store)
	toolsSvc = tools.New(logger, traceId, store)
	assistantsSvc = assistants.New(logger, traceId, store, []kithttp.ClientOption{
		//kithttp.ClientBefore(func(ctx context.Context, request *http.Request) context.Context {
		//	dump, _ := httputil.DumpRequest(request, true)
		//	fmt.Println(string(dump))
		//	return ctx
		//}),
		//kithttp.ClientAfter(func(ctx context.Context, response *http.Response) context.Context {
		//	dump, _ := httputil.DumpResponse(response, true)
		//	fmt.Println(string(dump))
		//	return ctx
		//}),
	}, []openai.Option{
		openai.WithToken(serviceLocalAiToken),
		openai.WithBaseURL(serviceLocalAiHost),
	})

	if logger != nil {
		authSvc = auth.NewLogging(logger, logging.TraceId)(authSvc)
		fileSvc = files.NewLogging(logger, logging.TraceId)(fileSvc)
		channelSvc = channels.NewLogging(logger, logging.TraceId)(channelSvc)
		modelSvc = models.NewLogging(logger, logging.TraceId)(modelSvc)
		fineTuningSvc = finetuning.NewLogging(logger, logging.TraceId)(fineTuningSvc)
		sysSvc = sys.NewLogging(logger, logging.TraceId)(sysSvc)
		datasetSvc = datasets.NewLogging(logger, logging.TraceId)(datasetSvc)
		toolsSvc = tools.NewLogging(logger, logging.TraceId)(toolsSvc)
	}

	if tracer != nil {
		authSvc = auth.NewTracing(tracer)(authSvc)
		fileSvc = files.NewTracing(tracer)(fileSvc)
		channelSvc = channels.NewTracing(tracer)(channelSvc)
		modelSvc = models.NewTracing(tracer)(modelSvc)
		fineTuningSvc = finetuning.NewTracing(tracer)(fineTuningSvc)
		sysSvc = sys.NewTracing(tracer)(sysSvc)
		datasetSvc = datasets.NewTracing(tracer)(datasetSvc)
		toolsSvc = tools.NewTracing(tracer)(toolsSvc)
	}

	g := &group.Group{}

	initHttpHandler(ctx, g)
	//initGRPCHandler(g)
	initCancelInterrupt(ctx, g)

	_ = level.Error(logger).Log("server exit", g.Run())
	return nil
}

func accessControl(h http.Handler, logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, val := range corsHeaders {
			w.Header().Set(key, val)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Connection", "keep-alive")

		if r.Method == "OPTIONS" {
			return
		}
		_ = level.Info(logger).Log("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)

		h.ServeHTTP(w, r)
	})
}

func initHttpHandler(ctx context.Context, g *group.Group) {
	httpLogger := log.With(logger, "component", "http")

	opts = []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encode.JsonError),
		kithttp.ServerErrorHandler(logging.NewLogErrorHandler(level.Error(logger), apiSvc)),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			guid := request.Header.Get("X-Request-Id")
			//token := request.Header.Get("Authorization")
			token := request.Header.Get("X-Token")
			tenantId := request.Header.Get("X-Tenant-Id")
			if strings.EqualFold(token, "") {
				token = request.URL.Query().Get("X-Token")
			}
			if strings.EqualFold(tenantId, "") {
				tenantId = request.URL.Query().Get("X-Tenant-Id")
			}
			request.Header.Set("Authorization", token)
			ctx = context.WithValue(ctx, kithttp.ContextKeyRequestAuthorization, token)
			ctx = context.WithValue(ctx, logging.TraceId, guid)
			ctx = context.WithValue(ctx, middleware.ContextKeyPublicTenantId, tenantId)
			ctx = context.WithValue(ctx, middleware.ContextKeyChannelId, channelId)
			return ctx
		}),
	}

	if tracer != nil {
		opts = append(opts,
			kithttp.ServerBefore(
				opentracing.HTTPToContext(tracer, "HTTPToContext", logger),
				middleware.TracingServerBefore(tracer),
			))
	}

	ems := []endpoint.Middleware{
		middleware.TracingMiddleware(tracer),                                                      // 2
		middleware.TokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), rateBucketNum)), // 1
	}

	authEms := []endpoint.Middleware{
		middleware.AuditMiddleware(logger, store),
		middleware.CheckTenantMiddleware(logger, store, tracer),
		middleware.CheckAuthMiddleware(logger, rdb, tracer),
	}
	authEms = append(authEms, ems...)

	r := mux.NewRouter()
	// auth模块
	r.PathPrefix("/api/auth").Handler(http.StripPrefix("/api/auth", auth.MakeHTTPHandler(authSvc, authEms, opts)))

	// file模块
	r.PathPrefix("/api/files").Handler(http.StripPrefix("/api", files.MakeHTTPHandler(fileSvc, authEms, opts)))
	// channel模块
	r.PathPrefix("/api/channels").Handler(http.StripPrefix("/api", channels.MakeHTTPHandler(channelSvc, authEms, opts)))
	// Model模块
	r.PathPrefix("/api/models").Handler(http.StripPrefix("/api", models.MakeHTTPHandler(modelSvc, authEms, opts)))
	// FineTuning模块
	r.PathPrefix("/api/finetuning").Handler(http.StripPrefix("/api", finetuning.MakeHTTPHandler(fineTuningSvc, authEms, opts)))
	// Sys模块
	r.PathPrefix("/api/sys").Handler(http.StripPrefix("/api/sys", sys.MakeHTTPHandler(sysSvc, authEms, opts)))
	// Dataset模块
	r.PathPrefix("/api/datasets").Handler(http.StripPrefix("/api/datasets", datasets.MakeHTTPHandler(datasetSvc, authEms, opts)))
	// Tools模块
	r.PathPrefix("/api/tools").Handler(http.StripPrefix("/api/tools", tools.MakeHTTPHandler(toolsSvc, authEms, opts)))
	// Assistants模块
	r.PathPrefix("/api/assistants").Handler(http.StripPrefix("/api/assistants", assistants.MakeHTTPHandler(assistantsSvc, authEms, opts)))
	// 对外metrics
	r.Handle("/metrics", promhttp.Handler())
	// 心跳检测
	r.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("ok"))
	})
	// 文件存储
	r.PathPrefix("/storage/").Handler(http.StripPrefix("/storage/", http.FileServer(http.Dir(serverStoragePath))))

	// web页面
	if webEmbed {
		fe, fsErr := fs.Sub(WebFs, DefaultWebPath)
		if fsErr != nil {
			_ = level.Error(logger).Log("FailedToSubPath", "web", "err", fsErr.Error())
		}
		r.PathPrefix("/").Handler(http.FileServer(http.FS(fe)))
	} else {
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(webPath)))
	}

	if enableCORS {
		corsHeaders["Access-Control-Allow-Origin"] = corsAllowOrigins
		corsHeaders["Access-Control-Allow-Methods"] = corsAllowMethods
		corsHeaders["Access-Control-Allow-Headers"] = corsAllowHeaders
		corsHeaders["Access-Control-credentials"] = strconv.FormatBool(corsAllowCredentials)
	}

	http.Handle("/", accessControl(r, httpLogger))

	g.Add(func() error {
		_ = level.Debug(httpLogger).Log("transport", "HTTP", "addr", httpAddr)
		go func() {
			_ = apiSvc.Alarm().Push(ctx, "服务启动", "服务它又起来了...", "service_start", alarm.LevelInfo, 1)
		}()
		return http.ListenAndServe(httpAddr, nil)
	}, func(e error) {
		closeConnection(ctx)
		_ = level.Error(httpLogger).Log("transport", "HTTP", "httpListener.Close", "http", "err", e)
		_ = apiSvc.Alarm().Push(ctx, "服务停止", fmt.Sprintf("msg: %s, err: %v", "服务它停了,是不是挂了...", e), "service_start", alarm.LevelInfo, 1)
		_ = level.Debug(logger).Log("db", "close", "err", db.Close())
		if rdb != nil {
			_ = level.Debug(logger).Log("rdb", "close", "err", rdb.Close())
		}
		os.Exit(1)
	})
}

func initCancelInterrupt(ctx context.Context, g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		select {
		case sig := <-c:
			if err != nil {
				_ = level.Error(logger).Log("rocketmq", "close", "err", err)
				return err
			}
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(err error) {
		close(cancelInterrupt)
	})
}

var localAddr string

func getLocalAddr() string {
	if localAddr != "" {
		return localAddr
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				localAddr = ipNet.IP.String()
				return localAddr
			}
		}
	}

	return ""
}
