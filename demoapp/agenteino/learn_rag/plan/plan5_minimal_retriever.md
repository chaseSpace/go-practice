# Plan 5：最小检索器——真实 chunk 的过滤、排序与 Top-K

## 本节只解决什么

实现内存基线 `InMemoryRetriever`，验证真实资料经过 query embedding、租户过滤、metadata 过滤、排序和 `TopK` 后的结果。

```text
真实 query chunk -> HashEmbedder -> filter(tenant/metadata) -> cosine sort -> Top-K
```

## 当前代码映射

[`../testfiles/plan5_test.go`](../testfiles/plan5_test.go) 从真实 corpus 分块后：

1. 为真实 chunk 写入 `status=current` metadata；
2. 把一个真实 chunk 复制为 `other-course` 租户对照；
3. 用首个真实 chunk 的完整文本查询；
4. 断言自身 chunk 排第一且结果不泄露其他租户；
5. 用相同真实 query 验证 `TopK=0` 会失败。

`SearchRequest.MinScore` 在本课程测试设为 `0.9999`，使完整同一 chunk 的 hash query 成为可重复的精确基线；它不是生产阈值建议。

## 运行与验收

```bash
go test ./learn_rag/testfiles -run TestPlan5 -v
```

这一步尚未评测自然语言答案，只证明候选证据的权限范围与排序契约。
