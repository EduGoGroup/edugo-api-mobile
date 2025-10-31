package entity

import (
	"time"

	"github.com/edugo/api-mobile/internal/domain/valueobject"
	"github.com/edugo/shared/pkg/types/enum"
)

// User representa un usuario móvil (student, teacher, guardian)
type User struct {
	id           valueobject.UserID
	email        valueobject.Email
	passwordHash string
	firstName    string
	lastName     string
	role         enum.SystemRole
	isActive     bool
	createdAt    time.Time
	updatedAt    time.Time
}

// ReconstructUser reconstruye un User desde la base de datos
func ReconstructUser(
	id valueobject.UserID,
	email valueobject.Email,
	passwordHash string,
	firstName string,
	lastName string,
	role enum.SystemRole,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
) *User {
	return &User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		firstName:    firstName,
		lastName:     lastName,
		role:         role,
		isActive:     isActive,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// Getters

func (u *User) ID() valueobject.UserID     { return u.id }
func (u *User) Email() valueobject.Email   { return u.email }
func (u *User) PasswordHash() string       { return u.passwordHash }
func (u *User) FirstName() string          { return u.firstName }
func (u *User) LastName() string           { return u.lastName }
func (u *User) FullName() string           { return u.firstName + " " + u.lastName }
func (u *User) Role() enum.SystemRole      { return u.role }
func (u *User) IsActive() bool             { return u.isActive }
func (u *User) CreatedAt() time.Time       { return u.createdAt }
func (u *User) UpdatedAt() time.Time       { return u.updatedAt }
