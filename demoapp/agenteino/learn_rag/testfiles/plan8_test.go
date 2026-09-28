//go:build sqlitevec

package testfiles

import (
	"path/filepath"
	"testing"
)

func TestPlan8SQLiteVecStoresAndFiltersVectors(t *testing.T) {
	// 临时数据库使每次测试相互隔离，同时可在后面通过重新打开验证落盘持久化。
	path := filepath.Join(t.TempDir(), "rag.sqlite")
	// SQLite-vec 表的向量列固定为 3 维，写入和检索都必须遵守这个契约。
	store := newSQLiteVecStore(t, path, 3)
	version, err := store.VecVersion()
	if err != nil || version == "" {
		t.Fatalf("sqlite-vec extension is unavailable: version=%q err=%v", version, err)
	}
	// 测试向量来自真实课程文件切片，courseVector 是可复现的教学用向量函数。
	chunks := requireCourseChunks(t, ChunkPolicy{MaxRunes: 500, OverlapRunes: 50})
	// fixtures 同时保存切片和对应向量，方便按相同顺序完成写入与断言。
	fixtures := []struct {
		chunk  Chunk
		vector []float32
	}{
		{chunks[0], courseVector(chunks[0].Text, 3)},
		{chunks[1], courseVector(chunks[1].Text, 3)},
	}
	// 写入一个其他租户的真实切片，验证 SQL 预过滤不会混入无权资料。
	foreign := chunks[len(chunks)-1]
	foreign.TenantID = "other-course"
	fixtures = append(fixtures, struct {
		chunk  Chunk
		vector []float32
	}{foreign, courseVector(foreign.Text, 3)})
	// Upsert 使用稳定的 chunk ID：重复写同一 ID 应更新记录而不是制造重复向量。
	for _, fixture := range fixtures {
		if err := store.Upsert(fixture.chunk, fixture.vector); err != nil {
			t.Fatalf("upsert %s: %v", fixture.chunk.ID, err)
		}
	}
	// 维度错误必须在落库前失败，避免后续查询的距离计算不可解释。
	if err := store.Upsert(fixtures[0].chunk, courseVector(fixtures[0].chunk.Text, 2)); err == nil {
		t.Fatal("wrong vector dimension must fail before writing")
	}

	// 用第一个切片自身的向量查询；它应是当前租户中的 Top-1。
	results, err := store.Search(fixtures[0].vector, "course", 2)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(results) == 0 || results[0].ID != fixtures[0].chunk.ID {
		t.Fatalf("unexpected top-k order: %#v", results)
	}
	// 向量相近不等于有访问权限，结果集仍必须执行租户隔离。
	for _, result := range results {
		if result.TenantID != "course" {
			t.Fatalf("partition filter leaked another tenant: %#v", result)
		}
	}
	// 先关闭连接再重新打开，确保下面验证的是磁盘持久化而非进程内缓存。
	if err := store.Close(); err != nil {
		t.Fatalf("close first store: %v", err)
	}

	// 重开同一路径的数据库后，已写入的切片仍应能被同一查询向量找回。
	reopened := newSQLiteVecStore(t, path, 3)
	defer reopened.Close()
	reopenedResults, err := reopened.Search(fixtures[0].vector, "course", 1)
	if err != nil || len(reopenedResults) != 1 || reopenedResults[0].ID != fixtures[0].chunk.ID {
		t.Fatalf("persisted vector was not searchable after reopen: %#v, %v", reopenedResults, err)
	}
}
