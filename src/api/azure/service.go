package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/util"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
)

var tmpUpload = "/tmp/upload/wav"

type Middleware func(Service) Service

type Service interface {
	// TTS 文转音
	TTS(ctx context.Context, request TTSRequest) (response string, filePath string, err error)
}

type Config struct {
	ChatApi struct {
		Host   string
		ApiKey string
	}
}

type service struct {
	Config
	logger log.Logger
	opts   []kithttp.ClientOption
}

func (s *service) TTS(ctx context.Context, req TTSRequest) (response string, filePath string, err error) {

	apiUrl := fmt.Sprintf("%s/voice/tts", s.ChatApi.Host)

	u, _ := url.Parse(apiUrl)

	ep := kithttp.NewClient(http.MethodPost, u, func(ctx context.Context, r *http.Request, request interface{}) error {
		r.Header.Set("Content-Type", "application/json; charset=utf-8")
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.ChatApi.ApiKey))
		var b bytes.Buffer
		r.Body = io.NopCloser(&b)
		return json.NewEncoder(&b).Encode(request)
	}, func(ctx context.Context, response2 *http.Response) (response interface{}, err error) {
		b, err := io.ReadAll(response2.Body)
		if err != nil {
			return nil, errors.Wrap(err, "read body error")
		}
		if response2.StatusCode != http.StatusOK {
			return nil, errors.Errorf("http status code not 200, code: %d, body:%s", response2.StatusCode, b)
		}
		var apiResp TTSResult
		err = json.Unmarshal(b, &apiResp)
		if err != nil {
			return nil, errors.Wrap(err, "json.Unmarshal apiResp")
		}
		if !apiResp.Success {
			return nil, errors.New(apiResp.Message)
		}
		return apiResp.Data.Data, err
	}, s.opts...).Endpoint()

	r, err := ep(ctx, req)
	if err != nil {
		return "", "", err
	}

	//暂存本地
	filePath, err = util.WriteAudioFile(r.(string), tmpUpload, "wav")
	if err != nil {
		return "", "", err
	}

	return r.(string), filePath, nil
}

func New(logger log.Logger, cfg Config, opts []kithttp.ClientOption) Service {
	s := &service{
		Config: cfg,
		logger: logger,
		opts:   opts,
	}
	return s
}
