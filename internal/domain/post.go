package domain

import "time"

type Post struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PostWithAuthor struct {
	Post
	AuthorUsername string `json:"author_username"`
}

type CreatePostInput struct {
	UserID      int64   `json:"user_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type UpdatePostInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
