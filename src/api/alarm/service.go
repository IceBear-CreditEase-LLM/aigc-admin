package alarm

import (
	"context"
	"encoding/json"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type Config struct {
	Host                   string
	Namespace, ServiceName string
}

type pushRequest struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	Metrics       string `json:"metrics"`
	Project       string `json:"project"`
	Service       string `json:"service"`
	Level         string `json:"level"`
	SilencePeriod int    `json:"silencePeriod"`
	Timestamp     int64  `json:"timestamp"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type Level string

func (s Level) String() string {
	return string(s)
}

const (
	LevelInfo    Level = "info"
	LevelWarning Level = "warn"
	LevelError   Level = "error"
)

type Middleware func(Service) Service

// Service 统一告警中心接口
type Service interface {
	// Push 发送到统一告警中心
	// title: 预警标题 必填 长度45
	// content: 预警内容 必填 长度4000
	// metrics: 预警标识 必填 长度80 消息发送方定义的预警标识，作为压制key使用，不同预警按不同的项目、服务和标识压制
	// project: 项目英文名 必填 长度50 与dolphin的英文名做对应
	// service: 服务英文名 必填 长度50 与dolphin的英文名做对应
	// level: 预警级别 非必填 长度10 可指定级别：info、warn、error, 不填默认info
	// silencePeriod: 压制间隔 非必填 单位：分钟； 为空默认30分钟
	Push(ctx context.Context, title, content, metrics string, level Level, silencePeriod int) (err error)
}

type service struct {
	traceId                  string
	host                     string
	projectName, serviceName string
	clientOpts               []kithttp.ClientOption
}

func (s *service) Push(ctx context.Context, title, content, metrics string, level Level, silencePeriod int) (err error) {
	return
}

func New(traceId string, cfg Config, opts []kithttp.ClientOption) Service {
	return &service{host: cfg.Host, clientOpts: opts, traceId: traceId, projectName: cfg.Namespace, serviceName: cfg.ServiceName}
}

func encodeJsonResponse(ctx context.Context, res *http.Response) (response interface{}, err error) {
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	var resp Response
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}
	if strings.EqualFold(strings.ToUpper(resp.Status), "OK") {
		return nil, errors.New(resp.Message)
	}
	return resp, err
}
