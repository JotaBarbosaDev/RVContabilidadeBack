package handlers

import (
	"RVContabilidadeBack/config"
	"RVContabilidadeBack/models"
	"RVContabilidadeBack/utils"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterClient godoc
// @Summary      Registo de novo cliente
// @Description  Cria uma nova solicitação de registo para aprovação da contabilista
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.RegistrationRequestDTO  true  "Dados de registo completos"
// @Success      201      {object}  models.SuccessResponse
// @Router       /auth/register [post]
func RegisterClient(c *gin.Context) {
	var req models.RegistrationRequestDTO
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Verificar se já existe utilizador com este NIF
	var existingUser models.User
	userExists := config.DB.Where("nif = ?", req.NIF).First(&existingUser).Error == nil
	
	var user models.User
	
	if userExists {
		// Utilizador já existe, verificar status
		switch existingUser.Status {
		case string(models.StatusApproved):
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Success: false,
				Error:   "Já tem uma conta aprovada com este NIF. Faça login.",
			})
			return
		case string(models.StatusPending):
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Success: false,
				Error:   "Já tem uma solicitação pendente com este NIF. Aguarde aprovação.",
			})
			return
		case string(models.StatusRejected):
			// Pode submeter nova solicitação
			user = existingUser
		}
	} else {
		// Verificar username único
		if config.DB.Where("username = ?", req.Username).First(&models.User{}).Error == nil {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Success: false,
				Error:   "Username já está em uso.",
			})
			return
		}

		// Verificar email único
		if config.DB.Where("email = ?", req.Email).First(&models.User{}).Error == nil {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Success: false,
				Error:   "Email já está em uso.",
			})
			return
		}

		// Parse da data de nascimento
		dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "Data de nascimento inválida. Use formato YYYY-MM-DD.",
			})
			return
		}

		// Criar novo utilizador com dados mínimos
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   "Erro ao processar password",
			})
			return
		}

		user = models.User{
			Username:            req.Username,
			Name:                req.Name,
			Email:               req.Email,
			Phone:               req.Phone,
			NIF:                 req.NIF,
			Password:            string(hashedPassword),
			DateOfBirth:         &dateOfBirth,
			FiscalAddress:       req.FiscalAddress,
			FiscalPostalCode:    req.FiscalPostalCode,
			FiscalCity:          req.FiscalCity,
			TaxResidenceCountry: "Portugal",
			ReportFrequency:     req.ReportFrequency,
			Role:                "client",
			Status:              string(models.StatusPending),
		}

		if err := config.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   "Erro ao criar utilizador",
			})
			return
		}
	}

	// Parse da data de constituição da empresa
	foundingDate, err := time.Parse("2006-01-02", req.FoundingDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Data de constituição inválida. Use formato YYYY-MM-DD.",
		})
		return
	}

	// Criar empresa com dados mínimos
	company := models.Company{
		UserID:           user.ID,
		CompanyName:      req.CompanyName,
		NIPC:            req.NIPC,
		CAE:             req.CAE,
		LegalForm:       req.LegalForm,
		FoundingDate:    &foundingDate,
		AccountingRegime: req.AccountingRegime,
		VATRegime:       req.VATRegime,
		BusinessActivity: req.BusinessActivity,
		EstimatedRevenue: req.EstimatedRevenue,
		MonthlyInvoices: req.MonthlyInvoices,
		NumberEmployees: req.NumberEmployees,
		Country:         "Portugal",
	}

	if err := config.DB.Create(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Erro ao criar empresa",
		})
		return
	}

	// Converter request para JSON (dados mínimos para aprovação)
	requestData := map[string]interface{}{
		"user_id":     user.ID,
		"company_id":  company.ID,
		"user_data":   user,
		"company_data": company,
	}
	requestDataJSON, _ := json.Marshal(requestData)

	// Determinar tipo de request
	requestType := "new_client"
	if userExists {
		requestType = "existing_client"
	}

	// Criar solicitação de registo
	registrationRequest := models.RegistrationRequest{
		UserID:        user.ID,
		RequestType:   requestType,
		RequestData:   string(requestDataJSON),
		Status:        "pending",
		ApprovalToken: utils.GenerateRandomToken(),
	}

	if err := config.DB.Create(&registrationRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Erro ao criar solicitação",
		})
		return
	}

	// Atualizar status do utilizador para pending (caso tenha sido rejeitado antes)
	user.Status = string(models.StatusPending)
	config.DB.Save(&user)

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Success: true,
		Message: "Solicitação enviada com sucesso. Aguarde aprovação da contabilista.",
		Data: gin.H{
			"request_id": registrationRequest.ID,
			"user_id":    user.ID,
		},
	})
}

// Register godoc (mantido para compatibilidade)
// @Summary      Criar conta diretamente
// @Description  Regista novo utilizador diretamente (para uso interno)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      models.RegisterRequest  true  "Dados de registo"
// @Success      201   {object}  models.AuthResponse
// @Router       /auth/register-direct [post]
func Register(c *gin.Context) {
	var req models.RegisterRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar se username já existe
	var existingUser models.User
	if err := config.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username já está em uso"})
		return
	}

	// Verificar se email já existe
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email já está em uso"})
		return
	}

	// Verificar se NIF já existe
	if err := config.DB.Where("nif = ?", req.NIF).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "NIF já está em uso"})
		return
	}

	// Encriptar password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar password"})
		return
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar utilizador"})
		return
	}

	// Gerar token
	token, err := utils.GenerateToken(user.ID, user.Username, user.NIF, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	// Resposta com token
	response := models.AuthResponse{
		Token: token,
		User:  user,
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary      Entrar
// @Description  Login com username e password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      models.LoginRequest  true  "Username e password"
// @Success      200          {object}  models.AuthResponse
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var req models.LoginRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	// Procurar utilizador por username
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	// Verificar status do utilizador
	if user.Status != string(models.StatusApproved) {
		statusMessages := map[string]string{
			string(models.StatusPending):  "Conta aguarda aprovação da contabilista",
			string(models.StatusRejected): "Conta foi rejeitada. Contacte o suporte.",
			string(models.StatusBlocked):  "Conta foi bloqueada. Contacte o suporte.",
		}
		
		message := statusMessages[user.Status]
		if message == "" {
			message = "Acesso negado"
		}
		
		c.JSON(http.StatusForbidden, gin.H{"error": message})
		return
	}

	// Verificar password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	// Gerar token
	token, err := utils.GenerateToken(user.ID, user.Username, user.NIF, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	// Resposta com token
	response := models.AuthResponse{
		Token: token,
		User:  user,
	}

	c.JSON(http.StatusOK, response)
}

// Logout godoc
// @Summary      Logout do utilizador
// @Description  Remove token do lado do cliente
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Router       /auth/logout [post]
func Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}

	userUsername, _ := c.Get("user_username")
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout realizado com sucesso",
		"user_id": userID,
		"username": userUsername,
	})
}

