package transaction

import (
	"context"
	"log/slog"

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
	slog.InfoContext(ctx, "transactionRepo.ListTransactions started",
		slog.Group("extra", "filter", filter),
	)
	txs, total, err := u.repo.ListTransactions(ctx, filter)
	if err != nil {
		slog.ErrorContext(ctx, "transactionRepo.ListTransactions failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "transactionRepo.ListTransactions finished",
		slog.Group("extra", "count", len(txs), "total", total),
	)
	return &ListResult{Transactions: txs, Total: total}, nil
}
