package payments

import (
	"errors"
	"fmt"
	"time"

	"github.com/mordmora/expirapp/internal/domain"
	ordersRepo "github.com/mordmora/expirapp/internal/modules/orders"
)

type Service struct {
	repo       *Repository
	ordersRepo *ordersRepo.Repository
}

func NewService(repo *Repository, ordersRepo *ordersRepo.Repository) *Service {
	return &Service{
		repo:       repo,
		ordersRepo: ordersRepo,
	}
}

func (s *Service) Create(req CreatePaymentRequest) (*domain.Payment, error) {
	order, err := s.ordersRepo.FindByID(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("order with id %d not found", req.OrderID)
	}

	var orderTotal float64
	for _, item := range order.Details {
		orderTotal += float64(item.Quantity) * item.UnitPrice
	}

	totalPaid, err := s.repo.GetTotalPaidByOrderID(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("error calculating total paid: %w", err)
	}

	pending := orderTotal - totalPaid
	if req.Amount > pending {
		return nil, fmt.Errorf("amount exceeds pending balance. Requested: %.2f, Pending: %.2f", req.Amount, pending)
	}

	if req.PaymentMethodID != nil {
		_, err := s.repo.FindPaymentMethodByID(*req.PaymentMethodID)
		if err != nil {
			return nil, fmt.Errorf("payment method with id %d not found", *req.PaymentMethodID)
		}
	}

	payment := &domain.Payment{
		OrderID:         req.OrderID,
		PaymentMethodID: req.PaymentMethodID,
		Amount:          req.Amount,
		Date:            time.Now(),
	}

	if err := s.repo.Create(payment); err != nil {
		return nil, fmt.Errorf("error creating payment: %w", err)
	}

	return payment, nil
}

func (s *Service) GetById(id uint) (*domain.Payment, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Update(id uint, req UpdatePaymentRequest) (*domain.Payment, error) {
	payment, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Amount > 0 {
		order, err := s.ordersRepo.FindByID(payment.OrderID)
		if err != nil {
			return nil, fmt.Errorf("error getting order: %w", err)
		}

		var orderTotal float64
		for _, item := range order.Details {
			orderTotal += float64(item.Quantity) * item.UnitPrice
		}

		totalPaid, err := s.repo.GetTotalPaidByOrderID(payment.OrderID)
		if err != nil {
			return nil, fmt.Errorf("error calculating total paid: %w", err)
		}

		newTotalPaid := totalPaid - payment.Amount + req.Amount
		if newTotalPaid > orderTotal {
			return nil, fmt.Errorf("new amount would exceed order total. Order total: %.2f, New total paid: %.2f", orderTotal, newTotalPaid)
		}

		payment.Amount = req.Amount
	}

	if req.PaymentMethodID != nil {
		_, err := s.repo.FindPaymentMethodByID(*req.PaymentMethodID)
		if err != nil {
			return nil, fmt.Errorf("payment method with id %d not found", *req.PaymentMethodID)
		}
		payment.PaymentMethodID = req.PaymentMethodID
	}

	if err := s.repo.Update(payment); err != nil {
		return nil, fmt.Errorf("error updating payment: %w", err)
	}

	return payment, nil
}

func (s *Service) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *Service) List(page, limit int) ([]domain.Payment, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.repo.List(limit, offset)
}

func (s *Service) GetByOrderID(orderID uint) ([]domain.Payment, error) {
	return s.repo.FindByOrderID(orderID)
}

func (s *Service) GetPaymentStatusByOrderID(orderID uint) (*PaymentByOrderResponse, error) {
	order, err := s.ordersRepo.FindByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("order with id %d not found", orderID)
	}

	var orderTotal float64
	for _, item := range order.Details {
		orderTotal += float64(item.Quantity) * item.UnitPrice
	}

	payments, err := s.repo.FindByOrderID(orderID)
	if err != nil {
		return nil, fmt.Errorf("error getting payments: %w", err)
	}

	totalPaid, err := s.repo.GetTotalPaidByOrderID(orderID)
	if err != nil {
		return nil, fmt.Errorf("error calculating total paid: %w", err)
	}

	pending := orderTotal - totalPaid
	if pending < 0 {
		pending = 0
	}

	paymentResponses := make([]PaymentResponse, len(payments))
	for i, payment := range payments {
		paymentResponses[i] = s.ToPaymentResponse(&payment)
	}

	return &PaymentByOrderResponse{
		OrderID:    orderID,
		OrderTotal: orderTotal,
		TotalPaid:  totalPaid,
		Pending:    pending,
		Payments:   paymentResponses,
	}, nil
}

func (s *Service) CreatePaymentMethod(name string) (*domain.PaymentMethod, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	existing, _ := s.repo.FindPaymentMethodByName(name)
	if existing != nil {
		return nil, errors.New("payment method already exists")
	}

	method := &domain.PaymentMethod{
		Name: name,
	}

	if err := s.repo.CreatePaymentMethod(method); err != nil {
		return nil, fmt.Errorf("error creating payment method: %w", err)
	}

	return method, nil
}

func (s *Service) GetPaymentMethodByID(id uint) (*domain.PaymentMethod, error) {
	return s.repo.FindPaymentMethodByID(id)
}

func (s *Service) ListPaymentMethods() ([]domain.PaymentMethod, int64, error) {
	return s.repo.ListPaymentMethods()
}

func (s *Service) UpdatePaymentMethod(id uint, name string) (*domain.PaymentMethod, error) {
	method, err := s.repo.FindPaymentMethodByID(id)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.FindPaymentMethodByName(name)
	if err == nil && existing.ID != id {
		return nil, errors.New("payment method already exists")
	}

	method.Name = name

	if err := s.repo.UpdatePaymentMethod(method); err != nil {
		return nil, fmt.Errorf("error updating payment method: %w", err)
	}

	return method, nil
}

func (s *Service) DeletePaymentMethod(id uint) error {
	_, err := s.repo.FindPaymentMethodByID(id)
	if err != nil {
		return err
	}

	return s.repo.DeletePaymentMethod(id)
}

func (s *Service) ToPaymentResponse(payment *domain.Payment) PaymentResponse {
	return PaymentResponse{
		ID:              payment.ID,
		OrderID:         payment.OrderID,
		PaymentMethodID: payment.PaymentMethodID,
		Amount:          payment.Amount,
		Date:            payment.Date,
	}
}

func (s *Service) ToPaymentMethodResponse(method *domain.PaymentMethod) PaymentMethodResponse {
	return PaymentMethodResponse{
		ID:   method.ID,
		Name: method.Name,
	}
}
