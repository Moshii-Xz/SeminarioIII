package orders

import "time"

type CreateOrderRequest struct {
	ClientID uint               `json:"client_id" binding:"required"`
	StoreID *uint              `json:"store_id" binding:"omitempty"`
	Items    []OrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type OrderItemRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	UnitPrice float64 `json:"unit_price" binding:"required,min=0"`
}

type UpdateOrderRequest struct {
	StoreID *uint  `json:"store_id" binding:"omitempty"`
	Status  string `json:"estado" binding:"omitempty,oneof=En preparación Listo para recoger Entregado"`
}

type AddOrderItemRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	UnitPrice float64 `json:"unit_price" binding:"required,min=0"`
}

type UpdateOrderItemRequest struct {
	Quantity  int     `json:"quantity" binding:"omitempty,min=1"`
	UnitPrice float64 `json:"unit_price" binding:"omitempty,min=0"`
}

type OrderItemResponse struct {
	ID        uint    `json:"id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Subtotal  float64 `json:"subtotal"`
}

type OrderResponse struct {
	ID       uint                `json:"id"`
	ClientID uint                `json:"client_id"`
	StoreID *uint               `json:"store_id,omitempty"`
	Date     time.Time           `json:"date"`
	Status   string              `json:"estado"`
	Items    []OrderItemResponse `json:"items"`
	Total    float64             `json:"total"`
}

type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int64           `json:"total"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
}
