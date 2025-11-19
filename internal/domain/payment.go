package domain

import (
	"time"
)

type PaymentMethod struct {
	ID   uint   `gorm:"column:id_metodo_pago;primaryKey;autoIncrement"`
	Name string `gorm:"column:nombre;size:50;not null;unique"`
}

func (PaymentMethod) TableName() string {
	return "metodo_pago"
}

type Payment struct {
	ID              uint           `gorm:"column:id_pago;primaryKey;autoIncrement"`
	OrderID         uint           `gorm:"column:id_compra;not null"`
	Order           *Order         `gorm:"foreignKey:OrderID;references:ID"`
	PaymentMethodID *uint          `gorm:"column:id_metodo_pago"`
	PaymentMethod   *PaymentMethod `gorm:"foreignKey:PaymentMethodID;references:ID"`
	Amount          float64        `gorm:"column:monto;type:numeric(10,2);not null;check:monto >= 0"`
	Date            time.Time      `gorm:"column:fecha_pago;default:CURRENT_DATE;not null"`
}

func (Payment) TableName() string {
	return "pago"
}
