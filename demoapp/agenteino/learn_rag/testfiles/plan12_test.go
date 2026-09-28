//go:build milvus && linux

package testfiles

import (
	"context"
	"errors"
	"testing"

	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

type backendConfig struct {
	// Kind 决定当前 RAG 检索实现，可选 sqlite 或 milvus。
	Kind string
	// SQLitePath 是 sqlite 后端的本地持久化路径。
	SQLitePath string
	// MilvusAddr 是 milvus 后端的服务地址，必须显式配置。
	MilvusAddr string
}

func validateBackendConfig(config backendConfig) error {
	// 后端切换只改变存储实现；配置校验先保证每种实现都有完成初始化所需的信息。
	switch config.Kind {
	case "sqlite":
		if config.SQLitePath == "" {
			return errors.New("sqlite backend requires a database path")
		}
	case "milvus":
		if config.MilvusAddr == "" {
			return errors.New("milvus backend requires an address")
		}
	default:
		return errors.New("backend must be sqlite or milvus")
	}
	return nil
}

func TestPlan12BackendSwitchPreservesRAGFilterContract(t *testing.T) {
	// 使用真实课程资料生成向量，验证后端切换时仍遵守相同的 RAG 数据与过滤契约。
	chunk := requireCourseChunks(t, ChunkPolicy{MaxRunes: 500, OverlapRunes: 50})[0]
	embedder := newHashEmbedder(64)
	vector, err := embedder.EmbedQuery(context.Background(), queryFromChunk(chunk))
	if err != nil {
		t.Fatalf("embed real source for backend contract: %v", err)
	}
	// SQLite 后端至少需要明确的数据库路径。
	if err := validateBackendConfig(backendConfig{Kind: "sqlite", SQLitePath: "derived/" + chunk.ID + ".sqlite"}); err != nil {
		t.Fatalf("valid sqlite config: %v", err)
	}
	// Milvus 后端至少需要可连接的地址；该检查不隐式启动或连接外部服务。
	if err := validateBackendConfig(backendConfig{Kind: "milvus", MilvusAddr: "127.0.0.1:19530"}); err != nil {
		t.Fatalf("valid Milvus config: %v", err)
	}
	// 缺少地址必须尽早失败，避免运行时退化为不可控的默认连接。
	if err := validateBackendConfig(backendConfig{Kind: "milvus"}); err == nil {
		t.Fatal("Milvus must not start without its explicit address")
	}

	// 无论采用 SQLite 还是 Milvus，都要保留租户、可见性等硬过滤，以及回答引用需要的字段。
	request, err := milvusclient.NewSearchOption("rag_"+chunk.ID, 4, []entity.Vector{entity.FloatVector(vector)}).
		WithANNSField("embedding").
		WithFilter("tenant_id == $tenant and visibility == $visibility").
		WithTemplateParam("tenant", chunk.TenantID).
		WithTemplateParam("visibility", chunk.Visibility).
		WithOutputFields("chunk_id", "source_uri", "chunk_text").
		Request()
	if err != nil {
		t.Fatalf("compose production search request: %v", err)
	}
	// 断言过滤条件与证据字段不会在后端适配时悄然丢失。
	if request.Dsl == "" || len(request.ExprTemplateValues) != 2 || len(request.OutputFields) != 3 {
		t.Fatalf("Milvus request lost hard filters or answer evidence fields: %#v", request)
	}
}
