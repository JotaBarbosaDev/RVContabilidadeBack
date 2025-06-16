package config

import (
	"RVContabilidadeBack/models"
	"fmt"

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
	// Primeiro, fazer migração personalizada do username para utilizadores existentes
	if err := migrateUsernameField(); err != nil {
		fmt.Printf("❌ Erro na migração personalizada do username: %v\n", err)
		return
	}
	
	// Criar tabelas automaticamente baseadas nos models
	err := DB.AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.RegistrationRequest{},
	)
	if err != nil {
		fmt.Printf("❌ Erro na migração: %v\n", err)
	} else {
		fmt.Println("✅ Migração das tabelas concluída")
		
		// Criar utilizador admin se não existir
		createDefaultAdmin()
	}
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
	// Verificar se a coluna username já existe
	var count int64
	err := DB.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'username'").Scan(&count).Error
	if err != nil {
		return err
	}
	
	if count == 0 {
		// Coluna não existe, vamos criá-la
		fmt.Println("🔄 Adicionando coluna username...")
		
		// 1. Adicionar coluna sem NOT NULL
		if err := DB.Exec("ALTER TABLE users ADD COLUMN username VARCHAR(255)").Error; err != nil {
			return err
		}
		
		// 2. Atualizar registos existentes com valores baseados no email
		if err := DB.Exec("UPDATE users SET username = SPLIT_PART(email, '@', 1) WHERE username IS NULL").Error; err != nil {
			return err
		}
		
		// 3. Tornar a coluna NOT NULL e UNIQUE
		if err := DB.Exec("ALTER TABLE users ALTER COLUMN username SET NOT NULL").Error; err != nil {
			return err
		}
		
		if err := DB.Exec("ALTER TABLE users ADD CONSTRAINT users_username_key UNIQUE (username)").Error; err != nil {
			// Se já existe, ignorar erro
			fmt.Println("⚠️  Constraint UNIQUE já existe para username")
		}
		
		fmt.Println("✅ Coluna username adicionada e migrada com sucesso")
	}
	
	return nil
}
