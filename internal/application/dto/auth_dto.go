package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-shared/pkg/validator"
)

// LoginRequest solicitud de login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	v := validator.New()
	v.Required(r.Email, "email")
	v.Email(r.Email, "email")
	v.Required(r.Password, "password")
	v.MinLength(r.Password, 6, "password")
	return v.GetError()
}

// LoginResponse respuesta de login
type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
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
