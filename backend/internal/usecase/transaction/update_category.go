package transaction

import (
	"context"

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
	return u.repo.UpdateTransactionCategory(ctx, id, categoryID)
}
