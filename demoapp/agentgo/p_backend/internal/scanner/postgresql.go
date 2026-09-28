package scanner

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"sort"
	"strings"
	"time"

	"agentgo/p_backend/internal/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgreSQLCollector struct{ ConnectTimeout time.Duration }

func (c PostgreSQLCollector) Collect(ctx context.Context, cfg domain.ConnectionConfig, excluded map[string]struct{}) (domain.SchemaSnapshot, error) {
	if cfg.Dialect != "postgres" && cfg.Dialect != "postgresql" {
		return domain.SchemaSnapshot{}, fmt.Errorf("unsupported database dialect")
	}
	dsn := (&url.URL{Scheme: "postgres", User: url.UserPassword(cfg.Username, cfg.Password), Host: net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port)), Path: cfg.Database, RawQuery: "connect_timeout=" + fmt.Sprintf("%d", max(1, int(c.ConnectTimeout.Seconds())))}).String()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return domain.SchemaSnapshot{}, fmt.Errorf("open postgresql connection: %w", err)
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		return domain.SchemaSnapshot{}, fmt.Errorf("connect to database: %w", err)
	}

	snapshot := domain.SchemaSnapshot{Dialect: "postgres", Database: cfg.Database, CreatedAt: time.Now().UTC()}
	_ = db.QueryRowContext(ctx, `SELECT datcollate FROM pg_database WHERE datname = current_database()`).Scan(&snapshot.DefaultCollation)
	rows, err := db.QueryContext(ctx, `SELECT c.relname, COALESCE(obj_description(c.oid), ''), COALESCE(c.reltuples, 0)::bigint, COALESCE(coll.collname, '')
FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
LEFT JOIN pg_collation coll ON coll.oid = c.relcollation
WHERE n.nspname = current_schema() AND c.relkind IN ('r', 'p') ORDER BY c.relname`)
	if err != nil {
		return domain.SchemaSnapshot{}, fmt.Errorf("collect tables: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var table domain.Table
		if err := rows.Scan(&table.Name, &table.Comment, &table.EstimatedRows, &table.Collation); err != nil {
			return domain.SchemaSnapshot{}, err
		}
		if _, skip := excluded[table.Name]; skip {
			continue
		}
		if table.Columns, err = loadPostgresColumns(ctx, db, table.Name); err != nil {
			snapshot.UncheckedObjects = append(snapshot.UncheckedObjects, table.Name)
			continue
		}
		if table.Indexes, err = loadPostgresIndexes(ctx, db, table.Name); err != nil {
			snapshot.UncheckedObjects = append(snapshot.UncheckedObjects, table.Name)
			continue
		}
		if table.ForeignKeys, err = loadPostgresForeignKeys(ctx, db, table.Name); err != nil {
			snapshot.UncheckedObjects = append(snapshot.UncheckedObjects, table.Name)
			continue
		}
		snapshot.Tables = append(snapshot.Tables, table)
	}
	if err := rows.Err(); err != nil {
		return domain.SchemaSnapshot{}, err
	}
	return snapshot, nil
}

func loadPostgresColumns(ctx context.Context, db *sql.DB, table string) ([]domain.Column, error) {
	rows, err := db.QueryContext(ctx, `SELECT c.column_name, c.data_type, c.udt_name, c.is_nullable, c.column_default, COALESCE(c.collation_name, ''),
COALESCE(pg_catalog.col_description(pc.oid, c.ordinal_position), '')
FROM information_schema.columns c
JOIN pg_class pc ON pc.relname = c.table_name
JOIN pg_namespace pn ON pn.oid = pc.relnamespace AND pn.nspname = c.table_schema
WHERE c.table_schema = current_schema() AND c.table_name = ? ORDER BY c.ordinal_position`, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns := make([]domain.Column, 0)
	for rows.Next() {
		var column domain.Column
		var nullable string
		var defaultValue sql.NullString
		if err := rows.Scan(&column.Name, &column.DataType, &column.ColumnType, &nullable, &defaultValue, &column.Collation, &column.Comment); err != nil {
			return nil, err
		}
		column.Nullable = strings.EqualFold(nullable, "YES")
		column.DefaultIsNull = !defaultValue.Valid
		if defaultValue.Valid {
			column.DefaultValue = defaultValue.String
		}
		columns = append(columns, column)
	}
	return columns, rows.Err()
}

func loadPostgresIndexes(ctx context.Context, db *sql.DB, table string) ([]domain.Index, error) {
	rows, err := db.QueryContext(ctx, `SELECT idx.relname, NOT i.indisunique, i.indisprimary,
string_agg(att.attname, ',' ORDER BY ord.ordinality)
FROM pg_index i
JOIN pg_class tbl ON tbl.oid = i.indrelid
JOIN pg_namespace n ON n.oid = tbl.relnamespace
JOIN pg_class idx ON idx.oid = i.indexrelid
JOIN unnest(i.indkey) WITH ORDINALITY AS ord(attnum, ordinality) ON TRUE
JOIN pg_attribute att ON att.attrelid = tbl.oid AND att.attnum = ord.attnum
WHERE n.nspname = current_schema() AND tbl.relname = ?
GROUP BY idx.relname, i.indisunique, i.indisprimary ORDER BY idx.relname`, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	indexes := make([]domain.Index, 0)
	for rows.Next() {
		var index domain.Index
		var columns string
		if err := rows.Scan(&index.Name, &index.NonUnique, &index.Primary, &columns); err != nil {
			return nil, err
		}
		index.Columns = strings.Split(columns, ",")
		indexes = append(indexes, index)
	}
	sort.Slice(indexes, func(i, j int) bool { return indexes[i].Name < indexes[j].Name })
	return indexes, rows.Err()
}

func loadPostgresForeignKeys(ctx context.Context, db *sql.DB, table string) ([]domain.ForeignKey, error) {
	rows, err := db.QueryContext(ctx, `SELECT a.attname, rt.relname, ra.attname
FROM pg_constraint con
JOIN pg_class t ON t.oid = con.conrelid
JOIN pg_namespace n ON n.oid = t.relnamespace
JOIN pg_class rt ON rt.oid = con.confrelid
JOIN unnest(con.conkey, con.confkey) AS keys(local_attnum, ref_attnum) ON TRUE
JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = keys.local_attnum
JOIN pg_attribute ra ON ra.attrelid = rt.oid AND ra.attnum = keys.ref_attnum
WHERE con.contype = 'f' AND n.nspname = current_schema() AND t.relname = ?`, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	keys := make([]domain.ForeignKey, 0)
	for rows.Next() {
		var key domain.ForeignKey
		if err := rows.Scan(&key.Column, &key.ReferencedTable, &key.ReferencedColumn); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}

func max(left, right int) int {
	if left > right {
		return left
	}
	return right
}
