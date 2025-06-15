package config

import (
	"RVContabilidadeBack/models"
	"fmt"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
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
	// Criar tabelas automaticamente baseadas nos models
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Printf("❌ Erro na migração: %v\n", err)
	} else {
		fmt.Println("✅ Migração das tabelas concluída")
	}
}
