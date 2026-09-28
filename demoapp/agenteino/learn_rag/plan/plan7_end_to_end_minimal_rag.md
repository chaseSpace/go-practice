# Plan 7：端到端最小 RAG——真实资料到带引用答案

## 本节只解决什么

将真实 TXT/PDF/DOCX 的解析、分块、离线 embedding、检索、证据上下文和确定性回答器连成闭环。

```text
learn_rag/data -> loadSourceDocument -> chunkDocument -> LocalRAG
真实 chunk query -> Retrieve -> evidence -> answer + [S1]
```

## 当前代码映射

[`../testfiles/plan7_test.go`](../testfiles/plan7_test.go) 使用全部真实 documents，按 `MaxRunes=300` 分块，以第一个真实 chunk 的完整内容作为 query。它断言回答含 `[S1]`，且首条 evidence URI 属于原始第一份资料。随后移除这份真实 source、重建 `LocalRAG`，用同一 query 断言返回“资料不足”。

当前入口是 `newLocalRAG` 和 `LocalRAG.Ask`，不是 CLI 命令；它的职责是验证链路边界。真实 embedding 与 LLM adapter 尚未在本节实现。

## 运行与验收

```bash
go test ./learn_rag/testfiles -run TestPlan7 -v
```

通过表示源文档确实影响回答；删除 source 后旧证据不能继续被回答器使用。
