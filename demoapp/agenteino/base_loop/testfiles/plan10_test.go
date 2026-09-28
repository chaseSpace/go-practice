package testfiles

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

func TestPlan10ConversationHistory(t *testing.T) {
	ctx := context.Background()
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "conversation_helper",
		Description: "演示调用方传递历史消息的 Agent",
		Instruction: "根据本次输入消息回答。",
		Model:       nameRecallModel{},
	})
	if err != nil {
		t.Fatalf("create agent: %v", err)
	}

	history := []*schema.Message{schema.UserMessage("我叫小明。")}
	firstReply := runOnce(t, ctx, agent, history)
	history = append(history, firstReply, schema.UserMessage("我叫什么名字？"))
	secondReply := runOnce(t, ctx, agent, history)
	if secondReply.Content != "你叫小明。" {
		t.Fatalf("answer without expected history use: %q", secondReply.Content)
	}
}

func runOnce(t *testing.T, ctx context.Context, agent adk.Agent, messages []*schema.Message) *schema.Message {
	t.Helper()
	iterator := agent.Run(ctx, &adk.AgentInput{Messages: messages})
	for {
		event, ok := iterator.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			t.Fatalf("agent event error: %v", event.Err)
		}
		if event.Output != nil && event.Output.MessageOutput != nil && event.Output.MessageOutput.Message != nil {
			return event.Output.MessageOutput.Message
		}
	}
	t.Fatal("agent returned no assistant message")
	return nil
}
