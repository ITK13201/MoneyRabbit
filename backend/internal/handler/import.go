package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	csvService "github.com/itk13201/money-rabbit/internal/service/csv"
	txUC "github.com/itk13201/money-rabbit/internal/usecase/transaction"
)

// importUsecase is the interface the import handler depends on.
type importUsecase interface {
	Confirm(ctx context.Context, inputs []txUC.CreateInput) (*txUC.ImportResult, error)
}

// ImportHandler handles /api/import endpoints.
type ImportHandler struct {
	uc importUsecase
}

func NewImportHandler(uc importUsecase) *ImportHandler {
	return &ImportHandler{uc: uc}
}

// ListFormats returns all supported import formats.
func ListFormats(c *gin.Context) {
	type formatResponse struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		ImportType string `json:"import_type"`
	}
	resp := make([]formatResponse, 0, len(csvService.Formats))
	for _, f := range csvService.Formats {
		resp = append(resp, formatResponse{
			ID:         string(f.ID),
			Name:       f.Name,
			ImportType: string(f.ImportType),
		})
	}
	c.JSON(http.StatusOK, resp)
}

// Preview parses a CSV file and returns rows without saving them.
func Preview(c *gin.Context) {
	formatID := c.PostForm("import_format_id")
	if formatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "import_format_id is required"})
		return
	}

	format, ok := csvService.Formats[csvService.ImportFormatID(formatID)]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported import_format_id"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer f.Close()

	result, err := csvService.Parse(f, format)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "csv parse error: " + err.Error()})
		return
	}

	type rowResponse struct {
		Date        string `json:"date"`
		Description string `json:"description"`
		Amount      int    `json:"amount"`
	}
	rows := make([]rowResponse, len(result.Rows))
	for i, row := range result.Rows {
		rows[i] = rowResponse{
			Date:        row.Date.Format(time.DateOnly),
			Description: row.Description,
			Amount:      row.Amount,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"rows":         rows,
		"skipped_rows": result.SkippedRows,
	})
}

// Confirm saves the user-confirmed transaction rows.
func (h *ImportHandler) Confirm(c *gin.Context) {
	type txRow struct {
		Date           string `json:"date" binding:"required"`
		Description    string `json:"description"`
		Amount         int    `json:"amount"`
	}
	var req struct {
		ImportFormatID string  `json:"import_format_id" binding:"required"`
		Transactions   []txRow `json:"transactions" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, ok := csvService.Formats[csvService.ImportFormatID(req.ImportFormatID)]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported import_format_id"})
		return
	}

	inputs := make([]txUC.CreateInput, 0, len(req.Transactions))
	for _, row := range req.Transactions {
		date, err := time.Parse(time.DateOnly, row.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date: " + row.Date})
			return
		}
		inputs = append(inputs, txUC.CreateInput{
			Date:           date,
			Description:    row.Description,
			Amount:         row.Amount,
			ImportFormatID: req.ImportFormatID,
		})
	}

	result, err := h.uc.Confirm(c.Request.Context(), inputs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
