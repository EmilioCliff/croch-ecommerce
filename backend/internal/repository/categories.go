package repository

import (
	"context"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
)

type Category struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *Category) Validate() error {
	if c.Name == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "content is required")
	}

	return nil
}

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	GetCategory(ctx context.Context, id uint32) (*Category, error)
	ListCategories(ctx context.Context) ([]*Category, error)
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, id uint32) error
}
