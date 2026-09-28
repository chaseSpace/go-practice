package scanner

import (
	"context"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"agentgo/p_backend/internal/domain"
)

func TestMySQLCollectorIntegration(t *testing.T) {
	rawURL := os.Getenv("AGENTGO_TEST_MYSQL_URL")
	if rawURL == "" {
		t.Skip("set AGENTGO_TEST_MYSQL_URL to run MySQL integration test")
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatal(err)
	}
	port := 3306
	if parsed.Port() != "" {
		port, err = strconv.Atoi(parsed.Port())
		if err != nil {
			t.Fatal(err)
		}
	}
	password, _ := parsed.User.Password()
	config := domain.ConnectionConfig{Dialect: "mysql", Host: parsed.Hostname(), Port: port, Database: strings.TrimPrefix(parsed.Path, "/"), Username: parsed.User.Username(), Password: password}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	snapshot, err := (MySQLCollector{ConnectTimeout: 5 * time.Second}).Collect(ctx, config, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Tables) == 0 {
		t.Fatal("expected collected tables")
	}
}
