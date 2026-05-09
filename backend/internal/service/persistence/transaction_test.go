package persistence_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
	"github.com/itk13201/money-rabbit/internal/service/persistence"
	categoryUC "github.com/itk13201/money-rabbit/internal/usecase/category"
	txUC "github.com/itk13201/money-rabbit/internal/usecase/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionRepository_BulkCreate(t *testing.T) {
	client := setupTestDB(t)
	repo := persistence.NewTransactionRepository(client)
	ctx := context.Background()

	inputs := []txUC.CreateInput{
		{Date: time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC), Description: "スーパー", Amount: -1480, ImportFormatID: "smbc_bank"},
		{Date: time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC), Description: "給料振込", Amount: 300000, ImportFormatID: "smbc_bank"},
	}

	txs, err := repo.BulkCreateTransactions(ctx, inputs)
	require.NoError(t, err)
	require.Len(t, txs, 2)
	assert.NotEqual(t, uuid.Nil, txs[0].ID)
	assert.Equal(t, "スーパー", txs[0].Description)
	assert.Equal(t, -1480, txs[0].Amount)
}

func TestTransactionRepository_FindDuplicates(t *testing.T) {
	client := setupTestDB(t)
	repo := persistence.NewTransactionRepository(client)
	ctx := context.Background()

	date := time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC)
	input := txUC.CreateInput{Date: date, Description: "スーパー", Amount: -1480, ImportFormatID: "smbc_bank"}

	// 先に1件登録
	_, err := repo.BulkCreateTransactions(ctx, []txUC.CreateInput{input})
	require.NoError(t, err)

	// 同一内容 → 重複、異なる内容 → 新規
	inputs := []txUC.CreateInput{
		input, // 重複
		{Date: date, Description: "別の取引", Amount: -500, ImportFormatID: "smbc_bank"}, // 新規
	}
	dups, err := repo.FindDuplicates(ctx, inputs)
	require.NoError(t, err)
	assert.True(t, dups[0])
	assert.False(t, dups[1])
}

func TestTransactionRepository_List_MonthFilter(t *testing.T) {
	client := setupTestDB(t)
	repo := persistence.NewTransactionRepository(client)
	ctx := context.Background()

	inputs := []txUC.CreateInput{
		{Date: time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC), Description: "4月", Amount: -100, ImportFormatID: "smbc_bank"},
		{Date: time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), Description: "5月", Amount: -200, ImportFormatID: "smbc_bank"},
	}
	_, err := repo.BulkCreateTransactions(ctx, inputs)
	require.NoError(t, err)

	year, month := 2026, 4
	txs, total, err := repo.ListTransactions(ctx, txUC.ListFilter{Year: &year, Month: &month, PageSize: 50})
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	require.Len(t, txs, 1)
	assert.Equal(t, "4月", txs[0].Description)
}

func TestTransactionRepository_List_CategoryFilter(t *testing.T) {
	client := setupTestDB(t)
	catRepo := persistence.NewCategoryRepository(client)
	txRepo := persistence.NewTransactionRepository(client)
	ctx := context.Background()

	// カテゴリ作成
	cat, err := catRepo.CreateCategory(ctx, categoryUC.CreateInput{
		Name:  "食費",
		Color: "#ff0000",
		Icon:  "🍜",
		Type:  entity.CategoryTypeExpense,
	})
	require.NoError(t, err)

	// 取引作成（1件はカテゴリあり）
	txs, err := txRepo.BulkCreateTransactions(ctx, []txUC.CreateInput{
		{Date: time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC), Description: "スーパー", Amount: -1000, ImportFormatID: "smbc_bank"},
		{Date: time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC), Description: "給料", Amount: 300000, ImportFormatID: "smbc_bank"},
	})
	require.NoError(t, err)

	// スーパーにカテゴリを設定
	_, err = txRepo.UpdateTransactionCategory(ctx, txs[0].ID, &cat.ID)
	require.NoError(t, err)

	// カテゴリでフィルタ
	catID := cat.ID
	result, total, err := txRepo.ListTransactions(ctx, txUC.ListFilter{CategoryID: &catID, PageSize: 50})
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	require.Len(t, result, 1)
	assert.Equal(t, "スーパー", result[0].Description)
}

func TestTransactionRepository_UpdateCategory_Clear(t *testing.T) {
	client := setupTestDB(t)
	catRepo := persistence.NewCategoryRepository(client)
	txRepo := persistence.NewTransactionRepository(client)
	ctx := context.Background()

	cat, err := catRepo.CreateCategory(ctx, categoryUC.CreateInput{
		Name:  "食費",
		Color: "#ff0000",
		Icon:  "🍜",
		Type:  entity.CategoryTypeExpense,
	})
	require.NoError(t, err)

	txs, err := txRepo.BulkCreateTransactions(ctx, []txUC.CreateInput{
		{Date: time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC), Description: "スーパー", Amount: -1000, ImportFormatID: "smbc_bank", CategoryID: &cat.ID},
	})
	require.NoError(t, err)

	// カテゴリをクリア
	updated, err := txRepo.UpdateTransactionCategory(ctx, txs[0].ID, nil)
	require.NoError(t, err)
	assert.Nil(t, updated.CategoryID)
}
