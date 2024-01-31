package paaschat

import (
	"context"
	"encoding/json"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Middleware func(Service) Service

type Config struct {
	Debug  bool
	ApiKey string
	Host   string
}

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"-"`
	Message string      `json:"message"`
	TraceId string      `json:"traceId"`
}

type DeployModelRequest struct {
	ModelName    string `json:"modelName"`
	Replicas     int    `json:"replicas"`
	Label        string `json:"label"`
	Gpu          int    `json:"gpu"`
	Quantization string `json:"quantization"`
	Vllm         bool   `json:"vllm"`
	MaxGpuMemory int    `json:"maxGpuMemory"`
}

type Service interface {
	ChatCompletionStream(ctx context.Context, request openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, err error)
	DeployModel(ctx context.Context, request DeployModelRequest) (err error)
	UndeployModel(ctx context.Context, modelName string) (err error)
	CancelFineTuningJob(ctx context.Context, jobId string) (err error)
	Wav2lipSynthesisCancel(ctx context.Context, uuid string) (err error)
	Wav2lipSynthesisCheck(ctx context.Context) (err error)
}

type service struct {
	logger log.Logger
	apiKey string
	host   string
	debug  bool
	opts   []kithttp.ClientOption
}

func (s *service) CancelFineTuningJob(ctx context.Context, jobId string) (err error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/fine_tuning/jobs/%s/cancel", s.host, jobId))
	if err != nil {
		return err
	}
	var res Response
	ep := kithttp.NewClient(http.MethodPatch, u, kithttp.EncodeJSONRequest, decodeJsonResponse(&res), s.opts...).Endpoint()
	_, err = ep(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeployModel(ctx context.Context, request DeployModelRequest) (err error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/deployment/%s", s.host, request.ModelName))
	if err != nil {
		return err
	}
	var res Response
	ep := kithttp.NewClient(http.MethodPatch, u, kithttp.EncodeJSONRequest, decodeJsonResponse(&res), s.opts...).Endpoint()
	_, err = ep(ctx, request)
	if err != nil {
		return err
	}
	if !res.Success {
		return errors.New(res.Message)
	}
	return nil
}

func (s *service) UndeployModel(ctx context.Context, modelName string) (err error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/deployment/%s", s.host, modelName))
	if err != nil {
		return err
	}
	var res Response
	ep := kithttp.NewClient(http.MethodDelete, u, kithttp.EncodeJSONRequest, decodeJsonResponse(&res), s.opts...).Endpoint()
	_, err = ep(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ChatCompletionStream(ctx context.Context, request openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, err error) {
	config := openai.DefaultConfig(s.apiKey)
	config.BaseURL = fmt.Sprintf("%s/v1", s.host)
	httpClient := http.DefaultClient
	if s.debug {
		httpClient = &http.Client{
			Transport: &proxyRoundTripper{},
		}
	}
	config.HTTPClient = httpClient
	client := openai.NewClientWithConfig(config)
	stream, err = client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (s *service) Wav2lipSynthesisCancel(ctx context.Context, uuid string) (err error) {
	u, err := url.Parse(fmt.Sprintf("%s/wav2lip/synthesis/%s/cancel", s.host, uuid))
	if err != nil {
		return err
	}
	var res Response
	ep := kithttp.NewClient(http.MethodPut, u, kithttp.EncodeJSONRequest, decodeJsonResponse(&res), s.opts...).Endpoint()
	_, err = ep(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Wav2lipSynthesisCheck(ctx context.Context) (err error) {
	u, err := url.Parse(fmt.Sprintf("%s/wav2lip/synthesis/check", s.host))
	if err != nil {
		return err
	}
	var res Response
	ep := kithttp.NewClient(http.MethodGet, u, kithttp.EncodeJSONRequest, decodeJsonResponse(&res), s.opts...).Endpoint()
	_, err = ep(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

type proxyRoundTripper struct {
	traceId string
	before  []kithttp.RequestFunc
	after   []kithttp.ClientResponseFunc
}

func (s *proxyRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	dump, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(dump))
	defer func() {
		if res != nil {
			dump, _ = httputil.DumpResponse(res, true)
			fmt.Println(string(dump))
		}
	}()
	return http.DefaultTransport.RoundTrip(req)
}

func decodeJsonResponse(data interface{}) func(ctx context.Context, res *http.Response) (response interface{}, err error) {
	return func(ctx context.Context, res *http.Response) (response interface{}, err error) {
		if res.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(res.Body)
			return res, fmt.Errorf("http status code is %d, body %s", res.StatusCode, string(body))
		}
		if data == nil {
			return res, nil
		}
		if err = json.NewDecoder(res.Body).Decode(data); err != nil {
			return res, errors.Wrap(err, "json decode")
		}
		return res, nil
	}
}

func New(logger log.Logger, cfg Config, opts []kithttp.ClientOption) Service {
	s := &service{
		logger: logger,
		apiKey: cfg.ApiKey,
		host:   cfg.Host,
		debug:  cfg.Debug,
		opts:   opts,
	}
	s.opts = append(s.opts, kithttp.ClientBefore(func(ctx context.Context, request *http.Request) context.Context {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
		return ctx
	}))
	return s
}
