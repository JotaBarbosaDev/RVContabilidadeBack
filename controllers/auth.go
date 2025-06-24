package controllers

import (
	"RVContabilidadeBack/models"
	"RVContabilidadeBack/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
authService = services.NewAuthService()
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

registrationRequest, err := authService.RegisterClient(req)
if err != nil {
statusCode := http.StatusInternalServerError
if err.Error() == "já existe uma solicitação pendente com este NIF" ||
   err.Error() == "já existe uma solicitação pendente com este email" ||
   err.Error() == "já existe uma conta aprovada com este NIF" ||
   err.Error() == "já existe uma conta aprovada com este email" {
statusCode = http.StatusConflict
}

c.JSON(statusCode, models.ErrorResponse{
Success: false,
Error:   err.Error(),
})
return
}

c.JSON(http.StatusCreated, models.SuccessResponse{
Success: true,
Message: "Solicitação enviada com sucesso. Aguarde aprovação da contabilista.",
Data:    registrationRequest,
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
c.JSON(http.StatusBadRequest, models.ErrorResponse{
Success: false,
Error:   err.Error(),
})
return
}

response, err := authService.Register(req)
if err != nil {
statusCode := http.StatusInternalServerError
if err.Error() == "username já está em uso" ||
   err.Error() == "email já está em uso" ||
   err.Error() == "NIF já está em uso" {
statusCode = http.StatusConflict
}

c.JSON(statusCode, models.ErrorResponse{
Success: false,
Error:   err.Error(),
})
return
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
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response, err := authService.LoginWithCredentials(req)
	if err != nil {
		statusCode := http.StatusUnauthorized
		if err.Error() == "utilizador não encontrado" {
			statusCode = http.StatusNotFound
		}
		
		c.JSON(statusCode, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
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
c.JSON(http.StatusOK, gin.H{
"success": true,
"message": "Logout realizado com sucesso",
})
}
