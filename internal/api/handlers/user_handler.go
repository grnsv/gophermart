package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/grnsv/gophermart/internal/api/requests"
	"github.com/grnsv/gophermart/internal/api/responses"
	"github.com/grnsv/gophermart/internal/logger"
	"github.com/grnsv/gophermart/internal/services"
)

type UserHandler struct {
	logger      logger.Logger
	userService services.UserService
	jwtService  services.JWTService
}

func NewUserHandler(logger logger.Logger, userService services.UserService, jwtService services.JWTService) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
		jwtService:  jwtService,
	}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req requests.RegisterRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Infoln(err)
		responses.WriteJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Message: fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	validate := requests.NewValidator()
	if err := validate.Struct(req); err != nil {
		h.logger.Infoln(err)
		responses.WriteJSON(w, http.StatusBadRequest, responses.NewErrorsResponse("Invalid request body", err))
		return
	}

	exists, err := h.userService.IsLoginExists(r.Context(), req.Login)
	if err != nil {
		h.logger.Errorln(err)
		responses.WriteJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Failed to register user",
		})
		return
	}
	if exists {
		responses.WriteJSON(w, http.StatusConflict, responses.ErrorResponse{
			Message: fmt.Sprintf("Username is already taken: %s", req.Login),
		})
		return
	}

	user, err := h.userService.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		h.logger.Errorln(err)
		responses.WriteJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Failed to register user",
		})
		return
	}

	cookie, err := h.jwtService.BuildCookie(user.ID)
	if err != nil {
		h.logger.Errorln(err)
		responses.WriteJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Server error",
		})
		return
	}

	http.SetCookie(w, cookie)
	responses.WriteJSON(w, http.StatusOK, responses.Response{
		Data: user,
	})
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req requests.LoginRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Infoln(err)
		responses.WriteJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Message: fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	validate := requests.NewValidator()
	if err := validate.Struct(req); err != nil {
		h.logger.Infoln(err)
		responses.WriteJSON(w, http.StatusBadRequest, responses.NewErrorsResponse("Invalid request body", err))
		return
	}

	user, err := h.userService.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			responses.WriteJSON(w, http.StatusUnauthorized, responses.ErrorResponse{
				Message: "Login failed",
			})
			return
		}
		h.logger.Errorln(err)
		responses.WriteJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Login failed",
		})
		return
	}

	cookie, err := h.jwtService.BuildCookie(user.ID)
	if err != nil {
		h.logger.Errorln(err)
		responses.WriteJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Server error",
		})
		return
	}

	http.SetCookie(w, cookie)
	responses.WriteJSON(w, http.StatusOK, responses.Response{
		Data: user,
	})
}

func (h *UserHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
}

func (h *UserHandler) WithdrawPoints(w http.ResponseWriter, r *http.Request) {
}
