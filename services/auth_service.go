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
		if nif != "" && existingRequest.NIF == nif {
			return errors.New("já existe uma solicitação pendente com este NIF")
		}
		if existingRequest.Email == email {
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
		
		// Dados obrigatórios do user
		Username:         req.Username,
		Name:             req.Name,
		Email:            req.Email,
		Phone:            req.Phone,
		NIF:              req.NIF,
		PasswordHash:     hashedPassword,
		FiscalAddress:    req.FiscalAddress,
		FiscalPostalCode: req.FiscalPostalCode,
		FiscalCity:       req.FiscalCity,
		
		// Dados obrigatórios da company
		CompanyName: req.CompanyName,
		NIPC:        req.NIPC,
		LegalForm:   req.LegalForm,
	}

	// Processar campos opcionais do User
	if req.DateOfBirth != nil {
		if dateOfBirth, err := time.Parse("2006-01-02", *req.DateOfBirth); err == nil {
			registrationRequest.DateOfBirth = &dateOfBirth
		}
	}
	
	registrationRequest.MaritalStatus = req.MaritalStatus
	registrationRequest.CitizenCardNumber = req.CitizenCardNumber
	
	if req.CitizenCardExpiry != nil {
		if citizenCardExpiry, err := time.Parse("2006-01-02", *req.CitizenCardExpiry); err == nil {
			registrationRequest.CitizenCardExpiry = &citizenCardExpiry
		}
	}
	
	registrationRequest.TaxResidenceCountry = req.TaxResidenceCountry
	registrationRequest.FixedPhone = req.FixedPhone
	registrationRequest.FiscalCounty = req.FiscalCounty
	registrationRequest.FiscalDistrict = req.FiscalDistrict
	registrationRequest.OfficialEmail = req.OfficialEmail
	registrationRequest.BillingSoftware = req.BillingSoftware
	registrationRequest.PreferredFormat = req.PreferredFormat
	registrationRequest.ReportFrequency = req.ReportFrequency
	registrationRequest.PreferredContactHours = req.PreferredContactHours

	// Processar campos opcionais da Company
	registrationRequest.CAE = req.CAE
	
	if req.FoundingDate != nil {
		if foundingDate, err := time.Parse("2006-01-02", *req.FoundingDate); err == nil {
			registrationRequest.FoundingDate = &foundingDate
		}
	}
	
	registrationRequest.AccountingRegime = req.AccountingRegime
	registrationRequest.VATRegime = req.VATRegime
	registrationRequest.BusinessActivity = req.BusinessActivity
	registrationRequest.EstimatedRevenue = req.EstimatedRevenue
	registrationRequest.MonthlyInvoices = req.MonthlyInvoices
	registrationRequest.NumberEmployees = req.NumberEmployees
	registrationRequest.TradeName = req.TradeName
	registrationRequest.CorporateObject = req.CorporateObject
	registrationRequest.CompanyAddress = req.CompanyAddress
	registrationRequest.CompanyPostalCode = req.CompanyPostalCode
	registrationRequest.CompanyCity = req.CompanyCity
	registrationRequest.CompanyCounty = req.CompanyCounty
	registrationRequest.CompanyDistrict = req.CompanyDistrict
	registrationRequest.CompanyCountry = req.CompanyCountry
	registrationRequest.ShareCapital = req.ShareCapital
	
	if req.GroupStartDate != nil {
		if groupStartDate, err := time.Parse("2006-01-02", *req.GroupStartDate); err == nil {
			registrationRequest.GroupStartDate = &groupStartDate
		}
	}
	
	registrationRequest.BankName = req.BankName
	registrationRequest.IBAN = req.IBAN
	registrationRequest.BIC = req.BIC
	registrationRequest.AnnualRevenue = req.AnnualRevenue
	registrationRequest.HasStock = req.HasStock
	registrationRequest.MainClients = req.MainClients
	registrationRequest.MainSuppliers = req.MainSuppliers

	return registrationRequest
}
