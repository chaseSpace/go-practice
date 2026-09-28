package testfiles

import (
	"context"
	"strings"
	"testing"
)

func TestPlan7LocalRAGAnswersOnlyFromIndexedSources(t *testing.T) {
	// 同一 Context 会贯穿切片、建索引和问答，生产中可由它统一控制超时和取消。
	ctx := context.Background()
	// 课程要求同时加载 data/ 下的真实 TXT、PDF、DOCX，而非临时构造文本。
	documents := requireCourseDocuments(t)
	// 300 rune 的切片和 50 rune 重叠用于展示端到端 RAG 的可调分块策略。
	policy := ChunkPolicy{MaxRunes: 300, OverlapRunes: 50}
	// 先从第一个真实文档得到一个查询片段，后续检查答案能否准确引用该部署资料。
	queryChunks, err := chunkDocument(documents[0], policy)
	if err != nil {
		t.Fatalf("chunk query source: %v", err)
	}
	query := queryFromChunk(queryChunks[0])
	// LocalRAG 封装“文档 → 分块 → 向量索引 → 检索 → 证据回答”的最小闭环。
	rag, err := newLocalRAG(ctx, documents, policy, newHashEmbedder(64))
	if err != nil {
		t.Fatalf("build local rag: %v", err)
	}
	// Ask 的结果需要包含带标签的答案及其证据，调用方可据此展示来源或继续审计。
	answer, err := rag.Ask(ctx, SearchRequest{Query: query, TenantID: "course", TopK: 2, MinScore: 0.9999}, 2000)
	if err != nil {
		t.Fatalf("ask rag: %v", err)
	}
	if !strings.Contains(answer.Text, "[S1]") || len(answer.Evidence) == 0 || answer.Evidence[0].Chunk.URI != documents[0].URI {
		t.Fatalf("answer did not use the deployment source: %#v", answer)
	}

	// 删除承载答案的源文档后重新建索引，验证索引更新确实影响可回答范围。
	withoutSource := documents[1:]
	rebuilt, err := newLocalRAG(ctx, withoutSource, policy, newHashEmbedder(64))
	if err != nil {
		t.Fatalf("rebuild local rag: %v", err)
	}
	answer, err = rebuilt.Ask(ctx, SearchRequest{Query: query, TenantID: "course", TopK: 2, MinScore: 0.9999}, 2000)
	if err != nil {
		t.Fatalf("ask rebuilt rag: %v", err)
	}
	// 已删除的内容不得被旧索引“幽灵命中”；没有证据时必须走拒答分支。
	if !strings.Contains(answer.Text, "资料不足") {
		t.Fatalf("removed source must not remain answerable: %#v", answer)
	}
}
