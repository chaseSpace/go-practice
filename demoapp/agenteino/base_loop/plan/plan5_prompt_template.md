# Plan 5：使用变量化 Prompt 模板

## 本节只解决什么

使用 `prompt.FromMessages` 和 `schema.FString` 把变量填入消息模板；只学习渲染，不调用模型。

## 理论基础

模板解决的是“复用输入结构”，不是让模型更聪明。把固定的 system 指令写一次，把每次变化的 `question`、`language` 等值作为变量传入，能让代码更清楚，也方便测试渲染结果。

Eino 的 `DefaultChatTemplate` 接收一组消息模板。`Format(ctx, variables)` 后得到普通 `[]*schema.Message`，然后才能交给 ChatModel。也就是说，模板位于“构造输入”阶段，模型位于“执行推理”阶段，两者应当分开测试。

本节采用 `schema.FString`，占位符写作 `{question}`。模板变量缺失通常是程序错误，应尽早让格式化报错；不要让空字符串悄悄进入模型。多轮历史也可以通过 `schema.MessagesPlaceholder("history", true)` 放入模板，但本轮先只使用一个问题变量。

## 动手步骤

```bash
go test ./testfiles -run TestPlan5PromptTemplate -v
```

测试断言渲染后第二条消息中出现变量内容。把变量名故意拼错一次，观察格式化错误；然后恢复正确名称。

## 常见误区

- 模板不替代输入校验。用户提供的数据仍要按业务规则长度限制、过滤或转义。
- 不要把用户输入拼进 system 指令来伪装成规则；它仍是不可信输入。
- 渲染成功不代表模型能严格输出 JSON 或调用工具，那是后面章节的独立能力。

## 本节验收

你能画出流程：变量 map -> Prompt 模板 -> 消息列表 -> ChatModel。下一节让模型文本能被 Go 程序解析为结构化数据。
