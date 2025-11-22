package users

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

func (r *Repository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) CreateClient(client *domain.Client) error {
	return r.db.Create(client).Error
}

func (r *Repository) CreateSeller(seller *domain.Seller) error {
	return r.db.Create(seller).Error
}

func (r *Repository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User

	err := r.db.Preload("Roles").Where("correo = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *Repository) Delete(id int) error {
	return r.db.Delete(&domain.User{}, id).Error
}

func (r *Repository) List(limit, offset int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	if err := r.db.Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *Repository) ExistsByEmail(email string) (bool, error) {
	var c int64
	err := r.db.Model(&domain.User{}).Where("correo = ?", email).Count(&c).Error
	return c > 0, err
}

func (r *Repository) FindRoleByName(name string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Where("nombre = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
