# Plan 9：组合模型和工具，创建基础 Agent

## 本节只解决什么

使用 Eino ADK 的 `ChatModelAgent` 把一个模型与一个工具组合成基础 ReAct Agent，并理解运行时发生的最小循环。

## 理论基础

当工具配置存在时，`ChatModelAgent` 的工作序列是：

```text
用户消息 -> 模型决定 ToolCall -> 执行已注册 Tool -> ToolMessage 回填 -> 模型最终回答
```

这就是 ReAct 的简化版本：Reason（模型决定）、Action（工具调用请求）、Act（执行工具）、Observation（工具结果成为下一轮上下文）。若没有工具，`ChatModelAgent` 会退化为一次普通模型调用。Agent 的价值不是替你编写业务函数，而是管理这个“模型—工具—模型”的重复流程。

创建 Agent 时至少说明名称、描述、指令和模型；工具放入 `adk.ToolsConfig`。运行时可直接调用 `agent.Run`，也可用 `adk.Runner` 统一入口。输出是事件迭代器，因而一次运行可能看到模型消息、工具结果和最终消息，而不是只有一个字符串。

`plan9_test.go` 使用本地脚本模型：第一轮它确定会请求 `add(2,3)`，收到工具结果后第二轮固定回答 `2 + 3 = 5`。因此测试验证的是 Eino 的 Agent 编排，不是“模型是否足够聪明”。真实模型可能选择不调用工具、生成错误参数或调用不合适的工具，这些都要通过 Prompt、工具描述、限制和测试来改善。

## 动手步骤

```bash
go test ./testfiles -run TestPlan9BasicReActAgent -v
```

随后阅读 `newAddTool` 和 `reactStudyModel`。前者是真实 Eino Tool，后者是替身模型。删除 `ToolsConfig` 后思考：为什么模型调用请求再也没有执行者？

## 本节验收

能够画出 ReAct 循环，并知道 Agent 只会执行你显式注册的工具。下一节在不引入持久化存储的前提下，处理最基础的多轮上下文。
