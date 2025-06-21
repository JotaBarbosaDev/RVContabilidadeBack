package models

import (
	"time"
)

// User representa um utilizador do sistema
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
    Username  string    `json:"username" gorm:"unique;not null" example:"joao.silva"`
    Email     string    `json:"email" gorm:"unique;not null" example:"joao@exemplo.com"`
    Password  string    `json:"-" gorm:"not null"`
    Name      string    `json:"name" gorm:"not null" example:"João Silva"` 
	Phone     string    `json:"phone" gorm:"not null" example:"912345678"`
	NIF       string    `json:"nif" gorm:"unique;not null" example:"123456789"`
    Role      string    `json:"role" gorm:"default:'client'" example:"client"`
	Status    string    `json:"status" gorm:"default:'pending'" example:"pending"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime" example:"2023-01-01T00:00:00Z"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2023-01-01T00:00:00Z"`

	// Dados pessoais completos (preenchidos depois da aprovação)
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
	PortalFinancasUser     string `json:"portal_financas_user"`
	PortalFinancasPassword string `json:"portal_financas_password"`
	EFaturaUser           string `json:"e_fatura_user"`
	EFaturaPassword       string `json:"e_fatura_password"`
	SSDirectUser          string `json:"ss_direct_user"`
	SSDirectPassword      string `json:"ss_direct_password"`
	OfficialEmail         string `json:"official_email"`
	BillingSoftware       string `json:"billing_software"`
	
	// Preferências de comunicação
	PreferredFormat       string `json:"preferred_format" gorm:"default:'digital'"`
	ReportFrequency       string `json:"report_frequency" gorm:"default:'mensal'"`
	PreferredContactHours string `json:"preferred_contact_hours"`
	
	// Relacionamentos
	Company  *Company              `json:"company,omitempty" gorm:"foreignKey:UserID"`
	Requests []RegistrationRequest `json:"requests,omitempty" gorm:"foreignKey:UserID"`
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
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Dados básicos (obrigatórios no registo)
	CompanyName  string     `json:"company_name" gorm:"not null"`
	NIPC         string     `json:"nipc" gorm:"unique;not null"`
	CAE          string     `json:"cae" gorm:"not null"`
	LegalForm    string     `json:"legal_form" gorm:"not null"`
	FoundingDate *time.Time `json:"founding_date"`
	
	// Regime contabilístico e fiscal (obrigatórios)
	AccountingRegime string `json:"accounting_regime" gorm:"not null"`
	VATRegime       string `json:"vat_regime" gorm:"not null"`
	
	// Dados operacionais básicos
	BusinessActivity string  `json:"business_activity" gorm:"not null"`
	EstimatedRevenue float64 `json:"estimated_revenue"`
	MonthlyInvoices  int     `json:"monthly_invoices"`
	NumberEmployees  int     `json:"number_employees"`
	
	// Dados completos (preenchidos depois da aprovação)
	TradeName       string  `json:"trade_name"`
	CorporateObject string  `json:"corporate_object"`
	Address         string  `json:"address"`
	PostalCode      string  `json:"postal_code"`
	City            string  `json:"city"`
	County          string  `json:"county"`
	District        string  `json:"district"`
	Country         string  `json:"country" gorm:"default:'Portugal'"`
	ShareCapital    float64 `json:"share_capital"`
	GroupStartDate  *time.Time `json:"group_start_date"`
	
	// Informação bancária
	BankName string `json:"bank_name"`
	IBAN     string `json:"iban"`
	BIC      string `json:"bic"`
	
	// Dados operacionais completos
	AnnualRevenue float64 `json:"annual_revenue"`
	HasStock      bool    `json:"has_stock"`
	MainClients   string  `json:"main_clients"`
	MainSuppliers string  `json:"main_suppliers"`
	Status        string  `json:"status" gorm:"default:'active'"`
	
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

// RegistrationRequestDTO para registo de clientes (dados mínimos)
type RegistrationRequestDTO struct {
	// Dados pessoais mínimos
	Username    string `json:"username" binding:"required" example:"joao.silva"`
	Name        string `json:"name" binding:"required" example:"João Silva"`
	Email       string `json:"email" binding:"required,email" example:"joao@exemplo.com"`
	Phone       string `json:"phone" binding:"required" example:"912345678"`
	NIF         string `json:"nif" binding:"required" example:"123456789"`
	Password    string `json:"password" binding:"required,min=6" example:"password123"`
	DateOfBirth string `json:"date_of_birth" binding:"required" example:"1990-01-15"`
	
	// Morada fiscal mínima
	FiscalAddress    string `json:"fiscal_address" binding:"required" example:"Rua das Flores, 123"`
	FiscalPostalCode string `json:"fiscal_postal_code" binding:"required" example:"1000-001"`
	FiscalCity       string `json:"fiscal_city" binding:"required" example:"Lisboa"`
	
	// Dados empresa mínimos
	CompanyName   string `json:"company_name" binding:"required" example:"Silva & Associados Lda"`
	NIPC          string `json:"nipc" binding:"required" example:"123456789"`
	CAE           string `json:"cae" binding:"required" example:"69200"`
	LegalForm     string `json:"legal_form" binding:"required" example:"Sociedade por Quotas"`
	FoundingDate  string `json:"founding_date" binding:"required" example:"2024-01-15"`
	
	// Regimes (obrigatórios para aprovação)
	AccountingRegime string `json:"accounting_regime" binding:"required,oneof=organizada simplificada" example:"organizada"`
	VATRegime        string `json:"vat_regime" binding:"required,oneof=normal isento_art53 pequeno_retalhista" example:"normal"`
	
	// Dados operacionais básicos
	BusinessActivity  string  `json:"business_activity" binding:"required" example:"Consultoria em gestão"`
	EstimatedRevenue  float64 `json:"estimated_revenue" binding:"required" example:"50000.00"`
	MonthlyInvoices   int     `json:"monthly_invoices" binding:"required" example:"10"`
	NumberEmployees   int     `json:"number_employees" example:"2"`
	ReportFrequency   string  `json:"report_frequency" binding:"required,oneof=mensal trimestral" example:"mensal"`
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

// CompleteCompanyDataDTO para completar dados da empresa após aprovação
type CompleteCompanyDataDTO struct {
	TradeName       string  `json:"trade_name" example:"Silva Consultoria"`
	CorporateObject string  `json:"corporate_object" example:"Prestação de serviços de consultoria"`
	Address         string  `json:"address" example:"Rua das Flores, 123"`
	PostalCode      string  `json:"postal_code" example:"1000-001"`
	City            string  `json:"city" example:"Lisboa"`
	County          string  `json:"county" example:"Lisboa"`
	District        string  `json:"district" example:"Lisboa"`
	ShareCapital    float64 `json:"share_capital" example:"5000.00"`
	GroupStartDate  string  `json:"group_start_date" example:"2024-01-01"`
	BankName        string  `json:"bank_name" example:"Banco Comercial Português"`
	IBAN           string  `json:"iban" example:"PT50000201231234567890154"`
	BIC            string  `json:"bic" example:"BCOMPTPL"`
	AnnualRevenue  float64 `json:"annual_revenue" example:"100000.00"`
	HasStock       bool    `json:"has_stock" example:"false"`
	MainClients    string  `json:"main_clients" example:"Cliente A, Cliente B"`
	MainSuppliers  string  `json:"main_suppliers" example:"Fornecedor X, Fornecedor Y"`
}
