package models

import (
	"encoding/json"
	"strconv"
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
	// Dados pessoais obrigatórios apenas username e password
	Username            string  `json:"username" gorm:"not null"`
	Name                *string `json:"name"`
	Email               *string `json:"email"`
	Phone               *string `json:"phone"`
	NIF                 *string `json:"nif"`
	PasswordHash        string  `json:"-" gorm:"not null"`
	
	// Dados pessoais opcionais
	DateOfBirth         *time.Time `json:"date_of_birth"`
	MaritalStatus       *string    `json:"marital_status"`
	CitizenCardNumber   *string    `json:"citizen_card_number"`
	CitizenCardExpiry   *time.Time `json:"citizen_card_expiry"`
	TaxResidenceCountry *string    `json:"tax_residence_country"`
	FixedPhone          *string    `json:"fixed_phone"`
	
	// Morada fiscal opcional
	FiscalAddress       *string `json:"fiscal_address"`
	FiscalPostalCode    *string `json:"fiscal_postal_code"`
	FiscalCity          *string `json:"fiscal_city"`
	FiscalCounty        *string `json:"fiscal_county"`
	FiscalDistrict      *string `json:"fiscal_district"`
	
	// Morada da empresa (campos adicionais do frontend)
	Address     *string `json:"address"`
	PostalCode  *string `json:"postal_code"`
	City        *string `json:"city"`
	Country     *string `json:"country"`
	
	// Preferências
	OfficialEmail         *string `json:"official_email"`
	BillingSoftware       *string `json:"billing_software"`
	PreferredFormat       *string `json:"preferred_format"`
	ReportFrequency       *string `json:"report_frequency"`
	PreferredContactHours *string `json:"preferred_contact_hours"`
	
	// === DADOS DA COMPANY (armazenados até aprovação) ===
	// Dados básicos opcionais (apenas LegalForm obrigatório)
	CompanyName     *string `json:"company_name"`
	NIPC            string  `json:"nipc"`
	LegalForm       string  `json:"legal_form" gorm:"not null"`
	
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
	// Obrigatórios apenas username e password
	Username    string `json:"username" binding:"required" example:"joao.silva"`
	Name        string `json:"name,omitempty" example:"João Silva"`
	Email       string `json:"email,omitempty" example:"joao@exemplo.com"`
	Phone       string `json:"phone,omitempty" example:"912345678"`
	NIF         string `json:"nif,omitempty" example:"123456789"`
	Password    string `json:"password" binding:"required,min=6" example:"password123"`
	
	// Morada fiscal opcional
	FiscalAddress    string `json:"fiscal_address,omitempty" example:"Rua das Flores, 123"`
	FiscalPostalCode string `json:"fiscal_postal_code,omitempty" example:"1000-001"`
	FiscalCity       string `json:"fiscal_city,omitempty" example:"Lisboa"`
	
	// Morada da empresa (campos adicionais do frontend)
	Address     string `json:"address,omitempty" example:"Rua da Empresa, 456"`
	PostalCode  string `json:"postal_code,omitempty" example:"1000-002"`
	City        string `json:"city,omitempty" example:"Porto"`
	Country     string `json:"country,omitempty" example:"Portugal"`
	
	// === DADOS EMPRESA ===
	// Obrigatórios
	LegalForm   string `json:"legal_form" binding:"required" example:"Sociedade por Quotas"`
	
	// Opcionais principais (campos que o frontend envia)
	CompanyName   string                `json:"company_name,omitempty" example:"Silva & Associados Lda"`
	TradeName     string                `json:"trade_name,omitempty" example:"Silva Consultoria"`
	NIPC          string                `json:"nipc,omitempty" validate:"omitempty" example:"123456789"`
	CAE           string                `json:"cae,omitempty" example:"69200"`
	FoundingDate  string                `json:"founding_date,omitempty" example:"2024-01-15"`
	ShareCapital  *FlexibleFloat64      `json:"share_capital,omitempty" swaggertype:"number" example:"5000.00"`
	
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
	
	// Outros campos da empresa (sem duplicatas)
	AccountingRegime   *string         `json:"accounting_regime,omitempty" example:"organizada"`
	VATRegime          *string         `json:"vat_regime,omitempty" example:"normal"`
	BusinessActivity   *string         `json:"business_activity,omitempty" example:"Consultoria em gestão"`
	EstimatedRevenue   *FlexibleFloat64 `json:"estimated_revenue,omitempty" swaggertype:"number" example:"50000.00"`
	MonthlyInvoices    *FlexibleInt     `json:"monthly_invoices,omitempty" swaggertype:"integer" example:"10"`
	NumberEmployees    *FlexibleInt     `json:"number_employees,omitempty" swaggertype:"integer" example:"2"`
	CorporateObject    *string         `json:"corporate_object,omitempty" example:"Prestação de serviços de consultoria"`
	CompanyAddress     *string         `json:"company_address,omitempty" example:"Rua das Flores, 123"`
	CompanyPostalCode  *string         `json:"company_postal_code,omitempty" example:"1000-001"`
	CompanyCity        *string         `json:"company_city,omitempty" example:"Lisboa"`
	CompanyCounty      *string         `json:"company_county,omitempty" example:"Lisboa"`
	CompanyDistrict    *string         `json:"company_district,omitempty" example:"Lisboa"`
	CompanyCountry     *string         `json:"company_country,omitempty" example:"Portugal"`
	GroupStartDate     *string         `json:"group_start_date,omitempty" example:"2024-01-01"`
	BankName           *string         `json:"bank_name,omitempty" example:"Banco Comercial Português"`
	IBAN              *string         `json:"iban,omitempty" example:"PT50000201231234567890154"`
	BIC               *string         `json:"bic,omitempty" example:"BCOMPTPL"`
	AnnualRevenue     *FlexibleFloat64 `json:"annual_revenue,omitempty" swaggertype:"number" example:"100000.00"`
	HasStock          *bool           `json:"has_stock,omitempty" example:"false"`
	MainClients       *string         `json:"main_clients,omitempty" example:"Cliente A, Cliente B"`
	MainSuppliers     *string         `json:"main_suppliers,omitempty" example:"Fornecedor X, Fornecedor Y"`
	
	// Campo adicional enviado pelo frontend
	RegistrationDate  *string         `json:"registration_date,omitempty" example:"2024-01-01"`
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
	
	// Dados essenciais do utilizador
	Username     string `json:"username" example:"joao.silva"`
	Name         string `json:"name" example:"João Silva"`
	Email        string `json:"email" example:"joao@exemplo.com"`
	Phone        string `json:"phone" example:"912345678"`
	NIF          string `json:"nif" example:"123456789"`
	
	// Dados da empresa
	CompanyName  string `json:"company_name" example:"Silva & Associados Lda"`
	NIPC         string `json:"nipc" example:"123456789"`
	LegalForm    string `json:"legal_form" example:"Sociedade por Quotas"`
	
	// Morada fiscal
	FiscalAddress    string `json:"fiscal_address" example:"Rua das Flores, 123"`
	FiscalPostalCode string `json:"fiscal_postal_code" example:"1000-001"`
	FiscalCity       string `json:"fiscal_city" example:"Lisboa"`
	
	// Dados completos se necessário (para visualização detalhada)
	RequestData RegistrationRequestDTO `json:"request_data,omitempty"`
}

// DashboardDataDTO para resposta do dashboard da contabilista
type DashboardDataDTO struct {
	// Estatísticas gerais
	TotalPendingRequests int `json:"total_pending_requests" example:"5"`
	TotalApprovedClients int `json:"total_approved_clients" example:"25"`
	TotalRejectedRequests int `json:"total_rejected_requests" example:"2"`
	
	// Solicitações pendentes recentes (últimas 10)
	RecentPendingRequests []PendingRequestResponseDTO `json:"recent_pending_requests"`
	
	// Estatísticas mensais
	MonthlyStats struct {
		NewRequests     int `json:"new_requests" example:"8"`
		ApprovedClients int `json:"approved_clients" example:"6"`
		RejectedRequests int `json:"rejected_requests" example:"1"`
	} `json:"monthly_stats"`
}

// ClientsOverviewDTO para visão geral de todos os clientes
type ClientsOverviewDTO struct {
	// Clientes pendentes
	PendingClients []PendingRequestResponseDTO `json:"pending_clients"`
	
	// Clientes aprovados (resumo)
	ApprovedClients []struct {
		ID          uint   `json:"id" example:"1"`
		Username    string `json:"username" example:"joao.silva"`
		Name        string `json:"name" example:"João Silva"`
		Email       string `json:"email" example:"joao@exemplo.com"`
		Phone       string `json:"phone" example:"912345678"`
		NIF         string `json:"nif" example:"123456789"`
		Status      string `json:"status" example:"approved"`
		CompanyName string `json:"company_name" example:"Silva & Associados Lda"`
		NIPC        string `json:"nipc" example:"123456789"`
		CreatedAt   string `json:"created_at" example:"2024-01-01T00:00:00Z"`
	} `json:"approved_clients"`
	
	// Estatísticas
	Stats struct {
		TotalPending  int `json:"total_pending"`
		TotalApproved int `json:"total_approved"`
		TotalRejected int `json:"total_rejected"`
	} `json:"stats"`
}

// === TIPOS PERSONALIZADOS PARA CONVERSÃO ===

// FlexibleFloat64 aceita tanto string quanto float64 no JSON
type FlexibleFloat64 struct {
	Value *float64
}

// UnmarshalJSON implementa conversão personalizada
func (f *FlexibleFloat64) UnmarshalJSON(data []byte) error {
	// Tentar primeiro como número
	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		f.Value = &num
		return nil
	}
	
	// Tentar como string
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if str == "" {
			f.Value = nil
			return nil
		}
		if parsed, err := strconv.ParseFloat(str, 64); err == nil {
			f.Value = &parsed
			return nil
		}
	}
	
	// Se não conseguir fazer parsing, deixar como nil
	f.Value = nil
	return nil
}

// MarshalJSON serializa de volta para JSON
func (f FlexibleFloat64) MarshalJSON() ([]byte, error) {
	if f.Value == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(*f.Value)
}

// FlexibleInt aceita tanto string quanto int no JSON
type FlexibleInt struct {
	Value *int
}

// UnmarshalJSON implementa conversão personalizada para int
func (f *FlexibleInt) UnmarshalJSON(data []byte) error {
	// Tentar primeiro como número
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		f.Value = &num
		return nil
	}
	
	// Tentar como string
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if str == "" {
			f.Value = nil
			return nil
		}
		if parsed, err := strconv.Atoi(str); err == nil {
			f.Value = &parsed
			return nil
		}
	}
	
	// Se não conseguir fazer parsing, deixar como nil
	f.Value = nil
	return nil
}

// MarshalJSON serializa de volta para JSON
func (f FlexibleInt) MarshalJSON() ([]byte, error) {
	if f.Value == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(*f.Value)
}

// CompleteUserOverviewDTO para endpoint que combina dados das 3 tabelas
type CompleteUserOverviewDTO struct {
	// === IDENTIFICAÇÃO ÚNICA ===
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"joao.silva"`
	
	// === STATUS E TIPO ===
	Status string `json:"status" example:"approved"` // pending, approved, rejected, blocked
	Role   string `json:"role" example:"client"`     // client, accountant, admin
	Source string `json:"source" example:"both"`     // registration_request, user_only, both
	
	// === DADOS PESSOAIS (prioridade: User > RegistrationRequest) ===
	Name                *string    `json:"name" example:"João Silva"`
	Email               *string    `json:"email" example:"joao@exemplo.com"`
	Phone               *string    `json:"phone" example:"912345678"`
	NIF                 *string    `json:"nif" example:"123456789"`
	DateOfBirth         *time.Time `json:"date_of_birth"`
	MaritalStatus       *string    `json:"marital_status" example:"Solteiro"`
	CitizenCardNumber   *string    `json:"citizen_card_number" example:"12345678"`
	CitizenCardExpiry   *time.Time `json:"citizen_card_expiry"`
	TaxResidenceCountry *string    `json:"tax_residence_country" example:"Portugal"`
	FixedPhone          *string    `json:"fixed_phone" example:"213456789"`
	
	// === MORADA FISCAL ===
	FiscalAddress    *string `json:"fiscal_address" example:"Rua das Flores, 123"`
	FiscalPostalCode *string `json:"fiscal_postal_code" example:"1000-001"`
	FiscalCity       *string `json:"fiscal_city" example:"Lisboa"`
	FiscalCounty     *string `json:"fiscal_county" example:"Lisboa"`
	FiscalDistrict   *string `json:"fiscal_district" example:"Lisboa"`
	
	// === MORADA PESSOAL/EMPRESA (de registration_requests) ===
	Address    *string `json:"address" example:"Rua da Empresa, 456"`
	PostalCode *string `json:"postal_code" example:"1000-002"`
	City       *string `json:"city" example:"Porto"`
	Country    *string `json:"country" example:"Portugal"`
	
	// === PREFERÊNCIAS ===
	OfficialEmail         *string `json:"official_email" example:"geral@empresa.com"`
	BillingSoftware       *string `json:"billing_software" example:"Moloni"`
	PreferredFormat       *string `json:"preferred_format" example:"digital"`
	ReportFrequency       *string `json:"report_frequency" example:"mensal"`
	PreferredContactHours *string `json:"preferred_contact_hours" example:"9h-17h"`
	
	// === DADOS DA EMPRESA (prioridade: Company > RegistrationRequest) ===
	CompanyID       *uint      `json:"company_id,omitempty"`
	CompanyName     *string    `json:"company_name" example:"Silva & Associados Lda"`
	TradeName       *string    `json:"trade_name" example:"Silva Consultoria"`
	NIPC            *string    `json:"nipc" example:"123456789"`
	LegalForm       *string    `json:"legal_form" example:"Sociedade por Quotas"`
	CAE             *string    `json:"cae" example:"69200"`
	FoundingDate    *time.Time `json:"founding_date"`
	ShareCapital    *float64   `json:"share_capital" example:"5000.00"`
	CompanyStatus   *string    `json:"company_status" example:"active"` // Campo exclusivo de Company
	
	// === CONFIGURAÇÕES CONTABILÍSTICAS ===
	AccountingRegime *string  `json:"accounting_regime" example:"organizada"`
	VATRegime        *string  `json:"vat_regime" example:"normal"`
	BusinessActivity *string  `json:"business_activity" example:"Consultoria empresarial"`
	EstimatedRevenue *float64 `json:"estimated_revenue" example:"50000.00"`
	MonthlyInvoices  *int     `json:"monthly_invoices" example:"10"`
	NumberEmployees  *int     `json:"number_employees" example:"2"`
	
	// === DETALHES DA EMPRESA ===
	CorporateObject *string `json:"corporate_object" example:"Prestação de serviços"`
	
	// === MORADA DA EMPRESA (Company table - pode ser diferente da morada pessoal) ===
	CompanyAddress    *string `json:"company_address" example:"Rua das Flores, 123"`
	CompanyPostalCode *string `json:"company_postal_code" example:"1000-001"`
	CompanyCity       *string `json:"company_city" example:"Lisboa"`
	CompanyCounty     *string `json:"company_county" example:"Lisboa"`
	CompanyDistrict   *string `json:"company_district" example:"Lisboa"`
	CompanyCountry    *string `json:"company_country" example:"Portugal"`
	GroupStartDate    *time.Time `json:"group_start_date"`
	
	// === INFORMAÇÃO BANCÁRIA ===
	BankName *string `json:"bank_name" example:"Banco Comercial Português"`
	IBAN     *string `json:"iban" example:"PT50000201231234567890154"`
	BIC      *string `json:"bic" example:"BCOMPTPL"`
	
	// === DADOS OPERACIONAIS ===
	AnnualRevenue *float64 `json:"annual_revenue" example:"100000.00"`
	HasStock      *bool    `json:"has_stock" example:"false"`
	MainClients   *string  `json:"main_clients" example:"Cliente A, Cliente B"`
	MainSuppliers *string  `json:"main_suppliers" example:"Fornecedor X, Fornecedor Y"`
	
	// === DADOS ESPECÍFICOS DE REGISTRATION_REQUEST ===
	RequestID       *uint      `json:"request_id,omitempty"`
	RequestType     *string    `json:"request_type" example:"new_client"`
	RequestStatus   *string    `json:"request_status" example:"approved"` // Status da request
	SubmittedAt     *time.Time `json:"submitted_at"`
	ReviewedAt      *time.Time `json:"reviewed_at"`
	ReviewedBy      *uint      `json:"reviewed_by"`
	ReviewNotes     *string    `json:"review_notes"`
	ApprovalToken   *string    `json:"approval_token,omitempty"`
	ReviewedByName  *string    `json:"reviewed_by_name,omitempty"`
	
	// === TIMESTAMPS ===
	UserCreatedAt    *time.Time `json:"user_created_at,omitempty"`    // De users
	UserUpdatedAt    *time.Time `json:"user_updated_at,omitempty"`    // De users
	CompanyCreatedAt *time.Time `json:"company_created_at,omitempty"` // De companies
	CompanyUpdatedAt *time.Time `json:"company_updated_at,omitempty"` // De companies
	RequestCreatedAt *time.Time `json:"request_created_at,omitempty"` // De registration_requests
	RequestUpdatedAt *time.Time `json:"request_updated_at,omitempty"` // De registration_requests
}
