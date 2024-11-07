package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
	"github.com/gin-gonic/gin"
)

// shopping cart > billing info > shipping info > shipping method > preview order > payment > confirmation

type orderResponse struct {
	OrderID         uint32              `json:"order_id"`
	UserID          uint32              `json:"user_id"`
	Amount          float64             `json:"amount"`
	ShippingAmount  float64             `json:"shipping_amount"`
	ShippingAddress string              `json:"shipping_address"`
	Status          string              `json:"status"`
	Data            []orderItemResponse `json:"data"`
	UpdatedBy       uint32              `json:"updated_by"`
	UpdatedAt       time.Time           `json:"updated_at"`
	CreatedAt       time.Time           `json:"created_at"`
}

type orderItemResponse struct {
	ProductID          uint32  `json:"product_id"`
	ProductName        string  `json:"product_name"`
	ProductDescription string  `json:"product_description"`
	Quantity           uint32  `json:"quantity"`
	Price              float64 `json:"price"`
	Color              *string `json:"color"`
	Size               *string `json:"size"`
}

type createOrderRequest struct {
	Amount          float64             `binding:"required"                    json:"amount"`
	ShippingAddress string              `binding:"required"                    json:"shipping_address"`
	ShippingAmount  float64             `binding:"required"                    json:"shipping_amount"`
	OrderItems      []orderItemsRequest `binding:"required"                    json:"order_items"`
	PaymentMethod   string              `binding:"required,oneof=MPESA STRIPE" json:"payment_method"`
}

// [{product_id: 1, quantity: 2, price: 300, color: red, size: 32}, {product_id: 1, quantity: 2, price: 300, color: red, size: 32}]

type orderItemsRequest struct {
	ProductID uint32  `binding:"required" json:"product_id"`
	Quantity  uint32  `binding:"required" json:"quantity"`
	Price     float64 `binding:"required" json:"price"`
	Color     string  `                   json:"color"`
	Size      string  `                   json:"size"`
}

func (s *HttpServer) createOrder(ctx *gin.Context) {
	payload, err := getPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	userId, err := getParam(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	if userId != payload.UserID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(pkg.Errorf(pkg.AUTHENTICATION_ERROR, "cannot create another user order")))

		return
	}
	// create order
	var req createOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	order := &repository.Order{
		UserID:          payload.UserID,
		Amount:          req.Amount,
		ShippingAmount:  req.ShippingAmount,
		ShippingAddress: req.ShippingAddress,
		UpdatedBy:       payload.UserID,
	}

	orderItems := []*repository.OrderItem{}
	for _, orderItem := range req.OrderItems {
		orderItems = append(orderItems, &repository.OrderItem{
			ProductID: orderItem.ProductID,
			Quantity:  orderItem.Quantity,
			Price:     orderItem.Price,
			Color:     pkg.StringPtr(orderItem.Color),
			Size:      pkg.StringPtr(orderItem.Size),
		})
	}

	orderCreated, err := s.repo.o.CreateOrder(ctx, order, orderItems)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	// send payment (MPESA or STRIPE)
	log.Println("Sending stk or something")

	ctx.JSON(http.StatusOK, orderCreated)
}

func (s *HttpServer) getOrder(ctx *gin.Context) {
	payload, err := getPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	orderId, err := getParam(ctx.Param("orderId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	order, err := s.repo.o.GetOrder(ctx, orderId)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	if order.UserID != payload.UserID && payload.Role != "ADMIN" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(pkg.Errorf(pkg.AUTHENTICATION_ERROR, "unauthorized to view this order")))

		return
	}

	rsp, err := s.structureOrderResponse(ctx, order)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (s *HttpServer) listOrders(ctx *gin.Context) {
	payload, err := getPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	if isAdmin, err := isAdmin(payload); !isAdmin {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))

		return
	}

	orders, err := s.repo.o.ListOrders(ctx)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	rsp := []*orderResponse{}

	for _, order := range orders {
		result, err := s.structureOrderResponse(ctx, order)
		if err != nil {
			ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

			return
		}

		rsp = append(rsp, result)
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (s *HttpServer) listUserOrders(ctx *gin.Context) {
	id, err := getParam(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	orders, err := s.repo.o.ListUserOrders(ctx, id)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	rsp := []*orderResponse{}

	for _, order := range orders {
		result, err := s.structureOrderResponse(ctx, order)
		if err != nil {
			ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

			return
		}

		rsp = append(rsp, result)
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (s *HttpServer) listOrderWithStatus(ctx *gin.Context) {
	payload, err := getPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	if isAdmin, err := isAdmin(payload); !isAdmin {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))

		return
	}

	orderStatus := ctx.Query("type")
	orderStatus = strings.ToUpper(orderStatus)

	if orderStatus != "PENDING" && orderStatus != "PROCESSING" && orderStatus != "SHIPPED" && orderStatus != "DELIVERED" {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "unknown order")))

		return
	}

	orders, err := s.repo.o.ListOrderWithStatus(ctx, orderStatus)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	rsp := []*orderResponse{}

	for _, order := range orders {
		result, err := s.structureOrderResponse(ctx, order)
		if err != nil {
			ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

			return
		}

		rsp = append(rsp, result)
	}

	ctx.JSON(http.StatusOK, rsp)
}

type updateOrderStatusRequest struct {
	Status string `binding:"required,oneof=PENDING PROCESSING SHIPPED DELIVERED" json:"status"`
}

func (s *HttpServer) updateOrderStatus(ctx *gin.Context) {
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

	body, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	var req updateOrderStatusRequest
	if err := json.Unmarshal(body, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	if req.Status != "PENDING" && req.Status != "PROCESSING" && req.Status != "SHIPPED" && req.Status != "DELIVERED" {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "unknown order")))

		return
	}

	updateOrder := &repository.UpdateOrder{
		ID:        id,
		Status:    req.Status,
		UpdatedBy: &payload.UserID,
	}

	if err := s.repo.o.UpdateOrder(ctx, updateOrder); err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *HttpServer) deleteOrder(ctx *gin.Context) {
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

	if err := s.repo.o.DeleteOrder(ctx, id); err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *HttpServer) structureOrderResponse(ctx context.Context, order *repository.Order) (*orderResponse, error) {
	result := orderResponse{
		OrderID:         order.ID,
		UserID:          order.UserID,
		Amount:          order.Amount,
		ShippingAmount:  order.ShippingAmount,
		ShippingAddress: order.ShippingAddress,
		Status:          order.Status,
		UpdatedBy:       order.UpdatedBy,
		UpdatedAt:       order.UpdatedAt,
		CreatedAt:       order.CreatedAt,
	}

	orderItems, err := s.repo.o.ListOrderOrderItems(ctx, result.OrderID)
	if err != nil {
		return nil, err
	}

	data := []orderItemResponse{}

	for _, orderItem := range orderItems {
		product, err := s.repo.p.GetProduct(ctx, orderItem.ProductID)
		if err != nil {
			return nil, err
		}

		rsp := orderItemResponse{
			ProductID:          orderItem.ProductID,
			ProductName:        product.Name,
			ProductDescription: product.Description,
			Quantity:           orderItem.Quantity,
			Price:              orderItem.Price,
		}

		if orderItem.Color != nil {
			rsp.Color = orderItem.Color
		}

		if orderItem.Size != nil {
			rsp.Size = orderItem.Size
		}

		data = append(data, rsp)
	}

	result.Data = data

	return &result, nil
}
