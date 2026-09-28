package testfiles

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	pdfreader "github.com/ledongthuc/pdf"
	worddocx "github.com/nguyenthenguyen/docx"
)

// SourceDocument 表示从一个真实源文件提取出的标准文档。
// 它不保存向量：向量属于可重建索引，文档才是事实来源。
type SourceDocument struct {
	ID          string            // 稳定文档标识；chunk、更新和删除都通过它关联源文档。
	URI         string            // 可回查的源地址；当前本地文件使用 file:// 路径。
	Title       string            // 展示给引用结果的标题；当前取源文件名。
	Format      string            // 规范化的小写格式：txt、pdf 或 docx。
	Text        string            // 从源文件提取并清理后的完整文本。
	ContentHash string            // 提取文本的 SHA-256；用于检测可索引内容是否变化。
	TenantID    string            // 检索前必须应用的租户隔离字段。
	Visibility  string            // 文档可见范围；当前课程统一为 internal。
	Metadata    map[string]string // 可过滤或追踪的扩展字段，如 format、source_path。
	Pages       []PDFPage         // PDF 的逐页文本；非 PDF 文档为空。
}

// PDFPage 保存 PDF 的单页文本，使后续 chunk 和引用能够回到具体页码。
type PDFPage struct {
	Number int    // 从 1 开始的 PDF 页码。
	Text   string // 从该页内容流提取出的纯文本。
}

// Chunk 是检索和向量化的最小文本单元，同时保留源文档追踪字段。
type Chunk struct {
	ID         string            // 由文档 ID、顺序和正文生成的稳定 chunk ID。
	DocumentID string            // 所属 SourceDocument.ID，用于级联更新和删除。
	URI        string            // 原始文件地址，生成引用时直接使用。
	Title      string            // 原始文档标题。
	Text       string            // 实际参与 embedding、检索和证据拼装的正文。
	Index      int               // chunk 在当前文档中的零基顺序。
	PageStart  int               // PDF 起始页；TXT/DOCX 为 0。
	PageEnd    int               // PDF 结束页；当前按页分块时与 PageStart 相同。
	TenantID   string            // 从源文档继承的硬隔离字段。
	Visibility string            // 从源文档继承的可见范围。
	Metadata   map[string]string // 从源文档复制，避免后续修改共享 map。
}

// RetrievedChunk 是一个带相关性分数的检索结果。
type RetrievedChunk struct {
	Chunk         // 保留完整 chunk，保证结果仍然可引用和鉴权。
	Score float32 // 排序分数；其方向和范围由具体检索实现定义。
}

// ManifestRecord 记录一次格式提取的结果和质量信号。
type ManifestRecord struct {
	DocumentID     string // 与生成的 SourceDocument.ID 一致。
	RelativePath   string // 本次载入使用的规范化源路径。
	Format         string // 实际走过的解析器格式。
	PageCount      int    // PDF 总页数；TXT/DOCX 为 0。
	ExtractedRunes int    // 成功提取的 Unicode 字符数，用于发现空文档。
	OCRRequired    bool   // PDF 文字过少时为 true，阻止其进入普通文本索引。
	ContentHash    string // 提取文本的 SHA-256。
}

// CorpusEntry 是建立索引前的文件盘点记录。
// 它只说明发现了什么文件，不假定每个扩展名都已有文本提取器。
type CorpusEntry struct {
	RelativePath string // 相对于 corpus 根目录的路径，用于可移植 manifest。
	Format       string // 从扩展名识别的 txt、pdf 或 docx。
	ContentHash  string // 原始文件字节的 SHA-256，区别于提取文本哈希。
	ParserState  string // 对应格式解析器是否可用；当前允许格式均为 ready。
}

// inventoryCorpus 扫描根目录中的 TXT/PDF/DOCX，计算原文件哈希并稳定排序。
// 它只做盘点，不执行文本提取。
func inventoryCorpus(root string) ([]CorpusEntry, error) {
	entries := make([]CorpusEntry, 0)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || strings.EqualFold(entry.Name(), "README.md") {
			return nil
		}

		// todo
		if !strings.HasSuffix(entry.Name(), "pdf") {
			println("skip entry", entry.Name())
			return nil
		}
		format := strings.TrimPrefix(strings.ToLower(filepath.Ext(entry.Name())), ".")
		if format != "txt" && format != "pdf" && format != "docx" {
			return nil
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		entries = append(entries, CorpusEntry{
			RelativePath: filepath.ToSlash(relativePath),
			Format:       format,
			ContentHash:  contentHash(string(contents)),
			ParserState:  "ready",
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(left, right int) bool {
		return entries[left].RelativePath < entries[right].RelativePath
	})
	return entries, nil
}

// requireCourseDocuments 是每个 Plan 测试加载真实语料的统一入口。
// 课程刻意从学习者提供的 TXT/PDF/DOCX 开始；手写字符串只允许用于局部错误对照。
func requireCourseDocuments(t *testing.T) []SourceDocument {
	t.Helper()
	root := os.Getenv("RAG_SOURCE_DIR")
	if root == "" {
		root = filepath.Join("..", "data")
	}
	entries, err := inventoryCorpus(root)
	if err != nil {
		t.Fatalf("inventory course corpus %q: %v", root, err)
	}
	if len(entries) == 0 {
		t.Fatalf("course corpus %q has no TXT/PDF/DOCX material", root)
	}

	documents := make([]SourceDocument, 0, len(entries))
	formats := map[string]bool{}
	for _, entry := range entries {
		document, manifest, err := loadSourceDocument(context.Background(), filepath.Join(root, filepath.FromSlash(entry.RelativePath)))
		if err != nil {
			t.Fatalf("extract course source %s: %v", entry.RelativePath, err)
		}
		if manifest.ExtractedRunes == 0 || (entry.Format == "pdf" && manifest.OCRRequired) {
			t.Fatalf("course source %s is not indexable: %#v", entry.RelativePath, manifest)
		}
		document.TenantID = "course"
		document.Visibility = "internal"
		document.Metadata = cloneMetadata(document.Metadata)
		document.Metadata["format"] = entry.Format
		document.Metadata["source_path"] = entry.RelativePath
		documents = append(documents, document)
		formats[entry.Format] = true
	}
	return documents
}

// requireCourseChunks 加载全部真实课程文档，并使用指定策略汇总其 chunk。
// 任意文档解析/分块失败都会立即终止当前测试。
func requireCourseChunks(t *testing.T, policy ChunkPolicy) []Chunk {
	t.Helper()
	var chunks []Chunk
	for _, document := range requireCourseDocuments(t) {
		documentChunks, err := chunkDocument(document, policy)
		if err != nil {
			t.Fatalf("chunk course source %s: %v", document.URI, err)
		}
		chunks = append(chunks, documentChunks...)
	}
	if len(chunks) < 2 {
		t.Fatal("course corpus needs at least two chunks for retrieval comparisons")
	}
	return chunks
}

// requireCourseDocumentFormat 返回课程语料中指定格式的第一份真实文档。
func requireCourseDocumentFormat(t *testing.T, format string) SourceDocument {
	t.Helper()
	for _, document := range requireCourseDocuments(t) {
		if document.Format == format {
			return document
		}
	}
	t.Fatalf("course corpus has no %s document", format)
	return SourceDocument{}
}

// queryFromChunk 使用真实 chunk 全文作为确定性检索 query。
// 这样离线测试可以断言“同一文本最接近自身”，不依赖自然语言模型随机性。
func queryFromChunk(chunk Chunk) string {
	return strings.TrimSpace(chunk.Text)
}

// courseVector 从真实 chunk 文本确定性导出存储测试向量。
// 它用于验证向量持久化，不冒充语义 embedding，也不硬编码演示向量。
func courseVector(text string, dimension int) []float32 {
	if dimension <= 0 || strings.TrimSpace(text) == "" {
		panic("course vector needs a dimension and non-empty real source text")
	}
	vector := make([]float32, dimension)
	for index := range vector {
		sum := sha256.Sum256([]byte(fmt.Sprintf("%d\x00%s", index, text)))
		vector[index] = float32(int(sum[0])<<8|int(sum[1]))/65535 + 0.001
	}
	return vector
}

// anotherChunkFromDocument 查找同一真实文档中的另一个 chunk，供更新测试使用。
func anotherChunkFromDocument(t *testing.T, chunks []Chunk, original Chunk) Chunk {
	t.Helper()
	for _, chunk := range chunks {
		if chunk.DocumentID == original.DocumentID && chunk.ID != original.ID && chunk.Text != original.Text {
			return chunk
		}
	}
	t.Fatalf("source document %s needs a second real chunk for reindex comparison", original.DocumentID)
	return Chunk{}
}

// newSourceDocument 统一规范化提取文本并填充文档默认值。
func newSourceDocument(id, uri, title, format, text, tenantID string) SourceDocument {
	text = strings.TrimSpace(text)
	return SourceDocument{
		ID:          id,
		URI:         uri,
		Title:       title,
		Format:      format,
		Text:        text,
		ContentHash: contentHash(text),
		TenantID:    tenantID,
		Visibility:  "internal",
		Metadata:    map[string]string{},
	}
}

// contentHash 返回文本内容的十六进制 SHA-256。
func contentHash(text string) string {
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:])
}

// stableID 将多个稳定字段用不可见分隔符连接后生成短 ID。
func stableID(parts ...string) string {
	return contentHash(strings.Join(parts, "\x00"))[:24]
}

// cloneMetadata 复制 metadata，避免 SourceDocument 与 Chunk 共享可变 map。
func cloneMetadata(values map[string]string) map[string]string {
	cloned := make(map[string]string, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}

// parseTXT 读取 UTF-8 TXT，并同时返回标准文档和提取清单。
func parseTXT(path string) (SourceDocument, ManifestRecord, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return SourceDocument{}, ManifestRecord{}, fmt.Errorf("read txt: %w", err)
	}
	if !utf8.Valid(contents) {
		return SourceDocument{}, ManifestRecord{}, errors.New("txt is not valid UTF-8")
	}
	text := strings.TrimSpace(string(contents))
	if text == "" {
		return SourceDocument{}, ManifestRecord{}, errors.New("txt is empty")
	}
	cleanPath := filepath.ToSlash(filepath.Clean(path))
	document := newSourceDocument(
		stableID("txt", cleanPath),
		"file://"+cleanPath,
		filepath.Base(path),
		"txt",
		text,
		"default",
	)
	return document, ManifestRecord{
		DocumentID:     document.ID,
		RelativePath:   cleanPath,
		Format:         "txt",
		ExtractedRunes: utf8.RuneCountInString(text),
		ContentHash:    document.ContentHash,
	}, nil
}

// PDFTextExtractor 隔离 PDF 解析库与下游 RAG 数据模型。
// 返回值必须保留从 1 开始的页码，即使某页没有可提取文字也不能悄悄改号。
type PDFTextExtractor interface {
	Extract(ctx context.Context, path string) ([]PDFPage, error)
}

// ledongthucPDFExtractor 使用 github.com/ledongthuc/pdf 提取逐页文本。
type ledongthucPDFExtractor struct{}

// Extract 按 PDF 视觉行顺序拼接每页文本，并在页与页之间保持明确边界。
func (ledongthucPDFExtractor) Extract(ctx context.Context, path string) ([]PDFPage, error) {
	file, reader, err := pdfreader.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if reader.NumPage() == 0 {
		return nil, errors.New("pdf has no pages")
	}

	pages := make([]PDFPage, 0, reader.NumPage())
	for pageNumber := 1; pageNumber <= reader.NumPage(); pageNumber++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		page := reader.Page(pageNumber)
		if page.V.IsNull() || page.V.Key("Contents").Kind() == pdfreader.Null {
			pages = append(pages, PDFPage{Number: pageNumber})
			continue
		}
		rows, err := page.GetTextByRow()
		if err != nil {
			return nil, fmt.Errorf("read page %d: %w", pageNumber, err)
		}
		lines := make([]string, 0, len(rows))
		for _, row := range rows {
			var line strings.Builder
			for _, text := range row.Content {
				line.WriteString(text.S)
			}
			if value := strings.TrimSpace(line.String()); value != "" {
				lines = append(lines, value)
			}
		}
		pages = append(pages, PDFPage{Number: pageNumber, Text: strings.Join(lines, "\n")})
	}
	return pages, nil
}

// parsePDF 将逐页提取结果转换为 SourceDocument，并计算页数、字符数和 OCR 信号。
func parsePDF(ctx context.Context, path string, extractor PDFTextExtractor) (SourceDocument, ManifestRecord, error) {
	pages, err := extractor.Extract(ctx, path)
	if err != nil {
		return SourceDocument{}, ManifestRecord{}, fmt.Errorf("extract pdf: %w", err)
	}
	if len(pages) == 0 {
		return SourceDocument{}, ManifestRecord{}, errors.New("pdf has no pages")
	}
	var textParts []string
	for _, page := range pages {
		if page.Number <= 0 {
			return SourceDocument{}, ManifestRecord{}, errors.New("pdf page number must be positive")
		}
		textParts = append(textParts, strings.TrimSpace(page.Text))
	}
	text := strings.TrimSpace(strings.Join(textParts, "\n\f\n"))
	document := newSourceDocument(
		stableID("pdf", filepath.ToSlash(filepath.Clean(path))),
		"file://"+filepath.ToSlash(filepath.Clean(path)),
		filepath.Base(path),
		"pdf",
		text,
		"default",
	)
	document.Pages = append([]PDFPage(nil), pages...)
	extractedRunes := utf8.RuneCountInString(text)
	return document, ManifestRecord{
		DocumentID:     document.ID,
		RelativePath:   filepath.ToSlash(filepath.Clean(path)),
		Format:         "pdf",
		PageCount:      len(pages),
		ExtractedRunes: extractedRunes,
		OCRRequired:    extractedRunes < 20,
		ContentHash:    document.ContentHash,
	}, nil
}

// parseDOCX 使用 docx 库打开 OOXML 包，再将 document.xml 转为段落文本。
func parseDOCX(ctx context.Context, path string) (SourceDocument, ManifestRecord, error) {
	if err := ctx.Err(); err != nil {
		return SourceDocument{}, ManifestRecord{}, err
	}
	replaceDoc, err := worddocx.ReadDocxFile(path)
	if err != nil {
		return SourceDocument{}, ManifestRecord{}, fmt.Errorf("open docx: %w", err)
	}
	defer replaceDoc.Close()
	text, err := wordXMLToText(replaceDoc.Editable().GetContent())
	if err != nil {
		return SourceDocument{}, ManifestRecord{}, fmt.Errorf("extract docx text: %w", err)
	}
	if strings.TrimSpace(text) == "" {
		return SourceDocument{}, ManifestRecord{}, errors.New("docx has no extractable text")
	}
	cleanPath := filepath.ToSlash(filepath.Clean(path))
	document := newSourceDocument(
		stableID("docx", cleanPath),
		"file://"+cleanPath,
		filepath.Base(path),
		"docx",
		text,
		"default",
	)
	return document, ManifestRecord{
		DocumentID:     document.ID,
		RelativePath:   cleanPath,
		Format:         "docx",
		ExtractedRunes: utf8.RuneCountInString(text),
		ContentHash:    document.ContentHash,
	}, nil
}

// wordXMLToText 将 docx 库返回的 WordprocessingML 转为按段落分隔的纯文本。
// OOXML 包由库负责打开；这个轻量规范化器只保留能成为 RAG 证据的文本，
// 并保留段落和表格行边界供后续分块使用。
func wordXMLToText(documentXML string) (string, error) {
	decoder := xml.NewDecoder(strings.NewReader(documentXML))
	var paragraphs []string
	var current strings.Builder
	flush := func() {
		if text := strings.TrimSpace(current.String()); text != "" {
			paragraphs = append(paragraphs, text)
		}
		current.Reset()
	}
	inText := false
	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", err
		}
		switch value := token.(type) {
		case xml.StartElement:
			switch value.Name.Local {
			case "t":
				inText = true
			case "tab":
				current.WriteByte('\t')
			case "br", "cr":
				current.WriteByte('\n')
			}
		case xml.CharData:
			if inText {
				current.Write([]byte(value))
			}
		case xml.EndElement:
			switch value.Name.Local {
			case "t":
				inText = false
			case "p", "tr":
				flush()
			}
		}
	}
	flush()
	return strings.Join(paragraphs, "\n\n"), nil
}

// loadSourceDocument 根据文件扩展名路由到 TXT、PDF 或 DOCX 解析器。
func loadSourceDocument(ctx context.Context, path string) (SourceDocument, ManifestRecord, error) {
	switch strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".") {
	case "txt":
		return parseTXT(path)
	case "pdf":
		return parsePDF(ctx, path, ledongthucPDFExtractor{})
	case "docx":
		return parseDOCX(ctx, path)
	default:
		return SourceDocument{}, ManifestRecord{}, fmt.Errorf("unsupported source extension: %s", filepath.Ext(path))
	}
}

// ChunkPolicy 定义确定性字符分块参数。
type ChunkPolicy struct {
	MaxRunes     int // 每个 chunk 允许的最大 Unicode 字符数。
	OverlapRunes int // 相邻 chunk 重复的字符数；必须小于 MaxRunes。
}

// chunkDocument 将标准文档切为稳定 chunk；PDF 按页独立切分并携带页码。
func chunkDocument(document SourceDocument, policy ChunkPolicy) ([]Chunk, error) {
	if policy.MaxRunes <= 0 {
		return nil, errors.New("chunk max runes must be positive")
	}
	if policy.OverlapRunes < 0 || policy.OverlapRunes >= policy.MaxRunes {
		return nil, errors.New("chunk overlap must be non-negative and smaller than max runes")
	}

	type segment struct {
		text string // 当前待切分的文档正文或单页正文。
		page int    // PDF 页码；非 PDF 文档为 0。
	}
	segments := []segment{{text: document.Text}}
	if len(document.Pages) > 0 {
		segments = make([]segment, 0, len(document.Pages))
		for _, page := range document.Pages {
			segments = append(segments, segment{text: page.Text, page: page.Number})
		}
	}

	chunks := make([]Chunk, 0)
	for _, segment := range segments {
		for _, text := range splitText(segment.text, policy.MaxRunes, policy.OverlapRunes) {
			index := len(chunks)
			chunks = append(chunks, Chunk{
				ID:         stableID(document.ID, fmt.Sprint(index), text),
				DocumentID: document.ID,
				URI:        document.URI,
				Title:      document.Title,
				Text:       text,
				Index:      index,
				PageStart:  segment.page,
				PageEnd:    segment.page,
				TenantID:   document.TenantID,
				Visibility: document.Visibility,
				Metadata:   cloneMetadata(document.Metadata),
			})
		}
	}
	if len(chunks) == 0 {
		return nil, errors.New("document produced no chunks")
	}
	return chunks, nil
}

// splitText 先规范化空白，再按 rune 而不是字节切分，并尽量落在空白边界。
func splitText(text string, maxRunes, overlapRunes int) []string {
	runes := []rune(strings.Join(strings.Fields(text), " "))
	var pieces []string
	for start := 0; start < len(runes); {
		end := start + maxRunes
		if end >= len(runes) {
			pieces = append(pieces, string(runes[start:]))
			break
		}
		// 优先落在词边界，但不能因为长词跨越边界而产生空 chunk 或极小 chunk。
		boundary := end
		for boundary > start+maxRunes/2 && !unicode.IsSpace(runes[boundary]) {
			boundary--
		}
		if boundary == start+maxRunes/2 {
			boundary = end
		}
		pieces = append(pieces, strings.TrimSpace(string(runes[start:boundary])))
		start = boundary - overlapRunes
	}
	return pieces
}

// Embedder 描述文档批量向量化和查询向量化必须共享的契约。
// Dimension 与 ModelID 用于阻止不同模型或维度的数据混入同一索引。
type Embedder interface {
	EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error)
	EmbedQuery(ctx context.Context, text string) ([]float32, error)
	Dimension() int
	ModelID() string
}

// HashEmbedder 是课程测试使用的确定性本地适配器。
// 它把真实源文本转换为字符 n-gram 向量，因此无需 API Key 也能测试中英文语料检索。
// 它不是语义 embedding 模型；在真实集成步骤应替换为实际 embedding 提供商。
type HashEmbedder struct {
	dimension int // 所有文档向量和查询向量的固定维度。
}

// newHashEmbedder 创建仅用于离线课程测试的字符 n-gram 哈希向量器。
func newHashEmbedder(dimension int) *HashEmbedder {
	if dimension <= 0 {
		panic("hash embedder dimension must be positive")
	}
	return &HashEmbedder{dimension: dimension}
}

// Dimension 返回当前向量契约的固定维度。
func (e *HashEmbedder) Dimension() int { return e.dimension }

// ModelID 返回可记录到索引 manifest 的确定性实现标识。
func (e *HashEmbedder) ModelID() string {
	return fmt.Sprintf("course-hash-ngram-%d", e.dimension)
}

// EmbedDocuments 按输入顺序批量生成文档向量；任一文本失败则整体返回错误。
func (e *HashEmbedder) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	vectors := make([][]float32, len(texts))
	for index, text := range texts {
		vector, err := e.embed(ctx, text)
		if err != nil {
			return nil, err
		}
		vectors[index] = vector
	}
	return vectors, nil
}

// EmbedQuery 使用与文档相同的算法生成单个查询向量。
func (e *HashEmbedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	return e.embed(ctx, text)
}

// embed 将规范化文本的 1-gram/2-gram 哈希到固定维度并保留正负符号。
func (e *HashEmbedder) embed(ctx context.Context, text string) ([]float32, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	runes := []rune(strings.ToLower(strings.Join(strings.Fields(text), " ")))
	if len(runes) == 0 {
		return nil, errors.New("cannot embed empty text")
	}
	vector := make([]float32, e.dimension)
	for width := 1; width <= 2; width++ {
		for start := 0; start+width <= len(runes); start++ {
			hasher := fnv.New64a()
			_, _ = hasher.Write([]byte(string(runes[start : start+width])))
			sum := hasher.Sum64()
			value := float32(1)
			if sum&1 == 1 {
				value = -1
			}
			vector[int(sum%uint64(e.dimension))] += value
		}
	}
	if _, err := cosineSimilarity(vector, vector); err != nil {
		return nil, fmt.Errorf("hash embedding is empty: %w", err)
	}
	return vector, nil
}

// cosineSimilarity 计算两个同维非零向量的余弦相似度。
func cosineSimilarity(left, right []float32) (float32, error) {
	if len(left) == 0 || len(left) != len(right) {
		return 0, errors.New("vectors must have the same non-zero dimension")
	}
	var dot, leftNorm, rightNorm float64
	for index := range left {
		dot += float64(left[index] * right[index])
		leftNorm += float64(left[index] * left[index])
		rightNorm += float64(right[index] * right[index])
	}
	if leftNorm == 0 || rightNorm == 0 {
		return 0, errors.New("zero vector has no cosine similarity")
	}
	return float32(dot / math.Sqrt(leftNorm*rightNorm)), nil
}

// SearchRequest 是后端无关的检索输入。
type SearchRequest struct {
	Query    string            // 用户问题或课程中取自真实 chunk 的查询文本。
	TenantID string            // 硬租户范围；非空时先过滤再计算候选。
	TopK     int               // 最多返回的候选数量，必须大于 0。
	MinScore float32           // 最低相似度阈值；低于它的候选直接丢弃。
	Filters  map[string]string // 对 Chunk.Metadata 的精确匹配条件。
}

// Retriever 是内存、SQLite-vec、Milvus 等检索实现共享的最小接口。
type Retriever interface {
	Retrieve(ctx context.Context, request SearchRequest) ([]RetrievedChunk, error)
}

// indexedChunk 将可引用的 chunk 与其索引向量绑定。
type indexedChunk struct {
	chunk  Chunk     // 可回查来源的原始 chunk。
	vector []float32 // 与 embedder 契约一致的索引向量。
}

// InMemoryRetriever 是用于解释和验证排序逻辑的穷举检索基线。
type InMemoryRetriever struct {
	embedder Embedder       // 同时用于索引文档和向量化查询。
	items    []indexedChunk // 当前内存中的所有 chunk/vector 对。
}

// newInMemoryRetriever 批量向量化真实 chunk，并检查数量和维度契约。
func newInMemoryRetriever(ctx context.Context, embedder Embedder, chunks []Chunk) (*InMemoryRetriever, error) {
	if len(chunks) == 0 {
		return nil, errors.New("cannot index zero chunks")
	}
	texts := make([]string, len(chunks))
	for index, chunk := range chunks {
		texts[index] = chunk.Text
	}
	vectors, err := embedder.EmbedDocuments(ctx, texts)
	if err != nil {
		return nil, fmt.Errorf("embed chunks: %w", err)
	}
	if len(vectors) != len(chunks) {
		return nil, errors.New("embedder returned an unexpected vector count")
	}
	retriever := &InMemoryRetriever{embedder: embedder, items: make([]indexedChunk, len(chunks))}
	for index, vector := range vectors {
		if len(vector) != embedder.Dimension() {
			return nil, fmt.Errorf("chunk %s has dimension %d, want %d", chunks[index].ID, len(vector), embedder.Dimension())
		}
		retriever.items[index] = indexedChunk{chunk: chunks[index], vector: vector}
	}
	return retriever, nil
}

// Retrieve 依次执行查询向量化、租户/metadata 过滤、余弦打分和 Top-K 截断。
func (r *InMemoryRetriever) Retrieve(ctx context.Context, request SearchRequest) ([]RetrievedChunk, error) {
	if strings.TrimSpace(request.Query) == "" {
		return nil, errors.New("search query is empty")
	}
	if request.TopK <= 0 {
		return nil, errors.New("top-k must be positive")
	}
	queryVector, err := r.embedder.EmbedQuery(ctx, request.Query)
	if err != nil {
		return nil, fmt.Errorf("embed query: %w", err)
	}
	if len(queryVector) != r.embedder.Dimension() {
		return nil, errors.New("query vector has an unexpected dimension")
	}
	results := make([]RetrievedChunk, 0, len(r.items))
	for _, item := range r.items {
		if request.TenantID != "" && item.chunk.TenantID != request.TenantID {
			continue
		}
		if !matchesFilters(item.chunk, request.Filters) {
			continue
		}
		score, err := cosineSimilarity(queryVector, item.vector)
		if err != nil {
			return nil, fmt.Errorf("score chunk %s: %w", item.chunk.ID, err)
		}
		if score < request.MinScore {
			continue
		}
		results = append(results, RetrievedChunk{Chunk: item.chunk, Score: score})
	}
	sort.Slice(results, func(left, right int) bool {
		if results[left].Score == results[right].Score {
			return results[left].ID < results[right].ID
		}
		return results[left].Score > results[right].Score
	})
	if len(results) > request.TopK {
		results = results[:request.TopK]
	}
	return results, nil
}

// matchesFilters 对 chunk metadata 执行全部条件都必须成立的精确匹配。
func matchesFilters(chunk Chunk, filters map[string]string) bool {
	for key, want := range filters {
		if chunk.Metadata[key] != want {
			return false
		}
	}
	return true
}

// Evidence 为检索结果分配稳定引用标签，供上下文和最终答案共同使用。
type Evidence struct {
	Label string         // 本次回答内的引用标签，如 S1、S2。
	Chunk RetrievedChunk // 标签对应的真实检索结果。
}

// buildEvidenceContext 按检索顺序拼装带来源标签的上下文，并严格遵守 rune 预算。
func buildEvidenceContext(results []RetrievedChunk, maxRunes int) (string, []Evidence, error) {
	if maxRunes <= 0 {
		return "", nil, errors.New("evidence budget must be positive")
	}
	var builder strings.Builder
	usedRunes := 0
	evidence := make([]Evidence, 0, len(results))
	for _, result := range results {
		label := fmt.Sprintf("S%d", len(evidence)+1)
		page := ""
		if result.PageStart > 0 {
			page = fmt.Sprintf(" page=%d", result.PageStart)
		}
		block := fmt.Sprintf("[%s] title=%s uri=%s%s\n%s\n\n", label, result.Title, result.URI, page, result.Text)
		blockRunes := utf8.RuneCountInString(block)
		if usedRunes+blockRunes > maxRunes {
			continue
		}
		builder.WriteString(block)
		usedRunes += blockRunes
		evidence = append(evidence, Evidence{Label: label, Chunk: result})
	}
	return builder.String(), evidence, nil
}

// Answer 是课程回答器的结构化结果。
type Answer struct {
	Text      string     // 面向调用方的答案文本。
	Citations []string   // 答案实际引用的 Evidence.Label 集合。
	Evidence  []Evidence // 生成答案时允许使用的完整证据集合。
}

// answerFromEvidence 故意保持确定性。
// 它只测试证据与引用边界，不假装真实 LLM 的输出是确定的。
func answerFromEvidence(question string, evidence []Evidence) Answer {
	if len(evidence) == 0 {
		return Answer{Text: "资料不足：没有检索到可用于回答的问题证据。"}
	}
	first := evidence[0]
	sentence := strings.Split(strings.TrimSpace(first.Chunk.Text), "。")[0]
	if sentence == "" {
		sentence = first.Chunk.Text
	}
	return Answer{
		Text:      fmt.Sprintf("根据 [%s]，%s。", first.Label, sentence),
		Citations: []string{first.Label},
		Evidence:  evidence,
	}
}

// LocalRAG 是“分块后索引 + 单次检索回答”的最小端到端封装。
type LocalRAG struct {
	retriever *InMemoryRetriever // 当前课程使用的离线穷举检索器。
}

// newLocalRAG 串联文档分块、批量向量化与内存索引构建。
func newLocalRAG(ctx context.Context, documents []SourceDocument, policy ChunkPolicy, embedder Embedder) (*LocalRAG, error) {
	var chunks []Chunk
	for _, document := range documents {
		documentChunks, err := chunkDocument(document, policy)
		if err != nil {
			return nil, fmt.Errorf("chunk %s: %w", document.ID, err)
		}
		chunks = append(chunks, documentChunks...)
	}
	retriever, err := newInMemoryRetriever(ctx, embedder, chunks)
	if err != nil {
		return nil, err
	}
	return &LocalRAG{retriever: retriever}, nil
}

// Ask 执行一次检索、证据预算组装和确定性回答，不保存任何会话记忆。
func (r *LocalRAG) Ask(ctx context.Context, request SearchRequest, budget int) (Answer, error) {
	results, err := r.retriever.Retrieve(ctx, request)
	if err != nil {
		return Answer{}, err
	}
	_, evidence, err := buildEvidenceContext(results, budget)
	if err != nil {
		return Answer{}, err
	}
	return answerFromEvidence(request.Query, evidence), nil
}

// EvalCase 描述一个检索评测样本及其允许命中的 gold chunk。
type EvalCase struct {
	Name             string        // 用于报告和错误定位的样本名。
	Request          SearchRequest // 本次评测使用的 query、权限范围和 Top-K。
	ExpectedChunkIDs []string      // 任一命中即可视为召回成功的 gold chunk ID。
	Answerable       bool          // false 表示预期零候选，例如无权限范围。
}

// EvalReport 汇总确定性检索指标。
type EvalReport struct {
	RecallAtK          float64 // 可答样本中至少召回一个 gold chunk 的比例。
	MRR                float64 // 可答样本首个 gold chunk 倒数排名的平均值。
	UnanswerablePassed int     // 正确返回零候选的不可答样本数量。
}

// evaluateRetriever 顺序运行评测样本，并分别统计可答和不可答结果。
func evaluateRetriever(ctx context.Context, retriever Retriever, cases []EvalCase) (EvalReport, error) {
	var report EvalReport
	answerableCount := 0
	for _, testCase := range cases {
		results, err := retriever.Retrieve(ctx, testCase.Request)
		if err != nil {
			return EvalReport{}, fmt.Errorf("evaluate %s: %w", testCase.Name, err)
		}
		if !testCase.Answerable {
			if len(results) == 0 {
				report.UnanswerablePassed++
			}
			continue
		}
		answerableCount++
		for rank, result := range results {
			if contains(testCase.ExpectedChunkIDs, result.ID) {
				report.RecallAtK++
				report.MRR += 1 / float64(rank+1)
				break
			}
		}
	}
	if answerableCount > 0 {
		report.RecallAtK /= float64(answerableCount)
		report.MRR /= float64(answerableCount)
	}
	return report, nil
}

// contains 判断字符串切片是否包含目标值，用于匹配 gold chunk ID。
func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
