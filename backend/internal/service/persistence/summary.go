package persistence

import (
	"context"
	"database/sql"
	"log/slog"

	sumUC "github.com/itk13201/money-rabbit/internal/usecase/summary"
)

// SummaryRepository implements usecase/summary.Repository using raw SQL.
type SummaryRepository struct {
	db *sql.DB
}

func NewSummaryRepository(db *sql.DB) *SummaryRepository {
	return &SummaryRepository{db: db}
}

// CategoryAnnualSummary returns aggregated expense per category for the given year.
func (r *SummaryRepository) CategoryAnnualSummary(ctx context.Context, year int) ([]sumUC.CategoryAnnualItem, error) {
	const q = `
		SELECT
			t.category_transactions,
			COALESCE(c.name,  '未分類')   AS category_name,
			COALESCE(c.color, '#94a3b8') AS category_color,
			SUM(-t.amount)               AS total_expense
		FROM transactions t
		LEFT JOIN categories c ON c.id = t.category_transactions
		WHERE YEAR(t.date) = ?
		  AND t.amount < 0
		GROUP BY t.category_transactions, c.name, c.color
		ORDER BY total_expense DESC
	`

	slog.InfoContext(ctx, "db.CategoryAnnualSummary.Query started",
		slog.Group("extra", "year", year),
	)
	rows, err := r.db.QueryContext(ctx, q, year)
	if err != nil {
		slog.ErrorContext(ctx, "db.CategoryAnnualSummary.Query failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	defer rows.Close()

	var result []sumUC.CategoryAnnualItem
	for rows.Next() {
		var item sumUC.CategoryAnnualItem
		var categoryID sql.NullString
		if err := rows.Scan(&categoryID, &item.CategoryName, &item.CategoryColor, &item.TotalExpense); err != nil {
			return nil, err
		}
		if categoryID.Valid {
			item.CategoryID = &categoryID.String
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "db.CategoryAnnualSummary.Query finished",
		slog.Group("extra", "count", len(result)),
	)
	return result, nil
}

// MonthlySummary returns aggregated income and expense per month for the given year.
func (r *SummaryRepository) MonthlySummary(ctx context.Context, year int) ([]sumUC.MonthSummary, error) {
	const q = `
		SELECT
			YEAR(date)  AS year,
			MONTH(date) AS month,
			COALESCE(SUM(CASE WHEN amount > 0 THEN amount  ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN amount < 0 THEN -amount ELSE 0 END), 0) AS expense
		FROM transactions
		WHERE YEAR(date) = ?
		GROUP BY YEAR(date), MONTH(date)
		ORDER BY MONTH(date)
	`

	slog.InfoContext(ctx, "db.MonthlySummary.Query started",
		slog.Group("extra", "year", year),
	)
	rows, err := r.db.QueryContext(ctx, q, year)
	if err != nil {
		slog.ErrorContext(ctx, "db.MonthlySummary.Query failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	defer rows.Close()

	var result []sumUC.MonthSummary
	for rows.Next() {
		var m sumUC.MonthSummary
		if err := rows.Scan(&m.Year, &m.Month, &m.Income, &m.Expense); err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "db.MonthlySummary.Query finished",
		slog.Group("extra", "count", len(result)),
	)
	return result, nil
}
