package handlers

import (
	"RVContabilidadeBack/config"
	"RVContabilidadeBack/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProfile godoc
// @Summary      Meu perfil
// @Description  Obtém dados do utilizador logado (requer token no header Authorization)
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SuccessResponse
// @Router       /profile [get]
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Utilizador não autenticado",
		})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "Utilizador não encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Perfil obtido com sucesso",
		Data:    user,
	})
}

// GetAllUsers godoc
// @Summary      Listar utilizadores
// @Description  Lista todos os utilizadores (só admins - requer token de admin no header Authorization)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SuccessResponse
// @Router       /admin/users [get]
func GetAllUsers(c *gin.Context) {
	var users []models.User
	
	// Obter todos os utilizadores da base de dados
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Erro ao obter utilizadores",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Utilizadores obtidos com sucesso",
		Data: map[string]interface{}{
			"users": users,
			"total": len(users),
		},
	})
}
