package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
	"github.com/gin-gonic/gin"
)

type cartResponse struct {
	ID   uint32          `json:"id"`
	Data []productInCart `json:"data"`
}

type productInCart struct {
	ProductName     string   `json:"product_name"`
	ProductDesc     string   `json:"product_desc"`
	ProductColor    []string `json:"product_color"`
	ProductSize     []string `json:"product_size"`
	ImgUrls         []string `json:"img_urls"`
	Quantity        uint32   `json:"quantity"`
	RegularPrice    float64  `json:"regular_price"`
	DiscountedPrice float64  `json:"discounted_price"`
}

type createCart struct {
	ProductID uint32 `json:"product_id"`
	Quantity  uint32 `json:"quantity"`
}

// type createCartRequest struct {
// 	Products []createCart `json:"products"`
// }

func (s *HttpServer) createCart(ctx *gin.Context) {
	// get id from token
	var userId uint32 = 1

	var req []createCart
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	var data []*repository.Cart

	// create cart
	for _, product := range req {
		cart, err := s.repo.cart.CreateCart(ctx, &repository.Cart{
			UserID:    userId,
			ProductID: product.ProductID,
			Quantity:  product.Quantity,
		})
		if err != nil {
			ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

			return
		}

		data = append(data, cart)
	}

	rsp, err := s.structureCart(ctx, data)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	rsp.ID = userId

	ctx.JSON(http.StatusOK, rsp)
}

func (s *HttpServer) updateCart(ctx *gin.Context) {
	// get id from token
	var userId uint32 = 3

	body, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	var req []createCart
	if err := json.Unmarshal(body, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	// create cart
	for _, product := range req {
		err := s.repo.cart.UpdateCart(ctx, product.Quantity, userId, product.ProductID)
		if err != nil {
			ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

			return
		}
	}

	carts, err := s.repo.cart.ListUserCarts(ctx, userId)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	rsp, err := s.structureCart(ctx, carts)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	rsp.ID = userId

	ctx.JSON(http.StatusOK, rsp)
}

func (s *HttpServer) listCarts(ctx *gin.Context) {
	usersCart, err := s.repo.cart.ListCarts(ctx)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	var listRsps []cartResponse

	for _, cart := range usersCart {
		rsp, err := s.structureCart(ctx, cart.Products)
		if err != nil {
			ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

			return
		}

		rsp.ID = cart.UserID
		listRsps = append(listRsps, rsp)
	}

	ctx.JSON(http.StatusOK, listRsps)
}

func (s *HttpServer) getCart(ctx *gin.Context) {
	// get id from token
	var userId uint32 = 3

	carts, err := s.repo.cart.ListUserCarts(ctx, userId)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	rsp, err := s.structureCart(ctx, carts)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	rsp.ID = userId

	ctx.JSON(http.StatusOK, rsp)
}

func (s *HttpServer) deleteCart(ctx *gin.Context) {
	var userId uint32 = 3

	err := s.repo.cart.DeleteCart(ctx, userId)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (s *HttpServer) structureCart(ctx *gin.Context, data []*repository.Cart) (cartResponse, error) {
	var result cartResponse

	for _, cart := range data {
		// get the product
		product, err := s.repo.p.GetProduct(ctx, cart.ProductID)
		if err != nil {
			return cartResponse{}, err
		}

		size, color, imgUrls, err := product.UnmarshalOptions()
		if err != nil {
			return cartResponse{}, pkg.Errorf(pkg.INTERNAL_ERROR, "could unmarshal products: %v", err)
		}

		result.Data = append(result.Data, productInCart{
			ProductName:     product.Name,
			ProductDesc:     product.Description,
			ProductColor:    color,
			ProductSize:     size,
			ImgUrls:         imgUrls,
			Quantity:        cart.Quantity,
			RegularPrice:    product.RegularPrice,
			DiscountedPrice: product.DiscountedPrice,
		})
	}

	return result, nil
}
