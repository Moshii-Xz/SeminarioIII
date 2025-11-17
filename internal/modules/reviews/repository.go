package reviews

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

func (r *Repository) Create(review *domain.Review) error {
	return r.db.Create(review).Error
}

func (r *Repository) FindByID(id uint) (*domain.Review, error) {
	var review domain.Review
	err := r.db.Preload("Product").Preload("Client").First(&review, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("review not found")
		}
		return nil, err
	}

	return &review, nil
}

func (r *Repository) Update(review *domain.Review) error {
	return r.db.Save(review).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&domain.Review{}, id).Error
}

func (r *Repository) ExistsByClientAndProduct(clientID, productID uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Review{}).
		Where("id_cliente = ? AND id_producto = ?", clientID, productID).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) List(limit, offset int) ([]domain.Review, int64, error) {
	var reviews []domain.Review
	var total int64

	if err := r.db.Model(&domain.Review{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Product").Preload("Client").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, total, err
}

func (r *Repository) ListByProduct(productID uint, limit, offset int) ([]domain.Review, int64, error) {
	var reviews []domain.Review
	var total int64

	if err := r.db.Model(&domain.Review{}).
		Where("id_producto = ?", productID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Client").
		Where("id_producto = ?", productID).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&reviews).Error

	return reviews, total, err
}

type RatingSummary struct {
	Average float64
	Count   int64
}

func (r *Repository) GetProductRatingSummary(productID uint) (*RatingSummary, error) {
	var summary RatingSummary
	err := r.db.Model(&domain.Review{}).
		Where("id_producto = ?", productID).
		Select("COALESCE(AVG(calificacion), 0) AS average, COUNT(id_resena) AS count").
		Scan(&summary).Error
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

