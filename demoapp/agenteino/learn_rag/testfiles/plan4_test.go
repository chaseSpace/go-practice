package testfiles

import (
	"context"
	"testing"
)

func TestPlan4EmbeddingContractAndCosineRanking(t *testing.T) {
	// 使用真实 TXT/PDF/DOCX 语料生成 chunk，而不是手写演示文本。
	chunks := requireCourseChunks(t, ChunkPolicy{MaxRunes: 500, OverlapRunes: 50})
	// 64 是课程离线 HashEmbedder 的固定向量维度；索引和查询必须保持一致。
	embedder := newHashEmbedder(64)

	// 以第一个真实 chunk 的完整文本作为 query，构造可重复的“自身命中”基线。
	query, err := embedder.EmbedQuery(context.Background(), queryFromChunk(chunks[0]))
	if err != nil {
		t.Fatalf("embed query: %v", err)
	}

	// 同一段正文向量应与 query 完全一致，因此余弦相似度应为 1。
	near, err := embedder.EmbedQuery(context.Background(), chunks[0].Text)
	if err != nil {
		t.Fatalf("embed real source chunk: %v", err)
	}
	// 另一段真实正文作为对照；它不是人为构造的“无关文本”。
	far, err := embedder.EmbedQuery(context.Background(), chunks[len(chunks)-1].Text)
	if err != nil {
		t.Fatalf("embed another real source chunk: %v", err)
	}
	// 余弦相似度比较要求两个向量维度相同且都不是零向量。
	nearScore, err := cosineSimilarity(query, near)
	if err != nil {
		t.Fatalf("score near vector: %v", err)
	}
	farScore, err := cosineSimilarity(query, far)
	if err != nil {
		t.Fatalf("score far vector: %v", err)
	}
	// 自身命中必须排在另一真实 chunk 之前或至少并列，证明排序方向正确。
	if nearScore != 1 || nearScore < farScore {
		t.Fatalf("real source query must score itself highest: self=%v other=%v", nearScore, farScore)
	}
	// 维度不一致不能静默计算；真实向量库同样会拒绝这种写入或查询。
	if _, err := cosineSimilarity(query, []float32{1, 0}); err == nil {
		t.Fatal("mixed vector dimensions must be rejected")
	}
	// 模型标识和维度是索引 manifest 的必要契约，模型切换时据此触发重建。
	if embedder.Dimension() != len(query) || embedder.ModelID() == "" {
		t.Fatalf("embedding contract is incomplete: dim=%d vector=%d model=%q", embedder.Dimension(), len(query), embedder.ModelID())
	}
}
