package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"column:id_usuario;primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name     string `gorm:"column:nombre;type:varchar(100);not null"`
	Email    string `gorm:"column:correo;type:varchar(100);uniqueIndex;not null"`
	Password string `gorm:"column:contrasena;type:varchar(255);not null"`

	Roles []Role `gorm:"many2many:usuario_rol;foreignKey:ID;joinForeignKey:id_usuario;References:ID;joinReferences:id_rol"`
}

func (User) TableName() string {
	return "usuario"
}

type Client struct {
	ID      uint   `gorm:"column:id_cliente;primaryKey"`
	Address string `gorm:"column:direccion"`
	Phone   string `gorm:"column:telefono"`
}

func (Client) TableName() string {
	return "cliente"
}

type Store struct {
	ID              uint   `gorm:"column:id_tienda;primaryKey"`
	ResponsibleArea string `gorm:"column:area_responsable"`
	Address         string `gorm:"column:direccion"`
	Phone           string `gorm:"column:telefono"`
}

func (Store) TableName() string {
	return "tienda"
}

type Admin struct {
	ID              uint   `gorm:"column:id_admin;primaryKey"`
	SpecialPermissions string `gorm:"column:permisos_especiales;type:text"`
}

func (Admin) TableName() string {
	return "administrador"
}