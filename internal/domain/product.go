package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID        uint           `gorm:"column:id_producto;primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name           string    `gorm:"column:nombre;type:varchar(100);not null"`
	Description    string    `gorm:"column:descripcion;type:text"`
	ImageURL       string    `gorm:"column:imagen_url;type:varchar(500)"`
	Price          float64   `gorm:"column:precio;type:numeric(10,2);not null;check:precio >= 0"`
	ExpirationDate time.Time `gorm:"column:fecha_vencimiento;type:date;not null"`
	Stock          int       `gorm:"column:stock;default:0;check:stock >= 0"`
	Status      string    `gorm:"column:estado;type:varchar(50);default:En preparación"`

	// Foreign Keys
	StoreID    uint     `gorm:"column:id_tienda;not null"`
	Store      Store    `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CategoryID *uint    `gorm:"column:id_categoria"`
	Category   *Category `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (Product) TableName() string {
	return "producto"
}
