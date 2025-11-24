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

	// Validar precios según la etiqueta
	if req.Badge != nil {
		if *req.Badge == "Oferta" {
			// Productos en oferta deben tener precio_original y precio_descuento
			if req.OriginalPrice == nil || req.DiscountPrice == nil {
				return nil, errors.New("productos con etiqueta 'Oferta' deben tener precio_original y precio_descuento")
			}
			if *req.OriginalPrice <= 0 || *req.DiscountPrice <= 0 {
				return nil, errors.New("precio_original y precio_descuento deben ser mayores a 0")
			}
			if *req.DiscountPrice >= *req.OriginalPrice {
				return nil, errors.New("precio_descuento debe ser menor que precio_original")
			}
		} else if *req.Badge == "Donación" {
			// Productos de donación no deben tener precio
			if req.Price != nil || req.OriginalPrice != nil || req.DiscountPrice != nil {
				return nil, errors.New("productos con etiqueta 'Donación' no deben tener precio")
			}
		}
	} else {
		// Productos sin etiqueta deben tener precio normal
		if req.Price == nil || *req.Price <= 0 {
			return nil, errors.New("productos sin etiqueta deben tener un precio válido mayor a 0")
		}
		// Asegurar que no tengan precios de oferta
		req.OriginalPrice = nil
		req.DiscountPrice = nil
	}

	product := &domain.Product{
		Name:           req.Name,
		Description:    req.Description,
		ImageURL:       req.ImageURL,
		Price:          req.Price,
		OriginalPrice:  req.OriginalPrice,
		DiscountPrice:  req.DiscountPrice,
		ExpirationDate: req.ExpirationDate,
		Stock:          req.Stock,
		Badge:          req.Badge,
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
	
	// Determinar la etiqueta final (la nueva o la existente)
	finalBadge := product.Badge
	if req.Badge != nil {
		finalBadge = req.Badge
	}
	
	// Validar y actualizar precios según la etiqueta
	if finalBadge != nil {
		if *finalBadge == "Oferta" {
			// Productos en oferta deben tener precio_original y precio_descuento
			if req.OriginalPrice != nil {
				product.OriginalPrice = req.OriginalPrice
			}
			if req.DiscountPrice != nil {
				product.DiscountPrice = req.DiscountPrice
			}
			// Validar que ambos precios estén presentes y sean válidos
			if product.OriginalPrice == nil || product.DiscountPrice == nil {
				return nil, errors.New("productos con etiqueta 'Oferta' deben tener precio_original y precio_descuento")
			}
			if *product.OriginalPrice <= 0 || *product.DiscountPrice <= 0 {
				return nil, errors.New("precio_original y precio_descuento deben ser mayores a 0")
			}
			if *product.DiscountPrice >= *product.OriginalPrice {
				return nil, errors.New("precio_descuento debe ser menor que precio_original")
			}
			// Limpiar precio normal
			product.Price = nil
		} else if *finalBadge == "Donación" {
			// Productos de donación no deben tener precio
			product.Price = nil
			product.OriginalPrice = nil
			product.DiscountPrice = nil
		}
	} else {
		// Productos sin etiqueta deben tener precio normal
		if req.Price != nil {
			product.Price = req.Price
		}
		// Limpiar precios de oferta
		product.OriginalPrice = nil
		product.DiscountPrice = nil
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
	// Actualizar etiqueta (badge) si se proporciona
	if req.Badge != nil {
		product.Badge = req.Badge
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
		OriginalPrice:  product.OriginalPrice,
		DiscountPrice:  product.DiscountPrice,
		ExpirationDate: product.ExpirationDate,
		Stock:          product.Stock,
		Badge:          product.Badge,
		CategoryID:     product.CategoryID,
	}
	
	if product.Category != nil {
		response.CategoryName = &product.Category.Name
	}
	
	return response
}
