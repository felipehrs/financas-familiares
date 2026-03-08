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
	categoriaRepo := repository.NewCategoriaRepository(db)
	cartaoRepo := repository.NewCartaoCreditoRepository(db)

	// Inicializar serviços
	authService := service.NewAuthService(authRepo, cfg.JWT.Secret)
	membroSvc := service.NewMembroService(membroRepo)
	categoriaSvc := service.NewCategoriaService(categoriaRepo)
	cartaoSvc := service.NewCartaoCreditoService(cartaoRepo)

	// Inicializar handlers
	authHandler := handler.NewAuthHandler(authService)
	membroHandler := handler.NewMembroHandler(membroSvc)
	categoriaHandler := handler.NewCategoriaHandler(categoriaSvc)
	cartaoHandler := handler.NewCartaoCreditoHandler(cartaoSvc)

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

			// Categorias
			protected.GET("/categorias", categoriaHandler.Listar)
			protected.POST("/categorias", categoriaHandler.Criar)
			protected.PUT("/categorias/:id", categoriaHandler.Atualizar)
			protected.DELETE("/categorias/:id", categoriaHandler.Excluir)

			// Cartões de crédito
			protected.GET("/cartoes", cartaoHandler.Listar)
			protected.POST("/cartoes", cartaoHandler.Criar)
			protected.GET("/cartoes/:id", cartaoHandler.BuscarPorID)
			protected.PUT("/cartoes/:id", cartaoHandler.Atualizar)
			protected.PATCH("/cartoes/:id/inativar", cartaoHandler.Inativar)
		}
	}

	addr := ":" + cfg.Server.Port
	log.Printf("servidor iniciando na porta %s", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
