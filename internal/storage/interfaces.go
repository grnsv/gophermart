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
	WithdrawalRepository
}

type UserRepository interface {
	io.Closer
	IsLoginExists(ctx context.Context, login string) (bool, error)
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByLogin(ctx context.Context, login string) (*models.User, error)
}

type OrderRepository interface {
	io.Closer
	CreateOrder(ctx context.Context, order *models.Order) error
	FindOrderByID(ctx context.Context, orderID int) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID string) ([]*models.Order, error)
	UpdateOrder(ctx context.Context, order *models.Order) error
}

type WithdrawalRepository interface {
	io.Closer
	GetBalance(ctx context.Context, userID string) (*models.Balance, error)
	CreateWithdrawal(ctx context.Context, withdrawal *models.Withdrawal) error
	GetWithdrawalsByUserID(ctx context.Context, userID string) ([]*models.Withdrawal, error)
}
