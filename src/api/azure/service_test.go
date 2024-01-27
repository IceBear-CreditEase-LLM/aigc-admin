package azure

import (
	"context"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
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
		ChatApi: struct {
			Host   string
			ApiKey string
		}{Host: "http://chat-api:8080", ApiKey: "sk-q8oSE4F7ANJPI7L60NBEENAGYXbYdS6J7gPFDPIFx24"},
	}, opts)
}

func TestService_TTS(t *testing.T) {
	svc := initSvc()
	ctx := context.Background()
	req := TTSRequest{
		Lang: "zh-CN",
		Name: "zh-CN-YunxiNeural",
		Text: "在一个宁静的小镇上，有一家独特的书店。店主是个热爱阅读的老者，他对每本书都了如指掌。这家书店没有固定的营业时间，无论何时走进去，都能感受到那份宁静和平和。",
	}
	_, filePath, err := svc.TTS(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(filePath)
}
