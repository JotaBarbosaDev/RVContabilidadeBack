package services

import (
	"RVContabilidadeBack/config"
	"RVContabilidadeBack/models"
	"errors"
	"time"
)

type CompanyService struct{}

func NewCompanyService() *CompanyService {
	return &CompanyService{}
}

// GetCompanyByUserID obtém empresa pelo ID do utilizador
func (s *CompanyService) GetCompanyByUserID(userID uint) (*models.Company, error) {
	var company models.Company
	if err := config.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		return nil, errors.New("empresa não encontrada")
	}
	return &company, nil
}

// GetByUserID é um alias para GetCompanyByUserID
func (s *CompanyService) GetByUserID(userID uint) (*models.Company, error) {
	return s.GetCompanyByUserID(userID)
}

// UpdateCompany atualiza dados da empresa (campos limitados para cliente)
func (s *CompanyService) UpdateCompany(userID uint, req models.UpdateCompanyDTO) (*models.Company, error) {
	var company models.Company
	if err := config.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		return nil, errors.New("empresa não encontrada")
	}

	// Atualizar apenas campos permitidos para cliente
	if req.TradeName != "" {
		company.TradeName = req.TradeName
	}
	if req.Address != "" {
		company.Address = req.Address
	}
	if req.PostalCode != "" {
		company.PostalCode = req.PostalCode
	}
	if req.City != "" {
		company.City = req.City
	}

	if err := config.DB.Save(&company).Error; err != nil {
		return nil, errors.New("erro ao atualizar empresa")
	}

	return &company, nil
}

// CompleteCompanyData completa dados da empresa após aprovação
func (s *CompanyService) CompleteCompanyData(userID uint, req models.CompleteCompanyDataDTO) (*models.Company, error) {
	var company models.Company
	if err := config.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		return nil, errors.New("empresa não encontrada")
	}

	// Atualizar todos os campos
	company.TradeName = req.TradeName
	company.CorporateObject = req.CorporateObject
	company.Address = req.Address
	company.PostalCode = req.PostalCode
	company.City = req.City
	company.County = req.County
	company.District = req.District
	company.ShareCapital = req.ShareCapital
	company.BankName = req.BankName
	company.IBAN = req.IBAN
	company.BIC = req.BIC
	company.AnnualRevenue = req.AnnualRevenue
	company.HasStock = req.HasStock
	company.MainClients = req.MainClients
	company.MainSuppliers = req.MainSuppliers

	// Processar data de início do grupo
	if req.GroupStartDate != "" {
		if startDate, err := time.Parse("2006-01-02", req.GroupStartDate); err == nil {
			company.GroupStartDate = &startDate
		}
	}

	if err := config.DB.Save(&company).Error; err != nil {
		return nil, errors.New("erro ao completar dados da empresa")
	}

	return &company, nil
}

// AdminUpdateCompany atualiza empresa (para admin)
func (s *CompanyService) AdminUpdateCompany(clientID uint, req models.AdminUpdateCompanyDTO) (*models.Company, error) {
	// Verificar se o cliente existe
	var client models.User
	if err := config.DB.Where("id = ? AND role = ? AND status = ?", clientID, "client", "approved").First(&client).Error; err != nil {
		return nil, errors.New("cliente não encontrado ou não aprovado")
	}

	// Encontrar a empresa do cliente
	var company models.Company
	if err := config.DB.Where("user_id = ?", clientID).First(&company).Error; err != nil {
		return nil, errors.New("empresa não encontrada")
	}

	// Preparar dados para atualização
	updateData := make(map[string]interface{})
	
	if req.CompanyName != nil {
		updateData["company_name"] = *req.CompanyName
	}
	if req.NIPC != nil {
		updateData["nipc"] = *req.NIPC
	}
	if req.CAE != nil {
		updateData["cae"] = *req.CAE
	}
	if req.LegalForm != nil {
		updateData["legal_form"] = *req.LegalForm
	}
	if req.FoundingDate != nil {
		if foundingDate, err := time.Parse("2006-01-02", *req.FoundingDate); err == nil {
			updateData["founding_date"] = foundingDate
		}
	}
	if req.TradeName != nil {
		updateData["trade_name"] = *req.TradeName
	}
	if req.CorporateObject != nil {
		updateData["corporate_object"] = *req.CorporateObject
	}
	if req.Address != nil {
		updateData["address"] = *req.Address
	}
	if req.PostalCode != nil {
		updateData["postal_code"] = *req.PostalCode
	}
	if req.City != nil {
		updateData["city"] = *req.City
	}
	if req.County != nil {
		updateData["county"] = *req.County
	}
	if req.District != nil {
		updateData["district"] = *req.District
	}
	if req.ShareCapital != nil {
		updateData["share_capital"] = *req.ShareCapital
	}
	if req.BankName != nil {
		updateData["bank_name"] = *req.BankName
	}
	if req.IBAN != nil {
		updateData["iban"] = *req.IBAN
	}
	if req.BIC != nil {
		updateData["bic"] = *req.BIC
	}
	if req.AccountingRegime != nil {
		updateData["accounting_regime"] = *req.AccountingRegime
	}
	if req.VATRegime != nil {
		updateData["vat_regime"] = *req.VATRegime
	}
	if req.BusinessActivity != nil {
		updateData["business_activity"] = *req.BusinessActivity
	}
	if req.EstimatedRevenue != nil {
		updateData["estimated_revenue"] = *req.EstimatedRevenue
	}
	if req.AnnualRevenue != nil {
		updateData["annual_revenue"] = *req.AnnualRevenue
	}
	if req.MonthlyInvoices != nil {
		updateData["monthly_invoices"] = *req.MonthlyInvoices
	}
	if req.NumberEmployees != nil {
		updateData["number_employees"] = *req.NumberEmployees
	}
	if req.ReportFrequency != nil {
		updateData["report_frequency"] = *req.ReportFrequency
	}
	if req.HasStock != nil {
		updateData["has_stock"] = *req.HasStock
	}
	if req.MainClients != nil {
		updateData["main_clients"] = *req.MainClients
	}
	if req.MainSuppliers != nil {
		updateData["main_suppliers"] = *req.MainSuppliers
	}

	if err := config.DB.Model(&company).Updates(updateData).Error; err != nil {
		return nil, errors.New("erro ao atualizar dados da empresa")
	}

	// Recarregar empresa atualizada
	if err := config.DB.First(&company, company.ID).Error; err != nil {
		return nil, errors.New("erro ao recarregar empresa")
	}

	return &company, nil
}

// CreateCompany cria nova empresa para um utilizador
func (s *CompanyService) CreateCompany(userID uint, companyData models.Company) (*models.Company, error) {
	// Verificar se o utilizador já tem empresa
	var existingCompany models.Company
	if config.DB.Where("user_id = ?", userID).First(&existingCompany).Error == nil {
		return nil, errors.New("utilizador já tem empresa")
	}

	companyData.UserID = userID

	if err := config.DB.Create(&companyData).Error; err != nil {
		return nil, errors.New("erro ao criar empresa")
	}

	return &companyData, nil
}
