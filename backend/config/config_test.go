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
				err := os.Setenv("CORS_ALLOWED_ORIGINS", tt.envValue)
				assert.NoError(t, err)
			} else {
				err := os.Unsetenv("CORS_ALLOWED_ORIGINS")
				assert.NoError(t, err)
			}
			defer func() {
				err := os.Unsetenv("CORS_ALLOWED_ORIGINS")
				assert.NoError(t, err)
			}()

			// Env vars obrigatórias para Load() não falhar
			err := os.Setenv("DATABASE_URL", "postgres://test")
			assert.NoError(t, err)
			err = os.Setenv("JWT_SECRET", "test-secret")
			assert.NoError(t, err)
			defer func() {
				err := os.Unsetenv("DATABASE_URL")
				assert.NoError(t, err)
			}()
			defer func() {
				err := os.Unsetenv("JWT_SECRET")
				assert.NoError(t, err)
			}()

			// Execute
			cfg, err := config.Load()

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOrigins, cfg.CORS.AllowedOrigins)
		})
	}
}
