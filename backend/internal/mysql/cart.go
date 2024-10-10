package mysql

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/mysql/generated"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
	"github.com/google/uuid"
)

var _ repository.CartRepository = (*CartRepository)(nil)

type CartRepository struct {
	db      *Store
	queries generated.Querier
}

func NewCartRepository(store *Store) *CartRepository {
	queries := generated.New(store.db)

	return &CartRepository{
		db:      store,
		queries: queries,
	}
}

func (c *CartRepository) CreateCart(ctx context.Context, cart *repository.Cart) (*repository.Cart, error) {
	if err := cart.Validate(); err != nil {
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "%v", err)
	}

	_, err := c.queries.CreateCart(ctx, generated.CreateCartParams{
		UserID:    cart.UserID.String(),
		ProductID: cart.ProductID.String(),
		Quantity:  cart.Quantity,
	})
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create cart: %v", err)
	}

	return cart, nil
}

func (c *CartRepository) ListCarts(ctx context.Context) ([]*repository.Cart, error) {
	carts, err := c.queries.ListCart(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no cart available")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list cart: %v", err)
	}

	var result []*repository.Cart
	for _, cart := range carts {
		result = append(result, &repository.Cart{
			UserID:    uuid.MustParse(cart.UserID),
			ProductID: uuid.MustParse(cart.ProductID),
			Quantity:  cart.Quantity,
			CreatedAt: cart.CreatedAt,
		})
	}

	return result, nil
}

func (c *CartRepository) ListUserCarts(ctx context.Context, userID uuid.UUID) ([]*repository.Cart, error) {
	carts, err := c.queries.ListUserCarts(ctx, userID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user cart is empty")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list users cart: %v", err)
	}

	var result []*repository.Cart
	for _, cart := range carts {
		result = append(result, &repository.Cart{
			UserID:    uuid.MustParse(cart.UserID),
			ProductID: uuid.MustParse(cart.ProductID),
			Quantity:  cart.Quantity,
			CreatedAt: cart.CreatedAt,
		})
	}

	return result, nil
}

func (c *CartRepository) ListProductInCarts(ctx context.Context, productID uuid.UUID) ([]*repository.Cart, error) {
	carts, err := c.queries.ListProductInCarts(ctx, productID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "product is not in cart")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list products in cart: %v", err)
	}

	var result []*repository.Cart
	for _, cart := range carts {
		result = append(result, &repository.Cart{
			UserID:    uuid.MustParse(cart.UserID),
			ProductID: uuid.MustParse(cart.ProductID),
			Quantity:  cart.Quantity,
			CreatedAt: cart.CreatedAt,
		})
	}

	return result, nil
}

func (c *CartRepository) UpdateCart(ctx context.Context, quantity uint32, userID uuid.UUID, productID uuid.UUID) error {
	err := c.queries.UpdateUserCart(ctx, generated.UpdateUserCartParams{
		Quantity:  quantity,
		UserID:    userID.String(),
		ProductID: productID.String(),
	})
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update cart: %v", err)
	}

	return nil
}

func (c *CartRepository) DeleteCart(ctx context.Context, userID uuid.UUID) error {
	err := c.queries.DeleteUserCart(ctx, userID.String())
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to delete users cart: %v", err)
	}

	return nil
}
