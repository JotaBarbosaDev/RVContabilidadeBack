package controllers

import (
	"RVContabilidadeBack/models"
	"RVContabilidadeBack/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	adminService = services.NewAdminService()
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
	requests, err := adminService.GetPendingRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Solicitações pendentes obtidas com sucesso",
		Data:    requests,
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

	// Converter reviewerID para uint
	reviewerIDUint, ok := reviewerID.(uint)
	if !ok {
		// Tentar converter de float64 para uint (caso venha do JWT como número)
		if f, ok := reviewerID.(float64); ok {
			reviewerIDUint = uint(f)
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   "Erro ao obter ID do revisor",
			})
			return
		}
	}

	request, err := adminService.ApproveRequest(req, reviewerIDUint)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "solicitação não encontrada" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "solicitação já foi processada" {
			statusCode = http.StatusBadRequest
		} else if err.Error() == "já existe uma empresa com este NIPC" {
			statusCode = http.StatusConflict
		}
		
		c.JSON(statusCode, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	message := "Solicitação aprovada com sucesso"
	if req.Status == "rejected" {
		message = "Solicitação rejeitada"
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: message,
		Data:    request,
	})
}

// GetAllRequests godoc
// @Summary      Listar pedidos de registo pendentes
// @Description  Lista todos os pedidos de registo pendentes para aprovação
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SuccessResponse
// @Router       /admin/requests [get]
func GetAllRequests(c *gin.Context) {
	requests, err := adminService.GetAllRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Lista de pedidos de registo obtida com sucesso",
		Data:    requests,
	})
}

// GetRequestDetails godoc
// @Summary      Detalhes de um pedido de registo
// @Description  Obtém detalhes completos de um pedido de registo específico
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID do pedido"
// @Success      200  {object}  models.SuccessResponse
// @Router       /admin/requests/{id} [get]
func GetRequestDetails(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "ID do pedido inválido",
		})
		return
	}

	request, err := adminService.GetRequestDetails(uint(requestID))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Detalhes do pedido de registo obtidos com sucesso",
		Data:    request,
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

	user, err := adminService.UpdateUserStatus(uint(userID), req.Status)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "utilizador não encontrado" {
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
	status := c.Query("status")
	role := c.Query("role")
	
	users, stats, err := adminService.GetAllUsers(status, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Utilizadores obtidos com sucesso",
		Data: gin.H{
			"users": users,
			"total": len(users),
			"stats": stats,
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

	user, err := adminService.GetUserDetails(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Detalhes do utilizador obtidos com sucesso",
		Data:    user,
	})
}

// GetApprovedClients godoc
// @Summary      Listar clientes aprovados
// @Description  Lista todos os clientes com status aprovado (apenas contabilista/admin)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SuccessResponse
// @Router       /admin/clients [get]
func GetApprovedClients(c *gin.Context) {
	clients, err := adminService.GetApprovedClients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Clientes aprovados obtidos com sucesso",
		Data: gin.H{
			"clients": clients,
			"total":   len(clients),
		},
	})
}

// UpdateClientData godoc
// @Summary      Atualizar dados de cliente
// @Description  Atualiza dados pessoais de um cliente (apenas contabilista/admin)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                        true  "ID do cliente"
// @Param        request  body      models.AdminUpdateClientDTO true  "Dados para atualizar"
// @Success      200      {object}  models.SuccessResponse
// @Router       /admin/clients/{id} [put]
func UpdateClientData(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "ID do cliente inválido",
		})
		return
	}

	var req models.AdminUpdateClientDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Dados inválidos: " + err.Error(),
		})
		return
	}

	err = adminService.UpdateClientData(uint(clientID), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "cliente aprovado não encontrado" {
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
		Message: "Dados do cliente atualizados com sucesso",
		Data:    gin.H{"client_id": clientID},
	})
}

// AdminUpdateClientCompany godoc
// @Summary      Atualizar empresa do cliente
// @Description  Atualiza dados da empresa de um cliente (apenas contabilista/admin)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                           true  "ID do cliente"
// @Param        request  body      models.AdminUpdateCompanyDTO true  "Dados da empresa para atualizar"
// @Success      200      {object}  models.SuccessResponse
// @Router       /admin/clients/{id}/company [put]
func AdminUpdateClientCompany(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "ID do cliente inválido",
		})
		return
	}

	var req models.AdminUpdateCompanyDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Dados inválidos: " + err.Error(),
		})
		return
	}

	company, err := adminService.UpdateClientCompany(uint(clientID), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "cliente aprovado não encontrado" ||
		   err.Error() == "empresa do cliente não encontrada" {
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
		Data:    gin.H{"client_id": clientID, "company_id": company.ID},
	})
}

// DeleteClient godoc
// @Summary      Eliminar cliente
// @Description  Elimina um cliente e a sua empresa (apenas contabilista/admin)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID do cliente"
// @Success      200  {object}  models.SuccessResponse
// @Router       /admin/clients/{id} [delete]
func DeleteClient(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "ID do cliente inválido",
		})
		return
	}

	err = adminService.DeleteClient(uint(clientID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "cliente não encontrado" {
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
		Message: "Cliente eliminado com sucesso",
		Data:    gin.H{"deleted_client_id": clientID},
	})
}

// GetUsersCount godoc
// @Summary      Contar utilizadores
// @Description  Conta quantos utilizadores existem por status e role
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SuccessResponse
// @Router       /admin/users/count [get]
func GetUsersCount(c *gin.Context) {
	counts, err := adminService.GetUsersCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Contagem de utilizadores obtida com sucesso",
		Data:    counts,
	})
}

// GetAllUsersSimple godoc
// @Summary      Listar utilizadores (simples)
// @Description  Lista utilizadores sem relacionamentos para debug
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SuccessResponse
// @Router       /admin/users/simple [get]
func GetAllUsersSimple(c *gin.Context) {
	users, err := adminService.GetAllUsersSimple()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
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
