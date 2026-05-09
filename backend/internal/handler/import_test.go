package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itk13201/money-rabbit/internal/handler"
	txUC "github.com/itk13201/money-rabbit/internal/usecase/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- mock ---

type mockImportUsecase struct {
	result *txUC.ImportResult
}

func (m *mockImportUsecase) Confirm(_ context.Context, _ []txUC.CreateInput) (*txUC.ImportResult, error) {
	return m.result, nil
}

// --- helpers ---

func newImportRouter(uc *mockImportUsecase) *gin.Engine {
	r := gin.New()
	h := handler.NewImportHandler(uc)
	r.GET("/api/import-formats", handler.ListFormats)
	r.POST("/api/import/preview", handler.Preview)
	r.POST("/api/import/confirm", h.Confirm)
	return r
}

// multipartRequest builds a multipart/form-data request with a file field.
func multipartRequest(formatID, csvContent string) *http.Request {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("import_format_id", formatID)
	fw, _ := w.CreateFormFile("file", "test.csv")
	_, _ = fw.Write([]byte(csvContent))
	w.Close()

	req, _ := http.NewRequest(http.MethodPost, "/api/import/preview", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

// --- tests ---

func TestImportHandler_ListFormats(t *testing.T) {
	r := newImportRouter(&mockImportUsecase{})

	w := doRequest(r, http.MethodGet, "/api/import-formats", nil)
	require.Equal(t, http.StatusOK, w.Code)

	var formats []map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &formats))
	assert.NotEmpty(t, formats)
	for _, f := range formats {
		assert.NotEmpty(t, f["id"])
		assert.NotEmpty(t, f["name"])
		assert.NotEmpty(t, f["import_type"])
	}
}

func TestImportHandler_Preview_MissingFormatID(t *testing.T) {
	r := newImportRouter(&mockImportUsecase{})

	var buf bytes.Buffer
	w2 := multipart.NewWriter(&buf)
	w2.Close()
	req, _ := http.NewRequest(http.MethodPost, "/api/import/preview", &buf)
	req.Header.Set("Content-Type", w2.FormDataContentType())

	rec := doRawRequest(r, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestImportHandler_Preview_UnsupportedFormatID(t *testing.T) {
	r := newImportRouter(&mockImportUsecase{})

	req := multipartRequest("unknown_format", "header\n2026/4/2,100,,desc,1000")
	rec := doRawRequest(r, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestImportHandler_Preview_MissingFile(t *testing.T) {
	r := newImportRouter(&mockImportUsecase{})

	var buf bytes.Buffer
	w2 := multipart.NewWriter(&buf)
	_ = w2.WriteField("import_format_id", "smbc_bank")
	w2.Close()
	req, _ := http.NewRequest(http.MethodPost, "/api/import/preview", &buf)
	req.Header.Set("Content-Type", w2.FormDataContentType())

	rec := doRawRequest(r, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestImportHandler_Confirm_OK(t *testing.T) {
	uc := &mockImportUsecase{
		result: &txUC.ImportResult{Imported: 2, Skipped: 0},
	}
	r := newImportRouter(uc)

	now := time.Now().Format(time.DateOnly)
	body := map[string]any{
		"import_format_id": "smbc_bank",
		"transactions": []map[string]any{
			{"date": now, "description": "スーパー", "amount": -1000},
			{"date": now, "description": "給料", "amount": 300000},
		},
	}
	w := doRequest(r, http.MethodPost, "/api/import/confirm", body)
	require.Equal(t, http.StatusOK, w.Code)

	var result txUC.ImportResult
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 2, result.Imported)
	assert.Equal(t, 0, result.Skipped)
}

func TestImportHandler_Confirm_UnsupportedFormat(t *testing.T) {
	r := newImportRouter(&mockImportUsecase{})

	now := time.Now().Format(time.DateOnly)
	body := map[string]any{
		"import_format_id": "unknown_bank",
		"transactions": []map[string]any{
			{"date": now, "description": "test", "amount": -100},
		},
	}
	w := doRequest(r, http.MethodPost, "/api/import/confirm", body)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestImportHandler_Confirm_InvalidDate(t *testing.T) {
	r := newImportRouter(&mockImportUsecase{})

	body := map[string]any{
		"import_format_id": "smbc_bank",
		"transactions": []map[string]any{
			{"date": "not-a-date", "description": "test", "amount": -100},
		},
	}
	w := doRequest(r, http.MethodPost, "/api/import/confirm", body)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestImportHandler_Preview_ValidCSV(t *testing.T) {
	r := newImportRouter(&mockImportUsecase{})

	// smbc_bank: ColDate=0, ColWithdraw=1, ColDeposit=2, ColDesc=3, ColBalance=4
	// ASCII-only content is valid Shift-JIS
	csv := strings.Join([]string{
		"header",
		"2026/4/2,1480,,purchase,4600092",
	}, "\n")
	req := multipartRequest("smbc_bank", csv)
	rec := doRawRequest(r, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Rows        []map[string]any `json:"rows"`
		SkippedRows int              `json:"skipped_rows"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Rows, 1)
	assert.Equal(t, "purchase", resp.Rows[0]["description"])
	assert.Equal(t, float64(-1480), resp.Rows[0]["amount"])
}
