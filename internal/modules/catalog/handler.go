package catalog

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateProduct crea un nuevo producto
// POST /api/v1/catalog/products
func (h *Handler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": err.Error(),
		})
		return
	}

	product, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "error creating product",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToResponse(product)
	c.JSON(http.StatusCreated, gin.H{
		"data":    response,
		"message": "product created successfully",
	})
}

// GetProduct obtiene un producto por ID
// GET /api/v1/catalog/products/:id
func (h *Handler) GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid product id",
			"message": "id must be a valid number",
		})
		return
	}

	product, err := h.service.GetById(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "product not found",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToResponse(product)
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// UpdateProduct actualiza un producto
// PUT /api/v1/catalog/products/:id
func (h *Handler) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid product id",
			"message": "id must be a valid number",
		})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": err.Error(),
		})
		return
	}

	product, err := h.service.Update(uint(id), req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "product not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error":   "error updating product",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToResponse(product)
	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "product updated successfully",
	})
}

// DeleteProduct elimina un producto
// DELETE /api/v1/catalog/products/:id
func (h *Handler) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid product id",
			"message": "id must be a valid number",
		})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "product not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error":   "error deleting product",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product deleted successfully",
	})
}

// ListProducts lista productos con paginación
// GET /api/v1/catalog/products?page=1&limit=10
func (h *Handler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, total, err := h.service.List(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "error listing products",
			"message": err.Error(),
		})
		return
	}

	responses := make([]ProductResponse, len(products))
	for i, product := range products {
		responses[i] = h.service.ToResponse(&product)
	}

	response := ProductListResponse{
		Products: responses,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// GetProductByName obtiene un producto por nombre
// GET /api/v1/catalog/products/name/:name
func (h *Handler) GetProductByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": "product name is required",
		})
		return
	}

	product, err := h.service.GetByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "product not found",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToResponse(product)
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// GetProductsByExpirationDate obtiene productos por fecha de vencimiento
// GET /api/v1/catalog/products/expiration/:date
func (h *Handler) GetProductsByExpirationDate(c *gin.Context) {
	dateParam := c.Param("date")
	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid date format",
			"message": "date must be in format YYYY-MM-DD",
		})
		return
	}

	products, err := h.service.GetByExpirationDate(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "error getting products",
			"message": err.Error(),
		})
		return
	}

	responses := make([]ProductResponse, len(products))
	for i, product := range products {
		responses[i] = h.service.ToResponse(&product)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
	})
}

// GetExpiringSoon obtiene productos que vencen pronto
// GET /api/v1/catalog/products/expiring-soon?days=7
func (h *Handler) GetExpiringSoon(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	products, err := h.service.GetExpiringSoon(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "error getting products",
			"message": err.Error(),
		})
		return
	}

	responses := make([]ProductResponse, len(products))
	for i, product := range products {
		responses[i] = h.service.ToResponse(&product)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
	})
}

// UpdateStock actualiza el stock de un producto
// PUT /api/v1/catalog/products/:id/stock
func (h *Handler) UpdateStock(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid product id",
			"message": "id must be a valid number",
		})
		return
	}

	var req struct {
		Quantity int `json:"cantidad" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": err.Error(),
		})
		return
	}

	if err := h.service.UpdateStock(uint(id), req.Quantity); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "product not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error":   "error updating stock",
			"message": err.Error(),
		})
		return
	}

	product, _ := h.service.GetById(uint(id))
	response := h.service.ToResponse(product)

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "stock updated successfully",
	})
}
