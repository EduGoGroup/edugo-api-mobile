package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

// ClientInfoMiddleware extrae información del cliente (IP, User-Agent) y la agrega al contexto
// Esta información es utilizada por el servicio de autenticación para rate limiting y auditoría
func ClientInfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extraer IP del cliente
		// ClientIP() respeta X-Forwarded-For y X-Real-IP headers
		clientIP := c.ClientIP()

		// Extraer User-Agent
		userAgent := c.Request.UserAgent()
		if userAgent == "" {
			userAgent = "unknown"
		}

		// Crear un map con la información del cliente
		clientInfo := map[string]interface{}{
			"client_ip":  clientIP,
			"user_agent": userAgent,
		}

		// Agregar al contexto para que el servicio pueda accederlo
		ctx := context.WithValue(c.Request.Context(), "gin_context", clientInfo)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
