// Plan 5：评估—优化循环。
package testfiles

import (
	"context"
	"fmt"
	"testing"

	"github.com/cloudwego/eino/adk"
)

type draftState struct {
	version  int
	feedback string
}

func TestEvaluatorOptimizerWithLoopAgent(t *testing.T) {
	ctx := context.Background()
	state := &draftState{}
	writer := &functionAgent{
		name:        "writer",
		description: "根据反馈改写草稿",
		run: func(_ context.Context, _ *adk.AgentInput) *adk.AgentEvent {
			state.version++
			if state.feedback == "" {
				return assistantEvent(fmt.Sprintf("草稿 v%d", state.version), nil)
			}
			return assistantEvent(fmt.Sprintf("草稿 v%d，已采纳：%s", state.version, state.feedback), nil)
		},
	}
	evaluator := &functionAgent{
		name:        "evaluator",
		description: "评估草稿是否达标",
		run: func(_ context.Context, _ *adk.AgentInput) *adk.AgentEvent {
			if state.version < 2 {
				state.feedback = "补充事实依据"
				return assistantEvent("评估：需要"+state.feedback, nil)
			}
			return assistantEvent("评估：通过", adk.NewBreakLoopAction("evaluator"))
		},
	}
	loop, err := adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{
		Name:          "writer_evaluator_loop",
		Description:   "写作与评估交替执行",
		SubAgents:     []adk.Agent{writer, evaluator},
		MaxIterations: 3, // 限制整个循环的次数，即每个sub agent可以运行的次数
	})
	if err != nil {
		t.Fatalf("create loop agent: %v", err)
	}

	events := runWithRunner(t, ctx, loop, "写一条产品广告文案")
	got := contents(events)
	for _, expected := range []string{"草稿 v1", "需要补充事实依据", "草稿 v2，已采纳：补充事实依据", "评估：通过"} {
		if !hasContent(got, expected) {
			t.Fatalf("missing %q in %#v", expected, got)
		}
	}
	if state.version != 2 {
		t.Fatalf("writer ran %d times, want 2", state.version)
	}
}
