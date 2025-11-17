package payments

import (
	"errors"
	"fmt"
	"time"

	"github.com/mordmora/expirapp/internal/domain"
	ordersRepo "github.com/mordmora/expirapp/internal/modules/orders"
)

type Service struct {
	repo        *Repository
	ordersRepo  *ordersRepo.Repository
}

func NewService(repo *Repository, ordersRepo *ordersRepo.Repository) *Service {
	return &Service{
		repo:       repo,
		ordersRepo: ordersRepo,
	}
}

func (s *Service) Create(req CreatePaymentRequest) (*domain.Payment, error) {
	order, err := s.ordersRepo.FindByID(req.IDCompra)
	if err != nil {
		return nil, fmt.Errorf("orden con id %d no encontrada", req.IDCompra)
	}

	var orderTotal float64
	for _, item := range order.Items {
		orderTotal += float64(item.Cantidad) * item.PrecioUnitario
	}

	totalPaid, err := s.repo.GetTotalPaidByOrderID(req.IDCompra)
	if err != nil {
		return nil, fmt.Errorf("error calculando total pagado: %w", err)
	}

	pending := orderTotal - totalPaid
	if req.Monto > pending {
		return nil, fmt.Errorf("el monto excede el pendiente. Monto solicitado: %.2f, Pendiente: %.2f", req.Monto, pending)
	}

	if req.IDMetodoPago != nil {
		_, err := s.repo.FindPaymentMethodByID(*req.IDMetodoPago)
		if err != nil {
			return nil, fmt.Errorf("método de pago con id %d no encontrado", *req.IDMetodoPago)
		}
	}

	payment := &domain.Payment{
		IDCompra:     req.IDCompra,
		IDMetodoPago: req.IDMetodoPago,
		Monto:        req.Monto,
		FechaPago:    time.Now(),
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

	if req.Monto > 0 {
		order, err := s.ordersRepo.FindByID(payment.IDCompra)
		if err != nil {
			return nil, fmt.Errorf("error obteniendo orden: %w", err)
		}

		var orderTotal float64
		for _, item := range order.Items {
			orderTotal += float64(item.Cantidad) * item.PrecioUnitario
		}

		totalPaid, err := s.repo.GetTotalPaidByOrderID(payment.IDCompra)
		if err != nil {
			return nil, fmt.Errorf("error calculando total pagado: %w", err)
		}

		newTotalPaid := totalPaid - payment.Monto + req.Monto
		if newTotalPaid > orderTotal {
			return nil, fmt.Errorf("el nuevo monto excedería el total de la orden. Total orden: %.2f, Nuevo total pagado: %.2f", orderTotal, newTotalPaid)
		}

		payment.Monto = req.Monto
	}

	if req.IDMetodoPago != nil {
		_, err := s.repo.FindPaymentMethodByID(*req.IDMetodoPago)
		if err != nil {
			return nil, fmt.Errorf("método de pago con id %d no encontrado", *req.IDMetodoPago)
		}
		payment.IDMetodoPago = req.IDMetodoPago
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
		return nil, fmt.Errorf("orden con id %d no encontrada", orderID)
	}

	var orderTotal float64
	for _, item := range order.Items {
		orderTotal += float64(item.Cantidad) * item.PrecioUnitario
	}

	payments, err := s.repo.FindByOrderID(orderID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo pagos: %w", err)
	}

	totalPaid, err := s.repo.GetTotalPaidByOrderID(orderID)
	if err != nil {
		return nil, fmt.Errorf("error calculando total pagado: %w", err)
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
		IDCompra:    orderID,
		TotalOrden:  orderTotal,
		TotalPagado: totalPaid,
		Pendiente:   pending,
		Payments:    paymentResponses,
	}, nil
}

func (s *Service) CreatePaymentMethod(nombre string) (*domain.PaymentMethod, error) {
	_, err := s.repo.FindPaymentMethodByName(nombre)
	if err == nil {
		return nil, errors.New("ya existe un método de pago con ese nombre")
	}

	method := &domain.PaymentMethod{
		Nombre: nombre,
	}

	if err := s.repo.CreatePaymentMethod(method); err != nil {
		return nil, fmt.Errorf("error creating payment method: %w", err)
	}

	return method, nil
}

func (s *Service) GetPaymentMethodByID(id uint) (*domain.PaymentMethod, error) {
	return s.repo.FindPaymentMethodByID(id)
}

func (s *Service) ListPaymentMethods() ([]domain.PaymentMethod, error) {
	return s.repo.ListPaymentMethods()
}

func (s *Service) UpdatePaymentMethod(id uint, nombre string) (*domain.PaymentMethod, error) {
	method, err := s.repo.FindPaymentMethodByID(id)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.FindPaymentMethodByName(nombre)
	if err == nil && existing.ID != id {
		return nil, errors.New("ya existe un método de pago con ese nombre")
	}

	method.Nombre = nombre

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
		IDPago:       payment.ID,
		IDCompra:     payment.IDCompra,
		IDMetodoPago: payment.IDMetodoPago,
		Monto:        payment.Monto,
		FechaPago:    payment.FechaPago,
	}
}

func (s *Service) ToPaymentMethodResponse(method *domain.PaymentMethod) PaymentMethodResponse {
	return PaymentMethodResponse{
		IDMetodoPago: method.ID,
		Nombre:       method.Nombre,
	}
}
