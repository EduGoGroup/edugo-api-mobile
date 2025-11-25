package dto

import (
	"github.com/EduGoGroup/edugo-shared/common/validator"
)

// LoginRequest solicitud de login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

func (r *LoginRequest) Validate() error {
	v := validator.New()
	v.Required(r.Email, "email")
	v.Email(r.Email, "email")
	v.Required(r.Password, "password")
	v.MinLength(r.Password, 6, "password")
	return v.GetError()
}

// LoginResponse respuesta de login (OAuth2 compatible)
type LoginResponse struct {
	AccessToken  string   `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT access token
	RefreshToken string   `json:"refresh_token" example:"550e8400-e29b-41d4-a716-446655440000"`   // Refresh token para renovar
	ExpiresIn    int      `json:"expires_in" example:"3600"`                                      // Segundos hasta expiración
	TokenType    string   `json:"token_type" example:"Bearer"`                                    // Siempre "Bearer"
	User         UserInfo `json:"user"`                                                           // Información del usuario
}

// RefreshRequest solicitud para refrescar access token
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// RefreshResponse respuesta al refrescar token (OAuth2 compatible)
type RefreshResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // Nuevo JWT access token
	ExpiresIn   int    `json:"expires_in" example:"3600"`                                      // Segundos hasta expiración
	TokenType   string `json:"token_type" example:"Bearer"`                                    // Siempre "Bearer"
}

// UserInfo información básica del usuario
type UserInfo struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string `json:"email" example:"user@example.com"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	FullName  string `json:"full_name" example:"John Doe"`
	Role      string `json:"role" example:"student"`
}
