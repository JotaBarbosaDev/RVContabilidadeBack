package handlers

import (
	"RVContabilidadeBack/config"
	"RVContabilidadeBack/models"
	"RVContabilidadeBack/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary      Criar conta
// @Description  Regista novo utilizador (não precisa de token)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      models.RegisterRequest  true  "Dados: email, username, password, name, role, is_active"
// @Success      201   {object}  models.AuthResponse
// @Router       /auth/register [post]
func Register(c *gin.Context) {
	var req models.RegisterRequest
	
	// Validar dados recebidos
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar se email já existe
	var existingUser models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email já está em uso"})
		return
	}

	// Verificar se username já existe
	if err := config.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username já está em uso"})
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
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Username: req.Username,
		Role:     req.Role,
		IsActive: req.IsActive,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar utilizador"})
		return
	}

	// Gerar token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Username, user.Role)
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
// @Description  Login com username e password (não precisa de token)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      models.LoginRequest  true  "Dados: username e password"
// @Success      200          {object}  models.AuthResponse
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var req models.LoginRequest
	
	// Validar dados recebidos
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

	// Verificar se o utilizador está ativo
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Conta desativada"})
		return
	}

	// Verificar password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	// Gerar token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Username, user.Role)
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
// @Description  Remove token do lado do cliente (logout simples - requer token)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Router       /auth/logout [post]
func Logout(c *gin.Context) {
	// Para logout simples, só precisamos confirmar que o token é válido
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}

	username, _ := c.Get("user_username")
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout realizado com sucesso",
		"user_id": userID,
		"username": username,
	})
}

