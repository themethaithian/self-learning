package config

import (
	"strings"
	"testing"
)

func TestFromEnv(t *testing.T) {
	validEnv := map[string]string{
		"DB_HOST":           "127.0.0.1",
		"DB_PORT":           "3306",
		"DB_NAME":           "self_learning",
		"DB_USER":           "app",
		"DB_PASSWORD":       "secret",
		"ANTHROPIC_API_KEY": "sk-ant-test",
		"API_BEARER_TOKEN":  "bearer-test",
	}

	tests := []struct {
		name    string
		mutate  func(map[string]string)
		wantErr string
	}{
		{
			name:   "valid config",
			mutate: func(map[string]string) {},
		},
		{
			name:    "missing DB_HOST",
			mutate:  func(env map[string]string) { delete(env, "DB_HOST") },
			wantErr: "DB_HOST is required",
		},
		{
			name:    "missing DB_PORT",
			mutate:  func(env map[string]string) { delete(env, "DB_PORT") },
			wantErr: "DB_PORT is required",
		},
		{
			name:    "non-numeric DB_PORT",
			mutate:  func(env map[string]string) { env["DB_PORT"] = "not-a-number" },
			wantErr: `DB_PORT must be a number, got "not-a-number"`,
		},
		{
			name:    "missing ANTHROPIC_API_KEY",
			mutate:  func(env map[string]string) { delete(env, "ANTHROPIC_API_KEY") },
			wantErr: "ANTHROPIC_API_KEY is required",
		},
		{
			name:    "missing API_BEARER_TOKEN",
			mutate:  func(env map[string]string) { delete(env, "API_BEARER_TOKEN") },
			wantErr: "API_BEARER_TOKEN is required",
		},
		{
			name: "multiple missing fields reported together",
			mutate: func(env map[string]string) {
				delete(env, "DB_USER")
				delete(env, "DB_PASSWORD")
			},
			wantErr: "DB_USER is required; DB_PASSWORD is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := make(map[string]string, len(validEnv))
			for k, v := range validEnv {
				env[k] = v
			}
			tt.mutate(env)
			getenv := func(key string) string { return env[key] }

			cfg, err := FromEnv(getenv)

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("FromEnv() unexpected error: %v", err)
				}
				if cfg.DB.Host != env["DB_HOST"] || cfg.DB.Port != 3306 || cfg.DB.Name != env["DB_NAME"] ||
					cfg.DB.User != env["DB_USER"] || cfg.DB.Password != env["DB_PASSWORD"] ||
					cfg.AnthropicAPIKey != env["ANTHROPIC_API_KEY"] || cfg.APIBearerToken != env["API_BEARER_TOKEN"] {
					t.Fatalf("FromEnv() = %+v, fields do not match env", cfg)
				}
				return
			}

			if err == nil {
				t.Fatalf("FromEnv() expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("FromEnv() error = %q, want it to contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}
