package catalog

import (
	"errors"
	"time"

	"github.com/mordmora/expirapp/internal/domain"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *Repository) FindByID(id uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

func (r *Repository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

func (r *Repository) List(limit, offset int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	if err := r.db.Model(&domain.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Limit(limit).Offset(offset).Find(&products).Error
	return products, total, err
}

func (r *Repository) FindByName(name string) (*domain.Product, error) {
	var product domain.Product

	err := r.db.Where("nombre = ?", name).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

func (r *Repository) FindByExpirationDate(date time.Time) ([]domain.Product, error) {
	var products []domain.Product

	err := r.db.Where("fecha_vencimiento = ?", date).Find(&products).Error
	return products, err
}

func (r *Repository) FindExpiringSoon(days int) ([]domain.Product, error) {
	var products []domain.Product
	threshold := time.Now().AddDate(0, 0, days)

	err := r.db.Where("fecha_vencimiento <= ? AND fecha_vencimiento >= ?", threshold, time.Now()).Find(&products).Error
	return products, err
}

func (r *Repository) UpdateStock(id uint, quantity int) error {
	return r.db.Model(&domain.Product{}).Where("id_producto = ?", id).Update("stock", gorm.Expr("stock + ?", quantity)).Error
}
