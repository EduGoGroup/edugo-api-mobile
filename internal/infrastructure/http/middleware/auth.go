package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/edugo/shared/pkg/auth"
	"github.com/edugo/shared/pkg/logger"
)

// AuthRequired middleware que requiere autenticación JWT
func AuthRequired(jwtManager *auth.JWTManager, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn("missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization required", "code": "UNAUTHORIZED"})
			c.Abort()
			return
		}

		// Extraer token del header "Bearer {token}"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Warn("invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header", "code": "UNAUTHORIZED"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validar token usando shared/auth
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			log.Warn("invalid token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token", "code": "UNAUTHORIZED"})
			c.Abort()
			return
		}

		// Agregar claims al contexto para uso en handlers
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		log.Debug("auth successful", "user_id", claims.UserID, "role", claims.Role)

		c.Next()
	}
}
