package service

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// AuthService define las operaciones de autenticación
type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
	Logout(ctx context.Context, userID string, refreshToken string) error
	RevokeAllSessions(ctx context.Context, userID string) error
}

type authService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtManager       *auth.JWTManager
	logger           logger.Logger
}

func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtManager *auth.JWTManager,
	logger logger.Logger,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
		logger:           logger,
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

	// Verificar password con bcrypt
	err = auth.VerifyPassword(user.PasswordHash(), req.Password)
	if err != nil {
		s.logger.Warn("invalid password attempt", "email", req.Email)
		return nil, errors.NewUnauthorizedError("invalid credentials")
	}

	// Generar access token JWT (15 minutos)
	accessToken, err := s.jwtManager.GenerateToken(
		user.ID().String(),
		user.Email().String(),
		user.Role(),
		15*time.Minute,
	)
	if err != nil {
		s.logger.Error("failed to generate access token", "error", err)
		return nil, errors.NewInternalError("token generation failed", err)
	}

	// Generar refresh token (7 días)
	refreshToken, err := auth.GenerateRefreshToken(7 * 24 * time.Hour)
	if err != nil {
		s.logger.Error("failed to generate refresh token", "error", err)
		return nil, errors.NewInternalError("refresh token generation failed", err)
	}

	// Guardar refresh token en BD
	tokenData := repository.RefreshTokenData{
		ID:         valueobject.NewUserID().UUID().UUID,
		TokenHash:  refreshToken.TokenHash,
		UserID:     user.ID().UUID().UUID,
		ClientInfo: map[string]string{},
		ExpiresAt:  refreshToken.ExpiresAt,
		CreatedAt:  time.Now(),
	}

	if err := s.refreshTokenRepo.Store(ctx, tokenData); err != nil {
		s.logger.Error("failed to store refresh token", "error", err)
		return nil, errors.NewInternalError("token storage failed", err)
	}

	s.logger.Info("user logged in",
		"user_id", user.ID().String(),
		"email", user.Email().String(),
		"role", user.Role().String(),
	)

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		ExpiresIn:    900, // 15 minutos en segundos
		TokenType:    "Bearer",
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

func (s *authService) RefreshAccessToken(ctx context.Context, refreshTokenStr string) (*dto.RefreshResponse, error) {
	// 1. Hashear el token recibido
	tokenHash := auth.HashToken(refreshTokenStr)

	// 2. Buscar token en BD
	tokenData, err := s.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		s.logger.Error("error finding refresh token", "error", err)
		return nil, errors.NewInternalError("error verifying token", err)
	}
	if tokenData == nil {
		s.logger.Warn("refresh token not found", "token_hash", tokenHash[:8]+"...")
		return nil, errors.NewUnauthorizedError("invalid refresh token")
	}

	// 3. Verificar que no esté revocado
	if tokenData.RevokedAt != nil {
		s.logger.Warn("attempt to use revoked token", "user_id", tokenData.UserID.String())
		return nil, errors.NewUnauthorizedError("token has been revoked")
	}

	// 4. Verificar que no esté expirado
	if time.Now().After(tokenData.ExpiresAt) {
		s.logger.Warn("expired refresh token", "user_id", tokenData.UserID.String())
		return nil, errors.NewUnauthorizedError("refresh token expired")
	}

	// 5. Buscar usuario
	userID, err := valueobject.UserIDFromString(tokenData.UserID.String())
	if err != nil {
		return nil, errors.NewValidationError("invalid user ID")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error("error finding user", "error", err, "user_id", tokenData.UserID.String())
		return nil, errors.NewDatabaseError("find user", err)
	}
	if user == nil || !user.IsActive() {
		s.logger.Warn("user not found or inactive", "user_id", tokenData.UserID.String())
		return nil, errors.NewUnauthorizedError("user not found or inactive")
	}

	// 6. Generar nuevo access token
	newAccessToken, err := s.jwtManager.GenerateToken(
		user.ID().String(),
		user.Email().String(),
		user.Role(),
		15*time.Minute,
	)
	if err != nil {
		s.logger.Error("failed to generate access token", "error", err)
		return nil, errors.NewInternalError("token generation failed", err)
	}

	s.logger.Info("access token refreshed", "user_id", user.ID().String())

	return &dto.RefreshResponse{
		AccessToken: newAccessToken,
		ExpiresIn:   900, // 15 minutos en segundos
		TokenType:   "Bearer",
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID string, refreshTokenStr string) error {
	// Hashear el token recibido
	tokenHash := auth.HashToken(refreshTokenStr)

	// Revocar el token
	if err := s.refreshTokenRepo.Revoke(ctx, tokenHash); err != nil {
		s.logger.Error("failed to revoke token", "error", err, "user_id", userID)
		return errors.NewInternalError("logout failed", err)
	}

	s.logger.Info("user logged out", "user_id", userID)
	return nil
}

func (s *authService) RevokeAllSessions(ctx context.Context, userID string) error {
	// Parsear user ID
	uid, err := valueobject.UserIDFromString(userID)
	if err != nil {
		return errors.NewValidationError("invalid user ID")
	}

	// Revocar todos los tokens del usuario
	if err := s.refreshTokenRepo.RevokeAllByUserID(ctx, uid.UUID().UUID); err != nil {
		s.logger.Error("failed to revoke all sessions", "error", err, "user_id", userID)
		return errors.NewInternalError("revoke failed", err)
	}

	s.logger.Info("all sessions revoked", "user_id", userID)
	return nil
}

