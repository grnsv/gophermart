package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/grnsv/gophermart/internal/api/middlewares"
	"github.com/grnsv/gophermart/internal/api/responses"
	"github.com/grnsv/gophermart/internal/logger"
	"github.com/grnsv/gophermart/internal/services"
)

var ErrUserIDNotFound = errors.New("user ID not found in context")

type OrderHandler struct {
	logger       logger.Logger
	orderService services.OrderService
	validator    services.Validator
}

func NewOrderHandler(logger logger.Logger, service services.OrderService, validator services.Validator) *OrderHandler {
	return &OrderHandler{logger: logger, orderService: service, validator: validator}
}

func (h *OrderHandler) getUserID(w http.ResponseWriter, r *http.Request) (string, error) {
	userID, ok := r.Context().Value(middlewares.UserIDContextKey).(string)
	if !ok {
		h.logger.Errorln(ErrUserIDNotFound)
		responses.WriteJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Server error",
		})
		return userID, ErrUserIDNotFound
	}
	return userID, nil
}

func (h *OrderHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userID, err := h.getUserID(w, r)
	if err != nil {
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Infoln(err)
		responses.WriteJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Message: fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	orderID := strings.ReplaceAll(string(body), " ", "")
	if orderID == "" || !h.validator.IsValid(orderID) {
		responses.WriteJSON(w, http.StatusUnprocessableEntity, responses.ErrorResponse{
			Message: fmt.Sprintf("Invalid order number: '%s'", orderID),
		})
		return
	}

	if err := h.orderService.UploadOrder(r.Context(), userID, orderID); err != nil {
		if e, ok := err.(*services.OrderAlreadyExistsError); ok {
			if e.UserID == userID {
				responses.WriteJSON(w, http.StatusOK, responses.Response{
					Data: "Order already exists",
				})
				return
			} else {
				responses.WriteJSON(w, http.StatusConflict, responses.ErrorResponse{
					Message: "Order ID is not unique",
				})
				return
			}
		}

		h.logger.Errorln(err)
		responses.WriteJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Failed to upload order",
		})
		return
	}

	responses.WriteJSON(w, http.StatusAccepted, responses.Response{
		Data: "Order uploaded",
	})
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(w, r)
	if err != nil {
		return
	}
	orders, err := h.orderService.GetOrders(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
}
