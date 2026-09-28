package testfiles

import (
	"testing"

	"github.com/cloudwego/eino/schema"
)

func TestPlan3MessageRolesAndHistory(t *testing.T) {
	history := []*schema.Message{
		schema.SystemMessage("你是一名简洁的 Go 助教。"),
		schema.UserMessage("什么是 Eino？"),
		schema.AssistantMessage("Eino 是 Go 的 AI 应用开发框架。", nil),
		schema.UserMessage("再用一句话说明它的用途。"),
	}
	wantRoles := []schema.RoleType{schema.System, schema.User, schema.Assistant, schema.User}
	for i, message := range history {
		if message.Role != wantRoles[i] {
			t.Fatalf("message %d role = %s, want %s", i, message.Role, wantRoles[i])
		}
	}
	if history[2].Content == "" {
		t.Fatal("an assistant reply must be kept for the next turn's context")
	}
}
