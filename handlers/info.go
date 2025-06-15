package handlers
import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// GetAPIInfo godoc
// @Summary      Informações da API
// @Description  Como usar a API e autenticação
// @Tags         info
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /info [get]
func GetAPIInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"api":     "RV Contabilidade",
		"version": "1.0",
		"endpoints": gin.H{
			"públicos": []string{
				"POST /auth/register - Criar conta",
				"POST /auth/login - Entrar",
				"GET /info - Esta informação",
			},
			"protegidos": []string{
				"POST /auth/logout - Sair (requer token)",
				"GET /profile - Meu perfil (requer token)",
				"GET /admin/users - Listar utilizadores (requer token admin)",
			},
		},
		"roles": []string{"client", "accountant", "admin"},
		"como_usar_token": gin.H{
			"1": "Fazer login em /auth/login",
			"2": "Copiar o 'token' da resposta",
			"3": "Clicar 'Authorize' no Swagger",
			"4": "Inserir: Bearer SEU_TOKEN_AQUI",
			"5": "Agora pode usar endpoints protegidos",
		},
	})
}
