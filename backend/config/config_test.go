package config_test

import (
	"os"
	"testing"

	"github.com/felipehrs/financas-familiares/backend/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad_CORSAllowedOrigins(t *testing.T) {
	tests := []struct {
		name            string
		envValue        string
		expectedOrigins []string
	}{
		{
			name:            "múltiplas origens separadas por vírgula",
			envValue:        "http://localhost:5173,http://localhost:3000,https://app.vercel.app",
			expectedOrigins: []string{"http://localhost:5173", "http://localhost:3000", "https://app.vercel.app"},
		},
		{
			name:            "origem única",
			envValue:        "http://localhost:5173",
			expectedOrigins: []string{"http://localhost:5173"},
		},
		{
			name:            "origens com espaços extras (deve trimmar)",
			envValue:        "http://localhost:5173 , http://localhost:3000 , https://app.vercel.app",
			expectedOrigins: []string{"http://localhost:5173", "http://localhost:3000", "https://app.vercel.app"},
		},
		{
			name:            "env var vazia (deve usar default)",
			envValue:        "",
			expectedOrigins: []string{"http://localhost:5173"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup: configurar env vars
			if tt.envValue != "" {
				os.Setenv("CORS_ALLOWED_ORIGINS", tt.envValue)
			} else {
				os.Unsetenv("CORS_ALLOWED_ORIGINS")
			}
			defer os.Unsetenv("CORS_ALLOWED_ORIGINS")

			// Env vars obrigatórias para Load() não falhar
			os.Setenv("DATABASE_URL", "postgres://test")
			os.Setenv("JWT_SECRET", "test-secret")
			defer os.Unsetenv("DATABASE_URL")
			defer os.Unsetenv("JWT_SECRET")

			// Execute
			cfg, err := config.Load()

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOrigins, cfg.CORS.AllowedOrigins)
		})
	}
}
