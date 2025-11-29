package fixtures

import (
	"time"

	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

var (
	AdminID   = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	TeacherID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	StudentID = uuid.MustParse("33333333-3333-3333-3333-333333333333")
)

func GetDefaultUsers() map[string]*pgentities.User {
	now := time.Now()
	return map[string]*pgentities.User{
		AdminID.String(): {
			ID:           AdminID,
			Email:        "admin@edugo.com",
			PasswordHash: "$2a$10$hash",
			FirstName:    "Admin",
			LastName:     "Sistema",
			Role:         "admin",
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		TeacherID.String(): {
			ID:           TeacherID,
			Email:        "teacher@edugo.com",
			PasswordHash: "$2a$10$hash",
			FirstName:    "María",
			LastName:     "González",
			Role:         "teacher",
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		StudentID.String(): {
			ID:           StudentID,
			Email:        "student@edugo.com",
			PasswordHash: "$2a$10$hash",
			FirstName:    "Juan",
			LastName:     "Pérez",
			Role:         "student",
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}
}
