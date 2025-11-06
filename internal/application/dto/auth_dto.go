package dto

import (
	"github.com/EduGoGroup/edugo-shared/common/validator"
)

// LoginRequest solicitud de login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
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
	AccessToken  string   `json:"access_token"`  // JWT access token
	RefreshToken string   `json:"refresh_token"` // Refresh token para renovar
	ExpiresIn    int      `json:"expires_in"`    // Segundos hasta expiración
	TokenType    string   `json:"token_type"`    // Siempre "Bearer"
	User         UserInfo `json:"user"`          // Información del usuario
}

// RefreshRequest solicitud para refrescar access token
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshResponse respuesta al refrescar token (OAuth2 compatible)
type RefreshResponse struct {
	AccessToken string `json:"access_token"` // Nuevo JWT access token
	ExpiresIn   int    `json:"expires_in"`   // Segundos hasta expiración
	TokenType   string `json:"token_type"`   // Siempre "Bearer"
}

// UserInfo información básica del usuario
type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
}
