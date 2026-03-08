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
	despesaRepo := repository.NewDespesaCartaoRepository(db)
	rendaFixaRepo := repository.NewRendaFixaRepository(db)
	assinaturaRepo := repository.NewAssinaturaRepository(db)
	contaFixaRepo := repository.NewContaFixaRepository(db)

	// Inicializar serviços
	authService := service.NewAuthService(authRepo, cfg.JWT.Secret)
	membroSvc := service.NewMembroService(membroRepo)
	categoriaSvc := service.NewCategoriaService(categoriaRepo)
	cartaoSvc := service.NewCartaoCreditoService(cartaoRepo)
	despesaSvc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)
	rendaFixaSvc := service.NewRendaFixaService(rendaFixaRepo)
	assinaturaSvc := service.NewAssinaturaService(assinaturaRepo)
	contaFixaSvc := service.NewContaFixaService(contaFixaRepo)
	dashboardSvc := service.NewDashboardService(rendaFixaRepo, despesaRepo)

	// Inicializar handlers
	authHandler := handler.NewAuthHandler(authService)
	membroHandler := handler.NewMembroHandler(membroSvc)
	categoriaHandler := handler.NewCategoriaHandler(categoriaSvc)
	cartaoHandler := handler.NewCartaoCreditoHandler(cartaoSvc)
	despesaHandler := handler.NewDespesaCartaoHandler(despesaSvc)
	rendaFixaHandler := handler.NewRendaFixaHandler(rendaFixaSvc)
	assinaturaHandler := handler.NewAssinaturaHandler(assinaturaSvc)
	contaFixaHandler := handler.NewContaFixaHandler(contaFixaSvc)
	dashboardHandler := handler.NewDashboardHandler(dashboardSvc)

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

			// Despesas de cartão de crédito
			protected.GET("/cartoes/:id/despesas", despesaHandler.ListarPorCartao)
			protected.GET("/cartoes/:id/despesas/fatura", despesaHandler.ListarPorFatura)
			protected.POST("/cartoes/:id/despesas", despesaHandler.Criar)
			protected.DELETE("/despesas/:id", despesaHandler.Excluir)

			// Rendas fixas
			protected.GET("/rendas-fixas", rendaFixaHandler.Listar)
			protected.POST("/rendas-fixas", rendaFixaHandler.Criar)
			protected.GET("/rendas-fixas/:id", rendaFixaHandler.BuscarPorID)
			protected.PUT("/rendas-fixas/:id", rendaFixaHandler.Atualizar)
			protected.PATCH("/rendas-fixas/:id/inativar", rendaFixaHandler.Inativar)

			// Assinaturas
			protected.GET("/assinaturas", assinaturaHandler.Listar)
			protected.POST("/assinaturas", assinaturaHandler.Criar)
			protected.GET("/assinaturas/:id", assinaturaHandler.BuscarPorID)
			protected.PUT("/assinaturas/:id", assinaturaHandler.Atualizar)
			protected.PATCH("/assinaturas/:id/status", assinaturaHandler.AlterarStatus)

			// Contas Fixas
			protected.GET("/contas-fixas", contaFixaHandler.Listar)
			protected.POST("/contas-fixas", contaFixaHandler.Criar)
			protected.GET("/contas-fixas/:id", contaFixaHandler.BuscarPorID)
			protected.PUT("/contas-fixas/:id", contaFixaHandler.Atualizar)
			protected.PATCH("/contas-fixas/:id/ativo", contaFixaHandler.AlterarAtivo)

			// Dashboard
			protected.GET("/dashboard/resumo", dashboardHandler.ResumoMensal)
		}
	}

	addr := ":" + cfg.Server.Port
	log.Printf("servidor iniciando na porta %s", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
