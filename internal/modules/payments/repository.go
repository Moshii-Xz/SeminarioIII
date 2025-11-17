package payments

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

func (r *Repository) Create(payment *domain.Payment) error {
	return r.db.Create(payment).Error
}

func (r *Repository) FindByID(id uint) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.Preload("Order").Preload("PaymentMethod").First(&payment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	return &payment, nil
}

func (r *Repository) Update(payment *domain.Payment) error {
	return r.db.Save(payment).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&domain.Payment{}, id).Error
}

func (r *Repository) List(limit, offset int) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var total int64

	if err := r.db.Model(&domain.Payment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Order").Preload("PaymentMethod").Limit(limit).Offset(offset).Order("fecha_pago DESC").Find(&payments).Error
	return payments, total, err
}

func (r *Repository) FindByOrderID(orderID uint) ([]domain.Payment, error) {
	var payments []domain.Payment
	err := r.db.Preload("PaymentMethod").Where("id_compra = ?", orderID).Order("fecha_pago DESC").Find(&payments).Error
	return payments, err
}

func (r *Repository) GetTotalPaidByOrderID(orderID uint) (float64, error) {
	var total float64
	err := r.db.Model(&domain.Payment{}).
		Where("id_compra = ?", orderID).
		Select("COALESCE(SUM(monto), 0)").
		Scan(&total).Error
	return total, err
}

func (r *Repository) CountByOrderID(orderID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Payment{}).Where("id_compra = ?", orderID).Count(&count).Error
	return count, err
}

// PaymentMethod methods
func (r *Repository) CreatePaymentMethod(method *domain.PaymentMethod) error {
	return r.db.Create(method).Error
}

func (r *Repository) FindPaymentMethodByID(id uint) (*domain.PaymentMethod, error) {
	var method domain.PaymentMethod
	err := r.db.First(&method, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment method not found")
		}
		return nil, err
	}

	return &method, nil
}

func (r *Repository) FindPaymentMethodByName(name string) (*domain.PaymentMethod, error) {
	var method domain.PaymentMethod
	err := r.db.Where("nombre = ?", name).First(&method).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment method not found")
		}
		return nil, err
	}

	return &method, nil
}

func (r *Repository) ListPaymentMethods() ([]domain.PaymentMethod, error) {
	var methods []domain.PaymentMethod
	err := r.db.Order("nombre ASC").Find(&methods).Error
	return methods, err
}

func (r *Repository) UpdatePaymentMethod(method *domain.PaymentMethod) error {
	return r.db.Save(method).Error
}

func (r *Repository) DeletePaymentMethod(id uint) error {
	return r.db.Delete(&domain.PaymentMethod{}, id).Error
}
