package service

import (
	"testing"

	"agentgo/p_backend/internal/domain"
)

func TestParseConnection(t *testing.T) {
	connection, err := ParseConnection("mysql://readonly:secret@db.example:3307/app_db", nil)
	if err != nil {
		t.Fatal(err)
	}
	if connection.Host != "db.example" || connection.Port != 3307 || connection.Database != "app_db" || connection.Password != "secret" {
		t.Fatalf("unexpected connection: %#v", connection)
	}
	if _, err := ParseConnection("mysql://readonly@db.example/app_db", &domain.ConnectionConfig{}); err == nil {
		t.Fatal("expected mutually-exclusive connection validation error")
	}
	if _, err := ParseConnection("mysql://readonly:secret@169.254.169.254:3306/app", nil); err == nil {
		t.Fatal("expected metadata host validation error")
	}
}
