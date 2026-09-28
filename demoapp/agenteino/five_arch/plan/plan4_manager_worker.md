# Plan 4：经理—工人——动态拆分任务并委派

## 模式定义

经理—工人模式由一个协调者负责理解目标、拆分任务、选择 worker，最后汇总结果：

```text
用户目标 -> 经理拆分
              |-> 政策 worker -> 结果 --\
              |-> 市场 worker -> 结果 ----> 经理汇总
```

它和固定提示链最大的区别是：子任务通常在运行时生成，数量和内容可能随用户目标变化。它和路由也不同：路由通常选择一个专家；经理可以创建多个任务并调用多个
worker。

## 经理与 worker 的职责

经理应该负责：

- 将目标拆成可执行、可验证的子任务；
- 为每个任务选择允许的 worker；
- 设置边界、依赖关系和完成标准；
- 汇总结果，发现冲突时要求补充。

worker 应该只关注一个专业任务，并返回结构化结果。经理不要亲自完成所有专业工作，否则架构只是多了一层 Prompt；worker
也不要擅自扩大任务范围。

## Eino ADK 映射

测试中的 `managerWorkerAgent` 实现 `adk.Agent`，维护两份受控数据：

```go
assignments []assignment
workers     map[string]adk.Agent
```

任务清单代表经理的规划结果，map 是允许调用的 worker 注册表。`Run` 依次发出拆分事件、给每个 worker
独立输入、收集结果，最后生成汇总事件。测试用确定性任务代替 LLM 规划，使架构行为可重复；真实项目可让模型产生经过 schema 校验的
`assignments`。

对于 Eino ADK，还可以把 worker 通过 `adk.NewAgentTool` 暴露为经理 `ChatModelAgent` 的工具。这通常比让 Agent 任意 transfer
更容易限制行动空间。

## 测试阅读重点

[对应测试](../testfiles/plan4_test.go)验证四个阶段都存在：经理拆分、政策 worker、市场 worker、经理汇总。

运行：

```bash
go test ./five_arch/testfiles -run TestManagerWorkerSplitsDelegatesAndSummarizes -v
```

## 生产边界

- 限制子任务数量、递归深度和总模型调用预算。
- worker 必须白名单注册，不能执行模型临时编造的名字。
- 任务需要唯一 ID，便于重试、去重和追踪。
- 汇总前检查结果是否齐全、是否互相矛盾。
- 可并行的 worker 可以并行运行；存在依赖时必须显式排序。
