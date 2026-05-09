package category

import (
	"context"
	"log/slog"

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
	slog.InfoContext(ctx, "categoryRepo.ListCategories started")
	cats, err := u.repo.ListCategories(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "categoryRepo.ListCategories failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "categoryRepo.ListCategories finished",
		slog.Group("extra", "count", len(cats)),
	)
	return cats, nil
}

func (u *Usecase) Get(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	slog.InfoContext(ctx, "categoryRepo.GetCategory started",
		slog.Group("extra", "id", id),
	)
	cat, err := u.repo.GetCategory(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "categoryRepo.GetCategory failed",
			slog.Group("extra", "id", id, "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "categoryRepo.GetCategory finished",
		slog.Group("extra", "id", id),
	)
	return cat, nil
}

func (u *Usecase) Create(ctx context.Context, input CreateInput) (*entity.Category, error) {
	slog.InfoContext(ctx, "categoryRepo.CreateCategory started",
		slog.Group("extra", "name", input.Name),
	)
	cat, err := u.repo.CreateCategory(ctx, input)
	if err != nil {
		slog.ErrorContext(ctx, "categoryRepo.CreateCategory failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "categoryRepo.CreateCategory finished",
		slog.Group("extra", "id", cat.ID),
	)
	return cat, nil
}

func (u *Usecase) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*entity.Category, error) {
	slog.InfoContext(ctx, "categoryRepo.UpdateCategory started",
		slog.Group("extra", "id", id),
	)
	cat, err := u.repo.UpdateCategory(ctx, id, input)
	if err != nil {
		slog.ErrorContext(ctx, "categoryRepo.UpdateCategory failed",
			slog.Group("extra", "id", id, "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "categoryRepo.UpdateCategory finished",
		slog.Group("extra", "id", id),
	)
	return cat, nil
}

func (u *Usecase) Delete(ctx context.Context, id uuid.UUID) error {
	slog.InfoContext(ctx, "categoryRepo.DeleteCategory started",
		slog.Group("extra", "id", id),
	)
	if err := u.repo.DeleteCategory(ctx, id); err != nil {
		slog.ErrorContext(ctx, "categoryRepo.DeleteCategory failed",
			slog.Group("extra", "id", id, "error", err),
		)
		return err
	}
	slog.InfoContext(ctx, "categoryRepo.DeleteCategory finished",
		slog.Group("extra", "id", id),
	)
	return nil
}

func (u *Usecase) CreateRule(ctx context.Context, input CreateRuleInput) (*entity.CategoryRule, error) {
	slog.InfoContext(ctx, "categoryRepo.CreateRule started",
		slog.Group("extra", "category_id", input.CategoryID, "keyword", input.Keyword),
	)
	rule, err := u.repo.CreateRule(ctx, input)
	if err != nil {
		slog.ErrorContext(ctx, "categoryRepo.CreateRule failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "categoryRepo.CreateRule finished",
		slog.Group("extra", "id", rule.ID),
	)
	return rule, nil
}

func (u *Usecase) UpdateRule(ctx context.Context, id uuid.UUID, input UpdateRuleInput) (*entity.CategoryRule, error) {
	slog.InfoContext(ctx, "categoryRepo.UpdateRule started",
		slog.Group("extra", "id", id),
	)
	rule, err := u.repo.UpdateRule(ctx, id, input)
	if err != nil {
		slog.ErrorContext(ctx, "categoryRepo.UpdateRule failed",
			slog.Group("extra", "id", id, "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "categoryRepo.UpdateRule finished",
		slog.Group("extra", "id", id),
	)
	return rule, nil
}

func (u *Usecase) DeleteRule(ctx context.Context, id uuid.UUID) error {
	slog.InfoContext(ctx, "categoryRepo.DeleteRule started",
		slog.Group("extra", "id", id),
	)
	if err := u.repo.DeleteRule(ctx, id); err != nil {
		slog.ErrorContext(ctx, "categoryRepo.DeleteRule failed",
			slog.Group("extra", "id", id, "error", err),
		)
		return err
	}
	slog.InfoContext(ctx, "categoryRepo.DeleteRule finished",
		slog.Group("extra", "id", id),
	)
	return nil
}
