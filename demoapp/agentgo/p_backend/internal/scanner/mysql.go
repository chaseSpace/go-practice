package scanner

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	"agentgo/p_backend/internal/domain"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/go-sql-driver/mysql"
)

type MySQLCollector struct {
	ConnectTimeout time.Duration
}

type SchemaCollectorTool struct{ Collector Collector }

type schemaToolInput struct {
	Connection domain.ConnectionConfig `json:"connection"`
	Excluded   []string                `json:"excludedTables"`
}

var _ tool.InvokableTool = SchemaCollectorTool{}

func (t SchemaCollectorTool) Info(context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "fetch_db_schema", Desc: "Read MySQL schema metadata without reading business rows.", ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
		"connection": {Type: schema.Object, Desc: "MySQL connection parameters", Required: true},
	})}, nil
}

func (t SchemaCollectorTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	var input schemaToolInput
	if err := json.Unmarshal([]byte(argumentsInJSON), &input); err != nil {
		return "", err
	}
	snapshot, err := t.Collect(ctx, input.Connection, input.Excluded)
	if err != nil {
		return "", err
	}
	output, err := json.Marshal(snapshot)
	return string(output), err
}

func (t SchemaCollectorTool) Collect(ctx context.Context, connection domain.ConnectionConfig, excludedTables []string) (domain.SchemaSnapshot, error) {
	excluded := make(map[string]struct{}, len(excludedTables))
	for _, table := range excludedTables {
		excluded[table] = struct{}{}
	}
	return t.Collector.Collect(ctx, connection, excluded)
}

func (c MySQLCollector) Collect(ctx context.Context, cfg domain.ConnectionConfig, excluded map[string]struct{}) (domain.SchemaSnapshot, error) {
	if cfg.Dialect != "mysql" {
		return domain.SchemaSnapshot{}, fmt.Errorf("unsupported database dialect")
	}
	dsn := mysql.NewConfig()
	dsn.User = cfg.Username
	dsn.Passwd = cfg.Password
	dsn.Net = "tcp"
	dsn.Addr = net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	dsn.DBName = cfg.Database
	dsn.ParseTime = true
	dsn.Loc = time.UTC
	dsn.Timeout = c.ConnectTimeout
	dsn.ReadTimeout = c.ConnectTimeout
	db, err := sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		return domain.SchemaSnapshot{}, fmt.Errorf("open mysql connection: %w", err)
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		return domain.SchemaSnapshot{}, fmt.Errorf("connect to database: %w", err)
	}

	snapshot := domain.SchemaSnapshot{Dialect: "mysql", Database: cfg.Database, CreatedAt: time.Now().UTC()}
	_ = db.QueryRowContext(ctx, `SELECT COALESCE(DEFAULT_COLLATION_NAME, '') FROM information_schema.schemata WHERE schema_name = ?`, cfg.Database).Scan(&snapshot.DefaultCollation)
	rows, err := db.QueryContext(ctx, `SELECT TABLE_NAME, COALESCE(TABLE_COMMENT, ''), COALESCE(TABLE_ROWS, 0), COALESCE(TABLE_COLLATION, '')
FROM information_schema.tables WHERE table_schema = ? AND table_type = 'BASE TABLE' ORDER BY table_name`, cfg.Database)
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
		table.Columns, err = loadColumns(ctx, db, cfg.Database, table.Name)
		if err != nil {
			snapshot.UncheckedObjects = append(snapshot.UncheckedObjects, table.Name)
			continue
		}
		table.Indexes, err = loadIndexes(ctx, db, cfg.Database, table.Name)
		if err != nil {
			snapshot.UncheckedObjects = append(snapshot.UncheckedObjects, table.Name)
			continue
		}
		table.ForeignKeys, err = loadForeignKeys(ctx, db, cfg.Database, table.Name)
		if err != nil {
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

func loadColumns(ctx context.Context, db *sql.DB, database, table string) ([]domain.Column, error) {
	rows, err := db.QueryContext(ctx, `SELECT COLUMN_NAME, DATA_TYPE, COLUMN_TYPE, IS_NULLABLE, COLUMN_DEFAULT, COALESCE(COLLATION_NAME, ''), COALESCE(COLUMN_COMMENT, '')
FROM information_schema.columns WHERE table_schema = ? AND table_name = ? ORDER BY ORDINAL_POSITION`, database, table)
	if err != nil {
		return nil, fmt.Errorf("collect columns for %s: %w", table, err)
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

func loadForeignKeys(ctx context.Context, db *sql.DB, database, table string) ([]domain.ForeignKey, error) {
	rows, err := db.QueryContext(ctx, `SELECT COLUMN_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME
FROM information_schema.key_column_usage
WHERE table_schema = ? AND table_name = ? AND referenced_table_name IS NOT NULL
ORDER BY constraint_name, ordinal_position`, database, table)
	if err != nil {
		return nil, fmt.Errorf("collect foreign keys for %s: %w", table, err)
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

func loadIndexes(ctx context.Context, db *sql.DB, database, table string) ([]domain.Index, error) {
	rows, err := db.QueryContext(ctx, `SELECT INDEX_NAME, NON_UNIQUE, SEQ_IN_INDEX, COLUMN_NAME
FROM information_schema.statistics WHERE table_schema = ? AND table_name = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX`, database, table)
	if err != nil {
		return nil, fmt.Errorf("collect indexes for %s: %w", table, err)
	}
	defer rows.Close()
	byName := map[string]*domain.Index{}
	for rows.Next() {
		var name, column string
		var nonUnique, sequence int
		if err := rows.Scan(&name, &nonUnique, &sequence, &column); err != nil {
			return nil, err
		}
		index := byName[name]
		if index == nil {
			index = &domain.Index{Name: name, NonUnique: nonUnique == 1, Primary: name == "PRIMARY"}
			byName[name] = index
		}
		index.Columns = append(index.Columns, column)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	indexes := make([]domain.Index, 0, len(byName))
	for _, index := range byName {
		indexes = append(indexes, *index)
	}
	sort.Slice(indexes, func(i, j int) bool { return indexes[i].Name < indexes[j].Name })
	return indexes, nil
}
