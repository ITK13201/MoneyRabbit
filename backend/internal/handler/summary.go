package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	sumUC "github.com/itk13201/money-rabbit/internal/usecase/summary"
)

type monthlyUsecase interface {
	Monthly(ctx context.Context, year int) ([]sumUC.MonthSummary, error)
}

type categoryAnnualUsecase interface {
	CategoryAnnual(ctx context.Context, year int) (*sumUC.CategoryAnnualSummary, error)
}

// SummaryHandler handles /api/summary endpoints.
type SummaryHandler struct {
	monthlyUC        monthlyUsecase
	categoryAnnualUC categoryAnnualUsecase
}

func NewSummaryHandler(monthlyUC monthlyUsecase, categoryAnnualUC categoryAnnualUsecase) *SummaryHandler {
	return &SummaryHandler{monthlyUC: monthlyUC, categoryAnnualUC: categoryAnnualUC}
}

// Monthly godoc
// @Summary     月別収支サマリー
// @Tags        summary
// @Produce     json
// @Param       year  query     int  true  "年（例: 2026）"
// @Success     200   {object}  map[string]any
// @Failure     400   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /summary/monthly [get]
func (h *SummaryHandler) Monthly(c *gin.Context) {
	yearStr := c.Query("year")
	if yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year is required"})
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 || year > 2100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}

	result, err := h.monthlyUC.Monthly(c.Request.Context(), year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result == nil {
		result = []sumUC.MonthSummary{}
	}
	c.JSON(http.StatusOK, gin.H{"months": result})
}

// CategoryAnnual godoc
// @Summary     年間カテゴリ分布
// @Tags        summary
// @Produce     json
// @Param       year  query     int  true  "年（例: 2026）"
// @Success     200   {object}  sumUC.CategoryAnnualSummary
// @Failure     400   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /summary/category-annual [get]
func (h *SummaryHandler) CategoryAnnual(c *gin.Context) {
	yearStr := c.Query("year")
	if yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year is required"})
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 || year > 2100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}

	result, err := h.categoryAnnualUC.CategoryAnnual(c.Request.Context(), year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
