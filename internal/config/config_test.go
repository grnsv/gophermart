package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()

	assert.Equal(t, "localhost:8080", cfg.RunAddress)
	assert.Equal(t, "postgresql://user:password@localhost:5432/dbname?sslmode=disable", cfg.DatabaseURI)
	assert.Equal(t, "http://localhost:8081", cfg.AccrualSystemAddress)
	assert.Equal(t, "secret", cfg.JWTSecret)
}

func TestParseEnvVariables(t *testing.T) {
	os.Setenv("RUN_ADDRESS", "127.0.0.1:9090")
	os.Setenv("DATABASE_URI", "postgresql://test:test@localhost:5432/testdb?sslmode=disable")
	os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "http://localhost:8082")
	os.Setenv("JWT_SECRET", "testsecret")
	defer func() {
		os.Unsetenv("RUN_ADDRESS")
		os.Unsetenv("DATABASE_URI")
		os.Unsetenv("ACCRUAL_SYSTEM_ADDRESS")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg, err := New()
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1:9090", cfg.RunAddress)
	assert.Equal(t, "postgresql://test:test@localhost:5432/testdb?sslmode=disable", cfg.DatabaseURI)
	assert.Equal(t, "http://localhost:8082", cfg.AccrualSystemAddress)
	assert.Equal(t, "testsecret", cfg.JWTSecret)
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "ValidConfig",
			config: &Config{
				RunAddress:           "127.0.0.1:8080",
				DatabaseURI:          "postgresql://user:password@localhost:5432/dbname?sslmode=disable",
				AccrualSystemAddress: "http://localhost:8081",
				JWTSecret:            "secret",
			},
			wantErr: false,
		},
		{
			name: "InvalidRunAddress",
			config: &Config{
				RunAddress:           "invalid_address",
				DatabaseURI:          "postgresql://user:password@localhost:5432/dbname?sslmode=disable",
				AccrualSystemAddress: "http://localhost:8081",
				JWTSecret:            "secret",
			},
			wantErr: true,
		},
		{
			name: "EmptyDatabaseURI",
			config: &Config{
				RunAddress:           "127.0.0.1:8080",
				DatabaseURI:          "",
				AccrualSystemAddress: "http://localhost:8081",
				JWTSecret:            "secret",
			},
			wantErr: true,
		},
		{
			name: "EmptyAccrualSystemAddress",
			config: &Config{
				RunAddress:           "127.0.0.1:8080",
				DatabaseURI:          "postgresql://user:password@localhost:5432/dbname?sslmode=disable",
				AccrualSystemAddress: "",
				JWTSecret:            "secret",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWithOptions(t *testing.T) {
	cfg, err := New(
		WithRunAddress("0.0.0.0:8081"),
		WithDatabaseURI("postgresql://custom:custom@localhost:5432/customdb?sslmode=disable"),
		WithAccrualSystemAddress("http://localhost:8083"),
		WithJWTSecret("customsecret"),
	)
	assert.NoError(t, err)
	assert.Equal(t, "0.0.0.0:8081", cfg.RunAddress)
	assert.Equal(t, "postgresql://custom:custom@localhost:5432/customdb?sslmode=disable", cfg.DatabaseURI)
	assert.Equal(t, "http://localhost:8083", cfg.AccrualSystemAddress)
	assert.Equal(t, "customsecret", cfg.JWTSecret)
}
