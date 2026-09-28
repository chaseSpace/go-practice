package rules

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
	"time"

	"agentgo/p_backend/internal/domain"
)

type Rule struct {
	ID       string
	WorkerID string
	Dialect  string
	Enabled  bool
	Evaluate func(domain.SchemaSnapshot, domain.Table) []domain.Finding
}

type Engine struct{ rules []Rule }

type Option func(*Engine)

// WithWorkerRuleSets restricts the built-in rules to the explicitly declared
// Worker rule sets loaded from config/rule_and_prompt.yaml.
func WithWorkerRuleSets(ruleSets map[string][]string) Option {
	return func(engine *Engine) {
		enabled := make(map[string]map[string]struct{}, len(ruleSets))
		for workerID, ruleIDs := range ruleSets {
			enabled[workerID] = make(map[string]struct{}, len(ruleIDs))
			for _, ruleID := range ruleIDs {
				enabled[workerID][ruleID] = struct{}{}
			}
		}
		for index := range engine.rules {
			_, isEnabled := enabled[engine.rules[index].WorkerID][engine.rules[index].ID]
			engine.rules[index].Enabled = isEnabled
		}
	}
}

func NewDefault(options ...Option) *Engine {
	engine := &Engine{rules: []Rule{
		{ID: "NAME-001", WorkerID: "naming-convention", Dialect: "mysql", Enabled: true, Evaluate: missingTableComment},
		{ID: "NAME-002", WorkerID: "naming-convention", Dialect: "mysql", Enabled: true, Evaluate: tableNameConvention},
		{ID: "NAME-003", WorkerID: "naming-convention", Dialect: "mysql", Enabled: true, Evaluate: columnNameConvention},
		{ID: "INDEX-001", WorkerID: "index-performance", Dialect: "mysql", Enabled: true, Evaluate: missingPrimaryKey},
		{ID: "INDEX-002", WorkerID: "index-performance", Dialect: "mysql", Enabled: true, Evaluate: duplicateIndexes},
		{ID: "INDEX-003", WorkerID: "index-performance", Dialect: "mysql", Enabled: true, Evaluate: foreignKeyWithoutIndex},
		{ID: "INDEX-004", WorkerID: "index-performance", Dialect: "mysql", Enabled: true, Evaluate: wideIndexes},
		{ID: "INDEX-005", WorkerID: "index-performance", Dialect: "mysql", Enabled: true, Evaluate: redundantPrefixIndexes},
		{ID: "TYPE-001", WorkerID: "type-constraint", Dialect: "mysql", Enabled: true, Evaluate: floatForAmount},
		{ID: "TYPE-002", WorkerID: "type-constraint", Dialect: "mysql", Enabled: true, Evaluate: undocumentedTimeColumns},
		{ID: "TYPE-003", WorkerID: "type-constraint", Dialect: "mysql", Enabled: true, Evaluate: nullableDefaultNull},
		{ID: "TYPE-004", WorkerID: "type-constraint", Dialect: "mysql", Enabled: true, Evaluate: inconsistentCollation},
	}}
	for _, option := range options {
		option(engine)
	}
	return engine
}

func (e *Engine) Inspect(workerID string, snapshot domain.SchemaSnapshot, table domain.Table) []domain.Finding {
	findings := make([]domain.Finding, 0)
	for _, rule := range e.rules {
		if !rule.Enabled || rule.WorkerID != workerID || rule.Dialect != snapshot.Dialect {
			continue
		}
		findings = append(findings, rule.Evaluate(snapshot, table)...)
	}
	return findings
}

func missingTableComment(_ domain.SchemaSnapshot, table domain.Table) []domain.Finding {
	if table.Comment != "" {
		return nil
	}
	return []domain.Finding{finding("NAME-001", domain.SeveritySuggestion, "table", table.Name, "表缺少注释", "表注释为空", "补充表的业务含义和维护说明", "", "", "", false, "low")}
}

var identifier = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

func tableNameConvention(_ domain.SchemaSnapshot, table domain.Table) []domain.Finding {
	if identifier.MatchString(table.Name) {
		return nil
	}
	return []domain.Finding{finding("NAME-002", domain.SeveritySuggestion, "table", table.Name, "表命名不规范", "表名未使用小写下划线风格", "在迁移计划中统一为小写下划线命名", "", "", "", true, "medium")}
}

func columnNameConvention(_ domain.SchemaSnapshot, table domain.Table) []domain.Finding {
	findings := make([]domain.Finding, 0)
	for _, column := range table.Columns {
		if identifier.MatchString(column.Name) {
			continue
		}
		findings = append(findings, finding("NAME-003", domain.SeveritySuggestion, "column", table.Name+"."+column.Name, "字段命名不规范", "字段名未使用小写下划线风格", "在兼容迁移计划中统一为小写下划线命名，并同步更新应用映射", "", "", "变更前确认 API、ORM 与报表字段依赖", true, "medium"))
	}
	return findings
}

func missingPrimaryKey(snapshot domain.SchemaSnapshot, table domain.Table) []domain.Finding {
	for _, index := range table.Indexes {
		if index.Primary {
			return nil
		}
	}
	tableName := quote(snapshot.Dialect, table.Name)
	candidate := quote(snapshot.Dialect, "candidate_unique_key")
	precheck := fmt.Sprintf("SELECT COUNT(*) AS row_count FROM %s;", tableName)
	dml := fmt.Sprintf("-- 将 <candidate_unique_key> 替换为已确认的唯一业务键后执行\nUPDATE %s SET %s = %s WHERE %s IS NULL;", tableName, candidate, uuidExpression(snapshot.Dialect), candidate)
	result := finding("INDEX-001", domain.SeverityUrgent, "table", table.Name, "表缺少主键", "未发现 PRIMARY KEY 索引", "确认唯一业务键，先使用 DML 补齐候选键，再在维护窗口新增主键", dml, precheck, "执行前备份表；DML 仅在人工替换候选键后分批执行", true, "high")
	result.DDLSuggestion = fmt.Sprintf("-- 人工将 <candidate_unique_key> 替换为已确认且已回填的唯一业务键\nALTER TABLE %s ADD PRIMARY KEY (%s);", tableName, candidate)
	return []domain.Finding{result}
}

func duplicateIndexes(_ domain.SchemaSnapshot, table domain.Table) []domain.Finding {
	seen := map[string]string{}
	findings := make([]domain.Finding, 0)
	for _, index := range table.Indexes {
		if index.Primary {
			continue
		}
		signature := strings.Join(index.Columns, ",")
		if previous, ok := seen[signature]; ok {
			result := finding("INDEX-002", domain.SeverityNormal, "index", table.Name+"."+index.Name, "存在重复索引", fmt.Sprintf("索引 %s 与 %s 的列顺序相同：%s", index.Name, previous, signature), "确认查询计划后删除冗余索引", "", fmt.Sprintf("SHOW INDEX FROM `%s`;", table.Name), "删除索引前保留建索引 DDL", true, "medium")
			result.DDLSuggestion = fmt.Sprintf("DROP INDEX %s ON %s;", quote("mysql", index.Name), quote("mysql", table.Name))
			findings = append(findings, result)
			continue
		}
		seen[signature] = index.Name
	}
	return findings
}

func foreignKeyWithoutIndex(_ domain.SchemaSnapshot, table domain.Table) []domain.Finding {
	findings := make([]domain.Finding, 0)
	for _, foreignKey := range table.ForeignKeys {
		if hasLeadingIndex(table.Indexes, foreignKey.Column) {
			continue
		}
		result := finding("INDEX-003", domain.SeverityNormal, "column", table.Name+"."+foreignKey.Column, "外键列缺少索引", fmt.Sprintf("外键关联 %s.%s，但该列不是任何索引的首列", foreignKey.ReferencedTable, foreignKey.ReferencedColumn), "为外键列创建合适的普通索引", "", fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE `%s` IS NOT NULL;", table.Name, foreignKey.Column), "保留创建索引 DDL 以便回滚", true, "medium")
		indexName := "idx_" + table.Name + "_" + foreignKey.Column
		result.DDLSuggestion = fmt.Sprintf("CREATE INDEX %s ON %s (%s);", quote("mysql", indexName), quote("mysql", table.Name), quote("mysql", foreignKey.Column))
		findings = append(findings, result)
	}
	return findings
}

func wideIndexes(_ domain.SchemaSnapshot, table domain.Table) []domain.Finding {
	findings := make([]domain.Finding, 0)
	for _, index := range table.Indexes {
		if len(index.Columns) < 5 {
			continue
		}
		result := finding("INDEX-004", domain.SeveritySuggestion, "index", table.Name+"."+index.Name, "联合索引列过多", fmt.Sprintf("索引包含 %d 列：%s", len(index.Columns), strings.Join(index.Columns, ", ")), "结合真实查询条件和选择性拆分或缩减索引", "", fmt.Sprintf("SHOW INDEX FROM `%s` WHERE Key_name = '%s';", table.Name, index.Name), "保留原索引 DDL", true, "medium")
		result.DDLSuggestion = fmt.Sprintf("-- 仅在人工确认保留列后执行\nDROP INDEX %s ON %s;\n-- CREATE INDEX <new_index_name> ON %s (<confirmed_columns>);", quote("mysql", index.Name), quote("mysql", table.Name), quote("mysql", table.Name))
		findings = append(findings, result)
	}
	return findings
}

func redundantPrefixIndexes(snapshot domain.SchemaSnapshot, table domain.Table) []domain.Finding {
	findings := make([]domain.Finding, 0)
	for _, candidate := range table.Indexes {
		if candidate.Primary || !candidate.NonUnique || len(candidate.Columns) == 0 {
			continue
		}
		for _, covering := range table.Indexes {
			if candidate.Name == covering.Name || len(covering.Columns) <= len(candidate.Columns) || !samePrefix(candidate.Columns, covering.Columns) {
				continue
			}
			result := finding("INDEX-005", domain.SeveritySuggestion, "index", table.Name+"."+candidate.Name, "索引可能被更长前缀覆盖", fmt.Sprintf("索引 %s(%s) 是 %s(%s) 的左前缀", candidate.Name, strings.Join(candidate.Columns, ", "), covering.Name, strings.Join(covering.Columns, ", ")), "核对唯一性、执行计划与写入开销后，评估是否保留较短索引", "", fmt.Sprintf("SHOW INDEX FROM `%s`;", table.Name), "删除前保留建索引 DDL，并在维护窗口验证回滚方案", true, "medium")
			result.DDLSuggestion = fmt.Sprintf("DROP INDEX %s ON %s;", quote("mysql", candidate.Name), quote("mysql", table.Name))
			findings = append(findings, result)
			break
		}
	}
	return findings
}

func floatForAmount(snapshot domain.SchemaSnapshot, table domain.Table) []domain.Finding {
	findings := make([]domain.Finding, 0)
	for _, column := range table.Columns {
		typeName := strings.ToLower(column.DataType)
		if typeName != "float" && typeName != "double" && typeName != "real" {
			continue
		}
		object := table.Name + "." + column.Name
		tableName, columnName := quote(snapshot.Dialect, table.Name), quote(snapshot.Dialect, column.Name)
		dml := fmt.Sprintf("UPDATE %s SET %s = ROUND(%s, 2) WHERE %s <> ROUND(%s, 2);", tableName, columnName, columnName, columnName, columnName)
		precheck := fmt.Sprintf("SELECT COUNT(*) AS affected_rows FROM %s WHERE %s <> ROUND(%s, 2);", tableName, columnName, columnName)
		findings = append(findings, finding("TYPE-001", domain.SeverityUrgent, "column", object, "金额字段使用浮点类型", "字段类型为 "+column.ColumnType+"，浮点数不能精确表示十进制金额", "先核对受影响记录，人工审核 DML 后迁移为 DECIMAL", dml, precheck, "执行前导出受影响行；DML 在事务中小批执行", true, "high"))
	}
	return findings
}

func undocumentedTimeColumns(_ domain.SchemaSnapshot, table domain.Table) []domain.Finding {
	findings := make([]domain.Finding, 0)
	for _, column := range table.Columns {
		typeName := strings.ToLower(column.DataType)
		if (typeName == "datetime" || typeName == "timestamp" || typeName == "date") && column.Comment == "" {
			findings = append(findings, finding("TYPE-002", domain.SeveritySuggestion, "column", table.Name+"."+column.Name, "时间字段语义不清", "时间字段缺少注释，无法区分创建、更新、业务生效等语义", "补充字段语义、时区和写入规则说明", "", "", "", false, "low"))
		}
	}
	return findings
}

func nullableDefaultNull(_ domain.SchemaSnapshot, table domain.Table) []domain.Finding {
	findings := make([]domain.Finding, 0)
	for _, column := range table.Columns {
		if !column.Nullable || !column.DefaultIsNull {
			continue
		}
		findings = append(findings, finding("TYPE-003", domain.SeverityNormal, "column", table.Name+"."+column.Name, "可空字段默认 NULL", "字段允许 NULL 且默认值为 NULL", "确认 NULL 是否表示未知；如不需要三值逻辑，应改为 NOT NULL 并设置业务默认值", "", fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE `%s` IS NULL;", table.Name, column.Name), "变更前备份受影响数据", true, "medium"))
	}
	return findings
}

func inconsistentCollation(snapshot domain.SchemaSnapshot, table domain.Table) []domain.Finding {
	findings := make([]domain.Finding, 0)
	for _, column := range table.Columns {
		if column.Collation == "" || table.Collation == "" || column.Collation == table.Collation {
			continue
		}
		findings = append(findings, finding("TYPE-004", domain.SeverityNormal, "column", table.Name+"."+column.Name, "字符集或排序规则不一致", fmt.Sprintf("列排序规则 %s 与表排序规则 %s 不一致", column.Collation, table.Collation), "确认排序和比较需求后统一字符集/排序规则", "", "", "变更前评估索引和排序行为", true, "medium"))
	}
	if table.Collation != "" && snapshot.DefaultCollation != "" && table.Collation != snapshot.DefaultCollation {
		findings = append(findings, finding("TYPE-004", domain.SeverityNormal, "table", table.Name, "表排序规则与数据库默认值不一致", fmt.Sprintf("表排序规则 %s，数据库默认规则 %s", table.Collation, snapshot.DefaultCollation), "确认业务兼容性后统一排序规则", "", "", "变更前评估索引和排序行为", true, "medium"))
	}
	return findings
}

func hasLeadingIndex(indexes []domain.Index, column string) bool {
	for _, index := range indexes {
		if len(index.Columns) > 0 && index.Columns[0] == column {
			return true
		}
	}
	return false
}

func samePrefix(prefix, values []string) bool {
	if len(prefix) > len(values) {
		return false
	}
	for index := range prefix {
		if prefix[index] != values[index] {
			return false
		}
	}
	return true
}

func quote(dialect, value string) string {
	if dialect == "postgres" || dialect == "postgresql" {
		return "\"" + strings.ReplaceAll(value, "\"", "\"\"") + "\""
	}
	return "`" + strings.ReplaceAll(value, "`", "``") + "`"
}

func uuidExpression(dialect string) string {
	if dialect == "postgres" || dialect == "postgresql" {
		return "gen_random_uuid()"
	}
	return "UUID()"
}

func finding(ruleID string, severity domain.Severity, objectType, objectName, title, evidence, recommendation, dml, precheck, rollback string, review bool, risk string) domain.Finding {
	checksum := sha256.Sum256([]byte(ruleID + "|" + objectName + "|" + title))
	return domain.Finding{ID: fmt.Sprintf("%s-%x", strings.ToLower(ruleID), checksum[:6]), RuleID: ruleID, Severity: severity, SeverityReason: title, ObjectType: objectType, ObjectName: objectName, Title: title, Evidence: evidence, Recommendation: recommendation, DMLSuggestion: dml, PrecheckSQL: precheck, RollbackPlan: rollback, RequiresManualReview: review, SQLRisk: risk, CreatedAt: time.Now().UTC()}
}
