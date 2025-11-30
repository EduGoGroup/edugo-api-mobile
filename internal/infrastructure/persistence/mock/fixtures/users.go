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
	// Hashes bcrypt válidos para la contraseña "password123" (cost=10)
	return map[string]*pgentities.User{
		AdminID.String(): {
			ID:           AdminID,
			Email:        "admin@edugo.com",
			PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
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
			PasswordHash: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
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
			PasswordHash: "$2a$10$eImiTXuWVxfM37uY4JANjQ5oG7lKwK1ggP1CJE3.JiQ3h6.qE6pOu",
			FirstName:    "Juan",
			LastName:     "Pérez",
			Role:         "student",
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}
}
