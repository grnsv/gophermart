package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/grnsv/gophermart/internal/models"
)

var ErrUnauthorized = errors.New("unauthorized")

type UploadOrderStatus int

const (
	OrderUploaded UploadOrderStatus = iota
	OrderAlreadyExistsForUser
	OrderAlreadyExistsForAnotherUser
)

type UserService interface {
	IsLoginExists(ctx context.Context, login string) (bool, error)
	Register(ctx context.Context, login, password string) (*models.User, error)
	Login(ctx context.Context, login, password string) (*models.User, error)
}

type OrderService interface {
	UploadOrder(ctx context.Context, userID, orderID string) (UploadOrderStatus, error)
	GetOrders(ctx context.Context, userID string) ([]*models.Order, error)
}

type JWTService interface {
	ParseCookie(r *http.Request) (string, error)
	BuildCookie(userID string) (*http.Cookie, error)
}

type Validator interface {
	IsValid(number string) bool
}

type AccrualService interface {
	GetAccrual(ctx context.Context, order *models.Order) (*models.Order, error)
}
