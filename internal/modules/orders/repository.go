package orders

import (
	"errors"

	"github.com/mordmora/expirapp/internal/domain"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *Repository) CreateWithItems(order *domain.Order, items []domain.OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].IDCompra = order.ID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *Repository) FindByID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Items").Preload("Items.Product").First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return &order, nil
}

func (r *Repository) FindByIDWithItems(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Items").Preload("Items.Product").First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return &order, nil
}

func (r *Repository) Update(order *domain.Order) error {
	return r.db.Save(order).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&domain.Order{}, id).Error
}

func (r *Repository) List(limit, offset int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	if err := r.db.Model(&domain.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Items").Preload("Items.Product").Limit(limit).Offset(offset).Order("fecha_compra DESC").Find(&orders).Error
	return orders, total, err
}

func (r *Repository) FindByClientID(clientID uint, limit, offset int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	if err := r.db.Model(&domain.Order{}).Where("id_cliente = ?", clientID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Items").Preload("Items.Product").Where("id_cliente = ?", clientID).Limit(limit).Offset(offset).Order("fecha_compra DESC").Find(&orders).Error
	return orders, total, err
}

func (r *Repository) FindBySellerID(sellerID uint, limit, offset int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	if err := r.db.Model(&domain.Order{}).Where("id_vendedor = ?", sellerID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Items").Preload("Items.Product").Where("id_vendedor = ?", sellerID).Limit(limit).Offset(offset).Order("fecha_compra DESC").Find(&orders).Error
	return orders, total, err
}

func (r *Repository) CreateOrderItem(item *domain.OrderItem) error {
	return r.db.Create(item).Error
}

func (r *Repository) FindOrderItemByID(id uint) (*domain.OrderItem, error) {
	var item domain.OrderItem
	err := r.db.Preload("Product").First(&item, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order item not found")
		}
		return nil, err
	}

	return &item, nil
}

func (r *Repository) FindOrderItemsByOrderID(orderID uint) ([]domain.OrderItem, error) {
	var items []domain.OrderItem
	err := r.db.Preload("Product").Where("id_compra = ?", orderID).Find(&items).Error
	return items, err
}

func (r *Repository) UpdateOrderItem(item *domain.OrderItem) error {
	return r.db.Save(item).Error
}

func (r *Repository) DeleteOrderItem(id uint) error {
	return r.db.Delete(&domain.OrderItem{}, id).Error
}

func (r *Repository) DeleteOrderItemsByOrderID(orderID uint) error {
	return r.db.Where("id_compra = ?", orderID).Delete(&domain.OrderItem{}).Error
}
