package mysql

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/mysql/generated"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
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
		UserID:    cart.UserID,
		ProductID: cart.ProductID,
		Quantity:  cart.Quantity,
	})
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create cart: %v", err)
	}

	return cart, nil
}

func (c *CartRepository) ListCarts(ctx context.Context) ([]*repository.UserCart, error) {
	carts, err := c.queries.ListCartByUser(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no cart available")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list cart: %v", err)
	}

	var result []*repository.UserCart

	// Iterate through each user's cart data
	for _, cart := range carts {
		if cart.CartItems.Valid {
			// Split concatenated cart items string by the separator " | "
			cartItems := strings.Split(cart.CartItems.String, " | ")

			var userProducts []*repository.Cart

			for _, item := range cartItems {
				parts := strings.Split(item, ", ")

				// Assuming parts[0], parts[1], parts[2] contain Product ID, Quantity, and Created At respectively
				if len(parts) == 3 {
					productIDStr := strings.TrimPrefix(parts[0], "Product ID: ")
					quantityStr := strings.TrimPrefix(parts[1], "Quantity: ")
					createdAtStr := strings.TrimPrefix(parts[2], "Created At: ")

					productID, _ := strconv.ParseUint(productIDStr, 10, 32)
					quantity, _ := strconv.ParseUint(quantityStr, 10, 32)
					createdAt, _ := time.Parse("2006-01-02 15:04:05", createdAtStr)

					userProducts = append(userProducts, &repository.Cart{
						UserID:    cart.UserID,
						ProductID: uint32(productID),
						Quantity:  uint32(quantity),
						CreatedAt: createdAt,
					})
				}
			}

			// Append user and their cart products to the result
			result = append(result, &repository.UserCart{
				UserID:   cart.UserID,
				Products: userProducts,
			})
		}
	}

	return result, nil
}

func (c *CartRepository) ListUserCarts(ctx context.Context, userID uint32) ([]*repository.Cart, error) {
	carts, err := c.queries.ListUserCarts(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user cart is empty")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list users cart: %v", err)
	}

	var result []*repository.Cart
	for _, cart := range carts {
		result = append(result, &repository.Cart{
			UserID:    cart.UserID,
			ProductID: cart.ProductID,
			Quantity:  cart.Quantity,
			CreatedAt: cart.CreatedAt,
		})
	}

	return result, nil
}

func (c *CartRepository) ListProductInCarts(ctx context.Context, productID uint32) ([]*repository.Cart, error) {
	carts, err := c.queries.ListProductInCarts(ctx, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "product is not in cart")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list products in cart: %v", err)
	}

	var result []*repository.Cart
	for _, cart := range carts {
		result = append(result, &repository.Cart{
			UserID:    cart.UserID,
			ProductID: cart.ProductID,
			Quantity:  cart.Quantity,
			CreatedAt: cart.CreatedAt,
		})
	}

	return result, nil
}

func (c *CartRepository) UpdateCart(ctx context.Context, quantity uint32, userID uint32, productID uint32) error {
	err := c.queries.UpdateUserCart(ctx, generated.UpdateUserCartParams{
		Quantity:  quantity,
		UserID:    userID,
		ProductID: productID,
	})
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update cart: %v", err)
	}

	return nil
}

func (c *CartRepository) DeleteCart(ctx context.Context, userID uint32) error {
	err := c.queries.DeleteUserCart(ctx, userID)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to delete users cart: %v", err)
	}

	return nil
}
