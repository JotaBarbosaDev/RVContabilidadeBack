package controllers

import (
	"RVContabilidadeBack/models"
	"RVContabilidadeBack/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	completeUserService    = services.NewUserService()
	completeCompanyService = services.NewCompanyService()
)

// CompleteUserData godoc
// @Summary      Completar dados pessoais
// @Description  Completa dados pessoais adicionais após aprovação
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        data  body      models.CompleteUserDataDTO  true  "Dados pessoais completos"
// @Success      200   {object}  models.SuccessResponse
// @Router       /client/complete-user-data [put]
func CompleteUserData(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Token inválido",
		})
		return
	}

	var dto models.CompleteUserDataDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Dados inválidos: " + err.Error(),
		})
		return
	}

	user, err := completeUserService.CompleteUserData(userID.(uint), dto)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "utilizador não encontrado" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "conta ainda não aprovada" {
			statusCode = http.StatusForbidden
		}
		
		c.JSON(statusCode, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Dados pessoais atualizados com sucesso",
		Data:    user,
	})
}

// CompleteCompanyData godoc
// @Summary      Completar dados da empresa
// @Description  Completa dados da empresa adicionais após aprovação
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        data  body      models.CompleteCompanyDataDTO  true  "Dados da empresa completos"
// @Success      200   {object}  models.SuccessResponse
// @Router       /client/complete-company-data [put]
func CompleteCompanyData(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Token inválido",
		})
		return
	}

	var dto models.CompleteCompanyDataDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Dados inválidos: " + err.Error(),
		})
		return
	}

	company, err := completeCompanyService.CompleteCompanyData(userID.(uint), dto)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "utilizador não encontrado" || err.Error() == "empresa não encontrada" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "conta ainda não aprovada" {
			statusCode = http.StatusForbidden
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
