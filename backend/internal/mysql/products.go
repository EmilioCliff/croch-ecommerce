package mysql

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/mysql/generated"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
)

var _ repository.ProductRepository = (*ProductRepository)(nil)

type ProductRepository struct {
	db      *Store
	queries generated.Querier
}

func NewProductRepository(store *Store) *ProductRepository {
	queries := generated.New(store.db)

	return &ProductRepository{
		queries: queries,
		db:      store,
	}
}

func (p *ProductRepository) CreateProduct(ctx context.Context, product *repository.Product) (*repository.Product, error) {
	if err := product.Validate(); err != nil {
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "%v", err)
	}

	result, err := p.queries.CreateProduct(ctx, generated.CreateProductParams{
		Name:            product.Name,
		Description:     product.Description,
		RegularPrice:    product.RegularPrice,
		DiscountedPrice: product.DiscountedPrice,
		Quantity:        product.Quantity,
		CategoryID:      product.CategoryID,
		SizeOption:      product.SizeOption,
		ColorOption:     product.ColorOption,
		Seasonal:        product.Seasonal,
		Featured:        product.Featured,
		ImgUrls:         product.ImgUrls,
		UpdatedBy:       product.UpdatedBy,
	})
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create product: %v", err)
	}

	createdID, err := result.LastInsertId()
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last inserted id: %v", err)
	}

	product.ID = uint32(createdID)

	return product, nil
}

func (p *ProductRepository) GetProduct(ctx context.Context, id uint32) (*repository.Product, error) {
	product, err := p.queries.GetProduct(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "product not found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get product: %v", err)
	}

	return toRepositoryProduct(product), nil
}

func (p *ProductRepository) GetProductName(ctx context.Context, id uint32) (string, error) {
	email, err := p.queries.GetProductName(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", pkg.Errorf(pkg.NOT_FOUND_ERROR, "product not found")
		}

		return "", pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get product: %v", err)
	}

	return email, nil
}

func (p *ProductRepository) UpdateProduct(ctx context.Context, product *repository.UpdateProduct) error {
	_, err := p.GetProduct(ctx, product.ID)
	if err != nil {
		return err
	}

	if err := product.Validate(); err != nil {
		return pkg.Errorf(pkg.INVALID_ERROR, "%v", err)
	}

	var req generated.UpdateProductParams

	req.ID = product.ID
	req.UpdatedBy = product.UpdatedBy

	if product.Name != nil {
		req.Name = sql.NullString{
			Valid:  true,
			String: *product.Name,
		}
	}

	if product.Description != nil {
		req.Description = sql.NullString{
			Valid:  true,
			String: *product.Description,
		}
	}

	if product.RegularPrice != nil {
		req.RegularPrice = *product.RegularPrice
	}

	if product.DiscountedPrice != nil {
		req.DiscountedPrice = *product.DiscountedPrice
	}

	if product.Quantity != nil {
		req.Quantity = sql.NullInt32{
			Valid: true,
			Int32: int32(*product.Quantity),
		}
	}

	if product.CategoryID != nil {
		req.CategoryID = sql.NullInt32{
			Valid: true,
			Int32: int32(*product.CategoryID),
		}
	}

	if product.SizeOption != nil {
		req.SizeOption = *product.SizeOption
	}

	if product.ColorOption != nil {
		req.ColorOption = *product.ColorOption
	}

	if product.Seasonal != nil {
		req.Seasonal = sql.NullBool{
			Valid: true,
			Bool:  *product.Seasonal,
		}
	}

	if product.Featured != nil {
		req.Featured = sql.NullBool{
			Valid: true,
			Bool:  *product.Seasonal,
		}
	}

	if product.ImgUrls != nil {
		req.ImgUrls = *product.ImgUrls
	}

	if product.UpdatedBy <= 0 {
		req.UpdatedBy = product.UpdatedBy
	}

	err = p.queries.UpdateProduct(ctx, req)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update product: %v", err)
	}

	return nil
}

func (p *ProductRepository) UpdateProductQuantity(ctx context.Context, id uint32, quantity uint32) error {
	_, err := p.GetProduct(ctx, id)
	if err != nil {
		return err
	}

	err = p.queries.UpdateProductQuantity(ctx, generated.UpdateProductQuantityParams{
		ID:       id,
		Quantity: quantity,
	})

	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update product quantity: %v", err)
	}

	return nil
}

func (p *ProductRepository) ListProducts(ctx context.Context) ([]*repository.Product, error) {
	products, err := p.queries.ListProducts(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no products found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list products: %v", err)
	}

	var result []*repository.Product
	for _, product := range products {
		result = append(result, toRepositoryProduct(product))
	}

	return result, nil
}

func (p *ProductRepository) ListNewProducts(ctx context.Context) ([]*repository.Product, error) {
	products, err := p.queries.ListNewProducts(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no new products found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list products: %v", err)
	}

	var result []*repository.Product
	for _, product := range products {
		result = append(result, toRepositoryProduct(product))
	}

	return result, nil
}

func (p *ProductRepository) ListSeasonalProducts(ctx context.Context) ([]*repository.Product, error) {
	products, err := p.queries.ListSeasonalProducts(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no seasonal products found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list products: %v", err)
	}

	var result []*repository.Product
	for _, product := range products {
		result = append(result, toRepositoryProduct(product))
	}

	return result, nil
}

func (p *ProductRepository) ListFeaturedProducts(ctx context.Context) ([]*repository.Product, error) {
	products, err := p.queries.ListFeaturedProducts(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no featured products found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list products: %v", err)
	}

	var result []*repository.Product
	for _, product := range products {
		result = append(result, toRepositoryProduct(product))
	}

	return result, nil
}

func (p *ProductRepository) ListDiscountedProducts(ctx context.Context) ([]*repository.Product, error) {
	products, err := p.queries.ListDiscountedProducts(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no discounted products found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list products: %v", err)
	}

	var result []*repository.Product
	for _, product := range products {
		result = append(result, toRepositoryProduct(product))
	}

	return result, nil
}

func (p *ProductRepository) ListProductsByCategory(ctx context.Context, categoryID uint32) ([]*repository.Product, error) {
	products, err := p.queries.ListProductsByCategory(ctx, categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no products found with category %v", categoryID)
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list products: %v", err)
	}

	var result []*repository.Product
	for _, product := range products {
		result = append(result, toRepositoryProduct(product))
	}

	return result, nil
}

func (p *ProductRepository) DeleteProduct(ctx context.Context, id uint32) error {
	err := p.queries.DeleteProduct(ctx, id)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to delete product: %v", err)
	}

	return nil
}

func toRepositoryProduct(product generated.Product) *repository.Product {
	return &repository.Product{
		ID:              product.ID,
		Name:            product.Name,
		Description:     product.Description,
		RegularPrice:    product.RegularPrice,
		DiscountedPrice: product.DiscountedPrice,
		Quantity:        product.Quantity,
		CategoryID:      product.CategoryID,
		SizeOption:      product.SizeOption,
		ColorOption:     product.ColorOption,
		Rating:          product.Rating,
		Seasonal:        product.Seasonal,
		Featured:        product.Featured,
		ImgUrls:         product.ImgUrls,
		UpdatedBy:       product.UpdatedBy,
		UpdatedAt:       product.UpdatedAt,
		CreatedAt:       product.CreatedAt,
	}
}
