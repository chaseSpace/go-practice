//go:build milvus && linux

package testfiles

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

func milvusRAGSchema(dim int64) *entity.Schema {
	// schema 同时保存切片标识、可过滤字段、可引用的源地址和原文；只存向量无法完成 RAG 回答。
	return entity.NewSchema().WithDynamicFieldEnabled(false).
		// chunk_id 是幂等写入与证据引用的稳定主键。
		WithField(entity.NewField().WithName("chunk_id").WithDataType(entity.FieldTypeVarChar).WithMaxLength(128).WithIsPrimaryKey(true)).
		// document_id 支持按源文档整体删除或重建索引。
		WithField(entity.NewField().WithName("document_id").WithDataType(entity.FieldTypeVarChar).WithMaxLength(128)).
		// tenant_id 必须参与检索过滤，避免跨租户的相似文本泄漏。
		WithField(entity.NewField().WithName("tenant_id").WithDataType(entity.FieldTypeVarChar).WithMaxLength(64)).
		// source_uri 和 chunk_text 是生成答案和展示引用所需的原始证据。
		WithField(entity.NewField().WithName("source_uri").WithDataType(entity.FieldTypeVarChar).WithMaxLength(1024)).
		WithField(entity.NewField().WithName("chunk_text").WithDataType(entity.FieldTypeVarChar).WithMaxLength(8192)).
		// embedding 维度必须等于当前嵌入模型的输出维度，否则距离无定义。
		WithField(entity.NewField().WithName("embedding").WithDataType(entity.FieldTypeFloatVector).WithDim(dim))
}

func TestPlan11MilvusSchemaAndSearchOption(t *testing.T) {
	// 从真实课程切片产生查询向量，令 Milvus 的建模示例与课程语料保持一致。
	chunk := requireCourseChunks(t, ChunkPolicy{MaxRunes: 500, OverlapRunes: 50})[0]
	embedder := newHashEmbedder(64)
	vector, err := embedder.EmbedQuery(context.Background(), queryFromChunk(chunk))
	if err != nil {
		t.Fatalf("embed real source for Milvus: %v", err)
	}
	// 使用切片 ID 派生唯一集合名，防止集成测试与已有业务集合互相影响。
	collection := "rag_" + chunk.ID
	// COSINE 索引与课程内存检索器使用同一相似度语义，便于替换后比较效果。
	schema := milvusRAGSchema(int64(embedder.Dimension()))
	option := milvusclient.NewCreateCollectionOption(collection, schema).WithIndexOptions(
		milvusclient.NewCreateIndexOption(collection, "embedding", index.NewAutoIndex(entity.COSINE)),
	)
	request := option.Request()
	if request.CollectionName != collection || len(request.Schema) == 0 || len(option.Indexes()) != 1 {
		t.Fatalf("collection request is incomplete: %#v", request)
	}

	// 搜索请求将租户作为模板参数，而不是拼接到过滤字符串中；同时取回回答所需证据字段。
	search, err := milvusclient.NewSearchOption(collection, 3, []entity.Vector{entity.FloatVector(vector)}).
		WithANNSField("embedding").
		WithFilter("tenant_id == $tenant").
		WithTemplateParam("tenant", chunk.TenantID).
		WithOutputFields("chunk_id", "document_id", "source_uri", "chunk_text").
		Request()
	if err != nil {
		t.Fatalf("build Milvus search request: %v", err)
	}
	// 断言请求没有丢失集合名、硬过滤或引用字段，避免仅测试“请求能创建”。
	if search.CollectionName != collection || search.Dsl != "tenant_id == $tenant" || len(search.OutputFields) != 4 {
		t.Fatalf("search request lost RAG fields or filter: %#v", search)
	}
}

// 此测试刻意采用显式开关：它只验证配置的 Milvus 服务是否可达。
// collection 的创建和修改应放在唯一命名、明确启用的集成测试中执行。
func TestPlan11MilvusServiceConnection(t *testing.T) {
	// 取真实切片仅用于生成一个无冲突的集合名；本测试不修改服务端数据。
	chunk := requireCourseChunks(t, ChunkPolicy{MaxRunes: 500, OverlapRunes: 50})[0]
	// 外部服务测试必须由环境变量显式开启，保证默认 go test 不依赖本地 Milvus。
	if os.Getenv("EINO_RUN_MILVUS_TESTS") != "1" {
		t.Skip("set EINO_RUN_MILVUS_TESTS=1 and MILVUS_ADDRESS to connect a Milvus service")
	}
	// 可通过 MILVUS_ADDRESS 指向 CI 或 Docker 中的服务；未设置时使用常见本地端口。
	address := os.Getenv("MILVUS_ADDRESS")
	if address == "" {
		address = "127.0.0.1:19530"
	}
	// 连接设置超时，防止服务不可用时测试无限阻塞。
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := milvusclient.New(ctx, &milvusclient.ClientConfig{Address: address})
	if err != nil {
		t.Fatalf("connect Milvus at %s: %v", address, err)
	}
	defer client.Close(context.Background())
	// HasCollection 是只读健康请求，足以验证客户端、网络和认证配置是否可用。
	if _, err := client.HasCollection(ctx, milvusclient.NewHasCollectionOption("rag_"+chunk.ID)); err != nil {
		t.Fatalf("Milvus health request: %v", err)
	}
}
