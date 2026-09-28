package testfiles

import (
	"errors"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// loadLearningEnv reads the local learning configuration when it exists.
// godotenv.Load does not overwrite values already supplied by the shell, so
// CI and explicit command-line configuration take precedence over .env.
func loadLearningEnv(t *testing.T) {
	t.Helper()
	// go test normally runs this package from testfiles/, so .env is the
	// learner's testfiles/.env. The other paths make direct test-binary runs
	// and a root-level .env work too.
	for _, file := range []string{".env", "../.env", "testfiles/.env"} {
		_, err := os.Stat(file)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			t.Fatalf("inspect %s: %v", file, err)
		}
		if err := godotenv.Load(file); err != nil {
			t.Fatalf("load %s: %v", file, err)
		}
		t.Logf("loaded local LLM configuration from %s", file)
	}
}

func llmEnv(primary, legacy string) string {
	if value := os.Getenv(primary); value != "" {
		return value
	}
	return os.Getenv(legacy)
}
