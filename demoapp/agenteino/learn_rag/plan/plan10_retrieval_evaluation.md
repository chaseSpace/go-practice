# Plan 10：检索评测——用真实 chunk 计算 Recall@K 与 MRR

## 本节只解决什么

把“检索看起来不错”变为可重复计算的检索指标，并把无权限问题从可答问题中分开。

## 当前实现

`evaluateRetriever` 接收 `Retriever` 和 `[]EvalCase`，对每个 case 计算：

| 指标 | 当前含义 |
| --- | --- |
| `RecallAtK` | 目标 chunk 是否进入 Top-K |
| `MRR` | 第一个目标 chunk 的倒数排名 |
| `UnanswerablePassed` | 无权限/不可答 case 是否返回零候选 |

## 当前代码映射

[`../testfiles/plan10_test.go`](../testfiles/plan10_test.go) 从真实 corpus 选取首尾两个 chunk；每个完整 chunk 文本作为 query，并将自身 ID 作为 gold evidence。第三个 case 使用同一份真实 query，但给一个不存在的租户，作为无权限对照。测试期望 `RecallAtK=1`、`MRR=1`、`UnanswerablePassed=1`。

这是一条确定性 smoke baseline，不是人工标注的生产评测集。后续应把真实业务问题、多个允许 evidence、页码和人工答案写入版本化评测文件。

## 运行与验收

```bash
go test ./learn_rag/testfiles -run TestPlan10 -v
```

先用指标判断解析、分块、filter 或排序是否退化，再考虑混合检索、rerank 或更换模型。
