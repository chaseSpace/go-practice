# Eino Agent 入门练习

这是一个面向 Go 开发者的、低密度 Eino Agent 学习项目。基础路线刻意只覆盖单 Agent 的第一段路：从消息、Prompt 和结构化数据开始，到工具调用、一个基础 ReAct Agent，再到显式传入多轮上下文。多 Agent、工作流编排、评估与生产可观测性不在基础路线范围内；RAG 则作为独立专题放在 [learn_rag](learn_rag/README.md)，避免把检索知识库与对话记忆混成一个概念。

每一个 `plan/planN_*.md` 只讲一个概念，并由 `testfiles/planN_test.go` 配一段可执行代码。先读一篇、运行一个测试、改一个断言或输入，再进入下一篇；不要一次性通读十篇。

## Eino 是什么

[CloudWeGo Eino](https://www.cloudwego.io/docs/eino/overview/) 是 Go 的 LLM/AI 应用开发框架。它把模型、Prompt、工具、检索等能力抽象成有边界的 Component；核心仓库提供接口与编排，`eino-ext` 提供 OpenAI、Ark、Ollama 等具体集成。这样，业务代码大多依赖稳定的 Go 接口，而不是某个模型厂商的 SDK。

对本练习最重要的四层抽象是：

1. `schema.Message`：模型通信的消息，携带 `system`、`user`、`assistant`、`tool` 等角色。
2. Component：例如 `ChatModel`、`ChatTemplate` 和 `Tool`。它们各自有明确的输入、输出和选项。
3. 编排：Eino 提供 Chain、Graph 与 Workflow；本入门项目暂不实现它们。
4. ADK：Agent Development Kit。`ChatModelAgent` 用模型作决策、用 Tool 作行动，并通过 `Runner`/事件流运行。配置了工具时，它执行 ReAct 循环：模型决定调用什么工具，框架执行工具并把结果送回模型，直到模型给出最终答复。

上述定位和 ReAct 行为来自 Eino 的[概览](https://www.cloudwego.io/docs/eino/overview/)与 [ChatModelAgent 文档](https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/)。Eino 的官方[快速开始索引](https://github.com/cloudwego/eino/blob/main/llms.txt)也说明了核心库、扩展库和示例库的分工。

## Eino 的特点与优势

### 1. 面向 Go 的类型边界

Eino 用 Go 接口和泛型描述组件。例如模型的输入是消息列表、工具有清晰的参数 schema，`utils.InferTool` 可从 Go 结构体推导工具参数。连接不兼容的组件时，许多错误能在编译期或构造期暴露；单元测试也能直接替换模型接口，避免真实网络调用。

这不是“模型输出天然可靠”的保证：自然语言、工具选择和 JSON 内容仍然是不确定的。类型的价值是把**你的程序边界**收紧，而不是替模型做推理。

### 2. 组件可替换、组合可演进

业务逻辑依赖 `ChatModel`、`Tool` 等抽象，提供商适配放在扩展层。因此可以把模型构造代码集中在一个位置，而让 Prompt、工具和 Agent 逻辑保持不变。小需求可直接调用模型；步骤固定时可选 Chain/Graph；需要模型自己选择工具时再使用 Agent。

### 3. Agent 运行时是显式的

ADK 的 `ChatModelAgent` 是完整 Agent，`Runner` 是运行入口，输出是事件迭代器。这让流式输出、工具事件和后续中断/恢复有统一承载方式。初学者要先区分：`ChatModel` 只是“调用模型的组件”，`ChatModelAgent` 才是“协调模型和工具的运行单元”。

### 4. 适合已有 Go 服务的工程约束

若你的 Web 服务、并发任务和部署链路本来就使用 Go，Eino 可以减少跨语言服务与运行时的额外边界。静态二进制、`context.Context` 取消传播、Go 的并发模型和现有测试工具通常能直接复用。是否更有优势仍取决于团队语言、模型提供商和周边集成，而不是框架名字本身。

## 与 Python LangChain 的不同

两者都能统一模型接口、定义工具、做 Agent；它们不是简单的“谁替代谁”。LangChain 当前也把 Agent 建立在图运行时上：其 `create_agent` 会构建 LangGraph 运行时，工具在模型循环中被调用。参见 LangChain 官方 [Agents 文档](https://docs.langchain.com/oss/python/langchain/agents)。

| 维度 | Eino | Python LangChain |
| --- | --- | --- |
| 主要语言与风格 | Go；接口、结构体、泛型和显式错误处理 | Python；函数、类、装饰器与 Pydantic 等动态生态工具 |
| 基本可组合单元 | Component（如 `ChatModel`、`ChatTemplate`、`Tool`） | Model、Tool、Middleware、Runnable/Agent 等 |
| Agent 入口 | ADK 的 `ChatModelAgent` + `Runner`，输出 Agent 事件迭代器 | `create_agent`，其运行时基于 LangGraph |
| 工具定义体验 | Go 结构体可推导参数 schema，工具函数显式返回 `(value, error)` | Python 函数可用 `@tool` 装饰器转换为工具 |
| 提供商集成组织 | 核心 `eino` 与集成 `eino-ext` 分离 | 核心与独立的 `langchain-<provider>` 包；官方文档列出大量集成 |
| 更自然的团队场景 | 现有 Go 后端、重视静态边界和单二进制交付 | Python 数据/AI 工具链成熟、希望快速使用大量 Python 集成 |

LangChain 的优势不应被低估：官方资料指出其 Python 生态有 1000+ 集成，并把监控、评估等能力连接到 LangSmith；见 [集成概览](https://docs.langchain.com/oss/python/integrations/providers/overview) 和[项目说明](https://github.com/langchain-ai/langchain)。相对地，Eino 的优势也不等于“比 Python 更聪明”——它主要是对 Go 工程的契合和明确的组件/运行时边界。

选型时可用一个简单问题判断：团队是否已经用 Go 构建核心在线服务，并希望 Agent 成为同一服务中的一部分？若是，先学 Eino；若团队主要在 Python 中做数据处理、实验和大量第三方 AI 集成，先学 LangChain/LangGraph 往往更顺手。实际项目也可以用 HTTP、消息队列或 MCP 把不同语言的服务连接起来。

## 学习路线

| Plan | 主题 | 对应测试 |
| --- | --- | --- |
| 1 | 项目与依赖准备 | `plan1_test.go` |
| 2 | 首次调用模型 | `plan2_test.go` |
| 3 | 消息角色与对话记录 | `plan3_test.go` |
| 4 | 简单 Prompt | `plan4_test.go` |
| 5 | 变量化 Prompt 模板 | `plan5_test.go` |
| 6 | 结构化输出的解析 | `plan6_test.go` |
| 7 | 第一个本地工具 | `plan7_test.go` |
| 8 | 向模型声明可调用工具 | `plan8_test.go` |
| 9 | 基础 ReAct Agent | `plan9_test.go` |
| 10 | 显式传递多轮上下文 | `plan10_test.go` |

## RAG 专题路线

[learn_rag](learn_rag/README.md) 是基础路线之后的独立系列。它先用可离线验证的最小检索器建立 RAG 的数据、分块、向量、召回和受约束生成概念，再分别实现 SQLite-vec 与 Milvus 两种后端；学习重点是“针对本次问题检索外部语料并据此回答”，不是保存用户偏好或聊天历史的记忆系统。

## 环境与运行

当前项目使用 Go `1.25.6`，并锁定 `github.com/cloudwego/eino v0.9.19` 和 OpenAI 模型适配器。安装依赖后，先运行完全离线的练习：

```bash
go test ./...
go test ./testfiles -run TestPlan7
```

`testfiles/.env` 是本地模型配置文件，测试会自动加载它；已有的 shell/CI 环境变量优先，不会被文件覆盖。文件使用 OpenAI 兼容接口的通用名称：`LLM_ENDPOINT`、`LLM_API_KEY`、`LLM_MODEL`。它和根目录 `.env` 都被 `.gitignore` 忽略；可参照 [testfiles/.env.example](base_loop/testfiles/.env.example)。

`plan2_test.go` 还提供一个真实调用模型的可选冒烟测试。它默认跳过，避免测试消耗额度或因网络失败。只有在 `.env` 或 shell 中明确设置下列开关时才会运行：

```bash
export EINO_RUN_LIVE_TESTS=1
# testfiles/.env 中已配置 LLM_ENDPOINT、LLM_API_KEY、LLM_MODEL
go test ./testfiles -run TestPlan2FirstLLMCall -v
```

不要把 API Key 写进代码、测试断言或 Git。测试中的 `fake` 模型仅用于验证 Eino 的消息、工具和 Agent 编排；它不评估大语言模型的真实能力。

## 官方资料

- [Eino 概览](https://www.cloudwego.io/docs/eino/overview/)
- [Eino 快速开始与组件索引](https://github.com/cloudwego/eino/blob/main/llms.txt)
- [Eino ChatModelAgent / ReAct](https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/)
- [Eino Cookbook](https://www.cloudwego.io/docs/eino/cookbook/)
- [LangChain Agents](https://docs.langchain.com/oss/python/langchain/agents)
- [LangChain Python 集成概览](https://docs.langchain.com/oss/python/integrations/providers/overview)
