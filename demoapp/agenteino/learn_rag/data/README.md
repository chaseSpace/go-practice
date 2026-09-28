# 本地 RAG 语料目录

当前目录可以直接放实际的 `.txt`、`.pdf`、`.docx`；也可以放进 `raw/`，例如：

```text
learn_rag/data/raw/
├── product-manual.txt
├── policy-2026.pdf
└── archived-notes/
    └── faq.txt
```

直接放在 `data/` 根目录的资料遵循仓库现有的 Git 管理方式；`raw/` 与 `derived/` 已被 Git 忽略，适合存放私有原文，以及提取文本、manifest、SQLite 文件和评测运行结果。若资料不能复制到仓库，可在未来实现时用 `RAG_SOURCE_DIR` 指向原目录。

首次执行索引前先完成 Plan 2 的盘点：当前测试会列出并真实提取 TXT/PDF/DOCX，再把提取文本送入分块。PDF 要记录是否为可提取文本或扫描件、页数/提取字符数、来源 URI、权限范围和内容哈希；扫描件提取为空时明确标为 OCR 待处理，不进入向量索引。不要将 API key、账号凭据或不该被问答系统检索的资料放入这里。
