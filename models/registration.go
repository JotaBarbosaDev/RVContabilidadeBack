package models

import (
	"time"
)

// RegistrationRequest representa o histórico de solicitações com dados completos
type RegistrationRequest struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	RequestType       string    `json:"request_type" gorm:"default:'new_client'"` // new_client, existing_client
	Status            string    `json:"status" gorm:"default:'pending'"`
	SubmittedAt       time.Time `json:"submitted_at" gorm:"autoCreateTime"`
	ReviewedAt        *time.Time `json:"reviewed_at"`
	ReviewedBy        *uint     `json:"reviewed_by"`
	ReviewNotes       string    `json:"review_notes"`
	ApprovalToken     string    `json:"approval_token" gorm:"unique"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// === DADOS DO USER (armazenados até aprovação) ===
	// Dados pessoais obrigatórios
	Username            string  `json:"username" gorm:"not null"`
	Name                string  `json:"name" gorm:"not null"`
	Email               string  `json:"email" gorm:"not null"`
	Phone               string  `json:"phone" gorm:"not null"`
	NIF                 string  `json:"nif" gorm:"not null"`
	PasswordHash        string  `json:"-" gorm:"not null"`
	
	// Dados pessoais opcionais
	DateOfBirth         *time.Time `json:"date_of_birth"`
	MaritalStatus       *string    `json:"marital_status"`
	CitizenCardNumber   *string    `json:"citizen_card_number"`
	CitizenCardExpiry   *time.Time `json:"citizen_card_expiry"`
	TaxResidenceCountry *string    `json:"tax_residence_country"`
	FixedPhone          *string    `json:"fixed_phone"`
	
	// Morada fiscal
	FiscalAddress       string  `json:"fiscal_address" gorm:"not null"`
	FiscalPostalCode    string  `json:"fiscal_postal_code" gorm:"not null"`
	FiscalCity          string  `json:"fiscal_city" gorm:"not null"`
	FiscalCounty        *string `json:"fiscal_county"`
	FiscalDistrict      *string `json:"fiscal_district"`
	
	// Preferências
	OfficialEmail         *string `json:"official_email"`
	BillingSoftware       *string `json:"billing_software"`
	PreferredFormat       *string `json:"preferred_format"`
	ReportFrequency       *string `json:"report_frequency"`
	PreferredContactHours *string `json:"preferred_contact_hours"`
	
	// === DADOS DA COMPANY (armazenados até aprovação) ===
	// Dados básicos obrigatórios
	CompanyName     string `json:"company_name" gorm:"not null"`
	NIPC            string `json:"nipc"`
	LegalForm       string `json:"legal_form" gorm:"not null"`
	
	// Dados básicos opcionais
	CAE                 *string    `json:"cae"`
	FoundingDate        *time.Time `json:"founding_date"`
	AccountingRegime    *string    `json:"accounting_regime"`
	VATRegime           *string    `json:"vat_regime"`
	BusinessActivity    *string    `json:"business_activity"`
	EstimatedRevenue    *float64   `json:"estimated_revenue"`
	MonthlyInvoices     *int       `json:"monthly_invoices"`
	NumberEmployees     *int       `json:"number_employees"`
	
	// Dados completos opcionais
	TradeName         *string    `json:"trade_name"`
	CorporateObject   *string    `json:"corporate_object"`
	CompanyAddress    *string    `json:"company_address"`
	CompanyPostalCode *string    `json:"company_postal_code"`
	CompanyCity       *string    `json:"company_city"`
	CompanyCounty     *string    `json:"company_county"`
	CompanyDistrict   *string    `json:"company_district"`
	CompanyCountry    *string    `json:"company_country"`
	ShareCapital      *float64   `json:"share_capital"`
	GroupStartDate    *time.Time `json:"group_start_date"`
	
	// Informação bancária
	BankName *string `json:"bank_name"`
	IBAN     *string `json:"iban"`
	BIC      *string `json:"bic"`
	
	// Dados operacionais
	AnnualRevenue *float64 `json:"annual_revenue"`
	HasStock      *bool    `json:"has_stock"`
	MainClients   *string  `json:"main_clients"`
	MainSuppliers *string  `json:"main_suppliers"`
	
	// === RELACIONAMENTOS (só existem após aprovação) ===
	UserID    *uint `json:"user_id,omitempty" gorm:"index"`    // NULL até aprovação
	CompanyID *uint `json:"company_id,omitempty" gorm:"index"` // NULL até aprovação
	
	// Relacionamentos
	User           *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Company        *Company `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
	ReviewedByUser *User `json:"reviewed_by_user,omitempty" gorm:"foreignKey:ReviewedBy"`
}

// RegistrationRequestDTO para registo completo
type RegistrationRequestDTO struct {
	// === DADOS PESSOAIS ===
	// Obrigatórios
	Username    string `json:"username" binding:"required" example:"joao.silva"`
	Name        string `json:"name" binding:"required" example:"João Silva"`
	Email       string `json:"email" binding:"required,email" example:"joao@exemplo.com"`
	Phone       string `json:"phone" binding:"required" example:"912345678"`
	NIF         string `json:"nif" binding:"required" example:"123456789"`
	Password    string `json:"password" binding:"required,min=6" example:"password123"`
	
	// Morada fiscal obrigatória
	FiscalAddress    string `json:"fiscal_address" binding:"required" example:"Rua das Flores, 123"`
	FiscalPostalCode string `json:"fiscal_postal_code" binding:"required" example:"1000-001"`
	FiscalCity       string `json:"fiscal_city" binding:"required" example:"Lisboa"`
	
	// === DADOS EMPRESA ===
	// Obrigatórios
	CompanyName string `json:"company_name" binding:"required" example:"Silva & Associados Lda"`
	LegalForm   string `json:"legal_form" binding:"required" example:"Sociedade por Quotas"`
	
	// Opcionais
	NIPC        string `json:"nipc,omitempty" validate:"omitempty" example:"123456789"`
	
	// === TODOS OS CAMPOS OPCIONAIS ===
	DateOfBirth           *string  `json:"date_of_birth,omitempty" example:"1990-01-15"`
	MaritalStatus         *string  `json:"marital_status,omitempty" example:"Solteiro"`
	CitizenCardNumber     *string  `json:"citizen_card_number,omitempty" example:"12345678"`
	CitizenCardExpiry     *string  `json:"citizen_card_expiry,omitempty" example:"2030-12-31"`
	TaxResidenceCountry   *string  `json:"tax_residence_country,omitempty" example:"Portugal"`
	FixedPhone            *string  `json:"fixed_phone,omitempty" example:"213456789"`
	FiscalCounty          *string  `json:"fiscal_county,omitempty" example:"Lisboa"`
	FiscalDistrict        *string  `json:"fiscal_district,omitempty" example:"Lisboa"`
	OfficialEmail         *string  `json:"official_email,omitempty" example:"geral@empresa.com"`
	BillingSoftware       *string  `json:"billing_software,omitempty" example:"Moloni"`
	PreferredFormat       *string  `json:"preferred_format,omitempty" example:"digital"`
	ReportFrequency       *string  `json:"report_frequency,omitempty" example:"mensal"`
	PreferredContactHours *string  `json:"preferred_contact_hours,omitempty" example:"9h-17h"`
	
	CAE                *string  `json:"cae,omitempty" example:"69200"`
	FoundingDate       *string  `json:"founding_date,omitempty" example:"2024-01-15"`
	AccountingRegime   *string  `json:"accounting_regime,omitempty" example:"organizada"`
	VATRegime          *string  `json:"vat_regime,omitempty" example:"normal"`
	BusinessActivity   *string  `json:"business_activity,omitempty" example:"Consultoria em gestão"`
	EstimatedRevenue   *float64 `json:"estimated_revenue,omitempty" example:"50000.00"`
	MonthlyInvoices    *int     `json:"monthly_invoices,omitempty" example:"10"`
	NumberEmployees    *int     `json:"number_employees,omitempty" example:"2"`
	TradeName          *string  `json:"trade_name,omitempty" example:"Silva Consultoria"`
	CorporateObject    *string  `json:"corporate_object,omitempty" example:"Prestação de serviços de consultoria"`
	CompanyAddress     *string  `json:"company_address,omitempty" example:"Rua das Flores, 123"`
	CompanyPostalCode  *string  `json:"company_postal_code,omitempty" example:"1000-001"`
	CompanyCity        *string  `json:"company_city,omitempty" example:"Lisboa"`
	CompanyCounty      *string  `json:"company_county,omitempty" example:"Lisboa"`
	CompanyDistrict    *string  `json:"company_district,omitempty" example:"Lisboa"`
	CompanyCountry     *string  `json:"company_country,omitempty" example:"Portugal"`
	ShareCapital       *float64 `json:"share_capital,omitempty" example:"5000.00"`
	GroupStartDate     *string  `json:"group_start_date,omitempty" example:"2024-01-01"`
	BankName           *string  `json:"bank_name,omitempty" example:"Banco Comercial Português"`
	IBAN              *string  `json:"iban,omitempty" example:"PT50000201231234567890154"`
	BIC               *string  `json:"bic,omitempty" example:"BCOMPTPL"`
	AnnualRevenue     *float64 `json:"annual_revenue,omitempty" example:"100000.00"`
	HasStock          *bool    `json:"has_stock,omitempty" example:"false"`
	MainClients       *string  `json:"main_clients,omitempty" example:"Cliente A, Cliente B"`
	MainSuppliers     *string  `json:"main_suppliers,omitempty" example:"Fornecedor X, Fornecedor Y"`
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
	RequestType  string                 `json:"request_type" example:"new_client"`
	Status       string                 `json:"status" example:"pending"`
	SubmittedAt  time.Time              `json:"submitted_at"`
	
	// Dados do utilizador da solicitação
	Username     string `json:"username" example:"joao.silva"`
	Name         string `json:"name" example:"João Silva"`
	Email        string `json:"email" example:"joao@exemplo.com"`
	Phone        string `json:"phone" example:"912345678"`
	NIF          string `json:"nif" example:"123456789"`
	CompanyName  string `json:"company_name" example:"Silva & Associados Lda"`
	NIPC         string `json:"nipc" example:"123456789"`
	
	// Dados completos se necessário
	RequestData RegistrationRequestDTO `json:"request_data,omitempty"`
}
