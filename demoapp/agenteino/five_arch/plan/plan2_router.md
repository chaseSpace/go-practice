# Plan 2：路由——先分类，再选择一个专家

## 模式定义

路由模式先判断任务属于哪一类，再只调用对应的专家：

```text
                   -> 账单专家
用户问题 -> 路由器
                   -> 技术专家
```

一次请求通常只走一条分支。客服系统很适合：退款问题交给账单专家，故障问题交给技术专家。路由器负责“找谁处理”，专家负责“如何处理”。

## 与提示链的区别

提示链的下一步在编码时已经确定；路由的下一步取决于本次输入。路由可以减少无关上下文和模型调用，但引入了一个新风险：分类错误会把正确问题交给错误专家。

路由方式可以由简单到复杂逐步选择：

1. 关键字或业务规则：稳定、便宜、容易测试。
2. 小模型分类：适合自然语言边界较模糊的类别。
3. 结构化输出：模型返回受约束的 route 枚举与置信度。

不要让模型返回任意 Agent 名称后直接执行。必须用白名单将 route 映射到已经注册的 `adk.Agent`。

## Eino ADK 映射

测试实现了 `routerAgent`，它直接满足 `adk.Agent` 接口：

```go
type routerAgent struct {
    name    string
    experts map[string]adk.Agent
}
```

`Run` 先调用确定性的 `selectRoute`，记录路由事件，然后只运行 `experts[route]`。生产版本可以把 `selectRoute` 换成 `ChatModelAgent` 或结构化分类模型，但白名单映射仍然保留。

### 为什么这里自定义 `adk.Agent`

Eino 的 `adk.Agent` 是统一执行协议：它提供名称、说明和 `Run()`；`Run()` 不直接返回一个字符串，而是返回 Agent 事件流。框架因此可以用同一种方式消费模型回答、工具结果、路由决定、错误和中断。

路由决策本身是确定性的业务代码，使用完整的 LLM Agent 反而增加成本和不确定性；但路由之后要继续调用已有专家 Agent，并把专家输出交给 Runner。因此本例实现一个很薄的 `routerAgent`，只负责选择和转发，不重复实现专家能力。

```text
routerAgent 后台逻辑                 Runner / 调用方
generator.Send(路由事件)                   iterator.Next()
generator.Send(专家事件)        ->          iterator.Next()
generator.Close()                         Next() 返回 ok=false
```

关键代码是：

```go
iterator, generator := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
```

- `iterator` 是消费端，必须由 `Run()` 立即返回给框架。
- `generator` 是生产端，后台 goroutine 使用 `Send` 推送路由事件和专家事件。
- 两者共享 Eino 的内部异步事件通道；`Close()` 表示不会再产生事件。遗漏 `Close()` 会让调用方一直等待。

在当前 Eino 版本中，`AsyncIterator` 的内部通道没有导出；`NewAsyncIteratorPair` 是手写自定义 Agent 时创建该事件流的标准公共入口。它属于“异步生成器对”或“生产者—消费者”模式。

通常不需要直接使用它：`adk.NewChatModelAgent`、顺序/并行/循环 Agent 等内置实现已经在内部管理事件流。只有当控制流本身由你的 Go 代码决定，例如路由、权限判断、固定业务步骤或对现有 Agent 的事件转发时，才适合自定义 `adk.Agent`。

## 测试阅读重点

[对应测试](../testfiles/plan2_test.go)输入退款问题，并验证：

- 出现 `路由决策：billing`；
- 账单专家被调用；
- 技术专家没有被调用。

第三个断言很重要：只验证“正确专家运行”还不够，也要证明其他分支没有误运行。

运行：

```bash
go test ./five_arch/testfiles -run TestRouterSelectsOneExpert -v
```

## 生产边界

- 给无法分类的输入准备 `fallback`，不要默认落入高风险专家。
- 记录 route、置信度和最终处理结果，定期分析误分类。
- 权限检查必须在专家执行前再次完成，不能只相信路由结果。
- 专家很多时，应先做分层路由，避免把几十个说明一次性塞进模型上下文。
