// Plan 4：经理—工人。
package testfiles

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type assignment struct {
	worker string
	task   string
}

// managerWorkerAgent 展示经理先拆任务、再调用指定 worker、最后汇总的骨架。
type managerWorkerAgent struct {
	name        string
	assignments []assignment
	workers     map[string]adk.Agent
}

func (a *managerWorkerAgent) Name(_ context.Context) string { return a.name }

func (a *managerWorkerAgent) Description(_ context.Context) string {
	return "拆分任务、调用 worker 并汇总的经理 Agent"
}

func (a *managerWorkerAgent) Run(ctx context.Context, _ *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iterator, generator := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
	go func() {
		defer generator.Close()
		plan := make([]string, 0, len(a.assignments))
		for _, item := range a.assignments {
			plan = append(plan, item.worker+"="+item.task)
		}
		planEvent := assistantEvent("经理拆分："+strings.Join(plan, "；"), nil)
		planEvent.AgentName = a.name
		generator.Send(planEvent)

		results := make([]string, 0, len(a.assignments))
		for _, item := range a.assignments {
			worker, ok := a.workers[item.worker]
			if !ok {
				generator.Send(&adk.AgentEvent{AgentName: a.name, Err: fmt.Errorf("worker %q not found", item.worker)})
				return
			}
			workerInput := &adk.AgentInput{Messages: []*schema.Message{schema.UserMessage(item.task)}}
			child := worker.Run(ctx, workerInput, options...)
			for {
				event, ok := child.Next()
				if !ok {
					break
				}
				generator.Send(event)
				if event.Output != nil && event.Output.MessageOutput != nil && event.Output.MessageOutput.Message != nil {
					results = append(results, event.Output.MessageOutput.Message.Content)
				}
			}
		}
		summary := assistantEvent("经理汇总："+strings.Join(results, "；"), nil)
		summary.AgentName = a.name
		generator.Send(summary)
	}()
	return iterator
}

func taskWorker(name, prefix string) adk.Agent {
	return &functionAgent{
		name:        name,
		description: name + " worker",
		run: func(_ context.Context, input *adk.AgentInput) *adk.AgentEvent {
			return assistantEvent(prefix+lastUserText(input), nil)
		},
	}
}

func TestManagerWorkerSplitsDelegatesAndSummarizes(t *testing.T) {
	ctx := context.Background()
	manager := &managerWorkerAgent{
		name: "market_manager",
		assignments: []assignment{
			{worker: "policy_worker", task: "检查监管政策"},
			{worker: "market_worker", task: "分析市场需求"},
		},
		workers: map[string]adk.Agent{
			"policy_worker": taskWorker("policy_worker", "政策结论："),
			"market_worker": taskWorker("market_worker", "市场结论："),
		},
	}

	got := contents(runWithRunner(t, ctx, manager, "研究某项产品的市场机会"))
	for _, expected := range []string{"经理拆分", "政策结论：检查监管政策", "市场结论：分析市场需求", "经理汇总"} {
		if !hasContent(got, expected) {
			t.Fatalf("missing %q in %#v", expected, got)
		}
	}
}

var _ adk.Agent = (*managerWorkerAgent)(nil)
