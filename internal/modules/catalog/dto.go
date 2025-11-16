package catalog

import "time"

type CreateProductRequest struct {
	Nombre           string    `json:"nombre" binding:"required,min=1,max=100"`
	Descripcion      string    `json:"descripcion" binding:"omitempty"`
	Precio           float64   `json:"precio" binding:"required,min=0"`
	FechaVencimiento time.Time `json:"fecha_vencimiento" binding:"required"`
	Stock            int       `json:"stock" binding:"omitempty,min=0"`
}

type UpdateProductRequest struct {
	Nombre           string    `json:"nombre" binding:"omitempty,min=1,max=100"`
	Descripcion      string    `json:"descripcion" binding:"omitempty"`
	Precio           float64   `json:"precio" binding:"omitempty,min=0"`
	FechaVencimiento time.Time `json:"fecha_vencimiento" binding:"omitempty"`
	Stock            int       `json:"stock" binding:"omitempty,min=0"`
}

type ProductResponse struct {
	ID               uint      `json:"id_producto"`
	Nombre           string    `json:"nombre"`
	Descripcion      string    `json:"descripcion"`
	Precio           float64   `json:"precio"`
	FechaVencimiento time.Time `json:"fecha_vencimiento"`
	Stock            int       `json:"stock"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"productos"`
	Total    int64             `json:"total"`
	Page     int               `json:"pagina"`
	Limit    int               `json:"limite"`
}