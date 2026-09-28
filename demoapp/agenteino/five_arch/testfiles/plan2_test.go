// Plan 2：路由（分类后只调用一个专家）。
package testfiles

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/cloudwego/eino/adk"
)

// routerAgent 先确定任务类别，再把原始输入交给唯一的专家 Agent。
type routerAgent struct {
	name    string
	experts map[string]adk.Agent
}

func (a *routerAgent) Name(_ context.Context) string { return a.name }

func (a *routerAgent) Description(_ context.Context) string {
	return "按任务类别选择唯一专家的路由 Agent"
}

func (a *routerAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	// 两者共享同一条事件通道：iterator 交给调用方读取，generator 留给本 Agent 写入。
	// generator 生产者， 可以随时发送你想发的Event
	// iterator 消费者，Run() 立即返回它；Runner 或调用方不断调用 Next() 读取事件
	iterator, generator := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	// Agent.Run 必须立即返回 iterator；真正的路由和专家执行在后台完成。
	go func() {
		// 无论正常结束、路由失败或子 Agent 结束，都必须关闭通道；
		// 否则调用方的 iterator.Next() 会一直等待下一条不存在的事件。
		defer generator.Close()

		route := selectRoute(lastUserText(input))
		routeEvent := assistantEvent("路由决策："+route, nil)
		routeEvent.AgentName = a.name
		// 先把路由决策作为第一条事件发送给调用方，便于观察执行轨迹。
		generator.Send(routeEvent)

		expert, ok := a.experts[route]
		if !ok {
			errorEvent := &adk.AgentEvent{AgentName: a.name, Err: fmt.Errorf("no expert for route %q", route)}
			generator.Send(errorEvent)
			return
		}

		// 子 Agent 也返回自己的 iterator。路由 Agent 不处理其业务输出，
		// 只把每条事件原样转发到自己的 generator，形成一个统一事件流。
		child := expert.Run(ctx, input, options...)
		for {
			// Next 会阻塞到子 Agent 产生事件，或在子 Agent 关闭通道时返回 ok=false。
			event, ok := child.Next()
			if !ok {
				return
			}
			// 调用方最终从 routerAgent 返回的 iterator 中收到这条专家事件。
			generator.Send(event)
		}
	}()

	// 此时后台 goroutine 可能尚未执行；调用方立即获得可消费的异步事件流。
	return iterator
}

func selectRoute(question string) string {
	if strings.Contains(question, "账单") || strings.Contains(question, "退款") {
		return "billing"
	}
	return "technical"
}

func TestRouterSelectsOneExpert(t *testing.T) {
	ctx := context.Background()
	router := &routerAgent{
		name: "support_router",
		experts: map[string]adk.Agent{
			"billing":   outputAgent("billing_expert", "账单专家：已说明退款进度。"),
			"technical": outputAgent("technical_expert", "技术专家：请检查网络设置。"),
		},
	}

	events := runWithRunner(t, ctx, router, "我的账单退款什么时候到账？")
	got := contents(events)
	if !hasContent(got, "路由决策：billing") || !hasContent(got, "账单专家") {
		t.Fatalf("router did not select billing expert: %#v", got)
	}
	if hasContent(got, "技术专家") {
		t.Fatalf("router must not invoke two experts: %#v", got)
	}
}

var _ adk.Agent = (*routerAgent)(nil)
