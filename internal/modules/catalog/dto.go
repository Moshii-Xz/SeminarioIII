package catalog

import "time"

type CreateProductRequest struct {
	Name           string    `json:"nombre" binding:"required,min=1,max=100"`
	Description    string    `json:"descripcion" binding:"omitempty"`
	Price          float64   `json:"precio" binding:"required,min=0"`
	ExpirationDate time.Time `json:"fecha_vencimiento" binding:"required"`
	Stock          int       `json:"stock" binding:"omitempty,min=0"`
	SellerID       uint      `json:"id_vendedor" binding:"required"`
}

type UpdateProductRequest struct {
	Name           string    `json:"nombre" binding:"omitempty,min=1,max=100"`
	Description    string    `json:"descripcion" binding:"omitempty"`
	Price          float64   `json:"precio" binding:"omitempty,min=0"`
	ExpirationDate time.Time `json:"fecha_vencimiento" binding:"omitempty"`
	Stock          int       `json:"stock" binding:"omitempty,min=0"`
}

type ProductResponse struct {
	ID             uint      `json:"id_producto"`
	Name           string    `json:"nombre"`
	Description    string    `json:"descripcion"`
	Price          float64   `json:"precio"`
	ExpirationDate time.Time `json:"fecha_vencimiento"`
	Stock          int       `json:"stock"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"productos"`
	Total    int64             `json:"total"`
	Page     int               `json:"pagina"`
	Limit    int               `json:"limite"`
}
