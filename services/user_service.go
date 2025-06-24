package services

import (
	"RVContabilidadeBack/config"
	"RVContabilidadeBack/models"
	"errors"
	"time"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// GetProfile obtém o perfil do utilizador
func (s *UserService) GetProfile(userID uint) (*models.User, error) {
	var user models.User
	if err := config.DB.Preload("Company").First(&user, userID).Error; err != nil {
		return nil, errors.New("utilizador não encontrado")
	}
	return &user, nil
}

// UpdateProfile atualiza o perfil do utilizador
func (s *UserService) UpdateProfile(userID uint, req models.UpdateProfileDTO) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("utilizador não encontrado")
	}

	// Atualizar apenas campos permitidos
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return nil, errors.New("erro ao atualizar perfil")
	}

	return &user, nil
}

// GetRequestHistory obtém o histórico de solicitações do utilizador
func (s *UserService) GetRequestHistory(userID uint) (map[string]interface{}, error) {
	// Buscar dados do utilizador
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("erro ao obter dados do utilizador")
	}

	// Buscar solicitação de registo
	var registrationRequest models.RegistrationRequest
	if err := config.DB.Preload("ReviewedByUser").Where("user_id = ?", user.ID).First(&registrationRequest).Error; err != nil {
		// Utilizador criado diretamente
		return map[string]interface{}{
			"current_status": user.Status,
			"created_at":     user.CreatedAt,
			"updated_at":     user.UpdatedAt,
		}, nil
	}

	response := map[string]interface{}{
		"current_status": user.Status,
		"created_at":     user.CreatedAt,
		"updated_at":     user.UpdatedAt,
		"submitted_at":   registrationRequest.SubmittedAt,
		"request_status": registrationRequest.Status,
	}

	// Adicionar informações de revisão se existirem
	if registrationRequest.ReviewedAt != nil {
		response["reviewed_at"] = registrationRequest.ReviewedAt
		response["review_notes"] = registrationRequest.ReviewNotes
		if registrationRequest.ReviewedByUser != nil {
			response["reviewed_by"] = registrationRequest.ReviewedByUser.Name
		}
	}

	return response, nil
}

// GetUserRequestHistory obtém o histórico de solicitações do utilizador
func (s *UserService) GetUserRequestHistory(userID uint) (map[string]interface{}, error) {
	// Buscar dados do utilizador
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("utilizador não encontrado")
	}

	// Criar resposta com histórico de status
	response := map[string]interface{}{
		"current_status": user.Status,
		"created_at":     user.CreatedAt,
		"updated_at":     user.UpdatedAt,
	}

	// Buscar solicitação de registo do utilizador
	var registrationRequest models.RegistrationRequest
	if err := config.DB.Preload("ReviewedByUser").Where("user_id = ?", user.ID).First(&registrationRequest).Error; err != nil {
		// Se não encontrar solicitação, pode ser utilizador criado diretamente
		return response, nil
	}

	// Adicionar dados da solicitação
	response["submitted_at"] = registrationRequest.SubmittedAt
	response["request_status"] = registrationRequest.Status

	// Adicionar informações de revisão se existirem
	if registrationRequest.ReviewedAt != nil {
		response["reviewed_at"] = registrationRequest.ReviewedAt
		response["review_notes"] = registrationRequest.ReviewNotes
		if registrationRequest.ReviewedByUser != nil {
			response["reviewed_by"] = registrationRequest.ReviewedByUser.Name
		}
	}

	return response, nil
}

// UpdateStatus atualiza o status de um utilizador (apenas admin)
func (s *UserService) UpdateStatus(userID uint, req models.UpdateUserStatusDTO) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("utilizador não encontrado")
	}

	user.Status = req.Status

	if err := config.DB.Save(&user).Error; err != nil {
		return nil, errors.New("erro ao atualizar status")
	}

	return &user, nil
}

// CompleteUserData completa dados pessoais após aprovação
func (s *UserService) CompleteUserData(userID uint, req models.CompleteUserDataDTO) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("utilizador não encontrado")
	}

	// Atualizar campos
	user.MaritalStatus = req.MaritalStatus
	user.CitizenCardNumber = req.CitizenCardNumber
	user.FixedPhone = req.FixedPhone
	user.FiscalCounty = req.FiscalCounty
	user.FiscalDistrict = req.FiscalDistrict
	user.OfficialEmail = req.OfficialEmail
	user.BillingSoftware = req.BillingSoftware
	user.PreferredFormat = req.PreferredFormat
	user.PreferredContactHours = req.PreferredContactHours

	// Processar data de expiração do cartão de cidadão
	if req.CitizenCardExpiry != "" {
		if expiry, err := time.Parse("2006-01-02", req.CitizenCardExpiry); err == nil {
			user.CitizenCardExpiry = &expiry
		}
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return nil, errors.New("erro ao completar dados do utilizador")
	}

	return &user, nil
}



// GetAllClients obtém todos os clientes (para admins/contabilistas)
func (s *UserService) GetAllClients() ([]models.User, error) {
	var users []models.User
	if err := config.DB.Where("role = ?", "client").Preload("Company").Find(&users).Error; err != nil {
		return nil, errors.New("erro ao obter lista de clientes")
	}
	return users, nil
}

// GetClientByID obtém um cliente específico por ID
func (s *UserService) GetClientByID(clientID uint) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("role = ?", "client").Preload("Company").First(&user, clientID).Error; err != nil {
		return nil, errors.New("cliente não encontrado")
	}
	return &user, nil
}

// UpdateUserStatus atualiza o status de um utilizador (para admins)
func (s *UserService) UpdateUserStatus(userID uint, req models.UpdateUserStatusDTO) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("utilizador não encontrado")
	}

	user.Status = req.Status
	
	if err := config.DB.Save(&user).Error; err != nil {
		return nil, errors.New("erro ao atualizar status do utilizador")
	}

	return &user, nil
}

// AdminUpdateClient permite admin/contabilista editar dados de cliente
func (s *UserService) AdminUpdateClient(clientID uint, req models.AdminUpdateClientDTO) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("role = ?", "client").First(&user, clientID).Error; err != nil {
		return nil, errors.New("cliente não encontrado")
	}

	// Atualizar campos permitidos para admin
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.NIF != nil {
		user.NIF = *req.NIF
	}
	if req.MaritalStatus != nil {
		user.MaritalStatus = *req.MaritalStatus
	}
	if req.CitizenCardNumber != nil {
		user.CitizenCardNumber = *req.CitizenCardNumber
	}
	if req.FixedPhone != nil {
		user.FixedPhone = *req.FixedPhone
	}
	if req.FiscalAddress != nil {
		user.FiscalAddress = *req.FiscalAddress
	}
	if req.FiscalPostalCode != nil {
		user.FiscalPostalCode = *req.FiscalPostalCode
	}
	if req.FiscalCity != nil {
		user.FiscalCity = *req.FiscalCity
	}
	if req.FiscalCounty != nil {
		user.FiscalCounty = *req.FiscalCounty
	}
	if req.FiscalDistrict != nil {
		user.FiscalDistrict = *req.FiscalDistrict
	}
	if req.OfficialEmail != nil {
		user.OfficialEmail = *req.OfficialEmail
	}
	if req.BillingSoftware != nil {
		user.BillingSoftware = *req.BillingSoftware
	}
	if req.PreferredFormat != nil {
		user.PreferredFormat = *req.PreferredFormat
	}
	if req.PreferredContactHours != nil {
		user.PreferredContactHours = *req.PreferredContactHours
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	// Parse dates if provided
	if req.DateOfBirth != nil {
		if parsed, err := time.Parse("2006-01-02", *req.DateOfBirth); err == nil {
			user.DateOfBirth = &parsed
		}
	}
	if req.CitizenCardExpiry != nil {
		if parsed, err := time.Parse("2006-01-02", *req.CitizenCardExpiry); err == nil {
			user.CitizenCardExpiry = &parsed
		}
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return nil, errors.New("erro ao atualizar dados do cliente")
	}

	// Recarregar com relacionamentos
	if err := config.DB.Preload("Company").First(&user, user.ID).Error; err != nil {
		return nil, errors.New("erro ao recarregar dados do cliente")
	}

	return &user, nil
}

// DeleteClient elimina um cliente
func (s *UserService) DeleteClient(clientID uint) error {
	// Verificar se o cliente existe e é cliente
	var client models.User
	if err := config.DB.Where("id = ? AND role = ?", clientID, "client").First(&client).Error; err != nil {
		return errors.New("cliente não encontrado")
	}

	// Iniciar transação
	tx := config.DB.Begin()

	// Eliminar empresa(s) do cliente
	if err := tx.Where("user_id = ?", clientID).Delete(&models.Company{}).Error; err != nil {
		tx.Rollback()
		return errors.New("erro ao eliminar empresa do cliente")
	}

	// Eliminar cliente
	if err := tx.Delete(&client).Error; err != nil {
		tx.Rollback()
		return errors.New("erro ao eliminar cliente")
	}

	tx.Commit()
	return nil
}
