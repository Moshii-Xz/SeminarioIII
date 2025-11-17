package payments

import (
	"context"
	"errors"
	"fmt"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type PaymentGatewayRequest struct {
	OrderID     uint
	Amount      float64
	Currency    string
	Description string
	CustomerID  string
	Metadata    map[string]string
}

type PaymentGatewayResponse struct {
	TransactionID string
	Status        PaymentStatus
	Amount        float64
	Currency      string
	Message       string
	RawResponse   interface{}
}

type RefundRequest struct {
	TransactionID string
	Amount        float64
	Reason        string
}

type RefundResponse struct {
	RefundID    string
	Status      PaymentStatus
	Amount      float64
	Message     string
	RawResponse interface{}
}

type Gateway interface {
	ProcessPayment(ctx context.Context, req PaymentGatewayRequest) (*PaymentGatewayResponse, error)
	GetPaymentStatus(ctx context.Context, transactionID string) (*PaymentGatewayResponse, error)
	Refund(ctx context.Context, req RefundRequest) (*RefundResponse, error)
	VerifyWebhook(ctx context.Context, payload []byte, signature string) (bool, error)
	ParseWebhook(ctx context.Context, payload []byte) (*PaymentGatewayResponse, error)
}

type MockGateway struct {
	apiKey    string
	apiSecret string
	baseURL   string
}

func NewMockGateway(apiKey, apiSecret, baseURL string) *MockGateway {
	return &MockGateway{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   baseURL,
	}
}

func (g *MockGateway) ProcessPayment(ctx context.Context, req PaymentGatewayRequest) (*PaymentGatewayResponse, error) {
	if req.Amount <= 0 {
		return nil, errors.New("el monto debe ser mayor a cero")
	}

	if req.Currency == "" {
		req.Currency = "USD"
	}

	transactionID := fmt.Sprintf("mock_txn_%d_%d", req.OrderID, req.Amount)

	return &PaymentGatewayResponse{
		TransactionID: transactionID,
		Status:        PaymentStatusCompleted,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Message:       "Pago procesado exitosamente (mock)",
		RawResponse: map[string]interface{}{
			"transaction_id": transactionID,
			"status":         "completed",
			"gateway":        "mock",
		},
	}, nil
}

func (g *MockGateway) GetPaymentStatus(ctx context.Context, transactionID string) (*PaymentGatewayResponse, error) {
	if transactionID == "" {
		return nil, errors.New("transaction ID es requerido")
	}

	return &PaymentGatewayResponse{
		TransactionID: transactionID,
		Status:        PaymentStatusCompleted,
		Message:       "Pago completado (mock)",
		RawResponse: map[string]interface{}{
			"transaction_id": transactionID,
			"status":         "completed",
			"gateway":        "mock",
		},
	}, nil
}

func (g *MockGateway) Refund(ctx context.Context, req RefundRequest) (*RefundResponse, error) {
	if req.TransactionID == "" {
		return nil, errors.New("transaction ID es requerido")
	}

	if req.Amount <= 0 {
		return nil, errors.New("el monto del reembolso debe ser mayor a cero")
	}

	refundID := fmt.Sprintf("mock_refund_%s", req.TransactionID)

	return &RefundResponse{
		RefundID: refundID,
		Status:   PaymentStatusRefunded,
		Amount:   req.Amount,
		Message:  "Reembolso procesado exitosamente (mock)",
		RawResponse: map[string]interface{}{
			"refund_id": refundID,
			"status":    "refunded",
			"gateway":   "mock",
		},
	}, nil
}

func (g *MockGateway) VerifyWebhook(ctx context.Context, payload []byte, signature string) (bool, error) {
	// En un gateway real, esto verificaría la firma del webhook --> Si quieres hacerlo me avisas.
	// Para el mock, siempre retornamos true si hay una firma
	if signature == "" {
		return false, errors.New("signature es requerida")
	}

	return true, nil
}

// ParseWebhook parsea un webhook simulado
func (g *MockGateway) ParseWebhook(ctx context.Context, payload []byte) (*PaymentGatewayResponse, error) {
	if len(payload) == 0 {
		return nil, errors.New("payload vacío")
	}

	// En un gateway real, esto parsearía el payload del webhook
	// Para el mock, retornamos una respuesta genérica
	return &PaymentGatewayResponse{
		TransactionID: "mock_webhook_txn",
		Status:        PaymentStatusCompleted,
		Message:       "Webhook procesado (mock)",
		RawResponse: map[string]interface{}{
			"gateway": "mock",
			"source":  "webhook",
		},
	}, nil
}

type GatewayFactory struct{}

func NewGatewayFactory() *GatewayFactory {
	return &GatewayFactory{}
}

func (f *GatewayFactory) CreateGateway(gatewayType string, config map[string]string) (Gateway, error) {
	switch gatewayType {
	case "mock":
		return NewMockGateway(
			config["api_key"],
			config["api_secret"],
			config["base_url"],
		), nil
	case "stripe":
		// TODO: Implementar StripeGateway cuando se necesite
		return nil, fmt.Errorf("gateway tipo '%s' no implementado aún", gatewayType)
	case "paypal":
		// TODO: Implementar PayPalGateway cuando se necesite
		return nil, fmt.Errorf("gateway tipo '%s' no implementado aún", gatewayType)
	default:
		return nil, fmt.Errorf("tipo de gateway desconocido: %s", gatewayType)
	}
}
