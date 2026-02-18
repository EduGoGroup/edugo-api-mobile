package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

type postgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) FindByID(ctx context.Context, id valueobject.UserID) (*pgentities.User, error) {
	query := `
		SELECT
			u.id,
			u.email,
			u.password_hash,
			u.first_name,
			u.last_name,
			COALESCE((
				SELECT r.name
				FROM user_roles ur
				JOIN roles r ON r.id = ur.role_id
				WHERE ur.user_id = u.id AND ur.is_active = true
				ORDER BY ur.granted_at DESC
				LIMIT 1
			), '') AS role,
			u.is_active,
			u.created_at,
			u.updated_at
		FROM users u
		WHERE u.id = $1 AND u.deleted_at IS NULL
	`

	var (
		userID       uuid.UUID
		email        string
		passwordHash string
		firstName    string
		lastName     string
		role         string
		isActive     bool
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := r.db.QueryRowContext(ctx, query, id.UUID()).Scan(
		&userID, &email, &passwordHash, &firstName, &lastName, &role, &isActive, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &pgentities.User{
		ID:           userID,
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		IsActive:     isActive,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}

func (r *postgresUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*pgentities.User, error) {
	query := `
		SELECT
			u.id,
			u.email,
			u.password_hash,
			u.first_name,
			u.last_name,
			COALESCE((
				SELECT r.name
				FROM user_roles ur
				JOIN roles r ON r.id = ur.role_id
				WHERE ur.user_id = u.id AND ur.is_active = true
				ORDER BY ur.granted_at DESC
				LIMIT 1
			), '') AS role,
			u.is_active,
			u.created_at,
			u.updated_at
		FROM users u
		WHERE u.email = $1 AND u.deleted_at IS NULL
	`

	var (
		userID       uuid.UUID
		emailStr     string
		passwordHash string
		firstName    string
		lastName     string
		role         string
		isActive     bool
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(
		&userID, &emailStr, &passwordHash, &firstName, &lastName, &role, &isActive, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &pgentities.User{
		ID:           userID,
		Email:        emailStr,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		IsActive:     isActive,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}

func (r *postgresUserRepository) Update(ctx context.Context, user *pgentities.User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, is_active = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)

	return err
}
