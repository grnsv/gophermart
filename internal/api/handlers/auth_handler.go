package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/grnsv/gophermart/internal/api/requests"
	"github.com/grnsv/gophermart/internal/api/responses"
	"github.com/grnsv/gophermart/internal/logger"
	"github.com/grnsv/gophermart/internal/services"
)

type AuthHandler struct {
	baseHandler
	userService services.UserService
	jwtService  services.JWTService
}

func NewAuthHandler(logger logger.Logger, userService services.UserService, jwtService services.JWTService) *AuthHandler {
	return &AuthHandler{
		baseHandler: newBaseHandler(logger),
		userService: userService,
		jwtService:  jwtService,
	}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req requests.RegisterRequest
	if !h.decodeAndValidate(w, r, &req) {
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

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req requests.LoginRequest
	if !h.decodeAndValidate(w, r, &req) {
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
