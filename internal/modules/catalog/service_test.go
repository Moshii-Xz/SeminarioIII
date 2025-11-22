package catalog_test

import (
	"testing"
	"time"

	"github.com/mordmora/expirapp/internal/domain"
	"github.com/mordmora/expirapp/internal/modules/catalog"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Enable Foreign Keys for SQLite
	db.Exec("PRAGMA foreign_keys = ON")

	// Migrate the schema
	err = db.AutoMigrate(&domain.User{}, &domain.Seller{}, &domain.Product{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupService(db *gorm.DB) *catalog.Service {
	repo := catalog.NewRepository(db)
	return catalog.NewService(repo)
}

func createSeller(db *gorm.DB) uint {
	user := domain.User{
		Name:     "Seller",
		Email:    "seller@test.com",
		Password: "hash",
	}
	db.Create(&user)

	seller := domain.Seller{
		ID:              user.ID,
		ResponsibleArea: "Test Area",
	}
	db.Create(&seller)

	return seller.ID
}

func TestCreateProduct_Success(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)
	sellerID := createSeller(db)

	req := catalog.CreateProductRequest{
		Name:           "Valid Product",
		Price:          10.0,
		ExpirationDate: time.Now().Add(24 * time.Hour),
		SellerID:       sellerID,
	}

	product, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, req.Name, product.Name)
	assert.Equal(t, sellerID, product.SellerID)
}

func TestCreateProduct_PastExpirationDate(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)
	sellerID := createSeller(db)

	req := catalog.CreateProductRequest{
		Name:           "Expired Product",
		Price:          10.0,
		ExpirationDate: time.Now().Add(-24 * time.Hour), // Past date
		SellerID:       sellerID,
	}

	product, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "la fecha de vencimiento no puede ser en el pasado")
}

func TestCreateProduct_InvalidSeller(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)
	// No seller created

	req := catalog.CreateProductRequest{
		Name:           "Orphan Product",
		Price:          10.0,
		ExpirationDate: time.Now().Add(24 * time.Hour),
		SellerID:       999, // Non-existent seller
	}

	product, err := service.Create(req)

	// GORM might not return an error immediately if foreign key constraints aren't enforced by SQLite driver by default,
	// or if the repository doesn't check existence.
	// However, standard SQL databases would fail.
	// Let's see if the service or repo handles this.
	// If not, this test might fail (expecting error but getting success).
	// In a real DB, this would fail. In SQLite memory, we need to ensure FKs are on.

	// For now, let's assume the DB enforces it or we want to catch it.
	// If this fails, it means we need to enable FKs in SQLite or add a check in service.

	// Actually, let's check if the service validates seller existence. It probably doesn't explicitly.
	// But the DB insert should fail.

	if err == nil {
		// If no error, check if product was actually saved with that ID.
		// SQLite default might be lax.
		// Let's skip asserting error for now and see if it runs.
		// Or better, let's enforce FK in setup.
	}

	// To be safe and robust, let's assert Error, and if it fails we know we need to fix the setup or code.
	assert.Error(t, err)
	assert.Nil(t, product)
}
