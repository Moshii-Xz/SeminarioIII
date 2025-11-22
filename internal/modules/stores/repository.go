package stores

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

func (r *Repository) FindByID(id uint) (*domain.Store, error) {
	var store domain.Store
	err := r.db.First(&store, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("store not found")
		}
		return nil, err
	}
	return &store, nil
}

func (r *Repository) FindByIDWithUser(id uint) (*domain.Store, *domain.User, error) {
	var store domain.Store
	err := r.db.First(&store, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("store not found")
		}
		return nil, nil, err
	}

	var user domain.User
	err = r.db.First(&user, id).Error
	if err != nil {
		return nil, nil, err
	}

	return &store, &user, nil
}

func (r *Repository) Update(store *domain.Store) error {
	return r.db.Save(store).Error
}

func (r *Repository) List(limit, offset int) ([]domain.Store, []domain.User, int64, error) {
	var stores []domain.Store
	var total int64

	if err := r.db.Model(&domain.Store{}).Count(&total).Error; err != nil {
		return nil, nil, 0, err
	}

	err := r.db.Limit(limit).Offset(offset).Find(&stores).Error
	if err != nil {
		return nil, nil, 0, err
	}

	// Get users for all stores
	var storeIDs []uint
	for _, s := range stores {
		storeIDs = append(storeIDs, s.ID)
	}

	var users []domain.User
	if len(storeIDs) > 0 {
		err = r.db.Where("id_usuario IN ?", storeIDs).Find(&users).Error
		if err != nil {
			return nil, nil, 0, err
		}
	}

	return stores, users, total, nil
}

func (r *Repository) FindProductsByStoreID(storeID uint, limit, offset int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	if err := r.db.Model(&domain.Product{}).Where("id_tienda = ?", storeID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("id_tienda = ?", storeID).Limit(limit).Offset(offset).Find(&products).Error
	return products, total, err
}

func (r *Repository) FindOrdersByStoreID(storeID uint, limit, offset int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	if err := r.db.Model(&domain.Order{}).Where("id_tienda = ?", storeID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Details").Preload("Details.Product").Where("id_tienda = ?", storeID).Limit(limit).Offset(offset).Order("fecha_compra DESC").Find(&orders).Error
	return orders, total, err
}

