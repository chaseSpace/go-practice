# Plan 2：第一次调用聊天模型

## 本节只解决什么

构造一个 Eino 的 OpenAI ChatModel，并理解一次非流式 `Generate` 调用的输入和输出。不要在本节引入工具、Agent 或记忆。

## 理论基础

`ChatModel` 是“把消息发给模型并得到消息”的组件。调用形状可以概括为：

```text
[]*schema.Message  ->  ChatModel.Generate  ->  *schema.Message
```

模型适配器由 `eino-ext` 提供。本项目使用 `openai.NewChatModel` 构造适配器，并从 `testfiles/.env` 自动读取 `LLM_API_KEY`、`LLM_MODEL` 和可选的 `LLM_ENDPOINT`。

它们是 OpenAI 兼容接口的通用配置名；也兼容旧的 `OPENAI_*` 环境变量。密钥属于运行时配置，不属于源码；把它提交到仓库既不安全，也会让项目无法在别人的环境中运行。

测试区分两类检查。`TestPlan2CanConstructModel` 用假 Key 创建客户端，不产生网络请求；它验证配置代码的形状。`TestPlan2FirstLLMCall` 是真实冒烟测试，只有 `EINO_RUN_LIVE_TESTS=1` 

且设置 Key 后才运行。真实调用失败时，优先检查模型名、Base URL、网络和账号权限，而不是先改 Prompt。

## 关键代码

```go
chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
    APIKey:  os.Getenv("LLM_API_KEY"),
    Model:   os.Getenv("LLM_MODEL"),
    BaseURL: os.Getenv("LLM_ENDPOINT"),
})
answer, err := chatModel.Generate(ctx, []*schema.Message{
    schema.UserMessage("用一句话解释 Eino"),
})
```

`Generate` 返回一条 assistant 消息。真正的网络程序还应设定 `context.WithTimeout`，以免服务端长期无响应时一直等待。

## 动手步骤

先运行离线构造测试：

```bash
go test ./testfiles -run TestPlan2CanConstructModel -v
```

准备好模型账号后再运行真实测试：

```bash
# 在 testfiles/.env 中把 EINO_RUN_LIVE_TESTS 改为 1；
# LLM_ENDPOINT、LLM_API_KEY、LLM_MODEL 已由该文件提供。
export EINO_RUN_LIVE_TESTS=1
go test ./testfiles -run TestPlan2FirstLLMCall -v
```

## 本节验收

你能指出：消息是输入，`Generate` 的返回值仍是一条消息；客户端构造不等于请求已经发出。下一节学习这条消息的角色语义。
