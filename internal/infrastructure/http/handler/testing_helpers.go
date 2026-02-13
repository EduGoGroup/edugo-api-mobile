package handler

import (
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
)

// SetupTestRouter crea un router Gin en modo test
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// MockAuthMiddleware simula middleware de autenticación completo (user_id y school_id)
func MockAuthMiddleware(userID, schoolID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Set("school_id", schoolID)
		c.Next()
	}
}

// NewTestLogger crea un logger silencioso para tests (solo errores)
func NewTestLogger() logger.Logger {
	return logger.NewZapLogger("error", "json")
}
