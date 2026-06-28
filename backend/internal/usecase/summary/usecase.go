package summary

import "context"

// Repository is the interface for summary queries.
type Repository interface {
	MonthlySummary(ctx context.Context, year int) ([]MonthSummary, error)
	CategoryAnnualSummary(ctx context.Context, year int) ([]CategoryAnnualItem, error)
}

// MonthSummary holds aggregated income/expense for one month.
type MonthSummary struct {
	Year    int `json:"year"`
	Month   int `json:"month"`
	Income  int `json:"income"`
	Expense int `json:"expense"`
}

// CategoryAnnualItem holds aggregated expense for one category across a year.
type CategoryAnnualItem struct {
	CategoryID    *string `json:"category_id"`
	CategoryName  string  `json:"category_name"`
	CategoryColor string  `json:"category_color"`
	TotalExpense  int     `json:"total_expense"`
	Percentage    float64 `json:"percentage"`
}

// CategoryAnnualSummary is the top-level response for annual category distribution.
type CategoryAnnualSummary struct {
	Year         int                  `json:"year"`
	TotalExpense int                  `json:"total_expense"`
	Categories   []CategoryAnnualItem `json:"categories"`
}
