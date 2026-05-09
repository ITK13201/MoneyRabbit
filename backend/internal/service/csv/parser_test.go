package csv

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testFormat returns a UTF-8 bank_account format for testing parse logic.
func testBankFormat() ImportFormat {
	return ImportFormat{
		ID:          "test_bank",
		Name:        "テスト銀行",
		ImportType:  ImportTypeBank,
		Encoding:    "UTF-8",
		SkipRows:    1,
		ColDate:     0,
		ColDesc:     3,
		ColWithdraw: 1,
		ColDeposit:  2,
		ColBalance:  4,
		ColAmount:   -1,
		DateFormat:  "2006/1/2",
	}
}

func testCardFormat() ImportFormat {
	return ImportFormat{
		ID:          "test_card",
		Name:        "テストカード",
		ImportType:  ImportTypeCreditCard,
		Encoding:    "UTF-8",
		SkipRows:    1,
		ColDate:     0,
		ColDesc:     1,
		ColAmount:   2,
		ColWithdraw: -1,
		ColDeposit:  -1,
		ColBalance:  -1,
		DateFormat:  "2006/01/02",
	}
}

func TestParse_BankAccount(t *testing.T) {
	csv := strings.Join([]string{
		"年月日,お引出し,お預入れ,お取り扱い内容,残高",
		"2026/4/2,1480,,スーパーマーケット,4600092",
		"2026/4/24,,359562,給料振込,4957768",
	}, "\n")

	result, err := Parse(strings.NewReader(csv), testBankFormat())
	require.NoError(t, err)
	require.Len(t, result.Rows, 2)
	assert.Equal(t, 0, result.SkippedRows)

	assert.Equal(t, "2026-04-02", result.Rows[0].Date.Format("2006-01-02"))
	assert.Equal(t, "スーパーマーケット", result.Rows[0].Description)
	assert.Equal(t, -1480, result.Rows[0].Amount)

	assert.Equal(t, "2026-04-24", result.Rows[1].Date.Format("2006-01-02"))
	assert.Equal(t, "給料振込", result.Rows[1].Description)
	assert.Equal(t, 359562, result.Rows[1].Amount)
}

func TestParse_CreditCard(t *testing.T) {
	csv := strings.Join([]string{
		"氏名ヘッダー行",
		"2026/03/06,GOOGLE PLAY,1280",
		"2026/03/09,AMAZON,3500",
	}, "\n")

	result, err := Parse(strings.NewReader(csv), testCardFormat())
	require.NoError(t, err)
	require.Len(t, result.Rows, 2)

	assert.Equal(t, -1280, result.Rows[0].Amount)
	assert.Equal(t, "GOOGLE PLAY", result.Rows[0].Description)
	assert.Equal(t, -3500, result.Rows[1].Amount)
}

func TestParse_SkipsEmptyDate(t *testing.T) {
	csv := strings.Join([]string{
		"年月日,お引出し,お預入れ,お取り扱い内容,残高",
		"2026/4/2,1480,,スーパー,1000",
		",,,,,187313", // 合計行（日付なし）
	}, "\n")

	result, err := Parse(strings.NewReader(csv), testBankFormat())
	require.NoError(t, err)
	assert.Len(t, result.Rows, 1)
	assert.Equal(t, 1, result.SkippedRows)
}

func TestParse_SkipsZeroAmount(t *testing.T) {
	csv := strings.Join([]string{
		"年月日,お引出し,お預入れ,お取り扱い内容,残高",
		"2026/4/2,,,メモのみ,5000", // 出金も入金も空
	}, "\n")

	result, err := Parse(strings.NewReader(csv), testBankFormat())
	require.NoError(t, err)
	assert.Len(t, result.Rows, 0)
	assert.Equal(t, 1, result.SkippedRows)
}

func TestNormalizeNumber_FullWidth(t *testing.T) {
	assert.Equal(t, "1280", normalizeNumber("１２８０"))
	assert.Equal(t, "1000000", normalizeNumber("1,000,000"))
	assert.Equal(t, "500", normalizeNumber("５００"))
}
