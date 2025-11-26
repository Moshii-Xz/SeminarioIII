package domain

import (
	"time"
)

type Order struct {
	ID            uint      `gorm:"column:id_compra;primaryKey;autoIncrement"`
	ClientID      uint      `gorm:"column:id_cliente;not null"`
	StoreID       *uint     `gorm:"column:id_tienda"`
	Date          time.Time `gorm:"column:fecha_compra;default:CURRENT_DATE;not null"`
	Status        string    `gorm:"column:estado;type:varchar(50);default:En preparación"`
	PaymentMethod string    `gorm:"column:metodo_pago;type:varchar(50)"`

	// Relaciones
	Details []OrderDetail `gorm:"foreignKey:OrderID;references:ID"`
}

func (Order) TableName() string {
	return "compra"
}

type OrderDetail struct {
	ID        uint    `gorm:"column:id_detalle;primaryKey;autoIncrement"`
	OrderID   uint    `gorm:"column:id_compra;not null"`
	ProductID uint    `gorm:"column:id_producto;not null"`
	Quantity  int     `gorm:"column:cantidad;not null;check:cantidad > 0"`
	UnitPrice float64 `gorm:"column:precio_unitario;type:numeric(10,2);not null;check:precio_unitario >= 0"`

	// Relaciones
	Product Product `gorm:"foreignKey:ProductID;references:ID"`
}

func (OrderDetail) TableName() string {
	return "detalle_compra"
}
