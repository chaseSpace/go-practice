package testfiles

import (
	"context"
	"testing"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
)

func TestPlan8BindToolToModel(t *testing.T) {
	loadLearningEnv(t)
	chatModel, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		APIKey:  llmEnv("LLM_API_KEY", ""),
		Model:   llmEnv("LLM_MODEL", ""),
		BaseURL: llmEnv("LLM_ENDPOINT", ""),
	})
	if err != nil {
		t.Fatalf("construct chat model: %v", err)
	}
	addTool, err := newAddTool()
	if err != nil {
		t.Fatalf("create add tool: %v", err)
	}
	info, err := addTool.Info(context.Background())
	if err != nil {
		t.Fatalf("read tool info: %v", err)
	}
	withTools, err := chatModel.WithTools([]*schema.ToolInfo{info})
	if err != nil {
		t.Fatalf("bind tool to model: %v", err)
	}
	if withTools == nil {
		t.Fatal("expected a model variant with tools")
	}

	arr := []*schema.Message{
		schema.SystemMessage("你是一个问答工具，可以利用已有的工具回复问题"),
		schema.UserMessage("3和11的和是多少"),
	}
	a, err := withTools.Generate(context.Background(), arr)
	if err != nil {
		t.Fatalf("generate tool-call request: %v", err)
	}
	if a == nil {
		t.Fatal("expected a model message")
	}
	// 得到调用意图
	t.Logf("Answer: %s -- toolcalls: %+v", a.Content, a.ToolCalls)

	if len(a.ToolCalls) == 0 {
		t.Fatal("expected a tool call")
	}
	if a.ToolCalls[0].Type != "function" {
		t.Fatalf("expected a function call, got %s", a.ToolCalls[0].Type)
	}
	tc := a.ToolCalls[0]
	tr, err := addTool.InvokableRun(context.Background(), tc.Function.Arguments)
	if err != nil {
		t.Fatalf("run tool: %v", err)
	}

	// 回填顺序必须是：原始消息 -> assistant 的 ToolCall -> 对应 ToolMessage。
	// Generate 不会修改 arr，因此 a 必须由调用方手动追加。
	arr = append(arr,
		a,
		schema.ToolMessage(tr, tc.ID, schema.WithToolName(tc.Function.Name)),
	)

	final, err := withTools.Generate(context.Background(), arr)
	if err != nil {
		t.Fatalf("generate final answer: %v", err)
	}
	if final == nil || final.Content == "" {
		t.Fatalf("expected final text answer, got %#v", final)
	}
	t.Logf("Final Answer: %s", final.Content)
}
