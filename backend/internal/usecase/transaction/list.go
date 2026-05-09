package transaction

import (
	"context"

	"github.com/itk13201/money-rabbit/internal/domain/entity"
)

// ListUsecase handles transaction listing.
type ListUsecase struct {
	repo Repository
}

func NewListUsecase(repo Repository) *ListUsecase {
	return &ListUsecase{repo: repo}
}

type ListResult struct {
	Transactions []*entity.Transaction
	Total        int
}

func (u *ListUsecase) List(ctx context.Context, filter ListFilter) (*ListResult, error) {
	txs, total, err := u.repo.ListTransactions(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &ListResult{Transactions: txs, Total: total}, nil
}
