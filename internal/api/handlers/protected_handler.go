package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/grnsv/gophermart/internal/api/middlewares"
	"github.com/grnsv/gophermart/internal/api/requests"
	"github.com/grnsv/gophermart/internal/api/responses"
	"github.com/grnsv/gophermart/internal/logger"
	"github.com/grnsv/gophermart/internal/models"
	"github.com/grnsv/gophermart/internal/services"
	"github.com/grnsv/gophermart/internal/storage"
)

var ErrUserIDNotFound = errors.New("user ID not found in context")

type ProtectedHandler struct {
	logger               logger.Logger
	orderService         services.OrderService
	validator            services.Validator
	withdrawalRepository storage.WithdrawalRepository
}

func NewProtectedHandler(
	logger logger.Logger,
	service services.OrderService,
	validator services.Validator,
	withdrawalRepository storage.WithdrawalRepository,
) *ProtectedHandler {
	return &ProtectedHandler{
		logger:               logger,
		orderService:         service,
		validator:            validator,
		withdrawalRepository: withdrawalRepository,
	}
}

func (h *ProtectedHandler) getUserID(w http.ResponseWriter, r *http.Request) (string, error) {
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

func (h *ProtectedHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
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

func (h *ProtectedHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(w, r)
	if err != nil {
		return
	}
	orders, err := h.orderService.GetOrders(r.Context(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		responses.WriteJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Failed to get orders",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func (h *ProtectedHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(w, r)
	if err != nil {
		return
	}
	balance, err := h.withdrawalRepository.GetBalance(r.Context(), userID)
	if err != nil {
		responses.WriteJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Failed to get balance",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
}

func (h *ProtectedHandler) WithdrawPoints(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(w, r)
	if err != nil {
		return
	}

	var req requests.WithdrawRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Infoln(err)
		responses.WriteJSON(w, http.StatusUnprocessableEntity, responses.ErrorResponse{
			Message: fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	if err := requests.NewValidator().Struct(req); err != nil {
		h.logger.Infoln(err)
		responses.WriteJSON(w, http.StatusBadRequest, responses.NewErrorsResponse("Invalid request body", err))
		return
	}

	orderID := strings.ReplaceAll(req.Order, " ", "")
	if orderID == "" || !h.validator.IsValid(orderID) {
		responses.WriteJSON(w, http.StatusUnprocessableEntity, responses.ErrorResponse{
			Message: fmt.Sprintf("Invalid order number: '%s'", orderID),
		})
		return
	}

	balance, err := h.withdrawalRepository.GetBalance(r.Context(), userID)
	if err != nil {
		responses.WriteJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Failed to get balance",
		})
		return
	}

	if balance.Current < req.Sum {
		responses.WriteJSON(w, http.StatusPaymentRequired, responses.ErrorResponse{
			Message: "Insufficient funds",
		})
		return
	}

	id, _ := strconv.Atoi(req.Order)
	if err := h.withdrawalRepository.CreateWithdrawal(r.Context(), &models.Withdrawal{
		UserID:  userID,
		OrderID: id,
		Sum:     req.Sum,
	}); err != nil {
		h.logger.Errorln(err)
		responses.WriteJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Failed to create withdrawal",
		})
		return
	}

	responses.WriteJSON(w, http.StatusOK, responses.Response{
		Data: "Withdrawal successful",
	})
}

func (h *ProtectedHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(w, r)
	if err != nil {
		return
	}
	withdrawals, err := h.withdrawalRepository.GetWithdrawalsByUserID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		responses.WriteJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Failed to get withdrawals",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(withdrawals)
}
