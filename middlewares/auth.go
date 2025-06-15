package middlewares

import (
	"RVContabilidadeBack/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Obter token do header Authorization
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de acesso requerido"})
            c.Abort()
            return
        }

        // Formato: "Bearer seu-token-aqui"
        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        
        // Validar token
        claims, err := utils.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
            c.Abort()
            return
        }

        // Guardar dados do utilizador no contexto
        c.Set("user_id", claims.UserID)
        c.Set("user_email", claims.Email)
        c.Set("user_username", claims.Username)
        c.Set("user_role", claims.Role)
        
        c.Next() // Continuar para o próximo handler
    }
}

// Middleware apenas para admins
func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("user_role")
        if !exists || role != "admin" {
            c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado - Admin requerido"})
            c.Abort()
            return
        }
        c.Next()
    }
}
