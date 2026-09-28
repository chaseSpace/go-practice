package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoadWithRuleAndPromptValidatesRequiredLLM(t *testing.T) {
	configPath := filepath.Join("..", "..", "config", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	var raw struct {
		LLM struct {
			APIKeyEnv string `yaml:"api_key_env"`
		} `yaml:"llm"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	t.Setenv(raw.LLM.APIKeyEnv, "test-key")
	loaded, catalog, err := LoadWithRuleAndPrompt(configPath)
	if err != nil {
		t.Fatalf("load valid configuration: %v", err)
	}
	if !loaded.LLM.Enabled || catalog.PromptForAIRule("AI-REPORT-001") == "" {
		t.Fatal("expected enabled LLM and report AI rule prompt")
	}

	t.Setenv(raw.LLM.APIKeyEnv, "")
	_, _, err = LoadWithRuleAndPrompt(configPath)
	if err == nil || !strings.Contains(err.Error(), "llm api key environment variable") {
		t.Fatalf("expected missing LLM key error, got %v", err)
	}
}
