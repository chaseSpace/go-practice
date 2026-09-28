package testfiles

import (
	"context"
	"testing"
)

func TestPlan5InMemoryRetrieverRanksAndFilters(t *testing.T) {
	// 使用课程目录中的真实 TXT、PDF、DOCX 抽取结果构建索引，避免示例脱离实际数据。
	chunks := requireCourseChunks(t, ChunkPolicy{MaxRunes: 500, OverlapRunes: 50})
	// 为每个切片补充可过滤的业务元数据；向量相似度不能替代权限和状态约束。
	for index := range chunks {
		chunks[index].Metadata["status"] = "current"
	}
	// 将一条真实切片复制到其他租户，用于验证检索端不会跨租户泄漏资料。
	foreign := chunks[len(chunks)-1]
	foreign.TenantID = "other-course"
	chunks = append(chunks[:len(chunks)-1], foreign)
	// 以一个真实切片的首句作为查询，预期它应以最高分命中自身。
	queryChunk := chunks[0]
	// 64 维哈希向量仅服务于可重复的课程测试；生产应替换为真实 Embedding 模型。
	retriever, err := newInMemoryRetriever(context.Background(), newHashEmbedder(64), chunks)
	if err != nil {
		t.Fatalf("create retriever: %v", err)
	}
	// TopK、最小分数、租户和元数据过滤共同组成检索契约，缺一不可。
	results, err := retriever.Retrieve(context.Background(), SearchRequest{
		Query:    queryFromChunk(queryChunk),
		TenantID: "course",
		TopK:     2,
		MinScore: 0.9999,
		Filters:  map[string]string{"status": "current"},
	})
	if err != nil {
		t.Fatalf("retrieve: %v", err)
	}
	// 自查询应排在 Top-1，证明余弦排序和 Top-K 截断按预期协同工作。
	if len(results) == 0 || results[0].ID != queryChunk.ID {
		t.Fatalf("unexpected top-k result: %#v", results)
	}
	// 即使其他租户的文本相似，也绝不能进入当前用户的证据集合。
	for _, result := range results {
		if result.TenantID != "course" {
			t.Fatalf("cross-tenant result leaked: %#v", result)
		}
	}
	// 非法 TopK 在请求入口失败，避免静默返回“看似正常”的空结果。
	if _, err := retriever.Retrieve(context.Background(), SearchRequest{Query: queryFromChunk(queryChunk), TopK: 0}); err == nil {
		t.Fatal("invalid top-k must fail deterministically")
	}
}
