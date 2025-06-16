package middlewares

import (
	"RVContabilidadeBack/config"
	"RVContabilidadeBack/models"
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

        // Verificar status do utilizador na base de dados
        var user models.User
        if err := config.DB.First(&user, claims.UserID).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilizador não encontrado"})
            c.Abort()
            return
        }

        // Verificar se o utilizador está aprovado
        if user.Status != string(models.StatusApproved) {
            statusMessages := map[string]string{
                string(models.StatusPending):  "Conta aguarda aprovação da contabilista",
                string(models.StatusRejected): "Conta foi rejeitada. Contacte o suporte.",
                string(models.StatusBlocked):  "Conta foi bloqueada. Contacte o suporte.",
            }
            
            message := statusMessages[user.Status]
            if message == "" {
                message = "Acesso negado"
            }
            
            c.JSON(http.StatusForbidden, gin.H{"error": message})
            c.Abort()
            return
        }

        // Guardar dados do utilizador no contexto
        c.Set("user_id", claims.UserID)
        c.Set("user_username", claims.Username)
        c.Set("user_nif", claims.NIF)
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

// RequireRole middleware para verificar se o utilizador tem uma das roles permitidas
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role de utilizador não encontrada"})
			c.Abort()
			return
		}

		roleStr := userRole.(string)
		
		// Verificar se a role está nas permitidas
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Não tem permissões para aceder a este recurso"})
		c.Abort()
	}
}
