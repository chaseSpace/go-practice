# Plan 3：分块——保留真实 PDF 的来源和页码

## 本节只解决什么

将真实文档按字符上限与重叠长度分块，并保证 chunk 的 ID、顺序、来源和 PDF 页码稳定。

## 当前实现

`chunkDocument(document, ChunkPolicy)` 会：

```text
TXT/DOCX 正文 -> 规范化空白 -> 按 MaxRunes/OverlapRunes 切分
PDF 每一页    -> 各页独立切分 -> Chunk.PageStart/PageEnd
```

`Chunk.ID` 由 `document.ID + chunk index + text` 的稳定哈希生成。当前实现优先保证可验证的字符切分和页码归属；Markdown 标题、表格和复杂版式的语义分块属于后续增强，不应被误认为已经实现。

## 当前代码映射

[`../testfiles/plan3_test.go`](../testfiles/plan3_test.go) 通过 `requireCourseDocumentFormat(t, "pdf")` 使用实际 PDF，以 `MaxRunes=24`、`OverlapRunes=4` 分块两次。它断言：chunk 数量大于一、索引连续、`DocumentID`/URI 未丢失、每个 PDF chunk 有合法页码，且重复分块得到相同 ID。

## 运行与验收

```bash
go test ./learn_rag/testfiles -run TestPlan3 -v
```

通过后，引用 PDF chunk 时至少能回到正确 PDF 和页码。若页眉页脚或多栏顺序污染正文，应先修 PDF 提取器，再调整 chunk 长度。
