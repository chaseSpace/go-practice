## Agent 框架横评

| 名称               | 架构范式     | 核心原语                                         | 灵活性                                    | 特性                                                    | 难度    | 缺点                                           |
|------------------|----------|----------------------------------------------|----------------------------------------|-------------------------------------------------------|-------|----------------------------------------------|
| Langchain        | 图状态机     | node、edge、state                              | 开发者自行设计状态模式、编排逻辑、错误处理                  | Checkpointer（时间旅行）、Interrupt（人工介入） 、StreamingTransmit | ⭐⭐    | 过于灵活                                         |
| CrewAI           | 角色驱动范式   |                                              | 300+连接器；130+格式解析器；混合检索                 | 事件驱动、异步优先、步骤化执行；数据密集型Agent独一档                         | ⭐⭐⭐⭐⭐ | 适合市场调研，不适合长时间运行的工作流。控制力不如LangGraph；0.x版本，不稳定 |
| LlamaIndex       | 事件驱动范式   | 行 3-3                                        |                                        |                                                       |       |
| AgentScope（阿里巴巴） | 事件驱动范式   |                                              | 每个步骤可视化展示，每个提示词可编辑，模型调用链路可追溯，每个推理步骤可审计 | 支持模型微调、Python/Java双语言                                 |
| OpenAIAgent      | SDK 封装范式 | Agent（执行单元）、Handoff（控制器移交）、Guardrail（输入输出校验） |                                        | 5行代码启动；Responses API、实时语音；没有checkpointer等             |       |
| PydanticAI       | SDK 封装范式 | Input、Agent、Output                           | 25+模型提供商；TestModel ；CI/CD救命功能？         |                                                       |       |
| 微软Agent框架        | 企业统一类    |                                              |                                        |                                                       |       |
| Google ADK       | 企业统一类    |                                              |                                        |                                                       |       |
| Eino（Go）         |          |                                              |                                        | 字节跳动内部大模型应用的首选全代码开发框架                                 |       |                                              ||
| LangChainGo      |          |                                              |                                        |                                                       |       |

- LangGraph 胜在控制力
- CrewAI 胜在易用性
- LlamaIndex 胜在数据
- Pydantic AI 胜在类型安全
- Dify 胜在低门槛

**👉 Agent 框架几乎全部发生在 Python 生态（原因是 LLM tooling + research stack）。**

### Agent 基础信息和评价维度

- 基础信息：语言/出品方/开源年限/STARS/ISSUE数量/贡献者数量/上次提交时间/上手难度/抽象程度/优势
- 架构与功能：语言/出品方/架构范式/核心概念/状态管理/适用场景
- 生态建设：插件生态/LLM生态/多Agent支持
- 工程能力：语言/出品方/可观测性/DEBUG/并发性/重试机制/生产部署/成本控制

采集时间：2026-03-30

#### 基础信息横评

| 语言                  | 项目                                                            | 出品方           | 开源年限       | STARS  | ISSUE数量 | 贡献者数量 | 最后提交时间      | 上手难度 | 优势                       |
|---------------------|---------------------------------------------------------------|---------------|------------|--------|---------|-------|-------------|------|--------------------------|
| Python              | [LangGraph](https://github.com/langchain-ai/langgraph)        | Langchain     | 2023 年     | 27.9K+ | 245     | 288   | 9 hours ago | ⭐⭐⭐⭐ | 定位通用型复杂应用开发框架            |
| Python & TypeScript | [LlamaIndex](https://github.com/run-llama/llama_index)        | Meta          | 2022 年     | 48K+   | 189     | 1851  | yesterday   | ⭐⭐   | 数据连接器多/索引灵活/模块化设计        |
| Python              | [CrewAI](https://github.com/crewAIInc/crewai)                 | crewAIInc(巴西) | 2023年11月   | 47.6K+ | 98      | 302   | 2 days ago  | ⭐⭐   | 内置多Agent编排，无需构建状态机；无抽象概念 |
| Python              | [OpenAIAgent](https://github.com/openai/openai-agents-python) | OpenAI        | 2025 年 3 月 | 20.4K+ | 58      | 238   | 2 days ago  | ⭐⭐   |                          |                          

**关于LangChain与LangGraph的区别**

LangChain 是一个用于构建大模型应用的基础框架，以线性组件和链式（DAG）工作流为主，适合快速开发简单应用。LangGraph 是构建在
LangChain 之上的扩展库，专注于通过图结构（节点和边）管理循环、复杂逻辑和持久化状态，适合构建高级 Agent 系统。简单来说：LangChain
是组件拼装，LangGraph 是图工作流编排。

#### 架构与功能横评

| 语言                  | 项目                                                     | 架构范式 | 核心概念                                   | 状态管理               | 适用场景                 |
|---------------------|--------------------------------------------------------|------|----------------------------------------|--------------------|----------------------|
| Python              | [LangGraph](https://github.com/langchain-ai/langgraph) | 图状态机 | node、edge、state                        | 原生持久化              | 复杂Agent/多Agent写作/长流程 |
| Python & TypeScript | [LlamaIndex](https://github.com/run-llama/llama_index) | 数据驱动 | Data Connectors/Retriever/Query Engine | 插件式                | 大型知识库系统              |
| Python              | [CrewAI](https://github.com/crewAIInc/crewai)          | 角色驱动 | Agent/Task/Process/Crew                | 提供短期、长期、实体和上下文四层记忆 | 快速构建原型               |

### Agent 运行时可能出现的问题

#### 1. 幻觉

#### 2. 工具调用失败

#### 3. 状态丢失（长任务）

#### 4. 如何重试

### 底层能力决定上线

理解大模型本身，才能驾驭框架。

### 部署成本

