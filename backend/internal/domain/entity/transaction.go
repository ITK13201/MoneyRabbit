package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID             uuid.UUID  `json:"id"`
	Date           time.Time  `json:"date"`
	Description    string     `json:"description"`
	Amount         int        `json:"amount"`
	ImportFormatID string     `json:"import_format_id"`
	ImportedAt     time.Time  `json:"imported_at"`
	CategoryID     *uuid.UUID `json:"category_id"`
	Category       *Category  `json:"category,omitempty"`
}
