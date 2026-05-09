package transaction_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/internal/usecase/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateCategoryUsecase_SetCategory(t *testing.T) {
	repo := &mockTxRepo{}
	uc := transaction.NewUpdateCategoryUsecase(repo)

	txID := uuid.New()
	catID := uuid.New()
	result, err := uc.UpdateCategory(context.Background(), txID, &catID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, txID, result.ID)
	require.NotNil(t, result.CategoryID)
	assert.Equal(t, catID, *result.CategoryID)
}

func TestUpdateCategoryUsecase_ClearCategory(t *testing.T) {
	repo := &mockTxRepo{}
	uc := transaction.NewUpdateCategoryUsecase(repo)

	txID := uuid.New()
	result, err := uc.UpdateCategory(context.Background(), txID, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, txID, result.ID)
	assert.Nil(t, result.CategoryID)
}
