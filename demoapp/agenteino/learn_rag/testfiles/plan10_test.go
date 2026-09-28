package testfiles

import (
	"context"
	"testing"
)

func TestPlan10EvaluationSeparatesRecallFromUnanswerableQueries(t *testing.T) {
	// 用真实语料切片建立待评测索引；评估不应只在手写玩具文本上获得高分。
	chunks := requireCourseChunks(t, ChunkPolicy{MaxRunes: 500, OverlapRunes: 50})
	retriever, err := newInMemoryRetriever(context.Background(), newHashEmbedder(64), chunks)
	if err != nil {
		t.Fatalf("create retriever: %v", err)
	}
	// 评测集既包括两个“应该找回”的真实切片，也包括一个权限拒绝的不可回答问题。
	report, err := evaluateRetriever(context.Background(), retriever, []EvalCase{
		{
			// 第一条正例检查第一个源切片是否出现在 Top-K 内。
			Name:             chunks[0].URI,
			Request:          SearchRequest{Query: queryFromChunk(chunks[0]), TenantID: "course", TopK: 2, MinScore: 0.9999},
			ExpectedChunkIDs: []string{chunks[0].ID},
			Answerable:       true,
		},
		{
			// 第二条正例覆盖另一份真实资料，避免指标只反映单一文档。
			Name:             chunks[len(chunks)-1].URI,
			Request:          SearchRequest{Query: queryFromChunk(chunks[len(chunks)-1]), TenantID: "course", TopK: 2, MinScore: 0.9999},
			ExpectedChunkIDs: []string{chunks[len(chunks)-1].ID},
			Answerable:       true,
		},
		{
			// 文本本身存在，但租户无权访问，因此应计入“正确拒答”而非召回失败。
			Name:       "same real source under denied tenant",
			Request:    SearchRequest{Query: queryFromChunk(chunks[0]), TenantID: "denied-tenant", TopK: 2, MinScore: 0.9999},
			Answerable: false,
		},
	})
	if err != nil {
		t.Fatalf("evaluate retriever: %v", err)
	}
	// Recall@K 衡量是否找回正确证据，MRR 衡量其排序，UnanswerablePassed 衡量拒答边界。
	if report.RecallAtK != 1 || report.MRR != 1 || report.UnanswerablePassed != 1 {
		t.Fatalf("unexpected retrieval report: %#v", report)
	}
}
