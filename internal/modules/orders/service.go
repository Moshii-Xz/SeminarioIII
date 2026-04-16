package orders

import (
	"errors"
	"fmt"
	"time"

	"github.com/Moshii-Xz/SeminarioIII/internal/domain"
	catalogRepo "github.com/Moshii-Xz/SeminarioIII/internal/modules/catalog"
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
	orderDetails := make([]domain.OrderDetail, len(req.Items))

	for i, itemReq := range req.Items {
		product, err := s.catalogRepo.FindByID(itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product with id %d not found", itemReq.ProductID)
		}

		if product.Stock < itemReq.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s (id: %d). Available: %d, Requested: %d",
				product.Name, itemReq.ProductID, product.Stock, itemReq.Quantity)
		}

		orderDetails[i] = domain.OrderDetail{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			UnitPrice: itemReq.UnitPrice,
		}
	}

	order := &domain.Order{
		ClientID:      req.ClientID,
		StoreID:       req.StoreID,
		Date:          time.Now(),
		PaymentMethod: req.PaymentMethod,
		Details:       orderDetails,
	}

	if err := s.repo.CreateWithItems(order, orderDetails); err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	for _, item := range orderDetails {
		if err := s.catalogRepo.UpdateStock(item.ProductID, -item.Quantity); err != nil {
			return nil, fmt.Errorf("error updating stock for product %d: %w", item.ProductID, err)
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

	if req.StoreID != nil {
		order.StoreID = req.StoreID
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

	for _, item := range order.Details {
		if err := s.catalogRepo.UpdateStock(item.ProductID, item.Quantity); err != nil {
			fmt.Printf("warning: error restoring stock for product %d: %v\n", item.ProductID, err)
		}
	}

	return s.repo.Delete(id)
}

func (s *Service) List(page, limit int) ([]domain.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.List(limit, offset)
}

func (s *Service) ListByClient(clientID uint, page, limit int) ([]domain.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.FindByClientID(clientID, limit, offset)
}

func (s *Service) ListByStore(storeID uint, page, limit int) ([]domain.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.FindByStoreID(storeID, limit, offset)
}

func (s *Service) AddOrderItem(orderID uint, req AddOrderItemRequest) (*domain.OrderDetail, error) {
	_, err := s.repo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	product, err := s.catalogRepo.FindByID(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product with id %d not found", req.ProductID)
	}

	if product.Stock < req.Quantity {
		return nil, fmt.Errorf("insufficient stock for product %s. Available: %d, Requested: %d",
			product.Name, product.Stock, req.Quantity)
	}

	item := &domain.OrderDetail{
		OrderID:   orderID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		UnitPrice: req.UnitPrice,
	}

	if err := s.repo.CreateDetail(item); err != nil {
		return nil, fmt.Errorf("error adding item to order: %w", err)
	}

	if err := s.catalogRepo.UpdateStock(req.ProductID, -req.Quantity); err != nil {
		return nil, fmt.Errorf("error updating stock: %w", err)
	}

	// Verify order exists to ensure consistency
	_, _ = s.repo.FindByID(orderID)

	return item, nil
}

func (s *Service) UpdateOrderItem(orderID, itemID uint, req UpdateOrderItemRequest) (*domain.OrderDetail, error) {
	item, err := s.repo.FindDetailByID(itemID)
	if err != nil {
		return nil, err
	}

	if item.OrderID != orderID {
		return nil, errors.New("item does not belong to this order")
	}

	if req.Quantity > 0 {
		product, err := s.catalogRepo.FindByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product not found: %w", err)
		}

		stockDifference := req.Quantity - item.Quantity
		if stockDifference > 0 {
			if product.Stock < stockDifference {
				return nil, fmt.Errorf("insufficient stock. Available: %d, Needed: %d",
					product.Stock, stockDifference)
			}
			if err := s.catalogRepo.UpdateStock(item.ProductID, -stockDifference); err != nil {
				return nil, fmt.Errorf("error updating stock: %w", err)
			}
		} else if stockDifference < 0 {
			if err := s.catalogRepo.UpdateStock(item.ProductID, -stockDifference); err != nil {
				return nil, fmt.Errorf("error updating stock: %w", err)
			}
		}

		item.Quantity = req.Quantity
	}

	if req.UnitPrice > 0 {
		item.UnitPrice = req.UnitPrice
	}

	if err := s.repo.UpdateDetail(item); err != nil {
		return nil, fmt.Errorf("error updating order item: %w", err)
	}

	return item, nil
}

func (s *Service) DeleteOrderItem(orderID, itemID uint) error {
	item, err := s.repo.FindDetailByID(itemID)
	if err != nil {
		return err
	}

	if item.OrderID != orderID {
		return errors.New("item does not belong to this order")
	}

	if err := s.catalogRepo.UpdateStock(item.ProductID, item.Quantity); err != nil {
		return fmt.Errorf("error restoring stock: %w", err)
	}

	return s.repo.DeleteDetail(itemID)
}

func (s *Service) ToOrderItemResponse(item *domain.OrderDetail) OrderItemResponse {
	return OrderItemResponse{
		ID:        item.ID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		UnitPrice: item.UnitPrice,
		Subtotal:  float64(item.Quantity) * item.UnitPrice,
	}
}

func (s *Service) ToResponse(order *domain.Order) OrderResponse {
	items := make([]OrderItemResponse, len(order.Details))
	var total float64

	for i, item := range order.Details {
		items[i] = s.ToOrderItemResponse(&item)
		total += items[i].Subtotal
	}

	return OrderResponse{
		ID:            order.ID,
		ClientID:      order.ClientID,
		StoreID:       order.StoreID,
		Date:          order.Date,
		Status:        order.Status,
		PaymentMethod: order.PaymentMethod,
		Items:         items,
		Total:         total,
	}
}
