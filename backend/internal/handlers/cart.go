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
	Data map[int32]int32 `binding:"required" json:"data"` // {1: 32, 3: 40, 5: 10}
}

func (s *HttpServer) createCart(ctx *gin.Context) {
	payload, err := getPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	var req createCart
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	data := []*repository.Cart{}

	// create cart
	for productId, quantity := range req.Data {
		cart, err := s.repo.cart.CreateCart(ctx, &repository.Cart{
			UserID:    payload.UserID,
			ProductID: uint32(productId),
			Quantity:  uint32(quantity),
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

	rsp.ID = payload.UserID

	ctx.JSON(http.StatusOK, rsp)
}

func (s *HttpServer) updateCart(ctx *gin.Context) {
	payload, err := getPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	body, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	var req createCart
	if err := json.Unmarshal(body, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	// update cart
	for productId, quantity := range req.Data {
		err := s.repo.cart.UpdateCart(ctx, uint32(quantity), payload.UserID, uint32(productId))
		if err != nil {
			ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

			return
		}
	}

	carts, err := s.repo.cart.ListUserCarts(ctx, payload.UserID)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	rsp, err := s.structureCart(ctx, carts)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	rsp.ID = payload.UserID

	ctx.JSON(http.StatusOK, rsp)
}

func (s *HttpServer) listCarts(ctx *gin.Context) {
	usersCart, err := s.repo.cart.ListCarts(ctx)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	listRsps := []cartResponse{}

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
	payload, err := getPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	carts, err := s.repo.cart.ListUserCarts(ctx, payload.UserID)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	rsp, err := s.structureCart(ctx, carts)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	rsp.ID = payload.UserID

	ctx.JSON(http.StatusOK, rsp)
}

func (s *HttpServer) deleteCart(ctx *gin.Context) {
	payload, err := getPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	err = s.repo.cart.DeleteCart(ctx, payload.UserID)
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
