package summary

import (
	"context"
	"log/slog"
)

// CategoryAnnualUsecase aggregates annual expense by category.
type CategoryAnnualUsecase struct {
	repo Repository
}

func NewCategoryAnnualUsecase(repo Repository) *CategoryAnnualUsecase {
	return &CategoryAnnualUsecase{repo: repo}
}

func (u *CategoryAnnualUsecase) CategoryAnnual(ctx context.Context, year int) (*CategoryAnnualSummary, error) {
	slog.InfoContext(ctx, "summaryRepo.CategoryAnnualSummary started",
		slog.Group("extra", "year", year),
	)
	items, err := u.repo.CategoryAnnualSummary(ctx, year)
	if err != nil {
		slog.ErrorContext(ctx, "summaryRepo.CategoryAnnualSummary failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "summaryRepo.CategoryAnnualSummary finished",
		slog.Group("extra", "count", len(items)),
	)

	total := 0
	for _, item := range items {
		total += item.TotalExpense
	}
	if total > 0 {
		for i := range items {
			items[i].Percentage = float64(items[i].TotalExpense) / float64(total) * 100
		}
	}
	if items == nil {
		items = []CategoryAnnualItem{}
	}

	return &CategoryAnnualSummary{
		Year:         year,
		TotalExpense: total,
		Categories:   items,
	}, nil
}
