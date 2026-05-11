package transaction

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
)

// Repository is the interface that the persistence layer implements.
type Repository interface {
	FindDuplicates(ctx context.Context, inputs []CreateInput) (map[int]bool, error)
	BulkCreateTransactions(ctx context.Context, inputs []CreateInput) ([]*entity.Transaction, error)
	CreateManualTransaction(ctx context.Context, input CreateManualInput) (*entity.Transaction, error)
	ListTransactions(ctx context.Context, filter ListFilter) ([]*entity.Transaction, int, error)
	UpdateTransactionCategory(ctx context.Context, id uuid.UUID, categoryID *uuid.UUID) (*entity.Transaction, error)
	UpdateTransaction(ctx context.Context, id uuid.UUID, input UpdateInput) (*entity.Transaction, error)
	DeleteTransaction(ctx context.Context, id uuid.UUID) error
}

// Classifier is the interface for automatic category classification.
type Classifier interface {
	Classify(ctx context.Context, descriptions []string, categories []*entity.Category) (map[string]*uuid.UUID, error)
}

// CreateInput represents a single transaction to be inserted.
// CategoryID is set during classification and may be nil (uncategorized).
type CreateInput struct {
	Date           time.Time
	Description    string
	Amount         int
	ImportFormatID string
	CategoryID     *uuid.UUID
}

type CreateManualInput struct {
	Date        time.Time
	Description string
	Amount      int
	CategoryID  *uuid.UUID
}

type UpdateInput struct {
	Date        time.Time
	Description string
	Amount      int
	CategoryID  *uuid.UUID
}

type ListFilter struct {
	Year       *int
	Month      *int
	CategoryID *uuid.UUID
	Page       int
	PageSize   int
}
