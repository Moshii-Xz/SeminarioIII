package catalog

import (
	"fmt"
	"net/http"
	"path/filepath"
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

func (h *Handler) Create(c *gin.Context) {
	var req CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": err.Error(),
		})
		return
	}

	// Obtener el userID del token JWT (establecido por el middleware de autenticación)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "user ID not found in token",
		})
		return
	}

	// Convertir userID a uint (el userID del token es el id_tienda para usuarios con rol "tienda")
	storeID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal error",
			"message": "invalid user ID type",
		})
		return
	}

	product, err := h.service.Create(req, storeID)
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

func (h *Handler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	product, err := h.service.GetById(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": h.service.ToResponse(product)})
}

func (h *Handler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Obtener el userID del token JWT
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "user ID not found in token",
		})
		return
	}

	storeID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal error",
			"message": "invalid user ID type",
		})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.Update(uint(id), req, storeID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "no tienes permiso para modificar este producto" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "product not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": h.service.ToResponse(product)})
}

func (h *Handler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Obtener el userID del token JWT
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "user ID not found in token",
		})
		return
	}

	storeID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal error",
			"message": "invalid user ID type",
		})
		return
	}

	if err := h.service.Delete(uint(id), storeID); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "no tienes permiso para eliminar este producto" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "product not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, total, err := h.service.List(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]ProductResponse, len(products))
	for i, p := range products {
		responses[i] = h.service.ToResponse(&p)
	}

	c.JSON(http.StatusOK, ProductListResponse{
		Products: responses,
		Total:    total,
		Page:     page,
		Limit:    limit,
	})
}

// UploadImage handles image upload for products
func (h *Handler) UploadImage(c *gin.Context) {
	// Single file
	file, err := c.FormFile("imagen")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid file",
			"message": "imagen field is required",
		})
		return
	}

	// Validate file type
	ext := filepath.Ext(file.Filename)
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	validExt := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			validExt = true
			break
		}
	}
	if !validExt {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid file type",
			"message": "only jpg, jpeg, png, gif, and webp are allowed",
		})
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "file too large",
			"message": "file size must be less than 5MB",
		})
		return
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	
	// Save file to uploads/images/ directory
	// In Docker, working directory is /root/, so we use relative path
	uploadPath := filepath.Join("uploads", "images", filename)
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to save file",
			"message": err.Error(),
		})
		return
	}

	// Return the URL (accessible via static file server)
	imageURL := fmt.Sprintf("/uploads/images/%s", filename)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"imagen_url": imageURL,
			"filename":   filename,
		},
		"message": "image uploaded successfully",
	})
}
