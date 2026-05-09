package category_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
	"github.com/itk13201/money-rabbit/internal/usecase/category"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- mock ---

type mockRepo struct {
	categories []*entity.Category
	rules      []*entity.CategoryRule
}

func (m *mockRepo) ListCategories(_ context.Context) ([]*entity.Category, error) {
	return m.categories, nil
}
func (m *mockRepo) GetCategory(_ context.Context, id uuid.UUID) (*entity.Category, error) {
	for _, c := range m.categories {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, nil
}
func (m *mockRepo) CreateCategory(_ context.Context, input category.CreateInput) (*entity.Category, error) {
	c := &entity.Category{
		ID:        uuid.New(),
		Name:      input.Name,
		Color:     input.Color,
		Icon:      input.Icon,
		Type:      input.Type,
		SortOrder: input.SortOrder,
	}
	m.categories = append(m.categories, c)
	return c, nil
}
func (m *mockRepo) UpdateCategory(_ context.Context, id uuid.UUID, input category.UpdateInput) (*entity.Category, error) {
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
func (m *mockRepo) DeleteCategory(_ context.Context, id uuid.UUID) error {
	filtered := m.categories[:0]
	for _, c := range m.categories {
		if c.ID != id {
			filtered = append(filtered, c)
		}
	}
	m.categories = filtered
	return nil
}
func (m *mockRepo) ListAllRules(_ context.Context) ([]*entity.CategoryRule, error) {
	return m.rules, nil
}
func (m *mockRepo) CreateRule(_ context.Context, input category.CreateRuleInput) (*entity.CategoryRule, error) {
	r := &entity.CategoryRule{
		ID:         uuid.New(),
		Keyword:    input.Keyword,
		Priority:   input.Priority,
		CategoryID: input.CategoryID,
	}
	m.rules = append(m.rules, r)
	return r, nil
}
func (m *mockRepo) UpdateRule(_ context.Context, id uuid.UUID, input category.UpdateRuleInput) (*entity.CategoryRule, error) {
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
func (m *mockRepo) DeleteRule(_ context.Context, id uuid.UUID) error {
	filtered := m.rules[:0]
	for _, r := range m.rules {
		if r.ID != id {
			filtered = append(filtered, r)
		}
	}
	m.rules = filtered
	return nil
}

// --- tests ---

func TestList(t *testing.T) {
	repo := &mockRepo{
		categories: []*entity.Category{
			{ID: uuid.New(), Name: "食費", Type: entity.CategoryTypeExpense},
			{ID: uuid.New(), Name: "給与", Type: entity.CategoryTypeIncome},
		},
	}
	uc := category.New(repo)

	cats, err := uc.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, cats, 2)
	assert.Equal(t, "食費", cats[0].Name)
}

func TestCreate(t *testing.T) {
	repo := &mockRepo{}
	uc := category.New(repo)

	cat, err := uc.Create(context.Background(), category.CreateInput{
		Name:  "交通費",
		Color: "#3b82f6",
		Icon:  "🚃",
		Type:  entity.CategoryTypeExpense,
	})
	require.NoError(t, err)
	assert.Equal(t, "交通費", cat.Name)
	assert.Equal(t, entity.CategoryTypeExpense, cat.Type)
	assert.NotEqual(t, uuid.Nil, cat.ID)

	// リポジトリに追加されていること
	assert.Len(t, repo.categories, 1)
}

func TestDelete(t *testing.T) {
	id := uuid.New()
	repo := &mockRepo{
		categories: []*entity.Category{
			{ID: id, Name: "削除対象"},
		},
	}
	uc := category.New(repo)

	err := uc.Delete(context.Background(), id)
	require.NoError(t, err)
	assert.Empty(t, repo.categories)
}

func TestCreateRule(t *testing.T) {
	catID := uuid.New()
	repo := &mockRepo{}
	uc := category.New(repo)

	rule, err := uc.CreateRule(context.Background(), category.CreateRuleInput{
		CategoryID: catID,
		Keyword:    "スーパー",
		Priority:   10,
	})
	require.NoError(t, err)
	assert.Equal(t, "スーパー", rule.Keyword)
	assert.Equal(t, 10, rule.Priority)
	assert.Equal(t, catID, rule.CategoryID)
}

func TestGet(t *testing.T) {
	id := uuid.New()
	repo := &mockRepo{
		categories: []*entity.Category{
			{ID: id, Name: "食費", Type: entity.CategoryTypeExpense},
		},
	}
	uc := category.New(repo)

	cat, err := uc.Get(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, cat)
	assert.Equal(t, id, cat.ID)
	assert.Equal(t, "食費", cat.Name)
}

func TestGet_NotFound(t *testing.T) {
	repo := &mockRepo{}
	uc := category.New(repo)

	cat, err := uc.Get(context.Background(), uuid.New())
	require.NoError(t, err)
	assert.Nil(t, cat)
}

func TestUpdate(t *testing.T) {
	id := uuid.New()
	repo := &mockRepo{
		categories: []*entity.Category{
			{ID: id, Name: "旧名前", Type: entity.CategoryTypeExpense},
		},
	}
	uc := category.New(repo)

	newName := "新名前"
	cat, err := uc.Update(context.Background(), id, category.UpdateInput{Name: &newName})
	require.NoError(t, err)
	require.NotNil(t, cat)
	assert.Equal(t, newName, cat.Name)
}

func TestUpdateRule(t *testing.T) {
	id := uuid.New()
	repo := &mockRepo{
		rules: []*entity.CategoryRule{
			{ID: id, Keyword: "旧キーワード", Priority: 5},
		},
	}
	uc := category.New(repo)

	newKw := "新キーワード"
	rule, err := uc.UpdateRule(context.Background(), id, category.UpdateRuleInput{Keyword: &newKw})
	require.NoError(t, err)
	require.NotNil(t, rule)
	assert.Equal(t, newKw, rule.Keyword)
}

func TestDeleteRule(t *testing.T) {
	id := uuid.New()
	repo := &mockRepo{
		rules: []*entity.CategoryRule{
			{ID: id, Keyword: "削除対象"},
		},
	}
	uc := category.New(repo)

	err := uc.DeleteRule(context.Background(), id)
	require.NoError(t, err)
	assert.Empty(t, repo.rules)
}
