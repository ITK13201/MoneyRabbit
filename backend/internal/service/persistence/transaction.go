package persistence

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/ent"
	entcategory "github.com/itk13201/money-rabbit/ent/category"
	"github.com/itk13201/money-rabbit/ent/transaction"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
	txUC "github.com/itk13201/money-rabbit/internal/usecase/transaction"
)

// TransactionRepository implements usecase/transaction.Repository using ent.
type TransactionRepository struct {
	client *ent.Client
}

func NewTransactionRepository(client *ent.Client) *TransactionRepository {
	return &TransactionRepository{client: client}
}

// FindDuplicates returns a map of input index → true for inputs that already exist in the DB.
// Duplicate key: (import_format_id, date, description, amount).
func (r *TransactionRepository) FindDuplicates(ctx context.Context, inputs []txUC.CreateInput) (map[int]bool, error) {
	result := make(map[int]bool, len(inputs))
	slog.InfoContext(ctx, "db.Transaction.Exist started",
		slog.Group("extra", "count", len(inputs)),
	)
	for i, inp := range inputs {
		exists, err := r.client.Transaction.
			Query().
			Where(
				transaction.ImportFormatIDEQ(transaction.ImportFormatID(inp.ImportFormatID)),
				transaction.DateEQ(inp.Date),
				transaction.DescriptionEQ(inp.Description),
				transaction.AmountEQ(inp.Amount),
			).
			Exist(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "db.Transaction.Exist failed",
				slog.Group("extra", "index", i, "error", err),
			)
			return nil, err
		}
		result[i] = exists
	}
	slog.InfoContext(ctx, "db.Transaction.Exist finished",
		slog.Group("extra", "count", len(inputs)),
	)
	return result, nil
}

// BulkCreateTransactions inserts multiple transactions in a single operation.
func (r *TransactionRepository) BulkCreateTransactions(ctx context.Context, inputs []txUC.CreateInput) ([]*entity.Transaction, error) {
	builders := make([]*ent.TransactionCreate, len(inputs))
	for i, inp := range inputs {
		b := r.client.Transaction.
			Create().
			SetDate(inp.Date).
			SetDescription(inp.Description).
			SetAmount(inp.Amount).
			SetImportFormatID(transaction.ImportFormatID(inp.ImportFormatID)).
			SetImportedAt(time.Now())
		if inp.CategoryID != nil {
			b = b.SetCategoryID(*inp.CategoryID)
		}
		builders[i] = b
	}

	slog.InfoContext(ctx, "db.Transaction.CreateBulk started",
		slog.Group("extra", "count", len(builders)),
	)
	rows, err := r.client.Transaction.CreateBulk(builders...).Save(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.Transaction.CreateBulk failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "db.Transaction.CreateBulk finished",
		slog.Group("extra", "created", len(rows)),
	)

	txs := make([]*entity.Transaction, len(rows))
	for i, row := range rows {
		txs[i] = toTransactionEntity(row)
	}
	return txs, nil
}

// ListTransactions returns paginated transactions with optional filters.
func (r *TransactionRepository) ListTransactions(ctx context.Context, filter txUC.ListFilter) ([]*entity.Transaction, int, error) {
	q := r.client.Transaction.Query().WithCategory()

	if filter.Year != nil && filter.Month != nil {
		start := time.Date(*filter.Year, time.Month(*filter.Month), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, 0)
		q = q.Where(
			transaction.DateGTE(start),
			transaction.DateLT(end),
		)
	} else if filter.Year != nil {
		start := time.Date(*filter.Year, 1, 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(1, 0, 0)
		q = q.Where(
			transaction.DateGTE(start),
			transaction.DateLT(end),
		)
	}

	if filter.CategoryID != nil {
		q = q.Where(transaction.HasCategoryWith(entcategory.ID(*filter.CategoryID)))
	}

	slog.InfoContext(ctx, "db.Transaction.Query started",
		slog.Group("extra", "filter", filter),
	)
	total, err := q.Count(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.Transaction.Query failed",
			slog.Group("extra", "error", err),
		)
		return nil, 0, err
	}

	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}
	offset := filter.Page * pageSize

	rows, err := q.
		Order(ent.Desc(transaction.FieldDate)).
		Limit(pageSize).
		Offset(offset).
		All(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.Transaction.Query failed",
			slog.Group("extra", "error", err),
		)
		return nil, 0, err
	}
	slog.InfoContext(ctx, "db.Transaction.Query finished",
		slog.Group("extra", "count", len(rows), "total", total),
	)

	txs := make([]*entity.Transaction, len(rows))
	for i, row := range rows {
		txs[i] = toTransactionEntity(row)
	}
	return txs, total, nil
}

// UpdateTransactionCategory sets or clears the category of a transaction.
func (r *TransactionRepository) UpdateTransactionCategory(ctx context.Context, id uuid.UUID, categoryID *uuid.UUID) (*entity.Transaction, error) {
	upd := r.client.Transaction.UpdateOneID(id)
	if categoryID != nil {
		upd = upd.SetCategoryID(*categoryID)
	} else {
		upd = upd.ClearCategory()
	}
	slog.InfoContext(ctx, "db.Transaction.UpdateOneID started",
		slog.Group("extra", "id", id, "category_id", categoryID),
	)
	_, err := upd.Save(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.Transaction.UpdateOneID failed",
			slog.Group("extra", "id", id, "error", err),
		)
		return nil, err
	}
	row, err := r.client.Transaction.Query().
		Where(transaction.IDEQ(id)).
		WithCategory().
		Only(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.Transaction.UpdateOneID failed",
			slog.Group("extra", "id", id, "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "db.Transaction.UpdateOneID finished",
		slog.Group("extra", "id", id),
	)
	return toTransactionEntity(row), nil
}

// DeleteTransaction deletes a transaction by ID.
func (r *TransactionRepository) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	slog.InfoContext(ctx, "db.Transaction.DeleteOneID started",
		slog.Group("extra", "id", id),
	)
	err := r.client.Transaction.DeleteOneID(id).Exec(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.Transaction.DeleteOneID failed",
			slog.Group("extra", "id", id, "error", err),
		)
		return err
	}
	slog.InfoContext(ctx, "db.Transaction.DeleteOneID finished",
		slog.Group("extra", "id", id),
	)
	return nil
}

// CreateManualTransaction inserts a single manually-entered transaction (no import_format_id).
func (r *TransactionRepository) CreateManualTransaction(ctx context.Context, input txUC.CreateManualInput) (*entity.Transaction, error) {
	b := r.client.Transaction.
		Create().
		SetDate(input.Date).
		SetDescription(input.Description).
		SetAmount(input.Amount).
		SetImportedAt(time.Now())
	if input.CategoryID != nil {
		b = b.SetCategoryID(*input.CategoryID)
	}

	slog.InfoContext(ctx, "db.Transaction.Create started")
	row, err := b.Save(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.Transaction.Create failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "db.Transaction.Create finished",
		slog.Group("extra", "id", row.ID),
	)

	// Re-fetch with category edge
	row, err = r.client.Transaction.Query().
		Where(transaction.IDEQ(row.ID)).
		WithCategory().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return toTransactionEntity(row), nil
}

// UpdateTransaction updates date, description, amount, and category of a transaction.
func (r *TransactionRepository) UpdateTransaction(ctx context.Context, id uuid.UUID, input txUC.UpdateInput) (*entity.Transaction, error) {
	upd := r.client.Transaction.UpdateOneID(id).
		SetDate(input.Date).
		SetDescription(input.Description).
		SetAmount(input.Amount)
	if input.CategoryID != nil {
		upd = upd.SetCategoryID(*input.CategoryID)
	} else {
		upd = upd.ClearCategory()
	}

	slog.InfoContext(ctx, "db.Transaction.UpdateOneID started",
		slog.Group("extra", "id", id),
	)
	_, err := upd.Save(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.Transaction.UpdateOneID failed",
			slog.Group("extra", "id", id, "error", err),
		)
		return nil, err
	}

	row, err := r.client.Transaction.Query().
		Where(transaction.IDEQ(id)).
		WithCategory().
		Only(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "db.Transaction.Query failed",
			slog.Group("extra", "id", id, "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "db.Transaction.UpdateOneID finished",
		slog.Group("extra", "id", id),
	)
	return toTransactionEntity(row), nil
}

func toTransactionEntity(row *ent.Transaction) *entity.Transaction {
	tx := &entity.Transaction{
		ID:          row.ID,
		Date:        row.Date,
		Description: row.Description,
		Amount:      row.Amount,
		ImportedAt:  row.ImportedAt,
	}
	if row.ImportFormatID != nil {
		s := string(*row.ImportFormatID)
		tx.ImportFormatID = &s
	}
	if row.Edges.Category != nil {
		catID := row.Edges.Category.ID
		tx.CategoryID = &catID
		tx.Category = toCategoryEntity(row.Edges.Category)
	}
	return tx
}
