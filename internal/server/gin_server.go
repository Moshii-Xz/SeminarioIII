package server

/*
Este archivo contiene la implementación del servidor Gin
*/

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Moshii-Xz/SeminarioIII/internal/middleware"
	"github.com/Moshii-Xz/SeminarioIII/internal/modules/catalog"
	"github.com/Moshii-Xz/SeminarioIII/internal/modules/orders"
	"github.com/Moshii-Xz/SeminarioIII/internal/modules/payments"
	"github.com/Moshii-Xz/SeminarioIII/internal/modules/reports"
	"github.com/Moshii-Xz/SeminarioIII/internal/modules/reviews"
	"github.com/Moshii-Xz/SeminarioIII/internal/modules/stores"
	"github.com/Moshii-Xz/SeminarioIII/internal/modules/users"
	"gorm.io/gorm"
)

/*
# Server representa el servidor Gin
* httpServer: el servidor HTTP subyacente
* engine: el motor Gin
* db: la conexión a la base de datos
* config: la configuración del servidor
*/

type Server struct {
	httpServer *http.Server
	engine     *gin.Engine
	db         *gorm.DB
	config     Config
}

/*
# New crea una nueva instancia del servidor Gin
* db: la conexión a la base de datos
* cfg: la configuración del servidor
*/
func New(db *gorm.DB, cfg Config) *Server {
	gin.SetMode(cfg.Mode)

	engine := gin.New()

	/*
		# Crea una instancia de un servidorr gin
	*/
	server := &Server{
		engine: engine,
		db:     db,
		config: cfg,
		httpServer: &http.Server{
			Addr:         cfg.Port,
			Handler:      engine,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}

	server.setupMiddlewares()

	server.setupRouter()

	return server
}

func (s *Server) setupMiddlewares() {
	/*
		# evita que el servidor reviente
		# si hay un panic en alguna parte del código
	*/
	s.engine.Use(gin.Recovery())

	/*
		# registra las solicitudes entrantes
		# solo para depurar
	*/
	s.engine.Use(gin.Logger())
	s.engine.Use(corsMiddleware())
}

/*
# corsMiddleware maneja las solicitudes CORS
*/
func corsMiddleware() gin.HandlerFunc {
	/*
		# permite solicitudes CORS desde cualquier origen
		# con los métodos y encabezados especificados
	*/
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *Server) setupRouter() {

	s.engine.GET("/health", s.healthCheck)

	v1 := s.engine.Group("/api/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "API v1 is working!",
				"status":  "running",
			})
		})

		// Initialize Modules

		// Users
		usersRepo := users.NewRepository(s.db)
		usersService := users.NewService(usersRepo)
		usersHandler := users.NewHandler(usersService)

		usersGroup := v1.Group("/users")
		{
			usersGroup.POST("/", usersHandler.Create)
			usersGroup.POST("/login", usersHandler.Login)

			protected := usersGroup.Group("/")
			protected.Use(middleware.Auth())
			{
				protected.GET("/:id", usersHandler.GetByID)
				protected.PUT("/:id", usersHandler.Update)
				protected.DELETE("/:id", usersHandler.Delete)
				protected.GET("/", usersHandler.List)
			}
		}

		// Catalog
		catalogRepo := catalog.NewRepository(s.db)
		catalogService := catalog.NewService(catalogRepo)
		catalogHandler := catalog.NewHandler(catalogService)

		catalogGroup := v1.Group("/products")
		{
			catalogGroup.GET("/:id", catalogHandler.GetByID)
			catalogGroup.GET("/", catalogHandler.List)

			protected := catalogGroup.Group("/")
			protected.Use(middleware.Auth())
			{
				protected.POST("/", catalogHandler.Create)
				protected.POST("/upload-image", catalogHandler.UploadImage)
				protected.PUT("/:id", catalogHandler.Update)
				protected.DELETE("/:id", catalogHandler.Delete)
			}
		}

		// Serve static files (uploaded images)
		s.engine.Static("/uploads", "./uploads")

		// Orders
		ordersRepo := orders.NewRepository(s.db)
		ordersService := orders.NewService(ordersRepo, catalogRepo)
		ordersHandler := orders.NewHandler(ordersService)

		ordersGroup := v1.Group("/orders")
		ordersGroup.Use(middleware.Auth())
		{
			ordersGroup.POST("/", ordersHandler.Create)
			ordersGroup.GET("/:id", ordersHandler.GetByID)
			ordersGroup.PUT("/:id", ordersHandler.Update)
			ordersGroup.DELETE("/:id", ordersHandler.Delete)
			ordersGroup.GET("/", ordersHandler.List)
			ordersGroup.POST("/:id/items", ordersHandler.AddDetail)
			ordersGroup.PUT("/:id/items/:itemId", ordersHandler.UpdateDetail)
			ordersGroup.DELETE("/:id/items/:itemId", ordersHandler.RemoveDetail)
		}

		// Payments
		paymentsRepo := payments.NewRepository(s.db)
		paymentsService := payments.NewService(paymentsRepo, ordersRepo)
		paymentsHandler := payments.NewHandler(paymentsService)

		paymentsGroup := v1.Group("/payments")
		paymentsGroup.Use(middleware.Auth())
		{
			paymentsGroup.POST("/", paymentsHandler.Create)
			paymentsGroup.GET("/:id", paymentsHandler.GetByID)
			paymentsGroup.PUT("/:id", paymentsHandler.Update)
			paymentsGroup.DELETE("/:id", paymentsHandler.Delete)
			paymentsGroup.GET("/", paymentsHandler.List)
			paymentsGroup.GET("/order/:orderId", paymentsHandler.GetByOrder)
			paymentsGroup.GET("/order/:orderId/status", paymentsHandler.GetPaymentStatusByOrder)

			// Payment Methods
			paymentsGroup.POST("/methods", paymentsHandler.CreatePaymentMethod)
			paymentsGroup.GET("/methods", paymentsHandler.ListPaymentMethods)
			paymentsGroup.GET("/methods/:id", paymentsHandler.GetPaymentMethod)
			paymentsGroup.PUT("/methods/:id", paymentsHandler.UpdatePaymentMethod)
			paymentsGroup.DELETE("/methods/:id", paymentsHandler.DeletePaymentMethod)
		}

		// Reviews
		reviewsRepo := reviews.NewRepository(s.db)
		reviewsService := reviews.NewService(reviewsRepo)
		reviewsHandler := reviews.NewHandler(reviewsService)

		reviewsGroup := v1.Group("/reviews")
		{
			reviewsGroup.GET("/:id", reviewsHandler.GetByID)
			reviewsGroup.GET("/", reviewsHandler.List)
			reviewsGroup.GET("/product/:productId", reviewsHandler.ListByProduct)

			protected := reviewsGroup.Group("/")
			protected.Use(middleware.Auth())
			{
				protected.POST("/", reviewsHandler.Create)
				protected.PUT("/:id", reviewsHandler.Update)
				protected.DELETE("/:id", reviewsHandler.Delete)
			}
		}

		// Reports
		reportsRepo := reports.NewRepository(s.db)
		reportsService := reports.NewService(reportsRepo)
		reportsHandler := reports.NewHandler(reportsService)

		reportsGroup := v1.Group("/reports")
		reportsGroup.Use(middleware.Auth())
		{
			reportsGroup.GET("/sales", reportsHandler.GetSalesReport)
			reportsGroup.GET("/stock", reportsHandler.GetStockReport)
			reportsGroup.GET("/expiring", reportsHandler.GetExpiringProductsReport)
		}

		// Stores
		storesRepo := stores.NewRepository(s.db)
		storesService := stores.NewServiceWithUserRepo(storesRepo, usersRepo)
		storesHandler := stores.NewHandler(storesService)

		storesGroup := v1.Group("/stores")
		{
			// Public endpoint for store registration
			storesGroup.POST("/", storesHandler.Create)

			// Protected endpoints
			protected := storesGroup.Group("/")
			protected.Use(middleware.Auth())
			{
				protected.GET("/", storesHandler.List)
				protected.GET("/:id", storesHandler.GetByID)
				protected.PUT("/:id", storesHandler.Update)
				protected.GET("/:id/products", storesHandler.GetProducts)
				protected.GET("/:id/orders", storesHandler.GetOrders)
			}
		}
	}

}

func (s *Server) healthCheck(c *gin.Context) {

	sqlDB, err := s.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "database connection failed"})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "database ping failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "dabase is ok"})

}

func (s *Server) Start() error {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on %s", s.config.Port)

		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", s.config.Port, err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
	return nil
}

/*
# para debug
*/
func (s *Server) Engine() *gin.Engine {
	return s.engine
}
