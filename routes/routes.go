package routes

import (
	"RVContabilidadeBack/controllers"
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
            auth.POST("/register", controllers.RegisterClient)      // Novo endpoint principal
            auth.POST("/register-direct", controllers.Register)     // Registo direto (interno)
            auth.POST("/login", controllers.Login)
            // Logout (protegida - requer token)
            auth.POST("/logout", middlewares.AuthMiddleware(), controllers.Logout)
        }

        // Rotas protegidas gerais (todos os utilizadores autenticados)
        protected := api.Group("/")
        protected.Use(middlewares.AuthMiddleware())
        {
            protected.GET("/profile", controllers.GetProfile)
        }

        // Rotas para administração (contabilistas e admins)
        admin := api.Group("/admin")
        admin.Use(middlewares.AuthMiddleware())
        admin.Use(middlewares.RequireRole("accountant", "admin"))
        {
            // Dashboard
            admin.GET("/dashboard", controllers.GetDashboardData)
            
            // Gestão de solicitações
            admin.GET("/pending-requests", controllers.GetPendingRequests)
            admin.POST("/approve-request", controllers.ApproveRequest)
            admin.GET("/requests", controllers.GetAllRequests)
            admin.GET("/requests/:id", controllers.GetRequestDetails)
            
            // Gestão de utilizadores
            admin.GET("/users", controllers.GetAllUsers)
            admin.GET("/users/count", controllers.GetUsersCount)
            admin.GET("/users/simple", controllers.GetAllUsersSimple)
            admin.GET("/users/:id", controllers.GetUserDetails)
            
            // Gestão de clientes aprovados
            admin.GET("/clients", controllers.GetApprovedClients)
            admin.GET("/clients/overview", controllers.GetAllClientsOverview)
            admin.PUT("/clients/:id", controllers.UpdateClientData)
            admin.PUT("/clients/:id/company", controllers.AdminUpdateClientCompany) 
            admin.DELETE("/clients/:id", controllers.DeleteClient)
            
            // Visão completa de todos os clientes (combina users, registration_requests e companies)
            admin.GET("/complete-users-overview", controllers.GetCompleteUsersOverview)
        }

        // Rotas apenas para admins
        adminOnly := api.Group("/admin")
        adminOnly.Use(middlewares.AuthMiddleware())
        adminOnly.Use(middlewares.RequireRole("admin"))
        {
            adminOnly.PUT("/users/:id/status", controllers.UpdateUserStatus)
        }

        // Rotas para clientes (apenas clientes aprovados)
        client := api.Group("/client")
        client.Use(middlewares.AuthMiddleware())
        client.Use(middlewares.RequireRole("client"))
        {
            client.GET("/profile", controllers.GetClientProfile)
            client.PUT("/profile", controllers.UpdateClientProfile)
            client.GET("/company", controllers.GetClientCompany)
            client.PUT("/company", controllers.UpdateClientCompany)
            client.GET("/requests", controllers.GetClientRequests)
            
            // Novos endpoints para completar dados após aprovação
            client.POST("/complete-user-data", controllers.CompleteUserData)
            client.POST("/complete-company-data", controllers.CompleteCompanyData)
        }

        // Rota de informações da API (pública)
        api.GET("/info", controllers.GetAPIInfo)
    }

    // Rota de health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "OK", "message": "API está funcionando"})
    })

    return router
}

