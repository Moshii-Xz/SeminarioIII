package orders

import (
	"errors"
	"fmt"
	"time"

	"github.com/mordmora/expirapp/internal/domain"
	catalogRepo "github.com/mordmora/expirapp/internal/modules/catalog"
)

type Service struct {
	repo        *Repository
	catalogRepo *catalogRepo.Repository
}

func NewService(repo *Repository, catalogRepo *catalogRepo.Repository) *Service {
	return &Service{
		repo:        repo,
		catalogRepo: catalogRepo,
	}
}

func (s *Service) Create(req CreateOrderRequest) (*domain.Order, error) {
	orderItems := make([]domain.OrderItem, len(req.Items))
	
	for i, itemReq := range req.Items {
		product, err := s.catalogRepo.FindByID(itemReq.IDProducto)
		if err != nil {
			return nil, fmt.Errorf("producto con id %d no encontrado", itemReq.IDProducto)
		}

		if product.Stock < itemReq.Cantidad {
			return nil, fmt.Errorf("stock insuficiente para el producto %s (id: %d). Stock disponible: %d, solicitado: %d", 
				product.Nombre, itemReq.IDProducto, product.Stock, itemReq.Cantidad)
		}

		orderItems[i] = domain.OrderItem{
			IDProducto:     itemReq.IDProducto,
			Cantidad:       itemReq.Cantidad,
			PrecioUnitario: itemReq.PrecioUnitario,
		}
	}

	order := &domain.Order{
		IDCliente:   req.IDCliente,
		IDVendedor:  req.IDVendedor,
		FechaCompra: time.Now(),
		Items:       orderItems,
	}

	if err := s.repo.CreateWithItems(order, orderItems); err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	for _, item := range orderItems {
		if err := s.catalogRepo.UpdateStock(item.IDProducto, -item.Cantidad); err != nil {
			return nil, fmt.Errorf("error updating stock for product %d: %w", item.IDProducto, err)
		}
	}

	return order, nil
}

func (s *Service) GetById(id uint) (*domain.Order, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Update(id uint, req UpdateOrderRequest) (*domain.Order, error) {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.IDVendedor != nil {
		order.IDVendedor = req.IDVendedor
	}

	if err := s.repo.Update(order); err != nil {
		return nil, fmt.Errorf("error updating order: %w", err)
	}

	return order, nil
}

func (s *Service) Delete(id uint) error {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		if err := s.catalogRepo.UpdateStock(item.IDProducto, item.Cantidad); err != nil {
			fmt.Printf("warning: error restoring stock for product %d: %v\n", item.IDProducto, err)
		}
	}

	return s.repo.Delete(id)
}

func (s *Service) List(page, limit int) ([]domain.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.repo.List(limit, offset)
}

func (s *Service) ListByClient(clientID uint, page, limit int) ([]domain.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.repo.FindByClientID(clientID, limit, offset)
}

func (s *Service) ListBySeller(sellerID uint, page, limit int) ([]domain.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.repo.FindBySellerID(sellerID, limit, offset)
}

func (s *Service) AddOrderItem(orderID uint, req AddOrderItemRequest) (*domain.OrderItem, error) {
	order, err := s.repo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	product, err := s.catalogRepo.FindByID(req.IDProducto)
	if err != nil {
		return nil, fmt.Errorf("producto con id %d no encontrado", req.IDProducto)
	}

	if product.Stock < req.Cantidad {
		return nil, fmt.Errorf("stock insuficiente para el producto %s. Stock disponible: %d, solicitado: %d", 
			product.Nombre, product.Stock, req.Cantidad)
	}

	item := &domain.OrderItem{
		IDCompra:       orderID,
		IDProducto:     req.IDProducto,
		Cantidad:       req.Cantidad,
		PrecioUnitario: req.PrecioUnitario,
	}

	if err := s.repo.CreateOrderItem(item); err != nil {
		return nil, fmt.Errorf("error adding item to order: %w", err)
	}

	if err := s.catalogRepo.UpdateStock(req.IDProducto, -req.Cantidad); err != nil {
		return nil, fmt.Errorf("error updating stock: %w", err)
	}

	return item, nil
}

func (s *Service) UpdateOrderItem(orderID, itemID uint, req UpdateOrderItemRequest) (*domain.OrderItem, error) {
	item, err := s.repo.FindOrderItemByID(itemID)
	if err != nil {
		return nil, err
	}

	if item.IDCompra != orderID {
		return nil, errors.New("el item no pertenece a esta orden")
	}

	if req.Cantidad > 0 {
		product, err := s.catalogRepo.FindByID(item.IDProducto)
		if err != nil {
			return nil, fmt.Errorf("producto no encontrado: %w", err)
		}

		stockDifference := req.Cantidad - item.Cantidad
		if stockDifference > 0 {
			if product.Stock < stockDifference {
				return nil, fmt.Errorf("stock insuficiente. Stock disponible: %d, necesario: %d", 
					product.Stock, stockDifference)
			}
			if err := s.catalogRepo.UpdateStock(item.IDProducto, -stockDifference); err != nil {
				return nil, fmt.Errorf("error updating stock: %w", err)
			}
		} else if stockDifference < 0 {
			if err := s.catalogRepo.UpdateStock(item.IDProducto, -stockDifference); err != nil {
				return nil, fmt.Errorf("error updating stock: %w", err)
			}
		}

		item.Cantidad = req.Cantidad
	}

	if req.PrecioUnitario > 0 {
		item.PrecioUnitario = req.PrecioUnitario
	}

	if err := s.repo.UpdateOrderItem(item); err != nil {
		return nil, fmt.Errorf("error updating order item: %w", err)
	}

	return item, nil
}

func (s *Service) DeleteOrderItem(orderID, itemID uint) error {
	item, err := s.repo.FindOrderItemByID(itemID)
	if err != nil {
		return err
	}

	if item.IDCompra != orderID {
		return errors.New("el item no pertenece a esta orden")
	}

	if err := s.catalogRepo.UpdateStock(item.IDProducto, item.Cantidad); err != nil {
		return fmt.Errorf("error restoring stock: %w", err)
	}

	return s.repo.DeleteOrderItem(itemID)
}

func (s *Service) ToOrderItemResponse(item *domain.OrderItem) OrderItemResponse {
	return OrderItemResponse{
		IDDetalle:      item.ID,
		IDProducto:     item.IDProducto,
		Cantidad:       item.Cantidad,
		PrecioUnitario: item.PrecioUnitario,
		Subtotal:       float64(item.Cantidad) * item.PrecioUnitario,
	}
}

func (s *Service) ToOrderResponse(order *domain.Order) OrderResponse {
	items := make([]OrderItemResponse, len(order.Items))
	var total float64

	for i, item := range order.Items {
		items[i] = s.ToOrderItemResponse(&item)
		total += items[i].Subtotal
	}

	return OrderResponse{
		IDCompra:    order.ID,
		IDCliente:   order.IDCliente,
		IDVendedor:  order.IDVendedor,
		FechaCompra: order.FechaCompra,
		Items:       items,
		Total:       total,
	}
}
