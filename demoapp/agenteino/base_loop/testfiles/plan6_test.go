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

type studyAnswer struct {
	Answer     string  `json:"answer"`
	Confidence float64 `json:"confidence"`
}

func TestPlan6StructuredOutputParser(t *testing.T) {
	parser := schema.NewMessageJSONParser[studyAnswer](nil)
	result, err := parser.Parse(context.Background(), schema.AssistantMessage(
		`{"answer":"Eino 用 Go 构建 AI 应用。","confidence":0.8}`, nil,
	))
	if err != nil {
		t.Fatalf("parse JSON output: %v", err)
	}
	if result.Answer == "" || result.Confidence != 0.8 {
		t.Fatalf("unexpected structured result: %#v", result)
	}
}

// TestPlan6ParseRealChatReply 演示完整边界：获取真实 ChatModel 回复，
// 规范化后解析为 Go 结构体。它需要显式开启，因为会发起网络请求并可能消耗额度。
func TestPlan6ParseRealChatReply(t *testing.T) {
	loadLearningEnv(t)
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
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  apiKey,
		Model:   modelName,
		BaseURL: llmEnv("LLM_ENDPOINT", "OPENAI_BASE_URL"),
	})
	if err != nil {
		t.Fatalf("construct chat model: %v", err)
	}

	response, err := chatModel.Generate(ctx, []*schema.Message{
		schema.SystemMessage("你是 Eino 助教。只输出一个合法 JSON 对象，不要 Markdown 代码围栏或额外文字。"),
		schema.UserMessage(`用一句话解释 Eino 的 Component，并按此 schema 回复：{"answer":"非空字符串","confidence":0 到 1 之间的小数}`),
	})
	if err != nil {
		t.Fatalf("generate structured answer: %v", err)
	}
	if response == nil || response.Content == "" {
		t.Fatalf("expected a non-empty assistant message, got %#v", response)
	}

	parser := schema.NewMessageJSONParser[studyAnswer](nil)
	result, err := parser.Parse(ctx, schema.AssistantMessage(jsonOnlyContent(response.Content), nil))
	if err != nil {
		t.Fatalf("parse real model reply %q: %v", response.Content, err)
	}
	if result.Answer == "" || result.Confidence < 0 || result.Confidence > 1 {
		t.Fatalf("invalid parsed result: %#v", result)
	}
	t.Logf("parsed real reply: answer=%q confidence=%.2f", result.Answer, result.Confidence)
}

// jsonOnlyContent 处理常见的 Markdown JSON 代码围栏。
// 它只做最小规范化，不会修复非法 JSON；解析器仍会校验模型回复是否符合 studyAnswer。
func jsonOnlyContent(content string) string {
	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, "```") {
		return content
	}
	if lineEnd := strings.IndexByte(content, '\n'); lineEnd >= 0 {
		content = content[lineEnd+1:]
	}
	content = strings.TrimSpace(content)
	content = strings.TrimSuffix(content, "```")
	return strings.TrimSpace(content)
}
