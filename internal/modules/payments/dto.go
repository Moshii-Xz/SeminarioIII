package payments

import "time"

type CreatePaymentRequest struct {
	IDCompra     uint    `json:"id_compra" binding:"required"`
	IDMetodoPago *uint   `json:"id_metodo_pago" binding:"omitempty"`
	Monto        float64 `json:"monto" binding:"required,min=0"`
}

type UpdatePaymentRequest struct {
	IDMetodoPago *uint   `json:"id_metodo_pago" binding:"omitempty"`
	Monto        float64 `json:"monto" binding:"omitempty,min=0"`
}

type PaymentResponse struct {
	IDPago       uint      `json:"id_pago"`
	IDCompra     uint      `json:"id_compra"`
	IDMetodoPago *uint     `json:"id_metodo_pago,omitempty"`
	Monto        float64   `json:"monto"`
	FechaPago    time.Time `json:"fecha_pago"`
}

type PaymentListResponse struct {
	Payments []PaymentResponse `json:"pagos"`
	Total    int64             `json:"total"`
	Page     int               `json:"pagina"`
	Limit    int               `json:"limite"`
}

type PaymentMethodResponse struct {
	IDMetodoPago uint   `json:"id_metodo_pago"`
	Nombre       string `json:"nombre"`
}

type PaymentMethodListResponse struct {
	Methods []PaymentMethodResponse `json:"metodos_pago"`
	Total   int64                   `json:"total"`
}

type PaymentByOrderResponse struct {
	IDCompra      uint                `json:"id_compra"`
	TotalOrden    float64             `json:"total_orden"`
	TotalPagado   float64             `json:"total_pagado"`
	Pendiente     float64             `json:"pendiente"`
	Payments      []PaymentResponse   `json:"pagos"`
}
