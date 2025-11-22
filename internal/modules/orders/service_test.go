package orders_test

import (
	"testing"
	"time"

	"github.com/mordmora/expirapp/internal/domain"
	"github.com/mordmora/expirapp/internal/modules/catalog"
	"github.com/mordmora/expirapp/internal/modules/orders"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Exec("PRAGMA foreign_keys = ON")

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Client{},
		&domain.Seller{},
		&domain.Product{},
		&domain.Order{},
		&domain.OrderDetail{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupService(db *gorm.DB) *orders.Service {
	ordersRepo := orders.NewRepository(db)
	catalogRepo := catalog.NewRepository(db)
	return orders.NewService(ordersRepo, catalogRepo)
}

func createClient(db *gorm.DB) uint {
	user := domain.User{Name: "Client", Email: "client@test.com", Password: "hash"}
	db.Create(&user)
	client := domain.Client{ID: user.ID}
	db.Create(&client)
	return client.ID
}

func createSellerAndProduct(db *gorm.DB, stock int) (uint, uint) {
	user := domain.User{Name: "Seller", Email: "seller@test.com", Password: "hash"}
	db.Create(&user)
	seller := domain.Seller{ID: user.ID}
	db.Create(&seller)

	product := domain.Product{
		Name:           "Product",
		Price:          100.0,
		ExpirationDate: time.Now().Add(24 * time.Hour),
		Stock:          stock,
		SellerID:       seller.ID,
	}
	db.Create(&product)

	return seller.ID, product.ID
}

func TestCreateOrder_Success(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)
	clientID := createClient(db)
	_, productID := createSellerAndProduct(db, 10)

	req := orders.CreateOrderRequest{
		ClientID: clientID,
		Items: []orders.OrderItemRequest{
			{ProductID: productID, Quantity: 2, UnitPrice: 100.0},
		},
	}

	order, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, clientID, order.ClientID)
	assert.Len(t, order.Details, 1)

	// Verify stock reduction
	var product domain.Product
	db.First(&product, productID)
	assert.Equal(t, 8, product.Stock)
}

func TestCreateOrder_InsufficientStock(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)
	clientID := createClient(db)
	_, productID := createSellerAndProduct(db, 5)

	req := orders.CreateOrderRequest{
		ClientID: clientID,
		Items: []orders.OrderItemRequest{
			{ProductID: productID, Quantity: 10, UnitPrice: 100.0},
		},
	}

	order, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Contains(t, err.Error(), "insufficient stock")
}

func TestCreateOrder_ProductNotFound(t *testing.T) {
	db, _ := setupTestDB()
	service := setupService(db)
	clientID := createClient(db)

	req := orders.CreateOrderRequest{
		ClientID: clientID,
		Items: []orders.OrderItemRequest{
			{ProductID: 999, Quantity: 1, UnitPrice: 100.0},
		},
	}

	order, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Contains(t, err.Error(), "product with id 999 not found")
}
