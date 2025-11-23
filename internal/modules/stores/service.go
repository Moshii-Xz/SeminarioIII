package stores

import (
	"errors"
	"fmt"

	"github.com/mordmora/expirapp/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	repo     *Repository
	userRepo UserRepository
}

// UserRepository defines the interface for user operations needed by stores service
type UserRepository interface {
	ExistsByEmail(email string) (bool, error)
	FindRoleByName(name string) (*domain.Role, error)
	GetDB() *gorm.DB
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// NewServiceWithUserRepo creates a service with user repository for creating stores
func NewServiceWithUserRepo(repo *Repository, userRepo UserRepository) *Service {
	return &Service{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *Service) GetById(id uint) (*StoreResponse, error) {
	store, user, err := s.repo.FindByIDWithUser(id)
	if err != nil {
		return nil, err
	}

	return s.ToResponse(store, user), nil
}

func (s *Service) Update(id uint, req UpdateStoreRequest) (*StoreResponse, error) {
	store, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.ResponsibleArea != "" {
		store.ResponsibleArea = req.ResponsibleArea
	}
	if req.Address != "" {
		store.Address = req.Address
	}
	if req.Phone != "" {
		store.Phone = req.Phone
	}

	if err := s.repo.Update(store); err != nil {
		return nil, fmt.Errorf("error updating store: %w", err)
	}

	// Get user info
	_, user, err := s.repo.FindByIDWithUser(id)
	if err != nil {
		return nil, err
	}

	return s.ToResponse(store, user), nil
}

func (s *Service) List(page, limit int) ([]StoreResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	stores, users, total, err := s.repo.List(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Create a map of user ID to user for quick lookup
	userMap := make(map[uint]*domain.User)
	for i := range users {
		userMap[users[i].ID] = &users[i]
	}

	// Build response
	responses := make([]StoreResponse, len(stores))
	for i, store := range stores {
		user := userMap[store.ID]
		if user != nil {
			responses[i] = *s.ToResponse(&store, user)
		}
	}

	return responses, total, nil
}

func (s *Service) GetProducts(storeID uint, page, limit int) (*StoreProductsResponse, error) {
	store, user, err := s.repo.FindByIDWithUser(storeID)
	if err != nil {
		return nil, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	products, total, err := s.repo.FindProductsByStoreID(storeID, limit, offset)
	if err != nil {
		return nil, err
	}

	productInfos := make([]ProductInfo, len(products))
	for i, p := range products {
		productInfos[i] = ProductInfo{
			ID:             p.ID,
			Name:           p.Name,
			Description:    p.Description,
			ImageURL:       p.ImageURL,
			Price:          p.Price,
			ExpirationDate: p.ExpirationDate,
			Stock:          p.Stock,
		}
	}

	return &StoreProductsResponse{
		Store:    *s.ToResponse(store, user),
		Products: productInfos,
		Total:    total,
	}, nil
}

func (s *Service) GetOrders(storeID uint, page, limit int) (*StoreOrdersResponse, error) {
	store, user, err := s.repo.FindByIDWithUser(storeID)
	if err != nil {
		return nil, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	orders, total, err := s.repo.FindOrdersByStoreID(storeID, limit, offset)
	if err != nil {
		return nil, err
	}

	orderInfos := make([]OrderInfo, len(orders))
	for i, o := range orders {
		var total float64
		for _, detail := range o.Details {
			total += detail.UnitPrice * float64(detail.Quantity)
		}

		orderInfos[i] = OrderInfo{
			ID:       o.ID,
			ClientID: o.ClientID,
			Date:     o.Date,
			Total:    total,
		}
	}

	return &StoreOrdersResponse{
		Store:  *s.ToResponse(store, user),
		Orders: orderInfos,
		Total:  total,
	}, nil
}

func (s *Service) Create(req CreateStoreRequest) (*StoreResponse, error) {
	if s.userRepo == nil {
		return nil, errors.New("user repository not configured")
	}

	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking email existence: %w", err)
	}
	if exists {
		return nil, errors.New("email already in use")
	}

	// Hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	role, err := s.userRepo.FindRoleByName("tienda")
	if err != nil {
		return nil, fmt.Errorf("role 'tienda' not found: %w", err)
	}

	db := s.userRepo.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPass),
		Roles:    []domain.Role{*role},
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	responsibleArea := req.ResponsibleArea
	if responsibleArea == "" {
		responsibleArea = "General"
	}

	store := &domain.Store{
		ID:              user.ID,
		ResponsibleArea: responsibleArea,
		Address:         req.Address,
		Phone:           req.Phone,
	}

	if err := tx.Create(store).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating store profile: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return s.ToResponse(store, user), nil
}

func (s *Service) ToResponse(store *domain.Store, user *domain.User) *StoreResponse {
	return &StoreResponse{
		ID:              store.ID,
		ResponsibleArea: store.ResponsibleArea,
		Address:         store.Address,
		Phone:           store.Phone,
		User: UserInfo{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
		CreatedAt: user.CreatedAt,
	}
}

