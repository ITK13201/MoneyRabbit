package entity

import "github.com/google/uuid"

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
	CategoryTypeBoth    CategoryType = "both"
)

type Category struct {
	ID        uuid.UUID
	Name      string
	Color     string
	Icon      string
	Type      CategoryType
	SortOrder int
	Rules     []*CategoryRule
}

type CategoryRule struct {
	ID         uuid.UUID
	Keyword    string
	Priority   int
	CategoryID uuid.UUID
}
