package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/tests"
	"github.com/go-kit/log"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"strings"
	"testing"
)

var logger log.Logger

func initSvc() Service {
	rdb, apiSvc, err := tests.Init()
	if err != nil {
		panic(err)
	}

	return New(logger, "", tests.Store, rdb, apiSvc)
}

func TestService_FunctionCall(t *testing.T) {
	var functions = []llms.FunctionDefinition{
		{
			Name:        "getCurrentWeather",
			Description: "Get the current weather in a given location",
			Parameters:  json.RawMessage(`{"type": "object", "properties": {"location": {"type": "string", "description": "The city and state, e.g. San Francisco, CA"}, "unit": {"type": "string", "enum": ["celsius", "fahrenheit"]}}, "required": ["location"]}`),
		},
	}

	llm, err := openai.NewChat(
		openai.WithModel("gpt-3.5-turbo"),
		openai.WithBaseURL("http://aigc-admin:8080/v1"),
		//openai.WithToken("sk-4W6vS2nG8mC1pT3kX5rH7fJ9bQ0dZ4lY2cV1xN3aM9gB6qD8"),
		openai.WithToken("sk-4W6vS2nG8mC1pT3kX5rH7fJ9bQ0dZ4lY2cV1xN3aM9gB6qD8"),
	)
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	completion, err := llm.Call(ctx, []schema.ChatMessage{
		schema.HumanChatMessage{Content: "What is the weather like in Boston?"},
	}, llms.WithFunctions(functions),
		llms.WithMaxTokens(2048))
	if err != nil {
		t.Error(err)
		return
	}

	if completion.FunctionCall != nil {
		fmt.Printf("Function call: %v\n", completion.FunctionCall)
		if strings.EqualFold(completion.FunctionCall.Name, "getCurrentWeather") {
			var args map[string]string
			_ = json.Unmarshal([]byte(completion.FunctionCall.Arguments), &args)
			res, err := getCurrentWeather(args["location"], "")
			if err != nil {
				t.Error(err)
				return
			}
			t.Log(res)
		}
	}

	b, _ := json.Marshal(completion)

	t.Log(string(b))
	t.Log(completion.GetContent())
}

func getCurrentWeather(location string, unit string) (string, error) {
	weatherInfo := map[string]interface{}{
		"location":    location,
		"temperature": "72",
		"unit":        unit,
		"forecast":    []string{"sunny", "windy"},
	}
	b, err := json.Marshal(weatherInfo)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
