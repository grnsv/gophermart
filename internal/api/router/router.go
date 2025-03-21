package router

import (
	"net/http"

	"github.com/grnsv/gophermart/internal/api/handlers"
	"github.com/grnsv/gophermart/internal/api/middlewares"
	"github.com/grnsv/gophermart/internal/logger"
	"github.com/grnsv/gophermart/internal/services"
)

type groupBuilder func(muxConfigurator)
type muxConfigurator func(*http.ServeMux)
type middlewareFunc func(http.Handler) http.Handler
type groupOption func(*groupConfig)
type groupConfig struct {
	prefix      string
	middlewares []middlewareFunc
}

func withPrefix(prefix string) groupOption {
	return func(cfg *groupConfig) {
		cfg.prefix = prefix
	}
}

func withMiddlewares(middlewares ...middlewareFunc) groupOption {
	return func(cfg *groupConfig) {
		cfg.middlewares = middlewares
	}
}

func buildMux(fn muxConfigurator) http.Handler {
	mux := http.NewServeMux()
	fn(mux)
	return mux
}

func useMiddlewares(handler http.Handler, middlewares []middlewareFunc) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func mountWithPrefix(mux *http.ServeMux, prefix string, handler http.Handler) {
	if prefix == "/" || prefix == "" {
		mux.Handle("/", handler)
	} else {
		mux.Handle(prefix+"/", http.StripPrefix(prefix, handler))
	}
}

func group(mux *http.ServeMux, opts ...groupOption) groupBuilder {
	cfg := &groupConfig{prefix: "/"}
	for _, opt := range opts {
		opt(cfg)
	}
	return func(fn muxConfigurator) {
		handler := buildMux(fn)
		handler = useMiddlewares(handler, cfg.middlewares)
		mountWithPrefix(mux, cfg.prefix, handler)
	}
}

func NewRouter(
	log logger.Logger,
	authHandler *handlers.AuthHandler,
	protectedHandler *handlers.ProtectedHandler,
	jwtService services.JWTService,
) http.Handler {
	return buildMux(func(mux *http.ServeMux) {
		group(mux,
			withPrefix("/api/user"),
			withMiddlewares(
				middlewares.WithLogging(log),
				middlewares.WithCompressing(log),
			),
		)(func(mux *http.ServeMux) {
			mux.HandleFunc("POST /register", authHandler.RegisterUser)
			mux.HandleFunc("POST /login", authHandler.LoginUser)

			group(mux,
				withMiddlewares(
					middlewares.Authenticate(log, jwtService),
				),
			)(func(mux *http.ServeMux) {
				mux.HandleFunc("POST /orders", protectedHandler.UploadOrder)
				mux.HandleFunc("GET /orders", protectedHandler.GetOrders)
				mux.HandleFunc("GET /balance", protectedHandler.GetBalance)
				mux.HandleFunc("POST /balance/withdraw", protectedHandler.WithdrawPoints)
				mux.HandleFunc("GET /withdrawals", protectedHandler.GetWithdrawals)
			})
		})
	})
}
