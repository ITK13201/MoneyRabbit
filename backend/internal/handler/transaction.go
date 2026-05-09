package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
	txUC "github.com/itk13201/money-rabbit/internal/usecase/transaction"
)

// listUsecase is the interface for listing transactions.
type listUsecase interface {
	List(ctx context.Context, filter txUC.ListFilter) (*txUC.ListResult, error)
}

// updateCategoryUsecase is the interface for updating transaction categories.
type updateCategoryUsecase interface {
	UpdateCategory(ctx context.Context, id uuid.UUID, categoryID *uuid.UUID) (*entity.Transaction, error)
}

// TransactionHandler handles /api/transactions endpoints.
type TransactionHandler struct {
	listUC   listUsecase
	updateUC updateCategoryUsecase
}

func NewTransactionHandler(listUC listUsecase, updateUC updateCategoryUsecase) *TransactionHandler {
	return &TransactionHandler{
		listUC:   listUC,
		updateUC: updateUC,
	}
}

// transactionListResponse is the response body for List.
type transactionListResponse struct {
	Transactions []*entity.Transaction `json:"transactions"`
	Total        int                   `json:"total"`
}

// List godoc
// @Summary     取引一覧
// @Tags        transactions
// @Produce     json
// @Param       year         query     int     false  "年（例: 2025）"
// @Param       month        query     int     false  "月（1-12）"
// @Param       category_id  query     string  false  "カテゴリID（UUID）"
// @Param       page         query     int     false  "ページ番号（0始まり）"
// @Param       page_size    query     int     false  "ページサイズ（デフォルト: 50）"
// @Success     200  {object}  transactionListResponse
// @Failure     400  {object}  map[string]string
// @Failure     500  {object}  map[string]string
// @Router      /transactions [get]
// List returns a paginated, filtered list of transactions.
func (h *TransactionHandler) List(c *gin.Context) {
	filter := txUC.ListFilter{
		Page:     0,
		PageSize: 50,
	}

	if v := c.Query("year"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
			return
		}
		filter.Year = &n
	}
	if v := c.Query("month"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
			return
		}
		filter.Month = &n
	}
	if v := c.Query("category_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id"})
			return
		}
		filter.CategoryID = &id
	}
	if v := c.Query("page"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
			return
		}
		filter.Page = n
	}
	if v := c.Query("page_size"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page_size"})
			return
		}
		filter.PageSize = n
	}

	result, err := h.listUC.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": result.Transactions,
		"total":        result.Total,
	})
}

// updateTransactionCategoryRequest is the request body for UpdateCategory.
type updateTransactionCategoryRequest struct {
	CategoryID *string `json:"category_id"`
}

// UpdateCategory godoc
// @Summary     取引のカテゴリを変更
// @Tags        transactions
// @Accept      json
// @Produce     json
// @Param       id    path      string                 true  "Transaction UUID"
// @Param       body  body      updateTransactionCategoryRequest  true  "Category ID（nullでクリア）"
// @Success     200   {object}  entity.Transaction
// @Failure     400   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /transactions/{id}/category [patch]
// UpdateCategory sets or clears the category of a transaction.
func (h *TransactionHandler) UpdateCategory(c *gin.Context) {
	id, err := parseUUID(c, "id")
	if err != nil {
		return
	}

	var req updateTransactionCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var catID *uuid.UUID
	if req.CategoryID != nil {
		parsed, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id"})
			return
		}
		catID = &parsed
	}

	tx, err := h.updateUC.UpdateCategory(c.Request.Context(), id, catID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tx)
}
