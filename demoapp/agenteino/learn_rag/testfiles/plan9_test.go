//go:build sqlitevec

package testfiles

import (
	"path/filepath"
	"testing"
)

func TestPlan9SQLiteVecReindexAndDeleteLifecycle(t *testing.T) {
	// 以独立临时库演示一个文档在“首次索引 → 内容更新 → 删除”时的完整生命周期。
	store := newSQLiteVecStore(t, filepath.Join(t.TempDir(), "rag.sqlite"), 3)
	defer store.Close()

	// 所有写入对象均来自真实课程资料的分块结果。
	chunks := requireCourseChunks(t, ChunkPolicy{MaxRunes: 300, OverlapRunes: 50})
	chunk := chunks[0]
	// 初次 Upsert 建立“chunk ID → 文本与向量”的索引记录。
	if err := store.Upsert(chunk, courseVector(chunk.Text, 3)); err != nil {
		t.Fatalf("initial index: %v", err)
	}
	// 用同一文档的另一切片模拟内容更新，并保留原 ID 以触发替换语义。
	replacement := anotherChunkFromDocument(t, chunks, chunk)
	replacement.ID = chunk.ID
	if err := store.Upsert(replacement, courseVector(replacement.Text, 3)); err != nil {
		t.Fatalf("reindex changed chunk: %v", err)
	}
	// 更新后只应存在一条新文本；旧版本既不能重复出现，也不能继续被召回。
	results, err := store.Search(courseVector(replacement.Text, 3), "course", 5)
	if err != nil || len(results) != 1 || results[0].Text != replacement.Text {
		t.Fatalf("upsert must replace, not duplicate, changed chunks: %#v, %v", results, err)
	}

	// 删除以 document ID 为边界，避免某个源文件被移除后遗留可检索的敏感切片。
	if err := store.DeleteByDocumentID(chunk.DocumentID); err != nil {
		t.Fatalf("delete source document: %v", err)
	}
	// 删除后再次查询，结果必须为空，证明索引与源文档生命周期同步。
	results, err = store.Search(courseVector(replacement.Text, 3), "course", 5)
	if err != nil {
		t.Fatalf("search after delete: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("deleted source remains retrievable: %#v", results)
	}
}
