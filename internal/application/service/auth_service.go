package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/pkg/auth"
	"github.com/EduGoGroup/edugo-shared/pkg/errors"
	"github.com/EduGoGroup/edugo-shared/pkg/logger"
)

// AuthService define las operaciones de autenticación
type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager
	logger     logger.Logger
}

func NewAuthService(
	userRepo repository.UserRepository,
	jwtManager *auth.JWTManager,
	logger logger.Logger,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Validar request
	if err := req.Validate(); err != nil {
		s.logger.Warn("validation failed", "error", err)
		return nil, err
	}

	// Buscar usuario por email
	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		s.logger.Error("failed to find user", "error", err)
		return nil, errors.NewDatabaseError("find user", err)
	}

	if user == nil || !user.IsActive() {
		return nil, errors.NewUnauthorizedError("invalid credentials")
	}

	// Verificar password (hash simple para ejemplo, en prod usar bcrypt)
	passwordHash := hashPassword(req.Password)
	if user.PasswordHash() != passwordHash {
		s.logger.Warn("invalid password attempt", "email", req.Email)
		return nil, errors.NewUnauthorizedError("invalid credentials")
	}

	// Generar JWT usando shared/auth
	token, err := s.jwtManager.GenerateToken(
		user.ID().String(),
		user.Email().String(),
		user.Role(),
		24*time.Hour, // Token expira en 24h
	)
	if err != nil {
		s.logger.Error("failed to generate token", "error", err)
		return nil, errors.NewInternalError("token generation failed", err)
	}

	// Generar refresh token (expira en 7 días)
	refreshToken, err := s.jwtManager.GenerateToken(
		user.ID().String(),
		user.Email().String(),
		user.Role(),
		7*24*time.Hour,
	)
	if err != nil {
		s.logger.Error("failed to generate refresh token", "error", err)
		return nil, errors.NewInternalError("token generation failed", err)
	}

	s.logger.Info("user logged in",
		"user_id", user.ID().String(),
		"email", user.Email().String(),
		"role", user.Role().String(),
	)

	return &dto.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		User: dto.UserInfo{
			ID:        user.ID().String(),
			Email:     user.Email().String(),
			FirstName: user.FirstName(),
			LastName:  user.LastName(),
			FullName:  user.FullName(),
			Role:      user.Role().String(),
		},
	}, nil
}

// hashPassword crea un hash simple del password (en prod usar bcrypt)
func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}
