package payments

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

// CreatePayment crea un nuevo pago
// POST /api/v1/payments
func (h *Handler) CreatePayment(c *gin.Context) {
	var req CreatePaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": err.Error(),
		})
		return
	}

	payment, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "error creating payment",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToPaymentResponse(payment)
	c.JSON(http.StatusCreated, gin.H{
		"data":    response,
		"message": "payment created successfully",
	})
}

// GetPayment obtiene un pago por ID
// GET /api/v1/payments/:id
func (h *Handler) GetPayment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid payment id",
			"message": "id must be a valid number",
		})
		return
	}

	payment, err := h.service.GetById(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "payment not found",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToPaymentResponse(payment)
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// UpdatePayment actualiza un pago
// PUT /api/v1/payments/:id
func (h *Handler) UpdatePayment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid payment id",
			"message": "id must be a valid number",
		})
		return
	}

	var req UpdatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": err.Error(),
		})
		return
	}

	payment, err := h.service.Update(uint(id), req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "payment not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error":   "error updating payment",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToPaymentResponse(payment)
	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "payment updated successfully",
	})
}

// DeletePayment elimina un pago
// DELETE /api/v1/payments/:id
func (h *Handler) DeletePayment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid payment id",
			"message": "id must be a valid number",
		})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "payment not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error":   "error deleting payment",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "payment deleted successfully",
	})
}

// ListPayments lista pagos con paginación
// GET /api/v1/payments?page=1&limit=10
func (h *Handler) ListPayments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	payments, total, err := h.service.List(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "error listing payments",
			"message": err.Error(),
		})
		return
	}

	responses := make([]PaymentResponse, len(payments))
	for i, payment := range payments {
		responses[i] = h.service.ToPaymentResponse(&payment)
	}

	response := PaymentListResponse{
		Payments: responses,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// GetPaymentsByOrder obtiene todos los pagos de una orden
// GET /api/v1/payments/order/:orderId
func (h *Handler) GetPaymentsByOrder(c *gin.Context) {
	orderIDParam := c.Param("orderId")
	orderID, err := strconv.ParseUint(orderIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid order id",
			"message": "order id must be a valid number",
		})
		return
	}

	payments, err := h.service.GetByOrderID(uint(orderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "error getting payments",
			"message": err.Error(),
		})
		return
	}

	responses := make([]PaymentResponse, len(payments))
	for i, payment := range payments {
		responses[i] = h.service.ToPaymentResponse(&payment)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
	})
}

// GetPaymentStatusByOrder obtiene el estado de pagos de una orden
// GET /api/v1/payments/order/:orderId/status
func (h *Handler) GetPaymentStatusByOrder(c *gin.Context) {
	orderIDParam := c.Param("orderId")
	orderID, err := strconv.ParseUint(orderIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid order id",
			"message": "order id must be a valid number",
		})
		return
	}

	status, err := h.service.GetPaymentStatusByOrderID(uint(orderID))
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "orden con id "+orderIDParam+" no encontrada" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error":   "error getting payment status",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": status,
	})
}

// CreatePaymentMethod crea un nuevo método de pago
// POST /api/v1/payments/methods
func (h *Handler) CreatePaymentMethod(c *gin.Context) {
	var req struct {
		Nombre string `json:"nombre" binding:"required,min=1,max=50"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": err.Error(),
		})
		return
	}

	method, err := h.service.CreatePaymentMethod(req.Nombre)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "error creating payment method",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToPaymentMethodResponse(method)
	c.JSON(http.StatusCreated, gin.H{
		"data":    response,
		"message": "payment method created successfully",
	})
}

// GetPaymentMethod obtiene un método de pago por ID
// GET /api/v1/payments/methods/:id
func (h *Handler) GetPaymentMethod(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid payment method id",
			"message": "id must be a valid number",
		})
		return
	}

	method, err := h.service.GetPaymentMethodByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "payment method not found",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToPaymentMethodResponse(method)
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// ListPaymentMethods lista todos los métodos de pago
// GET /api/v1/payments/methods
func (h *Handler) ListPaymentMethods(c *gin.Context) {
	methods, err := h.service.ListPaymentMethods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "error listing payment methods",
			"message": err.Error(),
		})
		return
	}

	responses := make([]PaymentMethodResponse, len(methods))
	for i, method := range methods {
		responses[i] = h.service.ToPaymentMethodResponse(&method)
	}

	response := PaymentMethodListResponse{
		Methods: responses,
		Total:   int64(len(responses)),
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// UpdatePaymentMethod actualiza un método de pago
// PUT /api/v1/payments/methods/:id
func (h *Handler) UpdatePaymentMethod(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid payment method id",
			"message": "id must be a valid number",
		})
		return
	}

	var req struct {
		Nombre string `json:"nombre" binding:"required,min=1,max=50"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": err.Error(),
		})
		return
	}

	method, err := h.service.UpdatePaymentMethod(uint(id), req.Nombre)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "payment method not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error":   "error updating payment method",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToPaymentMethodResponse(method)
	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "payment method updated successfully",
	})
}

// DeletePaymentMethod elimina un método de pago
// DELETE /api/v1/payments/methods/:id
func (h *Handler) DeletePaymentMethod(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid payment method id",
			"message": "id must be a valid number",
		})
		return
	}

	if err := h.service.DeletePaymentMethod(uint(id)); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "payment method not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error":   "error deleting payment method",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "payment method deleted successfully",
	})
}
