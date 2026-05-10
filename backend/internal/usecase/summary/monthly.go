package summary

import (
	"context"
	"log/slog"
)

// MonthlyUsecase handles monthly summary aggregation.
type MonthlyUsecase struct {
	repo Repository
}

func NewMonthlyUsecase(repo Repository) *MonthlyUsecase {
	return &MonthlyUsecase{repo: repo}
}

func (u *MonthlyUsecase) Monthly(ctx context.Context, year int) ([]MonthSummary, error) {
	slog.InfoContext(ctx, "summaryRepo.MonthlySummary started",
		slog.Group("extra", "year", year),
	)
	result, err := u.repo.MonthlySummary(ctx, year)
	if err != nil {
		slog.ErrorContext(ctx, "summaryRepo.MonthlySummary failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "summaryRepo.MonthlySummary finished",
		slog.Group("extra", "count", len(result)),
	)
	return result, nil
}
