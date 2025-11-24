package catalog

import "time"

type CreateProductRequest struct {
	Name           string    `json:"nombre" binding:"required,min=1,max=100"`
	Description    string    `json:"descripcion" binding:"omitempty"`
	ImageURL       string    `json:"imagen_url" binding:"omitempty,max=500"`
	Price          float64   `json:"precio" binding:"required,min=0"`
	ExpirationDate time.Time `json:"fecha_vencimiento" binding:"required"`
	Stock          int       `json:"stock" binding:"omitempty,min=0"`
	Status         string    `json:"estado" binding:"omitempty"`
	StoreID       uint      `json:"id_tienda" binding:"required"`
}

type UpdateProductRequest struct {
	Name           string    `json:"nombre" binding:"omitempty,min=1,max=100"`
	Description    string    `json:"descripcion" binding:"omitempty"`
	ImageURL       string    `json:"imagen_url" binding:"omitempty,max=500"`
	Price          float64   `json:"precio" binding:"omitempty,min=0"`
	ExpirationDate time.Time `json:"fecha_vencimiento" binding:"omitempty"`
	Stock          int       `json:"stock" binding:"omitempty,min=0"`
	Status         string    `json:"estado" binding:"omitempty"`
}

type ProductResponse struct {
	ID             uint      `json:"id_producto"`
	Name           string    `json:"nombre"`
	Description    string    `json:"descripcion"`
	ImageURL       string    `json:"imagen_url"`
	Price          float64   `json:"precio"`
	ExpirationDate time.Time `json:"fecha_vencimiento"`
	Stock          int       `json:"stock"`
	Status         string    `json:"estado"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"productos"`
	Total    int64             `json:"total"`
	Page     int               `json:"pagina"`
	Limit    int               `json:"limite"`
}
