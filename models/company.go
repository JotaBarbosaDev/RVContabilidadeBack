package models

import (
	"time"
)

// Company representa uma empresa aprovada
type Company struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;unique"` // One-to-One com User
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Dados copiados da RegistrationRequest após aprovação
	CompanyName      string     `json:"company_name" gorm:"not null"`
	NIPC             string     `json:"nipc" gorm:"unique"`
	CAE              string     `json:"cae"`
	LegalForm        string     `json:"legal_form" gorm:"not null"`
	FoundingDate     *time.Time `json:"founding_date"`
	AccountingRegime string     `json:"accounting_regime"`
	VATRegime        string     `json:"vat_regime"`
	BusinessActivity string     `json:"business_activity"`
	EstimatedRevenue float64    `json:"estimated_revenue"`
	MonthlyInvoices  int        `json:"monthly_invoices"`
	NumberEmployees  int        `json:"number_employees"`
	
	// Dados completos
	TradeName       string     `json:"trade_name"`
	CorporateObject string     `json:"corporate_object"`
	Address         string     `json:"address"`
	PostalCode      string     `json:"postal_code"`
	City            string     `json:"city"`
	County          string     `json:"county"`
	District        string     `json:"district"`
	Country         string     `json:"country" gorm:"default:'Portugal'"`
	ShareCapital    float64    `json:"share_capital"`
	GroupStartDate  *time.Time `json:"group_start_date"`
	
	// Informação bancária
	BankName string `json:"bank_name"`
	IBAN     string `json:"iban"`
	BIC      string `json:"bic"`
	
	// Dados operacionais
	AnnualRevenue float64 `json:"annual_revenue"`
	HasStock      bool    `json:"has_stock"`
	MainClients   string  `json:"main_clients"`
	MainSuppliers string  `json:"main_suppliers"`
	Status        string  `json:"status" gorm:"default:'active'"`
	
	// Relacionamentos
	User                *User                `json:"user,omitempty" gorm:"foreignKey:UserID"`
	RegistrationRequest *RegistrationRequest `json:"registration_request,omitempty" gorm:"foreignKey:CompanyID"`
}

// UpdateCompanyDTO para atualização de empresa (campos limitados)
type UpdateCompanyDTO struct {
	TradeName  string `json:"trade_name" example:"Silva Consultoria"`
	Address    string `json:"address" example:"Rua das Flores, 123"`
	PostalCode string `json:"postal_code" example:"1000-001"`
	City       string `json:"city" example:"Lisboa"`
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

// AdminUpdateCompanyDTO para contabilistas/admins editarem dados de empresas
type AdminUpdateCompanyDTO struct {
	CompanyName     *string  `json:"company_name,omitempty"`
	NIPC           *string  `json:"nipc,omitempty"`
	CAE            *string  `json:"cae,omitempty"`
	LegalForm      *string  `json:"legal_form,omitempty"`
	FoundingDate   *string  `json:"founding_date,omitempty"`
	TradeName      *string  `json:"trade_name,omitempty"`
	CorporateObject *string `json:"corporate_object,omitempty"`
	
	// Morada da empresa
	Address        *string  `json:"address,omitempty"`
	PostalCode     *string  `json:"postal_code,omitempty"`
	City           *string  `json:"city,omitempty"`
	County         *string  `json:"county,omitempty"`
	District       *string  `json:"district,omitempty"`
	
	// Dados financeiros
	ShareCapital   *float64 `json:"share_capital,omitempty"`
	BankName       *string  `json:"bank_name,omitempty"`
	IBAN           *string  `json:"iban,omitempty"`
	BIC            *string  `json:"bic,omitempty"`
	
	// Regimes
	AccountingRegime *string `json:"accounting_regime,omitempty"`
	VATRegime       *string  `json:"vat_regime,omitempty"`
	
	// Operacionais
	BusinessActivity   *string  `json:"business_activity,omitempty"`
	EstimatedRevenue   *float64 `json:"estimated_revenue,omitempty"`
	AnnualRevenue      *float64 `json:"annual_revenue,omitempty"`
	MonthlyInvoices    *int     `json:"monthly_invoices,omitempty"`
	NumberEmployees    *int     `json:"number_employees,omitempty"`
	ReportFrequency    *string  `json:"report_frequency,omitempty"`
	HasStock           *bool    `json:"has_stock,omitempty"`
	MainClients        *string  `json:"main_clients,omitempty"`
	MainSuppliers      *string  `json:"main_suppliers,omitempty"`
}
