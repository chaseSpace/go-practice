# Plan 4：向量契约——同一真实 chunk 必须最接近自身

## 本节只解决什么

理解查询和文档向量必须共享模型、维度和度量约定，并验证余弦相似度的基本边界。

## 当前实现

`HashEmbedder` 对 `learn_rag/data/` 已提取的中英文真实 chunk 生成字符 1-gram/2-gram 哈希向量。它只用于无 API Key 的离线课程测试，不是语义 embedding 模型。它实现：

```go
type Embedder interface {
    EmbedDocuments(context.Context, []string) ([][]float32, error)
    EmbedQuery(context.Context, string) ([]float32, error)
    Dimension() int
    ModelID() string
}
```

`cosineSimilarity` 拒绝空向量和不同维度。将来接入真实 embedding 时，替换 `HashEmbedder`，同时重建该模型版本的全部索引。

## 当前代码映射

[`../testfiles/plan4_test.go`](../testfiles/plan4_test.go) 从真实语料取首尾两个 chunk，以 64 维 `HashEmbedder` 比较：首 chunk 的 query 与自身相似度为 1，且不低于另一个真实 chunk；最后用不同维度数组作为边界对照，断言报错。

## 运行与验收

```bash
go test ./learn_rag/testfiles -run TestPlan4 -v
```

不要把不同 embedding 模型或不同维度的分数直接比较，更换模型不是改配置而是索引迁移。
