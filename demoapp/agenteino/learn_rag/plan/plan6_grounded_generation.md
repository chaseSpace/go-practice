# Plan 6：证据上下文与引用——只回答检索到的真实内容

## 本节只解决什么

将真实检索结果组装为有预算、带来源的证据上下文，并验证无证据时的拒答路径。

## 当前实现

`buildEvidenceContext` 为每个候选生成：

```text
[S1] title=<真实标题> uri=<真实文件 URI> page=<PDF 页码>
<真实 chunk 正文>
```

当某段会超过 rune 预算时不加入上下文。`answerFromEvidence` 是确定性课程回答器：有证据时只引用首个 `[S*]`，无证据时返回“资料不足”。它不是 LLM，也不模拟真实生成质量。

## 当前代码映射

[`../testfiles/plan6_test.go`](../testfiles/plan6_test.go) 对真实首 chunk 检索，使用 2000 rune 预算构建 context，断言 `[S1]` 和真实 URI 出现；随后验证回答引用只包含 `S1`，并以空 evidence 作为局部拒答对照。

## 运行与验收

```bash
go test ./learn_rag/testfiles -run TestPlan6 -v
```

生产中用真实 ChatModel 替换 `answerFromEvidence` 时，仍须保留“只用证据、无证据拒答、引用只能来自输入标签”的约束与测试。
