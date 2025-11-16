package orders

import "time"

type CreateOrderRequest struct {
	IDCliente  uint                `json:"id_cliente" binding:"required"`
	IDVendedor *uint               `json:"id_vendedor" binding:"omitempty"`
	Items      []OrderItemRequest  `json:"items" binding:"required,min=1,dive"`
}

type OrderItemRequest struct {
	IDProducto     uint    `json:"id_producto" binding:"required"`
	Cantidad       int     `json:"cantidad" binding:"required,min=1"`
	PrecioUnitario float64 `json:"precio_unitario" binding:"required,min=0"`
}

type UpdateOrderRequest struct {
	IDVendedor *uint `json:"id_vendedor" binding:"omitempty"`
}

type AddOrderItemRequest struct {
	IDProducto     uint    `json:"id_producto" binding:"required"`
	Cantidad       int     `json:"cantidad" binding:"required,min=1"`
	PrecioUnitario float64 `json:"precio_unitario" binding:"required,min=0"`
}

type UpdateOrderItemRequest struct {
	Cantidad       int     `json:"cantidad" binding:"omitempty,min=1"`
	PrecioUnitario float64 `json:"precio_unitario" binding:"omitempty,min=0"`
}

type OrderItemResponse struct {
	IDDetalle      uint    `json:"id_detalle"`
	IDProducto     uint    `json:"id_producto"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
	Subtotal       float64 `json:"subtotal"`
}

type OrderResponse struct {
	IDCompra      uint                `json:"id_compra"`
	IDCliente     uint                `json:"id_cliente"`
	IDVendedor    *uint               `json:"id_vendedor,omitempty"`
	FechaCompra   time.Time           `json:"fecha_compra"`
	Items         []OrderItemResponse `json:"items"`
	Total         float64             `json:"total"`
}

type OrderListResponse struct {
	Orders []OrderResponse `json:"ordenes"`
	Total  int64           `json:"total"`
	Page   int             `json:"pagina"`
	Limit  int             `json:"limite"`
}
