package domain

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID          uint           `gorm:"column:id_compra;primaryKey;autoIncrement"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	IDCliente   uint       `gorm:"column:id_cliente;not null"`
	IDVendedor  *uint      `gorm:"column:id_vendedor"`
	FechaCompra time.Time  `gorm:"column:fecha_compra;type:date;not null;default:CURRENT_DATE"`

	Items []OrderItem `gorm:"foreignKey:IDCompra;references:ID;constraint:OnDelete:CASCADE"`
}

func (Order) TableName() string {
	return "compra"
}

type OrderItem struct {
	ID             uint           `gorm:"column:id_detalle;primaryKey;autoIncrement"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	IDCompra       uint    `gorm:"column:id_compra;not null"`
	IDProducto     uint    `gorm:"column:id_producto;not null"`
	Cantidad       int     `gorm:"column:cantidad;not null;check:cantidad > 0"`
	PrecioUnitario float64 `gorm:"column:precio_unitario;type:numeric(10,2);not null;check:precio_unitario >= 0"`

	Order   Order   `gorm:"foreignKey:IDCompra;references:ID"`
	Product Product `gorm:"foreignKey:IDProducto;references:ID"`
}

func (OrderItem) TableName() string {
	return "detalle_compra"
}