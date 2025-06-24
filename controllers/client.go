package controllers

import (
	"RVContabilidadeBack/models"
	"RVContabilidadeBack/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	clientUserService    = services.NewUserService()
	clientCompanyService = services.NewCompanyService()
)

// GetClientProfile godoc
// @Summary      Perfil do cliente
// @Description  Obtém perfil do cliente logado
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SuccessResponse
// @Router       /client/profile [get]
func GetClientProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Utilizador não autenticado",
		})
		return
	}

	user, err := clientUserService.GetProfile(userID.(uint))
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

// UpdateClientProfile godoc
// @Summary      Atualizar perfil do cliente
// @Description  Atualiza dados pessoais do cliente
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        profile  body      models.UpdateProfileDTO  true  "Dados a atualizar"
// @Success      200      {object}  models.SuccessResponse
// @Router       /client/profile [put]
func UpdateClientProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Utilizador não autenticado",
		})
		return
	}

	var req models.UpdateProfileDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user, err := clientUserService.UpdateProfile(userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Perfil atualizado com sucesso",
		Data:    user,
	})
}

// GetClientCompany godoc
// @Summary      Dados da empresa do cliente
// @Description  Obtém dados da empresa do cliente logado
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SuccessResponse
// @Router       /client/company [get]
func GetClientCompany(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Utilizador não autenticado",
		})
		return
	}

	company, err := clientCompanyService.GetByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Dados da empresa obtidos com sucesso",
		Data:    company,
	})
}

// UpdateClientCompany godoc
// @Summary      Atualizar dados da empresa
// @Description  Atualiza dados da empresa do cliente (campos limitados)
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        company  body      models.UpdateCompanyDTO  true  "Dados a atualizar"
// @Success      200      {object}  models.SuccessResponse
// @Router       /client/company [put]
func UpdateClientCompany(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Utilizador não autenticado",
		})
		return
	}

	var req models.UpdateCompanyDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	company, err := clientCompanyService.UpdateCompany(userID.(uint), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "empresa não encontrada" {
			statusCode = http.StatusNotFound
		}
		
		c.JSON(statusCode, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Dados da empresa atualizados com sucesso",
		Data:    company,
	})
}

// GetClientRequests godoc
// @Summary      Histórico de solicitações do cliente
// @Description  Lista todas as solicitações do cliente logado
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SuccessResponse
// @Router       /client/requests [get]
func GetClientRequests(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Utilizador não autenticado",
		})
		return
	}

	history, err := clientUserService.GetUserRequestHistory(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Histórico do utilizador obtido com sucesso",
		Data:    history,
	})
}
