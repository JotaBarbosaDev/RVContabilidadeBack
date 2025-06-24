package services

import (
	"RVContabilidadeBack/config"
	"RVContabilidadeBack/models"
	"errors"
	"time"
)

type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

// GetPendingRequests obtém todas as solicitações pendentes
func (s *AdminService) GetPendingRequests() ([]models.PendingRequestResponseDTO, error) {
	var requests []models.RegistrationRequest
	if err := config.DB.Where("status = ?", "pending").Find(&requests).Error; err != nil {
		return nil, errors.New("erro ao obter solicitações pendentes")
	}

	// Converter para DTO
	var response []models.PendingRequestResponseDTO
	for _, req := range requests {
		dto := models.PendingRequestResponseDTO{
			ID:          req.ID,
			RequestType: req.RequestType,
			Status:      req.Status,
			SubmittedAt: req.SubmittedAt,
			Username:    req.Username,
			Name:        req.Name,
			Email:       req.Email,
			Phone:       req.Phone,
			NIF:         req.NIF,
			CompanyName: req.CompanyName,
			NIPC:        req.NIPC,
		}
		response = append(response, dto)
	}

	return response, nil
}

// GetAllRequests obtém todos os pedidos de registo
func (s *AdminService) GetAllRequests() ([]models.RegistrationRequest, error) {
	var requests []models.RegistrationRequest
	if err := config.DB.Find(&requests).Error; err != nil {
		return nil, errors.New("erro ao obter pedidos de registo")
	}
	return requests, nil
}

// GetRequestDetails obtém detalhes de um pedido específico
func (s *AdminService) GetRequestDetails(requestID uint) (*models.RegistrationRequest, error) {
	var request models.RegistrationRequest
	if err := config.DB.First(&request, requestID).Error; err != nil {
		return nil, errors.New("pedido não encontrado")
	}
	return &request, nil
}

// ApproveRequest aprova ou rejeita uma solicitação
func (s *AdminService) ApproveRequest(req models.ApprovalRequestDTO, reviewerID uint) (*models.RegistrationRequest, error) {
	// Buscar solicitação
	var request models.RegistrationRequest
	if err := config.DB.First(&request, req.RequestID).Error; err != nil {
		return nil, errors.New("solicitação não encontrada")
	}

	// Verificar se ainda está pendente
	if request.Status != "pending" {
		return nil, errors.New("solicitação já foi processada")
	}

	// Atualizar dados de review
	now := time.Now()
	request.Status = req.Status
	request.ReviewedAt = &now
	request.ReviewedBy = &reviewerID
	request.ReviewNotes = req.ReviewNotes

	// Se aprovado, criar User e Company
	if req.Status == "approved" {
		userID, companyID, err := s.createUserAndCompany(request)
		if err != nil {
			return nil, err
		}

		// Atualizar RegistrationRequest com os IDs
		request.UserID = &userID
		request.CompanyID = &companyID
	}

	// Salvar alterações
	if err := config.DB.Save(&request).Error; err != nil {
		return nil, errors.New("erro ao salvar alterações na solicitação")
	}

	return &request, nil
}

// GetAllUsers obtém todos os utilizadores com filtros
func (s *AdminService) GetAllUsers(status, role string) ([]models.User, map[string]interface{}, error) {
	var users []models.User
	
	query := config.DB
	
	// Aplicar filtros
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	
	// Buscar usuários
	if err := query.Find(&users).Error; err != nil {
		return nil, nil, errors.New("erro ao obter utilizadores")
	}

	// Para cada usuário, buscar a empresa separadamente
	for i := range users {
		if users[i].Role == "client" {
			var company models.Company
			if config.DB.Where("user_id = ?", users[i].ID).First(&company).Error == nil {
				users[i].Company = &company
			}
		}
	}

	// Calcular estatísticas
	stats := s.calculateUserStats()

	return users, stats, nil
}

// GetUserDetails obtém detalhes de um utilizador
func (s *AdminService) GetUserDetails(userID uint) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("utilizador não encontrado")
	}

	// Buscar empresa se for cliente
	if user.Role == "client" {
		var company models.Company
		if config.DB.Where("user_id = ?", userID).First(&company).Error == nil {
			user.Company = &company
		}
	}

	return &user, nil
}

// GetApprovedClients obtém todos os clientes aprovados
func (s *AdminService) GetApprovedClients() ([]models.User, error) {
	var users []models.User
	
	if err := config.DB.Where("role = ? AND status = ?", "client", "approved").Find(&users).Error; err != nil {
		return nil, errors.New("erro ao obter clientes aprovados")
	}

	// Para cada cliente, buscar a empresa
	for i := range users {
		var company models.Company
		if config.DB.Where("user_id = ?", users[i].ID).First(&company).Error == nil {
			users[i].Company = &company
		}
	}

	return users, nil
}

// GetUsersCount conta utilizadores por status e role
func (s *AdminService) GetUsersCount() (map[string]int64, error) {
	counts := make(map[string]int64)

	var total, approved, pending, clients, admins int64

	// Contar todos
	config.DB.Model(&models.User{}).Count(&total)
	counts["total"] = total
	
	// Contar por status
	config.DB.Model(&models.User{}).Where("status = ?", "approved").Count(&approved)
	counts["approved"] = approved
	
	config.DB.Model(&models.User{}).Where("status = ?", "pending").Count(&pending)
	counts["pending"] = pending
	
	// Contar por role
	config.DB.Model(&models.User{}).Where("role = ?", "client").Count(&clients)
	counts["clients"] = clients
	
	config.DB.Model(&models.User{}).Where("role IN ?", []string{"admin", "accountant"}).Count(&admins)
	counts["admins"] = admins

	return counts, nil
}

// ===== MÉTODOS PRIVADOS =====

func (s *AdminService) createUserAndCompany(request models.RegistrationRequest) (uint, uint, error) {
	// Criar User
	user := models.User{
		Username:            request.Username,
		Email:               request.Email,
		Password:            request.PasswordHash,
		Name:                request.Name,
		Phone:               request.Phone,
		NIF:                 request.NIF,
		Role:                "client",
		Status:              "approved",
		FiscalAddress:       request.FiscalAddress,
		FiscalPostalCode:    request.FiscalPostalCode,
		FiscalCity:          request.FiscalCity,
		TaxResidenceCountry: "Portugal",
	}

	// Campos opcionais do User
	if request.DateOfBirth != nil {
		user.DateOfBirth = request.DateOfBirth
	}
	if request.MaritalStatus != nil {
		user.MaritalStatus = *request.MaritalStatus
	}
	if request.CitizenCardNumber != nil {
		user.CitizenCardNumber = *request.CitizenCardNumber
	}
	if request.CitizenCardExpiry != nil {
		user.CitizenCardExpiry = request.CitizenCardExpiry
	}
	if request.TaxResidenceCountry != nil {
		user.TaxResidenceCountry = *request.TaxResidenceCountry
	}
	if request.FixedPhone != nil {
		user.FixedPhone = *request.FixedPhone
	}
	if request.FiscalCounty != nil {
		user.FiscalCounty = *request.FiscalCounty
	}
	if request.FiscalDistrict != nil {
		user.FiscalDistrict = *request.FiscalDistrict
	}
	if request.OfficialEmail != nil {
		user.OfficialEmail = *request.OfficialEmail
	}
	if request.BillingSoftware != nil {
		user.BillingSoftware = *request.BillingSoftware
	}
	if request.PreferredFormat != nil {
		user.PreferredFormat = *request.PreferredFormat
	} else {
		user.PreferredFormat = "digital"
	}
	if request.ReportFrequency != nil {
		user.ReportFrequency = *request.ReportFrequency
	} else {
		user.ReportFrequency = "mensal"
	}
	if request.PreferredContactHours != nil {
		user.PreferredContactHours = *request.PreferredContactHours
	}

	// Salvar User
	if err := config.DB.Create(&user).Error; err != nil {
		return 0, 0, errors.New("erro ao criar utilizador: " + err.Error())
	}

	// Verificar se já existe empresa com este NIPC
	if request.NIPC != "" {
		var existingCompany models.Company
		if config.DB.Where("nipc = ?", request.NIPC).First(&existingCompany).Error == nil {
			return 0, 0, errors.New("já existe uma empresa com este NIPC")
		}
	}

	// Criar Company
	company := models.Company{
		UserID:      user.ID,
		CompanyName: request.CompanyName,
		NIPC:        request.NIPC,
		LegalForm:   request.LegalForm,
		Country:     "Portugal",
		Status:      "active",
	}

	// Campos opcionais da Company
	if request.CAE != nil {
		company.CAE = *request.CAE
	}
	if request.FoundingDate != nil {
		company.FoundingDate = request.FoundingDate
	}
	if request.AccountingRegime != nil {
		company.AccountingRegime = *request.AccountingRegime
	}
	if request.VATRegime != nil {
		company.VATRegime = *request.VATRegime
	}
	if request.BusinessActivity != nil {
		company.BusinessActivity = *request.BusinessActivity
	}
	if request.EstimatedRevenue != nil {
		company.EstimatedRevenue = *request.EstimatedRevenue
	}
	if request.MonthlyInvoices != nil {
		company.MonthlyInvoices = *request.MonthlyInvoices
	}
	if request.NumberEmployees != nil {
		company.NumberEmployees = *request.NumberEmployees
	}
	if request.TradeName != nil {
		company.TradeName = *request.TradeName
	}
	if request.CorporateObject != nil {
		company.CorporateObject = *request.CorporateObject
	}
	if request.CompanyAddress != nil {
		company.Address = *request.CompanyAddress
	}
	if request.CompanyPostalCode != nil {
		company.PostalCode = *request.CompanyPostalCode
	}
	if request.CompanyCity != nil {
		company.City = *request.CompanyCity
	}
	if request.CompanyCounty != nil {
		company.County = *request.CompanyCounty
	}
	if request.CompanyDistrict != nil {
		company.District = *request.CompanyDistrict
	}
	if request.CompanyCountry != nil {
		company.Country = *request.CompanyCountry
	}
	if request.ShareCapital != nil {
		company.ShareCapital = *request.ShareCapital
	}
	if request.GroupStartDate != nil {
		company.GroupStartDate = request.GroupStartDate
	}
	if request.BankName != nil {
		company.BankName = *request.BankName
	}
	if request.IBAN != nil {
		company.IBAN = *request.IBAN
	}
	if request.BIC != nil {
		company.BIC = *request.BIC
	}
	if request.AnnualRevenue != nil {
		company.AnnualRevenue = *request.AnnualRevenue
	}
	if request.HasStock != nil {
		company.HasStock = *request.HasStock
	}
	if request.MainClients != nil {
		company.MainClients = *request.MainClients
	}
	if request.MainSuppliers != nil {
		company.MainSuppliers = *request.MainSuppliers
	}

	// Salvar Company
	if err := config.DB.Create(&company).Error; err != nil {
		// Se falhar, eliminar o utilizador criado
		config.DB.Delete(&user)
		return 0, 0, errors.New("erro ao criar empresa: " + err.Error())
	}

	return user.ID, company.ID, nil
}

func (s *AdminService) calculateUserStats() map[string]interface{} {
	stats := map[string]interface{}{}

	var totalClients, approvedClients, pendingClients, rejectedClients, blockedClients, clientsWithCompany int64

	// Contar clientes por status
	config.DB.Model(&models.User{}).Where("role = ?", "client").Count(&totalClients)
	config.DB.Model(&models.User{}).Where("role = ? AND status = ?", "client", "approved").Count(&approvedClients)
	config.DB.Model(&models.User{}).Where("role = ? AND status = ?", "client", "pending").Count(&pendingClients)
	config.DB.Model(&models.User{}).Where("role = ? AND status = ?", "client", "rejected").Count(&rejectedClients)
	config.DB.Model(&models.User{}).Where("role = ? AND status = ?", "client", "blocked").Count(&blockedClients)
	
	// Contar clientes com empresa
	config.DB.Table("users").
		Joins("JOIN companies ON companies.user_id = users.id").
		Where("users.role = ?", "client").
		Count(&clientsWithCompany)

	stats["total_clients"] = totalClients
	stats["approved_clients"] = approvedClients
	stats["pending_clients"] = pendingClients
	stats["rejected_clients"] = rejectedClients
	stats["blocked_clients"] = blockedClients
	stats["clients_with_company"] = clientsWithCompany

	return stats
}

// UpdateUserStatus atualiza o status de um utilizador
func (s *AdminService) UpdateUserStatus(userID uint, newStatus string) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("utilizador não encontrado")
	}

	user.Status = newStatus
	if err := config.DB.Save(&user).Error; err != nil {
		return nil, errors.New("erro ao atualizar utilizador")
	}

	return &user, nil
}

// UpdateClientData atualiza dados pessoais de um cliente
func (s *AdminService) UpdateClientData(clientID uint, req models.AdminUpdateClientDTO) error {
	// Verificar se o cliente existe e é cliente aprovado
	var client models.User
	if err := config.DB.Where("id = ? AND role = ? AND status = ?", clientID, "client", "approved").First(&client).Error; err != nil {
		return errors.New("cliente aprovado não encontrado")
	}

	// Atualizar apenas os campos fornecidos
	updateData := make(map[string]interface{})
	
	if req.Name != nil {
		updateData["name"] = *req.Name
	}
	if req.Email != nil {
		updateData["email"] = *req.Email
	}
	if req.Phone != nil {
		updateData["phone"] = *req.Phone
	}
	if req.NIF != nil {
		updateData["nif"] = *req.NIF
	}
	if req.DateOfBirth != nil {
		if parsedDate, err := time.Parse("2006-01-02", *req.DateOfBirth); err == nil {
			updateData["date_of_birth"] = parsedDate
		}
	}
	if req.MaritalStatus != nil {
		updateData["marital_status"] = *req.MaritalStatus
	}
	if req.CitizenCardNumber != nil {
		updateData["citizen_card_number"] = *req.CitizenCardNumber
	}
	if req.CitizenCardExpiry != nil {
		if parsedDate, err := time.Parse("2006-01-02", *req.CitizenCardExpiry); err == nil {
			updateData["citizen_card_expiry"] = parsedDate
		}
	}
	if req.FixedPhone != nil {
		updateData["fixed_phone"] = *req.FixedPhone
	}
	if req.FiscalAddress != nil {
		updateData["fiscal_address"] = *req.FiscalAddress
	}
	if req.FiscalPostalCode != nil {
		updateData["fiscal_postal_code"] = *req.FiscalPostalCode
	}
	if req.FiscalCity != nil {
		updateData["fiscal_city"] = *req.FiscalCity
	}
	if req.FiscalCounty != nil {
		updateData["fiscal_county"] = *req.FiscalCounty
	}
	if req.FiscalDistrict != nil {
		updateData["fiscal_district"] = *req.FiscalDistrict
	}
	if req.OfficialEmail != nil {
		updateData["official_email"] = *req.OfficialEmail
	}
	if req.BillingSoftware != nil {
		updateData["billing_software"] = *req.BillingSoftware
	}
	if req.PreferredFormat != nil {
		updateData["preferred_format"] = *req.PreferredFormat
	}
	if req.PreferredContactHours != nil {
		updateData["preferred_contact_hours"] = *req.PreferredContactHours
	}
	if req.Status != nil {
		updateData["status"] = *req.Status
	}

	if err := config.DB.Model(&client).Updates(updateData).Error; err != nil {
		return errors.New("erro ao atualizar dados do cliente")
	}

	return nil
}

// UpdateClientCompany atualiza dados da empresa de um cliente
func (s *AdminService) UpdateClientCompany(clientID uint, req models.AdminUpdateCompanyDTO) (*models.Company, error) {
	// Verificar se o cliente existe e é cliente aprovado
	var client models.User
	if err := config.DB.Where("id = ? AND role = ? AND status = ?", clientID, "client", "approved").First(&client).Error; err != nil {
		return nil, errors.New("cliente aprovado não encontrado")
	}

	// Encontrar a empresa do cliente
	var company models.Company
	if err := config.DB.Where("user_id = ?", clientID).First(&company).Error; err != nil {
		return nil, errors.New("empresa do cliente não encontrada")
	}

	// Atualizar apenas os campos fornecidos
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
		if parsedDate, err := time.Parse("2006-01-02", *req.FoundingDate); err == nil {
			updateData["founding_date"] = parsedDate
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

	return &company, nil
}

// DeleteClient elimina um cliente e a sua empresa
func (s *AdminService) DeleteClient(clientID uint) error {
	// Verificar se o cliente existe e é cliente
	var client models.User
	if err := config.DB.Where("id = ? AND role = ?", clientID, "client").First(&client).Error; err != nil {
		return errors.New("cliente não encontrado")
	}

	// Iniciar transação para eliminar cliente e empresa
	tx := config.DB.Begin()

	// Eliminar empresa(s) do cliente
	if err := tx.Where("user_id = ?", clientID).Delete(&models.Company{}).Error; err != nil {
		tx.Rollback()
		return errors.New("erro ao eliminar empresa do cliente")
	}

	// Eliminar cliente
	if err := tx.Delete(&client).Error; err != nil {
		tx.Rollback()
		return errors.New("erro ao eliminar cliente")
	}

	tx.Commit()
	return nil
}

// GetAllUsersSimple obtém dados básicos dos utilizadores
func (s *AdminService) GetAllUsersSimple() ([]models.User, error) {
	var users []models.User
	
	// Buscar apenas os dados básicos dos usuários
	if err := config.DB.Select("id, username, name, email, role, status, created_at").Find(&users).Error; err != nil {
		return nil, errors.New("erro ao obter utilizadores: " + err.Error())
	}

	return users, nil
}
