package mysql

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/mysql/generated"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
)

var _ repository.CategoryRepository = (*CategoryRepository)(nil)

type CategoryRepository struct {
	db      *Store
	queries generated.Querier
}

func NewCategoryRepository(db *Store) *CategoryRepository {
	q := generated.New(db.db)

	return &CategoryRepository{
		db:      db,
		queries: q,
	}
}

func (c *CategoryRepository) CreateCategory(ctx context.Context, category *repository.Category) (*repository.Category, error) {
	if err := category.Validate(); err != nil {
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "%v", err)
	}

	result, err := c.queries.CreateCategory(ctx, generated.CreateCategoryParams{
		Name:        category.Name,
		Description: category.Description,
	})
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create category: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last inserted id: %v", err)
	}

	category.ID = uint32(id)

	return category, nil
}

func (c *CategoryRepository) GetCategory(ctx context.Context, id uint32) (*repository.Category, error) {
	category, err := c.queries.GetCategory(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no category found with id %d", id)
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get category: %v", err)
	}

	return &repository.Category{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (c *CategoryRepository) ListCategories(ctx context.Context) ([]*repository.Category, error) {
	categories, err := c.queries.ListCategories(ctx)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list categories: %v", err)
	}

	var result []*repository.Category
	for _, category := range categories {
		result = append(result, &repository.Category{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}

	return result, nil
}

func (c *CategoryRepository) UpdateCategory(ctx context.Context, category *repository.Category) error {
	if err := category.Validate(); err != nil {
		return pkg.Errorf(pkg.INVALID_ERROR, "%v", err)
	}

	err := c.queries.UpdateCategory(ctx, generated.UpdateCategoryParams{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	})
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update category: %v", err)
	}

	return nil
}

func (c *CategoryRepository) DeleteCategory(ctx context.Context, id uint32) error {
	err := c.queries.DeleteCategory(ctx, id)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to delete category: %v", err)
	}

	return nil
}
