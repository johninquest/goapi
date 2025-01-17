// internal/dto/auth_dto.go
package dto

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
    Name     string `json:"name" validate:"required"`
}

type AuthResponse struct {
    Token string `json:"token"`
}