package transaction_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
	"github.com/itk13201/money-rabbit/internal/usecase/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- mocks ---

type mockTxRepo struct {
	created []*entity.Transaction
}

func (m *mockTxRepo) FindDuplicates(_ context.Context, inputs []transaction.CreateInput) (map[int]bool, error) {
	return map[int]bool{}, nil
}

func (m *mockTxRepo) BulkCreateTransactions(_ context.Context, inputs []transaction.CreateInput) ([]*entity.Transaction, error) {
	txs := make([]*entity.Transaction, len(inputs))
	for i, inp := range inputs {
		fmtID := inp.ImportFormatID
		txs[i] = &entity.Transaction{
			ID:             uuid.New(),
			Date:           inp.Date,
			Description:    inp.Description,
			Amount:         inp.Amount,
			ImportFormatID: &fmtID,
			CategoryID:     inp.CategoryID,
		}
	}
	m.created = txs
	return txs, nil
}

func (m *mockTxRepo) CreateManualTransaction(_ context.Context, input transaction.CreateManualInput) (*entity.Transaction, error) {
	tx := &entity.Transaction{
		ID:          uuid.New(),
		Date:        input.Date,
		Description: input.Description,
		Amount:      input.Amount,
		CategoryID:  input.CategoryID,
	}
	m.created = append(m.created, tx)
	return tx, nil
}

func (m *mockTxRepo) UpdateTransaction(_ context.Context, id uuid.UUID, input transaction.UpdateInput) (*entity.Transaction, error) {
	return &entity.Transaction{
		ID:          id,
		Date:        input.Date,
		Description: input.Description,
		Amount:      input.Amount,
		CategoryID:  input.CategoryID,
	}, nil
}

func (m *mockTxRepo) ListTransactions(_ context.Context, _ transaction.ListFilter) ([]*entity.Transaction, int, error) {
	return nil, 0, nil
}

func (m *mockTxRepo) UpdateTransactionCategory(_ context.Context, id uuid.UUID, catID *uuid.UUID) (*entity.Transaction, error) {
	return &entity.Transaction{ID: id, CategoryID: catID}, nil
}

func (m *mockTxRepo) DeleteTransaction(_ context.Context, _ uuid.UUID) error {
	return nil
}

type mockCatRepo struct {
	categories []*entity.Category
	rules      []*entity.CategoryRule
}

func (m *mockCatRepo) ListCategories(_ context.Context) ([]*entity.Category, error) {
	return m.categories, nil
}
func (m *mockCatRepo) ListAllRules(_ context.Context) ([]*entity.CategoryRule, error) {
	return m.rules, nil
}

type mockClassifier struct {
	results map[string]*uuid.UUID
}

func (m *mockClassifier) Classify(_ context.Context, descs []string, _ []*entity.Category) (map[string]*uuid.UUID, error) {
	return m.results, nil
}

// --- tests ---

func TestConfirm_SavesTransactions(t *testing.T) {
	txRepo := &mockTxRepo{}
	catRepo := &mockCatRepo{}
	uc := transaction.NewImportUsecase(txRepo, catRepo, nil)

	inputs := []transaction.CreateInput{
		{Date: time.Now(), Description: "スーパー", Amount: -1000, ImportFormatID: "smbc_bank"},
		{Date: time.Now(), Description: "給料", Amount: 300000, ImportFormatID: "smbc_bank"},
	}

	result, err := uc.Confirm(context.Background(), inputs)
	require.NoError(t, err)
	assert.Equal(t, 2, result.Imported)
	assert.Equal(t, 0, result.Skipped)
	assert.Empty(t, result.Errors)
}

func TestConfirm_KeywordRuleAssignsCategory(t *testing.T) {
	catID := uuid.New()
	txRepo := &mockTxRepo{}
	catRepo := &mockCatRepo{
		categories: []*entity.Category{
			{ID: catID, Name: "食費"},
		},
		rules: []*entity.CategoryRule{
			{ID: uuid.New(), Keyword: "スーパー", Priority: 10, CategoryID: catID},
		},
	}
	uc := transaction.NewImportUsecase(txRepo, catRepo, nil)

	inputs := []transaction.CreateInput{
		{Date: time.Now(), Description: "スーパーマーケット購入", Amount: -2000, ImportFormatID: "smbc_bank"},
	}

	result, err := uc.Confirm(context.Background(), inputs)
	require.NoError(t, err)
	require.Equal(t, 1, result.Imported)

	// キーワードルールでカテゴリが自動設定されること
	require.Len(t, txRepo.created, 1)
	require.NotNil(t, txRepo.created[0].CategoryID)
	assert.Equal(t, catID, *txRepo.created[0].CategoryID)
}

func TestConfirm_SkipsDuplicates(t *testing.T) {
	dupRepo := &mockTxRepoDup{}
	catRepo := &mockCatRepo{}
	uc := transaction.NewImportUsecase(dupRepo, catRepo, nil)

	inputs := []transaction.CreateInput{
		{Date: time.Now(), Description: "重複", Amount: -500, ImportFormatID: "smbc_bank"},
		{Date: time.Now(), Description: "新規", Amount: -1000, ImportFormatID: "smbc_bank"},
	}

	result, err := uc.Confirm(context.Background(), inputs)
	require.NoError(t, err)
	assert.Equal(t, 1, result.Imported)
	assert.Equal(t, 1, result.Skipped)
}

func TestConfirm_EmptyInputReturnsZero(t *testing.T) {
	uc := transaction.NewImportUsecase(&mockTxRepo{}, &mockCatRepo{}, nil)

	result, err := uc.Confirm(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, 0, result.Imported)
}

// mockTxRepoDup marks index 0 as a duplicate.
type mockTxRepoDup struct {
	mockTxRepo
}

func (m *mockTxRepoDup) FindDuplicates(_ context.Context, inputs []transaction.CreateInput) (map[int]bool, error) {
	result := make(map[int]bool)
	if len(inputs) > 0 {
		result[0] = true // 先頭を重複扱い
	}
	return result, nil
}
