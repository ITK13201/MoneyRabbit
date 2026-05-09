package transaction

import (
	"context"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/itk13201/money-rabbit/internal/domain/entity"
)

type categoryRepository interface {
	ListCategories(ctx context.Context) ([]*entity.Category, error)
	ListAllRules(ctx context.Context) ([]*entity.CategoryRule, error)
}

// ImportUsecase handles the CSV import confirm flow.
type ImportUsecase struct {
	repo         Repository
	categoryRepo categoryRepository
	classifier   Classifier
}

func NewImportUsecase(repo Repository, categoryRepo categoryRepository, classifier Classifier) *ImportUsecase {
	return &ImportUsecase{
		repo:         repo,
		categoryRepo: categoryRepo,
		classifier:   classifier,
	}
}

type ImportResult struct {
	Imported int
	Skipped  int
	Errors   []string
}

// Confirm deduplicates, classifies, and saves a list of parsed transactions.
func (u *ImportUsecase) Confirm(ctx context.Context, inputs []CreateInput) (*ImportResult, error) {
	if len(inputs) == 0 {
		return &ImportResult{}, nil
	}

	// 1. Duplicate detection
	slog.InfoContext(ctx, "transactionRepo.FindDuplicates started",
		slog.Group("extra", "count", len(inputs)),
	)
	dupMap, err := u.repo.FindDuplicates(ctx, inputs)
	if err != nil {
		slog.ErrorContext(ctx, "transactionRepo.FindDuplicates failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "transactionRepo.FindDuplicates finished",
		slog.Group("extra", "duplicates", len(dupMap)),
	)

	toInsert := make([]CreateInput, 0, len(inputs))
	skipped := 0
	for i, inp := range inputs {
		if dupMap[i] {
			skipped++
		} else {
			toInsert = append(toInsert, inp)
		}
	}

	if len(toInsert) == 0 {
		return &ImportResult{Skipped: skipped}, nil
	}

	// 2. Load keyword rules and categories
	slog.InfoContext(ctx, "categoryRepo.ListAllRules started")
	rules, err := u.categoryRepo.ListAllRules(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "categoryRepo.ListAllRules failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "categoryRepo.ListAllRules finished",
		slog.Group("extra", "count", len(rules)),
	)

	slog.InfoContext(ctx, "categoryRepo.ListCategories started")
	categories, err := u.categoryRepo.ListCategories(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "categoryRepo.ListCategories failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "categoryRepo.ListCategories finished",
		slog.Group("extra", "count", len(categories)),
	)

	// 3. Keyword rule matching (rules ordered by priority desc in persistence)
	var unclassifiedIdxs []int
	for i := range toInsert {
		catID := matchKeywordRule(toInsert[i].Description, rules)
		if catID != nil {
			toInsert[i].CategoryID = catID
		} else {
			unclassifiedIdxs = append(unclassifiedIdxs, i)
		}
	}

	// 4. Claude API classification for unmatched transactions
	if len(unclassifiedIdxs) > 0 && u.classifier != nil && len(categories) > 0 {
		descs := make([]string, len(unclassifiedIdxs))
		for j, idx := range unclassifiedIdxs {
			descs[j] = toInsert[idx].Description
		}

		slog.InfoContext(ctx, "classifier.Classify service started",
			slog.Group("extra", "unclassified_count", len(descs)),
		)
		classified, err := u.classifier.Classify(ctx, descs, categories)
		if err != nil {
			slog.ErrorContext(ctx, "classifier.Classify service failed",
				slog.Group("extra", "error", err),
			)
			// On classifier error, leave CategoryID as nil (uncategorized)
		} else {
			slog.InfoContext(ctx, "classifier.Classify service finished",
				slog.Group("extra", "classified_count", len(classified)),
			)
			for j, idx := range unclassifiedIdxs {
				if catID, ok := classified[descs[j]]; ok {
					toInsert[idx].CategoryID = catID
				}
			}
		}
	}

	// 5. Bulk insert
	slog.InfoContext(ctx, "transactionRepo.BulkCreateTransactions started",
		slog.Group("extra", "count", len(toInsert)),
	)
	created, err := u.repo.BulkCreateTransactions(ctx, toInsert)
	if err != nil {
		slog.ErrorContext(ctx, "transactionRepo.BulkCreateTransactions failed",
			slog.Group("extra", "error", err),
		)
		return nil, err
	}
	slog.InfoContext(ctx, "transactionRepo.BulkCreateTransactions finished",
		slog.Group("extra", "created", len(created)),
	)

	return &ImportResult{
		Imported: len(created),
		Skipped:  skipped,
		Errors:   []string{},
	}, nil
}

// matchKeywordRule returns the category ID of the highest-priority matching rule,
// or nil if no rule matches.
func matchKeywordRule(description string, rules []*entity.CategoryRule) *uuid.UUID {
	for _, r := range rules {
		if strings.Contains(description, r.Keyword) {
			id := r.CategoryID
			return &id
		}
	}
	return nil
}
