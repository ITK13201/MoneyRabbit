package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
	"github.com/itk13201/money-rabbit/internal/handler"
	categoryUC "github.com/itk13201/money-rabbit/internal/usecase/category"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// --- mock ---

type mockCategoryUsecase struct {
	categories []*entity.Category
	rules      []*entity.CategoryRule
}

func (m *mockCategoryUsecase) List(_ context.Context) ([]*entity.Category, error) {
	return m.categories, nil
}

func (m *mockCategoryUsecase) Get(_ context.Context, id uuid.UUID) (*entity.Category, error) {
	for _, c := range m.categories {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, nil
}

func (m *mockCategoryUsecase) Create(_ context.Context, input categoryUC.CreateInput) (*entity.Category, error) {
	c := &entity.Category{
		ID:    uuid.New(),
		Name:  input.Name,
		Color: input.Color,
		Icon:  input.Icon,
		Type:  input.Type,
	}
	m.categories = append(m.categories, c)
	return c, nil
}

func (m *mockCategoryUsecase) Update(_ context.Context, id uuid.UUID, input categoryUC.UpdateInput) (*entity.Category, error) {
	for _, c := range m.categories {
		if c.ID == id {
			if input.Name != nil {
				c.Name = *input.Name
			}
			return c, nil
		}
	}
	return nil, nil
}

func (m *mockCategoryUsecase) Delete(_ context.Context, id uuid.UUID) error {
	filtered := m.categories[:0]
	for _, c := range m.categories {
		if c.ID != id {
			filtered = append(filtered, c)
		}
	}
	m.categories = filtered
	return nil
}

func (m *mockCategoryUsecase) CreateRule(_ context.Context, input categoryUC.CreateRuleInput) (*entity.CategoryRule, error) {
	r := &entity.CategoryRule{
		ID:         uuid.New(),
		Keyword:    input.Keyword,
		Priority:   input.Priority,
		CategoryID: input.CategoryID,
	}
	m.rules = append(m.rules, r)
	return r, nil
}

func (m *mockCategoryUsecase) UpdateRule(_ context.Context, id uuid.UUID, input categoryUC.UpdateRuleInput) (*entity.CategoryRule, error) {
	for _, r := range m.rules {
		if r.ID == id {
			if input.Keyword != nil {
				r.Keyword = *input.Keyword
			}
			return r, nil
		}
	}
	return nil, nil
}

func (m *mockCategoryUsecase) DeleteRule(_ context.Context, id uuid.UUID) error {
	filtered := m.rules[:0]
	for _, r := range m.rules {
		if r.ID != id {
			filtered = append(filtered, r)
		}
	}
	m.rules = filtered
	return nil
}

// --- helpers ---

func newCategoryRouter(uc *mockCategoryUsecase) *gin.Engine {
	r := gin.New()
	h := handler.NewCategoryHandler(uc)
	r.GET("/api/categories", h.List)
	r.GET("/api/categories/:id", h.Get)
	r.POST("/api/categories", h.Create)
	r.PUT("/api/categories/:id", h.Update)
	r.DELETE("/api/categories/:id", h.Delete)
	r.POST("/api/category-rules", h.CreateRule)
	r.PUT("/api/category-rules/:id", h.UpdateRule)
	r.DELETE("/api/category-rules/:id", h.DeleteRule)
	return r
}

func doRequest(r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func doRawRequest(r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// --- tests ---

func TestCategoryHandler_List(t *testing.T) {
	uc := &mockCategoryUsecase{
		categories: []*entity.Category{
			{ID: uuid.New(), Name: "食費", Type: entity.CategoryTypeExpense},
			{ID: uuid.New(), Name: "給与", Type: entity.CategoryTypeIncome},
		},
	}
	r := newCategoryRouter(uc)

	w := doRequest(r, http.MethodGet, "/api/categories", nil)
	require.Equal(t, http.StatusOK, w.Code)

	var cats []*entity.Category
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &cats))
	assert.Len(t, cats, 2)
	assert.Equal(t, "食費", cats[0].Name)
}

func TestCategoryHandler_Create(t *testing.T) {
	uc := &mockCategoryUsecase{}
	r := newCategoryRouter(uc)

	w := doRequest(r, http.MethodPost, "/api/categories", map[string]any{
		"name":  "交通費",
		"color": "#3b82f6",
		"icon":  "🚃",
		"type":  "expense",
	})
	require.Equal(t, http.StatusCreated, w.Code)

	var cat entity.Category
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &cat))
	assert.Equal(t, "交通費", cat.Name)
	assert.NotEqual(t, uuid.Nil, cat.ID)
}

func TestCategoryHandler_Delete(t *testing.T) {
	id := uuid.New()
	uc := &mockCategoryUsecase{
		categories: []*entity.Category{
			{ID: id, Name: "削除対象"},
		},
	}
	r := newCategoryRouter(uc)

	w := doRequest(r, http.MethodDelete, "/api/categories/"+id.String(), nil)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, uc.categories)
}

func TestCategoryHandler_Get(t *testing.T) {
	id := uuid.New()
	uc := &mockCategoryUsecase{
		categories: []*entity.Category{
			{ID: id, Name: "食費", Type: entity.CategoryTypeExpense},
		},
	}
	r := newCategoryRouter(uc)

	w := doRequest(r, http.MethodGet, "/api/categories/"+id.String(), nil)
	require.Equal(t, http.StatusOK, w.Code)

	var cat entity.Category
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &cat))
	assert.Equal(t, id, cat.ID)
	assert.Equal(t, "食費", cat.Name)
}

func TestCategoryHandler_Update(t *testing.T) {
	id := uuid.New()
	uc := &mockCategoryUsecase{
		categories: []*entity.Category{
			{ID: id, Name: "旧名前", Type: entity.CategoryTypeExpense},
		},
	}
	r := newCategoryRouter(uc)

	newName := "新名前"
	w := doRequest(r, http.MethodPut, "/api/categories/"+id.String(), map[string]any{
		"name": newName,
	})
	require.Equal(t, http.StatusOK, w.Code)

	var cat entity.Category
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &cat))
	assert.Equal(t, newName, cat.Name)
}

func TestCategoryHandler_UpdateRule(t *testing.T) {
	catID := uuid.New()
	ruleID := uuid.New()
	uc := &mockCategoryUsecase{
		rules: []*entity.CategoryRule{
			{ID: ruleID, Keyword: "旧キーワード", Priority: 10, CategoryID: catID},
		},
	}
	r := newCategoryRouter(uc)

	newKw := "新キーワード"
	w := doRequest(r, http.MethodPut, "/api/category-rules/"+ruleID.String(), map[string]any{
		"keyword": newKw,
	})
	require.Equal(t, http.StatusOK, w.Code)

	var rule entity.CategoryRule
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &rule))
	assert.Equal(t, newKw, rule.Keyword)
}

func TestCategoryHandler_DeleteRule(t *testing.T) {
	id := uuid.New()
	uc := &mockCategoryUsecase{
		rules: []*entity.CategoryRule{
			{ID: id, Keyword: "削除対象"},
		},
	}
	r := newCategoryRouter(uc)

	w := doRequest(r, http.MethodDelete, "/api/category-rules/"+id.String(), nil)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, uc.rules)
}

func TestCategoryHandler_CreateRule(t *testing.T) {
	catID := uuid.New()
	uc := &mockCategoryUsecase{}
	r := newCategoryRouter(uc)

	w := doRequest(r, http.MethodPost, "/api/category-rules", map[string]any{
		"category_id": catID.String(),
		"keyword":     "スーパー",
		"priority":    10,
	})
	require.Equal(t, http.StatusCreated, w.Code)

	var rule entity.CategoryRule
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &rule))
	assert.Equal(t, "スーパー", rule.Keyword)
	assert.Equal(t, catID, rule.CategoryID)
}
