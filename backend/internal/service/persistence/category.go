package persistence

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/ent"
	entcategory "github.com/itk13201/money-rabbit/ent/category"
	"github.com/itk13201/money-rabbit/ent/categoryrule"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
	categoryUC "github.com/itk13201/money-rabbit/internal/usecase/category"
)

// CategoryRepository implements usecase/category.Repository using ent.
type CategoryRepository struct {
	client *ent.Client
}

func NewCategoryRepository(client *ent.Client) *CategoryRepository {
	return &CategoryRepository{client: client}
}

func (r *CategoryRepository) ListCategories(ctx context.Context) ([]*entity.Category, error) {
	slog.InfoContext(ctx, "db.Category.Query started")
	rows, err := r.client.Category.
		Query().
		WithRules(func(q *ent.CategoryRuleQuery) {
			q.Order(ent.Asc(categoryrule.FieldPriority))
		}).
		Order(ent.Asc(entcategory.FieldSortOrder)).
		All(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.Category.Query failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "db.Category.Query finished",
		slog.Group("extra", "count", len(rows)),
	)
	cats := make([]*entity.Category, len(rows))
	for i, row := range rows {
		cats[i] = toCategoryEntity(row)
	}
	return cats, nil
}

func (r *CategoryRepository) GetCategory(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	slog.InfoContext(ctx, "db.Category.Query(one) started",
		slog.Group("extra", "id", id),
	)
	row, err := r.client.Category.
		Query().
		Where(entcategory.ID(id)).
		WithRules().
		Only(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.Category.Query(one) failed",
			slog.Group("extra", "id", id, "error", err),
		)
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("category not found")
		}
		return nil, err
	}
	slog.InfoContext(ctx, "db.Category.Query(one) finished",
		slog.Group("extra", "id", id),
	)
	return toCategoryEntity(row), nil
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, input categoryUC.CreateInput) (*entity.Category, error) {
	slog.InfoContext(ctx, "db.Category.Create started",
		slog.Group("extra", "name", input.Name),
	)
	row, err := r.client.Category.
		Create().
		SetName(input.Name).
		SetColor(input.Color).
		SetIcon(input.Icon).
		SetType(entcategory.Type(input.Type)).
		SetSortOrder(input.SortOrder).
		Save(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.Category.Create failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "db.Category.Create finished",
		slog.Group("extra", "id", row.ID),
	)
	return toCategoryEntity(row), nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, id uuid.UUID, input categoryUC.UpdateInput) (*entity.Category, error) {
	upd := r.client.Category.UpdateOneID(id)
	if input.Name != nil {
		upd.SetName(*input.Name)
	}
	if input.Color != nil {
		upd.SetColor(*input.Color)
	}
	if input.Icon != nil {
		upd.SetIcon(*input.Icon)
	}
	if input.Type != nil {
		upd.SetType(entcategory.Type(*input.Type))
	}
	if input.SortOrder != nil {
		upd.SetSortOrder(*input.SortOrder)
	}
	slog.InfoContext(ctx, "db.Category.UpdateOneID started",
		slog.Group("extra", "id", id),
	)
	row, err := upd.Save(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.Category.UpdateOneID failed",
			slog.Group("extra", "id", id, "error", err),
		)
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("category not found")
		}
		return nil, err
	}
	slog.InfoContext(ctx, "db.Category.UpdateOneID finished",
		slog.Group("extra", "id", id),
	)
	return toCategoryEntity(row), nil
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	slog.InfoContext(ctx, "db.Category.DeleteOneID started",
		slog.Group("extra", "id", id),
	)
	err := r.client.Category.DeleteOneID(id).Exec(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.Category.DeleteOneID failed",
			slog.Group("extra", "id", id, "error", err),
		)
		if ent.IsNotFound(err) {
			return fmt.Errorf("category not found")
		}
		return err
	}
	slog.InfoContext(ctx, "db.Category.DeleteOneID finished",
		slog.Group("extra", "id", id),
	)
	return nil
}

func (r *CategoryRepository) ListAllRules(ctx context.Context) ([]*entity.CategoryRule, error) {
	slog.InfoContext(ctx, "db.CategoryRule.Query started")
	rows, err := r.client.CategoryRule.
		Query().
		Order(ent.Desc(categoryrule.FieldPriority)).
		WithCategory().
		All(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.CategoryRule.Query failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "db.CategoryRule.Query finished",
		slog.Group("extra", "count", len(rows)),
	)
	rules := make([]*entity.CategoryRule, len(rows))
	for i, row := range rows {
		rules[i] = toRuleEntity(row)
	}
	return rules, nil
}

func (r *CategoryRepository) CreateRule(ctx context.Context, input categoryUC.CreateRuleInput) (*entity.CategoryRule, error) {
	slog.InfoContext(ctx, "db.CategoryRule.Create started",
		slog.Group("extra", "category_id", input.CategoryID, "keyword", input.Keyword),
	)
	row, err := r.client.CategoryRule.
		Create().
		SetKeyword(input.Keyword).
		SetPriority(input.Priority).
		SetCategoryID(input.CategoryID).
		Save(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.CategoryRule.Create failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "db.CategoryRule.Create finished",
		slog.Group("extra", "id", row.ID),
	)
	return toRuleEntity(row), nil
}

func (r *CategoryRepository) UpdateRule(ctx context.Context, id uuid.UUID, input categoryUC.UpdateRuleInput) (*entity.CategoryRule, error) {
	upd := r.client.CategoryRule.UpdateOneID(id)
	if input.Keyword != nil {
		upd.SetKeyword(*input.Keyword)
	}
	if input.Priority != nil {
		upd.SetPriority(*input.Priority)
	}
	slog.InfoContext(ctx, "db.CategoryRule.UpdateOneID started",
		slog.Group("extra", "id", id),
	)
	row, err := upd.Save(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.CategoryRule.UpdateOneID failed",
			slog.Group("extra", "id", id, "error", err),
		)
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("rule not found")
		}
		return nil, err
	}
	slog.InfoContext(ctx, "db.CategoryRule.UpdateOneID finished",
		slog.Group("extra", "id", id),
	)
	return toRuleEntity(row), nil
}

func (r *CategoryRepository) DeleteRule(ctx context.Context, id uuid.UUID) error {
	slog.InfoContext(ctx, "db.CategoryRule.DeleteOneID started",
		slog.Group("extra", "id", id),
	)
	err := r.client.CategoryRule.DeleteOneID(id).Exec(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.CategoryRule.DeleteOneID failed",
			slog.Group("extra", "id", id, "error", err),
		)
		if ent.IsNotFound(err) {
			return fmt.Errorf("rule not found")
		}
		return err
	}
	slog.InfoContext(ctx, "db.CategoryRule.DeleteOneID finished",
		slog.Group("extra", "id", id),
	)
	return nil
}

func toCategoryEntity(row *ent.Category) *entity.Category {
	cat := &entity.Category{
		ID:        row.ID,
		Name:      row.Name,
		Color:     row.Color,
		Icon:      row.Icon,
		Type:      entity.CategoryType(row.Type),
		SortOrder: row.SortOrder,
		Rules:     []*entity.CategoryRule{},
	}
	for _, r := range row.Edges.Rules {
		cat.Rules = append(cat.Rules, toRuleEntity(r))
	}
	return cat
}

func toRuleEntity(row *ent.CategoryRule) *entity.CategoryRule {
	r := &entity.CategoryRule{
		ID:       row.ID,
		Keyword:  row.Keyword,
		Priority: row.Priority,
	}
	if row.Edges.Category != nil {
		r.CategoryID = row.Edges.Category.ID
	}
	return r
}
