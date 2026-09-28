package testfiles

import (
	"runtime"
	"testing"

	"github.com/cloudwego/eino/schema"
)

func TestPlan1ProjectSetup(t *testing.T) {
	message := schema.UserMessage("你好，Eino")
	if message.Role != schema.User {
		t.Fatalf("unexpected role: %s", message.Role)
	}
	if message.Content != "你好，Eino" {
		t.Fatalf("unexpected content: %q", message.Content)
	}
	t.Logf("Go runtime: %s; Eino schema import works", runtime.Version())
}
