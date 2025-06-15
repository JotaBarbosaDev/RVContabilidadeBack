package models

import (
	"time"
)

// User representa um utilizador do sistema
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
    Email     string    `json:"email" gorm:"unique;not null" example:"joao@exemplo.com"`
    Password  string    `json:"-" gorm:"not null"`
    Name      string    `json:"name" gorm:"not null" example:"João Silva"` 
	Username  string    `json:"username" gorm:"unique;not null" example:"joaosilva"` // nome de utilizador único
    Role      string    `json:"role" gorm:"default:'user'" example:"user"` //cliente - 0, contabilista (gestor) - 1, admin - 2
	IsActive  bool      `json:"is_active" gorm:"default:true" example:"true"` // indica se tem permissão para aceder à aplicação, pode estar bloqueado
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime" example:"2023-01-01T00:00:00Z"` // data de criação do utilizador
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2023-01-01T00:00:00Z"` // data da última atualização do utilizador
}

// Dados para login
type LoginRequest struct {
    Username string `json:"username" binding:"required,min=2" example:"joaosilva"`
    Password string `json:"password" binding:"required,min=6" example:"123456"`
}

// Dados para registo
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email" example:"joao@exemplo.com"`
    Password string `json:"password" binding:"required,min=6" example:"123456"`
    Name     string `json:"name" binding:"required,min=2" example:"João Silva"`
    Username string `json:"username" binding:"required,min=2" example:"joaosilva"`
	Role     string `json:"role" binding:"required,oneof=client accountant admin" example:"user"`
	IsActive bool   `json:"is_active" binding:"required" example:"true"`
}

// Resposta com token
type AuthResponse struct {
    Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
    User  User   `json:"user"`
}

// Resposta padrão de sucesso
type SuccessResponse struct {
    Success bool        `json:"success" example:"true"`
    Message string      `json:"message" example:"Operação realizada com sucesso"`
    Data    interface{} `json:"data,omitempty"`
}

// Resposta padrão de erro
type ErrorResponse struct {
    Success bool   `json:"success" example:"false"`
    Error   string `json:"error" example:"Descrição do erro"`
}
