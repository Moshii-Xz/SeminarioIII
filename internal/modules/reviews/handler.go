package reviews

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateReview(c *gin.Context) {
	var req CreateReviewRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": err.Error(),
		})
		return
	}

	review, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "error creating review",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToResponse(review)
	c.JSON(http.StatusCreated, gin.H{
		"data":    response,
		"message": "review created successfully",
	})
}

func (h *Handler) GetReview(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid review id",
			"message": "id must be a valid number",
		})
		return
	}

	review, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "review not found",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToResponse(review)
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (h *Handler) UpdateReview(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid review id",
			"message": "id must be a valid number",
		})
		return
	}

	var req UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": err.Error(),
		})
		return
	}

	review, err := h.service.Update(uint(id), req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "review not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error":   "error updating review",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToResponse(review)
	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "review updated successfully",
	})
}

func (h *Handler) DeleteReview(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid review id",
			"message": "id must be a valid number",
		})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "review not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error":   "error deleting review",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "review deleted successfully",
	})
}

func (h *Handler) ListReviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reviews, total, err := h.service.List(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "error listing reviews",
			"message": err.Error(),
		})
		return
	}

	responses := make([]ReviewResponse, len(reviews))
	for i, review := range reviews {
		responses[i] = h.service.ToResponse(&review)
	}

	response := ReviewListResponse{
		Reviews: responses,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (h *Handler) ListReviewsByProduct(c *gin.Context) {
	productIDParam := c.Param("productId")
	productID, err := strconv.ParseUint(productIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid product id",
			"message": "id must be a valid number",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reviews, total, err := h.service.ListByProduct(uint(productID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "error listing reviews",
			"message": err.Error(),
		})
		return
	}

	responses := make([]ReviewResponse, len(reviews))
	for i, review := range reviews {
		responses[i] = h.service.ToResponse(&review)
	}

	response := ReviewListResponse{
		Reviews: responses,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (h *Handler) GetProductRatingSummary(c *gin.Context) {
	productIDParam := c.Param("productId")
	productID, err := strconv.ParseUint(productIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid product id",
			"message": "id must be a valid number",
		})
		return
	}

	summary, err := h.service.GetProductRatingSummary(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "error getting rating summary",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": summary,
	})
}
