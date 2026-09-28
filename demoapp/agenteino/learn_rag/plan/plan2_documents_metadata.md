# Plan 2：提取 TXT/PDF/DOCX——把真实资料变成可追溯文档

## 本节只解决什么

从 `learn_rag/data/`（或 `RAG_SOURCE_DIR`）真实提取 TXT、PDF、DOCX，构造含来源、内容哈希、格式和页码的 `SourceDocument`。本节不做向量检索。

## 当前实现

| 格式 | 当前代码 | 保留的信息 |
| --- | --- | --- |
| TXT | `parseTXT` | UTF-8 正文、文件 URI、内容哈希 |
| PDF | `ledongthucPDFExtractor` + `parsePDF` | 按页文字、`PDFPage.Number`、OCR 状态 |
| DOCX | `parseDOCX` + `nguyenthenguyen/docx` + `wordXMLToText` | WordprocessingML 段落、文件 URI、内容哈希 |

`inventoryCorpus` 只发现允许的扩展名并计算原文件哈希；`loadSourceDocument` 才按扩展名调用真实提取器。扫描 PDF 若文字过少会被标为 `OCRRequired`，`requireCourseDocuments` 会拒绝把它送进后续索引测试。

## 当前代码映射

[`../testfiles/plan2_test.go`](../testfiles/plan2_test.go) 通过 `requireCourseDocuments` 实际加载当前目录中的三种格式，并断言：每份资料都有 ID、URI、内容哈希、正文、`source_path` 和 `format` metadata；PDF 还必须有页列表。

## 运行与验收

```bash
go test ./learn_rag/testfiles -run TestPlan2 -v
```

日志只输出格式统计，不输出私有正文。当前资料应显示 `txt`、`pdf`、`docx` 均被提取。

## 下一步与边界

本实现提取可复制文字的 PDF；扫描件需要单独接 OCR，并把 OCR 版本和页码写入 metadata。不要把二进制文件内容、权限规则或 API Key 塞进 chunk 正文。
