package repository

import (
	"context"
	"database/sql"

	"github.com/edugo/api-mobile/internal/domain/entity"
	"github.com/edugo/api-mobile/internal/domain/repository"
	"github.com/edugo/api-mobile/internal/domain/valueobject"
	"github.com/edugo/shared/pkg/types/enum"
)

type postgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) FindByID(ctx context.Context, id valueobject.UserID) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var (
		idStr        string
		email        string
		passwordHash string
		firstName    string
		lastName     string
		roleStr      string
		isActive     bool
		createdAt    sql.NullTime
		updatedAt    sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &email, &passwordHash, &firstName, &lastName, &roleStr, &isActive, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	userID, _ := valueobject.UserIDFromString(idStr)
	emailVO, _ := valueobject.NewEmail(email)

	return entity.ReconstructUser(
		userID,
		emailVO,
		passwordHash,
		firstName,
		lastName,
		enum.SystemRole(roleStr),
		isActive,
		createdAt.Time,
		updatedAt.Time,
	), nil
}

func (r *postgresUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var (
		idStr        string
		emailStr     string
		passwordHash string
		firstName    string
		lastName     string
		roleStr      string
		isActive     bool
		createdAt    sql.NullTime
		updatedAt    sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(
		&idStr, &emailStr, &passwordHash, &firstName, &lastName, &roleStr, &isActive, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	userID, _ := valueobject.UserIDFromString(idStr)
	emailVO, _ := valueobject.NewEmail(emailStr)

	return entity.ReconstructUser(
		userID,
		emailVO,
		passwordHash,
		firstName,
		lastName,
		enum.SystemRole(roleStr),
		isActive,
		createdAt.Time,
		updatedAt.Time,
	), nil
}

func (r *postgresUserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query,
		user.FirstName(),
		user.LastName(),
		user.UpdatedAt(),
		user.ID().String(),
	)

	return err
}
