package summary

import "context"

// Repository is the interface for monthly summary queries.
type Repository interface {
	MonthlySummary(ctx context.Context, year int) ([]MonthSummary, error)
}

// MonthSummary holds aggregated income/expense for one month.
type MonthSummary struct {
	Year    int `json:"year"`
	Month   int `json:"month"`
	Income  int `json:"income"`
	Expense int `json:"expense"`
}
