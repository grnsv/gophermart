package responses

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Data any `json:"data"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type ErrorsResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func NewErrorsResponse(message string, err error) ErrorsResponse {
	resp := ErrorsResponse{
		Message: message,
		Errors:  make(map[string]string),
	}
	for _, err := range err.(validator.ValidationErrors) {
		resp.Errors[err.Field()] = err.Tag()
	}
	return resp
}

func WriteJSON[T Response | ErrorResponse | ErrorsResponse](w http.ResponseWriter, statusCode int, v T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}
