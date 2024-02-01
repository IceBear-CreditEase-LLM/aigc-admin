/**
 * @Time: 2023/02/15 15:12
 * @Author: solacowa@gmail.com
 * @File: test_init.go
 * @Software: GoLand
 */

package tests

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api/fastchat"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api/ldapcli"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/encode"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/logging"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-redis/redis/v8"
	redisclient "github.com/icowan/redis-client"
	"github.com/sashabaranov/go-openai"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultHttpPort   = ":8080"
	DefaultConfigPath = "/usr/local/aigc-admin/etc/app.cfg"
	DefaultWebPath    = "/usr/local/aigc-admin/web"

	// [DB相关]
	EnvNameDbDrive       = "AIGC_DB_DRIVE"
	EnvNameMysqlHost     = "AIGC_MYSQL_HOST"
	EnvNameMysqlPort     = "AIGC_MYSQL_PORT"
	EnvNameMysqlUser     = "AIGC_MYSQL_USER"
	EnvNameMysqlPassword = "AIGC_MYSQL_PASSWORD"
	EnvNameMysqlDatabase = "AIGC_MYSQL_DATABASE"
	EnvNameRedisHosts    = "AIGC_REDIS_HOSTS"
	EnvNameRedisDb       = "AIGC_REDIS_DB"
	EnvNameRedisPassword = "AIGC_REDIS_PASSWORD"
	EnvNameRedisPrefix   = "AIGC_REDIS_PREFIX"

	// [跨域]
	EnvNameEnableCORS           = "AIGC_ENABLE_CORS"
	EnvNameCORSAllowMethods     = "AIGC_CORS_ALLOW_METHODS"
	EnvNameCORSAllowHeaders     = "AIGC_CORS_ALLOW_HEADERS"
	EnvNameCORSAllowCredentials = "AIGC_CORS_ALLOW_CREDENTIALS"
	EnvNameCORSAllowOrigins     = "AIGC_CORS_ALLOW_ORIGINS"

	// [Trace相关]
	EnvNameTracerEnable         = "AIGC_TRACER_ENABLE"
	EnvNameTracerDrive          = "AIGC_TRACER_DRIVE"
	EnvNameTracerJaegerHost     = "AIGC_TRACER_JAEGER_HOST"
	EnvNameTracerJaegerParam    = "AIGC_TRACER_JAEGER_PARAM"
	EnvNameTracerJaegerType     = "AIGC_TRACER_JAEGER_TYPE"
	EnvNameTracerJaegerLogSpans = "AIGC_TRACER_JAEGER_LOG_SPANS"

	// [外部Service相关]
	EnvNameServerHttpProxy       = "AIGC_SERVER_HTTP_PROXY"
	EnvNameServiceAlarmHost      = "AIGC_SERVICE_ALARM_HOST"    // 告警相关
	EnvNameServiceGptHost        = "AIGC_SERVICE_CHAT_API_HOST" // chat-api相关
	EnvNameServiceOpenAiHost     = "AIGC_SERVICE_OPENAI_HOST"
	EnvNameServiceOpenAiToken    = "AIGC_SERVICE_OPENAI_TOKEN"
	EnvNameServiceOpenAiModel    = "AIGC_SERVICE_OPENAI_MODEL"
	EnvNameServiceOpenAiOrgId    = "AIGC_SERVICE_OPENAI_ORG_ID"
	EnvNameServiceS3Host         = "AIGC_SERVICE_S3_HOST" // S3对象存储相当
	EnvNameServiceS3AccessKey    = "AIGC_SERVICE_S3_ACCESS_KEY"
	EnvNameServiceS3SecretKey    = "AIGC_SERVICE_S3_SECRET_KEY"
	EnvNameServiceS3S3Url        = "AIGC_SERVICE_S3_S3URL"
	EnvNameServiceS3Region       = "AIGC_SERVICE_S3_REGION"
	EnvNameServiceS3Bucket       = "AIGC_SERVICE_S3_BUCKET"
	EnvNameServiceS3BucketPublic = "AIGC_SERVICE_S3_BUCKET_PUBLIC"
	EnvNameServiceS3DownloadUrl  = "AIGC_SERVICE_S3_DOWNLOAD_URL"
	EnvNameServiceS3ProjectName  = "AIGC_SERVICE_S3_PROJECT_NAME"
	EnvNameServiceS3Cluster      = "AIGC_SERVICE_S3_CLUSTER"
	EnvNameDockerWorkspace       = "AIGC_Docker_WORKSPACE" // chat-api 相关

	// [LDAP 相关]
	EnvNameLdapHost        = "AIGC_LDAP_HOST"
	EnvNameLdapPort        = "AIGC_LDAP_PORT"
	EnvNameLdapUseSSL      = "AIGC_LDAP_USE_SSL"
	EnvNameLdapBaseDN      = "AIGC_LDAP_BASE_DN"
	EnvNameLdapBindUser    = "AIGC_LDAP_BIND_USER"
	EnvNameLdapBindPass    = "AIGC_LDAP_BIND_PASS"
	EnvNameLdapUserFilter  = "AIGC_LDAP_USER_FILTER"
	EnvNameLdapGroupFilter = "AIGC_LDAP_GROUP_FILTER"
	EnvNameLdapUserAttr    = "AIGC_LDAP_USER_ATTR"

	// [以下是aigc-admin模块配置]
	EnvHttpPort           = "AIGC_SERVER_HTTP_PORT"
	EnvNameServerLogDrive = "AIGC_SERVER_LOG_DRIVE"
	EnvNameServerLogPath  = "AIGC_SERVER_LOG_PATH"
	EnvNameServerName     = "AIGC_SERVER_NAME"
	EnvNameServerDebug    = "AIGC_SERVER_DEBUG"
	EnvNameServerKey      = "AIGC_SERVER_KEY"
	EnvNameServerLogLevel = "AIGC_SERVER_LOG_LEVEL"
	EnvNameServerLogName  = "AIGC_SERVER_LOG_NAME"

	DefaultDbDrive       = "mysql"
	DefaultMysqlHost     = "localhost"
	DefaultMysqlPort     = 3306
	DefaultMysqlUser     = "aigc"
	DefaultMysqlPassword = ""
	DefaultMysqlDatabase = "aigc"
	DefaultRedisHosts    = "localhost:6379"
	DefaultRedisDb       = 0
	DefaultRedisPassword = ""
	DefaultRedisPrefix   = "aigc"

	DefaultServerName      = "aigc-admin"
	DefaultServerKey       = ""
	DefaultServerLogLevel  = "all"
	DefaultServerLogDrive  = "term"
	DefaultServerLogPath   = ""
	DefaultServerLogName   = "aigc-admin.log"
	DefaultServerDebug     = false
	DefaultEnableCORS      = false
	DefaultServerHttpProxy = ""

	DefaultCORSAllowOrigins     = "*"
	DefaultCORSAllowMethods     = "GET,POST,PUT,DELETE,OPTIONS"
	DefaultCORSAllowHeaders     = "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization"
	DefaultCORSAllowCredentials = true
	DefaultCORSExposeHeaders    = "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type"

	DefaultJaegerEnable           = false
	DefaultJaegerDrive            = "jaeger"
	DefaultJaegerHost             = "jaeger:6832"
	DefaultJaegerParam    float64 = 1
	DefaultJaegerType             = "const"
	DefaultJaegerLogSpans         = false

	DefaultServiceAlarmHost = ""
	// [gpt]相关
	DefaultServiceChatApiHost = "http://chat-api:8080"
	DefaultServiceOpenAiHost  = "https://api.openai.com/v1"
	DefaultServiceOpenAiToken = "sk-001"
	DefaultServiceOpenAiModel = openai.GPT3Dot5Turbo
	DefaultServiceOpenAiOrgId = ""

	// [ldap相关]
	DefaultLdapHost        = "ldap"
	DefaultLdapPort        = 389
	DefaultLdapBaseDn      = "OU=HABROOT,DC=corp"
	DefaultLdapBindUser    = "aigc_ldap"
	DefaultLdapBindPass    = ""
	DefaultLdapUserFilter  = "(userPrincipalName=%s)"
	DefaultLdapGroupFilter = ""
	DefaultLdapAttributes  = "name,mail,userPrincipalName,displayName,sAMAccountName"

	// [s3]
	DefaultServiceS3Host         = "http://s3"
	DefaultServiceS3AccessKey    = ""
	DefaultServiceS3SecretKey    = ""
	DefaultServiceS3Bucket       = "aigc"
	DefaultServiceS3BucketPublic = "aigc"
	DefaultServiceS3Region       = "default"
	DefaultServiceS3Cluster      = "ceph-c2"
	DefaultDockerWorkspace       = "/tmp"

	DefaultServiceChatToken = ""
)

var (
	httpAddr, configPath string
	webPath              string
	logger               log.Logger
	Logger               log.Logger
	gormDB               *gorm.DB
	db                   *sql.DB
	err                  error
	Store                repository.Repository
	namespace            string
	webEmbed             bool
)

var (
	rdb    redis.UniversalClient
	apiSvc api.Service
	//hashId   hashids.HashIds
	dbDrive, mysqlHost, mysqlUser, mysqlPassword, mysqlDatabase                                        string
	mysqlPort, redisDb, ormPort                                                                        int
	redisAuth, redisHosts, redisPrefix, serverHttpProxy                                                string
	serverName, serverKey, serverLogLevel, serverLogDrive, serverLogPath, serverLogName                string
	corsAllowOrigins, corsAllowMethods, corsAllowHeaders, corsExposeHeaders                            string
	serverDebug, enableCORS, corsAllowCredentials, tracerEnable, tracerJaegerLogSpans, mysqlOrmMetrics bool
	tracerDrive, tracerJaegerHost, tracerJaegerType                                                    string
	tracerJaegerParam                                                                                  float64
	serviceAlarmHost                                                                                   string

	// [gpt]
	serviceGPTHost, serviceGPTModel                                               string
	serviceOpenAiHost, serviceOpenAiToken, serviceOpenAiModel, serviceOpenAiOrgId string

	// [s3]
	serviceS3Host, serviceS3AccessKey, serviceS3SecretKey, serviceS3Bucket, serviceS3Region, serviceS3Cluster string

	// [docker]
	dockerWorkspace string

	// [ldap]相关
	ldapHost, ldapBaseDn, ldapBindUser, ldapBindPass, ldapUserFilter, ldapGroupFilter string
	ldapPort                                                                          int
	ldapUserAttr                                                                      []string
	ldapUseSsl                                                                        bool

	corsHeaders   = make(map[string]string, 3)
	rateBucketNum = 50000
	traceId       = logging.TraceId

	goOS                                     = runtime.GOOS
	goArch                                   = runtime.GOARCH
	goVersion                                = runtime.Version()
	compiler                                 = runtime.Compiler
	version, buildDate, gitCommit, gitBranch string
)

func preRun() {
	webPath = envString("WEB_PATH", DefaultWebPath)

	httpAddr = envString(EnvHttpPort, DefaultHttpPort)
	namespace = envString("POD_NAMESPACE", envString("NAMESPACE", namespace))

	// [database]
	dbDrive = envString(EnvNameDbDrive, DefaultDbDrive)
	mysqlHost = envString(EnvNameMysqlHost, DefaultMysqlHost)
	mysqlPort, _ = strconv.Atoi(envString(EnvNameMysqlPort, strconv.Itoa(DefaultMysqlPort)))
	mysqlUser = envString(EnvNameMysqlUser, DefaultMysqlUser)
	mysqlPassword = envString(EnvNameMysqlPassword, DefaultMysqlPassword)
	mysqlDatabase = envString(EnvNameMysqlDatabase, DefaultMysqlDatabase)

	// [redis]
	redisHosts = envString(EnvNameRedisHosts, DefaultRedisHosts)
	redisDb, _ = strconv.Atoi(envString(EnvNameRedisDb, strconv.Itoa(DefaultRedisDb)))
	redisAuth = envString(EnvNameRedisPassword, DefaultRedisPassword)
	redisPrefix = envString(EnvNameRedisPrefix, DefaultRedisPrefix)

	// [cors]
	enableCORS, _ = strconv.ParseBool(envString(EnvNameEnableCORS, strconv.FormatBool(DefaultEnableCORS)))
	corsAllowMethods = envString(EnvNameCORSAllowMethods, DefaultCORSAllowMethods)
	corsAllowHeaders = envString(EnvNameCORSAllowHeaders, DefaultCORSAllowHeaders)
	corsAllowOrigins = envString(EnvNameCORSAllowOrigins, DefaultCORSAllowOrigins)
	corsAllowCredentials, _ = strconv.ParseBool(envString(EnvNameCORSAllowCredentials, strconv.FormatBool(DefaultCORSAllowCredentials)))

	// [trace]
	tracerEnable, _ = strconv.ParseBool(envString(EnvNameTracerEnable, strconv.FormatBool(DefaultJaegerEnable)))
	tracerDrive = envString(EnvNameTracerDrive, DefaultJaegerDrive)
	tracerJaegerParam, _ = strconv.ParseFloat(envString(EnvNameTracerJaegerParam, strconv.FormatFloat(tracerJaegerParam, 'f', -1, 64)), 64)
	tracerJaegerHost = envString(EnvNameTracerJaegerHost, DefaultJaegerHost)
	tracerJaegerType = envString(EnvNameTracerJaegerType, DefaultJaegerType)
	tracerJaegerLogSpans, _ = strconv.ParseBool(envString(EnvNameTracerJaegerLogSpans, strconv.FormatBool(DefaultJaegerLogSpans)))

	// [server]
	serverName = envString(EnvNameServerName, DefaultServerName)
	serverKey = envString(EnvNameServerKey, DefaultServerKey)
	serverLogLevel = envString(EnvNameServerLogLevel, DefaultServerLogLevel)
	serverLogDrive = envString(EnvNameServerLogDrive, DefaultServerLogDrive)
	serverLogPath = envString(EnvNameServerLogPath, DefaultServerLogPath)
	serverLogName = envString(EnvNameServerLogName, DefaultServerLogName)
	serverHttpProxy = envString(EnvNameServerHttpProxy, DefaultServerHttpProxy)
	serverDebug, _ = strconv.ParseBool(envString(EnvNameServerDebug, strconv.FormatBool(DefaultServerDebug)))

	// 以下是[service] 模块配置
	serviceAlarmHost = envString(EnvNameServiceAlarmHost, DefaultServiceAlarmHost)

	// [service.gpt]
	serviceGPTHost = envString(EnvNameServiceGptHost, DefaultServiceChatApiHost)
	serviceOpenAiHost = envString(EnvNameServiceOpenAiHost, DefaultServiceOpenAiHost)
	serviceOpenAiToken = envString(EnvNameServiceOpenAiToken, DefaultServiceOpenAiToken)
	serviceOpenAiModel = envString(EnvNameServiceOpenAiModel, DefaultServiceOpenAiModel)
	serviceOpenAiOrgId = envString(EnvNameServiceOpenAiOrgId, DefaultServiceOpenAiOrgId)

	// [ldap]
	ldapHost = envString(EnvNameLdapHost, DefaultLdapHost)
	ldapPort, _ = strconv.Atoi(envString(EnvNameLdapPort, strconv.Itoa(DefaultLdapPort)))
	ldapUseSsl, _ = strconv.ParseBool(envString(EnvNameLdapUseSSL, "false"))
	ldapBindUser = envString(EnvNameLdapBindUser, DefaultLdapBindUser)
	ldapBindPass = envString(EnvNameLdapBindPass, DefaultLdapBindPass)
	ldapBaseDn = envString(EnvNameLdapBaseDN, DefaultLdapBaseDn)
	ldapUserFilter = envString(EnvNameLdapUserFilter, DefaultLdapUserFilter)
	ldapUserAttr = strings.Split(envString(EnvNameLdapUserAttr, DefaultLdapAttributes), ",")

	// [service.s3]
	serviceS3Host = envString(EnvNameServiceS3Host, DefaultServiceS3Host)
	serviceS3AccessKey = envString(EnvNameServiceS3AccessKey, DefaultServiceS3AccessKey)
	serviceS3SecretKey = envString(EnvNameServiceS3SecretKey, DefaultServiceS3SecretKey)
	serviceS3Bucket = envString(EnvNameServiceS3Bucket, DefaultServiceS3Bucket)
	serviceS3Region = envString(EnvNameServiceS3Region, DefaultServiceS3Region)
	serviceS3Cluster = envString(EnvNameServiceS3Cluster, DefaultServiceS3Cluster)

	// [docker]
	dockerWorkspace = envString(EnvNameDockerWorkspace, DefaultDockerWorkspace)

}

func Init() (rdb redis.UniversalClient, apiSvc api.Service, err error) {
	preRun()
	err = prepare(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	return rdb, apiSvc, nil
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func prepare(ctx context.Context) error {
	logger = log.NewLogfmtLogger(os.Stdout)

	logger = logging.SetLogging(logger, serverLogPath, serverLogName, serverLogLevel, serverName, serverLogDrive)

	// 连接数据库
	if strings.EqualFold(dbDrive, "mysql") {
		dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
			mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
		var dbErr error
		sqlDB, err := sql.Open("mysql", dbUrl)
		if err != nil {
			_ = level.Error(logger).Log("sql", "Open", "err", err.Error())
			return err
		}
		gormDB, err = gorm.Open(mysql.New(mysql.Config{
			Conn:              sqlDB,
			DefaultStringSize: 255,
		}), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if dbErr != nil {
			_ = level.Error(logger).Log("db", "connect", "err", dbErr.Error())
			dbErr = encode.ErrServerStartDbConnect.Wrap(dbErr)
			return dbErr
		}
		//gormDB.Statement.Clauses["soft_delete_enabled"] = clause.Clause{}
		db, dbErr = gormDB.DB()
		if dbErr != nil {
			_ = level.Error(logger).Log("gormDB", "DB", "err", dbErr.Error())
			dbErr = encode.ErrServerStartDbConnect.Wrap(dbErr)
			return dbErr
		}
		_ = level.Debug(logger).Log("mysql", "connect", "success", true)
	}

	if !strings.EqualFold(serverLogPath, "") {
		gormDB.Logger = logging.NewGormLogging(logger)
	} else {
		gormDB.Logger = gormlogger.Default.LogMode(gormlogger.Info)
	}
	if mysqlOrmMetrics {
		//if err = gormDB.Use(gormprometheus.New(gormprometheus.Config{
		//	DBName:          mysqlDatabase,
		//	RefreshInterval: 15,
		//	//PushAddr:        prometheusHost,  // 如果配置了 `PushAddr`，则推送指标
		//	StartServer: false, // 启用一个 http 服务来暴露指标
		//	//HTTPServerPort: uint32(ormPort), // 配置 http 服务监听端口，默认端口为 8080 （如果您配置了多个，只有第一个 `HTTPServerPort` 会被使用）
		//	MetricsCollector: []gormprometheus.MetricsCollector{
		//		&gormprometheus.MySQL{
		//			VariableNames: []string{"Threads_running"},
		//		},
		//	}, // 用户自定义指标
		//})); err != nil {
		//	_ = level.Error(logger).Log("gormDB", "Use", "plugin", "prometheus", "err", err.Error())
		//}
	}

	// 链路追踪
	if tracerEnable {

	}

	// 实例化redis
	rdb, err = redisclient.NewRedisClient(redisHosts, redisAuth, redisPrefix, redisDb, nil)
	if err != nil {
		_ = level.Error(logger).Log("redis", "connect", "err", err.Error())
		return err
	}
	_ = level.Debug(logger).Log("redis", "connect", "success", true)

	var clientOpts []kithttp.ClientOption
	dialer := &net.Dialer{
		Timeout:   10 * time.Minute,
		KeepAlive: 10 * time.Minute,
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Minute,
		Transport: &http.Transport{
			DialContext: dialer.DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}
	clientOpts = []kithttp.ClientOption{
		kithttp.SetClient(httpClient),
		//kithttp.ClientBefore(kittracing.ContextToHTTP(nil, logger)),
	}

	// 实例化仓库
	Store = repository.New(gormDB, logger, logging.TraceId, nil)

	// 实例化外部API
	apiSvc = api.NewApi(ctx, logger, logging.TraceId, serverDebug, nil, &api.Config{
		Namespace: namespace, ServiceName: serverName,
		FastChat: fastchat.Config{
			OpenAiEndpoint: serviceOpenAiHost,
			OpenAiToken:    serviceOpenAiToken,
			OpenAiModel:    serviceOpenAiModel,
			OpenAiOrgId:    serviceOpenAiOrgId,
			//chatapiEndpoint: serviceGPTHost,
			//chat - apiModel:    serviceGPTModel,
		},
		Ldap: ldapcli.Config{
			Host:         ldapHost,
			Port:         ldapPort,
			UseSSL:       ldapUseSsl,
			BindUser:     ldapBindUser,
			BindPassword: ldapBindPass,
			BindDN:       ldapBaseDn,
			Attributes:   ldapUserAttr,
			Filter:       ldapUserFilter,
		},
		Alarm: struct {
			Host                   string
			Namespace, ServiceName string
		}{Host: serviceAlarmHost, Namespace: namespace, ServiceName: serverName},
	}, clientOpts, rdb, dockerWorkspace)

	Logger = logger
	return err
}
