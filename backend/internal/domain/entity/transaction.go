package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID             uuid.UUID
	Date           time.Time
	Description    string
	Amount         int
	ImportFormatID string
	ImportedAt     time.Time
	CategoryID     *uuid.UUID
	Category       *Category
}
