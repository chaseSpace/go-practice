# Plan 8：让模型知道有哪些工具

## 本节只解决什么

把 Tool 的 `ToolInfo` 绑定到支持工具调用的 ChatModel。此节只建立“模型可以看见的行动清单”，不自己实现 ReAct 循环。

## 理论基础

工具调用至少有两个阶段，不能混为一谈：

1. **声明/选择**：将 `ToolInfo` 交给模型；模型返回某个工具名及 JSON 参数的 `ToolCall`。
2. **执行/回填**：程序执行真正的工具，把结果作为 tool 消息放回下一轮模型输入。

`ChatModel.WithTools` 返回一个带工具定义的模型变体。Eino 文档建议优先使用这种返回新实例的方式，而不是可变地在共享模型上绑定工具，因为并发请求可拥有不同的工具集合。Plan 8 的离线测试仅验证工具信息能够被绑定；没有调用 `Generate`，所以不会联网。

即便模型声称要调用 `add`，也不要直接信任工具名和参数。执行层必须只允许已注册工具、校验 JSON、限制资源。Eino 的 ToolsNode/Agent 会帮助完成路由和回填，但安全策略仍由应用负责。

## 动手步骤

```bash
go test ./testfiles -run TestPlan8BindToolToModel -v
```

阅读测试：先从 Plan 7 工具取得 `Info`，再调用 `chatModel.WithTools([]*schema.ToolInfo{info})`。尝试查看 `info.Name` 与 `info.Desc`，并解释为什么模型需要它们。

## 本节验收

能指出“绑定工具”不会执行工具，也不会保证模型选择它。下一节由 `ChatModelAgent` 把选择、执行和结果回填连成一个循环。
