package requests

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type RegisterRequest struct {
	Login    string `json:"login" validate:"required,min=4,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func NewValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" {
			return fld.Name
		}
		return name
	})

	return validate
}
