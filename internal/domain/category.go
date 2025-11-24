package domain

import (
	"time"
)

type Category struct {
	ID        uint      `gorm:"column:id_categoria;primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Name string `gorm:"column:nombre;type:varchar(100);uniqueIndex;not null"`
}

func (Category) TableName() string {
	return "categoria"
}

