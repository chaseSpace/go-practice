package testfiles

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/pkg/errors"
)

func TestPlan9BasicReActAgent(t *testing.T) {
	ctx := context.Background()
	addTool, err := newAddTool()
	if err != nil {
		t.Fatalf("create add tool: %v", err)
	}
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "math_helper",
		Description: "一个能使用加法工具的学习 Agent",
		Instruction: "遇到加法问题时使用 add 工具。",
		Model:       &reactStudyModel{},
		ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{addTool},
		}},
	})
	if err != nil {
		t.Fatalf("create agent: %v", err)
	}

	iterator := agent.Run(ctx, &adk.AgentInput{Messages: []*schema.Message{
		schema.UserMessage("2 加 3 等于多少？"),
	}})
	var finalAnswer string
	for {
		event, ok := iterator.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			t.Fatalf("agent event error: %v", event.Err)
		}
		if event.Output != nil && event.Output.MessageOutput != nil && event.Output.MessageOutput.Message != nil {
			content := event.Output.MessageOutput.Message.Content
			if content != "" {
				finalAnswer = content
			}
		}
	}
	if finalAnswer != "2 + 3 = 5" {
		t.Fatalf("final answer = %q, want %q", finalAnswer, "2 + 3 = 5")
	}
}

func TestPlan9BasicReActAgentManual(t *testing.T) {
	manualImplAgentLoop(t)
}

type ManualAgent struct {
	Model         model.ToolCallingChatModel
	Tools         map[string]tool.InvokableTool
	Messages      []*schema.Message
	MaxModelCalls int
	Trace         []TraceEvent
}

type TraceEvent struct {
	At         time.Time
	Round      int
	Stage      string // model_request / model_response / tool_start / tool_end / final
	ToolName   string
	ToolCallID string
	Duration   time.Duration
	Detail     string
}

func NewManualAgent(t *testing.T) *ManualAgent {
	cm, addTool := newChatModel(t)
	return &ManualAgent{
		Model:         cm,
		Tools:         map[string]tool.InvokableTool{"add": addTool},
		MaxModelCalls: 8,
	}
}

func (a *ManualAgent) record(event TraceEvent) {
	event.At = time.Now()
	a.Trace = append(a.Trace, event)
}

func (a *ManualAgent) printEvents(t *testing.T) {
	for _, event := range a.Trace {
		t.Logf(
			"[%s] round=%d stage=%s tool=%s cost=%s %s",
			event.At.Format(time.RFC3339),
			event.Round,
			event.Stage,
			event.ToolName,
			event.Duration,
			event.Detail,
		)
	}
}

// Run 开始并实现Loop
func (a *ManualAgent) Run(ctx context.Context, inputMsg []*schema.Message, newSession ...bool) (string, error) {
	if len(newSession) > 0 && newSession[0] {
		a.Messages = nil
		a.Trace = nil
	}

	var finalAnswer string
	var chatStart = time.Now()

	a.AppendMessage(inputMsg...)

	for round := 0; round < a.MaxModelCalls; round++ {
		start := time.Now()
		a.record(TraceEvent{
			Round:  round,
			Stage:  "model_req",
			Detail: fmt.Sprintf("msg-len=%d", len(a.Messages)),
		})
		mg, err := a.Model.Generate(ctx, a.Messages)
		a.record(TraceEvent{
			Round:    round,
			Stage:    "model_rsp",
			Duration: time.Since(start),
		})
		if err != nil {
			return "", errors.Wrapf(err, "generate message (round %d)", round)
		}

		if len(mg.ToolCalls) == 0 {
			finalAnswer = mg.Content
			a.record(TraceEvent{
				Round:    round,
				Stage:    "final",
				Duration: time.Since(chatStart),
				Detail:   fmt.Sprintf("msg-len=%d", len(a.Messages)),
			})
			break
		}
		a.AppendMessage(mg)

		for _, tc := range mg.ToolCalls {
			switch tc.Type {
			case "function":
				a.record(TraceEvent{
					Round:      round,
					Stage:      "tool_call",
					ToolName:   tc.Function.Name,
					ToolCallID: tc.ID,
					Detail:     fmt.Sprintf("msg-len=%d tool_calls=%d", len(a.Messages), len(mg.ToolCalls)),
				})
				if tc.Function.Name == "add" {
					tl, ok := a.Tools[tc.Function.Name]
					if !ok {
						a.AppendMessage(schema.ToolMessage("tool not found", tc.ID, schema.WithToolName(tc.Function.Name)))
						break
					}

					start = time.Now()

					toolCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
					ret, err := tl.InvokableRun(toolCtx, tc.Function.Arguments)
					cancel()
					a.record(TraceEvent{
						Round:      round,
						Stage:      "tool_call_end",
						ToolName:   tc.Function.Name,
						ToolCallID: tc.ID,
						Duration:   time.Since(start),
						Detail:     fmt.Sprintf("tool_calls=%d result_len=%d", len(mg.ToolCalls), len(ret)),
					})
					if err != nil {
						a.AppendMessage(schema.ToolMessage("Tool call failed: "+err.Error(), tc.ID, schema.WithToolName(tc.Function.Name)))
						break
					}
					// 工具输出 追加 到消息数组
					a.AppendMessage(schema.ToolMessage(ret, tc.ID, schema.WithToolName(tc.Function.Name)))
				}
			default:
				a.AppendMessage(schema.ToolMessage("unexpected tool call type: "+tc.Type, tc.ID, schema.WithToolName(tc.Function.Name)))
			}
		}

	}
	return finalAnswer, nil
}

func (a *ManualAgent) AppendMessage(msg ...*schema.Message) []*schema.Message {
	a.Messages = append(a.Messages, msg...)
	return a.Messages
}

func manualImplAgentLoop(t *testing.T) {
	// input - think - tool - tool-result - think - output

	var inputArr = []*schema.Message{
		schema.SystemMessage("你是一个简单学习Agent，主要解决加法问题，tool最多调用1次"),
		schema.UserMessage("2 加 3 等于多少？"),
	}

	agent := NewManualAgent(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	final, err := agent.Run(ctx, inputArr)
	if err != nil {
		t.Fatalf("run agent: %v", err)
	}

	agent.printEvents(t)

	if !strings.Contains(final, "5") {
		t.Fatalf("final answer = %q, want %q", final, "2 + 3 = 5")
	}
	t.Logf("final answer = %s", final)
}

func newChatModel(t *testing.T) (model.ToolCallingChatModel, tool.InvokableTool) {
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
	return withTools, addTool
}
