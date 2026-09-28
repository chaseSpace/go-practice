package testfiles

import (
	"context"
	"strings"
	"testing"
)

func TestPlan6EvidenceBudgetAndGroundedAnswer(t *testing.T) {
	// 从真实课程资料建索引，后续证据、引用和回答都必须能回溯到这些切片。
	chunks := requireCourseChunks(t, ChunkPolicy{MaxRunes: 500, OverlapRunes: 50})
	retriever, err := newInMemoryRetriever(context.Background(), newHashEmbedder(64), chunks)
	if err != nil {
		t.Fatalf("index real course chunks: %v", err)
	}
	// 以真实切片生成查询，确保检索结果可稳定作为受约束回答的依据。
	query := queryFromChunk(chunks[0])
	results, err := retriever.Retrieve(context.Background(), SearchRequest{Query: query, TenantID: "course", TopK: 2, MinScore: 0.9999})
	if err != nil {
		t.Fatalf("retrieve real evidence: %v", err)
	}
	// 2,000 rune 预算模拟模型上下文窗口：证据按分数装入，超预算时停止而非无限拼接。
	contextText, evidence, err := buildEvidenceContext(results, 2000)
	if err != nil {
		t.Fatalf("build evidence: %v", err)
	}
	// 证据上下文必须同时保留引用标签和源 URI，才能让模型输出可核验的引用。
	if len(evidence) == 0 || !strings.Contains(contextText, "[S1]") || !strings.Contains(contextText, chunks[0].URI) {
		t.Fatalf("evidence context lacks source labels: %q", contextText)
	}
	// 演示回答器只依据传入的 Evidence 生成文字和 [S1]，不允许凭训练知识补全事实。
	answer := answerFromEvidence(query, evidence)
	if len(answer.Citations) != 1 || answer.Citations[0] != "S1" || !strings.Contains(answer.Text, "[S1]") {
		t.Fatalf("answer is not bound to supplied evidence: %#v", answer)
	}
	// 没有任何证据时，RAG 应明确拒答（资料不足），而不是编造一个貌似合理的答案。
	noEvidence := answerFromEvidence("unknown", nil)
	if len(noEvidence.Citations) != 0 || !strings.Contains(noEvidence.Text, "资料不足") {
		t.Fatalf("missing-evidence path must refuse instead of inventing: %#v", noEvidence)
	}
}
