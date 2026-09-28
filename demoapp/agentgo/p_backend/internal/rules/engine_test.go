package rules

import (
	"testing"

	"agentgo/p_backend/internal/domain"
)

func TestDefaultEngineCoversInitialTenRules(t *testing.T) {
	snapshot := domain.SchemaSnapshot{Dialect: "mysql", DefaultCollation: "utf8mb4_0900_ai_ci"}
	table := domain.Table{
		Name: "Orders-Log", Collation: "utf8mb4_general_ci",
		Columns: []domain.Column{
			{Name: "amount", DataType: "float", ColumnType: "float", Collation: "utf8mb4_general_ci"},
			{Name: "created_at", DataType: "timestamp", ColumnType: "timestamp", Collation: "utf8mb4_general_ci"},
			{Name: "note", DataType: "varchar", ColumnType: "varchar(32)", Nullable: true, DefaultIsNull: true, Collation: "latin1_swedish_ci"},
			{Name: "Bad-Column", DataType: "varchar", ColumnType: "varchar(32)", Collation: "utf8mb4_general_ci"},
		},
		Indexes: []domain.Index{
			{Name: "idx_a", Columns: []string{"code"}},
			{Name: "idx_a_dup", Columns: []string{"code"}},
			{Name: "idx_short", NonUnique: true, Columns: []string{"a", "b"}},
			{Name: "idx_wide", NonUnique: true, Columns: []string{"a", "b", "c", "d", "e"}},
		},
		ForeignKeys: []domain.ForeignKey{{Column: "user_id", ReferencedTable: "users", ReferencedColumn: "id"}},
	}
	engine := NewDefault()
	seen := map[string]bool{}
	for _, workerID := range []string{"naming-convention", "index-performance", "type-constraint"} {
		for _, finding := range engine.Inspect(workerID, snapshot, table) {
			seen[finding.RuleID] = true
		}
	}
	for _, id := range []string{"NAME-001", "NAME-002", "NAME-003", "INDEX-001", "INDEX-002", "INDEX-003", "INDEX-004", "INDEX-005", "TYPE-001", "TYPE-002", "TYPE-003", "TYPE-004"} {
		if !seen[id] {
			t.Errorf("expected finding for %s", id)
		}
	}
}
