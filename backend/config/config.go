package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config agrupa todas as configurações da aplicação.
type Config struct {
	Database  DatabaseConfig
	JWT       JWTConfig
	Server    ServerConfig
	Seed      SeedConfig
	CORS      CORSConfig
	RateLimit RateLimitConfig
}

// DatabaseConfig contém as configurações de conexão com o banco de dados.
type DatabaseConfig struct {
	URL string // DATABASE_URL
}

// JWTConfig contém o segredo para assinatura dos JWTs.
type JWTConfig struct {
	Secret string // JWT_SECRET
}

// ServerConfig contém configurações do servidor HTTP.
type ServerConfig struct {
	Port string // SERVER_PORT, default "8080"
}

// SeedConfig contém os dados dos usuários iniciais para seed.
type SeedConfig struct {
	User1Email string // SEED_USER1_EMAIL
	User1Nome  string // SEED_USER1_NOME
	User1Senha string // SEED_USER1_SENHA (plaintext, será hasheada)
	User2Email string // SEED_USER2_EMAIL
	User2Nome  string // SEED_USER2_NOME
	User2Senha string // SEED_USER2_SENHA
}

// CORSConfig contém as configurações de CORS (Cross-Origin Resource Sharing).
type CORSConfig struct {
	AllowedOrigins []string // CORS_ALLOWED_ORIGINS (separadas por vírgula)
}

// RateLimitConfig contém as configurações de rate limiting.
type RateLimitConfig struct {
	LoginRate string // RATE_LIMIT_LOGIN (formato "X-Y", ex: "5-M")
}

// Load carrega a configuração a partir de variáveis de ambiente e do arquivo .env.
// Variáveis de ambiente têm precedência sobre o arquivo .env.
func Load() (*Config, error) {
	v := viper.New()

	// Tenta carregar .env se existir (ignora erro se não existir)
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	_ = v.ReadInConfig()

	// Permite substituição por variáveis de ambiente reais
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Valores padrão
	v.SetDefault("SERVER_PORT", "8080")

	cfg := &Config{
		Database: DatabaseConfig{
			URL: v.GetString("DATABASE_URL"),
		},
		JWT: JWTConfig{
			Secret: v.GetString("JWT_SECRET"),
		},
		Server: ServerConfig{
			Port: v.GetString("SERVER_PORT"),
		},
		Seed: SeedConfig{
			User1Email: v.GetString("SEED_USER1_EMAIL"),
			User1Nome:  v.GetString("SEED_USER1_NOME"),
			User1Senha: v.GetString("SEED_USER1_SENHA"),
			User2Email: v.GetString("SEED_USER2_EMAIL"),
			User2Nome:  v.GetString("SEED_USER2_NOME"),
			User2Senha: v.GetString("SEED_USER2_SENHA"),
		},
	}

	// Parse CORS allowed origins
	originsStr := v.GetString("CORS_ALLOWED_ORIGINS")
	if originsStr == "" {
		// Default para desenvolvimento local
		originsStr = "http://localhost:5173"
	}

	// Split por vírgula e trim espaços
	cfg.CORS.AllowedOrigins = strings.Split(originsStr, ",")
	for i := range cfg.CORS.AllowedOrigins {
		cfg.CORS.AllowedOrigins[i] = strings.TrimSpace(cfg.CORS.AllowedOrigins[i])
	}

	// Parse rate limiting config
	cfg.RateLimit.LoginRate = v.GetString("RATE_LIMIT_LOGIN")
	if cfg.RateLimit.LoginRate == "" {
		cfg.RateLimit.LoginRate = "5-M" // Default: 5 por minuto
	}

	return cfg, nil
}
