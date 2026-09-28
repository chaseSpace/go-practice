# Plan 8：SQLite-vec（一）——真实 chunk 的本地持久化检索

## 本节只解决什么

将真实 chunk 和其确定性测试向量写入 sqlite-vec `vec0` 表，验证租户隔离、Top-K 与关闭重开后的持久化。

## 当前实现

[`../testfiles/sqlite_vec_helpers_test.go`](../testfiles/sqlite_vec_helpers_test.go) 定义当前真实 DDL：

```text
rag_vectors(
  chunk_id text primary key,
  tenant_id text partition key,
  document_id text,
  embedding float[3],
  +source_uri text,
  +chunk_text text
)
```

`courseVector` 从真实 chunk 文本确定性导出三维存储测试向量；它不是生产 embedding。`sqliteVecStore.Upsert` 使用“先 delete 再 insert”，避免把 `INSERT OR REPLACE` 当作 sqlite-vec 的安全 upsert。

## 当前代码映射

[`../testfiles/plan8_test.go`](../testfiles/plan8_test.go) 取三个真实 chunk，其中一个仅改变租户作为隔离对照。它验证 `vec_version()`、写入、错误维度、以真实向量查询的首命中、租户隔离，以及关闭重开后仍能检索。

## 运行与验收

该测试有 CGO 依赖：

```bash
CGO_ENABLED=1 go test -tags=sqlitevec ./learn_rag/testfiles -run TestPlan8 -v
```

运行前需安装 SQLite 开发头文件和 C toolchain。默认 `go test ./learn_rag/testfiles` 不编译此 tag；当前环境不满足该系统前置条件时，这是预期行为。
