package repository

import (
	"context"
	"time"
)

// LoginAttemptReader define operaciones de lectura para intentos de login
// Principio ISP: Separar lectura de escritura
type LoginAttemptReader interface {
	// CountFailedAttempts cuenta intentos fallidos recientes de un identifier
	CountFailedAttempts(ctx context.Context, identifier string, windowMinutes int) (int, error)

	// IsRateLimited verifica si un identifier está bloqueado por rate limiting
	IsRateLimited(ctx context.Context, identifier string, maxAttempts int, windowMinutes int) (bool, error)
}

// LoginAttemptWriter define operaciones de escritura para intentos de login
type LoginAttemptWriter interface {
	// RecordAttempt registra un intento de login
	RecordAttempt(ctx context.Context, attempt LoginAttemptData) error
}

// LoginAttemptRepository agrega todas las capacidades de gestión de intentos de login
// Las implementaciones deben cumplir con todas las interfaces segregadas
type LoginAttemptRepository interface {
	LoginAttemptReader
	LoginAttemptWriter
}

// LoginAttemptData representa los datos de un intento de login
type LoginAttemptData struct {
	Identifier  string // Email o IP
	AttemptType string // "email" o "ip"
	Successful  bool
	UserAgent   string
	IPAddress   string
	AttemptedAt time.Time
}
