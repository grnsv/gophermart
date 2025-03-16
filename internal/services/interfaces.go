package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/grnsv/gophermart/internal/models"
)

var ErrUnauthorized = errors.New("unauthorized")

type OrderAlreadyExistsError struct {
	UserID string
}

func (e *OrderAlreadyExistsError) Error() string {
	return fmt.Sprintf("order already uploaded by user: %s", e.UserID)
}

type UserService interface {
	IsLoginExists(ctx context.Context, login string) (bool, error)
	Register(ctx context.Context, login, password string) (*models.User, error)
	Login(ctx context.Context, login, password string) (*models.User, error)
}

type OrderService interface {
	UploadOrder(ctx context.Context, userID, orderID string) error
	GetOrders(ctx context.Context, userID string) ([]*models.Order, error)
}

type JWTService interface {
	ParseCookie(r *http.Request) (string, error)
	BuildCookie(userID string) (*http.Cookie, error)
}

type Validator interface {
	IsValid(number string) bool
}
