package models

import (
	"time"
)

// User representa um utilizador do sistema
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
    Username  string    `json:"username" gorm:"unique" example:"joao.silva"`
    Email     string    `json:"email" gorm:"unique;not null" example:"joao@exemplo.com"`
    Password  string    `json:"-" gorm:"not null"`
    Name      string    `json:"name" gorm:"not null" example:"João Silva"` 
	Phone     string    `json:"phone" example:"912345678"`
	NIF       string    `json:"nif" gorm:"unique;not null" example:"123456789"` // Chave única para detetar duplicados
    Role      string    `json:"role" gorm:"default:'client'" example:"client"` // client, accountant, admin
	Status    string    `json:"status" gorm:"default:'pending'" example:"pending"` // pending, approved, rejected, blocked
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime" example:"2023-01-01T00:00:00Z"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2023-01-01T00:00:00Z"`
	
	// Relacionamentos
	Companies []Company `json:"companies,omitempty" gorm:"foreignKey:UserID"`
	Requests  []RegistrationRequest `json:"requests,omitempty" gorm:"foreignKey:UserID"`
}

// UserStatus constants
type UserStatus string

const (
	StatusPending  UserStatus = "pending"   // Aguarda aprovação
	StatusApproved UserStatus = "approved"  // Aprovado, pode aceder
	StatusRejected UserStatus = "rejected"  // Rejeitado, sem acesso
	StatusBlocked  UserStatus = "blocked"   // Bloqueado
)

// Company representa uma empresa
type Company struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	UserID           uint      `json:"user_id" gorm:"not null"`
	CompanyName      string    `json:"company_name" gorm:"not null"`
	TradeName        string    `json:"trade_name"`
	NIPC             string    `json:"nipc" gorm:"unique;not null"`
	Address          string    `json:"address"`
	PostalCode       string    `json:"postal_code"`
	City             string    `json:"city"`
	Country          string    `json:"country" gorm:"default:'Portugal'"`
	CAE              string    `json:"cae"`
	LegalForm        string    `json:"legal_form"`
	ShareCapital     float64   `json:"share_capital"`
	RegistrationDate time.Time `json:"registration_date"`
	Status           string    `json:"status" gorm:"default:'active'"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// RegistrationRequest representa o histórico de solicitações
type RegistrationRequest struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	UserID            uint      `json:"user_id" gorm:"not null"`
	RequestType       string    `json:"request_type" gorm:"default:'new_client'"` // new_client, existing_client
	RequestData       string    `json:"request_data" gorm:"type:text"`
	Status            string    `json:"status" gorm:"default:'pending'"`
	SubmittedAt       time.Time `json:"submitted_at" gorm:"autoCreateTime"`
	ReviewedAt        *time.Time `json:"reviewed_at"`
	ReviewedBy        *uint     `json:"reviewed_by"`
	ReviewNotes       string    `json:"review_notes"`
	ApprovalToken     string    `json:"approval_token" gorm:"unique"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	User           User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ReviewedByUser *User `json:"reviewed_by_user,omitempty" gorm:"foreignKey:ReviewedBy"`
}

// DTOs para Request/Response

// RegistrationRequestDTO para registo de clientes
type RegistrationRequestDTO struct {
	// Dados Pessoais
	Username string `json:"username" binding:"required" example:"joao.silva"`
	Name     string `json:"name" binding:"required" example:"João Silva"`
	Email    string `json:"email" binding:"required,email" example:"joao@exemplo.com"`
	Phone    string `json:"phone" binding:"required" example:"912345678"`
	NIF      string `json:"nif" binding:"required" example:"123456789"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
	
	// Dados da Empresa
	CompanyName      string  `json:"company_name" binding:"required" example:"Silva & Associados Lda"`
	TradeName        string  `json:"trade_name" example:"Silva Consultoria"`
	NIPC             string  `json:"nipc" binding:"required" example:"123456789"`
	Address          string  `json:"address" binding:"required" example:"Rua das Flores, 123"`
	PostalCode       string  `json:"postal_code" binding:"required" example:"1000-001"`
	City             string  `json:"city" binding:"required" example:"Lisboa"`
	Country          string  `json:"country" example:"Portugal"`
	CAE              string  `json:"cae" example:"69200"`
	LegalForm        string  `json:"legal_form" example:"Sociedade por Quotas"`
	ShareCapital     float64 `json:"share_capital" example:"5000.00"`
	RegistrationDate string  `json:"registration_date" example:"2024-01-15"`
}

// ApprovalRequestDTO para aprovação de solicitações
type ApprovalRequestDTO struct {
	RequestID    uint   `json:"request_id" binding:"required" example:"1"`
	Status       string `json:"status" binding:"required,oneof=approved rejected" example:"approved"`
	ReviewNotes  string `json:"review_notes" example:"Documentação em ordem"`
}

// PendingRequestResponseDTO para resposta de solicitações pendentes
type PendingRequestResponseDTO struct {
	ID           uint                   `json:"id" example:"1"`
	User         User                   `json:"user"`
	RequestType  string                 `json:"request_type" example:"new_client"`
	RequestData  RegistrationRequestDTO `json:"request_data"`
	SubmittedAt  time.Time              `json:"submitted_at"`
	Status       string                 `json:"status" example:"pending"`
}

// UpdateUserStatusDTO para alterar status de utilizador
type UpdateUserStatusDTO struct {
	Status string `json:"status" binding:"required,oneof=approved rejected blocked" example:"approved"`
	Notes  string `json:"notes" example:"Motivo da alteração"`
}

// Dados para login
type LoginRequest struct {
    Username string `json:"username" binding:"required" example:"joao.silva"`
    Password string `json:"password" binding:"required,min=6" example:"123456"`
}

// Dados para registo simples (não usado no novo fluxo)
type RegisterRequest struct {
    Username string `json:"username" binding:"required" example:"joao.silva"`
    Email    string `json:"email" binding:"required,email" example:"joao@exemplo.com"`
    Password string `json:"password" binding:"required,min=6" example:"123456"`
    Name     string `json:"name" binding:"required,min=2" example:"João Silva"`
    Phone    string `json:"phone" binding:"required" example:"912345678"`
    NIF      string `json:"nif" binding:"required" example:"123456789"`
	Role     string `json:"role" binding:"required,oneof=client accountant admin" example:"client"`
	Status   string `json:"status" binding:"required,oneof=pending approved rejected blocked" example:"pending"`
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

// UpdateProfileDTO para atualização de perfil
type UpdateProfileDTO struct {
	Name  string `json:"name" example:"João Silva"`
	Phone string `json:"phone" example:"912345678"`
}

// UpdateCompanyDTO para atualização de empresa (campos limitados)
type UpdateCompanyDTO struct {
	TradeName  string `json:"trade_name" example:"Silva Consultoria"`
	Address    string `json:"address" example:"Rua das Flores, 123"`
	PostalCode string `json:"postal_code" example:"1000-001"`
	City       string `json:"city" example:"Lisboa"`
}
