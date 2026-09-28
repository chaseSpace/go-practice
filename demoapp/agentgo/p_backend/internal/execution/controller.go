package execution

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"agentgo/p_backend/internal/domain"
	"agentgo/p_backend/internal/service"

	"github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrDisabled         = errors.New("controlled execution is disabled")
	ErrConfirmation     = errors.New("execution confirmation is required")
	ErrTargetNotAllowed = errors.New("database target is not allowlisted")
	ErrUnsafeStatement  = errors.New("statement is not an allowed single DML statement")
)

type Config struct {
	Enabled      bool
	AllowedHosts map[string]struct{}
}
type Controller struct{ config Config }
type Result struct {
	RowsAffected int64     `json:"rowsAffected"`
	ExecutedAt   time.Time `json:"executedAt"`
}

func New(config Config) *Controller { return &Controller{config: config} }

func (c *Controller) Execute(ctx context.Context, finding domain.Finding, confirmation, connectionURL string, connection *domain.ConnectionConfig) (Result, error) {
	if !c.config.Enabled {
		return Result{}, ErrDisabled
	}
	if confirmation != "EXECUTE" {
		return Result{}, ErrConfirmation
	}
	config, err := service.ParseConnection(connectionURL, connection)
	if err != nil {
		return Result{}, err
	}
	if _, ok := c.config.AllowedHosts[config.Host]; !ok {
		return Result{}, ErrTargetNotAllowed
	}
	statement, err := allowedDML(finding.DMLSuggestion)
	if err != nil {
		return Result{}, err
	}
	db, err := openDatabase(config)
	if err != nil {
		return Result{}, err
	}
	defer db.Close()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return Result{}, err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, statement)
	if err != nil {
		return Result{}, err
	}
	if err := tx.Commit(); err != nil {
		return Result{}, err
	}
	count, _ := result.RowsAffected()
	return Result{RowsAffected: count, ExecutedAt: time.Now().UTC()}, nil
}

func allowedDML(value string) (string, error) {
	lines := make([]string, 0)
	for _, line := range strings.Split(value, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "--") {
			lines = append(lines, line)
		}
	}
	statement := strings.TrimSpace(strings.TrimSuffix(strings.Join(lines, "\n"), ";"))
	upper := strings.ToUpper(statement)
	if statement == "" || strings.Contains(statement, ";") || !(strings.HasPrefix(upper, "UPDATE ") || strings.HasPrefix(upper, "INSERT ") || strings.HasPrefix(upper, "DELETE ")) {
		return "", ErrUnsafeStatement
	}
	return statement, nil
}

func openDatabase(config domain.ConnectionConfig) (*sql.DB, error) {
	if config.Dialect == "mysql" {
		dsn := mysql.NewConfig()
		dsn.User, dsn.Passwd, dsn.Net, dsn.Addr, dsn.DBName = config.Username, config.Password, "tcp", net.JoinHostPort(config.Host, fmt.Sprintf("%d", config.Port)), config.Database
		return sql.Open("mysql", dsn.FormatDSN())
	}
	if config.Dialect == "postgres" || config.Dialect == "postgresql" {
		dsn := (&url.URL{Scheme: "postgres", User: url.UserPassword(config.Username, config.Password), Host: net.JoinHostPort(config.Host, fmt.Sprintf("%d", config.Port)), Path: config.Database}).String()
		return sql.Open("pgx", dsn)
	}
	return nil, errors.New("unsupported database dialect")
}
