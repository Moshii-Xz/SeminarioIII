package catalog

import "time"

type CreateProductRequest struct {
	Name           string    `json:"nombre" form:"nombre" binding:"required,min=1,max=100"`
	Description    string    `json:"descripcion" form:"descripcion" binding:"omitempty"`
	ImageURL       string    `json:"imagen_url" form:"imagen_url" binding:"omitempty,max=500"`
	Price          float64   `json:"precio" form:"precio" binding:"omitempty,min=0"`
	OriginalPrice  float64   `json:"precio_original" form:"precio_original" binding:"omitempty,min=0"`
	DiscountPrice  float64   `json:"precio_descuento" form:"precio_descuento" binding:"omitempty,min=0"`
	ExpirationDate time.Time `json:"fecha_vencimiento" form:"fecha_vencimiento" time_format:"2006-01-02" binding:"required"`
	Stock          int       `json:"stock" form:"stock" binding:"omitempty,min=0"`
	Badge          *string   `json:"badge" form:"badge" binding:"omitempty,oneof=Oferta Donación"`
	CategoryID     *uint     `json:"id_categoria" form:"id_categoria" binding:"omitempty"`
	StoreID        uint      `json:"id_tienda" form:"-" binding:"-"` // Se ignora si viene en el request, se obtiene del token JWT
}

type UpdateProductRequest struct {
	Name           string    `json:"nombre" binding:"omitempty,min=1,max=100"`
	Description    string    `json:"descripcion" binding:"omitempty"`
	ImageURL       string    `json:"imagen_url" binding:"omitempty,max=500"`
	Price          float64   `json:"precio" binding:"omitempty,min=0"`
	OriginalPrice  float64   `json:"precio_original" binding:"omitempty,min=0"`
	DiscountPrice  float64   `json:"precio_descuento" binding:"omitempty,min=0"`
	ExpirationDate time.Time `json:"fecha_vencimiento" binding:"omitempty"`
	Stock          int       `json:"stock" binding:"omitempty,min=0"`
	Badge          *string   `json:"badge" binding:"omitempty,oneof=Oferta Donación"`
	CategoryID     *uint     `json:"id_categoria" binding:"omitempty"`
}

type ProductResponse struct {
	ID             uint      `json:"id_producto"`
	Name           string    `json:"nombre"`
	Description    string    `json:"descripcion"`
	ImageURL       string    `json:"imagen_url"`
	Price          float64   `json:"precio,omitempty"`
	OriginalPrice  float64   `json:"precio_original,omitempty"`
	DiscountPrice  float64   `json:"precio_descuento,omitempty"`
	ExpirationDate time.Time `json:"fecha_vencimiento"`
	Stock          int       `json:"stock"`
	Badge          *string   `json:"badge,omitempty"`
	CategoryID     *uint     `json:"id_categoria,omitempty"`
	// Status removido - ahora está en pedidos (orders)
	CategoryName *string `json:"nombre_categoria,omitempty"`
	StoreName    string  `json:"nombre_tienda,omitempty"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"productos"`
	Total    int64             `json:"total"`
	Page     int               `json:"pagina"`
	Limit    int               `json:"limite"`
}
