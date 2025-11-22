package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/google/uuid"
)

// ============================================
// Mocks
// ============================================

type MockUserReader struct {
	mock.Mock
}

func (m *MockUserReader) FindByID(ctx context.Context, id valueobject.UserID) (*pgentities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pgentities.User), args.Error(1)
}

func (m *MockUserReader) FindByEmail(ctx context.Context, email valueobject.Email) (*pgentities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pgentities.User), args.Error(1)
}

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Store(ctx context.Context, token repository.RefreshTokenData) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*repository.RefreshTokenData, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.RefreshTokenData), args.Error(1)
}

func (m *MockRefreshTokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteExpired(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

type MockLoginAttemptRepository struct {
	mock.Mock
}

func (m *MockLoginAttemptRepository) RecordAttempt(ctx context.Context, attempt repository.LoginAttemptData) error {
	args := m.Called(ctx, attempt)
	return args.Error(0)
}

func (m *MockLoginAttemptRepository) CountFailedAttempts(ctx context.Context, identifier string, windowMinutes int) (int, error) {
	args := m.Called(ctx, identifier, windowMinutes)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockLoginAttemptRepository) IsRateLimited(ctx context.Context, identifier string, maxAttempts int, windowMinutes int) (bool, error) {
	args := m.Called(ctx, identifier, maxAttempts, windowMinutes)
	return args.Get(0).(bool), args.Error(1)
}

// ============================================
// Helper Functions
// ============================================

func setupAuthServiceTest(t *testing.T) (
	*authService,
	*MockUserReader,
	*MockRefreshTokenRepository,
	*MockLoginAttemptRepository,
) {
	t.Helper()

	userReader := new(MockUserReader)
	refreshTokenRepo := new(MockRefreshTokenRepository)
	loginAttemptRepo := new(MockLoginAttemptRepository)

	// JWT Manager real para tests
	jwtManager := auth.NewJWTManager("test-secret-key-for-testing-only", "edugo-mobile-api")

	log := new(MockLogger)
	// Setup logger to accept any calls
	log.On("Info", mock.Anything, mock.Anything).Maybe()
	log.On("Warn", mock.Anything, mock.Anything).Maybe()
	log.On("Error", mock.Anything, mock.Anything).Maybe()
	log.On("Debug", mock.Anything, mock.Anything).Maybe()

	service := NewAuthService(
		userReader,
		refreshTokenRepo,
		loginAttemptRepo,
		jwtManager,
		log,
	).(*authService)

	return service, userReader, refreshTokenRepo, loginAttemptRepo
}

func createTestUser(t *testing.T, email string, password string, active bool) *pgentities.User {
	t.Helper()

	emailVO, err := valueobject.NewEmail(email)
	require.NoError(t, err)

	// Hash password con bcrypt
	passwordHash, err := auth.HashPassword(password)
	require.NoError(t, err)

	user := &pgentities.User{
		ID:           valueobject.NewUserID().UUID().UUID,
		Email:        emailVO.String(),
		PasswordHash: passwordHash,
		FirstName:    "Test",
		LastName:     "User",
		Role:         string(enum.SystemRoleStudent),
		IsActive:     active,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return user
}

// ============================================
// Tests: Login
// ============================================

func TestAuthService_Login_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	service, userReader, refreshTokenRepo, loginAttemptRepo := setupAuthServiceTest(t)
	ctx := context.Background()

	testUser := createTestUser(t, "test@example.com", "password123", true)
	email, _ := valueobject.NewEmail("test@example.com")

	// Mock expectations
	loginAttemptRepo.On("IsRateLimited", ctx, "test@example.com", 5, 15).Return(false, nil)
	loginAttemptRepo.On("IsRateLimited", ctx, "unknown", 5, 15).Return(false, nil)
	userReader.On("FindByEmail", ctx, email).Return(testUser, nil)
	refreshTokenRepo.On("Store", ctx, mock.AnythingOfType("repository.RefreshTokenData")).Return(nil)
	loginAttemptRepo.On("RecordAttempt", ctx, mock.AnythingOfType("repository.LoginAttemptData")).Return(nil)

	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Act
	response, err := service.Login(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.Equal(t, 900, response.ExpiresIn)
	assert.Equal(t, "Bearer", response.TokenType)
	assert.Equal(t, testUser.ID.String(), response.User.ID)
	assert.Equal(t, "test@example.com", response.User.Email)

	userReader.AssertExpectations(t)
	refreshTokenRepo.AssertExpectations(t)
	loginAttemptRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidEmail(t *testing.T) {
	t.Parallel()

	service, _, _, _ := setupAuthServiceTest(t)
	ctx := context.Background()

	req := dto.LoginRequest{
		Email:    "invalid-email",
		Password: "password123",
	}

	response, err := service.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestAuthService_Login_RateLimitExceeded(t *testing.T) {
	t.Parallel()

	service, _, _, loginAttemptRepo := setupAuthServiceTest(t)
	ctx := context.Background()

	loginAttemptRepo.On("IsRateLimited", ctx, "test@example.com", 5, 15).Return(true, nil)

	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	response, err := service.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, response)
	loginAttemptRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	t.Parallel()

	service, userReader, _, loginAttemptRepo := setupAuthServiceTest(t)
	ctx := context.Background()

	email, _ := valueobject.NewEmail("notfound@example.com")

	loginAttemptRepo.On("IsRateLimited", ctx, "notfound@example.com", 5, 15).Return(false, nil)
	loginAttemptRepo.On("IsRateLimited", ctx, "unknown", 5, 15).Return(false, nil)
	userReader.On("FindByEmail", ctx, email).Return(nil, nil)
	loginAttemptRepo.On("RecordAttempt", ctx, mock.AnythingOfType("repository.LoginAttemptData")).Return(nil)

	req := dto.LoginRequest{
		Email:    "notfound@example.com",
		Password: "password123",
	}

	response, err := service.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, response)
	userReader.AssertExpectations(t)
	loginAttemptRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	t.Parallel()

	service, userReader, _, loginAttemptRepo := setupAuthServiceTest(t)
	ctx := context.Background()

	testUser := createTestUser(t, "test@example.com", "correctpassword", true)
	email, _ := valueobject.NewEmail("test@example.com")

	loginAttemptRepo.On("IsRateLimited", ctx, "test@example.com", 5, 15).Return(false, nil)
	loginAttemptRepo.On("IsRateLimited", ctx, "unknown", 5, 15).Return(false, nil)
	userReader.On("FindByEmail", ctx, email).Return(testUser, nil)
	loginAttemptRepo.On("RecordAttempt", ctx, mock.AnythingOfType("repository.LoginAttemptData")).Return(nil)

	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	response, err := service.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, response)
	userReader.AssertExpectations(t)
	loginAttemptRepo.AssertExpectations(t)
}

// ============================================
// Tests: RefreshAccessToken
// ============================================

func TestAuthService_RefreshAccessToken_Success(t *testing.T) {
	t.Parallel()

	service, userReader, refreshTokenRepo, _ := setupAuthServiceTest(t)
	ctx := context.Background()

	testUser := createTestUser(t, "test@example.com", "password123", true)
	refreshToken := "test-refresh-token"
	tokenHash := auth.HashToken(refreshToken)

	tokenData := &repository.RefreshTokenData{
		ID:        uuid.New(),
		TokenHash: tokenHash,
		UserID:    testUser.ID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
		RevokedAt: nil,
	}

	refreshTokenRepo.On("FindByTokenHash", ctx, tokenHash).Return(tokenData, nil)
	userReader.On("FindByID", ctx, testUser.ID).Return(testUser, nil)

	response, err := service.RefreshAccessToken(ctx, refreshToken)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	assert.Equal(t, 900, response.ExpiresIn)
	assert.Equal(t, "Bearer", response.TokenType)

	refreshTokenRepo.AssertExpectations(t)
	userReader.AssertExpectations(t)
}

func TestAuthService_RefreshAccessToken_TokenNotFound(t *testing.T) {
	t.Parallel()

	service, _, refreshTokenRepo, _ := setupAuthServiceTest(t)
	ctx := context.Background()

	refreshToken := "invalid-token"
	tokenHash := auth.HashToken(refreshToken)

	refreshTokenRepo.On("FindByTokenHash", ctx, tokenHash).Return(nil, nil)

	response, err := service.RefreshAccessToken(ctx, refreshToken)

	assert.Error(t, err)
	assert.Nil(t, response)
	refreshTokenRepo.AssertExpectations(t)
}

func TestAuthService_RefreshAccessToken_TokenRevoked(t *testing.T) {
	t.Parallel()

	service, _, refreshTokenRepo, _ := setupAuthServiceTest(t)
	ctx := context.Background()

	refreshToken := "revoked-token"
	tokenHash := auth.HashToken(refreshToken)
	revokedAt := time.Now().Add(-1 * time.Hour)

	tokenData := &repository.RefreshTokenData{
		ID:        uuid.New(),
		TokenHash: tokenHash,
		UserID:    uuid.New(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
		RevokedAt: &revokedAt,
	}

	refreshTokenRepo.On("FindByTokenHash", ctx, tokenHash).Return(tokenData, nil)

	response, err := service.RefreshAccessToken(ctx, refreshToken)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "revoked")
	refreshTokenRepo.AssertExpectations(t)
}

func TestAuthService_RefreshAccessToken_TokenExpired(t *testing.T) {
	t.Parallel()

	service, _, refreshTokenRepo, _ := setupAuthServiceTest(t)
	ctx := context.Background()

	refreshToken := "expired-token"
	tokenHash := auth.HashToken(refreshToken)

	tokenData := &repository.RefreshTokenData{
		ID:        uuid.New(),
		TokenHash: tokenHash,
		UserID:    uuid.New(),
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		CreatedAt: time.Now().Add(-8 * 24 * time.Hour),
		RevokedAt: nil,
	}

	refreshTokenRepo.On("FindByTokenHash", ctx, tokenHash).Return(tokenData, nil)

	response, err := service.RefreshAccessToken(ctx, refreshToken)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "expired")
	refreshTokenRepo.AssertExpectations(t)
}

// ============================================
// Tests: Logout
// ============================================

func TestAuthService_Logout_Success(t *testing.T) {
	t.Parallel()

	service, _, refreshTokenRepo, _ := setupAuthServiceTest(t)
	ctx := context.Background()

	userID := uuid.New().String()
	refreshToken := "test-refresh-token"
	tokenHash := auth.HashToken(refreshToken)

	refreshTokenRepo.On("Revoke", ctx, tokenHash).Return(nil)

	err := service.Logout(ctx, userID, refreshToken)

	assert.NoError(t, err)
	refreshTokenRepo.AssertExpectations(t)
}

func TestAuthService_Logout_RevokeFailure(t *testing.T) {
	t.Parallel()

	service, _, refreshTokenRepo, _ := setupAuthServiceTest(t)
	ctx := context.Background()

	userID := uuid.New().String()
	refreshToken := "test-refresh-token"
	tokenHash := auth.HashToken(refreshToken)

	refreshTokenRepo.On("Revoke", ctx, tokenHash).Return(errors.NewDatabaseError("revoke", nil))

	err := service.Logout(ctx, userID, refreshToken)

	assert.Error(t, err)
	refreshTokenRepo.AssertExpectations(t)
}

// ============================================
// Tests: RevokeAllSessions
// ============================================

func TestAuthService_RevokeAllSessions_Success(t *testing.T) {
	t.Parallel()

	service, _, refreshTokenRepo, _ := setupAuthServiceTest(t)
	ctx := context.Background()

	userID := valueobject.NewUserID()

	refreshTokenRepo.On("RevokeAllByUserID", ctx, userID.UUID().UUID).Return(nil)

	err := service.RevokeAllSessions(ctx, userID.String())

	assert.NoError(t, err)
	refreshTokenRepo.AssertExpectations(t)
}

func TestAuthService_RevokeAllSessions_InvalidUserID(t *testing.T) {
	t.Parallel()

	service, _, _, _ := setupAuthServiceTest(t)
	ctx := context.Background()

	err := service.RevokeAllSessions(ctx, "invalid-uuid")

	assert.Error(t, err)
}

func TestAuthService_RevokeAllSessions_RevokeFailure(t *testing.T) {
	t.Parallel()

	service, _, refreshTokenRepo, _ := setupAuthServiceTest(t)
	ctx := context.Background()

	userID := valueobject.NewUserID()

	refreshTokenRepo.On("RevokeAllByUserID", ctx, userID.UUID().UUID).Return(errors.NewDatabaseError("revoke all", nil))

	err := service.RevokeAllSessions(ctx, userID.String())

	assert.Error(t, err)
	refreshTokenRepo.AssertExpectations(t)
}
