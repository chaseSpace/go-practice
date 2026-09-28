# 如何在本仓库编写一个学习 Plan

本仓库不是把知识点写成一篇长文，而是把一个专题拆成连续、可执行、可验证的 Plan。学习者每次只读一个 `plan/planN_*.md`，运行同编号的 `testfiles/planN_test.go`，改一个输入或断言，再继续下一节。

Plan 的目标是让读者看见“概念 → 最小代码 → 自动验收”的因果关系；它不是产品需求文档，也不是一次性堆满所有高级特性的示例。

## 现有专题与边界

| 目录 | 学习重点 | 代码策略 |
| --- | --- | --- |
| `base_loop/` | Eino 单 Agent 基础：消息、Prompt、结构化输出、Tool、ReAct、显式上下文 | 直接学习 Eino Component/ADK；真实模型调用必须可控 |
| `five_arch/` | 五种 Agent 编排：顺序链、路由、并行、经理—工人、评估—优化 | 用确定性 fake Agent 验证编排；需要真实模型/HTTP 时单独开启 |
| `learn_rag/` | 文档 RAG：解析、分块、向量契约、检索、引用、SQLite-vec、Milvus | 核心链路用普通 Go 接口和官方存储 SDK；真实 TXT/PDF/DOCX 是主语料，Eino 只在真实模型适配层按需使用 |

不要把“对话记忆”“用户画像”“长期状态”混进 RAG Plan。RAG 的对象是可追溯、可删除、可重建、可评测的外部文档证据。

## 新专题的目录契约

新增专题使用下面的结构；名称采用小写 snake_case，例如 `learn_rag/`。

```text
<topic>/
├── README.md                  # 专题定位、总路线、运行方式、依赖/服务边界
├── plan/
│   ├── plan1_<slug>.md
│   ├── plan2_<slug>.md
│   └── ...
├── testfiles/
│   ├── plan1_test.go
│   ├── plan2_test.go
│   └── <shared>_helpers_test.go
└── data/                      # 仅专题确实需要真实文件语料时添加
    ├── README.md
    ├── raw/                   # 私有原始资料；通常由 .gitignore 忽略
    └── derived/               # manifest、索引、临时产物；通常由 .gitignore 忽略
```

规则：

1. 每个 `plan/planN_*.md` 必须且只对应一个 `testfiles/planN_test.go`。
2. `N` 从 1 连续递增；不要跳号、不要用 `plan01`，也不要把多个知识点塞进同一个编号。
3. 专题 README 的路线表必须列出 Plan、主题、验收目标和精确的测试文件路径。
4. 共享的构造器、fake、解析器、评测器放入 `testfiles/*_helpers_test.go`；不要复制到每个 Plan 测试中。
5. 若代码已不只是教学测试，提取到普通 `.go` package；测试仍以 `planN_test.go` 作为该节的入口。

## Plan 文档的最低结构

每份 Plan 只推进一个新能力，并至少包含下列小节。

```markdown
# Plan N：<动词 + 单一主题>

## 本节只解决什么

说明本节的唯一目标，以及明确不解决的相邻问题。

## 核心概念

用最短的必要理论、数据流或状态图解释为什么需要它。

## 代码映射

指出 `../testfiles/planN_test.go` 中哪些函数、接口或断言体现该概念。

## 动手步骤

给出从仓库根目录可直接复制执行的命令。

## 本节验收

写出通过条件、关键断言和失败时首先检查的位置。

## 容易混淆的点 / 生产边界

说明本节的简化假设，避免学习者把 demo 当生产结论。
```

推荐使用一个小型 ASCII 图解释有三个以上节点的流程，例如：

```text
源文档 -> 解析 -> 分块 -> 检索候选 -> 证据上下文 -> 回答 + 引用
```

图必须与同一 Plan 的代码一致，不能为了“看起来完整”添加未实现组件。

## 测试是课程代码，不是附录

`testfiles/planN_test.go` 要能独立阅读、独立运行。测试名称应表达行为，例如 `TestPlan5InMemoryRetrieverRanksAndFilters`，而不是 `TestSomething`。

每个测试遵守以下顺序：

1. 准备本节输入和依赖；
2. 执行一个清楚的动作；
3. 对本节新引入的行为做断言；
4. 失败信息说明哪个边界被破坏。

不要断言真实 LLM 的固定自然语言句子、浮点相似度的偶然小数，或事件内部实现细节。应断言来源、角色、调用顺序、结构、权限范围、是否引用了输入证据等稳定契约。

### 真实语料专题的额外规则

只要专题以用户提供的文件、数据库记录或其他语料为学习对象，**每一个 Plan 测试的主路径都必须使用该真实数据**。手写字符串仅可用于同一 Plan 内的局部对照，例如：

- 空输入、错误向量维度、无权限租户等错误路径；
- 用于证明删除/替换生效的最小差异；
- 无法由真实资料稳定复现的单一边界条件。

不能以几条手写的“部署端口”“数据库超时”文本代替已有语料。`learn_rag/testfiles` 的统一入口是 `requireCourseDocuments` / `requireCourseChunks`：它会加载 `learn_rag/data/`（或 `RAG_SOURCE_DIR`），真实提取 TXT、PDF、DOCX，再把结果交给每个 Plan。新增数据驱动专题应建立同样的单一入口，保证所有 Plan 使用同一份语料、同一份 metadata 与相同的权限模型。

真实语料测试不得把正文、密钥或整段 Prompt 打印到日志。日志只输出格式数量、文档 ID、chunk ID、页码、耗时、分数或聚合统计。

## 依赖与外部服务分层

默认测试必须离线、确定性、可重复：

```bash
go test ./<topic>/testfiles -v
```

将不稳定或有成本的能力放到明确的 opt-in 层，而不是让普通测试隐式访问网络。

| 能力 | 推荐做法 | 运行开关 / 前置条件 |
| --- | --- | --- |
| LLM、embedding、远程 HTTP | fake 或本地确定性实现覆盖核心契约；另建 live 冒烟测试 | 显式环境变量，如 `EINO_RUN_LIVE_TESTS=1` |
| 真实本地文件语料 | 默认读取已放入专题 `data/` 的材料；可用环境变量替换目录 | `RAG_SOURCE_DIR` 等路径配置，不提交私有资料 |
| SQLite-vec | 将 CGO/系统依赖测试置于明确 build tag | `CGO_ENABLED=1 go test -tags=sqlitevec ...`，并说明 SQLite 开发头文件要求 |
| Milvus | SDK 测试置于 Linux/WSL 专用 build tag；真实连接使用独立集成测试 | `-tags=milvus`、`EINO_RUN_MILVUS_TESTS=1`、`MILVUS_ADDRESS` |

环境变量只描述配置，绝不把 API Key、密码、内部 URL 写入 Plan、测试断言或 Git。`.env.example` 只保留变量名和无敏感样例。

## 如何决定是否使用 Eino

先看本节要学的边界，而不是专题名称：

- 学 `schema.Message`、Prompt、Tool、`ChatModelAgent` 或 ADK 编排时，使用 Eino。
- 学文档解析、分块、向量格式、检索评分、引用上下文、SQLite/Milvus 数据模型时，优先使用普通 Go interface 与对应官方库。
- 真实模型层可通过一个小 adapter 接入 Eino；不要为了“全项目都用 Eino”把数据层和存储层硬塞进 Agent。

这使底层检索链既可以被命令行或 HTTP 服务调用，也可以被 Eino Agent 调用。

## 编写顺序

新专题按以下顺序落地，避免先写长文、最后发现代码无法验证：

1. 明确专题边界、受众和完成后的可观察成果；创建 README 路线表。
2. 先写 Plan 1 的测试，再写能让它通过的最小代码和 Plan 文档。
3. 按编号重复：每次新增一个能力、一个测试入口、一份文档。
4. 若涉及真实语料，先实现盘点、格式提取、内容哈希、来源/页码 metadata；随后才开始分块和向量化。
5. 先实现可解释的内存/顺序基线，再接 SQLite-vec、Milvus、重排或 Agent。
6. 最后增加可选 live/integration 测试，并在 README 写明费用、网络、CGO、Docker 或服务前置条件。

## 提交前检查表

- [ ] 目录、Plan 编号、测试编号和 README 路线表完全一致。
- [ ] 每份 Plan 只引入一个新概念，且能链接到同编号测试。
- [ ] 默认 `go test ./<topic>/testfiles -v` 不使用网络、不需要 API Key。
- [ ] 数据驱动专题的每一个 Plan 主断言都从真实语料入口开始。
- [ ] 文档解析记录文件格式、稳定来源、内容哈希；PDF 记录页码；扫描 PDF 有 OCR 状态。
- [ ] 所有检索测试检查 metadata/租户等硬过滤，而不是只检查“看起来相关”。
- [ ] 生成测试验证证据与引用边界；无证据时有明确拒答路径。
- [ ] 外部服务测试有显式开关、超时、清理范围和不含敏感正文的日志。
- [ ] `gofmt`、专题测试、`go mod tidy`、`go mod verify` 已执行。
