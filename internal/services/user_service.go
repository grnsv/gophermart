package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/grnsv/gophermart/internal/models"
	"github.com/grnsv/gophermart/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	storage storage.UserRepository
}

func NewUserService(storage storage.UserRepository) UserService {
	return &userService{storage: storage}
}

func (s *userService) IsLoginExists(ctx context.Context, login string) (bool, error) {
	return s.storage.IsLoginExists(ctx, login)
}

func (s *userService) Register(ctx context.Context, login, password string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	user := &models.User{
		ID:       id.String(),
		Login:    login,
		Password: string(hash),
	}
	if err = s.storage.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, login, password string) (*models.User, error) {
	user, err := s.storage.FindUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}
	return user, nil
}
