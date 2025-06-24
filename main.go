package main

import (
	"RVContabilidadeBack/config"
	_ "RVContabilidadeBack/docs" // Será gerado automaticamente
	"RVContabilidadeBack/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           RV Contabilidade API
// @version         2.0
// @description     Sistema de gestão contabilística com aprovação de clientes - Clean Architecture

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
    // Inicializar BD
    config.ConnectDatabase()

    // Configurar rotas
    router := routes.SetupRoutes()

    // Swagger endpoint
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Rota de documentação
    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "RV Contabilidade API",
            "version": "2.0",
            "docs":    "http://localhost:8080/swagger/index.html",
        })
    })

    router.Run(":8080")
}
