package classifier

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildPrompt_ContainsDescriptionsAndCategories(t *testing.T) {
	catID := uuid.New()
	cats := []*entity.Category{
		{ID: catID, Name: "食費", Type: entity.CategoryTypeExpense},
	}
	descs := []string{"スーパーで買い物", "コンビニ"}

	prompt := buildPrompt(descs, cats)

	assert.Contains(t, prompt, catID.String())
	assert.Contains(t, prompt, "食費")
	assert.Contains(t, prompt, "スーパーで買い物")
	assert.Contains(t, prompt, "コンビニ")
}

func TestParseResponse_ValidJSON(t *testing.T) {
	catID := uuid.New()
	cats := []*entity.Category{{ID: catID, Name: "食費"}}

	raw := buildRawResponse([]classifyResult{
		{Description: "スーパー", CategoryID: ptr(catID.String())},
	})

	result := parseResponse(raw, cats)
	require.Contains(t, result, "スーパー")
	require.NotNil(t, result["スーパー"])
	assert.Equal(t, catID, *result["スーパー"])
}

func TestParseResponse_NullCategoryID(t *testing.T) {
	cats := []*entity.Category{{ID: uuid.New(), Name: "食費"}}

	raw := buildRawResponse([]classifyResult{
		{Description: "不明な取引", CategoryID: ptr("null")},
	})

	result := parseResponse(raw, cats)
	require.Contains(t, result, "不明な取引")
	assert.Nil(t, result["不明な取引"])
}

func TestParseResponse_InvalidCategoryID(t *testing.T) {
	cats := []*entity.Category{{ID: uuid.New(), Name: "食費"}}

	raw := buildRawResponse([]classifyResult{
		{Description: "test", CategoryID: ptr(uuid.New().String())}, // 存在しないカテゴリID
	})

	result := parseResponse(raw, cats)
	assert.NotContains(t, result, "test")
}

func TestParseResponse_MalformedJSON(t *testing.T) {
	cats := []*entity.Category{{ID: uuid.New(), Name: "食費"}}
	result := parseResponse("this is not json", cats)
	assert.Empty(t, result)
}

func TestParseResponse_JSONWithSurroundingText(t *testing.T) {
	catID := uuid.New()
	cats := []*entity.Category{{ID: catID, Name: "食費"}}

	jsonPart := buildRawResponse([]classifyResult{
		{Description: "スーパー", CategoryID: ptr(catID.String())},
	})
	raw := "以下が結果です:\n" + jsonPart + "\n以上です。"

	result := parseResponse(raw, cats)
	require.Contains(t, result, "スーパー")
	assert.NotNil(t, result["スーパー"])
}

func TestParseResponse_EmptyDescriptions(t *testing.T) {
	cats := []*entity.Category{{ID: uuid.New(), Name: "食費"}}
	raw := buildRawResponse(nil)
	result := parseResponse(raw, cats)
	assert.Empty(t, result)
}

// --- helpers ---

func buildRawResponse(results []classifyResult) string {
	resp := classifyResponse{Results: results}
	b, _ := json.Marshal(resp)
	return string(b)
}

func ptr(s string) *string {
	return &s
}

func TestBuildPrompt_EmptyInputs(t *testing.T) {
	prompt := buildPrompt([]string{}, []*entity.Category{})
	assert.NotEmpty(t, prompt) // プロンプトテンプレート自体は生成される
}

func TestBuildPrompt_NumberedDescriptions(t *testing.T) {
	descs := []string{"item1", "item2", "item3"}
	prompt := buildPrompt(descs, []*entity.Category{})
	assert.True(t, strings.Contains(prompt, "1. item1"))
	assert.True(t, strings.Contains(prompt, "2. item2"))
	assert.True(t, strings.Contains(prompt, "3. item3"))
}
