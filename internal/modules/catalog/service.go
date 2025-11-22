package catalog

import (
	"errors"
	"fmt"
	"time"

	"github.com/mordmora/expirapp/internal/domain"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req CreateProductRequest) (*domain.Product, error) {
	if req.ExpirationDate.Before(time.Now()) {
		return nil, errors.New("la fecha de vencimiento no puede ser en el pasado")
	}

	product := &domain.Product{
		Name:           req.Name,
		Description:    req.Description,
		Price:          req.Price,
		ExpirationDate: req.ExpirationDate,
		Stock:          req.Stock,
		StoreID:       req.StoreID,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	return product, nil
}

func (s *Service) GetById(id uint) (*domain.Product, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Update(id uint, req UpdateProductRequest) (*domain.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if !req.ExpirationDate.IsZero() {
		if req.ExpirationDate.Before(time.Now()) {
			return nil, errors.New("la fecha de vencimiento no puede ser en el pasado")
		}
		product.ExpirationDate = req.ExpirationDate
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
	return s.repo.Delete(id)
}

func (s *Service) List(page, limit int) ([]domain.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.List(limit, offset)
}

func (s *Service) ToResponse(product *domain.Product) ProductResponse {
	return ProductResponse{
		ID:             product.ID,
		Name:           product.Name,
		Description:    product.Description,
		Price:          product.Price,
		ExpirationDate: product.ExpirationDate,
		Stock:          product.Stock,
	}
}
