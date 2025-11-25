// Package middleware proporciona middlewares HTTP para la aplicación.
// RemoteAuthMiddleware valida tokens JWT contra api-admin como autoridad central.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-mobile/internal/client"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// RemoteAuthConfig configuración del middleware de autenticación remota
type RemoteAuthConfig struct {
	AuthClient     *client.AuthClient                   // Cliente para validar tokens con api-admin
	Logger         logger.Logger                        // Logger para registrar eventos
	SkipPaths      []string                             // Paths que no requieren autenticación
	OnUnauthorized func(c *gin.Context, message string) // Handler personalizado para 401
}

// RemoteAuthMiddleware crea un middleware que valida tokens con api-admin
// Este middleware reemplaza la validación local de JWT por validación remota centralizada
func RemoteAuthMiddleware(config RemoteAuthConfig) gin.HandlerFunc {
	// Construir mapa de paths a saltar para búsqueda O(1)
	skipPaths := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		// 1. Verificar si el path debe saltarse
		if skipPaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		// 2. Extraer token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handleUnauthorized(c, config, "Token de autorización requerido")
			return
		}

		// 3. Parsear Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			handleUnauthorized(c, config, "Formato de token inválido, use: Bearer <token>")
			return
		}
		token := parts[1]

		// 4. Validar token con api-admin
		tokenInfo, err := config.AuthClient.ValidateToken(c.Request.Context(), token)
		if err != nil {
			if config.Logger != nil {
				config.Logger.Error("error validando token con api-admin", "error", err)
			}
			handleUnauthorized(c, config, "Error validando token")
			return
		}

		// 5. Verificar que el token es válido
		if !tokenInfo.Valid {
			if config.Logger != nil {
				config.Logger.Warn("token inválido", "error", tokenInfo.Error)
			}
			message := "Token inválido o expirado"
			if tokenInfo.Error != "" {
				message = tokenInfo.Error
			}
			handleUnauthorized(c, config, message)
			return
		}

		// 6. Inyectar información del usuario en el contexto de Gin
		// Estos valores estarán disponibles en los handlers
		c.Set("user_id", tokenInfo.UserID)
		c.Set("email", tokenInfo.Email)
		c.Set("role", tokenInfo.Role)
		c.Set("token_expires_at", tokenInfo.ExpiresAt)

		if config.Logger != nil {
			config.Logger.Debug("autenticación exitosa",
				"user_id", tokenInfo.UserID,
				"role", tokenInfo.Role,
			)
		}

		// 7. Continuar con el siguiente handler
		c.Next()
	}
}

// handleUnauthorized maneja respuestas 401 Unauthorized
func handleUnauthorized(c *gin.Context, config RemoteAuthConfig, message string) {
	if config.OnUnauthorized != nil {
		config.OnUnauthorized(c, message)
		return
	}

	// Respuesta por defecto
	c.JSON(http.StatusUnauthorized, gin.H{
		"error":   "unauthorized",
		"message": message,
		"code":    "UNAUTHORIZED",
	})
	c.Abort()
}

// ============================================
// Helper Functions para obtener datos del contexto
// ============================================

// GetUserID obtiene el user_id del contexto de Gin
func GetUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserEmail obtiene el email del contexto de Gin
func GetUserEmail(c *gin.Context) string {
	if email, exists := c.Get("email"); exists {
		if e, ok := email.(string); ok {
			return e
		}
	}
	return ""
}

// GetUserRole obtiene el role del contexto de Gin
func GetUserRole(c *gin.Context) string {
	if role, exists := c.Get("role"); exists {
		if r, ok := role.(string); ok {
			return r
		}
	}
	return ""
}

// RequireRole middleware que verifica que el usuario tenga uno de los roles permitidos
// Debe usarse DESPUÉS de RemoteAuthMiddleware
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	roleMap := make(map[string]bool)
	for _, role := range allowedRoles {
		roleMap[role] = true
	}

	return func(c *gin.Context) {
		userRole := GetUserRole(c)
		if userRole == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "No se encontró rol de usuario",
				"code":    "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		if !roleMap[userRole] {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "No tiene permisos para esta operación",
				"code":    "FORBIDDEN",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
