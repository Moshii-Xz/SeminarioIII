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

func (s *Service) Create(req CreateProductRequest, storeID uint) (*domain.Product, error) {
	if req.ExpirationDate.Before(time.Now()) {
		return nil, errors.New("la fecha de vencimiento no puede ser en el pasado")
	}

	status := req.Status
	if status == "" {
		status = "En preparación" // Default value
	} else {
		validStatuses := []string{"En preparación", "Listo para recoger", "Entregado"}
		valid := false
		for _, vs := range validStatuses {
			if status == vs {
				valid = true
				break
			}
		}
		if !valid {
			return nil, errors.New("estado inválido. Debe ser uno de: 'En preparación', 'Listo para recoger', 'Entregado'")
		}
	}

	product := &domain.Product{
		Name:           req.Name,
		Description:    req.Description,
		ImageURL:       req.ImageURL,
		Price:          req.Price,
		ExpirationDate: req.ExpirationDate,
		Stock:          req.Stock,
		Status:        status,
		StoreID:       storeID, // Usar el storeID del token JWT
		CategoryID:    req.CategoryID,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	return product, nil
}

func (s *Service) GetById(id uint) (*domain.Product, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Update(id uint, req UpdateProductRequest, storeID uint) (*domain.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Validar que el producto pertenezca a la tienda del usuario autenticado
	if product.StoreID != storeID {
		return nil, errors.New("no tienes permiso para modificar este producto")
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
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
	if req.Status != "" {
		validStatuses := []string{"En preparación", "Listo para recoger", "Entregado"}
		valid := false
		for _, vs := range validStatuses {
			if req.Status == vs {
				valid = true
				break
			}
		}
		if !valid {
			return nil, errors.New("estado inválido. Debe ser uno de: 'En preparación', 'Listo para recoger', 'Entregado'")
		}
		product.Status = req.Status
	}
	// Manejar actualización de categoría (puede ser nil para eliminar la categoría)
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	} else {
		// Si se envía explícitamente null, eliminar la categoría
		// Pero como es omitempty, solo se actualiza si viene en el request
		// Por ahora, solo actualizamos si viene un valor
	}

	if err := s.repo.Update(product); err != nil {
		return nil, fmt.Errorf("error updating product: %w", err)
	}

	// Recargar el producto con la categoría actualizada
	updatedProduct, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error reloading product: %w", err)
	}

	return updatedProduct, nil
}

func (s *Service) Delete(id uint, storeID uint) error {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Validar que el producto pertenezca a la tienda del usuario autenticado
	if product.StoreID != storeID {
		return errors.New("no tienes permiso para eliminar este producto")
	}

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
	response := ProductResponse{
		ID:             product.ID,
		Name:           product.Name,
		Description:    product.Description,
		ImageURL:       product.ImageURL,
		Price:          product.Price,
		ExpirationDate: product.ExpirationDate,
		Stock:          product.Stock,
		Status:         product.Status,
		CategoryID:     product.CategoryID,
	}
	
	if product.Category != nil {
		response.CategoryName = &product.Category.Name
	}
	
	return response
}
