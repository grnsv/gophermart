package config

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/caarlos0/env"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTSecret            string `env:"JWT_SECRET"`
}

type Option func(*Config)

func WithRunAddress(url string) Option {
	return func(c *Config) {
		c.RunAddress = url
	}
}

func WithDatabaseURI(url string) Option {
	return func(c *Config) {
		c.DatabaseURI = url
	}
}

func WithAccrualSystemAddress(url string) Option {
	return func(c *Config) {
		c.AccrualSystemAddress = url
	}
}

func WithJWTSecret(key string) Option {
	return func(c *Config) {
		c.JWTSecret = key
	}
}

func New(opts ...Option) (*Config, error) {
	config := &Config{
		RunAddress:           "localhost:8080",
		DatabaseURI:          "postgresql://user:password@localhost:5432/dbname?sslmode=disable",
		AccrualSystemAddress: "localhost:8081",
		JWTSecret:            "secret",
	}

	flag.StringVar(&config.RunAddress, "a", config.RunAddress, "RUN_ADDRESS")
	flag.StringVar(&config.DatabaseURI, "d", config.DatabaseURI, "DATABASE_URI")
	flag.StringVar(&config.AccrualSystemAddress, "r", config.AccrualSystemAddress, "ACCRUAL_SYSTEM_ADDRESS")

	flag.Parse()

	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	for _, opt := range opts {
		opt(config)
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.RunAddress == "" {
		return fmt.Errorf("run address cannot be empty")
	}
	hp := strings.Split(c.RunAddress, ":")
	if len(hp) != 2 || hp[0] == "" || hp[1] == "" {
		return fmt.Errorf("run address must be in format `host:port`, %s given", c.RunAddress)
	}
	if _, err := net.LookupHost(hp[0]); err != nil {
		return fmt.Errorf("host is invalid or unreachable: %w", err)
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return fmt.Errorf("invalid port: %w", err)
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, %d given", port)
	}

	if c.DatabaseURI == "" {
		return fmt.Errorf("database URI cannot be empty")
	}
	if _, err := url.Parse(c.DatabaseURI); err != nil {
		return fmt.Errorf("database URI has invalid format: %w", err)
	}

	if c.AccrualSystemAddress == "" {
		return fmt.Errorf("accrual system address cannot be empty")
	}
	if _, err := url.Parse(c.DatabaseURI); err != nil {
		return fmt.Errorf("accrual system address has invalid format: %w", err)
	}

	return nil
}
