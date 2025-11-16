package catalog

import (
	"errors"
	"fmt"
	"time"

	"github.com/mordmora/expirapp/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req CreateProductRequest) (*domain.Product, error) {
	if req.FechaVencimiento.Before(time.Now()) {
		return nil, errors.New("la fecha de vencimiento no puede ser en el pasado")
	}

	product := &domain.Product{
		Nombre:           req.Nombre,
		Descripcion:      req.Descripcion,
		Precio:           req.Precio,
		FechaVencimiento: req.FechaVencimiento,
		Stock:            req.Stock,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	return product, nil
}

func (s *Service) GetById(id uint) (*domain.Product, error) {
	return s.repo.FindByID(id)
}

func (s *Service) GetByName(name string) (*domain.Product, error) {
	return s.repo.FindByName(name)
}

func (s *Service) Update(id uint, req UpdateProductRequest) (*domain.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Nombre != "" {
		product.Nombre = req.Nombre
	}

	if req.Descripcion != "" {
		product.Descripcion = req.Descripcion
	}

	if req.Precio > 0 {
		product.Precio = req.Precio
	}

	if !req.FechaVencimiento.IsZero() {
		if req.FechaVencimiento.Before(time.Now()) {
			return nil, errors.New("la fecha de vencimiento no puede ser en el pasado")
		}
		product.FechaVencimiento = req.FechaVencimiento
	}

	if req.Stock >= 0 {
		product.Stock = req.Stock
	}

	if err := s.repo.Update(product); err != nil {
		return nil, fmt.Errorf("error updating product: %w", err)
	}

	return product, nil
}

func (s *Service) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *Service) List(page, limit int) ([]domain.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.repo.List(limit, offset)
}

func (s *Service) GetByExpirationDate(date time.Time) ([]domain.Product, error) {
	return s.repo.FindByExpirationDate(date)
}

func (s *Service) GetExpiringSoon(days int) ([]domain.Product, error) {
	if days < 1 {
		days = 7 
	}
	return s.repo.FindExpiringSoon(days)
}

func (s *Service) UpdateStock(id uint, quantity int) error {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	newStock := product.Stock + quantity
	if newStock < 0 {
		return errors.New("stock insuficiente: no se puede reducir el stock por debajo de 0")
	}

	return s.repo.UpdateStock(id, quantity)
}

func (s *Service) ToResponse(product *domain.Product) ProductResponse {
	return ProductResponse{
		ID:               product.ID,
		Nombre:           product.Nombre,
		Descripcion:      product.Descripcion,
		Precio:           product.Precio,
		FechaVencimiento: product.FechaVencimiento,
		Stock:            product.Stock,
	}
}
