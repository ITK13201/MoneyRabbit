package csv

// ImportFormatID identifies a supported CSV import format.
type ImportFormatID string

const (
	ImportFormatSMBCBank ImportFormatID = "smbc_bank" // 三井住友銀行 入出金明細
	ImportFormatSMBCCard ImportFormatID = "smbc_card" // SMBCカード 利用明細
)

// ImportType is the category of the import source.
type ImportType string

const (
	ImportTypeBank       ImportType = "bank_account"
	ImportTypeCreditCard ImportType = "credit_card"
)

// ImportFormat describes how to parse a specific bank/card CSV file.
type ImportFormat struct {
	ID          ImportFormatID
	Name        string
	ImportType  ImportType
	Encoding    string // "Shift_JIS" | "UTF-8"
	SkipRows    int
	ColDate     int
	ColDesc     int
	ColAmount   int // credit_card: amount column (-1 = unused)
	ColWithdraw int // bank_account: withdrawal column (-1 = unused)
	ColDeposit  int // bank_account: deposit column (-1 = unused)
	ColBalance  int // balance column (-1 = unused)
	DateFormat  string // Go time layout
}

// Formats is the registry of all supported import formats.
var Formats = map[ImportFormatID]ImportFormat{
	ImportFormatSMBCBank: {
		ID:          ImportFormatSMBCBank,
		Name:        "三井住友銀行",
		ImportType:  ImportTypeBank,
		Encoding:    "Shift_JIS",
		SkipRows:    1,
		ColDate:     0,
		ColDesc:     3,
		ColAmount:   -1,
		ColWithdraw: 1,
		ColDeposit:  2,
		ColBalance:  4,
		DateFormat:  "2006/1/2",
	},
	ImportFormatSMBCCard: {
		ID:          ImportFormatSMBCCard,
		Name:        "SMBCカード",
		ImportType:  ImportTypeCreditCard,
		Encoding:    "Shift_JIS",
		SkipRows:    1,
		ColDate:     0,
		ColDesc:     1,
		ColAmount:   2,
		ColWithdraw: -1,
		ColDeposit:  -1,
		ColBalance:  -1,
		DateFormat:  "2006/01/02",
	},
}
