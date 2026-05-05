package service

import (
	"errors"
	"fmt"

	"github.com/Asilbeek1/mini-twitter-api/internal/domain"
	"github.com/Asilbeek1/mini-twitter-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrDuplicateUser      = errors.New("username or email already taken")
	ErrNotFound           = errors.New("user not found")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(input domain.CreateUserInput) (*domain.User, error) {
	if len(input.Password) < 8 {
		return nil, fmt.Errorf("password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	u := *&domain.User{
		Username:     input.Username,
		Email:        input.Email,
		FirstName:    input.FirstName,
		PasswordHash: string(hash),
		Role:         domain.RoleUser,
	}

	id, err := s.repo.Create(&u)
	if errors.Is(err, repository.ErrDuplicateEntry) {
		return nil, ErrDuplicateUser
	}
	if err != nil {
		return nil, err
	}
	u.ID = id
	return &u, nil
}

func (s *UserService) Login(email, password string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(email)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}

func (s *UserService) GetProfile(id int64) (*domain.User, error) {
	profile, err := s.repo.GetByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *UserService) UpdateProfile(callerID, targetID int64, input domain.UpdateUserInput) error {
	if callerID != targetID {
		return ErrUnauthorized
	}
	err := s.repo.Update(targetID, &input)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	return nil
}
