package service

import (
	"errors"
	"fmt"

	"github.com/Asilbeek1/mini-twitter-api/internal/domain"
	"github.com/Asilbeek1/mini-twitter-api/internal/repository"
)

var ErrForbidden = errors.New("You don`t own this post")

type PostService struct {
	postRepo *repository.PostRepository
}

func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{postRepo: postRepo}
}

func (s *PostService) CreatePost(input domain.CreatePostInput) (*domain.Post, error) {
	if input.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if len(input.Title) > 100 {
		return nil, fmt.Errorf("Title length should be less than 100 characaters")
	}
	post := &domain.Post{
		UserID:      input.UserID,
		Title:       input.Title,
		Description: input.Description,
	}

	id, err := s.postRepo.Create(post)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, repository.ErrNotFound
	}
	post.ID = id

	return post, nil
}

func (s *PostService) GetFeed(page, pageSize int) ([]*domain.PostWithAuthor, error) {
	if pageSize > 100 {
		pageSize = 100 // cap page size
	}
	offset := (page - 1) * pageSize
	return s.postRepo.ListFeed(pageSize, offset)
}

func (s *PostService) GetUserPosts(userID int64, page, pageSize int) ([]*domain.Post, error) {
	offset := (page - 1) * pageSize
	return s.postRepo.ListByUser(userID, pageSize, offset)
}

func (s *PostService) UpdatePost(callerID, postID int64, input domain.UpdatePostInput) error {
	post, err := s.postRepo.GetByID(postID)
	if errors.Is(err, repository.ErrNotFound) {
		return repository.ErrNotFound
	}
	if err != nil {
		return err
	}

	if post.UserID != callerID {
		return ErrForbidden
	}

	return s.postRepo.Update(postID, &input)
}

func (s *PostService) DeletePost(callerID int64, callerRole domain.Role, postID int64) error {
	post, err := s.postRepo.GetByID(postID)
	if errors.Is(err, repository.ErrNotFound) {
		return repository.ErrNotFound
	}
	if err != nil {
		return err
	}

	if post.UserID != callerID && callerRole != domain.RoleAdmin {
		return ErrForbidden
	}

	return s.postRepo.Delete(postID)
}
