package transaction

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
)

// UpdateUsecase handles full transaction updates.
type UpdateUsecase struct {
	repo Repository
}

func NewUpdateUsecase(repo Repository) *UpdateUsecase {
	return &UpdateUsecase{repo: repo}
}

func (u *UpdateUsecase) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*entity.Transaction, error) {
	slog.InfoContext(ctx, "transactionRepo.UpdateTransaction started",
		slog.Group("extra", "id", id),
	)
	tx, err := u.repo.UpdateTransaction(ctx, id, input)
	if err != nil {
		slog.ErrorContext(ctx, "transactionRepo.UpdateTransaction failed",
			slog.Group("extra", "id", id, "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "transactionRepo.UpdateTransaction finished",
		slog.Group("extra", "id", id),
	)
	return tx, nil
}
