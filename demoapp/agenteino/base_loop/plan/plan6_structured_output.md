# Plan 6：让程序读取结构化输出

## 本节只解决什么

定义一个 Go 结构体，并用 `schema.NewMessageJSONParser` 把模型消息中的 JSON 文本解析为它。这里重点是“解析和验证输出”，不是保证模型一定产出合法 JSON。

## 理论基础

自然语言适合人阅读，程序通常需要字段明确的数据。例如学习助手可返回：

```json
{"answer":"Eino 是 Go 的 AI 应用框架", "confidence": 0.8}
```

对应的 Go 结构体给程序一个契约：字段名、类型和可选性都清楚。`MessageJSONParser[T]` 从 `schema.Message.Content`（默认来源）读取 JSON 并反序列化为 `T`。解析失败必须被当成正常错误处理：可能是模型没有遵守格式、内容被 Markdown 代码块包住、字段类型不对，或上游把错误消息当成了结果。

先有“请求模型只输出 JSON”的 Prompt 并不够。真正可靠的应用至少要解析、验证必填字段、处理失败；更严格时再使用模型提供商的 JSON Schema / structured output 能力或重试策略。

Plan 6 提供两个层次的测试。离线测试直接提供合法 JSON，专门验证程序端的最后一道解析边界；真实测试则发送一个要求 JSON 的 Prompt，取得真实 `ChatModel.Generate` 回复后再解析。真实模型有时仍会包上 Markdown 代码围栏，所以示例先做非常小的规范化，再让 JSON parser 决定内容是否合法。规范化不是“修复模型输出”：错误的 JSON、缺字段或类型不对仍必须失败。

## 动手步骤

```bash
go test ./testfiles -run TestPlan6StructuredOutputParser -v
```

准备好 `testfiles/.env` 后，把 `EINO_RUN_LIVE_TESTS=1`，可运行真实回复的完整示例：

```bash
go test ./testfiles -run TestPlan6ParseRealChatReply -v
```

该测试使用与 Plan 2 相同的 `LLM_ENDPOINT`、`LLM_API_KEY`、`LLM_MODEL` 配置，并会消耗一次模型请求。

把测试中的 `confidence` 改成字符串，确认解析失败。再思考：如果业务要求值必须在 0 到 1，JSON 解析成功后还应加什么校验？答案是应用层的数值范围验证。

## 本节验收

能区分三件事：要求 JSON 的 Prompt、模型提供的结构化输出能力、Go 端的解析/业务验证。下一节开始给 Agent 准备可执行的工具。
