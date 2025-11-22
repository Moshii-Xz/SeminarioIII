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
	Price          float64   `gorm:"column:precio;type:numeric(10,2);not null;check:precio >= 0"`
	ExpirationDate time.Time `gorm:"column:fecha_vencimiento;type:date;not null"`
	Stock          int       `gorm:"column:stock;default:0;check:stock >= 0"`

	// Foreign Key
	StoreID uint   `gorm:"column:id_tienda;not null"`
	Store   Store `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Product) TableName() string {
	return "producto"
}
