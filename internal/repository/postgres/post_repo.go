package postgresdb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Asilbeek1/mini-twitter-api/internal/domain"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(p *domain.Post) (int64, error) {
	query := `
        INSERT INTO posts (user_id, title, description)
        VALUES ($1, $2, $3)
        RETURNING id`

	var id int64
	err := r.db.QueryRow(query, p.UserID, p.Title, p.Description).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create post: %w", err)
	}
	return id, nil
}

func (r *PostRepository) GetByID(id int64) (*domain.Post, error) {
	query := `
        SELECT id, user_id, title, description, created_at, updated_at
        FROM posts WHERE id = $1`

	p := &domain.Post{}
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.UserID, &p.Title, &p.Description, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get post by id: %w", err)
	}
	return p, nil
}

func (r *PostRepository) ListFeed(limit, offset int) ([]*domain.PostWithAuthor, error) {
	query := `
        SELECT p.id, p.user_id, p.title, p.description, p.created_at, p.updated_at,
               u.username
        FROM posts p
        JOIN users u ON u.id = p.user_id
        ORDER BY p.created_at DESC
        LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list feed: %w", err)
	}
	defer rows.Close()

	var posts []*domain.PostWithAuthor
	for rows.Next() {
		p := &domain.PostWithAuthor{}
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.Title, &p.Description,
			&p.CreatedAt, &p.UpdatedAt, &p.AuthorUsername,
		); err != nil {
			return nil, fmt.Errorf("scan post: %w", err)
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *PostRepository) ListByUser(userID int64, limit, offset int) ([]*domain.Post, error) {
	query := `
        SELECT id, user_id, title, description, created_at, updated_at
        FROM posts WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list posts by user: %w", err)
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		p := &domain.Post{}
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan post: %w", err)
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *PostRepository) Update(id int64, input *domain.UpdatePostInput) error {
	query := "UPDATE posts SET updated_at = CURRENT_TIMESTAMP"
	args := []any{}
	placeholder := 1

	if input.Title != nil {
		query += fmt.Sprintf(", title = $%d", placeholder)
		args = append(args, *input.Title)
		placeholder++
	}
	if input.Description != nil {
		query += fmt.Sprintf(", description = $%d", placeholder)
		args = append(args, *input.Description)
		placeholder++
	}

	query += fmt.Sprintf(" WHERE id = $%d", placeholder)
	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("update post: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
