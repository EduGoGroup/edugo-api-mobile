package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired middleware para validar JWT token
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token de autorización requerido",
			})
			return
		}

		// Formato esperado: "Bearer {token}"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Formato de token inválido. Use: Bearer {token}",
			})
			return
		}

		token := parts[1]

		// TODO: Validar JWT token real
		// TODO: Extraer user_id y role del token
		// TODO: Setear en contexto

		// Mock: Aceptar cualquier token que no esté vacío
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token vacío",
			})
			return
		}

		// Setear usuario mock en contexto
		c.Set("user_id", "user-mock-123")
		c.Set("role", "teacher")

		c.Next()
	}
}

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RequestLogger middleware
func RequestLogger() gin.HandlerFunc {
	return gin.Logger()
}

// RateLimiter middleware (mock)
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implementar rate limiting real con Redis
		c.Next()
	}
}
