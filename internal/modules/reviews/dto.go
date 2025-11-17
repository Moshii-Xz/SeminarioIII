package reviews

import "time"

type CreateReviewRequest struct {
	ProductID uint   `json:"id_producto" binding:"required"`
	ClientID  uint   `json:"id_cliente" binding:"required"`
	Rating    int    `json:"calificacion" binding:"required,min=1,max=5"`
	Comment   string `json:"comentario" binding:"omitempty,max=500"`
}

type UpdateReviewRequest struct {
	Rating  int    `json:"calificacion" binding:"omitempty,min=1,max=5"`
	Comment string `json:"comentario" binding:"omitempty,max=500"`
}

type ReviewResponse struct {
	ID        uint      `json:"id_resena"`
	ProductID uint      `json:"id_producto"`
	ClientID  uint      `json:"id_cliente"`
	Rating    int       `json:"calificacion"`
	Comment   string    `json:"comentario"`
	CreatedAt time.Time `json:"fecha_creacion"`
	UpdatedAt time.Time `json:"fecha_actualizacion"`
}

type ReviewListResponse struct {
	Reviews []ReviewResponse `json:"resenas"`
	Total   int64            `json:"total"`
	Page    int              `json:"pagina"`
	Limit   int              `json:"limite"`
}

type ProductRatingSummary struct {
	ProductID    uint    `json:"id_producto"`
	AverageRating float64 `json:"calificacion_promedio"`
	ReviewsCount int64   `json:"total_resenas"`
}
