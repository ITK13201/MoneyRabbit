package persistence_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
	"github.com/itk13201/money-rabbit/internal/service/persistence"
	categoryUC "github.com/itk13201/money-rabbit/internal/usecase/category"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryRepository_CRUD(t *testing.T) {
	client := setupTestDB(t)
	repo := persistence.NewCategoryRepository(client)
	ctx := context.Background()

	// Create
	cat, err := repo.CreateCategory(ctx, categoryUC.CreateInput{
		Name:      "食費",
		Color:     "#ff0000",
		Icon:      "🍜",
		Type:      entity.CategoryTypeExpense,
		SortOrder: 1,
	})
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, cat.ID)
	assert.Equal(t, "食費", cat.Name)

	// List
	cats, err := repo.ListCategories(ctx)
	require.NoError(t, err)
	require.Len(t, cats, 1)
	assert.Equal(t, cat.ID, cats[0].ID)

	// Get
	got, err := repo.GetCategory(ctx, cat.ID)
	require.NoError(t, err)
	assert.Equal(t, cat.ID, got.ID)
	assert.Equal(t, "食費", got.Name)

	// Update
	newName := "外食費"
	updated, err := repo.UpdateCategory(ctx, cat.ID, categoryUC.UpdateInput{Name: &newName})
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)

	// Delete
	err = repo.DeleteCategory(ctx, cat.ID)
	require.NoError(t, err)

	cats, err = repo.ListCategories(ctx)
	require.NoError(t, err)
	assert.Empty(t, cats)
}

func TestCategoryRepository_GetNotFound(t *testing.T) {
	client := setupTestDB(t)
	repo := persistence.NewCategoryRepository(client)
	ctx := context.Background()

	_, err := repo.GetCategory(ctx, uuid.New())
	assert.Error(t, err)
}

func TestCategoryRepository_Rules(t *testing.T) {
	client := setupTestDB(t)
	repo := persistence.NewCategoryRepository(client)
	ctx := context.Background()

	// カテゴリを作成
	cat, err := repo.CreateCategory(ctx, categoryUC.CreateInput{
		Name:  "食費",
		Color: "#ff0000",
		Icon:  "🍜",
		Type:  entity.CategoryTypeExpense,
	})
	require.NoError(t, err)

	// ルール作成
	rule, err := repo.CreateRule(ctx, categoryUC.CreateRuleInput{
		CategoryID: cat.ID,
		Keyword:    "スーパー",
		Priority:   10,
	})
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, rule.ID)
	assert.Equal(t, "スーパー", rule.Keyword)

	// ルール一覧
	rules, err := repo.ListAllRules(ctx)
	require.NoError(t, err)
	require.Len(t, rules, 1)

	// ルール更新
	newKw := "コンビニ"
	updated, err := repo.UpdateRule(ctx, rule.ID, categoryUC.UpdateRuleInput{Keyword: &newKw})
	require.NoError(t, err)
	assert.Equal(t, newKw, updated.Keyword)

	// ルール削除
	err = repo.DeleteRule(ctx, rule.ID)
	require.NoError(t, err)

	rules, err = repo.ListAllRules(ctx)
	require.NoError(t, err)
	assert.Empty(t, rules)
}
