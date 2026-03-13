package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware extrai e valida o Bearer token do header Authorization.
// Se inválido ou ausente: retorna 401 e aborta a requisição.
// Se válido: adiciona o userID e familiaID ao contexto Gin.
func AuthMiddleware(authService service.AuthServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header ausente"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "formato de authorization inválido"})
			return
		}

		token := parts[1]
		usuarioID, familiaID, err := authService.ValidateAccessToken(token)
		if err != nil {
			if errors.Is(err, domain.ErrTokenExpirado) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expirado"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			return
		}

		c.Set("userID", usuarioID)
		c.Set("familiaID", familiaID)
		c.Next()
	}
}
