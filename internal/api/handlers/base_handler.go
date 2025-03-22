package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grnsv/gophermart/internal/api/requests"
	"github.com/grnsv/gophermart/internal/api/responses"
	"github.com/grnsv/gophermart/internal/logger"
)

type baseHandler struct {
	logger logger.Logger
}

func newBaseHandler(logger logger.Logger) baseHandler {
	return baseHandler{
		logger: logger,
	}
}

func (h *baseHandler) decodeAndValidate(w http.ResponseWriter, r *http.Request, req any) bool {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		h.logger.Infoln(err)
		responses.WriteJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Message: fmt.Sprintf("Invalid request body: %v", err),
		})
		return false
	}

	if err := requests.NewValidator().Struct(req); err != nil {
		h.logger.Infoln(err)
		responses.WriteJSON(w, http.StatusBadRequest, responses.NewErrorsResponse("Invalid request body", err))
		return false
	}
	return true
}
