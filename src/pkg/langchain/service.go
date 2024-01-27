package langchain

import (
	"context"
	"fmt"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/util"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type ChainType string

const (
	// ChainTypeStuff
	// 事物文档链（“stuff”这里指的是“填充”或“装满”）是最直接的文档链之一。它接收一个文档列表，将所有文档插入到一个提示中，并将该提示传递给一个LLM模型。
	// 这个链式结构非常适用于文档较小且大多数调用只传递少量文档的应用场景。
	ChainTypeStuff ChainType = "stuff"
	// ChainTypeRefine
	// 完善文档链通过对输入的各个文档进行循环，以迭代更新回答来构建响应。对于每个文档，它会将所有非文档输入、当前文档以及最新的中间回答传递给LLM链，以获得一个新的回答。
	// 由于完善链一次只传递一个文档给LLM模型，因此它非常适用于需要分析大量文档超出模型上下文范围的任务。显而易见的权衡是，相比“装填文档链”等其他链式结构而言，该链会产生更多LLM调用。同时，在某些任务中，迭代方法可能较难实现。例如，在经常相互交叉引用的文档或需要从多个文档获取详细信息的任务中，完善链可能表现不佳。
	ChainTypeRefine ChainType = "refine"
	// ChainTypeMapReduce
	// 映射-归约文档链首先对每个文档单独应用LLM链（映射步骤），将链的输出视为一个新的文档。然后，它将所有新文档传递给一个独立的合并文档链，以获得单一的输出（归约步骤）。可以选择首先对映射后的文档进行压缩或折叠，以确保其适应合并文档链（通常会将其传递给LLM）。如果需要，此压缩步骤会进行递归处理。
	ChainTypeMapReduce ChainType = "map_reduce"
)

type Service interface {
	// Summary 生成内容摘要信息
	// chainType: ChainType
	// prompt: 提示词，如不填则为“写出本篇内容的简洁摘要”
	// filePath: 文件路径，可以是url, url 必须传内部地址，生产环境无法直接访问.com域名只能访问到.idc
	// modelName: 模型名称
	// maxTokens: 模型最大返回的Tokens数量
	// temperature, topP: 控制输出内容的随机性
	// streamingFunc: 回调函数，会把内容持续输出的到这个函数，结束后才会整体return 到 res
	Summary(ctx context.Context, chainType ChainType, prompt, filePath, modelName string, maxTokens int, temperature, topP float64, streamingFunc func(ctx context.Context, chunk []byte) error) (res map[string]any, err error)
}

type service struct {
	logger                log.Logger
	traceId               string
	llm                   *openai.LLM
	defaultAllowedSpecial []string
}

func (s *service) Summary(ctx context.Context, chainType ChainType, prompt, filePath, modelName string, maxTokens int, temperature, topP float64, streamingFunc func(ctx context.Context, chunk []byte) error) (res map[string]any, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	if strings.EqualFold(prompt, "") {
		prompt = "写出本篇内容的摘要"
	}

	ts := textsplitter.NewTokenSplitter()
	ts.ChunkSize = 2048 // 应该根据 maxTokens和 modelMaxTokens进行计算
	ts.ChunkOverlap = 256
	ts.ModelName = modelName
	ts.AllowedSpecial = []string{"\n\n", "\n"}

	f, filePath, err := getLocalFile(filePath)
	if err != nil {
		_ = level.Error(logger).Log("get", "file", "err", err.Error())
		return nil, errors.Wrap(err, "getLocalFile")
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	//fType, err := util.GetFileTypeName(f)
	//if err != nil {
	//	err = errors.Wrap(err, "util.GetFileTypeName")
	//	_ = level.Error(logger).Log("util", "GetFileTypeName", "err", err.Error())
	//	return
	//}
	fType := strings.ToLower(strings.TrimPrefix(filepath.Ext(filePath), "."))
	var docs []schema.Document
	switch fType {
	case "pdf":
		docs, err = loadPdfDocs(f)
		break
	case "csv":
		docs, err = loadCsvDocs(f)
		break
	case "html":
		docs, err = loadHtmlDocs(f)
		break
	case "txt":
		docs, err = loadTextDocs(f)
		break
	}
	if err != nil {
		_ = level.Error(logger).Log("load", "file", "fType", fType, "err", err.Error())
		return nil, err
	}
	var docLen int
	for _, doc := range docs {
		docLen += len(doc.PageContent)
	}

	if strings.Contains(modelName, "-16k") && docLen < 16000-maxTokens ||
		strings.EqualFold(modelName, "gpt-4") && docLen < 7000-maxTokens ||
		docLen < 3800-maxTokens {
		chainType = ChainTypeStuff
	} else if docLen > 24000 && docLen < 640000 {
		chainType = ChainTypeMapReduce
	} else {
		chainType = ChainTypeRefine
	}

	docsPageNum := len(docs)
	_ = level.Info(logger).Log("chainType", chainType, "docLen", docLen, "docs", len(docs), "fType", fType)

	var options []chains.ChainCallOption
	options = append(options, chains.WithModel(modelName),
		chains.WithMaxTokens(maxTokens), chains.WithTemperature(temperature),
		chains.WithTopP(topP))

	//options = append(options, chains.WithStopWords([]string{"\nstop:\n\n"}))

	if streamingFunc != nil {
		prevChar := "-"
		currentChar := "-"
		var n = 1
		streamBlock := func(targetN int, ctx context.Context, chunk []byte) error {
			prevChar = currentChar
			currentChar = string(chunk)
			if prevChar == "" && currentChar == "" {
				n++
			}
			if n == targetN {
				return streamingFunc(ctx, chunk)
			}
			return streamingFunc(ctx, nil)
		}

		options = append(options, chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			if chainType == ChainTypeStuff {
				return streamingFunc(ctx, chunk)
			}
			if chainType == ChainTypeRefine {
				return streamBlock(docsPageNum, ctx, chunk)
			}
			return streamBlock(2, ctx, chunk)
		}))
	}
	llm := s.llm
	llmChain := chains.NewLLMChain(llm, prompts.NewPromptTemplate(
		_stuffSummarizationTemplate, []string{"context", "prompt"},
	))

	var chainCall map[string]any
	switch chainType {
	case ChainTypeStuff:
		chainCall, err = chains.NewStuffDocuments(llmChain).Call(ctx, map[string]any{
			"input_documents": docs,
			"prompt":          prompt,
		}, options...)
		if err != nil {
			return nil, errors.Wrap(err, "chains.LoadStuffSummarization.Call")
		}
	case ChainTypeRefine:
		refineLLMChain := chains.NewLLMChain(llm, prompts.NewPromptTemplate(
			_refineSummarizationTemplate, []string{"existing_answer", "context"},
		))
		chainCall, err = chains.NewRefineDocuments(llmChain, refineLLMChain).Call(ctx, map[string]any{
			"input_documents": docs,
			"prompt":          prompt,
		}, options...)
		if err != nil {
			return nil, errors.Wrap(err, "chains.LoadRefineSummarization.Call")
		}
	case ChainTypeMapReduce:
		mapChain := chains.NewLLMChain(llm, prompts.NewPromptTemplate(
			_stuffSummarizationTemplate, []string{"context", "prompt"},
		))
		chainCall, err = chains.NewMapReduceDocuments(mapChain, chains.NewStuffDocuments(llmChain)).
			Call(ctx, map[string]any{
				"input_documents": docs,
				"prompt":          prompt,
			}, options...)
		if err != nil {
			return nil, errors.Wrap(err, "chains.LoadMapReduceSummarization.Call")
		}
	}

	return chainCall, nil
}

func New(logger log.Logger, traceId string, llm *openai.LLM) Service {
	return &service{logger: logger, traceId: traceId, llm: llm}
}

func loadPdfDocs(f *os.File) (res []schema.Document, err error) {
	info, err := f.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "file.Stat")
	}
	p := documentloaders.NewPDF(f, info.Size() /*, documentloaders.WithPassword("password")*/)
	docs, err := p.Load(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "documentloaders.NewPDF")
	}
	return docs, nil
}

func loadTextDocs(f *os.File) (res []schema.Document, err error) {
	p := documentloaders.NewText(f)
	docs, err := p.Load(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "documentloaders.NewText")
	}
	return docs, nil
}

func loadHtmlDocs(f *os.File) (res []schema.Document, err error) {
	p := documentloaders.NewHTML(f)
	docs, err := p.Load(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "documentloaders.NewHTML")
	}
	return docs, nil
}

func loadCsvDocs(f *os.File) (res []schema.Document, err error) {
	p := documentloaders.NewCSV(f)
	docs, err := p.Load(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "documentloaders.NewCSV")
	}
	return docs, nil
}

func getLocalFile(path string) (file *os.File, filePath string, err error) {
	filePath = path
	if !util.IsUrl(path) {
		file, err = os.Open(filepath.Clean(path))
		return
	}
	destDir := "/tmp"
	// 下载文件并保存
	if !strings.EqualFold(runtime.GOOS, "linux") {
		dir, err := os.Getwd()
		if err != nil {
			return nil, filePath, errors.Wrap(err, "Error getting current directory")
		}
		destDir = dir + string(os.PathSeparator) // 其他系统使用当前目录
	}
	u, err := url.Parse(path)
	if err != nil {
		err = errors.Wrap(err, "Error parsing URL")
		return
	}

	fileExt := strings.TrimPrefix(filepath.Ext(u.Path), ".")
	file, err = downloadFile(path, fmt.Sprintf("%s/%s.%s", destDir, uuid.New().String(), strings.ToLower(fileExt)))
	return
}

// 下载文件并保存
func downloadFile(URL string, filePath string) (file *os.File, err error) {
	response, err := http.Get(URL)
	if err != nil {
		err = errors.Wrap(err, "failed to download file")
		return
	}

	defer response.Body.Close()
	file, err = os.Create(filePath)
	if err != nil {
		err = fmt.Errorf("failed to create file: %w", err)
		return

	}
	if _, err = io.Copy(file, response.Body); err != nil {
		err = errors.Wrap(err, "io.Copy")
	}
	return
}
