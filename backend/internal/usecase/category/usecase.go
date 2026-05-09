package category

import (
	"context"

	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
)

// Repository is the interface that the persistence layer implements.
type Repository interface {
	ListCategories(ctx context.Context) ([]*entity.Category, error)
	GetCategory(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	CreateCategory(ctx context.Context, input CreateInput) (*entity.Category, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, input UpdateInput) (*entity.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error

	ListAllRules(ctx context.Context) ([]*entity.CategoryRule, error)
	CreateRule(ctx context.Context, input CreateRuleInput) (*entity.CategoryRule, error)
	UpdateRule(ctx context.Context, id uuid.UUID, input UpdateRuleInput) (*entity.CategoryRule, error)
	DeleteRule(ctx context.Context, id uuid.UUID) error
}

type CreateInput struct {
	Name      string
	Color     string
	Icon      string
	Type      entity.CategoryType
	SortOrder int
}

type UpdateInput struct {
	Name      *string
	Color     *string
	Icon      *string
	Type      *entity.CategoryType
	SortOrder *int
}

type CreateRuleInput struct {
	CategoryID uuid.UUID
	Keyword    string
	Priority   int
}

type UpdateRuleInput struct {
	Keyword  *string
	Priority *int
}

// Usecase implements business logic for categories and rules.
type Usecase struct {
	repo Repository
}

func New(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) List(ctx context.Context) ([]*entity.Category, error) {
	return u.repo.ListCategories(ctx)
}

func (u *Usecase) Get(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	return u.repo.GetCategory(ctx, id)
}

func (u *Usecase) Create(ctx context.Context, input CreateInput) (*entity.Category, error) {
	return u.repo.CreateCategory(ctx, input)
}

func (u *Usecase) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*entity.Category, error) {
	return u.repo.UpdateCategory(ctx, id, input)
}

func (u *Usecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteCategory(ctx, id)
}

func (u *Usecase) CreateRule(ctx context.Context, input CreateRuleInput) (*entity.CategoryRule, error) {
	return u.repo.CreateRule(ctx, input)
}

func (u *Usecase) UpdateRule(ctx context.Context, id uuid.UUID, input UpdateRuleInput) (*entity.CategoryRule, error) {
	return u.repo.UpdateRule(ctx, id, input)
}

func (u *Usecase) DeleteRule(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteRule(ctx, id)
}
