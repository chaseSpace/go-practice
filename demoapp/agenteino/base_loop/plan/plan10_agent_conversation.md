# Plan 10：实现带上下文的简单 Agent 对话

## 本节只解决什么

理解并实现最小的多轮上下文：应用自己保存先前的 user/assistant 消息，并在下一次 `agent.Run` 时完整传入。

## 理论基础

“模型记得我叫什么”通常不是模型服务自动拥有的长期记忆。对一个无状态请求而言，模型只看到本次传入的消息。最小的会话管理就是维护一个切片：

```text
history = [user(自我介绍), assistant(答复)]
history = append(history, user(追问))
agent.Run(ctx, history)
```

拿到新的 assistant 回复后，再把它追加进历史。这个规则非常朴素，却是后续 Session、数据库存储、摘要压缩、权限隔离的基础。Eino 的 Agent/Runner 负责一次运行；本入门例子明确由调用方负责跨运行保留历史，这能避免误以为每次 `Runner.Query` 都自动“记住”前文。

Plan 10 的假模型只有在输入里发现“我叫小明”时，才会正确回答“你叫什么名字”。测试先运行第一轮、手动拼接历史、再运行第二轮，以可验证的方式展示上下文确实进入了模型。

## 动手步骤

```bash
go test ./testfiles -run TestPlan10ConversationHistory -v
```

尝试删掉历史中的第一条 user 消息，测试应失败。然后把历史替换为另一个名字，观察第二轮答案如何变化。

## 到这里为止，你已经掌握了什么

你已经能用 Go/Eino 建立模型输入、编写 Prompt、解析结果、定义工具、运行一个单 Agent ReAct 循环，并显式传递多轮上下文。下一步可按真实需求选择学习 RAG、Graph/Workflow、持久化 Session 或生产可观测性；它们都不应在基础概念尚未熟悉时一并塞进第一个 Agent。
