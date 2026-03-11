package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/felipehrs/financas-familiares/backend/config"
	"github.com/felipehrs/financas-familiares/backend/internal/handler"
	"github.com/felipehrs/financas-familiares/backend/internal/middleware"
	"github.com/felipehrs/financas-familiares/backend/internal/repository"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/felipehrs/financas-familiares/backend/migrations"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
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
	defer func() { _ = db.Close() }()

	log.Println("conexão com banco de dados estabelecida")

	// Rodar migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("erro ao rodar migrations: %v", err)
	}

	// Rodar seeds (SeedUsuarios já chama SeedCategorias internamente)
	if err := repository.SeedUsuarios(db, &cfg.Seed); err != nil {
		log.Printf("aviso: erro no seed de usuários: %v", err)
	}

	// Inicializar repositórios
	authRepo := repository.NewAuthRepository(db)
	familiaRepo := repository.NewFamiliaRepository(db)
	membroRepo := repository.NewMembroRepository(db)
	categoriaRepo := repository.NewCategoriaRepository(db)
	cartaoRepo := repository.NewCartaoCreditoRepository(db)
	despesaRepo := repository.NewDespesaCartaoRepository(db)
	rendaFixaRepo := repository.NewRendaFixaRepository(db)
	assinaturaRepo := repository.NewAssinaturaRepository(db)
	contaFixaRepo := repository.NewContaFixaRepository(db)
	despesaGeralRepo := repository.NewDespesaGeralRepository(db)
	rendaVariavelRepo := repository.NewRendaVariavelRepository(db)
	rendaExtraRepo := repository.NewRendaExtraRepository(db)
	rendimentoRepo := repository.NewRendimentoInvestimentoRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)

	// Inicializar serviços
	authService := service.NewAuthService(authRepo, familiaRepo, cfg.JWT.Secret)
	membroSvc := service.NewMembroService(membroRepo)
	categoriaSvc := service.NewCategoriaService(categoriaRepo)
	cartaoSvc := service.NewCartaoCreditoService(cartaoRepo)
	despesaSvc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)
	rendaFixaSvc := service.NewRendaFixaService(rendaFixaRepo)
	assinaturaSvc := service.NewAssinaturaService(assinaturaRepo)
	contaFixaSvc := service.NewContaFixaService(contaFixaRepo)
	despesaGeralSvc := service.NewDespesaGeralService(despesaGeralRepo)
	rendaVariavelSvc := service.NewRendaVariavelService(rendaVariavelRepo)
	rendaExtraSvc := service.NewRendaExtraService(rendaExtraRepo)
	rendimentoSvc := service.NewRendimentoInvestimentoService(rendimentoRepo)
	dashboardSvc := service.NewDashboardService(
		rendaFixaRepo,
		despesaRepo,
		rendaVariavelRepo,
		rendaExtraRepo,
		rendimentoRepo,
		assinaturaRepo,
		contaFixaRepo,
		despesaGeralRepo,
		dashboardRepo,
	)
	rendaHistoricoSvc := service.NewRendaHistoricoService(rendaFixaRepo, rendaVariavelRepo, rendaExtraRepo, rendimentoRepo)

	// Inicializar handlers
	authHandler := handler.NewAuthHandler(authService)
	membroHandler := handler.NewMembroHandler(membroSvc)
	categoriaHandler := handler.NewCategoriaHandler(categoriaSvc)
	cartaoHandler := handler.NewCartaoCreditoHandler(cartaoSvc)
	despesaHandler := handler.NewDespesaCartaoHandler(despesaSvc)
	rendaFixaHandler := handler.NewRendaFixaHandler(rendaFixaSvc)
	assinaturaHandler := handler.NewAssinaturaHandler(assinaturaSvc)
	contaFixaHandler := handler.NewContaFixaHandler(contaFixaSvc)
	despesaGeralHandler := handler.NewDespesaGeralHandler(despesaGeralSvc)
	rendaVariavelHandler := handler.NewRendaVariavelHandler(rendaVariavelSvc)
	rendaExtraHandler := handler.NewRendaExtraHandler(rendaExtraSvc)
	rendimentoHandler := handler.NewRendimentoInvestimentoHandler(rendimentoSvc)
	dashboardHandler := handler.NewDashboardHandler(dashboardSvc)
	rendaHistoricoHandler := handler.NewRendaHistoricoHandler(rendaHistoricoSvc)

	// Configurar rotas
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

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
			protected.GET("/rendas-fixas/vigentes", rendaFixaHandler.ListarVigentesPorMes)
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

			// Despesas Gerais
			protected.GET("/despesas-gerais", despesaGeralHandler.Listar)
			protected.POST("/despesas-gerais", despesaGeralHandler.Criar)
			protected.GET("/despesas-gerais/:id", despesaGeralHandler.BuscarPorID)
			protected.PUT("/despesas-gerais/:id", despesaGeralHandler.Atualizar)
			protected.DELETE("/despesas-gerais/:id", despesaGeralHandler.Excluir)

			// Rendas Variáveis
			protected.GET("/rendas-variaveis", rendaVariavelHandler.Listar)
			protected.POST("/rendas-variaveis", rendaVariavelHandler.Criar)
			protected.GET("/rendas-variaveis/:id", rendaVariavelHandler.BuscarPorID)
			protected.PUT("/rendas-variaveis/:id", rendaVariavelHandler.Atualizar)
			protected.DELETE("/rendas-variaveis/:id", rendaVariavelHandler.Excluir)

			// Rendas Extras
			protected.GET("/rendas-extras", rendaExtraHandler.Listar)
			protected.POST("/rendas-extras", rendaExtraHandler.Criar)
			protected.GET("/rendas-extras/:id", rendaExtraHandler.BuscarPorID)
			protected.PUT("/rendas-extras/:id", rendaExtraHandler.Atualizar)
			protected.DELETE("/rendas-extras/:id", rendaExtraHandler.Excluir)

			// Rendimentos de Investimento
			protected.GET("/rendimentos-investimento", rendimentoHandler.Listar)
			protected.POST("/rendimentos-investimento", rendimentoHandler.Criar)
			protected.GET("/rendimentos-investimento/:id", rendimentoHandler.BuscarPorID)
			protected.PUT("/rendimentos-investimento/:id", rendimentoHandler.Atualizar)
			protected.DELETE("/rendimentos-investimento/:id", rendimentoHandler.Excluir)

			// Dashboard
			protected.GET("/dashboard/resumo", dashboardHandler.ResumoMensal)
			protected.GET("/dashboard/categorias", dashboardHandler.DespesasPorCategoria)
			protected.GET("/dashboard/evolucao", dashboardHandler.EvolucaoMensal)
			protected.GET("/dashboard/projecao", dashboardHandler.Projecao)

			// Histórico de rendas
			protected.GET("/rendas/historico", rendaHistoricoHandler.Historico)
		}
	}

	addr := ":" + cfg.Server.Port
	log.Printf("servidor iniciando na porta %s", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}

func runMigrations(db *sqlx.DB) error {
	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("migrations aplicadas com sucesso")
	return nil
}
