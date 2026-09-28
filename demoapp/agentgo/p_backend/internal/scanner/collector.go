package scanner

import (
	"context"
	"fmt"

	"agentgo/p_backend/internal/domain"
)

type Collector interface {
	Collect(context.Context, domain.ConnectionConfig, map[string]struct{}) (domain.SchemaSnapshot, error)
}

type MultiCollector struct {
	MySQL      MySQLCollector
	PostgreSQL PostgreSQLCollector
}

func (c MultiCollector) Collect(ctx context.Context, config domain.ConnectionConfig, excluded map[string]struct{}) (domain.SchemaSnapshot, error) {
	switch config.Dialect {
	case "mysql":
		return c.MySQL.Collect(ctx, config, excluded)
	case "postgres", "postgresql":
		return c.PostgreSQL.Collect(ctx, config, excluded)
	default:
		return domain.SchemaSnapshot{}, fmt.Errorf("unsupported database dialect")
	}
}
