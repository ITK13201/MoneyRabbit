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

type mockUpdateCategoryUsecase struct{}

func (m *mockUpdateCategoryUsecase) UpdateCategory(_ context.Context, id uuid.UUID, catID *uuid.UUID) (*entity.Transaction, error) {
	return &entity.Transaction{ID: id, CategoryID: catID}, nil
}

type mockDeleteUsecase struct{}

func (m *mockDeleteUsecase) Delete(_ context.Context, _ uuid.UUID) error {
	return nil
}

// --- helpers ---

type mockCreateUsecase struct{}

func (m *mockCreateUsecase) Create(_ context.Context, _ txUC.CreateManualInput) (*entity.Transaction, error) {
	return &entity.Transaction{ID: uuid.New()}, nil
}

type mockUpdateUsecase struct{}

func (m *mockUpdateUsecase) Update(_ context.Context, id uuid.UUID, _ txUC.UpdateInput) (*entity.Transaction, error) {
	return &entity.Transaction{ID: id}, nil
}

func newTransactionRouter(listUC *mockListUsecase, updateUC *mockUpdateCategoryUsecase) *gin.Engine {
	r := gin.New()
	h := handler.NewTransactionHandler(listUC, updateUC, &mockDeleteUsecase{}, &mockCreateUsecase{}, &mockUpdateUsecase{})
	r.GET("/api/transactions", h.List)
	r.POST("/api/transactions", h.Create)
	r.PATCH("/api/transactions/:id/category", h.UpdateCategory)
	r.PUT("/api/transactions/:id", h.Update)
	r.DELETE("/api/transactions/:id", h.Delete)
	return r
}

// --- tests ---

func TestTransactionHandler_List(t *testing.T) {
	now := time.Now()
	fmtID := "smbc_bank"
	txs := []*entity.Transaction{
		{ID: uuid.New(), Date: now, Description: "スーパー", Amount: -1000, ImportFormatID: &fmtID},
		{ID: uuid.New(), Date: now, Description: "給料", Amount: 300000, ImportFormatID: &fmtID},
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
