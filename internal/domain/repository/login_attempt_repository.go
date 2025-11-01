package repository

import (
	"context"
	"time"
)

// LoginAttemptRepository define operaciones para gestión de intentos de login
type LoginAttemptRepository interface {
	// RecordAttempt registra un intento de login
	RecordAttempt(ctx context.Context, attempt LoginAttemptData) error

	// CountFailedAttempts cuenta intentos fallidos recientes de un identifier
	CountFailedAttempts(ctx context.Context, identifier string, windowMinutes int) (int, error)

	// IsRateLimited verifica si un identifier está bloqueado por rate limiting
	IsRateLimited(ctx context.Context, identifier string, maxAttempts int, windowMinutes int) (bool, error)
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
