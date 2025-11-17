package domain

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID          uint           `gorm:"column:id_pago;primaryKey;autoIncrement"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	IDCompra     uint       `gorm:"column:id_compra;not null"`
	IDMetodoPago *uint      `gorm:"column:id_metodo_pago"`
	Monto        float64    `gorm:"column:monto;type:numeric(10,2);not null;check:monto >= 0"`
	FechaPago    time.Time  `gorm:"column:fecha_pago;type:date;not null;default:CURRENT_DATE"`

	Order        Order         `gorm:"foreignKey:IDCompra;references:ID"`
	PaymentMethod *PaymentMethod `gorm:"foreignKey:IDMetodoPago;references:ID"`
}

func (Payment) TableName() string {
	return "pago"
}

type PaymentMethod struct {
	ID     uint   `gorm:"column:id_metodo_pago;primaryKey;autoIncrement"`
	Nombre string `gorm:"column:nombre;type:varchar(50);not null;uniqueIndex"`

	Payments []Payment `gorm:"foreignKey:IDMetodoPago;references:ID"`
}

func (PaymentMethod) TableName() string {
	return "metodo_pago"
}