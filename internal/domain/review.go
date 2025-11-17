package domain

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID        uint           `gorm:"column:id_resena;primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ProductID uint   `gorm:"column:id_producto;not null"`
	ClientID  uint   `gorm:"column:id_cliente;not null"`
	Rating    int    `gorm:"column:calificacion;type:int;not null;check:calificacion >= 1 AND calificacion <= 5"`
	Comment   string `gorm:"column:comentario;type:text"`

	Product Product `gorm:"foreignKey:ProductID;references:ID"`
	Client  User    `gorm:"foreignKey:ClientID;references:ID"`
}

func (Review) TableName() string {
	return "resena"
}

