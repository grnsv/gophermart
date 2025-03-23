package requests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterRequestValidation(t *testing.T) {
	validate := NewValidator()

	tests := []struct {
		name    string
		request RegisterRequest
		isValid bool
	}{
		{"ValidRequest", RegisterRequest{Login: "testuser", Password: "password123"}, true},
		{"MissingLogin", RegisterRequest{Password: "password123"}, false},
		{"ShortLogin", RegisterRequest{Login: "usr", Password: "password123"}, false},
		{"MissingPassword", RegisterRequest{Login: "testuser"}, false},
		{"ShortPassword", RegisterRequest{Login: "testuser", Password: "short"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestLoginRequestValidation(t *testing.T) {
	validate := NewValidator()

	tests := []struct {
		name    string
		request LoginRequest
		isValid bool
	}{
		{"ValidRequest", LoginRequest{Login: "testuser", Password: "password123"}, true},
		{"MissingLogin", LoginRequest{Password: "password123"}, false},
		{"MissingPassword", LoginRequest{Login: "testuser"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestWithdrawRequestValidation(t *testing.T) {
	validate := NewValidator()

	tests := []struct {
		name    string
		request WithdrawRequest
		isValid bool
	}{
		{"ValidRequest", WithdrawRequest{Order: "1234567890", Sum: 100.50}, true},
		{"MissingOrder", WithdrawRequest{Sum: 100.50}, false},
		{"MissingSum", WithdrawRequest{Order: "1234567890"}, false},
		{"ZeroSum", WithdrawRequest{Order: "1234567890", Sum: 0}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
