package reports

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

func (h *Handler) GetSalesSummary(c *gin.Context) {
	var req SalesSummaryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid query params",
			"message": err.Error(),
		})
		return
	}

	summary, err := h.service.GetSalesSummary(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "could not fetch sales summary",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": summary})
}

func (h *Handler) GetTopProducts(c *gin.Context) {
	var req ReportFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid query params",
			"message": err.Error(),
		})
		return
	}
	if req.Limit == 0 {
		req.Limit = 5
	}

	products, err := h.service.GetTopProducts(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "could not fetch top products",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func (h *Handler) GetLowStock(c *gin.Context) {
	threshold, err := strconv.Atoi(c.DefaultQuery("umbral", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid threshold",
			"message": "umbral must be a number",
		})
		return
	}

	products, err := h.service.GetLowStock(threshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "could not fetch low stock products",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func (h *Handler) GetDailySales(c *gin.Context) {
	var req SalesSummaryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid query params",
			"message": err.Error(),
		})
		return
	}

	sales, err := h.service.GetDailySales(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "could not fetch daily sales",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sales})
}

func (h *Handler) GetTopCustomers(c *gin.Context) {
	var req ReportFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid query params",
			"message": err.Error(),
		})
		return
	}
	if req.Limit == 0 {
		req.Limit = 5
	}

	customers, err := h.service.GetTopCustomers(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "could not fetch top customers",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customers})
}

func (h *Handler) GetPaymentMethodSummary(c *gin.Context) {
	var req SalesSummaryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid query params",
			"message": err.Error(),
		})
		return
	}

	summary, err := h.service.GetPaymentMethodSummary(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "could not fetch payment method summary",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": summary})
}

func (h *Handler) GetPendingPayments(c *gin.Context) {
	report, err := h.service.GetPendingPayments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "could not fetch pending payments",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": report})
}
