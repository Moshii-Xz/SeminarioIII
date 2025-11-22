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
	if req.Role == "vendedor" {
		roleName = "vendedor"
	}
	// Note: "admin" role cannot be assigned via public registration

	role, err := s.repo.FindRoleByName(roleName)
	var roles []domain.Role
	if err == nil {
		roles = append(roles, *role)
	} else {
		fmt.Printf("Warning: Role '%s' not found: %v\n", roleName, err)
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPass,
		Roles:    roles,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// Create specific profile based on role
	if roleName == "vendedor" {
		seller := &domain.Seller{
			ID:              user.ID,
			ResponsibleArea: "General", // Default value
		}
		if err := s.repo.CreateSeller(seller); err != nil {
			// Rollback user creation? Ideally yes, but for now let's log error
			fmt.Printf("Error creating seller profile: %v\n", err)
		}
	} else {
		// Default to client profile
		client := &domain.Client{
			ID:      user.ID,
			Address: "", // Empty by default
			Phone:   "", // Empty by default
		}
		if err := s.repo.CreateClient(client); err != nil {
			fmt.Printf("Error creating client profile: %v\n", err)
		}
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
