# Plan 12：Milvus（二）——后端配置与生产请求边界

## 本节只解决什么

验证 SQLite/Milvus 后端配置的最低要求，并用真实 chunk 构造包含租户和可见性硬过滤的 Milvus 搜索请求。

## 当前代码映射

[`../testfiles/plan12_test.go`](../testfiles/plan12_test.go) 从真实 chunk 派生 collection 名、SQLite 路径、64 维 query vector、tenant 和 visibility，然后：

- `validateBackendConfig` 验证 SQLite 需要路径、Milvus 需要地址；
- 用真实向量构造 Milvus `SearchOption`；
- 断言 `tenant_id == $tenant and visibility == $visibility`、两个模板参数和 `chunk_id/source_uri/chunk_text` 输出字段不丢失。

缺少 Milvus 地址是本节唯一的局部配置错误对照；业务证据、过滤值和向量均来自现有语料。

## 运行与验收

```bash
go test -tags=milvus ./learn_rag/testfiles -run TestPlan12 -v
```

该测试当前只在 Linux/WSL 编译，避免 Milvus 2.5.x 间接依赖在 Windows 的 `MemoryInfoExStat.RSS/Shared` 编译问题。通过只说明请求在发送前具有正确的后端和权限边界。当前代码**尚未**实现 SQLite/Milvus 的运行时切流、回填、灰度、重试、指标或真实 ANN 对比；这些应在对应实现和集成测试加入后再写入课程结论。

## 生产检查表

- 同一 source snapshot、chunk policy、embedding 版本与评测集再比较后端；
- 将 tenant/visibility 放在数据库 filter，而不是交给模型自觉遵守；
- 为 embedding、检索、重排、LLM 分别设置超时与日志；
- 记录索引版本、文档删除状态与引用缺失率；
- 外部服务不可用时明确降级，不伪造“有证据”的回答。
