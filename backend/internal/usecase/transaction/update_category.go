package transaction

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
)

// UpdateCategoryUsecase handles manual category assignment.
type UpdateCategoryUsecase struct {
	repo Repository
}

func NewUpdateCategoryUsecase(repo Repository) *UpdateCategoryUsecase {
	return &UpdateCategoryUsecase{repo: repo}
}

// UpdateCategory sets or clears the category of a transaction.
func (u *UpdateCategoryUsecase) UpdateCategory(ctx context.Context, id uuid.UUID, categoryID *uuid.UUID) (*entity.Transaction, error) {
	slog.InfoContext(ctx, "transactionRepo.UpdateTransactionCategory started",
		slog.Group("extra", "transaction_id", id, "category_id", categoryID),
	)
	tx, err := u.repo.UpdateTransactionCategory(ctx, id, categoryID)
	if err != nil {
		slog.ErrorContext(ctx, "transactionRepo.UpdateTransactionCategory failed",
			slog.Group("extra", "transaction_id", id, "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "transactionRepo.UpdateTransactionCategory finished",
		slog.Group("extra", "transaction_id", id),
	)
	return tx, nil
}
