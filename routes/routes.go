package routes

import (
	"RVContabilidadeBack/handlers"
	"RVContabilidadeBack/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
    router := gin.Default()

    // Middlewares globais
    router.Use(middlewares.LoggingMiddleware())
    router.Use(middlewares.CORSMiddleware())

    // API
    api := router.Group("/api")
    {
        // Rotas de autenticação (públicas)
        auth := api.Group("/auth")
        {
            auth.POST("/register", handlers.RegisterClient)      // Novo endpoint principal
            auth.POST("/register-direct", handlers.Register)     // Registo direto (interno)
            auth.POST("/login", handlers.Login)
            // Logout (protegida - requer token)
            auth.POST("/logout", middlewares.AuthMiddleware(), handlers.Logout)
        }

        // Rotas protegidas gerais (todos os utilizadores autenticados)
        protected := api.Group("/")
        protected.Use(middlewares.AuthMiddleware())
        {
            protected.GET("/profile", handlers.GetProfile)
        }

        // Rotas para administração (contabilistas e admins)
        admin := api.Group("/admin")
        admin.Use(middlewares.AuthMiddleware())
        admin.Use(middlewares.RequireRole("accountant", "admin"))
        {
            // Gestão de solicitações
            admin.GET("/pending-requests", handlers.GetPendingRequests)
            admin.POST("/approve-request", handlers.ApproveRequest)
            admin.GET("/requests", handlers.GetAllRequests)
            admin.GET("/requests/:id", handlers.GetRequestDetails)
            
            // Gestão de utilizadores
            admin.GET("/users", handlers.GetAllUsers)
            admin.GET("/users/:id", handlers.GetUserDetails)
        }

        // Rotas apenas para admins
        adminOnly := api.Group("/admin")
        adminOnly.Use(middlewares.AuthMiddleware())
        adminOnly.Use(middlewares.RequireRole("admin"))
        {
            adminOnly.PUT("/users/:id/status", handlers.UpdateUserStatus)
        }

        // Rotas para clientes (apenas clientes aprovados)
        client := api.Group("/client")
        client.Use(middlewares.AuthMiddleware())
        client.Use(middlewares.RequireRole("client"))
        {
            client.GET("/profile", handlers.GetClientProfile)
            client.PUT("/profile", handlers.UpdateClientProfile)
            client.GET("/company", handlers.GetClientCompany)
            client.PUT("/company", handlers.UpdateClientCompany)
            client.GET("/requests", handlers.GetClientRequests)
            
            // Novos endpoints para completar dados após aprovação
            client.POST("/complete-user-data", handlers.CompleteUserData)
            client.POST("/complete-company-data", handlers.CompleteCompanyData)
        }

        // Rota de informações da API (pública)
        api.GET("/info", handlers.GetAPIInfo)
    }

    // Rota de health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "OK", "message": "API está funcionando"})
    })

    return router
}

