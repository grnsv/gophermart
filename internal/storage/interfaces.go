package storage

import (
	"context"
	"errors"
	"io"

	"github.com/grnsv/gophermart/internal/models"
)

var ErrNotFound = errors.New("not found")

type Storage interface {
	UserRepository
	OrderRepository
}

type UserRepository interface {
	io.Closer
	IsLoginExists(ctx context.Context, login string) (bool, error)
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByLogin(ctx context.Context, login string) (*models.User, error)
	UpdateBalance(userID string, balance float64) error
}

type OrderRepository interface {
	io.Closer
	CreateOrder(ctx context.Context, order *models.Order) error
	FindOrderByID(ctx context.Context, orderID int) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID string) ([]*models.Order, error)
	UpdateOrder(ctx context.Context, order *models.Order) error
}
