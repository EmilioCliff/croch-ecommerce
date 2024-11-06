package handlers

import (
	"time"

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
	UpdatedAt       time.Time           `json:"updated_at"`
	CreatedAt       time.Time           `json:"created_at"`
}

type orderItemResponse struct {
	ProductID          uint32  `json:"product_id"`
	ProductName        string  `json:"product_name"`
	ProductDescription string  `json:"product_description"`
	ProductCategory    string  `json:"product_category"`
	Quantity           uint32  `json:"quantity"`
	Price              float64 `json:"price"`
	Color              string  `json:"color"`
	Size               string  `json:"size"`
}

func (s *HttpServer) createOrder(ctx *gin.Context) {}

func (s *HttpServer) getOrder(ctx *gin.Context) {}

func (s *HttpServer) listOrders(ctx *gin.Context) {}

func (s *HttpServer) listUserOrders(ctx *gin.Context) {}

func (s *HttpServer) listOrderWithStatus(ctx *gin.Context) {}

func (s *HttpServer) updateOrderStatus(ctx *gin.Context) {}

func (s *HttpServer) deleteOrder(ctx *gin.Context) {}

func (s *HttpServer) structureResponse() {}
