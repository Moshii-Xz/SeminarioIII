package reviews

import (
	"errors"
	"fmt"

	"github.com/mordmora/expirapp/internal/domain"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req CreateReviewRequest) (*domain.Review, error) {
	exists, err := s.repo.ExistsByClientAndProduct(req.ClientID, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("error checking existing review: %w", err)
	}
	if exists {
		return nil, errors.New("el cliente ya registró una reseña para este producto")
	}

	review := &domain.Review{
		ProductID: req.ProductID,
		ClientID:  req.ClientID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	if err := s.repo.Create(review); err != nil {
		return nil, fmt.Errorf("error creating review: %w", err)
	}

	return review, nil
}

func (s *Service) GetByID(id uint) (*domain.Review, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Update(id uint, req UpdateReviewRequest) (*domain.Review, error) {
	review, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Rating > 0 {
		review.Rating = req.Rating
	}

	if req.Comment != "" {
		review.Comment = req.Comment
	}

	if err := s.repo.Update(review); err != nil {
		return nil, fmt.Errorf("error updating review: %w", err)
	}

	return review, nil
}

func (s *Service) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *Service) List(page, limit int) ([]domain.Review, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.repo.List(limit, offset)
}

func (s *Service) ListByProduct(productID uint, page, limit int) ([]domain.Review, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.repo.ListByProduct(productID, limit, offset)
}

func (s *Service) GetProductRatingSummary(productID uint) (*ProductRatingSummary, error) {
	summary, err := s.repo.GetProductRatingSummary(productID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo resumen de calificaciones: %w", err)
	}

	return &ProductRatingSummary{
		ProductID:     productID,
		AverageRating: summary.Average,
		ReviewsCount:  summary.Count,
	}, nil
}

func (s *Service) ToResponse(review *domain.Review) ReviewResponse {
	return ReviewResponse{
		ID:        review.ID,
		ProductID: review.ProductID,
		ClientID:  review.ClientID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}
}

import (
	"errors"
	"fmt"

	"github.com/mordmora/expirapp/internal/domain"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req CreateReviewRequest) (*domain.Review, error) {
	exists, err := s.repo.ExistsByClientAndProduct(req.ClientID, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("error checking existing review: %w", err)
	}
	if exists {
		return nil, errors.New("el cliente ya registró una reseña para este producto")
	}

	review := &domain.Review{
		ProductID: req.ProductID,
		ClientID:  req.ClientID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	if err := s.repo.Create(review); err != nil {
		return nil, fmt.Errorf("error creating review: %w", err)
	}

	return review, nil
}

func (s *Service) GetByID(id uint) (*domain.Review, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Update(id uint, req UpdateReviewRequest) (*domain.Review, error) {
	review, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Rating > 0 {
		review.Rating = req.Rating
	}

	if req.Comment != "" {
		review.Comment = req.Comment
	}

	if err := s.repo.Update(review); err != nil {
		return nil, fmt.Errorf("error updating review: %w", err)
	}

	return review, nil
}

func (s *Service) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *Service) List(page, limit int) ([]domain.Review, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.repo.List(limit, offset)
}

func (s *Service) ListByProduct(productID uint, page, limit int) ([]domain.Review, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.repo.ListByProduct(productID, limit, offset)
}

func (s *Service) GetProductRatingSummary(productID uint) (*ProductRatingSummary, error) {
	summary, err := s.repo.GetProductRatingSummary(productID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo resumen de calificaciones: %w", err)
	}

	return &ProductRatingSummary{
		ProductID:      productID,
		AverageRating:  summary.Average,
		ReviewsCount:   summary.Count,
	}, nil
}

func (s *Service) ToResponse(review *domain.Review) ReviewResponse {
	return ReviewResponse{
		ID:        review.ID,
		ProductID: review.ProductID,
		ClientID:  review.ClientID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}
}

