package orders

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

// CreateOrder crea una nueva orden
// POST /api/v1/orders
func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": err.Error(),
		})
		return
	}

	order, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "error creating order",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToOrderResponse(order)
	c.JSON(http.StatusCreated, gin.H{
		"data":    response,
		"message": "order created successfully",
	})
}

// GetOrder obtiene una orden por ID
// GET /api/v1/orders/:id
func (h *Handler) GetOrder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid order id",
			"message": "id must be a valid number",
		})
		return
	}

	order, err := h.service.GetById(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "order not found",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToOrderResponse(order)
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// UpdateOrder actualiza una orden
// PUT /api/v1/orders/:id
func (h *Handler) UpdateOrder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid order id",
			"message": "id must be a valid number",
		})
		return
	}

	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": err.Error(),
		})
		return
	}

	order, err := h.service.Update(uint(id), req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "order not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error":   "error updating order",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToOrderResponse(order)
	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "order updated successfully",
	})
}

// DeleteOrder elimina una orden
// DELETE /api/v1/orders/:id
func (h *Handler) DeleteOrder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid order id",
			"message": "id must be a valid number",
		})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "order not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error":   "error deleting order",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order deleted successfully",
	})
}

// ListOrders lista órdenes con paginación
// GET /api/v1/orders?page=1&limit=10
func (h *Handler) ListOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, total, err := h.service.List(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "error listing orders",
			"message": err.Error(),
		})
		return
	}

	responses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = h.service.ToOrderResponse(&order)
	}

	response := OrderListResponse{
		Orders: responses,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// ListOrdersByClient lista órdenes de un cliente específico
// GET /api/v1/orders/client/:clientId?page=1&limit=10
func (h *Handler) ListOrdersByClient(c *gin.Context) {
	clientIDParam := c.Param("clientId")
	clientID, err := strconv.ParseUint(clientIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid client id",
			"message": "client id must be a valid number",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, total, err := h.service.ListByClient(uint(clientID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "error listing orders",
			"message": err.Error(),
		})
		return
	}

	responses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = h.service.ToOrderResponse(&order)
	}

	response := OrderListResponse{
		Orders: responses,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// ListOrdersBySeller lista órdenes de un vendedor específico
// GET /api/v1/orders/seller/:sellerId?page=1&limit=10
func (h *Handler) ListOrdersBySeller(c *gin.Context) {
	sellerIDParam := c.Param("sellerId")
	sellerID, err := strconv.ParseUint(sellerIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid seller id",
			"message": "seller id must be a valid number",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, total, err := h.service.ListBySeller(uint(sellerID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "error listing orders",
			"message": err.Error(),
		})
		return
	}

	responses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = h.service.ToOrderResponse(&order)
	}

	response := OrderListResponse{
		Orders: responses,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// AddOrderItem agrega un item a una orden
// POST /api/v1/orders/:id/items
func (h *Handler) AddOrderItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid order id",
			"message": "id must be a valid number",
		})
		return
	}

	var req AddOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": err.Error(),
		})
		return
	}

	item, err := h.service.AddOrderItem(uint(id), req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "order not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error":   "error adding item to order",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToOrderItemResponse(item)
	c.JSON(http.StatusCreated, gin.H{
		"data":    response,
		"message": "item added to order successfully",
	})
}

// UpdateOrderItem actualiza un item de una orden
// PUT /api/v1/orders/:id/items/:itemId
func (h *Handler) UpdateOrderItem(c *gin.Context) {
	idParam := c.Param("id")
	orderID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid order id",
			"message": "id must be a valid number",
		})
		return
	}

	itemIDParam := c.Param("itemId")
	itemID, err := strconv.ParseUint(itemIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid item id",
			"message": "item id must be a valid number",
		})
		return
	}

	var req UpdateOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": err.Error(),
		})
		return
	}

	item, err := h.service.UpdateOrderItem(uint(orderID), uint(itemID), req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "order item not found" || err.Error() == "el item no pertenece a esta orden" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error":   "error updating order item",
			"message": err.Error(),
		})
		return
	}

	response := h.service.ToOrderItemResponse(item)
	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "order item updated successfully",
	})
}

// DeleteOrderItem elimina un item de una orden
// DELETE /api/v1/orders/:id/items/:itemId
func (h *Handler) DeleteOrderItem(c *gin.Context) {
	idParam := c.Param("id")
	orderID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid order id",
			"message": "id must be a valid number",
		})
		return
	}

	itemIDParam := c.Param("itemId")
	itemID, err := strconv.ParseUint(itemIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid item id",
			"message": "item id must be a valid number",
		})
		return
	}

	if err := h.service.DeleteOrderItem(uint(orderID), uint(itemID)); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "order item not found" || err.Error() == "el item no pertenece a esta orden" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error":   "error deleting order item",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order item deleted successfully",
	})
}
