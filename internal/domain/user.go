package domain

import "time"

type User struct {
	ID              int64     `json:"id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	FirstName       string    `json:"first_name"`
	SecondName      *string   `json:"second_name,omitempty"`
	PasswordHash    string    `json:"-"`
	IsEmailVerified bool      `json:"is_email_verified"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateUserInput struct {
	Username   string  `json:"username"`
	Email      string  `json:"email"`
	FirstName  string  `json:"first_name"`
	SecondName *string `json:"second_name"`
	Password   string  `json:"password"`
}

type UpdateUserInput struct {
	FirstName  *string `json:"first_name"`
	SecondName *string `json:"second_name"`
	Email      *string `json:"email"`
}
