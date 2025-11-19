package users

import "time"

/*
#CreateUserRequest; body para POST /users
*/
type CreateUserRequest struct {
	Name     string `json:"nombre" binding:"required,min=2,max=100"`
	Email    string `json:"correo" binding:"required,email"`
	Password string `json:"contrasena" binding:"required,min=6"`
}

/*
# UpdateUserRequest; body para PUT /users/:id
*/
type UpdateUserRequest struct {
	Name  string `json:"nombre" binding:"omitempty,min=2,max=100"`
	Email string `json:"correo" binding:"omitempty,email"`
}

/*
# ChangePasswordRequest; body para PUT /users/:id/password
*/
type ChangePasswordRequest struct {
	CurrentPass string `json:"contrasena_actual" binding:"required"`
	NewPass     string `json:"contrasena_nueva" binding:"required,min=6"`
}

type UserResponse struct {
	ID        uint      `json:"id_usuario"`
	Name      string    `json:"nombre"`
	Email     string    `json:"correo"`
	CreatedAt time.Time `json:"fecha_registro"`
}

type UserListResponse struct {
	Users []UserResponse `json:"usuarios"`
	Total int64          `json:"total"`
	Page  int            `json:"pagina"`
	Limit int            `json:"limite"`
}

type LoginRequest struct {
	Email    string `json:"correo" binding:"required,email"`
	Password string `json:"contrasena" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"usuario"`
}
