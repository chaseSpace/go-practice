package testfiles

import (
	"context"
	"strings"
	"testing"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func TestPlan5PromptTemplate(t *testing.T) {
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一名 {language} 助教。"),
		schema.UserMessage("请简洁回答：{question}"),
	)
	messages, err := template.Format(context.Background(), map[string]any{
		"language": "Go",
		"question": "Eino 的 ChatModel 做什么？",
	})
	if err != nil {
		t.Fatalf("format prompt: %v", err)
	}
	if len(messages) != 2 {
		t.Fatalf("message count = %d, want 2", len(messages))
	}
	if !strings.Contains(messages[1].Content, "ChatModel") {
		t.Fatalf("question was not rendered: %q", messages[1].Content)
	}
}
