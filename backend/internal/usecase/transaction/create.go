package transaction

import (
	"context"
	"log/slog"

	"github.com/itk13201/money-rabbit/internal/domain/entity"
)

// CreateUsecase handles manual transaction creation.
type CreateUsecase struct {
	repo Repository
}

func NewCreateUsecase(repo Repository) *CreateUsecase {
	return &CreateUsecase{repo: repo}
}

func (u *CreateUsecase) Create(ctx context.Context, input CreateManualInput) (*entity.Transaction, error) {
	slog.InfoContext(ctx, "transactionRepo.CreateManualTransaction started",
		slog.Group("extra", "description", input.Description, "amount", input.Amount),
	)
	tx, err := u.repo.CreateManualTransaction(ctx, input)
	if err != nil {
		slog.ErrorContext(ctx, "transactionRepo.CreateManualTransaction failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "transactionRepo.CreateManualTransaction finished",
		slog.Group("extra", "id", tx.ID),
	)
	return tx, nil
}
