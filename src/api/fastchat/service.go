package fastchat

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api/alarm"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/util"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Debug bool
	// OpenAI地址
	OpenAiEndpoint string
	// OpenAI token
	OpenAiToken string
	// OpenAI 默认模型
	OpenAiModel string
	// OpenAI 组织ID
	OpenAiOrgId string
	// Chat-api 地址
	PaasGptEndpoint string
	// PaasTPT 默认模型
	PaasGptModel string
	// SdHost 地址
	SdHost     string
	SdRedisKey string
	// svc 当前服务信息
	SvcIp     string
	SvcPort   int
	AuthToken string
}

type CtxPlatform string
type Platform string

const (
	PlatformOpenAI Platform = "OpenAI"

	ContextKeyPlatform CtxPlatform = "ctx-platform-name"
	ContextKeyApiKey   CtxPlatform = "ctx-api-key"
)

type ApiErrResponse struct {
	Object  string `json:"object"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Middleware func(service Service) Service

type Service interface {
	// UploadFile
	// file string Required
	// Name of the JSON Lines file to be uploaded.
	// If the purpose is set to "fine-tune", the file will be used for fine-tuning.
	// purpose string Required
	// The intended purpose of the uploaded documents.
	// Use "fine-tune" for fine-tuning. This allows us to validate the format of the uploaded file.
	UploadFile(ctx context.Context, modelName, fileName, filePath, purpose string) (res openai.File, err error)
	// CreateFineTuningJob 创建微调任务
	CreateFineTuningJob(ctx context.Context, req openai.FineTuningJobRequest) (res openai.FineTuningJob, err error)
	// RetrieveFineTuningJob 获取微调任务
	RetrieveFineTuningJob(ctx context.Context, jobId string) (res openai.FineTuningJob, err error)
	ListFineTune(ctx context.Context, modelName string) (res openai.FineTuneList, err error)
	// CancelFineTuningJob 取消微调任务
	CancelFineTuningJob(ctx context.Context, modelName, jobId string) (err error)
	// ChatCompletion 聊天处理
	// model: str
	// messages: List[Dict[str, str]]
	// temperature: Optional[float] = 0.7
	// top_p: Optional[float] = 1.0
	// n: Optional[int] = 1
	// max_tokens: Optional[int] = None
	// stop: Optional[Union[str, List[str]]] = None
	// stream: Optional[bool] = False
	// presence_penalty: Optional[float] = 0.0
	// frequency_penalty: Optional[float] = 0.0
	// user: Optional[str] = None
	ChatCompletion(ctx context.Context, model string, messages []openai.ChatCompletionMessage, temperature, topP, presencePenalty, frequencyPenalty float64, maxToken, n int, stop []string, user string, functions []openai.FunctionDefinition, functionCall any) (res openai.ChatCompletionResponse, status int, err error)
	// ChatCompletionStream 聊天处理流传输
	ChatCompletionStream(ctx context.Context, model string, messages []openai.ChatCompletionMessage, temperature float64, topP, presencePenalty, frequencyPenalty float64, maxToken, n int, stop []string, user string, functions []openai.FunctionDefinition, functionCall any) (stream *openai.ChatCompletionStream, status int, err error)
	// Models 模型列表
	Models(ctx context.Context) (res []openai.Model, err error)
	// Embeddings 创建图片
	Embeddings(ctx context.Context, model string, documents any) (res openai.EmbeddingResponse, err error)
	// ModeRations 检测量否有不当内容
	ModeRations(ctx context.Context, model, input string) (res openai.ModerationResponse, err error)
	// CreateImage 创建图片
	// Deprecated: use CreateSdImage
	// 暂时还无法使用
	CreateImage(ctx context.Context, prompt, size, format string) (res []openai.ImageResponseDataInner, err error)
	// CheckLength 验证Token是否超过相应长度
	CheckLength(ctx context.Context, prompt string, maxToken int) (tokenNum int, err error)
	// CreateAndGetSdImage 调用stable diffusion 文字生成图片并同时获取图片生成过程
	CreateAndGetSdImage(ctx context.Context, prompt, negativePrompt, samplerIndex string, steps int) (res <-chan Txt2ImgResult, err error)
}

type chatGPTToken struct {
	Token          string `json:"token"`
	OrganizationId string `json:"organizationId"`
}

type service struct {
	logger                                                                     log.Logger
	host, chatGPTHost, chatGPTToken, chatGPTModel, chatGPTOrgId, sdHost, svcIp string
	svcPort                                                                    int
	model, authToken                                                           string
	opts                                                                       []kithttp.ClientOption
	debug                                                                      bool
	sdImg                                                                      sync.Mutex
	rdb                                                                        redis.UniversalClient
	sdApiRedisKey                                                              string
	alarmSvc                                                                   alarm.Service
	chatGPTTokens                                                              []chatGPTToken
}

func (s *service) CancelFineTuningJob(ctx context.Context, modelName, jobId string) (err error) {
	client, _ := s.getClient(ctx, modelName)
	_, err = client.CancelFineTuningJob(ctx, jobId)
	if err != nil {
		err = errors.Wrap(err, "cancel fine tune")
		return
	}
	return
}

func (s *service) ListFineTune(ctx context.Context, modelName string) (res openai.FineTuneList, err error) {
	client, _ := s.getClient(ctx, modelName)
	tunes, err := client.ListFineTunes(ctx)
	if err != nil {
		err = errors.Wrap(err, "list fine tunes")
		return
	}
	return tunes, nil
}

func (s *service) RetrieveFineTuningJob(ctx context.Context, jobId string) (res openai.FineTuningJob, err error) {
	client, _ := s.getClient(ctx, openai.GPT3Dot5Turbo)
	tune, err := client.RetrieveFineTuningJob(ctx, jobId)
	if err != nil {
		err = errors.Wrap(err, "get fine tune")
		return
	}
	return tune, nil
}

func (s *service) UploadFile(ctx context.Context, modelName, fileName, filePath, purpose string) (res openai.File, err error) {
	client, _ := s.getClient(ctx, modelName)

	file, err := client.CreateFile(ctx, openai.FileRequest{
		FileName: fileName,
		FilePath: filePath,
		Purpose:  "fine-tune",
	})
	if err != nil {
		err = errors.Wrap(err, "create file")
		return
	}
	b, _ := json.Marshal(file)
	fmt.Println(string(b))
	return file, nil
}

func (s *service) CreateFineTuningJob(ctx context.Context, req openai.FineTuningJobRequest) (res openai.FineTuningJob, err error) {
	client, _ := s.getClient(ctx, req.Model)

	tune, err := client.CreateFineTuningJob(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "create fine tune")
		return
	}
	//if tune.Status == "failed" {
	//	err = errors.New(tune.Error)
	//	return err
	//}
	b, _ := json.Marshal(tune)
	fmt.Println(string(b))
	return tune, nil
}

func (s *service) ModeRations(ctx context.Context, model, input string) (res openai.ModerationResponse, err error) {
	client, _ := s.getClient(ctx, model)
	res, err = client.Moderations(ctx, openai.ModerationRequest{
		Input: input, Model: model,
	})
	return res, err
}

func (s *service) Embeddings(ctx context.Context, model string, documents any) (res openai.EmbeddingResponse, err error) {
	if isOpenAiEmbeddingModel(model) {
		client, _ := s.getClient(ctx, model)
		res, err = client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
			Input: documents,
			Model: openai.AdaEmbeddingV2, // 暂时写死
		})
		if err != nil {
			err = errors.Wrap(err, "Embeddings")
			return openai.EmbeddingResponse{}, err
		}
		return res, nil
	}

	tgt, _ := url.Parse(fmt.Sprintf("%s/v1/embeddings", s.host))
	ep := kithttp.NewClient(http.MethodPost, tgt, kithttp.EncodeJSONRequest, decodeJsonResponse(&res), s.opts...).Endpoint()
	_, err = ep(ctx, map[string]any{
		"input": documents,
		"model": model,
	})
	if err != nil {
		err = errors.Wrap(err, "LocalAIEmbeddings")
		return openai.EmbeddingResponse{}, err
	}
	return res, nil
}

func (s *service) CreateSdImageV1(ctx context.Context, prompt, negativePrompt, samplerIndex string, steps int) (res string, err error) {

	return
}

func (s *service) GetImageProgress(ctx context.Context, idTask string, idLivePreview int) (res []byte, err error) {
	type req struct {
		IdTask        string `json:"id_task"`
		IdLivePreview int    `json:"id_live_preview"`
	}

	var resData ImageProgress

	tgt, _ := url.Parse(fmt.Sprintf("%s/sdapi/v1/progress", s.host))
	ep := kithttp.NewClient(http.MethodGet, tgt, kithttp.EncodeJSONRequest, decodeJsonResponse(&resData), s.opts...).Endpoint()

	_, err = ep(ctx, req{
		IdTask:        idTask,
		IdLivePreview: idLivePreview,
	})
	if err != nil {
		err = errors.Wrap(err, "GetImageProgress")
		return nil, err
	}
	fmt.Println(resData)
	return
}

func (s *service) ChatCompletionPaasStream(ctx context.Context, model string, messages []openai.ChatCompletionMessage, temperature, topP float64, maxToken int) (stream *openai.ChatCompletionStream, status int, err error) {
	httpClient := http.DefaultClient
	if s.debug {
		httpClient = &http.Client{
			Transport: &proxyRoundTripper{},
		}
	}
	client := openai.NewClientWithConfig(openai.ClientConfig{
		BaseURL:            fmt.Sprintf("%s/langchain/v1", s.host),
		EmptyMessagesLimit: 300,
		HTTPClient:         httpClient,
	})

	stream, err = client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:       model,
		MaxTokens:   maxToken,
		Messages:    messages,
		Temperature: float32(temperature),
		TopP:        float32(topP),
		Stream:      true,
	})
	if err != nil {
		er := &openai.RequestError{}
		if errors.As(err, &er) {
			status = er.HTTPStatusCode
			err = er.Err
		}
	}

	return stream, status, err
}

func (s *service) CheckLength(ctx context.Context, prompt string, maxToken int) (tokenNum int, err error) {
	tgt, _ := url.Parse(fmt.Sprintf("%s/worker/count_token", s.host))

	ep := kithttp.NewClient(http.MethodPost, tgt, kithttp.EncodeJSONRequest, func(ctx context.Context, response *http.Response) (response1 interface{}, e error) {
		if response.StatusCode != http.StatusOK {
			var res ApiErrResponse
			if err = json.NewDecoder(response.Body).Decode(&res); err != nil {
				return nil, err
			}
			return nil, errors.New(res.Message)
		}

		var res struct {
			Count     int `json:"count"`
			ErrorCode int `json:"error_code"`
		}
		if err = json.NewDecoder(response.Body).Decode(&res); err != nil {
			return nil, err
		}

		return res.Count, nil
	}, s.opts...).Endpoint()

	type req struct {
		Prompt string `json:"prompt"`
	}

	if res, err := ep(ctx, req{
		Prompt: prompt,
	}); err == nil {
		tokenNum = res.(int)
	}

	var contextLength = 2048 // TODO 调用接口获得

	if tokenNum+maxToken > contextLength {
		//return tokenNum, errors.New(fmt.Sprintf("token num %d + max token %d > context length %d", tokenNum, maxToken, contextLength))
	}

	return
}

func (s *service) CreateImage(ctx context.Context, prompt, size, format string) (res []openai.ImageResponseDataInner, err error) {
	c := openai.NewClientWithConfig(openai.ClientConfig{
		BaseURL:            fmt.Sprintf("%s/v1", s.host),
		EmptyMessagesLimit: 300,
		//HTTPClient:         http.DefaultClient,
		HTTPClient: &http.Client{
			Transport: &proxyRoundTripper{},
		},
	})

	//translation, err := c.CreateTranslation(ctx, openai.AudioRequest{
	//	Model:    openai.Whisper1,
	//	FilePath: "~/Downloads/langchain.wav",
	//})
	//if err != nil {
	//	return nil, err
	//}
	//
	//fmt.Println(translation.Text)

	image, err := c.CreateImage(ctx, openai.ImageRequest{
		Prompt:         prompt,
		N:              1,
		Size:           size,
		ResponseFormat: format,
	})
	if err != nil {
		err = errors.Wrap(err, "CreateImage")
		return nil, err
	}

	return image.Data, nil
}

func (s *service) Models(ctx context.Context) (res []openai.Model, err error) {
	client, _ := s.getClient(ctx, "")
	models, err := client.ListModels(ctx)
	if err != nil {
		err = errors.Wrap(err, "ListModels")
		return nil, err
	}

	return models.Models, nil
}

func (s *service) getClient(ctx context.Context, model string) (*openai.Client, string) {
	logger := log.With(s.logger, "method", "getClient")
	httpClient := http.DefaultClient
	if s.debug {
		httpClient = &http.Client{
			Transport: &proxyRoundTripper{},
		}
	}
	ran := rand.Intn(len(s.chatGPTTokens))
	token := s.chatGPTTokens[ran].Token
	_ = level.Info(logger).Log("traceId", ctx.Value("traceId"), "model", model, "ran", ran)
	if isOpenAiModel(model) || isOpenAiEmbeddingModel(model) {
		config := openai.DefaultConfig(token)
		config.BaseURL = s.chatGPTHost
		config.HTTPClient = httpClient
		return openai.NewClientWithConfig(config), model
	}

	platform, ok := ctx.Value(ContextKeyPlatform).(Platform)
	if !ok || platform == "" {
		apiKey := s.authToken
		if key, exists := ctx.Value(ContextKeyApiKey).(string); exists && key != "" {
			apiKey = key
		}
		config := openai.DefaultConfig(apiKey)
		config.BaseURL = fmt.Sprintf("%s/v1", s.host)
		config.HTTPClient = httpClient
		return openai.NewClientWithConfig(config), model
	}
	config := openai.DefaultConfig(token)
	config.BaseURL = s.chatGPTHost
	config.HTTPClient = httpClient
	return openai.NewClientWithConfig(config), model
}

func (s *service) ChatCompletionStream(ctx context.Context, model string, messages []openai.ChatCompletionMessage, temperature, topP, presencePenalty, frequencyPenalty float64, maxToken, n int, stop []string, user string, functions []openai.FunctionDefinition, functionCall any) (stream *openai.ChatCompletionStream, status int, err error) {
	var client *openai.Client
	client, model = s.getClient(ctx, model)
	req := openai.ChatCompletionRequest{
		Model:            model,
		MaxTokens:        maxToken,
		Messages:         messages,
		Stream:           true,
		TopP:             float32(topP),
		Temperature:      float32(temperature),
		N:                n,
		PresencePenalty:  float32(presencePenalty),
		FrequencyPenalty: float32(frequencyPenalty),
		Stop:             stop,
		User:             user,
		Functions:        functions,
	}
	if functionCall != nil {
		if callStr, ok := functionCall.(string); ok && callStr == "" {
			req.FunctionCall = nil
		} else {
			req.FunctionCall = functionCall
		}
	}
	stream, err = client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		var er *openai.APIError
		if errors.As(err, &er) {
			status = er.HTTPStatusCode
			err = er
		}
		return
	}

	return stream, status, nil
}

func (s *service) ChatCompletion(ctx context.Context, model string, messages []openai.ChatCompletionMessage, temperature, topP, presencePenalty, frequencyPenalty float64, maxToken, n int, stop []string, user string, functions []openai.FunctionDefinition, functionCall any) (res openai.ChatCompletionResponse, status int, err error) {
	var client *openai.Client
	client, model = s.getClient(ctx, model)

	req := openai.ChatCompletionRequest{
		Model:            model,
		MaxTokens:        maxToken,
		Messages:         messages,
		Stream:           false,
		TopP:             float32(topP),
		Temperature:      float32(temperature),
		N:                n,
		PresencePenalty:  float32(presencePenalty),
		FrequencyPenalty: float32(frequencyPenalty),
		Stop:             stop,
		User:             user,
		Functions:        functions,
	}
	if functionCall != nil {
		if callStr, ok := functionCall.(string); ok && callStr == "" {
			req.FunctionCall = nil
		} else {
			req.FunctionCall = functionCall
		}
	}
	res, err = client.CreateChatCompletion(ctx, req)
	if err != nil {
		er := &openai.RequestError{}
		//apiErr := &openai.APIError{}
		if errors.As(err, &er) {
			status = er.HTTPStatusCode
			err = errors.New(fmt.Sprintf("%s %s", er.Error(), "可能是Token超过该模型的最大token限制了"))
		}

		return res, status, err
	}
	return
}

func (s *service) CreateAndGetSdImage(ctx context.Context, prompt, negativePrompt, samplerIndex string, steps int) (res <-chan Txt2ImgResult, err error) {
	logger := log.With(s.logger, "method", "CreateAndGetSdImage")
	var sdApiHost string               //当前获取到的sdApiHost
	var dot = make(chan Txt2ImgResult) //返回内容的channel

	s.sdImg.Lock()
	defer func() {
		s.sdImg.Unlock()
		//将sdApiHost地址放回redis队列
		s.rdb.RPush(context.Background(), s.sdApiRedisKey, sdApiHost)
		_ = level.Info(logger).Log("CreateAndGetSdImage", "over")
	}()

	//从redis队列取出一个sdApi地址
	//BLPop 移出并获取列表的第一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止
	//BLPop 命令的返回值是一个包含两个元素的字符串切片。第一个元素是弹出的列表的名称，第二个元素是弹出的元素值。
	resRedis, err := s.rdb.BLPop(context.Background(), time.Minute*3, s.sdApiRedisKey).Result()
	if err != nil {
		_ = level.Error(logger).Log("s.rdb.LPop", s.sdApiRedisKey, "err", err.Error())

		//如果报错，则返回结束
		txt2ImgResult := Txt2ImgResult{
			Finish: true,
			Error:  err,
		}
		dot <- txt2ImgResult
		close(dot)

		//发送告警
		_ = s.alarmSvc.Push(ctx, "从Redis获取SdApiHost失败", err.Error(), "CreateAndGetSdImage", alarm.LevelError, 2)

		return nil, err
	}

	_ = level.Info(logger).Log("rdb", "BLPop", "sdApiRedisKey", s.sdApiRedisKey, "list", resRedis[0], "val", resRedis[1])

	//获取到的 sdApiHost 地址
	sdApiHost = resRedis[1]

	//参数默认值
	if steps == 0 {
		steps = 30
	}
	if strings.EqualFold(samplerIndex, "") {
		samplerIndex = "Euler a"
	}

	//创建 入参
	req := Txt2ImgRequest{
		Steps:          steps,
		Prompt:         prompt,
		NegativePrompt: negativePrompt,
		SamplerIndex:   samplerIndex,
	}

	reqStr, _ := json.Marshal(req)

	_ = level.Info(logger).Log("req", string(reqStr))

	parentCtx, cancel := context.WithCancel(context.Background())

	//通知拉流接口开始
	progressStart := make(chan bool)

	// 启动 A 请求
	go func() {
		createSdImageRep := make(chan CreateSdImage)
		go s.createSdImage(context.Background(), sdApiHost, req, createSdImageRep, progressStart)
		select {
		case <-parentCtx.Done():
			fmt.Println("createSdImage 请求结束")
			return
		case result := <-createSdImageRep:
			fmt.Println("createSdImage 请求返回结果:", result.SdSuccessRep.Info, "sdApiHost", sdApiHost)
			//将最终的图片数据流传进去  //todo 未来多图如何兼容

			if result.Error != nil {
				//接口请求失败
				txt2ImgResult := Txt2ImgResult{
					Finish: true,
					Error:  err,
				}
				dot <- txt2ImgResult
				close(dot)
			} else {
				//接口请求成功
				if len(result.SdSuccessRep.Images) > 0 {
					resProgress := ImageProgress{}
					resProgress.CurrentImage = result.SdSuccessRep.Images[0]
					txt2ImgResult := Txt2ImgResult{
						Finish:        true,
						ImageProgress: resProgress,
					}
					dot <- txt2ImgResult
					close(dot)
				}
			}
			// 取消父上下文，停止 B 请求的循环
			cancel()
		}
	}()

	<-progressStart

	go func() {
		// 循环发起 B 请求
		for {
			select {
			case <-parentCtx.Done():
				fmt.Println("getProgress 请求停止循环")
				cancel() //中断当前协程
				return
			default:
				resProgress, err := s.getProgress(context.Background(), sdApiHost)
				//fmt.Println("getProgress 请求结果", resProgress.State, err)
				time.Sleep(time.Millisecond * 50)
				if err == nil && resProgress.CurrentImage != "" && resProgress.State.SamplingStep > 5 && resProgress.State.SamplingStep < resProgress.State.SamplingSteps {
					txt2ImgResult := Txt2ImgResult{
						Finish:        false,
						ImageProgress: resProgress,
					}
					dot <- txt2ImgResult
				}
			}
		}
	}()
	return dot, nil
}

func (s *service) createSdImage(ctx context.Context, sdApi string, req Txt2ImgRequest, res chan<- CreateSdImage, progressStart chan<- bool) {
	createRep := new(CreateSdImage)
	// 发送 创建请求
	opts := s.opts
	opts = append(opts, kithttp.ClientBefore(func(ctx context.Context, request *http.Request) context.Context {
		//创建请求开始后，通知拉流接口开始
		time.Sleep(time.Millisecond * 100)
		progressStart <- true
		return ctx
	}), kithttp.ClientAfter(func(ctx context.Context, response *http.Response) context.Context {
		return ctx
	}), kithttp.ClientFinalizer(func(ctx context.Context, err error) {

	}))

	//生成图片返回成功信息
	sdSuccessRep := new(SdSuccessRep)
	u, _ := url.Parse(fmt.Sprintf("%s/sdapi/v1/txt2img", sdApi))
	c := kithttp.NewClient("POST", u, kithttp.EncodeJSONRequest, decodeJsonResponse(sdSuccessRep), opts...).Endpoint()
	_, err := c(ctx, req)
	if err != nil {
		fmt.Println("createSdImage err:", err)
		err = errors.Wrap(err, "createSdImage")
		createRep.Error = err
		res <- *createRep
	} else {
		createRep.SdSuccessRep = *sdSuccessRep
		res <- *createRep
	}
	return
}

func (s *service) getProgress(ctx context.Context, sdApi string) (res ImageProgress, err error) {
	progressRep := new(ImageProgress)
	u, _ := url.Parse(fmt.Sprintf("%s/sdapi/v1/progress", sdApi))
	c := kithttp.NewClient("GET", u, kithttp.EncodeJSONRequest, decodeJsonResponse(progressRep), []kithttp.ClientOption{}...).Endpoint()
	_, err = c(ctx, nil)
	if err != nil {
		fmt.Println("getProgress err:", err)
		err = errors.Wrap(err, "getProgress")
		return
	} else {
		//fmt.Println("getProgress success")
		res = *progressRep
	}
	return
}

func base64ToMultipartFile(base64Str string, fileName string) (multipart.File, error) {
	// 解码Base64字符串
	dec, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}
	// 创建临时文件
	file, err := os.CreateTemp("/tmp", fileName) // 可以根据实际需求指定文件名和扩展名
	if err != nil {
		return nil, err
	}

	// 将解码后的数据写入临时文件
	_, err = file.Write(dec)
	if err != nil {
		_ = file.Close()
		_ = os.Remove(file.Name())
		return nil, err
	}

	// 将文件指针定位到文件开头
	_, err = file.Seek(0, 0)
	if err != nil {
		_ = file.Close()
		_ = os.Remove(file.Name())
		return nil, err
	}
	return file, nil
}

func decodeJsonResponse(data interface{}) func(ctx context.Context, res *http.Response) (response interface{}, err error) {
	return func(ctx context.Context, res *http.Response) (response interface{}, err error) {
		if res.StatusCode == 422 {
			body, _ := io.ReadAll(res.Body)
			return res, fmt.Errorf("http status code is %d, body %s", res.StatusCode, string(body))
		}
		/*		if res.StatusCode != 200 {
				body, _ := io.ReadAll(res.Body)
				return res, fmt.Errorf("http status code is %d, body %s", res.StatusCode, string(body))
			}*/
		if data == nil {
			return res, nil
		}
		if err = json.NewDecoder(res.Body).Decode(data); err != nil {
			return res, errors.Wrap(err, "json decode")
		}
		return res, nil
	}
}

func New(logger log.Logger, cfg Config, opts []kithttp.ClientOption, rdb redis.UniversalClient, alarmSvc alarm.Service) Service {
	logger = log.With(logger, "api", "fastchat")
	chatGPTTokens := parseTokens(cfg.OpenAiToken)
	return &service{
		logger:        logger,
		host:          cfg.PaasGptEndpoint,
		model:         cfg.PaasGptModel,
		sdHost:        cfg.SdHost,
		svcIp:         cfg.SvcIp,
		svcPort:       cfg.SvcPort,
		opts:          opts,
		chatGPTHost:   cfg.OpenAiEndpoint,
		chatGPTToken:  cfg.OpenAiToken,
		chatGPTOrgId:  cfg.OpenAiOrgId,
		debug:         cfg.Debug,
		rdb:           rdb,
		sdApiRedisKey: cfg.SdRedisKey, //sd api list 队列
		alarmSvc:      alarmSvc,
		authToken:     cfg.AuthToken,
		chatGPTTokens: chatGPTTokens,
	}
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

func TokensNumFromMessages(messages []openai.ChatCompletionMessage, model string) (numTokens int) {
	if strings.EqualFold(model, "") {
		model = openai.GPT3Dot5Turbo
	}
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("EncodingForModel: %v", err)
		fmt.Println(err)
		return
	}

	var tokensPerMessage int
	var tokensPerName int
	if model == "gpt-3.5-turbo-0301" || model == "gpt-3.5-turbo" {
		tokensPerMessage = 4
		tokensPerName = -1
	} else if model == "gpt-4-0314" || model == "gpt-4" {
		tokensPerMessage = 3
		tokensPerName = 1
	} else {
		fmt.Println("Warning: model not found. Using cl100k_base encoding.")
		tokensPerMessage = 3
		tokensPerName = 1
	}

	for _, message := range messages {
		numTokens += tokensPerMessage
		numTokens += len(tkm.Encode(message.Content, nil, nil))
		numTokens += len(tkm.Encode(message.Role, nil, nil))
		numTokens += len(tkm.Encode(message.Name, nil, nil))
		if message.Name != "" {
			numTokens += tokensPerName
		}
	}
	numTokens += 3
	return numTokens
}

// TokenizerGetWord 只支持英文
func TokenizerGetWord(text string, model string) []string {
	if strings.EqualFold(model, "") {
		model = openai.GPT3Dot5Turbo
	}
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("EncodingForModel: %v", err)
		fmt.Println(err)
		return nil
	}
	var words []string
	tokens := tkm.Encode(text, nil, nil)
	for _, v := range tokens {
		words = append(words, tkm.Decode([]int{v}))
	}
	return words
}

// InputToString 将输入转换为字符串
func InputToString(input any) (res string, err error) {
	tk, err := tiktoken.EncodingForModel(openai.GPT3Dot5Turbo)
	if err != nil {
		return "", err
	}
	switch input.(type) {
	case string:
		return input.(string), nil
	case []string:
		return tk.Decode(tk.Encode(strings.Join(input.([]string), " "), nil, nil)), nil
	case [][]int:
		var allInput []int
		for _, val := range input.([][]int) {
			for _, num := range val {
				allInput = append(allInput, num)
			}
		}
		return tk.Decode(allInput), nil
	}
	return res, nil
}

func isOpenAiModel(model string) bool {
	return util.StringContainsArray([]string{
		openai.GPT3Dot5Turbo,
		openai.GPT432K0613,
		openai.GPT432K0314,
		openai.GPT432K,
		openai.GPT40613,
		openai.GPT40314,
		openai.GPT4,
		openai.GPT3Dot5Turbo0613,
		openai.GPT3Dot5Turbo0301,
		openai.GPT3Dot5Turbo16K,
		openai.GPT3Dot5Turbo16K0613,
		openai.GPT3Dot5Turbo,
		openai.GPT3Dot5TurboInstruct,
		openai.GPT3Davinci,
		openai.GPT3Davinci002,
	}, model)
}

func isOpenAiEmbeddingModel(model string) bool {
	return util.StringInArray([]string{
		"text-similarity-ada-001",
		"text-similarity-babbage-001",
		"text-similarity-curie-001",
		"text-similarity-davinci-001",
		"text-search-ada-doc-001",
		"text-search-ada-query-001",
		"text-search-babbage-doc-001",
		"text-search-babbage-query-001",
		"text-search-curie-doc-001",
		"text-search-curie-query-001",
		"text-search-davinci-doc-001",
		"text-search-davinci-query-001",
		"code-search-ada-code-001",
		"code-search-ada-text-001",
		"code-search-babbage-code-001",
		"code-search-babbage-text-001",
		"text-embedding-ada-002",
	}, model)
}

func parseTokens(cfgToken string) []chatGPTToken {
	var chatGPTTokens []chatGPTToken

	parseToken := func(token string) {
		split := strings.SplitN(token, ":", 2)
		chatToken := chatGPTToken{OrganizationId: "", Token: split[0]}

		if len(split) > 1 {
			chatToken.OrganizationId = split[0]
			chatToken.Token = split[1]
		}

		chatGPTTokens = append(chatGPTTokens, chatToken)
	}

	for _, token := range strings.Split(cfgToken, ",") {
		parseToken(token)
	}

	return chatGPTTokens
}
