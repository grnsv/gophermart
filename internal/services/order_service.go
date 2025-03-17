package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/grnsv/gophermart/internal/logger"
	"github.com/grnsv/gophermart/internal/models"
	"github.com/grnsv/gophermart/internal/storage"
)

type orderService struct {
	logger         logger.Logger
	storage        storage.OrderRepository
	accrualService AccrualService
}

func NewOrderService(
	logger logger.Logger,
	storage storage.OrderRepository,
	accrualService AccrualService,
) OrderService {
	return &orderService{logger: logger, storage: storage, accrualService: accrualService}
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

	order = &models.Order{
		ID:     id,
		UserID: userID,
		Status: models.StatusNew,
	}
	if err = s.storage.CreateOrder(ctx, order); err != nil {
		return err
	}

	go s.updateOrder(order)

	return nil
}

func (s *orderService) updateOrder(order *models.Order) {
	ctx := context.Background()
	order, err := s.accrualService.GetAccrual(ctx, order)
	if err != nil {
		s.logger.Errorln(err)
		return
	}
	err = s.storage.UpdateOrder(ctx, order)
	if err != nil {
		s.logger.Errorln(err)
	}
}

func (s *orderService) GetOrders(ctx context.Context, userID string) ([]*models.Order, error) {
	return s.storage.GetOrdersByUserID(ctx, userID)
}
