// Package config loads and validates process configuration from environment
// variables. Variable names match .env.example.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds everything the API process needs at startup.
type Config struct {
	DB              DB
	AnthropicAPIKey string
	APIBearerToken  string
	APIPort         int
}

const defaultAPIPort = 8080

// DB holds the connection parameters for the MySQL pool.
type DB struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

// Load reads Config from the current process environment.
func Load() (Config, error) {
	return FromEnv(os.Getenv)
}

// FromEnv reads Config using getenv, letting callers inject a fake lookup in tests.
func FromEnv(getenv func(string) string) (Config, error) {
	var errs []string

	cfg := Config{
		DB: DB{
			Host:     getenv("DB_HOST"),
			Name:     getenv("DB_NAME"),
			User:     getenv("DB_USER"),
			Password: getenv("DB_PASSWORD"),
		},
		AnthropicAPIKey: getenv("ANTHROPIC_API_KEY"),
		APIBearerToken:  getenv("API_BEARER_TOKEN"),
		APIPort:         defaultAPIPort,
	}

	portRaw := getenv("DB_PORT")
	switch port, err := strconv.Atoi(portRaw); {
	case portRaw == "":
		errs = append(errs, "DB_PORT is required")
	case err != nil:
		errs = append(errs, fmt.Sprintf("DB_PORT must be a number, got %q", portRaw))
	default:
		cfg.DB.Port = port
	}

	if apiPortRaw := getenv("API_PORT"); apiPortRaw != "" {
		port, err := strconv.Atoi(apiPortRaw)
		if err != nil {
			errs = append(errs, fmt.Sprintf("API_PORT must be a number, got %q", apiPortRaw))
		} else {
			cfg.APIPort = port
		}
	}

	for _, req := range []struct {
		name  string
		value string
	}{
		{"DB_HOST", cfg.DB.Host},
		{"DB_NAME", cfg.DB.Name},
		{"DB_USER", cfg.DB.User},
		{"DB_PASSWORD", cfg.DB.Password},
		{"ANTHROPIC_API_KEY", cfg.AnthropicAPIKey},
		{"API_BEARER_TOKEN", cfg.APIBearerToken},
	} {
		if req.value == "" {
			errs = append(errs, req.name+" is required")
		}
	}

	if len(errs) > 0 {
		return Config{}, errors.New("config: " + strings.Join(errs, "; "))
	}
	return cfg, nil
}
