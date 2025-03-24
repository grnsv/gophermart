package middlewares

import (
	"context"
	"net/http"

	"github.com/grnsv/gophermart/internal/api/responses"
	"github.com/grnsv/gophermart/internal/logger"
	"github.com/grnsv/gophermart/internal/services"
)

type contextKey string

const UserIDContextKey contextKey = "userID"

func Authenticate(logger logger.Logger, jwtService services.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := jwtService.ParseCookie(r)
			if err != nil {
				logger.Infoln(err)
				unauthorized(w)
				return
			}
			if userID == "" {
				unauthorized(w)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func unauthorized(w http.ResponseWriter) {
	responses.WriteJSON(w, http.StatusUnauthorized, responses.ErrorResponse{
		Message: "Unauthorized",
	})
}
