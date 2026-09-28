# Plan 1：提示链——固定顺序的 Agent 流水线

## 模式定义

提示链把一个大任务拆成顺序固定的小步骤，每个步骤只承担一种职责。后一步依赖前一步的产物，整体像工厂流水线：

```text
用户目标 -> 收集素材 -> 整理提纲 -> 生成初稿 -> 最终结果
```

图片中的“写分析报告”就是典型例子。与自主 Agent 不同，提示链的路径由开发者提前决定，模型没有权力跳过阶段或临时改变步骤。

## 为什么要拆链

把所有要求塞进一个 Prompt，模型需要同时搜索、筛选、规划和写作，任何环节出错都难以定位。拆链后，每个阶段都有明确输入输出，可以独立测试、重试或替换。它的主要收益是可控性，而不是让模型获得更多自主性。

适合提示链的条件：

- 任务步骤稳定，而且顺序不能改变。
- 每个阶段的输出可以被下一阶段消费。
- 希望知道失败发生在素材、提纲还是写作阶段。

不适合的情况是：下一步必须由运行时信息动态决定。此时更适合路由或真正的 Agent 循环。

## Eino ADK 映射

测试使用：

```go
adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
    SubAgents: []adk.Agent{collector, outlineEditor, draftWriter},
})
```

三个子 Agent 都是 `adk.ChatModelAgent`，但职责与指令严格分开：

1. `official_material_collector`：必须调用一次 `collect_eino_official_material`，并通过 `ReturnDirectly` 将真实 Tool 结果直接交给下一阶段，不再自行总结或重复调用。
2. `outline_editor`：读取已收集素材，为 Go Agent 初学者设计由浅入深的大纲，不提前写正文。
3. `draft_writer`：严格依据素材和大纲写中文初稿，不输出推理过程或额外说明。

资料 Tool 使用真实 HTTP 请求抓取两个 CloudWeGo Eino 官方页面，将 HTML 清理为长度受限的正文摘录。URL 在 Go 代码中使用白名单，不接受模型生成的任意地址，避免把通用网页抓取变成 SSRF 入口。

Eino v0.9.19 源码对 Workflow Agent 标记了生产使用提示：它基于共享完整上下文的 Agent transfer；多数复杂多 Agent 场景优先考虑
`ChatModelAgent + AgentTool` 或 DeepAgent。本例使用 `SequentialAgent` 是为了清楚学习顺序架构。若每一步只是普通组件而不是独立
Agent，Chain/Graph 往往更轻量。

## 测试阅读重点

[对应测试](../testfiles/plan1_test.go)是一个自包含的真实集成示例。环境加载、HTTP Tool、三个 Agent、顺序编排与事件读取都在同一个文件中，不依赖公共测试 helper。

模型输出具有随机性，因此这个示例不判断文章内容，也不要求固定措辞。它只处理构造和运行错误，并遍历 ADK 事件；顺序链中最后一条非空 assistant 消息被视为初稿。

运行时只消费 ADK 的 `AgentEvent`：它天然携带所属 Agent 名称，因此能直接记录当前 Agent、消息角色、工具请求/结果、工具调用数量和内容长度。轨迹不打印完整 Prompt、中间素材或密钥，结束后再单独打印最终初稿。需要排查框架内部节点或模型级耗时时，再单独接入 callbacks，不混入这个入门示例。

运行：

```bash
export EINO_RUN_LIVE_TESTS=1
go test ./five_arch/testfiles -run '^TestPlan1RealPromptChain$' -v
```

模型配置复用 `base_loop/testfiles/.env` 中的 `LLM_ENDPOINT`、`LLM_API_KEY` 和 `LLM_MODEL`。该测试会访问外网并产生多次 LLM 请求，默认跳过，避免普通 `go test ./...` 意外消耗额度。

## 设计检查表

- 每一步是否只有一个清楚职责？
- 上一步输出是否经过格式和业务校验？
- 某一步失败时，是终止、重试还是允许降级？
- 是否真的需要 Agent；普通 Go 函数或 Eino Component 是否更简单？
