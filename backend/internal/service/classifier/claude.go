package classifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
)

const (
	anthropicAPIURL  = "https://api.anthropic.com/v1/messages"
	model            = "claude-sonnet-4-6"
	anthropicVersion = "2023-06-01"
)

// Classifier classifies transaction descriptions using the Claude API.
type Classifier struct {
	apiKey string
	client *http.Client
}

func New(apiKey string) *Classifier {
	return &Classifier{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// Classify sends transaction descriptions to Claude and returns a map of
// description → category ID. Descriptions without a confident match are omitted.
func (c *Classifier) Classify(ctx context.Context, descriptions []string, categories []*entity.Category) (map[string]*uuid.UUID, error) {
	if len(descriptions) == 0 || len(categories) == 0 {
		return map[string]*uuid.UUID{}, nil
	}

	prompt := buildPrompt(descriptions, categories)
	slog.InfoContext(ctx, "claude.api started",
		slog.Group("extra",
			"description_count", len(descriptions),
			"category_count", len(categories),
		),
	)
	raw, usage, err := c.callAPI(ctx, prompt)
	if err != nil {
		slog.ErrorContext(ctx, "claude.api failed",
			slog.Group("extra", "error", err),
		)
		return nil, fmt.Errorf("claude api: %w", err)
	}
	slog.InfoContext(ctx, "claude.api finished",
		slog.Group("extra",
			"input_tokens", usage.InputTokens,
			"output_tokens", usage.OutputTokens,
		),
	)

	return parseResponse(raw, categories), nil
}

func buildPrompt(descriptions []string, categories []*entity.Category) string {
	var sb strings.Builder
	sb.WriteString("以下のカテゴリ一覧から、各取引の摘要に最も適したカテゴリIDを選んでください。\n\n")
	sb.WriteString("カテゴリ一覧:\n")
	for _, cat := range categories {
		sb.WriteString(fmt.Sprintf("- id:%s name:%s type:%s\n", cat.ID, cat.Name, cat.Type))
	}
	sb.WriteString("\n取引摘要一覧:\n")
	for i, desc := range descriptions {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, desc))
	}
	sb.WriteString("\n以下のJSON形式で回答してください。マッチしない場合は null を指定してください:\n")
	sb.WriteString(`{"results":[{"description":"摘要テキスト","category_id":"UUID or null"}]}`)
	return sb.String()
}

type apiRequest struct {
	Model     string       `json:"model"`
	MaxTokens int          `json:"max_tokens"`
	Messages  []apiMessage `json:"messages"`
}

type apiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type apiUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type apiResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Usage apiUsage `json:"usage"`
}

func (c *Classifier) callAPI(ctx context.Context, prompt string) (string, apiUsage, error) {
	body, err := json.Marshal(apiRequest{
		Model:     model,
		MaxTokens: 4096,
		Messages:  []apiMessage{{Role: "user", Content: prompt}},
	})
	if err != nil {
		return "", apiUsage{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, anthropicAPIURL, bytes.NewReader(body))
	if err != nil {
		return "", apiUsage{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", apiUsage{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", apiUsage{}, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", apiUsage{}, err
	}

	for _, block := range apiResp.Content {
		if block.Type == "text" {
			return block.Text, apiResp.Usage, nil
		}
	}
	return "", apiUsage{}, fmt.Errorf("no text content in response")
}

type classifyResult struct {
	Description string  `json:"description"`
	CategoryID  *string `json:"category_id"`
}

type classifyResponse struct {
	Results []classifyResult `json:"results"`
}

func parseResponse(raw string, categories []*entity.Category) map[string]*uuid.UUID {
	// Extract JSON from the response (Claude may include surrounding text)
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start == -1 || end == -1 || end <= start {
		return map[string]*uuid.UUID{}
	}
	raw = raw[start : end+1]

	var resp classifyResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		return map[string]*uuid.UUID{}
	}

	// Build a valid category ID set for validation
	validIDs := make(map[string]uuid.UUID, len(categories))
	for _, cat := range categories {
		validIDs[cat.ID.String()] = cat.ID
	}

	result := make(map[string]*uuid.UUID, len(resp.Results))
	for _, r := range resp.Results {
		if r.CategoryID == nil || *r.CategoryID == "null" || *r.CategoryID == "" {
			result[r.Description] = nil
			continue
		}
		if id, ok := validIDs[*r.CategoryID]; ok {
			idCopy := id
			result[r.Description] = &idCopy
		}
	}
	return result
}
