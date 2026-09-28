# Plan 5：评估—优化——基于反馈循环改进

## 模式定义

评估—优化模式让生成者提交草稿，由评估者按照标准给出反馈；未达标就带着反馈重写，达标后退出：

```text
生成草稿 -> 评估
   ^          |
   |-- 反馈 --| 未通过
              \-> 通过 -> 最终结果
```

图片中的广告文案示例属于这一类。它适合“质量可以被明确标准判断”的任务，例如格式完整性、事实引用、代码测试结果或文案禁用词。

## 评估者不只是说好不好

有效反馈至少包含：

- 是否通过；
- 未通过的具体原因；
- 下一版可执行的修改建议；
- 可选的分项评分。

如果评估者只说“再优化一下”，生成者很难产生稳定改进。能用程序验证的标准优先使用代码，例如 JSON Schema、单元测试和规则检查；模型评估适合处理风格、清晰度等难以完全规则化的部分。

## Eino ADK 映射

测试使用：

```go
adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{
    SubAgents:     []adk.Agent{writer, evaluator},
    MaxIterations: 3,
})
```

`writer` 与 `evaluator` 交替执行。第一次评估写入反馈“补充事实依据”，第二版 writer 明确读取并采纳反馈；第二次评估通过后，通过 `adk.NewBreakLoopAction("evaluator")` 提前结束循环。

示例用共享 `draftState` 表达草稿与反馈，便于离线测试。真实系统应使用每次运行隔离的会话状态，避免并发请求共享同一个可变对象。

Eino v0.9.19 对 `LoopAgent` 的生产多 Agent 使用也有同样提示；本例把它作为理解循环和 `BreakLoopAction` 的最小教材。

## 测试阅读重点

[对应测试](../testfiles/plan5_test.go)验证事件顺序包含：

```text
草稿 v1
评估：需要补充事实依据
草稿 v2，已采纳：补充事实依据
评估：通过
```

并断言 writer 只运行两次，证明 `BreakLoopAction` 生效。

运行：

```bash
go test ./five_arch/testfiles -run TestEvaluatorOptimizerWithLoopAgent -v
```

## 生产边界

- 必须设置最大迭代次数、总超时和 token/费用预算。
- 防止评估者每轮改变标准；评估 rubric 应固定且可版本化。
- 保存每一版草稿与反馈，才能判断质量是否真的改善。
- 达到上限仍未通过时，应返回“未达标”的明确状态，而不是把最后一版伪装为成功。
