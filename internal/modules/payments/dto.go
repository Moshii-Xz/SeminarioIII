package payments

import "time"

type CreatePaymentRequest struct {
	OrderID         uint    `json:"order_id" binding:"required"`
	PaymentMethodID *uint   `json:"payment_method_id" binding:"omitempty"`
	Amount          float64 `json:"amount" binding:"required,min=0"`
}

type UpdatePaymentRequest struct {
	PaymentMethodID *uint   `json:"payment_method_id" binding:"omitempty"`
	Amount          float64 `json:"amount" binding:"omitempty,min=0"`
}

type PaymentResponse struct {
	ID              uint      `json:"id"`
	OrderID         uint      `json:"order_id"`
	PaymentMethodID *uint     `json:"payment_method_id,omitempty"`
	Amount          float64   `json:"amount"`
	Date            time.Time `json:"date"`
}

type PaymentListResponse struct {
	Payments []PaymentResponse `json:"payments"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

type PaymentMethodResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PaymentMethodListResponse struct {
	Methods []PaymentMethodResponse `json:"methods"`
	Total   int64                   `json:"total"`
}

type PaymentByOrderResponse struct {
	OrderID    uint              `json:"order_id"`
	OrderTotal float64           `json:"order_total"`
	TotalPaid  float64           `json:"total_paid"`
	Pending    float64           `json:"pending"`
	Payments   []PaymentResponse `json:"payments"`
}
