package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID               uint           `gorm:"column:id_producto;primaryKey;autoIncrement"`
	CreatedAt        time.Time      `gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`

	Nombre           string    `gorm:"column:nombre;type:varchar(100);not null"`
	Descripcion      string    `gorm:"column:descripcion;type:text"`
	Precio           float64   `gorm:"column:precio;type:numeric(10,2);not null;check:precio >= 0"`
	FechaVencimiento time.Time `gorm:"column:fecha_vencimiento;type:date;not null"`
	Stock            int       `gorm:"column:stock;type:int;default:0;check:stock >= 0"`
}

func (Product) TableName() string {
	return "producto"
}