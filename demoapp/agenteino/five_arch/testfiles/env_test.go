package testfiles

import (
	"errors"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// loadLearningEnv 加载本地学习配置。
// godotenv.Load 不会覆盖 shell 中已经存在的变量，因此命令行和 CI 配置优先。
func loadLearningEnv(t *testing.T) {
	t.Helper()
	for _, file := range []string{
		".env",
		"../../base_loop/testfiles/.env",
		"../.env",
		"testfiles/.env",
	} {
		_, err := os.Stat(file)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			t.Fatalf("检查 %s：%v", file, err)
		}
		if err := godotenv.Load(file); err != nil {
			t.Fatalf("加载 %s：%v", file, err)
		}
		t.Logf("已从 %s 加载本地 LLM 配置", file)
	}
}

func llmEnv(primary, legacy string) string {
	if value := os.Getenv(primary); value != "" {
		return value
	}
	return os.Getenv(legacy)
}
