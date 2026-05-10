package transaction

import (
	"context"

	"github.com/google/uuid"
)

type DeleteUsecase struct {
	repo Repository
}

func NewDeleteUsecase(repo Repository) *DeleteUsecase {
	return &DeleteUsecase{repo: repo}
}

func (u *DeleteUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteTransaction(ctx, id)
}
