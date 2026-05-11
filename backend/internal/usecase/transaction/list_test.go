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

// mockListTxRepo supports configurable ListTransactions results.
type mockListTxRepo struct {
	mockTxRepo
	txs   []*entity.Transaction
	total int
}

func (m *mockListTxRepo) ListTransactions(_ context.Context, _ transaction.ListFilter) ([]*entity.Transaction, int, error) {
	return m.txs, m.total, nil
}

func TestListUsecase_ReturnsPaginatedResults(t *testing.T) {
	now := time.Now()
	fmtID := "smbc_bank"
	txs := []*entity.Transaction{
		{ID: uuid.New(), Date: now, Description: "スーパー", Amount: -1000, ImportFormatID: &fmtID},
		{ID: uuid.New(), Date: now, Description: "給料", Amount: 300000, ImportFormatID: &fmtID},
	}
	repo := &mockListTxRepo{txs: txs, total: 2}
	uc := transaction.NewListUsecase(repo)

	result, err := uc.List(context.Background(), transaction.ListFilter{Page: 0, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.Transactions, 2)
	assert.Equal(t, "スーパー", result.Transactions[0].Description)
}

func TestListUsecase_EmptyResult(t *testing.T) {
	repo := &mockListTxRepo{txs: nil, total: 0}
	uc := transaction.NewListUsecase(repo)

	result, err := uc.List(context.Background(), transaction.ListFilter{Page: 0, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, 0, result.Total)
	assert.Empty(t, result.Transactions)
}
