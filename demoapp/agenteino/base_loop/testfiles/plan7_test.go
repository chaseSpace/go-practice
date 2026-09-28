package testfiles

import (
	"context"
	"testing"
)

func TestPlan7FirstTool(t *testing.T) {
	addTool, err := newAddTool()
	if err != nil {
		t.Fatalf("create add tool: %v", err)
	}
	info, err := addTool.Info(context.Background())
	if err != nil {
		t.Fatalf("read tool info: %v", err)
	}
	if info.Name != "add" || info.ParamsOneOf == nil {
		t.Fatalf("unexpected tool info: %#v", info)
	}

	result, err := addTool.InvokableRun(context.Background(), `{"a":2,"b":3}`)
	if err != nil {
		t.Fatalf("run add tool: %v", err)
	}
	if result != "5" {
		t.Fatalf("tool result = %q, want 5", result)
	}
}
