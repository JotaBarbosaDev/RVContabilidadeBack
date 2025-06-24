package models

import (
	"time"
)

// User representa um utilizador aprovado do sistema
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
    Username  string    `json:"username" gorm:"unique;not null" example:"joao.silva"`
    Email     string    `json:"email" gorm:"unique;not null" example:"joao@exemplo.com"`
    Password  string    `json:"-" gorm:"not null"`
    Name      string    `json:"name" gorm:"not null" example:"João Silva"` 
	Phone     string    `json:"phone" gorm:"not null" example:"912345678"`
	NIF       string    `json:"nif" gorm:"unique;not null" example:"123456789"`
    Role      string    `json:"role" gorm:"default:'client'" example:"client"`
	Status    string    `json:"status" gorm:"default:'approved'" example:"approved"` // Sempre approved quando criado
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime" example:"2023-01-01T00:00:00Z"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2023-01-01T00:00:00Z"`

	// Dados copiados da RegistrationRequest após aprovação
	DateOfBirth         *time.Time `json:"date_of_birth"`
	MaritalStatus       string     `json:"marital_status"`
	CitizenCardNumber   string     `json:"citizen_card_number"`
	CitizenCardExpiry   *time.Time `json:"citizen_card_expiry"`
	TaxResidenceCountry string     `json:"tax_residence_country" gorm:"default:'Portugal'"`
	
	// Contactos adicionais
	FixedPhone string `json:"fixed_phone"`
	
	// Morada fiscal completa
	FiscalAddress    string `json:"fiscal_address"`
	FiscalPostalCode string `json:"fiscal_postal_code"`
	FiscalCity       string `json:"fiscal_city"`
	FiscalCounty     string `json:"fiscal_county"`
	FiscalDistrict   string `json:"fiscal_district"`
	
	// Acessos e credenciais (guardados encriptados)
	OfficialEmail         string `json:"official_email"`
	BillingSoftware       string `json:"billing_software"`
	
	// Preferências de comunicação
	PreferredFormat       string `json:"preferred_format" gorm:"default:'digital'"`
	ReportFrequency       string `json:"report_frequency" gorm:"default:'mensal'"`
	PreferredContactHours string `json:"preferred_contact_hours"`
	
	// Relacionamentos
	Company             *Company              `json:"company,omitempty" gorm:"foreignKey:UserID"`
	RegistrationRequest *RegistrationRequest  `json:"registration_request,omitempty" gorm:"foreignKey:UserID"`
	ReviewedRequests    []RegistrationRequest `json:"reviewed_requests,omitempty" gorm:"foreignKey:ReviewedBy"`
}

// UserStatus constants
type UserStatus string

const (
	StatusPending  UserStatus = "pending"   // Aguarda aprovação
	StatusApproved UserStatus = "approved"  // Aprovado, pode aceder
	StatusRejected UserStatus = "rejected"  // Rejeitado, sem acesso
	StatusBlocked  UserStatus = "blocked"   // Bloqueado
)

// DTOs para Request/Response - User específicos

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

// CompleteUserDataDTO para completar dados pessoais após aprovação
type CompleteUserDataDTO struct {
	MaritalStatus           string `json:"marital_status" example:"Solteiro"`
	CitizenCardNumber       string `json:"citizen_card_number" example:"12345678"`
	CitizenCardExpiry       string `json:"citizen_card_expiry" example:"2030-12-31"`
	FixedPhone             string `json:"fixed_phone" example:"213456789"`
	FiscalCounty           string `json:"fiscal_county" example:"Lisboa"`
	FiscalDistrict         string `json:"fiscal_district" example:"Lisboa"`
	PortalFinancasUser     string `json:"portal_financas_user" example:"123456789"`
	PortalFinancasPassword string `json:"portal_financas_password" example:"password123"`
	EFaturaUser           string `json:"e_fatura_user" example:"123456789"`
	EFaturaPassword       string `json:"e_fatura_password" example:"password123"`
	SSDirectUser          string `json:"ss_direct_user" example:"123456789"`
	SSDirectPassword      string `json:"ss_direct_password" example:"password123"`
	OfficialEmail         string `json:"official_email" example:"geral@empresa.com"`
	BillingSoftware       string `json:"billing_software" example:"Moloni"`
	PreferredFormat       string `json:"preferred_format" example:"digital"`
	PreferredContactHours string `json:"preferred_contact_hours" example:"9h-17h"`
}

// AdminUpdateClientDTO para contabilistas/admins editarem dados de clientes
type AdminUpdateClientDTO struct {
	// Dados pessoais que admin pode editar
	Name                  *string `json:"name,omitempty"`
	Email                 *string `json:"email,omitempty"`
	Phone                 *string `json:"phone,omitempty"`
	NIF                   *string `json:"nif,omitempty"`
	DateOfBirth          *string `json:"date_of_birth,omitempty"`
	MaritalStatus        *string `json:"marital_status,omitempty"`
	CitizenCardNumber    *string `json:"citizen_card_number,omitempty"`
	CitizenCardExpiry    *string `json:"citizen_card_expiry,omitempty"`
	FixedPhone           *string `json:"fixed_phone,omitempty"`
	
	// Morada fiscal
	FiscalAddress        *string `json:"fiscal_address,omitempty"`
	FiscalPostalCode     *string `json:"fiscal_postal_code,omitempty"`
	FiscalCity           *string `json:"fiscal_city,omitempty"`
	FiscalCounty         *string `json:"fiscal_county,omitempty"`
	FiscalDistrict       *string `json:"fiscal_district,omitempty"`
	
	// Configurações
	OfficialEmail        *string `json:"official_email,omitempty"`
	BillingSoftware      *string `json:"billing_software,omitempty"`
	PreferredFormat      *string `json:"preferred_format,omitempty"`
	PreferredContactHours *string `json:"preferred_contact_hours,omitempty"`
	Status               *string `json:"status,omitempty"`
}
