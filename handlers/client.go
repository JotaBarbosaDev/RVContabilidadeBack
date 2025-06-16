package handlers

import (
	"RVContabilidadeBack/config"
	"RVContabilidadeBack/models"
	"net/http"

	"github.com/gin-gonic/gin"
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

	var user models.User
	if err := config.DB.Preload("Companies").First(&user, userID).Error; err != nil {
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

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "Utilizador não encontrado",
		})
		return
	}

	// Atualizar apenas campos permitidos
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	// Email e NIF não podem ser alterados pelo cliente

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Erro ao atualizar perfil",
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

	var company models.Company
	if err := config.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "Empresa não encontrada",
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

	var company models.Company
	if err := config.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "Empresa não encontrada",
		})
		return
	}

	// Atualizar apenas campos permitidos (cliente não pode alterar NIPC, nome oficial, etc.)
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
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Erro ao atualizar empresa",
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

	var requests []models.RegistrationRequest
	if err := config.DB.Where("user_id = ?", userID).Preload("ReviewedByUser").Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Erro ao obter histórico",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Histórico de solicitações obtido com sucesso",
		Data:    requests,
	})
}
