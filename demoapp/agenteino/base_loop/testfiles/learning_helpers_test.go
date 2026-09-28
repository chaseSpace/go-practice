package testfiles

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	toolutils "github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type addInput struct {
	A int `json:"a" jsonschema_description:"第一个整数"`
	B int `json:"b" jsonschema_description:"第二个整数"`
}

func newAddTool() (tool.InvokableTool, error) {
	return toolutils.InferTool("add", "计算两个整数之和。", func(_ context.Context, input addInput) (int, error) {
		println("证明使用了加法 tool")
		return input.A + input.B, nil
	})
}

// reactStudyModel is deterministic: it asks for add(2,3), then answers after
// the tool result is present. It lets Plan 9 test Eino's loop without a model API.
type reactStudyModel struct {
	mu        sync.Mutex
	callCount int
}

func (m *reactStudyModel) Generate(_ context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	options := model.GetCommonOptions(nil, opts...)
	if len(options.Tools) != 1 || options.Tools[0].Name != "add" {
		return nil, fmt.Errorf("expected add tool to be supplied to the model")
	}
	m.mu.Lock()
	m.callCount++
	call := m.callCount
	m.mu.Unlock()
	if call == 1 {
		return schema.AssistantMessage("", []schema.ToolCall{{
			ID: "add-2-and-3",
			Function: schema.FunctionCall{
				Name:      "add",
				Arguments: `{"a":2,"b":3}`,
			},
		}}), nil
	}
	for _, message := range input {
		if message.Role == schema.Tool && message.Content == "5" {
			return schema.AssistantMessage("2 + 3 = 5", nil), nil
		}
	}
	return nil, fmt.Errorf("tool result was not returned to the model")
}

func (m *reactStudyModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	message, err := m.Generate(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return schema.StreamReaderFromArray([]*schema.Message{message}), nil
}

// nameRecallModel answers correctly only when the caller includes a previous
// self-introduction in the messages it passes to Agent.Run.
type nameRecallModel struct{}

func (nameRecallModel) Generate(_ context.Context, input []*schema.Message, _ ...model.Option) (*schema.Message, error) {
	last := input[len(input)-1]
	if last.Content == "我叫什么名字？" {
		for _, message := range input {
			if message.Role == schema.User && message.Content == "我叫小明。" {
				return schema.AssistantMessage("你叫小明。", nil), nil
			}
		}
		return schema.AssistantMessage("我不知道。", nil), nil
	}
	return schema.AssistantMessage("收到，我会记住这句话出现在本次上下文中。", nil), nil
}

func (m nameRecallModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	message, err := m.Generate(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return schema.StreamReaderFromArray([]*schema.Message{message}), nil
}
