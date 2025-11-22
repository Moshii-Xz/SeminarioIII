package stores

import "time"

type UpdateStoreRequest struct {
	ResponsibleArea string `json:"area_responsable" binding:"omitempty,max=100"`
	Address         string `json:"direccion" binding:"omitempty,max=200"`
	Phone           string `json:"telefono" binding:"omitempty,max=20"`
}

type StoreResponse struct {
	ID              uint      `json:"id_tienda"`
	ResponsibleArea string    `json:"area_responsable"`
	Address         string    `json:"direccion"`
	Phone           string    `json:"telefono"`
	User            UserInfo  `json:"usuario"`
	CreatedAt       time.Time `json:"fecha_creacion"`
}

type UserInfo struct {
	ID        uint      `json:"id_usuario"`
	Name      string    `json:"nombre"`
	Email     string    `json:"correo"`
	CreatedAt time.Time `json:"fecha_registro"`
}

type StoreListResponse struct {
	Stores []StoreResponse `json:"tiendas"`
	Total  int64           `json:"total"`
	Page   int             `json:"pagina"`
	Limit  int             `json:"limite"`
}

type StoreProductsResponse struct {
	Store    StoreResponse `json:"tienda"`
	Products []ProductInfo `json:"productos"`
	Total    int64         `json:"total"`
}

type ProductInfo struct {
	ID             uint      `json:"id_producto"`
	Name           string    `json:"nombre"`
	Description    string    `json:"descripcion"`
	Price          float64   `json:"precio"`
	ExpirationDate time.Time `json:"fecha_vencimiento"`
	Stock          int       `json:"stock"`
}

type StoreOrdersResponse struct {
	Store  StoreResponse `json:"tienda"`
	Orders []OrderInfo   `json:"ordenes"`
	Total  int64         `json:"total"`
}

type OrderInfo struct {
	ID       uint      `json:"id"`
	ClientID uint      `json:"client_id"`
	Date     time.Time `json:"date"`
	Total    float64   `json:"total"`
}

