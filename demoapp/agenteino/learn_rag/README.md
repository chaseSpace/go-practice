# Eino RAG 由浅入深学习计划

本专题承接 `base_loop`：假定你已经能构造 Prompt、调用模型、解析输出，并理解一次请求的消息上下文。这里的 RAG（Retrieval-Augmented Generation）是一个**按问题从受管语料中取回证据，再让模型基于证据作答**的系统：

```text
文档 -> 清洗/分块 -> 向量化 -> 写入索引
问题 -> 向量化 -> 检索候选 -> 组装证据 -> 受约束生成 -> 带来源的回答
```

它不是“记忆”教程。本专题不实现聊天记录压缩、用户画像、跨会话偏好、长期记忆写入策略或自动回忆；这些属于状态/记忆设计，和以文档为中心、可重建可评估的 RAG 索引有不同的权限、更新和验证要求。

每份 `plan/planN_*.md` 只推进一个能力。你的 TXT/PDF/DOCX 是本专题的主语料：先从它们中挑一个小而有代表性的样本跑通，再扩展到全量。每个计划已对应 `testfiles/planN_test.go`，且主断言都经由 `requireCourseDocuments` / `requireCourseChunks` 读取 `learn_rag/data/` 中的真实资料；只有错误维度、空 evidence、无权限租户等同一 Plan 内的边界对照是人工构造。离线 hash vector 只替代有成本的 embedding API，不替代语料。真实 embedding/LLM 验证使用单独开关，避免日常 `go test ./...` 产生费用或依赖外网。

## 放置你的 TXT/PDF

当前已放入 `learn_rag/data/` 的 TXT/PDF/DOCX 会被 Plan 2 盘点；根目录内资料遵循仓库现有的 Git 管理方式。私有原文、提取文本和本地索引应放进 `raw/`、`derived/`（两者已被 Git 忽略），或通过 `RAG_SOURCE_DIR` 指向原有资料目录。文件命名尽量稳定，首次盘点后生成 manifest，后续用内容哈希而不是文件名判断是否需要重建。

Plan 2 使用真实库读取三种格式：UTF-8 TXT 直接读取，PDF 通过 `github.com/ledongthuc/pdf` 逐页提取，DOCX 通过 `github.com/nguyenthenguyen/docx` 打开并转换 WordprocessingML 段落。建议先选 10–30 份可公开或已脱敏的代表性资料：包含正常可复制文字的 PDF、长短不同文本、相近主题和一两个无答案的问题。扫描 PDF 不会因为扩展名是 `.pdf` 就有可检索文本；需要记录页数、提取字符数和 OCR 状态，提取为空时走 OCR/人工处理，而不是进入向量化。

## 路线

| Plan | 主题 | 交付物 / 验收焦点 | 可运行代码 |
| --- | --- | --- | --- |
| 1 | RAG 是什么，以及它不是什么 | 清楚区分知识检索、上下文和记忆 | `testfiles/plan1_test.go` |
| 2 | TXT/PDF/DOCX 语料、解析与元数据 | 真实库提取、可追溯、可删除、可过滤的文档模型 | `testfiles/plan2_test.go` |
| 3 | 文本清洗与分块 | 保留来源、页码和标题路径的稳定 chunk 方案 | `testfiles/plan3_test.go` |
| 4 | Embedding、相似度与向量 | 维度/模型/度量一致的向量契约 | `testfiles/plan4_test.go` |
| 5 | 最小检索器与检索指标 | 离线可验证的 Top-K 召回基线 | `testfiles/plan5_test.go` |
| 6 | 证据上下文与受约束生成 | 有引用、会拒答的回答组装 | `testfiles/plan6_test.go` |
| 7 | 端到端最小 RAG | 本地 TXT/PDF 到答案的完整闭环 | `testfiles/plan7_test.go` |
| 8 | SQLite-vec：本地向量索引 | `vec0` 建表、写入、相似检索与过滤 | `testfiles/plan8_test.go` |
| 9 | SQLite-vec：可重建 RAG 应用 | 幂等索引、删除更新、回答与来源 | `testfiles/plan9_test.go` |
| 10 | 检索质量诊断与进阶检索 | 评测集、失败分类、重排/混合检索决策 | `testfiles/plan10_test.go` |
| 11 | Milvus：集合、索引与过滤 | 从 SQLite 数据模型映射到 Milvus | `testfiles/plan11_test.go` |
| 12 | Milvus：服务化 RAG 与上线边界 | 同一 RAG 契约下的后端切换与运行检查 | `testfiles/plan12_test.go` |

## 两个后端的定位

| 后端 | 在本系列中的位置 | 适合先学什么 | 不解决什么 |
| --- | --- | --- | --- |
| SQLite-vec | 单进程、本地优先的第一种真实索引 | SQL、数据生命周期、精确可观察的检索链路 | 多节点服务、集群容量和高可用 |
| Milvus | 服务化向量库的第二种实现 | collection/schema、ANN 索引、过滤、部署边界 | 自动提升答案正确性或替代评测 |

SQLite-vec 是 pre-v1 项目，API 可能有破坏性变更；真正实现 Plan 8 时应把 bindings 与 SQLite driver 的版本一起锁进 `go.mod`，并记录 CGO/WASM 选择。Milvus 实战则固定一套 Milvus 服务镜像与 Go SDK 的兼容组合。两种后端只替换“存储与检索实现”，不会改变 chunk、embedding、证据组装和评测的基本责任。

## 测试与 Eino 的边界

`testfiles` 的核心 RAG 代码不依赖 Eino：文档盘点、分块、embedding 契约、检索、引用、SQLite-vec 与 Milvus schema 都是普通 Go 接口或各自官方 SDK。这样同一检索链能在命令行、HTTP 服务或 Eino 应用中复用。只有接入真实 embedding 模型和 ChatModel 时，才按需增加 Eino adapter；离线学习测试不应为了调用框架而引入模型费用。

默认运行不访问模型或 Milvus：

```bash
go test ./learn_rag/testfiles -v
```

从真实资料开始时，先单独运行 Plan 2；它会扫描 `learn_rag/data/`，并对 TXT、PDF、DOCX 逐一提取后分块：

```bash
go test ./learn_rag/testfiles -run TestPlan2 -v
```

Plan 8–9 的真实 SQLite-vec 测试使用 CGO bindings，需先安装 SQLite 开发头文件/C toolchain，然后显式运行：

```bash
CGO_ENABLED=1 go test -tags=sqlitevec ./learn_rag/testfiles -run 'TestPlan(8|9)' -v
```

Plan 11–12 直接引用 Milvus Go SDK，当前 SDK 的间接依赖不能在 Windows 上编译；它们仅在 Linux/WSL 中且显式启用 `milvus` tag 时参与编译。Plan 11 的服务连通性测试还需要 `EINO_RUN_MILVUS_TESTS=1` 和可访问的 `MILVUS_ADDRESS`：

```bash
go test -tags=milvus ./learn_rag/testfiles -run TestPlan11MilvusSchemaAndSearchOption -v
EINO_RUN_MILVUS_TESTS=1 MILVUS_ADDRESS=127.0.0.1:19530 \
  go test -tags=milvus ./learn_rag/testfiles -run TestPlan11MilvusServiceConnection -v
```

## 推荐推进方式

不要从数据库安装开始。先完成 Plan 1–7，得到一个能解释“为什么从你的 TXT/PDF 取回这些 chunk”的离线 RAG；再顺序完成 Plan 8–9 与 Plan 11–12。Plan 10 可在 SQLite-vec 闭环完成后执行，也应在迁移 Milvus 前完成，因为没有指标的迁移无法证明检索变好或至少没有退化。

## 参考资料

- [Eino Document Parser 接口](https://www.cloudwego.io/docs/eino/core_modules/components/document_loader_guide/document_parser_interface_guide/)
- [Eino Retriever 生态组件](https://www.cloudwego.io/docs/eino/ecosystem_integration/retriever/)
- [ledongthuc/pdf：Go PDF 文本与按页读取](https://github.com/ledongthuc/pdf)
- [nguyenthenguyen/docx：Go DOCX 打开与读取](https://github.com/nguyenthenguyen/docx)
- [sqlite-vec 项目与 vec0 表](https://github.com/asg017/sqlite-vec)
- [sqlite-vec Go bindings（CGO 与 ncruces/WASM）](https://github.com/asg017/sqlite-vec-go-bindings)
- [Milvus 创建 collection / 索引](https://milvus.io/docs/create-collection.md)
- [Milvus Go SDK 安装](https://milvus.io/docs/install-go.md)
