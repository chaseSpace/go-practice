package testfiles

import (
	"context"
	"strings"
	"testing"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

// functionAgent 是用于学习编排的确定性 Agent。
// 生产代码可将它替换为 adk.ChatModelAgent，而测试不应依赖真实模型服务。
type functionAgent struct {
	name        string
	description string
	run         func(context.Context, *adk.AgentInput) *adk.AgentEvent
}

func (a *functionAgent) Name(_ context.Context) string {
	return a.name
}

func (a *functionAgent) Description(_ context.Context) string {
	return a.description
}

func (a *functionAgent) Run(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iterator, generator := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
	go func() {
		defer generator.Close()
		event := a.run(ctx, input)
		if event.AgentName == "" {
			event.AgentName = a.name
		}
		generator.Send(event)
	}()
	return iterator
}

func assistantEvent(content string, action *adk.AgentAction) *adk.AgentEvent {
	event := adk.EventFromMessage(schema.AssistantMessage(content, nil), nil, schema.Assistant, "")
	event.Action = action
	return event
}

func outputAgent(name, content string) adk.Agent {
	return &functionAgent{
		name:        name,
		description: name + " 子 Agent",
		run: func(_ context.Context, _ *adk.AgentInput) *adk.AgentEvent {
			return assistantEvent(content, nil)
		},
	}
}

func runWithRunner(t *testing.T, ctx context.Context, agent adk.Agent, query string) []*adk.AgentEvent {
	t.Helper()
	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent})
	iterator := runner.Query(ctx, query)
	return collectEvents(t, iterator)
}

func collectEvents(t *testing.T, iterator *adk.AsyncIterator[*adk.AgentEvent]) []*adk.AgentEvent {
	t.Helper()
	var events []*adk.AgentEvent
	for {
		event, ok := iterator.Next()
		if !ok {
			return events
		}
		if event.Err != nil {
			t.Fatalf("agent event error: %v", event.Err)
		}
		events = append(events, event)
	}
}

func contents(events []*adk.AgentEvent) []string {
	result := make([]string, 0, len(events))
	for _, event := range events {
		if event.Output == nil || event.Output.MessageOutput == nil || event.Output.MessageOutput.Message == nil {
			continue
		}
		result = append(result, event.Output.MessageOutput.Message.Content)
	}
	return result
}

func lastUserText(input *adk.AgentInput) string {
	if input == nil {
		return ""
	}
	for index := len(input.Messages) - 1; index >= 0; index-- {
		message := input.Messages[index]
		if message != nil && message.Role == schema.User {
			return message.Content
		}
	}
	return ""
}

func hasContent(values []string, want string) bool {
	for _, value := range values {
		if strings.Contains(value, want) {
			return true
		}
	}
	return false
}
