# Status da Implementação — Sprint 9 TT-07 (Data Isolation por familia_id)

**Branch de trabalho:** `feat/sprint9-familia-isolation`
**Data do último ponto:** 2026-03-12

---

## Objetivo

Adicionar isolamento de dados por `familia_id` em todas as camadas do backend (migrations, domain, repository, service, handler) seguindo o plano em `docs/tecnico/03-plano-sprint9-seguranca.md` seção 2.

---

## Padrão de mudança (para referência)

- Toda interface de repositório ganha `familiaID string` como **primeiro parâmetro** em todos os métodos.
- Toda interface de serviço reflete o mesmo.
- Toda implementação de service passa `familiaID` para o repo.
- Mocks de teste atualizam as assinaturas.
- Chamadas de teste passam `"familia-id"` como primeiro argumento.

---

## Progresso por arquivo

### Migrations (CONCLUÍDO)
- [x] `backend/migrations/000017_add_familia_isolation.up.sql` — criado
- [x] `backend/migrations/000017_add_familia_isolation.down.sql` — criado

### Domain (CONCLUÍDO)
- [x] `backend/internal/domain/familia.go` — criado (`Familia` struct, `ErrFamiliaNaoEncontrada`)

### Repository (CONCLUÍDO)
- [x] `backend/internal/repository/familia_repository.go` — criado (`BuscarFamiliaPorUsuario`)
- [x] `backend/internal/repository/membro_repository.go` — `familiaID` em todos os métodos, INSERT inclui `familia_id`, WHERE com `AND familia_id = $N`
- [x] `backend/internal/repository/categoria_repository.go` — idem; `Excluir` sub-queries usam `id=$1, familiaID=$2` em cada branch
- [x] `backend/internal/repository/cartao_credito_repository.go` — `familiaID` em todos os métodos
- [x] `backend/internal/repository/despesa_cartao_repository.go` — `ListarPorFaturaGlobal(familiaID, mes, ano)`: `WHERE familia_id=$1 AND fatura_mes=$2 AND fatura_ano=$3`
- [x] `backend/internal/repository/assinatura_repository.go` — `ListarAtivas(familiaID)`, `AlterarStatus(familiaID, id, status)`
- [x] `backend/internal/repository/conta_fixa_repository.go` — `ListarAtivas(familiaID)`, `AlterarAtivo(familiaID, id, ativa)`
- [x] `backend/internal/repository/despesa_geral_repository.go` — `ListarPorMes(familiaID, mes, ano)`, `Atualizar(familiaID, d)`
- [x] `backend/internal/repository/renda_fixa_repository.go` — `ListarAtivas(familiaID)`, `ListarVigentesPorMes(familiaID, mes, ano)`, `Inativar(familiaID, id)`
- [x] `backend/internal/repository/renda_variavel_repository.go` — `ListarPorMes(familiaID, mes, ano)`
- [x] `backend/internal/repository/renda_extra_repository.go` — `ListarPorMes(familiaID, mes, ano)`
- [x] `backend/internal/repository/rendimento_investimento_repository.go` — `ListarPorMes(familiaID, mes, ano)`
- [x] `backend/internal/repository/dashboard_repository.go` — `DespesasPorCategoria(familiaID, mes, ano)`: $1=familiaID, $2=mes, $3=ano em todo o UNION ALL

### Middleware (CONCLUÍDO)
- [x] `backend/internal/middleware/auth.go` — atualizado (`ValidateAccessToken` agora retorna `usuarioID, familiaID, err`; seta `"familiaID"` no contexto)

### Handler helpers (CONCLUÍDO)
- [x] `backend/internal/handler/helpers.go` — criado (`getFamiliaID(c *gin.Context) (string, bool)`)

### Service — interfaces + implementações (CONCLUÍDO)
- [x] `backend/internal/service/auth_service.go` — `FamiliaRepository` interface; `NewAuthService` recebe `familiaRepo`; `Login` e `RefreshToken` buscam `familiaID`; JWT inclui `familia_id`; `ValidateAccessToken` retorna 3 valores
- [x] `backend/internal/service/membro_service.go`
- [x] `backend/internal/service/categoria_service.go`
- [x] `backend/internal/service/cartao_credito_service.go`
- [x] `backend/internal/service/despesa_cartao_service.go` — inclui `CartaoRepositoryForDespesa.BuscarPorID(familiaID, id)`
- [x] `backend/internal/service/assinatura_service.go`
- [x] `backend/internal/service/conta_fixa_service.go`
- [x] `backend/internal/service/despesa_geral_service.go`
- [x] `backend/internal/service/renda_fixa_service.go`
- [x] `backend/internal/service/renda_variavel_service.go`
- [x] `backend/internal/service/renda_extra_service.go`
- [x] `backend/internal/service/rendimento_investimento_service.go`
- [x] `backend/internal/service/dashboard_service.go` — CONCLUÍDO: interfaces + 4 métodos públicos com `familiaID`
- [x] `backend/internal/service/renda_historico_service.go` — CONCLUÍDO: interfaces + `BuscarHistorico` com `familiaID`

### Service — testes (CONCLUÍDO)
- [x] `backend/internal/service/auth_service_test.go`
- [x] `backend/internal/service/membro_service_test.go`
- [x] `backend/internal/service/categoria_service_test.go`
- [x] `backend/internal/service/cartao_credito_service_test.go`
- [x] `backend/internal/service/despesa_cartao_service_test.go`
- [x] `backend/internal/service/assinatura_service_test.go`
- [x] `backend/internal/service/conta_fixa_service_test.go`
- [x] `backend/internal/service/despesa_geral_service_test.go`
- [x] `backend/internal/service/renda_fixa_service_test.go`
- [x] `backend/internal/service/renda_variavel_service_test.go`
- [x] `backend/internal/service/renda_extra_service_test.go`
- [x] `backend/internal/service/rendimento_investimento_service_test.go` — CONCLUÍDO: mocks e chamadas de teste atualizados
- [x] `backend/internal/service/dashboard_service_test.go` — CONCLUÍDO: 9 mocks + chamadas de teste atualizados
- [x] `backend/internal/service/renda_historico_service_test.go` — CONCLUÍDO: 8 mocks + chamadas de teste atualizados

### Handlers (CONCLUÍDO)
- [x] `backend/internal/handler/membro_handler.go`
- [x] `backend/internal/handler/categoria_handler.go`
- [x] `backend/internal/handler/cartao_credito_handler.go`
- [x] `backend/internal/handler/despesa_cartao_handler.go`
- [x] `backend/internal/handler/assinatura_handler.go`
- [x] `backend/internal/handler/conta_fixa_handler.go`
- [x] `backend/internal/handler/despesa_geral_handler.go`
- [x] `backend/internal/handler/renda_fixa_handler.go`
- [x] `backend/internal/handler/renda_variavel_handler.go`
- [x] `backend/internal/handler/renda_extra_handler.go`
- [x] `backend/internal/handler/rendimento_investimento_handler.go`
- [x] `backend/internal/handler/dashboard_handler.go`
- [x] `backend/internal/handler/renda_historico_handler.go`
- [x] `backend/internal/handler/auth_handler.go` — `ValidateAccessToken` retorna 3 valores

### Handler testes (CONCLUÍDO)
- [x] `backend/internal/handler/membro_handler_test.go`
- [x] `backend/internal/handler/categoria_handler_test.go`
- [x] `backend/internal/handler/cartao_credito_handler_test.go`
- [x] `backend/internal/handler/despesa_cartao_handler_test.go`
- [x] `backend/internal/handler/assinatura_handler_test.go`
- [x] `backend/internal/handler/conta_fixa_handler_test.go`
- [x] `backend/internal/handler/despesa_geral_handler_test.go`
- [x] `backend/internal/handler/renda_fixa_handler_test.go`
- [x] `backend/internal/handler/renda_variavel_handler_test.go`
- [x] `backend/internal/handler/renda_extra_handler_test.go`
- [x] `backend/internal/handler/rendimento_investimento_handler_test.go`
- [x] `backend/internal/handler/dashboard_handler_test.go`
- [x] `backend/internal/handler/renda_historico_handler_test.go`

### Outros (CONCLUÍDO)
- [x] `backend/internal/repository/seed.go` — `SeedCategorias(db, familiaID)`, `SeedUsuarios` cria família compartilhada entre user1 (admin) e user2 (membro)
- [x] `backend/cmd/server/main.go` — cria `familiaRepo`, passa para `NewAuthService`, remove `SeedCategorias` separado

---

## Status atual: ✅ CONCLUÍDO E TESTADO

Todas as camadas foram atualizadas com isolamento por `familiaID`.

**Verificação final (2026-03-12):**
- ✅ `go build ./...` — compilação OK (sem erros)
- ✅ `go test ./...` — todos os testes passando
  - `internal/handler`: ok (cached)
  - `internal/service`: ok (cached)

**Próximos passos:** commit final + PR para master

---

## dashboard_service.go — detalhe da mudança

As interfaces for-dashboard precisam de `familiaID` como primeiro param:

```go
type RendaFixaRepositoryForDashboard interface {
    ListarVigentesPorMes(familiaID string, mes, ano int) ([]*domain.RendaFixa, error)
}
type DespesaCartaoRepositoryForDashboard interface {
    ListarPorFaturaGlobal(familiaID string, mes, ano int) ([]*domain.DespesaCartao, error)
}
type RendaVariavelRepositoryForDashboard interface {
    ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendaVariavel, error)
}
type RendaExtraRepositoryForDashboard interface {
    ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendaExtra, error)
}
type RendimentoRepositoryForDashboard interface {
    ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendimentoInvestimento, error)
}
type AssinaturaRepositoryForDashboard interface {
    ListarAtivas(familiaID string) ([]*domain.Assinatura, error)
}
type ContaFixaRepositoryForDashboard interface {
    ListarAtivas(familiaID string) ([]*domain.ContaFixa, error)
}
type DespesaGeralRepositoryForDashboard interface {
    ListarPorMes(familiaID string, mes, ano int) ([]*domain.DespesaGeral, error)
}
type DashboardRepositoryForCategorias interface {
    DespesasPorCategoria(familiaID string, mes, ano int) ([]domain.CategoriaTotalRaw, error)
}
```

Métodos públicos do DashboardService também ganham `familiaID string` como primeiro param:
- `ResumoMensal(familiaID string, mes, ano int)`
- `EvolucaoMensal(familiaID string, qtdMeses int)`
- `ProjecaoProximosMeses(familiaID string, qtdMeses int)`
- `DespesasPorCategoria(familiaID string, mes, ano int)`

Internamente, cada chamada ao repo recebe `familiaID`.

---

## renda_historico_service.go — detalhe da mudança

Interfaces for-historico ganham `familiaID`:

```go
type RendaFixaRepositoryForHistorico interface {
    Listar(familiaID string) ([]*domain.RendaFixa, error)
    ListarVigentesPorMes(familiaID string, mes, ano int) ([]*domain.RendaFixa, error)
}
type RendaVariavelRepositoryForHistorico interface {
    Listar(familiaID string) ([]*domain.RendaVariavel, error)
    ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendaVariavel, error)
}
type RendaExtraRepositoryForHistorico interface {
    Listar(familiaID string) ([]*domain.RendaExtra, error)
    ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendaExtra, error)
}
type RendimentoRepositoryForHistorico interface {
    Listar(familiaID string) ([]*domain.RendimentoInvestimento, error)
    ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendimentoInvestimento, error)
}
```

Método público:
- `BuscarHistorico(familiaID string, filtro domain.FiltroHistoricoRendas)`

Internamente, todas as chamadas ao repo recebem `familiaID`.

---

## Ordem recomendada para continuar

1. **rendimento_investimento_service_test.go** — atualizar chamadas de teste (mocks já OK)
2. **dashboard_service.go** — atualizar interfaces + implementação
3. **dashboard_service_test.go** — atualizar mocks + chamadas
4. **renda_historico_service.go** — atualizar interfaces + implementação
5. **renda_historico_service_test.go** — atualizar mocks + chamadas
6. **Todos os repository files** (13 arquivos) — adicionar `AND familia_id = $N` nas queries SQL
7. **Todos os handler files** (13 arquivos) — extrair `familiaID` via `getFamiliaID(c)` e passar para service
8. **Todos os handler test files** (13 arquivos) — atualizar mocks de serviço e setupRouter
9. **seed.go** — adaptar para família
10. **main.go** — wiring final
11. Verificação: `go build ./...` e `go test ./...`
