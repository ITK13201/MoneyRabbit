package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// ParsedRow is a single normalized transaction row from a CSV file.
type ParsedRow struct {
	Date        time.Time
	Description string
	Amount      int // positive = income, negative = expense
}

// ParseResult is the result of parsing a CSV file.
type ParseResult struct {
	Rows        []ParsedRow
	SkippedRows int
}

// Parse decodes a CSV file according to the given ImportFormat and returns normalized rows.
func Parse(r io.Reader, format ImportFormat) (*ParseResult, error) {
	decoded, err := decode(r, format.Encoding)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	reader := csv.NewReader(strings.NewReader(decoded))
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1 // 行ごとに列数が異なっても許容する

	all, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("csv parse: %w", err)
	}

	result := &ParseResult{}
	for i, row := range all {
		if i < format.SkipRows {
			continue
		}
		if len(row) <= format.ColDate {
			result.SkippedRows++
			continue
		}

		rawDate := strings.TrimSpace(row[format.ColDate])
		if rawDate == "" {
			result.SkippedRows++
			continue
		}

		date, err := time.Parse(format.DateFormat, rawDate)
		if err != nil {
			result.SkippedRows++
			continue
		}

		description := ""
		if len(row) > format.ColDesc {
			description = strings.TrimSpace(normalizeFullWidth(row[format.ColDesc]))
		}

		amount, skip := parseAmount(row, format)
		if skip {
			result.SkippedRows++
			continue
		}

		result.Rows = append(result.Rows, ParsedRow{
			Date:        date,
			Description: description,
			Amount:      amount,
		})
	}

	return result, nil
}

// parseAmount extracts and normalizes the signed amount from a row based on the format.
// Returns (amount, shouldSkip).
func parseAmount(row []string, format ImportFormat) (int, bool) {
	getCol := func(col int) string {
		if col < 0 || col >= len(row) {
			return ""
		}
		return strings.TrimSpace(normalizeNumber(row[col]))
	}

	switch format.ImportType {
	case ImportTypeCreditCard:
		raw := getCol(format.ColAmount)
		if raw == "" {
			return 0, true
		}
		v, err := strconv.Atoi(raw)
		if err != nil || v == 0 {
			return 0, true
		}
		return -v, false // credit card purchases are expenses

	case ImportTypeBank:
		withdraw := getCol(format.ColWithdraw)
		deposit := getCol(format.ColDeposit)

		if deposit != "" {
			v, err := strconv.Atoi(deposit)
			if err == nil && v > 0 {
				return v, false
			}
		}
		if withdraw != "" {
			v, err := strconv.Atoi(withdraw)
			if err == nil && v > 0 {
				return -v, false
			}
		}
		return 0, true

	default:
		return 0, true
	}
}

// decode converts the reader to a UTF-8 string using the specified encoding.
func decode(r io.Reader, encoding string) (string, error) {
	switch encoding {
	case "Shift_JIS":
		decoded, err := io.ReadAll(transform.NewReader(r, japanese.ShiftJIS.NewDecoder()))
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	default:
		b, err := io.ReadAll(r)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
}

// normalizeNumber converts full-width digits/commas to ASCII and removes commas.
func normalizeNumber(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '０' && r <= '９' {
			b.WriteRune(r - 0xFEE0)
		} else if r == '，' || r == ',' {
			// skip thousand separators
		} else if unicode.IsDigit(r) || r == '-' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// normalizeFullWidth converts full-width ASCII characters to half-width.
func normalizeFullWidth(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '！' && r <= '～' {
			b.WriteRune(r - 0xFEE0)
		} else if r == '　' {
			b.WriteRune(' ')
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
