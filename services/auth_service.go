package services

import (
	"RVContabilidadeBack/config"
	"RVContabilidadeBack/models"
	"RVContabilidadeBack/utils"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// RegisterClient cria uma nova solicitação de registo
func (s *AuthService) RegisterClient(req models.RegistrationRequestDTO) (*models.RegistrationRequest, error) {
	// Verificar se já existe solicitação pendente
	if err := s.checkExistingRequest(req.NIF, req.Email); err != nil {
		return nil, err
	}

	// Verificar se já existe utilizador aprovado
	if err := s.checkExistingUser(req.NIF, req.Email); err != nil {
		return nil, err
	}

	// Hash da password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("erro ao processar password")
	}

	// Criar solicitação
	registrationRequest := s.buildRegistrationRequest(req, string(hashedPassword))

	// Salvar na base de dados
	if err := config.DB.Create(&registrationRequest).Error; err != nil {
		return nil, err
	}

	return &registrationRequest, nil
}

// Login autentica um utilizador
func (s *AuthService) Login(username, password string) (*models.AuthResponse, error) {
	var user models.User

	// Procurar utilizador
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Verificar status
	if user.Status != string(models.StatusApproved) {
		return nil, s.getStatusError(user.Status)
	}

	// Verificar password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Gerar token
	token, err := utils.GenerateToken(user.ID, user.Username, user.NIF, user.Role)
	if err != nil {
		return nil, errors.New("erro ao gerar token")
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// Register é um alias para CreateUserDirect (para compatibilidade)
func (s *AuthService) Register(req models.RegisterRequest) (*models.AuthResponse, error) {
	return s.CreateUserDirect(req)
}

// LoginWithCredentials com LoginRequest DTO
func (s *AuthService) LoginWithCredentials(req models.LoginRequest) (*models.AuthResponse, error) {
	return s.Login(req.Username, req.Password)
}

// CreateUserDirect cria utilizador diretamente (para uso interno)
func (s *AuthService) CreateUserDirect(req models.RegisterRequest) (*models.AuthResponse, error) {
	// Verificar duplicações
	if err := s.checkUserDuplicates(req.Username, req.Email, req.NIF); err != nil {
		return nil, err
	}

	// Hash da password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("erro ao processar password")
	}

	// Criar utilizador
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Phone:    req.Phone,
		NIF:      req.NIF,
		Role:     req.Role,
		Status:   req.Status,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return nil, errors.New("erro ao criar utilizador")
	}

	// Gerar token
	token, err := utils.GenerateToken(user.ID, user.Username, user.NIF, user.Role)
	if err != nil {
		return nil, errors.New("erro ao gerar token")
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// ===== MÉTODOS PRIVADOS =====

func (s *AuthService) checkExistingRequest(nif, email string) error {
	var existingRequest models.RegistrationRequest
	whereClause := "email = ? AND status = ?"
	whereArgs := []interface{}{email, "pending"}
	
	if nif != "" {
		whereClause = "(nif = ? OR email = ?) AND status = ?"
		whereArgs = []interface{}{nif, email, "pending"}
	}
	
	if config.DB.Where(whereClause, whereArgs...).First(&existingRequest).Error == nil {
		if nif != "" && existingRequest.NIF != nil && *existingRequest.NIF == nif {
			return errors.New("já existe uma solicitação pendente com este NIF")
		}
		if existingRequest.Email != nil && *existingRequest.Email == email {
			return errors.New("já existe uma solicitação pendente com este email")
		}
	}
	return nil
}

func (s *AuthService) checkExistingUser(nif, email string) error {
	var existingUser models.User
	whereClause := "email = ?"
	whereArgs := []interface{}{email}
	
	if nif != "" {
		whereClause = "nif = ? OR email = ?"
		whereArgs = []interface{}{nif, email}
	}
	
	if config.DB.Where(whereClause, whereArgs...).First(&existingUser).Error == nil {
		if nif != "" && existingUser.NIF == nif {
			return errors.New("já existe uma conta aprovada com este NIF")
		}
		if existingUser.Email == email {
			return errors.New("já existe uma conta aprovada com este email")
		}
	}
	return nil
}

func (s *AuthService) checkUserDuplicates(username, email, nif string) error {
	var existingUser models.User
	
	if err := config.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return errors.New("username já está em uso")
	}
	
	if err := config.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return errors.New("email já está em uso")
	}
	
	if err := config.DB.Where("nif = ?", nif).First(&existingUser).Error; err == nil {
		return errors.New("NIF já está em uso")
	}
	
	return nil
}

func (s *AuthService) getStatusError(status string) error {
	statusMessages := map[string]string{
		string(models.StatusPending):  "conta aguarda aprovação da contabilista",
		string(models.StatusRejected): "conta foi rejeitada. Contacte o suporte.",
		string(models.StatusBlocked):  "conta foi bloqueada. Contacte o suporte.",
	}
	
	message := statusMessages[status]
	if message == "" {
		message = "acesso negado"
	}
	
	return errors.New(message)
}

func (s *AuthService) buildRegistrationRequest(req models.RegistrationRequestDTO, hashedPassword string) models.RegistrationRequest {
	registrationRequest := models.RegistrationRequest{
		RequestType:   "new_client",
		Status:        "pending",
		ApprovalToken: utils.GenerateRandomToken(),
		
		// Dados obrigatórios
		Username:     req.Username,
		PasswordHash: hashedPassword,
		LegalForm:    req.LegalForm,
	}

	// Campos opcionais do usuário
	if req.Name != "" {
		registrationRequest.Name = &req.Name
	}
	if req.Email != "" {
		registrationRequest.Email = &req.Email
	}
	if req.Phone != "" {
		registrationRequest.Phone = &req.Phone
	}
	if req.NIF != "" {
		registrationRequest.NIF = &req.NIF
	}
	if req.FiscalAddress != "" {
		registrationRequest.FiscalAddress = &req.FiscalAddress
	}
	if req.FiscalPostalCode != "" {
		registrationRequest.FiscalPostalCode = &req.FiscalPostalCode
	}
	if req.FiscalCity != "" {
		registrationRequest.FiscalCity = &req.FiscalCity
	}
	
	// Novos campos de morada da empresa
	if req.Address != "" {
		registrationRequest.Address = &req.Address
	}
	if req.PostalCode != "" {
		registrationRequest.PostalCode = &req.PostalCode
	}
	if req.City != "" {
		registrationRequest.City = &req.City
	}
	if req.Country != "" {
		registrationRequest.Country = &req.Country
	}
	
	// Campos opcionais da empresa
	if req.CompanyName != "" {
		registrationRequest.CompanyName = &req.CompanyName
	}
	if req.TradeName != "" {
		registrationRequest.TradeName = &req.TradeName
	}
	if req.CAE != "" {
		registrationRequest.CAE = &req.CAE
	}
	if req.FoundingDate != "" {
		if foundingDate, err := time.Parse("2006-01-02", req.FoundingDate); err == nil {
			registrationRequest.FoundingDate = &foundingDate
		}
	}
	if req.ShareCapital != nil && req.ShareCapital.Value != nil {
		registrationRequest.ShareCapital = req.ShareCapital.Value
	}
	
	// Campos obrigatórios da empresa
	registrationRequest.NIPC = req.NIPC
	
	// Mapear todos os campos adicionais que o frontend está enviando
	if req.DateOfBirth != nil && *req.DateOfBirth != "" {
		if dateOfBirth, err := time.Parse("2006-01-02", *req.DateOfBirth); err == nil {
			registrationRequest.DateOfBirth = &dateOfBirth
		}
	}
	
	if req.AccountingRegime != nil {
		registrationRequest.AccountingRegime = req.AccountingRegime
	}
	if req.VATRegime != nil {
		registrationRequest.VATRegime = req.VATRegime
	}
	if req.BusinessActivity != nil {
		registrationRequest.BusinessActivity = req.BusinessActivity
	}
	if req.MonthlyInvoices != nil && req.MonthlyInvoices.Value != nil {
		registrationRequest.MonthlyInvoices = req.MonthlyInvoices.Value
	}
	if req.ReportFrequency != nil {
		registrationRequest.ReportFrequency = req.ReportFrequency
	}
	
	// Campos pessoais adicionais
	if req.MaritalStatus != nil {
		registrationRequest.MaritalStatus = req.MaritalStatus
	}
	if req.CitizenCardNumber != nil {
		registrationRequest.CitizenCardNumber = req.CitizenCardNumber
	}
	if req.CitizenCardExpiry != nil && *req.CitizenCardExpiry != "" {
		if cardExpiry, err := time.Parse("2006-01-02", *req.CitizenCardExpiry); err == nil {
			registrationRequest.CitizenCardExpiry = &cardExpiry
		}
	}
	if req.TaxResidenceCountry != nil {
		registrationRequest.TaxResidenceCountry = req.TaxResidenceCountry
	}
	if req.FixedPhone != nil {
		registrationRequest.FixedPhone = req.FixedPhone
	}
	if req.FiscalCounty != nil {
		registrationRequest.FiscalCounty = req.FiscalCounty
	}
	if req.FiscalDistrict != nil {
		registrationRequest.FiscalDistrict = req.FiscalDistrict
	}
	
	// Preferências
	if req.OfficialEmail != nil {
		registrationRequest.OfficialEmail = req.OfficialEmail
	}
	if req.BillingSoftware != nil {
		registrationRequest.BillingSoftware = req.BillingSoftware
	}
	if req.PreferredFormat != nil {
		registrationRequest.PreferredFormat = req.PreferredFormat
	}
	if req.PreferredContactHours != nil {
		registrationRequest.PreferredContactHours = req.PreferredContactHours
	}
	
	// Dados completos da empresa
	if req.EstimatedRevenue != nil && req.EstimatedRevenue.Value != nil {
		registrationRequest.EstimatedRevenue = req.EstimatedRevenue.Value
	}
	if req.NumberEmployees != nil && req.NumberEmployees.Value != nil {
		registrationRequest.NumberEmployees = req.NumberEmployees.Value
	}
	if req.CorporateObject != nil {
		registrationRequest.CorporateObject = req.CorporateObject
	}
	if req.CompanyAddress != nil {
		registrationRequest.CompanyAddress = req.CompanyAddress
	}
	if req.CompanyPostalCode != nil {
		registrationRequest.CompanyPostalCode = req.CompanyPostalCode
	}
	if req.CompanyCity != nil {
		registrationRequest.CompanyCity = req.CompanyCity
	}
	if req.CompanyCounty != nil {
		registrationRequest.CompanyCounty = req.CompanyCounty
	}
	if req.CompanyDistrict != nil {
		registrationRequest.CompanyDistrict = req.CompanyDistrict
	}
	if req.CompanyCountry != nil {
		registrationRequest.CompanyCountry = req.CompanyCountry
	}
	if req.GroupStartDate != nil && *req.GroupStartDate != "" {
		if groupStartDate, err := time.Parse("2006-01-02", *req.GroupStartDate); err == nil {
			registrationRequest.GroupStartDate = &groupStartDate
		}
	}
	
	// Informação bancária
	if req.BankName != nil {
		registrationRequest.BankName = req.BankName
	}
	if req.IBAN != nil {
		registrationRequest.IBAN = req.IBAN
	}
	if req.BIC != nil {
		registrationRequest.BIC = req.BIC
	}
	
	// Dados operacionais
	if req.AnnualRevenue != nil && req.AnnualRevenue.Value != nil {
		registrationRequest.AnnualRevenue = req.AnnualRevenue.Value
	}
	if req.HasStock != nil {
		registrationRequest.HasStock = req.HasStock
	}
	if req.MainClients != nil {
		registrationRequest.MainClients = req.MainClients
	}
	if req.MainSuppliers != nil {
		registrationRequest.MainSuppliers = req.MainSuppliers
	}

	return registrationRequest
}
