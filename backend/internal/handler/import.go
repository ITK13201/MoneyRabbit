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

// formatResponse is the response body for ListFormats.
type formatResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ImportType string `json:"import_type"`
}

// rowResponse is a single parsed CSV row.
type rowResponse struct {
	Date        string `json:"date"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
}

// previewResponse is the response body for Preview.
type previewResponse struct {
	Rows        []rowResponse `json:"rows"`
	SkippedRows int           `json:"skipped_rows"`
}

// importResultResponse is the response body for Confirm.
type importResultResponse struct {
	Imported int      `json:"Imported"`
	Skipped  int      `json:"Skipped"`
	Errors   []string `json:"Errors"`
}

// confirmRequest is the request body for Confirm.
type confirmRequest struct {
	ImportFormatID string       `json:"import_format_id" binding:"required"`
	Transactions   []txRowInput `json:"transactions"      binding:"required"`
}

// txRowInput is a single transaction row for confirm.
type txRowInput struct {
	Date        string `json:"date"        binding:"required"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
}

// ListFormats godoc
// @Summary     対応インポートフォーマット一覧
// @Tags        import
// @Produce     json
// @Success     200  {array}  formatResponse
// @Router      /import-formats [get]
// ListFormats returns all supported import formats.
func ListFormats(c *gin.Context) {
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

// Preview godoc
// @Summary     CSVプレビュー（DB書き込みなし）
// @Tags        import
// @Accept      multipart/form-data
// @Produce     json
// @Param       import_format_id  formData  string  true  "フォーマットID（例: smbc_bank）"
// @Param       file              formData  file    true  "CSVファイル"
// @Success     200  {object}  previewResponse
// @Failure     400  {object}  map[string]string
// @Failure     500  {object}  map[string]string
// @Router      /import/preview [post]
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

	rows := make([]rowResponse, len(result.Rows))
	for i, row := range result.Rows {
		rows[i] = rowResponse{
			Date:        row.Date.Format(time.DateOnly),
			Description: row.Description,
			Amount:      row.Amount,
		}
	}

	c.JSON(http.StatusOK, previewResponse{Rows: rows, SkippedRows: result.SkippedRows})
}

// Confirm godoc
// @Summary     CSVインポート確定（分類＋保存）
// @Tags        import
// @Accept      json
// @Produce     json
// @Param       body  body      confirmRequest    true  "インポートデータ"
// @Success     200   {object}  importResultResponse
// @Failure     400   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /import/confirm [post]
// Confirm saves the user-confirmed transaction rows.
func (h *ImportHandler) Confirm(c *gin.Context) {
	var req confirmRequest
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
