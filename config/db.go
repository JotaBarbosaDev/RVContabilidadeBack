package config

import (
	"RVContabilidadeBack/models"
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	_ = godotenv.Load()

	dsn := "host=localhost user=postgres password=postgres dbname=RVContabilidadeDB port=5432 sslmode=disable client_encoding=UTF8"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("❌ Erro ao ligar à base de dados: " + err.Error())
	}

	DB = db
	fmt.Println("✅ Ligação à BD estabelecida")
	
	// Auto-migrate models
	migrate()
}

func migrate() {
	// Primeiro, criar tabelas automaticamente baseadas nos models
	err := DB.AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.RegistrationRequest{},
	)
	if err != nil {
		fmt.Printf("❌ Erro na migração: %v\n", err)
		return
	}
	
	fmt.Println("✅ Migração das tabelas concluída")
	
	// Depois, fazer migração personalizada do username (se necessário)
	if err := migrateUsernameField(); err != nil {
		fmt.Printf("❌ Erro na migração personalizada do username: %v\n", err)
		// Não retornar aqui para permitir que continue
	}
	
	// Criar utilizador admin se não existir
	createDefaultAdmin()
}

func createDefaultAdmin() {
	var adminUser models.User
	if err := DB.Where("role = ? AND username = ?", "admin", "admin").First(&adminUser).Error; err != nil {
		// Admin não existe, criar
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		
		admin := models.User{
			Username: "admin",
			Name:     "Administrador",
			Email:    "admin@rvcontabilidade.com",
			Phone:    "000000000",
			NIF:      "000000000",
			Password: string(hashedPassword),
			Role:     "admin",
			Status:   string(models.StatusApproved),
		}
		
		if err := DB.Create(&admin).Error; err != nil {
			fmt.Printf("❌ Erro ao criar admin: %v\n", err)
		} else {
			fmt.Println("✅ Utilizador admin criado: admin / admin123")
		}
	}
	
	// Criar utilizador contabilista se não existir
	var accountantUser models.User
	if err := DB.Where("role = ? AND username = ?", "accountant", "contabilista").First(&accountantUser).Error; err != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("contabilista123"), bcrypt.DefaultCost)
		
		accountant := models.User{
			Username: "contabilista",
			Name:     "Contabilista",
			Email:    "contabilista@rvcontabilidade.com",
			Phone:    "111111111",
			NIF:      "111111111",
			Password: string(hashedPassword),
			Role:     "accountant",
			Status:   string(models.StatusApproved),
		}
		
		if err := DB.Create(&accountant).Error; err != nil {
			fmt.Printf("❌ Erro ao criar contabilista: %v\n", err)
		} else {
			fmt.Println("✅ Utilizador contabilista criado: contabilista / contabilista123")
		}
	}
}

// migrateUsernameField faz a migração do campo username de forma segura
func migrateUsernameField() error {
	// Como o GORM já criou a tabela, só precisamos verificar se há utilizadores sem username
	var usersWithoutUsername []models.User
	err := DB.Where("username IS NULL OR username = ''").Find(&usersWithoutUsername).Error
	if err != nil {
		return err
	}
	
	if len(usersWithoutUsername) > 0 {
		fmt.Printf("🔄 Migrando %d utilizadores sem username...\n", len(usersWithoutUsername))
		
		for _, user := range usersWithoutUsername {
			// Gerar username baseado no email
			username := strings.Split(user.Email, "@")[0]
			
			// Verificar se já existe, se sim adicionar número
			var count int64
			DB.Model(&models.User{}).Where("username = ?", username).Count(&count)
			if count > 0 {
				username = fmt.Sprintf("%s_%d", username, user.ID)
			}
			
			// Atualizar o utilizador
			DB.Model(&user).Update("username", username)
		}
		
		fmt.Println("✅ Migração de usernames concluída")
	}
	
	return nil
}
