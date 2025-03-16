package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/grnsv/gophermart/internal/models"
	"github.com/grnsv/gophermart/internal/storage"
)

type orderService struct {
	storage storage.OrderRepository
}

func NewOrderService(storage storage.OrderRepository) OrderService {
	return &orderService{storage: storage}
}

func (s *orderService) UploadOrder(ctx context.Context, userID, orderID string) error {
	id, err := strconv.Atoi(orderID)
	if err != nil {
		return err
	}

	order, err := s.storage.FindOrderByID(ctx, id)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return err
	}
	if order != nil {
		return &OrderAlreadyExistsError{UserID: order.UserID}
	}

	return s.storage.CreateOrder(ctx, &models.Order{
		ID:     id,
		UserID: userID,
		Status: models.StatusNew,
	})
}

func (s *orderService) GetOrders(ctx context.Context, userID string) ([]*models.Order, error) {
	return s.storage.GetOrdersByUserID(ctx, userID)
}
