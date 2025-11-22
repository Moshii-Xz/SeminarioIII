package users_test

import (
	"testing"

	"github.com/mordmora/expirapp/internal/domain"
	"github.com/mordmora/expirapp/internal/modules/users"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	err = db.AutoMigrate(&domain.User{}, &domain.Role{}, &domain.Client{}, &domain.Seller{})
	if err != nil {
		return nil, err
	}

	// Seed roles
	db.Create(&domain.Role{RoleName: "comprador"})
	db.Create(&domain.Role{RoleName: "vendedor"})

	return db, nil
}

func setupService(db *gorm.DB) *users.Service {
	repo := users.NewRepository(db)
	return users.NewService(repo)
}

func TestCreateUser_Success(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)

	req := users.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	user, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Name, user.Name)
	assert.Equal(t, req.Email, user.Email)
	assert.NotEqual(t, req.Password, user.Password) // Password should be hashed
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)

	req := users.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	_, _ = service.Create(req)
	user2, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, user2)
	assert.Contains(t, err.Error(), "email already in use")
}

func TestGetByID_Success(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)

	req := users.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	createdUser, _ := service.Create(req)

	foundUser, err := service.GetById(createdUser.ID)

	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID, foundUser.ID)
	assert.Equal(t, createdUser.Email, foundUser.Email)
}

func TestGetByID_NotFound(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)

	foundUser, err := service.GetById(999)

	assert.Error(t, err)
	assert.Nil(t, foundUser)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUpdateUser_Success(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)

	req := users.CreateUserRequest{
		Name:     "Original Name",
		Email:    "original@example.com",
		Password: "password123",
	}
	createdUser, _ := service.Create(req)

	updateReq := users.UpdateUserRequest{
		Name: "Updated Name",
	}

	updatedUser, err := service.Update(createdUser.ID, updateReq)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedUser.Name)
	assert.Equal(t, "original@example.com", updatedUser.Email) // Email shouldn't change
}

func TestUpdateUser_EmailConflict(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)

	// Create two users
	service.Create(users.CreateUserRequest{Name: "User 1", Email: "user1@example.com", Password: "pw1"})
	user2, _ := service.Create(users.CreateUserRequest{Name: "User 2", Email: "user2@example.com", Password: "pw2"})

	// Try to update user2's email to user1's email
	updateReq := users.UpdateUserRequest{
		Email: "user1@example.com",
	}

	_, err := service.Update(user2.ID, updateReq)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email already in use")
}

func TestDeleteUser_Success(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)

	user, _ := service.Create(users.CreateUserRequest{Name: "To Delete", Email: "delete@example.com", Password: "pw"})

	err := service.Delete(user.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = service.GetById(user.ID)
	assert.Error(t, err)
}

func TestListUsers_Pagination(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)

	// Create 15 users
	for i := 1; i <= 15; i++ {
		service.Create(users.CreateUserRequest{
			Name:     "User",
			Email:    "user" + string(rune(i)) + "@example.com", // Simple unique email generation
			Password: "pw",
		})
	}

	// Page 1, Limit 10
	list1, total, err := service.List(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(list1))
	assert.Equal(t, int64(15), total)

	// Page 2, Limit 10
	list2, _, err := service.List(2, 10)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(list2))
}

func TestLogin_Success(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)

	password := "securepassword"
	service.Create(users.CreateUserRequest{
		Name:     "Login User",
		Email:    "login@example.com",
		Password: password,
	})

	loginReq := users.LoginRequest{
		Email:    "login@example.com",
		Password: password,
	}

	response, err := service.Login(loginReq)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "login@example.com", response.User.Email)
}

func TestLogin_InvalidPassword(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)

	service.Create(users.CreateUserRequest{
		Name:     "Login User",
		Email:    "login@example.com",
		Password: "correctpassword",
	})

	loginReq := users.LoginRequest{
		Email:    "login@example.com",
		Password: "wrongpassword",
	}

	response, err := service.Login(loginReq)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestLogin_UserNotFound(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)

	loginReq := users.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	response, err := service.Login(loginReq)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestCreateUser_AsSeller(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)

	req := users.CreateUserRequest{
		Name:     "Seller User",
		Email:    "seller@example.com",
		Password: "password123",
		Role:     "vendedor",
	}

	user, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Verify Seller profile was created
	var seller domain.Seller
	err = db.First(&seller, user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, user.ID, seller.ID)
	assert.Equal(t, "General", seller.ResponsibleArea)
}

func TestCreateUser_AsBuyer(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)

	req := users.CreateUserRequest{
		Name:     "Buyer User",
		Email:    "buyer@example.com",
		Password: "password123",
		Role:     "comprador",
	}

	user, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Verify Client profile was created
	var client domain.Client
	err = db.First(&client, user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, user.ID, client.ID)
}
