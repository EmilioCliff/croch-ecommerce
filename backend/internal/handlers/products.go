package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
	"github.com/gin-gonic/gin"
)

type createProductRequest struct {
	Name            string   `binding:"required" json:"name"`
	Description     string   `binding:"required" json:"description"`
	RegularPrice    float64  `binding:"required" json:"regular_price"`
	DiscountedPrice float64  `binding:""         json:"discounted_price"`
	Quantity        uint32   `binding:""         json:"quantity"`
	CategoryID      uint32   `binding:"required" json:"category_id"`
	SizeOption      []string `binding:""         json:"size_option"`
	ColorOption     []string `binding:""         json:"color_option"`
	Seasonal        bool     `binding:""         json:"seasonal"`
	Featured        bool     `binding:""         json:"featured"`
	ImgUrls         []string `binding:""         json:"img_urls"`
}

func (s *HttpServer) createProduct(ctx *gin.Context) {
	payload, err := getPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	if isAdmin, err := isAdmin(payload); !isAdmin {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))

		return
	}

	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	if req.SizeOption == nil {
		req.SizeOption = []string{}
	}

	if req.ColorOption == nil {
		req.ColorOption = []string{""}
	}

	if req.ImgUrls == nil {
		req.ImgUrls = []string{}
	}

	reqProduct := &repository.Product{
		Name:            req.Name,
		Description:     req.Description,
		RegularPrice:    req.RegularPrice,
		DiscountedPrice: req.DiscountedPrice,
		Quantity:        req.Quantity,
		CategoryID:      req.CategoryID,
		Seasonal:        req.Seasonal,
		Featured:        req.Featured,
		UpdatedBy:       payload.UserID,
	}
	if err := reqProduct.MarshalOptions(req.SizeOption, req.ColorOption, req.ImgUrls); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	product, err := s.repo.p.CreateProduct(ctx, reqProduct)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusCreated, product)
}

func (s *HttpServer) updateProduct(ctx *gin.Context) {
	payload, err := getPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	if isAdmin, err := isAdmin(payload); !isAdmin {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))

		return
	}

	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	productId, err := getParam(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	reqProduct := &repository.UpdateProduct{
		UpdatedBy:       payload.UserID,
		ID:              productId,
		Name:            pkg.StringPtr(req.Name),
		Description:     pkg.StringPtr(req.Description),
		RegularPrice:    pkg.Float64Ptr(req.RegularPrice),
		DiscountedPrice: pkg.Float64Ptr(req.DiscountedPrice),
		Quantity:        pkg.Uint32Ptr(req.Quantity),
		CategoryID:      pkg.Uint32Ptr(req.CategoryID),
		Seasonal:        pkg.BoolPtr(req.Seasonal),
		Featured:        pkg.BoolPtr(req.Featured),
	}

	if err := reqProduct.MarshalOptions(req.SizeOption, req.ColorOption, req.ImgUrls); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	err = s.repo.p.UpdateProduct(ctx, reqProduct)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, reqProduct)
}

type updateProductQuantityRequest struct {
	Quantity uint32 `binding:"required" json:"quantity"`
}

func (s *HttpServer) updateProductQuantity(ctx *gin.Context) {
	payload, err := getPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	if isAdmin, err := isAdmin(payload); !isAdmin {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))

		return
	}

	body, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	var req updateProductQuantityRequest
	if err := json.Unmarshal(body, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	id, err := getParam(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	err = s.repo.p.UpdateProductQuantity(ctx, id, req.Quantity)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (s *HttpServer) getProduct(ctx *gin.Context) {
	id, err := getParam(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	product, err := s.repo.p.GetProduct(ctx, id)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (s *HttpServer) listProducts(ctx *gin.Context) {
	productType := ctx.Query("type")

	var result []*repository.Product

	var err error

	switch productType {
	case "new":
		result, err = s.repo.p.ListNewProducts(ctx)
	case "seasonal":
		result, err = s.repo.p.ListSeasonalProducts(ctx)
	case "featured":
		result, err = s.repo.p.ListFeaturedProducts(ctx)
	case "discounted":
		result, err = s.repo.p.ListDiscountedProducts(ctx)
	default:
		result, err = s.repo.p.ListProducts(ctx)
	}

	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *HttpServer) deleteProduct(ctx *gin.Context) {
	payload, err := getPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	if isAdmin, err := isAdmin(payload); !isAdmin {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))

		return
	}

	id, err := getParam(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	err = s.repo.p.DeleteProduct(ctx, id)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
