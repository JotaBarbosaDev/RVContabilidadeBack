package routes

import (
	"RVContabilidadeBack/handlers"
	"RVContabilidadeBack/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
    router := gin.Default()

    // Middlewares globais
    router.Use(middlewares.CORSMiddleware())

    // API
    api := router.Group("/api")
    {

        // Rotas de autenticação (públicas)
        auth := api.Group("/auth")
        {
            auth.POST("/register", handlers.Register)
            auth.POST("/login", handlers.Login)
            // Logout (protegida - requer token)
            auth.POST("/logout", middlewares.AuthMiddleware(), handlers.Logout)
        }

        // Rotas protegidas
        protected := api.Group("/")
        protected.Use(middlewares.AuthMiddleware())
        {
            protected.GET("/profile", handlers.GetProfile)
        }

        // Rotas apenas para admins
        admin := api.Group("/admin")
        admin.Use(middlewares.AuthMiddleware())
        admin.Use(middlewares.AdminMiddleware())
        {
            admin.GET("/users", handlers.GetAllUsers)
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

