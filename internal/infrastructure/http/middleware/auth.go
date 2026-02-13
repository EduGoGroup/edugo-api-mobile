package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// Constantes para keys del contexto
const (
	ContextKeyUserID   = "user_id"
	ContextKeyEmail    = "email"
	ContextKeyRole     = "role"
	ContextKeySchoolID = "school_id"
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
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyEmail, claims.Email)
		c.Set(ContextKeyRole, claims.Role)
		c.Set(ContextKeySchoolID, claims.SchoolID)

		// RBAC: Inyectar ActiveContext si está disponible
		if claims.ActiveContext != nil {
			c.Set(ContextKeyActiveContext, claims.ActiveContext)
			// Sobreescribir role con el del contexto RBAC
			c.Set(ContextKeyRole, claims.ActiveContext.RoleName)
		}

		log.Debug("auth successful",
			"user_id", claims.UserID,
			"role", claims.Role,
			"school_id", claims.SchoolID,
		)

		c.Next()
	}
}

// GetUserIDFromContext obtiene el user_id del contexto de Gin
// Retorna el UUID y true si existe, uuid.Nil y false si no existe o es inválido
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userIDStr, exists := c.Get(ContextKeyUserID)
	if !exists {
		return uuid.Nil, false
	}

	str, ok := userIDStr.(string)
	if !ok {
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil, false
	}

	return userID, true
}

// MustGetUserIDFromContext obtiene el user_id del contexto de Gin
// Hace panic si no existe o es inválido (usar solo después de AuthRequired middleware)
func MustGetUserIDFromContext(c *gin.Context) uuid.UUID {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		panic("user_id not found in context - ensure AuthRequired middleware is applied")
	}
	return userID
}

// GetSchoolIDFromContext obtiene el school_id del contexto de Gin
// Retorna el UUID y true si existe, uuid.Nil y false si no existe o es inválido
func GetSchoolIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	schoolIDStr, exists := c.Get(ContextKeySchoolID)
	if !exists {
		return uuid.Nil, false
	}

	str, ok := schoolIDStr.(string)
	if !ok || str == "" {
		return uuid.Nil, false
	}

	schoolID, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil, false
	}

	return schoolID, true
}

// MustGetSchoolIDFromContext obtiene el school_id del contexto de Gin
// Hace panic si no existe o es inválido (usar solo cuando school_id es obligatorio)
func MustGetSchoolIDFromContext(c *gin.Context) uuid.UUID {
	schoolID, ok := GetSchoolIDFromContext(c)
	if !ok {
		panic("school_id not found in context - ensure JWT contains school_id claim")
	}
	return schoolID
}

// GetRoleFromContext obtiene el role del contexto de Gin
func GetRoleFromContext(c *gin.Context) (string, bool) {
	role, exists := c.Get(ContextKeyRole)
	if !exists {
		return "", false
	}

	// El role puede venir como enum.SystemRole o string
	switch v := role.(type) {
	case string:
		return v, true
	default:
		return "", false
	}
}

// GetEmailFromContext obtiene el email del contexto de Gin
func GetEmailFromContext(c *gin.Context) (string, bool) {
	email, exists := c.Get(ContextKeyEmail)
	if !exists {
		return "", false
	}

	str, ok := email.(string)
	return str, ok
}

// Deprecated: Usar RequirePermission() en su lugar.
// IsAdminRole verifica si el usuario autenticado tiene rol admin o super_admin
func IsAdminRole(c *gin.Context) bool {
	role, ok := GetRoleFromContext(c)
	if !ok {
		return false
	}
	return role == "admin" || role == "super_admin"
}

// Deprecated: Usar RequirePermission() en su lugar.
// HasRole verifica si el usuario autenticado tiene alguno de los roles especificados
func HasRole(c *gin.Context, roles ...string) bool {
	role, ok := GetRoleFromContext(c)
	if !ok {
		return false
	}
	for _, r := range roles {
		if role == r {
			return true
		}
	}
	return false
}
