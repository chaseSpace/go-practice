package testfiles

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
)

func TestPlan4StaticPrompt(t *testing.T) {
	loadLearningEnv(t)
	promptMessages := []*schema.Message{
		schema.SystemMessage("你是 Go 助教。用中文、最多三句话回答；不知道就明确说不知道。"),
		schema.UserMessage("解释 Eino 中的 Tool 是什么。"),
	}
	if len(promptMessages) != 2 {
		t.Fatalf("prompt message count = %d, want 2", len(promptMessages))
	}
	if promptMessages[0].Role != schema.System || promptMessages[1].Role != schema.User {
		t.Fatalf("unexpected prompt roles: %s, %s", promptMessages[0].Role, promptMessages[1].Role)
	}
	if os.Getenv("EINO_RUN_LIVE_TESTS") != "1" {
		t.Skip("set EINO_RUN_LIVE_TESTS=1 in testfiles/.env to make a paid network request")
	}
	apiKey := llmEnv("LLM_API_KEY", "OPENAI_API_KEY")
	modelName := llmEnv("LLM_MODEL", "OPENAI_MODEL")
	if apiKey == "" || modelName == "" {
		t.Skip("LLM_API_KEY and LLM_MODEL are required for the live test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	cm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  apiKey,
		Model:   modelName,
		BaseURL: llmEnv("LLM_ENDPOINT", "OPENAI_BASE_URL"),
	})
	if err != nil {
		t.Fatalf("construct chat model: %v", err)
	}

	ans, err := cm.Generate(ctx, promptMessages)
	if err != nil {
		t.Fatalf("generate prompt answer: %v", err)
	}
	if ans == nil || ans.Content == "" {
		t.Fatalf("expected a non-empty assistant message, got %#v", ans)
	}
	t.Logf("Answer: %s", ans.Content)
}
