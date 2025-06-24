package controllers

import (
	"RVContabilidadeBack/models"
	"RVContabilidadeBack/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	userService = services.NewUserService()
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

	user, err := userService.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Perfil obtido com sucesso",
		Data:    user,
	})
}
