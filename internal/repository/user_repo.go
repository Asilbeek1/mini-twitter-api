package repository

import (
	"database/sql"
	"errors"
	"fmt"

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
        INSERT INTO users (username, email, first_name, second_name, password_hash, role)
        VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		u.Username, u.Email, u.FirstName, u.SecondName, u.PasswordHash, u.Role,
	)
	if err != nil {
		if isSQLiteUniqueErr(err) {
			return 0, ErrDuplicateEntry
		}
		return 0, fmt.Errorf("create user: %w", err)
	}

	return result.LastInsertId()
}

func (r *UserRepository) GetByID(id int64) (*domain.User, error) {
	query := `
        	SELECT id, username, email, first_name, second_name,
            password_hash, role, is_email_verified, created_at, updated_at
        	FROM users WHERE id = ?`
	u := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		u.ID, &u.Username, &u.Email, &u.FirstName, &u.SecondName,
		&u.PasswordHash, &u.Role, &u.IsEmailVerified, &u.CreatedAt, &u.UpdatedAt,
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
               password_hash, role, is_email_verified, created_at, updated_at
        FROM users WHERE email = ?`

	u := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&u.ID, &u.Username, &u.Email, &u.FirstName, &u.SecondName,
		&u.PasswordHash, &u.Role, &u.IsEmailVerified, &u.CreatedAt, &u.UpdatedAt,
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

	if input.FirstName != nil {
		query += ", first_name = ?"
		args = append(args, *input.FirstName)
	}
	if input.SecondName != nil {
		query += ", second_name = ?"
		args = append(args, *input.SecondName)
	}
	if input.Email != nil {
		query += ", email = ?"
		args = append(args, *input.Email)
	}

	query += " WHERE id = ?"
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
	result, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func isSQLiteUniqueErr(err error) bool {
	return err != nil && errors.Is(err, errors.New("UNIQUE constraint failed")) ||
		containsStr(err.Error(), "UNIQUE constraint failed")
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && findStr(s, substr))
}

func findStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
