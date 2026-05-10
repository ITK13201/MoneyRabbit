package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	sumUC "github.com/itk13201/money-rabbit/internal/usecase/summary"
)

// monthlyUsecase is the interface for monthly summary aggregation.
type monthlyUsecase interface {
	Monthly(ctx context.Context, year int) ([]sumUC.MonthSummary, error)
}

// SummaryHandler handles /api/summary endpoints.
type SummaryHandler struct {
	monthlyUC monthlyUsecase
}

func NewSummaryHandler(monthlyUC monthlyUsecase) *SummaryHandler {
	return &SummaryHandler{monthlyUC: monthlyUC}
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
// Monthly returns monthly income/expense aggregation for the given year.
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
