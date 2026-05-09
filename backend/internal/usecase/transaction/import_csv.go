package transaction

import (
	"context"
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
	dupMap, err := u.repo.FindDuplicates(ctx, inputs)
	if err != nil {
		return nil, err
	}

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
	rules, err := u.categoryRepo.ListAllRules(ctx)
	if err != nil {
		return nil, err
	}
	categories, err := u.categoryRepo.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

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

		classified, err := u.classifier.Classify(ctx, descs, categories)
		if err == nil {
			for j, idx := range unclassifiedIdxs {
				if catID, ok := classified[descs[j]]; ok {
					toInsert[idx].CategoryID = catID
				}
			}
		}
		// On classifier error, leave CategoryID as nil (uncategorized)
	}

	// 5. Bulk insert
	created, err := u.repo.BulkCreateTransactions(ctx, toInsert)
	if err != nil {
		return nil, err
	}

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
