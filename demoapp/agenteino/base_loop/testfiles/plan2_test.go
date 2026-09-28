package testfiles

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
)

func TestPlan2CanConstructModel(t *testing.T) {
	loadLearningEnv(t)
	// NewChatModel only constructs a client; it does not send a request.
	chatModel, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		APIKey: "test-key-not-sent-over-network",
		Model:  "test-model",
	})
	if err != nil {
		t.Fatalf("construct chat model: %v", err)
	}
	if chatModel == nil {
		t.Fatal("chat model is nil")
	}
}

func TestPlan2FirstLLMCall(t *testing.T) {
	loadLearningEnv(t)
	if os.Getenv("EINO_RUN_LIVE_TESTS") != "1" {
		t.Skip("set EINO_RUN_LIVE_TESTS=1 in testfiles/.env to make a paid network request")
	}
	apiKey := llmEnv("LLM_API_KEY", "OPENAI_API_KEY")
	modelName := llmEnv("LLM_MODEL", "OPENAI_MODEL")
	if apiKey == "" || modelName == "" {
		t.Skip("LLM_API_KEY and LLM_MODEL are required for the live test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  apiKey,
		Model:   modelName,
		BaseURL: llmEnv("LLM_ENDPOINT", "OPENAI_BASE_URL"),
	})
	if err != nil {
		t.Fatalf("construct chat model: %v", err)
	}

	answer, err := chatModel.Generate(ctx, []*schema.Message{
		schema.UserMessage("请只回复：Eino 学习测试成功"),
	})
	if err != nil {
		t.Fatalf("generate answer: %v", err)
	}
	if answer == nil || !strings.Contains(answer.Content, "Eino 学习测试成功") {
		t.Fatalf("expected a non-empty assistant message, got %#v", answer)
	}
}
