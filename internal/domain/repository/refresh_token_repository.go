package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// RefreshTokenRepository define operaciones para gestión de refresh tokens
type RefreshTokenRepository interface {
	// Store guarda un nuevo refresh token en la base de datos
	Store(ctx context.Context, token RefreshTokenData) error

	// FindByTokenHash busca un token por su hash SHA256
	// Retorna nil si el token no existe
	FindByTokenHash(ctx context.Context, tokenHash string) (*RefreshTokenData, error)

	// Revoke marca un token como revocado (logout)
	// Retorna error si el token no existe o ya está revocado
	Revoke(ctx context.Context, tokenHash string) error

	// RevokeAllByUserID revoca todos los tokens activos de un usuario
	// Útil para "cerrar todas las sesiones" o cuando cambia password
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error

	// DeleteExpired elimina tokens expirados (housekeeping)
	// Retorna el número de tokens eliminados
	DeleteExpired(ctx context.Context) (int64, error)
}

// RefreshTokenData representa los datos de un refresh token
type RefreshTokenData struct {
	ID         uuid.UUID
	TokenHash  string
	UserID     uuid.UUID
	ClientInfo map[string]string // IP, UserAgent, Device, etc.
	ExpiresAt  time.Time
	CreatedAt  time.Time
	RevokedAt  *time.Time
	ReplacedBy *uuid.UUID // ID del token que reemplazó a este (rotation)
}
