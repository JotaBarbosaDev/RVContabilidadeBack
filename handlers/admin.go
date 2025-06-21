package handlers

import (
	"RVContabilidadeBack/config"
	"RVContabilidadeBack/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetPendingRequests godoc
// @Summary      Listar solicitações pendentes
// @Description  Lista todas as solicitações de registo pendentes (apenas contabilistas/admin)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SuccessResponse
// @Router       /admin/pending-requests [get]
func GetPendingRequests(c *gin.Context) {
	var requests []models.RegistrationRequest
	if err := config.DB.Preload("User").Where("status = ?", "pending").Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Erro ao obter solicitações",
		})
		return
	}

	// Converter para DTO
	var response []models.PendingRequestResponseDTO
	for _, req := range requests {
		var requestData models.RegistrationRequestDTO
		if err := json.Unmarshal([]byte(req.RequestData), &requestData); err != nil {
			// Se não conseguir fazer unmarshal, criar um requestData vazio
			requestData = models.RegistrationRequestDTO{}
		}
		
		response = append(response, models.PendingRequestResponseDTO{
			ID:          req.ID,
			User:        req.User,
			RequestType: req.RequestType,
			RequestData: requestData,
			SubmittedAt: req.SubmittedAt,
			Status:      req.Status,
		})
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Solicitações obtidas com sucesso",
		Data:    response,
	})
}

// ApproveRequest godoc
// @Summary      Aprovar/rejeitar solicitação
// @Description  Aprova ou rejeita uma solicitação de registo (apenas contabilistas/admin)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.ApprovalRequestDTO  true  "Dados de aprovação"
// @Success      200      {object}  models.SuccessResponse
// @Router       /admin/approve-request [post]
func ApproveRequest(c *gin.Context) {
	reviewerID, _ := c.Get("user_id")

	var req models.ApprovalRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Buscar solicitação
	var registrationRequest models.RegistrationRequest
	if err := config.DB.Preload("User").First(&registrationRequest, req.RequestID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "Solicitação não encontrada",
		})
		return
	}

	// Verificar se ainda está pendente
	if registrationRequest.Status != "pending" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Solicitação já foi processada",
		})
		return
	}

	// Começar transação
	tx := config.DB.Begin()

	// Atualizar solicitação
	now := time.Now()
	reviewerIDUint, ok := reviewerID.(uint)
	if !ok {
		// Tentar converter de float64 para uint (caso venha do JWT como número)
		if f, ok := reviewerID.(float64); ok {
			reviewerIDUint = uint(f)
		} else {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   "Erro ao obter ID do revisor",
			})
			return
		}
	}
	
	registrationRequest.Status = req.Status
	registrationRequest.ReviewedAt = &now
	registrationRequest.ReviewedBy = &reviewerIDUint
	registrationRequest.ReviewNotes = req.ReviewNotes

	if err := tx.Save(&registrationRequest).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Erro ao atualizar solicitação",
		})
		return
	}

	// Atualizar utilizador
	user := registrationRequest.User
	if req.Status == "approved" {
		user.Status = string(models.StatusApproved)
	} else {
		user.Status = string(models.StatusRejected)
	}

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Erro ao atualizar utilizador",
		})
		return
	}

	tx.Commit()

	message := "Solicitação aprovada com sucesso"
	if req.Status == "rejected" {
		message = "Solicitação rejeitada"
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: message,
		Data:    registrationRequest,
	})
}

// GetAllRequests godoc
// @Summary      Histórico de solicitações
// @Description  Lista todas as solicitações (apenas contabilistas/admin)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SuccessResponse
// @Router       /admin/requests [get]
func GetAllRequests(c *gin.Context) {
	var requests []models.RegistrationRequest
	if err := config.DB.Preload("User").Preload("ReviewedByUser").Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Erro ao obter solicitações",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Histórico de solicitações obtido com sucesso",
		Data:    requests,
	})
}

// GetRequestDetails godoc
// @Summary      Detalhes de uma solicitação
// @Description  Obtém detalhes completos de uma solicitação específica
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID da solicitação"
// @Success      200  {object}  models.SuccessResponse
// @Router       /admin/requests/{id} [get]
func GetRequestDetails(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "ID da solicitação inválido",
		})
		return
	}

	var request models.RegistrationRequest
	if err := config.DB.Preload("User").Preload("ReviewedByUser").First(&request, requestID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "Solicitação não encontrada",
		})
		return
	}

	// Converter dados da solicitação
	var requestData models.RegistrationRequestDTO
	json.Unmarshal([]byte(request.RequestData), &requestData)

	response := gin.H{
		"request":      request,
		"request_data": requestData,
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Detalhes da solicitação obtidos com sucesso",
		Data:    response,
	})
}

// UpdateUserStatus godoc
// @Summary      Alterar status de utilizador
// @Description  Bloqueia/desbloqueia utilizador (apenas admin)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int  true  "ID do utilizador"
// @Param        request body      models.UpdateUserStatusDTO  true  "Novo status"
// @Success      200     {object}  models.SuccessResponse
// @Router       /admin/users/{id}/status [put]
func UpdateUserStatus(c *gin.Context) {
	userRole, _ := c.Get("user_role")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Success: false,
			Error:   "Apenas admins podem alterar status",
		})
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "ID do utilizador inválido",
		})
		return
	}

	var req models.UpdateUserStatusDTO
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

	user.Status = req.Status
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Erro ao atualizar utilizador",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Status do utilizador atualizado com sucesso",
		Data:    user,
	})
}

// GetAllUsers godoc
// @Summary      Listar todos os utilizadores
// @Description  Lista todos os utilizadores do sistema (apenas admin/contabilista)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SuccessResponse
// @Router       /admin/users [get]
func GetAllUsers(c *gin.Context) {
	var users []models.User
	
	if err := config.DB.Preload("Companies").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Erro ao obter utilizadores",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Utilizadores obtidos com sucesso",
		Data: gin.H{
			"users": users,
			"total": len(users),
		},
	})
}

// GetUserDetails godoc
// @Summary      Detalhes de um utilizador
// @Description  Obtém detalhes completos de um utilizador
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID do utilizador"
// @Success      200  {object}  models.SuccessResponse
// @Router       /admin/users/{id} [get]
func GetUserDetails(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "ID do utilizador inválido",
		})
		return
	}

	var user models.User
	if err := config.DB.Preload("Companies").Preload("Requests").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "Utilizador não encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Detalhes do utilizador obtidos com sucesso",
		Data:    user,
	})
}
