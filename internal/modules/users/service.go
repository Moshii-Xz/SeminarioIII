package users

import (
	"errors"
	"fmt"

	"github.com/mordmora/expirapp/internal/domain"
	"github.com/mordmora/expirapp/internal/platform/auth"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) hashPassword(pass string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *Service) Create(req CreateUserRequest) (*domain.User, error) {
	exists, err := s.repo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already in use")
	}

	hashedPass, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing passwrord %w", err)
	}

	// Determine role (default to "comprador")
	roleName := "comprador"
	if req.Role == "tienda" {
		roleName = "tienda"
	} else if req.Role == "admin" {
		roleName = "admin"
	}
	// Note: "admin" role should typically be assigned by existing admins, not via public registration

	// Find role - this is required for creating the user-role relationship
	role, err := s.repo.FindRoleByName(roleName)
	if err != nil {
		return nil, fmt.Errorf("role '%s' not found: %w", roleName, err)
	}
	var roles []domain.Role
	roles = append(roles, *role)

	// Use transaction to ensure atomicity
	db := s.repo.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := &domain.User{
		ID:       req.ID,
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPass,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// Associate roles after user is created to ensure ID is available
	if err := tx.Model(user).Association("Roles").Append(roles); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error associating roles: %w", err)
	}

	// Create specific profile based on role
	switch roleName {
	case "tienda":
		store := &domain.Store{
			ID:              user.ID,
			ResponsibleArea: "General", // Default value
			Address:         "",         // Empty by default
			Phone:           "",         // Empty by default
		}
		if err := tx.Create(store).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error creating store profile: %w", err)
		}
	case "admin":
		admin := &domain.Admin{
			ID:                user.ID,
			SpecialPermissions: "", // Empty by default
		}
		if err := tx.Create(admin).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error creating admin profile: %w", err)
		}
	default: // "comprador" or any other role defaults to client
		client := &domain.Client{
			ID:      user.ID,
			Address: "", // Empty by default
			Phone:   "", // Empty by default
		}
		if err := tx.Create(client).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error creating client profile: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return user, nil
}

func (s *Service) GetById(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *Service) GetByEmail(email string) (*domain.User, error) {
	return s.repo.FindByEmail(email)
}

func (s *Service) Update(id uint, req UpdateUserRequest) (*domain.User, error) {
	usr, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		usr.Name = req.Name
	}

	if req.Email != "" && req.Email != usr.Email {
		exists, err := s.repo.ExistsByEmail(req.Email)
		if err != nil {
			return nil, fmt.Errorf("error checking email existence: %w", err)
		}
		if exists {
			return nil, errors.New("email already in use")
		}
	}

	if err := s.repo.Update(usr); err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return usr, nil
}

func (s *Service) ChangePassword(id uint, req ChangePasswordRequest) error {

	usr, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if !s.verifyPassword(usr.Password, req.CurrentPass) {
		return errors.New("current password is incorrect")
	}

	hashedPass, err := s.hashPassword(req.NewPass)
	if err != nil {
		return fmt.Errorf("error to hash password: %w", err)
	}

	usr.Password = hashedPass
	return s.repo.Update(usr)
}

func (s *Service) Delete(id uint) error {

	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(int(id))
}

func (s *Service) List(page, limit int) ([]domain.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.repo.List(limit, offset)
}

func (s *Service) verifyPassword(hashedP, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedP), []byte(password))
	return err == nil
}

func (s *Service) ToResponse(usr *domain.User) UserResponse {
	return UserResponse{
		ID:        usr.ID,
		Name:      usr.Name,
		Email:     usr.Email,
		CreatedAt: usr.CreatedAt,
	}
}

func (s *Service) Login(req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !s.verifyPassword(user.Password, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	role := ""
	if len(user.Roles) > 0 {
		role = user.Roles[0].RoleName
	}

	token, err := auth.GenerateToken(user.ID, role)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  s.ToResponse(user),
	}, nil
}
