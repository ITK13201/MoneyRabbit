package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
	categoryUC "github.com/itk13201/money-rabbit/internal/usecase/category"
)

// categoryUsecase is the interface the handler depends on.
type categoryUsecase interface {
	List(ctx context.Context) ([]*entity.Category, error)
	Get(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	Create(ctx context.Context, input categoryUC.CreateInput) (*entity.Category, error)
	Update(ctx context.Context, id uuid.UUID, input categoryUC.UpdateInput) (*entity.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CreateRule(ctx context.Context, input categoryUC.CreateRuleInput) (*entity.CategoryRule, error)
	UpdateRule(ctx context.Context, id uuid.UUID, input categoryUC.UpdateRuleInput) (*entity.CategoryRule, error)
	DeleteRule(ctx context.Context, id uuid.UUID) error
}

// CategoryHandler handles /api/categories and /api/category-rules endpoints.
type CategoryHandler struct {
	uc categoryUsecase
}

func NewCategoryHandler(uc categoryUsecase) *CategoryHandler {
	return &CategoryHandler{uc: uc}
}

// --- Category CRUD ---

// List godoc
// @Summary     カテゴリ一覧
// @Tags        categories
// @Produce     json
// @Success     200  {array}   entity.Category
// @Failure     500  {object}  map[string]string
// @Router      /categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	cats, err := h.uc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cats)
}

// Get godoc
// @Summary     カテゴリ取得
// @Tags        categories
// @Produce     json
// @Param       id   path      string  true  "Category UUID"
// @Success     200  {object}  entity.Category
// @Failure     400  {object}  map[string]string
// @Failure     404  {object}  map[string]string
// @Router      /categories/{id} [get]
func (h *CategoryHandler) Get(c *gin.Context) {
	id, err := parseUUID(c, "id")
	if err != nil {
		return
	}
	cat, err := h.uc.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	c.JSON(http.StatusOK, cat)
}

type createCategoryRequest struct {
	Name      string              `json:"name" binding:"required"`
	Color     string              `json:"color" binding:"required"`
	Icon      string              `json:"icon" binding:"required"`
	Type      entity.CategoryType `json:"type" binding:"required,oneof=income expense both"`
	SortOrder int                 `json:"sort_order"`
}

// Create godoc
// @Summary     カテゴリ作成
// @Tags        categories
// @Accept      json
// @Produce     json
// @Param       body  body      createCategoryRequest  true  "Category"
// @Success     201   {object}  entity.Category
// @Failure     400   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req createCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat, err := h.uc.Create(c.Request.Context(), categoryUC.CreateInput{
		Name:      req.Name,
		Color:     req.Color,
		Icon:      req.Icon,
		Type:      req.Type,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

type updateCategoryRequest struct {
	Name      *string              `json:"name"`
	Color     *string              `json:"color"`
	Icon      *string              `json:"icon"`
	Type      *entity.CategoryType `json:"type"`
	SortOrder *int                 `json:"sort_order"`
}

// Update godoc
// @Summary     カテゴリ更新
// @Tags        categories
// @Accept      json
// @Produce     json
// @Param       id    path      string                 true  "Category UUID"
// @Param       body  body      updateCategoryRequest  true  "Fields to update"
// @Success     200   {object}  entity.Category
// @Failure     400   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := parseUUID(c, "id")
	if err != nil {
		return
	}
	var req updateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat, err := h.uc.Update(c.Request.Context(), id, categoryUC.UpdateInput{
		Name:      req.Name,
		Color:     req.Color,
		Icon:      req.Icon,
		Type:      req.Type,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cat)
}

// Delete godoc
// @Summary     カテゴリ削除
// @Tags        categories
// @Param       id  path  string  true  "Category UUID"
// @Success     204
// @Failure     400  {object}  map[string]string
// @Failure     500  {object}  map[string]string
// @Router      /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := parseUUID(c, "id")
	if err != nil {
		return
	}
	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// --- Category Rule CRUD ---

type createRuleRequest struct {
	CategoryID string `json:"category_id" binding:"required,uuid"`
	Keyword    string `json:"keyword" binding:"required"`
	Priority   int    `json:"priority"`
}

// CreateRule godoc
// @Summary     キーワードルール作成
// @Tags        category-rules
// @Accept      json
// @Produce     json
// @Param       body  body      createRuleRequest   true  "Rule"
// @Success     201   {object}  entity.CategoryRule
// @Failure     400   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /category-rules [post]
func (h *CategoryHandler) CreateRule(c *gin.Context) {
	var req createRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	catID, _ := uuid.Parse(req.CategoryID)
	rule, err := h.uc.CreateRule(c.Request.Context(), categoryUC.CreateRuleInput{
		CategoryID: catID,
		Keyword:    req.Keyword,
		Priority:   req.Priority,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rule)
}

type updateRuleRequest struct {
	Keyword  *string `json:"keyword"`
	Priority *int    `json:"priority"`
}

// UpdateRule godoc
// @Summary     キーワードルール更新
// @Tags        category-rules
// @Accept      json
// @Produce     json
// @Param       id    path      string             true  "Rule UUID"
// @Param       body  body      updateRuleRequest  true  "Fields to update"
// @Success     200   {object}  entity.CategoryRule
// @Failure     400   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /category-rules/{id} [put]
func (h *CategoryHandler) UpdateRule(c *gin.Context) {
	id, err := parseUUID(c, "id")
	if err != nil {
		return
	}
	var req updateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rule, err := h.uc.UpdateRule(c.Request.Context(), id, categoryUC.UpdateRuleInput{
		Keyword:  req.Keyword,
		Priority: req.Priority,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rule)
}

// DeleteRule godoc
// @Summary     キーワードルール削除
// @Tags        category-rules
// @Param       id  path  string  true  "Rule UUID"
// @Success     204
// @Failure     400  {object}  map[string]string
// @Failure     500  {object}  map[string]string
// @Router      /category-rules/{id} [delete]
func (h *CategoryHandler) DeleteRule(c *gin.Context) {
	id, err := parseUUID(c, "id")
	if err != nil {
		return
	}
	if err := h.uc.DeleteRule(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// parseUUID extracts and validates a UUID path parameter, writing a 400 response on error.
func parseUUID(c *gin.Context, param string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param(param))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid: " + param})
		return uuid.UUID{}, err
	}
	return id, nil
}
