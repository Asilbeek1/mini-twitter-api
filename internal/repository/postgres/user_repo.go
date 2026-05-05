package postgresdb

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Asilbeek1/mini-twitter-api/internal/domain"
)

var ErrNotFound = errors.New("record not found")
var ErrDuplicateEntry = errors.New("duplicate entry")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(u *domain.User) (int64, error) {
	query := `
        INSERT INTO users (username, email, first_name, second_name, password_hash)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	var id int64
	err := r.db.QueryRow(query,
		u.Username, u.Email, u.FirstName, u.SecondName, u.PasswordHash,
	).Scan(&id)
	if err != nil {
		if isPostgresUniqueErr(err) {
			return 0, ErrDuplicateEntry
		}
		return 0, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}

func (r *UserRepository) GetByID(id int64) (*domain.User, error) {
	query := `
        	SELECT id, username, email, first_name, second_name,
            password_hash, is_email_verified, created_at, updated_at
        	FROM users WHERE id = $1`
	u := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		u.ID, &u.Username, &u.Email, &u.FirstName, &u.SecondName,
		&u.PasswordHash, &u.IsEmailVerified, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return u, nil

}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
        SELECT id, username, email, first_name, second_name,
               password_hash, is_email_verified, created_at, updated_at
        FROM users WHERE email = $1`

	u := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&u.ID, &u.Username, &u.Email, &u.FirstName, &u.SecondName,
		&u.PasswordHash, &u.IsEmailVerified, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return u, nil
}

func (r *UserRepository) Update(id int64, input *domain.UpdateUserInput) error {
	query := "UPDATE users SET updated_at = CURRENT_TIMESTAMP"
	args := []any{}
	placeholder := 1

	if input.FirstName != nil {
		query += fmt.Sprintf(", first_name = $%d", placeholder)
		args = append(args, *input.FirstName)
		placeholder++
	}
	if input.SecondName != nil {
		query += fmt.Sprintf(", second_name = $%d", placeholder)
		args = append(args, *input.SecondName)
		placeholder++
	}
	if input.Email != nil {
		query += fmt.Sprintf(", email = $%d", placeholder)
		args = append(args, *input.Email)
		placeholder++
	}

	query += fmt.Sprintf(" WHERE id = $%d", placeholder)
	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *UserRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func isPostgresUniqueErr(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "duplicate key value") ||
		strings.Contains(err.Error(), "violates unique constraint")
}
