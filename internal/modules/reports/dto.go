package reports

import "time"

type SalesSummaryRequest struct {
	StartDate time.Time `form:"fecha_inicio" binding:"required"`
	EndDate   time.Time `form:"fecha_fin" binding:"required"`
}

type SalesSummaryResponse struct {
	TotalOrders       int64   `json:"total_ordenes"`
	TotalRevenue      float64 `json:"total_ingresos"`
	AverageOrderValue float64 `json:"ticket_promedio"`
	TotalItemsSold    int64   `json:"total_items_vendidos"`
}

type TopProductResponse struct {
	ProductID   uint    `json:"id_producto"`
	ProductName string  `json:"nombre_producto"`
	UnitsSold   int64   `json:"unidades_vendidas"`
	Revenue     float64 `json:"ingresos"`
}

type InventoryStatusResponse struct {
	ProductID   uint   `json:"id_producto"`
	ProductName string `json:"nombre_producto"`
	Stock       int    `json:"stock"`
}

type DailySalesResponse struct {
	Date        time.Time `json:"fecha"`
	TotalOrders int64     `json:"total_ordenes"`
	Revenue     float64   `json:"ingresos"`
}

type CustomerRankingResponse struct {
	CustomerID   uint    `json:"id_cliente"`
	CustomerName string  `json:"nombre_cliente"`
	OrdersCount  int64   `json:"total_ordenes"`
	TotalSpent   float64 `json:"total_gastado"`
}

type PaymentMethodSummaryResponse struct {
	MethodID    *uint   `json:"id_metodo_pago,omitempty"`
	MethodName  *string `json:"nombre_metodo_pago,omitempty"`
	TotalAmount float64 `json:"total_pagado"`
	Payments    int64   `json:"cantidad_pagos"`
}

type PendingPaymentResponse struct {
	OrderID       uint      `json:"id_orden"`
	CustomerName  string    `json:"nombre_cliente"`
	OrderTotal    float64   `json:"total_orden"`
	TotalPaid     float64   `json:"total_pagado"`
	PendingAmount float64   `json:"pendiente"`
	OrderDate     time.Time `json:"fecha_orden"`
}

type ReportFilterRequest struct {
	StartDate time.Time `form:"fecha_inicio" binding:"required"`
	EndDate   time.Time `form:"fecha_fin" binding:"required"`
	Limit     int       `form:"limite" binding:"omitempty,min=1,max=100"`
}

