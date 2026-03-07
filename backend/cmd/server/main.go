package main

import (
	"log"
	"net/http"

	"github.com/felipehrs/financas-familiares/backend/config"
	"github.com/felipehrs/financas-familiares/backend/internal/handler"
	"github.com/felipehrs/financas-familiares/backend/internal/middleware"
	"github.com/felipehrs/financas-familiares/backend/internal/repository"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("erro ao carregar configuração: %v", err)
	}

	if cfg.JWT.Secret == "" {
		log.Fatal("JWT_SECRET não configurado")
	}

	// Conectar ao banco de dados
	db, err := sqlx.Connect("pgx", cfg.Database.URL)
	if err != nil {
		log.Fatalf("erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	log.Println("conexão com banco de dados estabelecida")

	// Rodar seeds
	if err := repository.SeedUsuarios(db, &cfg.Seed); err != nil {
		log.Printf("aviso: erro no seed de usuários: %v", err)
	}
	if err := repository.SeedCategorias(db); err != nil {
		log.Printf("aviso: erro no seed de categorias: %v", err)
	}

	// Inicializar repositórios
	authRepo := repository.NewAuthRepository(db)
	membroRepo := repository.NewMembroRepository(db)

	// Inicializar serviços
	authService := service.NewAuthService(authRepo, cfg.JWT.Secret)
	membroSvc := service.NewMembroService(membroRepo)

	// Inicializar handlers
	authHandler := handler.NewAuthHandler(authService)
	membroHandler := handler.NewMembroHandler(membroSvc)

	// Configurar rotas
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		// Rotas públicas de autenticação
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.Refresh)
		}

		// Rotas protegidas por JWT
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// Membros da família
			protected.GET("/membros", membroHandler.Listar)
			protected.POST("/membros", membroHandler.Criar)
			protected.GET("/membros/:id", membroHandler.BuscarPorID)
			protected.PUT("/membros/:id", membroHandler.Atualizar)
			protected.PATCH("/membros/:id/inativar", membroHandler.Inativar)
		}
	}

	addr := ":" + cfg.Server.Port
	log.Printf("servidor iniciando na porta %s", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
