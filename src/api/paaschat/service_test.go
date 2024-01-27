package paaschat

import (
	"context"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"testing"
)

func initSvc() Service {
	logger := log.NewLogfmtLogger(os.Stdout)
	var opts []kithttp.ClientOption
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
	return New(logger, Config{
		Debug:  true,
		ApiKey: "sk-4W6vS2nG8mC1pT3kX5rH7fJ9bQ0dZ4lY2cV1xN3aM9gB6qD8",
		Host:   "http://chat-api:8080",
	}, opts)
}

func TestService_ChatCompletionStream(t *testing.T) {
	ctx := context.Background()
	svc := initSvc()
	stream, err := svc.ChatCompletionStream(ctx, openai.ChatCompletionRequest{
		MaxTokens: 300,
		Model:     "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Hello!",
			},
		},
		Stream: true,
	})
	if err != nil {
		t.Error(err)
	}
	defer stream.Close()
	for {
		choice, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			t.Log(err)
			break
		}
		if errors.Is(err, openai.ErrTooManyEmptyStreamMessages) {
			t.Error(err)
			break
		}
		if err != nil {
			t.Error(err)
			break
		}
		t.Log(choice.Choices[0].Delta.Content)
	}
	t.Log("ok")
}

func TestService_DeployModel(t *testing.T) {
	ctx := context.Background()
	svc := initSvc()
	err := svc.DeployModel(ctx, DeployModelRequest{
		ModelName:    "test",
		Replicas:     1,
		Label:        "test",
		Gpu:          1,
		Quantization: "float16",
		Vllm:         false,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log("ok")
}
