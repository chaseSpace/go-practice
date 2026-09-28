# Plan 9：SQLite-vec（二）——真实文档的更新与删除生命周期

## 本节只解决什么

验证真实文档 chunk 的替换和按 `document_id` 删除，避免索引留下旧内容。

## 当前代码映射

[`../testfiles/plan9_test.go`](../testfiles/plan9_test.go) 从真实 corpus 找到一个 chunk，并从**同一真实文档**找到另一个内容不同的 chunk。测试以原 chunk ID 写入第一段，再以相同 ID 写入第二段真实正文和对应向量，断言查询只返回新正文；最后 `DeleteByDocumentID` 并断言零命中。

这覆盖当前 `sqliteVecStore` 的真实 upsert/delete 行为。它尚未实现文件监听、manifest 增量比较、回答生成或通用 `Retriever` adapter；这些应在新增代码前再扩展 Plan，不应写成已完成能力。

## 运行与验收

```bash
CGO_ENABLED=1 go test -tags=sqlitevec ./learn_rag/testfiles -run TestPlan9 -v
```

与 Plan 8 相同，需要 SQLite 开发头文件和 C toolchain。
