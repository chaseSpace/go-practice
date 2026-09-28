// Plan 3：并行分支。
package testfiles

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/cloudwego/eino/adk"
)

// parallelMergeAgent 等待所有并行分支结束，再生成一个合并结果。
type parallelMergeAgent struct {
	name     string
	branches adk.Agent
}

func (a *parallelMergeAgent) Name(_ context.Context) string { return a.name }

func (a *parallelMergeAgent) Description(_ context.Context) string {
	return "并行执行后合并所有结果"
}

func (a *parallelMergeAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iterator, generator := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
	go func() {
		defer generator.Close()
		child := a.branches.Run(ctx, input, options...)
		var results []string
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
		// NOTE 这里不只能返回 AssistantMessage。只是因为它更合适，因为 “汇总结论” 本质上是 Agent 要说给用户听的一段话。
		// -- 扩展：Agent event还可以是其他语义：Tool Msg/运行失败：Err/控制流动作：Action: Exit / TransferToAgent / BreakLoop/非消息结构化数据：Output.CustomizedOutput
		summary := assistantEvent("汇总结论："+strings.Join(results, "；"), nil)
		summary.AgentName = a.name
		generator.Send(summary)
	}()
	return iterator
}

// blockingAgent 是测试专用 Agent：它在启动后暂停，直到测试统一放行。
// started/release 只是用于稳定验证并发启动，不是实际并行 Agent 的必需业务逻辑。
type blockingAgent struct {
	name string
	// started 用于向测试报告“当前分支已经进入执行”。
	started chan<- string
	// release 是测试创建的同步屏障；通道关闭后，所有分支同时继续完成。
	release <-chan struct{}
}

func (a *blockingAgent) Name(_ context.Context) string { return a.name }

func (a *blockingAgent) Description(_ context.Context) string { return "并行验证子 Agent" }

func (a *blockingAgent) Run(ctx context.Context, _ *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iterator, generator := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
	go func() {
		defer generator.Close()
		select {
		// 首先报告已启动；若两个分支都能走到这里，说明 ParallelAgent 已同时调度它们。
		case a.started <- a.name:
		// 测试超时时退出，避免子 goroutine 长时间等待。
		case <-ctx.Done():
			generator.Send(&adk.AgentEvent{AgentName: a.name, Err: ctx.Err()})
			return
		}
		select {
		// 真实业务不会等待该通道；这里刻意阻塞，防止某个分支过快结束而掩盖并发行为。
		case <-a.release:
			generator.Send(assistantEvent(a.name+"：分析完成", nil))
		case <-ctx.Done():
			generator.Send(&adk.AgentEvent{AgentName: a.name, Err: ctx.Err()})
		}
	}()
	return iterator
}

func TestParallelAgentStartsAllBranchesBeforeCompletion(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 带缓冲的 started 可容纳两个分支的启动信号，测试不需要先读取才能让分支继续运行。
	started := make(chan string, 2)
	// 关闭 release 会同时唤醒所有等待它的分支。
	release := make(chan struct{})
	parallel, err := adk.NewParallelAgent(ctx, &adk.ParallelAgentConfig{
		Name:        "risk_analysis",
		Description: "并行分析不同风险维度",
		SubAgents: []adk.Agent{
			&blockingAgent{name: "market", started: started, release: release},
			&blockingAgent{name: "regulation", started: started, release: release},
		},
	})
	if err != nil {
		t.Fatalf("create parallel agent: %v", err)
	}

	workflow := &parallelMergeAgent{name: "risk_merger", branches: parallel}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: workflow})
	iterator := runner.Query(ctx, "分析这项投资的市场与监管风险")
	seen := map[string]bool{}
	// 在放行任何分支前，先等待两个分支都报告已启动；
	// 这是本测试证明“并行启动”而非“碰巧先后完成”的关键。
	for len(seen) < 2 {
		select {
		case name := <-started:
			seen[name] = true
		case <-ctx.Done():
			t.Fatalf("parallel branches did not start together: %v", seen)
		}
	}
	// 两个分支均已启动后再统一放行，让它们产生真实输出并进入汇总阶段。
	close(release)

	got := contents(collectEvents(t, iterator))
	if !hasContent(got, "market：分析完成") || !hasContent(got, "regulation：分析完成") {
		t.Fatalf("parallel outputs = %#v", got)
	}
	if !hasContent(got, "汇总结论") {
		t.Fatalf("parallel results were not merged: %#v", got)
	}
}

var _ adk.Agent = (*blockingAgent)(nil)
var _ adk.Agent = (*parallelMergeAgent)(nil)
