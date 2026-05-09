package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
	"github.com/itk13201/money-rabbit/internal/handler"
	txUC "github.com/itk13201/money-rabbit/internal/usecase/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- mocks ---

type mockListUsecase struct {
	result *txUC.ListResult
}

func (m *mockListUsecase) List(_ context.Context, _ txUC.ListFilter) (*txUC.ListResult, error) {
	return m.result, nil
}

type mockUpdateCategoryUsecase struct {
	updated *entity.Transaction
}

func (m *mockUpdateCategoryUsecase) UpdateCategory(_ context.Context, id uuid.UUID, catID *uuid.UUID) (*entity.Transaction, error) {
	return &entity.Transaction{ID: id, CategoryID: catID}, nil
}

// --- helpers ---

func newTransactionRouter(listUC *mockListUsecase, updateUC *mockUpdateCategoryUsecase) *gin.Engine {
	r := gin.New()
	h := handler.NewTransactionHandler(listUC, updateUC)
	r.GET("/api/transactions", h.List)
	r.PATCH("/api/transactions/:id/category", h.UpdateCategory)
	return r
}

// --- tests ---

func TestTransactionHandler_List(t *testing.T) {
	now := time.Now()
	txs := []*entity.Transaction{
		{ID: uuid.New(), Date: now, Description: "スーパー", Amount: -1000, ImportFormatID: "smbc_bank"},
		{ID: uuid.New(), Date: now, Description: "給料",   Amount: 300000, ImportFormatID: "smbc_bank"},
	}
	uc := &mockListUsecase{result: &txUC.ListResult{Transactions: txs, Total: 2}}
	r := newTransactionRouter(uc, &mockUpdateCategoryUsecase{})

	w := doRequest(r, http.MethodGet, "/api/transactions", nil)
	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Transactions []*entity.Transaction `json:"transactions"`
		Total        int                   `json:"total"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 2, resp.Total)
	assert.Len(t, resp.Transactions, 2)
}

func TestTransactionHandler_UpdateCategory(t *testing.T) {
	txID := uuid.New()
	catID := uuid.New()
	r := newTransactionRouter(&mockListUsecase{result: &txUC.ListResult{}}, &mockUpdateCategoryUsecase{})

	w := doRequest(r, http.MethodPatch, "/api/transactions/"+txID.String()+"/category", map[string]any{
		"category_id": catID.String(),
	})
	require.Equal(t, http.StatusOK, w.Code)

	var tx entity.Transaction
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &tx))
	assert.Equal(t, txID, tx.ID)
	require.NotNil(t, tx.CategoryID)
	assert.Equal(t, catID, *tx.CategoryID)
}
