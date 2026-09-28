# Plan 11：Milvus（一）——用真实 chunk 构造 collection 与搜索请求

## 本节只解决什么

把真实 chunk 的 ID、来源、租户和向量映射到 Milvus collection schema 与参数化搜索请求；默认不修改外部 Milvus。

## 当前代码映射

[`../testfiles/plan11_test.go`](../testfiles/plan11_test.go) 从真实 chunk 取得文本、tenant 和 64 维 hash vector，然后：

1. `milvusRAGSchema(64)` 构造 `chunk_id`、`document_id`、`tenant_id`、`source_uri`、`chunk_text`、`embedding` 字段；
2. 用 `AUTOINDEX + COSINE` 构造创建 collection 请求；
3. 用真实向量、真实 tenant 和参数化 filter 构造 `SearchOption`；
4. 断言 collection 名、schema、filter、输出 evidence 字段都存在。

`TestPlan11MilvusServiceConnection` 也先读取真实 chunk，再在 `EINO_RUN_MILVUS_TESTS=1` 时只发 `HasCollection` 连通性请求；它不创建、写入或删除 collection。

## 运行与验收

```bash
go test -tags=milvus ./learn_rag/testfiles -run TestPlan11MilvusSchemaAndSearchOption -v
```

默认运行 schema/request 测试并跳过服务连接。连接已有服务：

```bash
EINO_RUN_MILVUS_TESTS=1 MILVUS_ADDRESS=127.0.0.1:19530 \
  go test -tags=milvus ./learn_rag/testfiles -run TestPlan11MilvusServiceConnection -v
```

这些测试当前只在 Linux/WSL 编译。Milvus 2.5.x 的间接依赖在 Windows 上存在 `MemoryInfoExStat.RSS/Shared` 编译问题，因此 Windows 默认运行 RAG 基础 Plan 时不会加载 Milvus SDK。真正的 `MilvusStore`、collection 创建、插入、load、search 和受限清理尚未实现；应先在独立、唯一命名的集成测试中实现，避免 drop 共享 collection。
