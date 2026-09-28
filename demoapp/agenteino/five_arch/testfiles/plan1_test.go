// Plan 1：使用真实 Tool 与真实 LLM 实现提示链。
package testfiles

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	toolutils "github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type collectOfficialMaterialInput struct {
	Topic string `json:"topic" jsonschema_description:"需要收集资料的主题，例如 Eino Agent 开发"`
}

var (
	scriptAndStylePattern = regexp.MustCompile(`(?is)<(script|style)[^>]*>.*?</(script|style)>`)
	htmlTagPattern        = regexp.MustCompile(`(?s)<[^>]+>`)
)

// EINO_RUN_LIVE_TESTS=1  go test -run=TestPlan1RealPromptChain -v
// TestPlan1RealPromptChain 的执行过程为：
// 资料收集 Agent（真实 HTTP Tool）-> 大纲 Agent（LLM）-> 初稿 Agent（LLM）。
// 这是集成示例，不断言模型的随机输出；运行时打印轨迹，结束后打印最终初稿。
func TestPlan1RealPromptChain(t *testing.T) {
	loadLearningEnv(t)
	if os.Getenv("EINO_RUN_LIVE_TESTS") != "1" {
		t.Skip("设置 EINO_RUN_LIVE_TESTS=1 后运行真实 Tool 与 LLM 提示链")
	}

	apiKey := llmEnv("LLM_API_KEY", "OPENAI_API_KEY")
	modelName := llmEnv("LLM_MODEL", "OPENAI_MODEL")
	if apiKey == "" || modelName == "" {
		t.Skip("缺少 LLM_API_KEY 或 LLM_MODEL")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  apiKey,
		Model:   modelName,
		BaseURL: llmEnv("LLM_ENDPOINT", "OPENAI_BASE_URL"),
	})
	if err != nil {
		t.Fatalf("创建 ChatModel：%v", err)
	}

	collectorTool, err := newOfficialMaterialTool()
	if err != nil {
		t.Fatalf("创建资料收集 Tool：%v", err)
	}
	collector, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "official_material_collector",
		Description: "从 Eino 官方文档收集可靠资料的研究员",
		Instruction: `你是严谨的技术资料研究员。你的唯一任务是为后续编辑收集 Eino 官方事实。
必须先调用 collect_eino_official_material 工具，并把用户主题原样作为 topic。

至多调用2次该工具。工具结果会直接传给后续编辑；不要凭记忆补充事实，不要写素材总结、大纲或初稿。`,
		Model: chatModel,
		ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{collectorTool},
		}, ReturnDirectly: map[string]bool{
			// Tool 完成后立刻结束收集 Agent，不再让模型继续发起 ToolCall。
			"collect_eino_official_material": true,
		}},
		// 收集 Agent 只需要一轮模型调用来选择 Tool。
		MaxIterations: 2,
	})
	if err != nil {
		t.Fatalf("创建资料收集 Agent：%v", err)
	}

	outlineEditor, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "outline_editor",
		Description: "把官方素材整理成教学文章大纲的资深编辑",
		Instruction: `你是资深技术编辑。阅读对话中资料收集员提供的素材卡片，为 Go Agent 初学者设计文章大纲。
大纲必须由浅入深，包含标题、目标读者、3 至 5 个一级章节，以及每章要回答的核心问题。
只输出大纲，不补充素材中没有的事实，不提前写正文。`,
		Model: chatModel,
	})
	if err != nil {
		t.Fatalf("创建大纲 Agent：%v", err)
	}

	draftWriter, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "draft_writer",
		Description: "依据官方素材和已确认大纲撰写初稿的技术作者",
		Instruction: `你是擅长 Go 与 Agent 教学的技术作者。严格依据前面的官方素材卡片和编辑大纲撰写中文初稿。
使用清晰短句解释术语，先说明问题再给概念，避免营销语言；重要结论附上素材中的来源 URL。
只输出完整初稿，不输出写作过程、评价、免责声明或额外说明。`,
		Model: chatModel,
	})
	if err != nil {
		t.Fatalf("创建初稿 Agent：%v", err)
	}

	chain, err := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
		Name:        "eino_article_chain",
		Description: "官方资料收集、编辑大纲、撰写初稿的顺序提示链",
		SubAgents:   []adk.Agent{collector, outlineEditor, draftWriter},
	})
	if err != nil {
		t.Fatalf("创建 SequentialAgent：%v", err)
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: chain})
	traceStarted := time.Now()
	fmt.Printf("[运行轨迹] +0s workflow=eino_article_chain stage=start\n")
	iterator := runner.Query(ctx, "面向有 Go 基础但没开发过 Agent 的读者，写一篇 Eino Agent 入门文章。")
	var finalDraft string
	eventIndex := 0
	for {
		event, ok := iterator.Next()
		if !ok {
			break
		}
		eventIndex++
		printPlan1AgentEvent(traceStarted, eventIndex, event)
		if event.Err != nil {
			t.Fatalf("提示链执行失败：%v", event.Err)
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}
		message := event.Output.MessageOutput.Message
		if message != nil && message.Role == schema.Assistant && message.Content != "" {
			// 顺序链中最后一条非空 assistant 消息就是初稿 Agent 的输出。
			finalDraft = message.Content
		}
	}

	fmt.Printf("[运行轨迹] +%s workflow=eino_article_chain stage=end events=%d\n", plan1Elapsed(traceStarted), eventIndex)
	t.Logf("\n========== 最终初稿 ==========\n%s", finalDraft)
}

// printPlan1AgentEvent 只使用 ADK 的 AgentEvent，因此始终能显示所属 Agent。
// 不打印中间正文，只记录阶段、工具和内容长度。
func printPlan1AgentEvent(started time.Time, index int, event *adk.AgentEvent) {
	stage := "event"
	role := ""
	toolName := ""
	contentLength := 0
	toolCalls := 0
	if event.Err != nil {
		stage = "error"
	}
	if event.Action != nil {
		stage = "action"
	}
	if event.Output != nil && event.Output.MessageOutput != nil {
		output := event.Output.MessageOutput
		role = string(output.Role)
		toolName = output.ToolName
		if output.Message != nil {
			contentLength = len([]rune(output.Message.Content))
			toolCalls = len(output.Message.ToolCalls)
			if toolCalls > 0 {
				stage = "tool_request"
			} else if output.Role == schema.Tool {
				stage = "tool_result"
			} else if output.Role == schema.Assistant {
				stage = "assistant_output"
			}
		}
	}
	fmt.Printf(
		"[运行轨迹] +%s event=%d agent=%s stage=%s role=%s tool=%s tool_calls=%d chars=%d\n",
		plan1Elapsed(started), index, event.AgentName, stage, role, toolName, toolCalls, contentLength,
	)
}

func plan1Elapsed(started time.Time) time.Duration {
	return time.Since(started).Round(time.Millisecond)
}

// newOfficialMaterialTool 创建真实的资料收集 Tool。
// URL 固定为官方白名单，不接受模型生成的任意地址，避免形成 SSRF 漏洞。
func newOfficialMaterialTool() (tool.InvokableTool, error) {
	return toolutils.InferTool(
		"collect_eino_official_material",
		"从 CloudWeGo Eino 官方页面收集资料。讨论 Eino、组件或 ADK 时必须使用此工具。",
		func(ctx context.Context, input collectOfficialMaterialInput) (string, error) {
			urls := []string{
				"https://www.cloudwego.io/docs/eino/overview/",
				"https://www.cloudwego.io/docs/eino/core_modules/eino_adk/",
			}
			client := &http.Client{Timeout: 15 * time.Second}
			sections := []string{"主题：" + input.Topic}
			for _, sourceURL := range urls {
				text, err := fetchOfficialPage(ctx, client, sourceURL)
				if err != nil {
					return "", err
				}
				sections = append(sections, fmt.Sprintf("来源：%s\n正文摘录：%s", sourceURL, text))
			}
			return strings.Join(sections, "\n\n"), nil
		},
	)
}

func fetchOfficialPage(ctx context.Context, client *http.Client, sourceURL string) (string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return "", fmt.Errorf("创建官方页面请求：%w", err)
	}
	request.Header.Set("User-Agent", "agenteino-learning/1.0")
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("访问 %s：%w", sourceURL, err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("访问 %s 返回状态码 %d", sourceURL, response.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, 512<<10))
	if err != nil {
		return "", fmt.Errorf("读取 %s：%w", sourceURL, err)
	}
	plainText := scriptAndStylePattern.ReplaceAllString(string(body), " ")
	plainText = htmlTagPattern.ReplaceAllString(plainText, " ")
	plainText = html.UnescapeString(plainText)
	plainText = strings.Join(strings.Fields(plainText), " ")
	return truncatePlan1Text(plainText, 4_000), nil
}

func truncatePlan1Text(text string, maxRunes int) string {
	runes := []rune(text)
	if len(runes) <= maxRunes {
		return text
	}
	return string(runes[:maxRunes]) + "……"
}
