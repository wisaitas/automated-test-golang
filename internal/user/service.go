package user

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error)
}

type service struct {
	repo Repository
}

func NewService(
	repo Repository,
) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	if err := ValidateCreate(req); err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return &CreateUserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	}, nil
}
