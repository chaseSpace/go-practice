package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type agentArchitectureReview struct {
	Recommendation  string   `json:"recommendation" jsonschema_description:"The recommended architecture option"`
	Rationale       string   `json:"rationale" jsonschema_description:"Why this option is recommended"`
	MainRisk        string   `json:"main_risk" jsonschema_description:"The main production risk"`
	ConfidenceScore int      `json:"confidence_score" jsonschema:"minimum=1,maximum=10" jsonschema_description:"Confidence score from 1 to 10"`
	NeedsEval       *bool    `json:"needs_eval" jsonschema_description:"Whether an evaluation suite is required before rollout"`
	KeyActions      []string `json:"key_actions" jsonschema_description:"Two to four practical next actions"`
}

func generateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	return reflector.Reflect(v)
}

var agentArchitectureReviewResponseSchema = generateSchema[agentArchitectureReview]()

const (
	maxStructuredOutputAttempts = 3
	structuredOutputTimeout     = 10 * time.Second
	agentArchitectureQuestion   = "为中型 SaaS 公司建设企业级 AI Agent，方案 A 是用现成平台快速上线，方案 B 是用 LangGraph 自研核心编排并接入内部 Agent Registry。请用 JSON 给出推荐方案、理由、主要风险、信心分、是否需要评估和关键行动。"
)

func TestOpenai(t *testing.T) {
	if err := loadEnvFile("logs/llm.env"); err != nil {
		t.Fatalf("load logs/llm.env: %v", err)
	}

	baseURL := os.Getenv("LLM_BASE_URL")
	apiKey := os.Getenv("LLM_API_KEY")
	model := os.Getenv("LLM_MODEL")
	if baseURL == "" || apiKey == "" || model == "" {
		t.Skip("LLM_BASE_URL, LLM_API_KEY, or LLM_MODEL is not set")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "agent_architecture_review",
		Description: openai.String("A structured review of an AI agent architecture decision."),
		Schema:      agentArchitectureReviewResponseSchema,
		Strict:      openai.Bool(true),
	}

	review, content, err := generateAgentArchitectureReview(
		context.TODO(),
		&client,
		openai.ChatModel(model),
		schemaParam,
		agentArchitectureQuestion,
		t.Logf,
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("raw structured JSON: %s", content)
	t.Logf("%#v", review)
}

func generateAgentArchitectureReview(
	ctx context.Context,
	client *openai.Client,
	model openai.ChatModel,
	schemaParam openai.ResponseFormatJSONSchemaJSONSchemaParam,
	question string,
	logf func(format string, args ...any),
) (agentArchitectureReview, string, error) {
	var lastContent string
	var lastErr error

	for attempt := 1; attempt <= maxStructuredOutputAttempts; attempt++ {
		userPrompt := question
		if attempt > 1 && lastContent != "" {
			userPrompt = repairPrompt(lastContent, lastErr)
		}
		logf("structured output attempt %d/%d, timeout=%s", attempt, maxStructuredOutputAttempts, structuredOutputTimeout)

		attemptCtx, cancel := context.WithTimeout(ctx, structuredOutputTimeout)
		chatCompletion, err := client.Chat.Completions.New(
			attemptCtx, openai.ChatCompletionNewParams{
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage("Return only JSON. Use exactly these keys: recommendation(string), rationale(string), main_risk(string), confidence_score(number 1-10), needs_eval(boolean), key_actions(array of strings). Do not add extra keys."),
					openai.UserMessage(userPrompt),
				},
				Model:     model,
				MaxTokens: openai.Int(300),
				ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
					OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{JSONSchema: schemaParam},
				},
			},
		)
		cancel()
		if err != nil {
			lastErr = err
			logf("structured output attempt %d failed request: %v", attempt, err)
			continue
		}

		if len(chatCompletion.Choices) == 0 {
			return agentArchitectureReview{}, lastContent, fmt.Errorf("expected at least one completion choice")
		}

		lastContent = chatCompletion.Choices[0].Message.Content
		review, err := parseAgentArchitectureReview(lastContent)
		if err == nil {
			logf("structured output attempt %d succeeded", attempt)
			return review, lastContent, nil
		}

		lastErr = err
		logf("structured output attempt %d failed validation: %v; content: %s", attempt, err, lastContent)
	}

	return agentArchitectureReview{}, lastContent, fmt.Errorf("structured output invalid after %d attempts: %w; last content: %s", maxStructuredOutputAttempts, lastErr, lastContent)
}

func parseAgentArchitectureReview(content string) (agentArchitectureReview, error) {
	var review agentArchitectureReview
	if err := json.Unmarshal([]byte(content), &review); err != nil {
		return agentArchitectureReview{}, fmt.Errorf("unmarshal structured output: %w", err)
	}

	if missing := missingAgentArchitectureReviewFields(review); len(missing) > 0 {
		return agentArchitectureReview{}, fmt.Errorf("missing required fields: %s", strings.Join(missing, ", "))
	}

	return review, nil
}

func missingAgentArchitectureReviewFields(review agentArchitectureReview) []string {
	var missing []string

	if review.Recommendation == "" {
		missing = append(missing, "recommendation")
	}
	if review.Rationale == "" {
		missing = append(missing, "rationale")
	}
	if review.MainRisk == "" {
		missing = append(missing, "main_risk")
	}
	if review.ConfidenceScore < 1 || review.ConfidenceScore > 10 {
		missing = append(missing, "confidence_score")
	}
	if review.NeedsEval == nil {
		missing = append(missing, "needs_eval")
	}
	if len(review.KeyActions) == 0 {
		missing = append(missing, "key_actions")
	}

	return missing
}

func repairPrompt(invalidContent string, err error) string {
	return fmt.Sprintf(`Fix this JSON so it has exactly these keys and types: recommendation string, rationale string, main_risk string, confidence_score number 1-10, needs_eval boolean, key_actions array of strings.
Validation error: %v
JSON: %s`, err, invalidContent)
}

func loadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key == "" {
			continue
		}

		if _, exists := os.LookupEnv(key); !exists {
			if err := os.Setenv(key, value); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}
