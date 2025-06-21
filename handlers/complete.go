package handlers

import (
	"RVContabilidadeBack/config"
	"RVContabilidadeBack/models"
	"RVContabilidadeBack/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	// Buscar utilizador
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "Utilizador não encontrado",
		})
		return
	}

	// Verificar se está aprovado
	if user.Status != "approved" {
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Success: false,
			Error:   "Conta ainda não aprovada",
		})
		return
	}

	// Atualizar dados opcionais
	if dto.MaritalStatus != "" {
		user.MaritalStatus = dto.MaritalStatus
	}
	if dto.CitizenCardNumber != "" {
		user.CitizenCardNumber = dto.CitizenCardNumber
	}
	if dto.CitizenCardExpiry != "" {
		if expiry, err := time.Parse("2006-01-02", dto.CitizenCardExpiry); err == nil {
			user.CitizenCardExpiry = &expiry
		}
	}
	if dto.FixedPhone != "" {
		user.FixedPhone = dto.FixedPhone
	}
	if dto.FiscalCounty != "" {
		user.FiscalCounty = dto.FiscalCounty
	}
	if dto.FiscalDistrict != "" {
		user.FiscalDistrict = dto.FiscalDistrict
	}
	
	// Credenciais encriptadas
	if dto.PortalFinancasUser != "" {
		user.PortalFinancasUser = dto.PortalFinancasUser
		if encrypted, err := utils.Encrypt(dto.PortalFinancasPassword); err == nil {
			user.PortalFinancasPassword = encrypted
		}
	}
	if dto.EFaturaUser != "" {
		user.EFaturaUser = dto.EFaturaUser
		if encrypted, err := utils.Encrypt(dto.EFaturaPassword); err == nil {
			user.EFaturaPassword = encrypted
		}
	}
	if dto.SSDirectUser != "" {
		user.SSDirectUser = dto.SSDirectUser
		if encrypted, err := utils.Encrypt(dto.SSDirectPassword); err == nil {
			user.SSDirectPassword = encrypted
		}
	}
	
	if dto.OfficialEmail != "" {
		user.OfficialEmail = dto.OfficialEmail
	}
	if dto.BillingSoftware != "" {
		user.BillingSoftware = dto.BillingSoftware
	}
	if dto.PreferredFormat != "" {
		user.PreferredFormat = dto.PreferredFormat
	}
	if dto.PreferredContactHours != "" {
		user.PreferredContactHours = dto.PreferredContactHours
	}

	// Salvar na base de dados
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Erro ao atualizar dados",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Dados pessoais atualizados com sucesso",
		Data:    nil,
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

	// Buscar utilizador
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "Utilizador não encontrado",
		})
		return
	}

	// Verificar se está aprovado
	if user.Status != "approved" {
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Success: false,
			Error:   "Conta ainda não aprovada",
		})
		return
	}

	// Buscar empresa
	var company models.Company
	if err := config.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "Empresa não encontrada",
		})
		return
	}

	// Atualizar dados opcionais da empresa
	if dto.TradeName != "" {
		company.TradeName = dto.TradeName
	}
	if dto.CorporateObject != "" {
		company.CorporateObject = dto.CorporateObject
	}
	if dto.Address != "" {
		company.Address = dto.Address
	}
	if dto.PostalCode != "" {
		company.PostalCode = dto.PostalCode
	}
	if dto.City != "" {
		company.City = dto.City
	}
	if dto.County != "" {
		company.County = dto.County
	}
	if dto.District != "" {
		company.District = dto.District
	}
	if dto.ShareCapital > 0 {
		company.ShareCapital = dto.ShareCapital
	}
	if dto.GroupStartDate != "" {
		if startDate, err := time.Parse("2006-01-02", dto.GroupStartDate); err == nil {
			company.GroupStartDate = &startDate
		}
	}
	if dto.BankName != "" {
		company.BankName = dto.BankName
	}
	if dto.IBAN != "" {
		company.IBAN = dto.IBAN
	}
	if dto.BIC != "" {
		company.BIC = dto.BIC
	}
	if dto.AnnualRevenue > 0 {
		company.AnnualRevenue = dto.AnnualRevenue
	}
	// HasStock é boolean, sempre atualiza
	company.HasStock = dto.HasStock
	if dto.MainClients != "" {
		company.MainClients = dto.MainClients
	}
	if dto.MainSuppliers != "" {
		company.MainSuppliers = dto.MainSuppliers
	}

	// Salvar na base de dados
	if err := config.DB.Save(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Erro ao atualizar dados da empresa",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Dados da empresa atualizados com sucesso",
		Data:    nil,
	})
}
