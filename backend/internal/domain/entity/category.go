package entity

import "github.com/google/uuid"

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
	CategoryTypeBoth    CategoryType = "both"
)

type Category struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Color     string         `json:"color"`
	Icon      string         `json:"icon"`
	Type      CategoryType   `json:"type"`
	SortOrder int            `json:"sort_order"`
	Rules     []*CategoryRule `json:"rules,omitempty"`
}

type CategoryRule struct {
	ID         uuid.UUID `json:"id"`
	Keyword    string    `json:"keyword"`
	Priority   int       `json:"priority"`
	CategoryID uuid.UUID `json:"category_id"`
}
