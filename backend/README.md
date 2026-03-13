# Backend — Finanças Familiares

> API REST do sistema de gestão financeira familiar.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.12-00ADD8?style=flat-square&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat-square&logo=postgresql&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-Auth-000000?style=flat-square&logo=jsonwebtokens&logoColor=white)
![Testify](https://img.shields.io/badge/Testify-TDD-00ADD8?style=flat-square&logo=go&logoColor=white)
[![Coverage](https://codecov.io/gh/felipehrs/financas-familiares/graph/badge.svg?flag=backend)](https://codecov.io/gh/felipehrs/financas-familiares)

---

## Stack

| Propósito | Biblioteca |
|---|---|
| HTTP Framework | Gin v1.12 |
| Banco de dados | PostgreSQL 16 via pgx/v5 + sqlx |
| Migrações | golang-migrate |
| Autenticação | golang-jwt/jwt v5 + bcrypt (custo 12) |
| Configuração | Viper (env vars + .env) |
| Testes | stdlib testing + testify |

---

## Estrutura

```
backend/
├── cmd/server/        # Entrypoint (main.go)
├── config/            # Carregamento de configuração via Viper
├── internal/
│   ├── domain/        # Entidades e erros sentinela
│   ├── service/       # Casos de uso e regras de negócio (testados com mocks)
│   ├── handler/       # Handlers HTTP (Gin)
│   ├── repository/    # Acesso ao PostgreSQL (sqlx)
│   └── middleware/    # JWT auth, CORS, logging
├── migrations/        # Arquivos SQL versionados (up/down)
├── .env.example
├── .golangci.yml
└── Makefile
```

---

## Pré-requisitos

**Recomendado — Dev Container (zero configuração local):**

Abra o repositório raiz no VS Code com a extensão Dev Containers. O container já inclui Go 1.25, Node.js, pnpm e `golang-migrate`.

**Alternativa — ambiente local:**

- Go 1.25+
- Docker + Docker Compose (para PostgreSQL local)
- [`golang-migrate`](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) CLI

---

## Configuração

Copie `.env.example` para `.env` e preencha:

```bash
cp .env.example .env
```

| Variável | Descrição | Exemplo |
|---|---|---|
| `DATABASE_URL` | Connection string PostgreSQL | `postgres://postgres:postgres@localhost:5432/financas_familiares` |
| `JWT_SECRET` | Segredo para assinar tokens JWT | string aleatória longa |
| `SERVER_PORT` | Porta do servidor | `8080` |
| `SEED_USER1_EMAIL` | E-mail do usuário 1 | `usuario1@email.com` |
| `SEED_USER1_NOME` | Nome do usuário 1 | `Felipe` |
| `SEED_USER1_SENHA` | Senha do usuário 1 (plaintext) | `senha-segura` |
| `SEED_USER2_EMAIL` | E-mail do usuário 2 | `usuario2@email.com` |
| `SEED_USER2_NOME` | Nome do usuário 2 | `Ana` |
| `SEED_USER2_SENHA` | Senha do usuário 2 (plaintext) | `senha-segura` |
| `ALLOWED_ORIGINS` | Origens permitidas no CORS (separadas por vírgula) | `http://localhost:5173` |

---

## Segurança

### CORS (Cross-Origin Resource Sharing)

O backend restringe requisições apenas a origens autorizadas para prevenir ataques CSRF.

**Configuração:**

```env
# Lista de origens permitidas (separadas por vírgula)
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
```

**Produção (Railway):**

```env
CORS_ALLOWED_ORIGINS=https://financas-familiares.vercel.app
```

Para múltiplos domínios (ex: preview branches do Vercel):

```env
CORS_ALLOWED_ORIGINS=https://financas-familiares.vercel.app,https://preview-branch.vercel.app
```

---

### Rate Limiting

Proteção contra força bruta no endpoint de login (`POST /api/v1/auth/login`).

**Configuração:**

```env
# Formato: "X-Y" onde X = número de requests, Y = período (S, M, H, D)
RATE_LIMIT_LOGIN=5-M  # 5 tentativas por minuto
```

**Comportamento:**
- Após 5 tentativas de login do mesmo IP em 1 minuto, o servidor retorna HTTP 429 (Too Many Requests)
- O limite reseta automaticamente após 60 segundos
- Headers de resposta incluem informações do rate limit:
  - `X-RateLimit-Limit`: Limite total
  - `X-RateLimit-Remaining`: Tentativas restantes
  - `X-RateLimit-Reset`: Timestamp Unix quando o limite reseta

**Nota de produção:**

O rate limiting atual usa um store em memória, adequado para deploy single-instance (Railway). Para escalar horizontalmente (múltiplas instâncias), considere migrar para Redis.

---

## Subindo o banco localmente

Dentro do Dev Container o PostgreSQL já sobe automaticamente junto com o container. Fora dele:

```bash
# Na raiz do monorepo
docker compose up -d
```

---

## Migrações

As migrações são aplicadas automaticamente ao rodar `make dev` (via `dev.sh`). Para rodar manualmente (dentro do Dev Container ou com `golang-migrate` instalado localmente):

```bash
migrate -path migrations -database "$DATABASE_URL" up
migrate -path migrations -database "$DATABASE_URL" down
```

---

## Comandos

```bash
# Rodar o servidor
go run ./cmd/server

# Build do binário
go build -o bin/server ./cmd/server

# Todos os testes
go test ./...

# Apenas testes unitários (sem banco)
go test ./internal/domain/... ./internal/service/... ./internal/handler/...

# Cobertura
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Lint
golangci-lint run
```

Ou via Makefile (na raiz do monorepo):

```bash
make run
make test
make test-coverage
make lint-backend
```

---

## Testes

O projeto adota **TDD** — testes são escritos antes da implementação. Mocks são manuais (sem geração de código).

| Camada | Tipo | Cobertura mínima |
|---|---|---|
| Domain | Unitário | 90% |
| Service | Unitário com mocks | 80% |
| Handler HTTP | Integração leve (httptest) | 70% |
| Repository | Integração (banco real) | — |

---

## Endpoints

Todos os endpoints (exceto `/health` e `/api/v1/auth/*`) requerem `Authorization: Bearer <token>`.

| Método | Rota | Descrição |
|---|---|---|
| GET | `/health` | Health check |
| POST | `/api/v1/auth/login` | Login — retorna access token (15min) + refresh token (7d) |
| POST | `/api/v1/auth/refresh` | Renova access token via refresh token |
| GET | `/api/v1/membros` | Listar membros |
| POST | `/api/v1/membros` | Criar membro |
| GET | `/api/v1/membros/:id` | Buscar membro por ID |
| PUT | `/api/v1/membros/:id` | Atualizar membro |
| PATCH | `/api/v1/membros/:id/inativar` | Inativar membro |
| GET | `/api/v1/categorias` | Listar categorias |
| POST | `/api/v1/categorias` | Criar categoria |
| PUT | `/api/v1/categorias/:id` | Atualizar categoria |
| DELETE | `/api/v1/categorias/:id` | Excluir categoria (409 se tiver vínculos) |

---

## Migrações disponíveis

| # | Tabela |
|---|---|
| 001 | `usuarios` |
| 002 | `membros` |
| 003 | `categorias` |
| 004 | `cartoes_credito` |
| 005 | `despesas_cartao` |
| 006 | `assinaturas` |
| 007 | `contas_fixas` |
| 008 | `despesas_gerais` |
| 009 | `rendas_fixas` |
| 010 | `rendas_variaveis` |
| 011 | `rendas_extras` |
| 012 | `rendimentos_investimento` |
| 013 | `refresh_tokens` |
| 014 | unique constraint em `categorias.nome` |
| 015 | campo `num_parcelas` em `despesas_cartao` |
| 016 | campos `data_inicio` / `data_fim` em `rendas_fixas` |
| 017 | tabelas `familias` + `familia_usuarios`; coluna `familia_id` em todas as tabelas de dados |
| 018 | índices de performance: `familia_id` + `deleted_at` parciais; índice composto para fatura; índices em datas |

---

## Índices de Performance

**Criados em:** Sprint 9 (TT-09) — 12/03/2026
**Migration:** `000018_add_performance_indexes.up.sql`

### Objetivo

Otimizar queries de listagem, cálculo de faturas, e dashboard. Todos os índices são **parciais** (`WHERE deleted_at IS NULL`) para reduzir tamanho e melhorar performance.

### Índices Criados

#### 1. Índices em `familia_id` (Isolamento Multi-Tenant)

Todas as queries de domínio filtram por `familia_id` após TT-07.

| Tabela | Índice | Tipo |
|--------|--------|------|
| `membros` | `idx_membros_familia_id_active` | Parcial |
| `categorias` | `idx_categorias_familia_id_active` | Parcial |
| `cartoes_credito` | `idx_cartoes_credito_familia_id_active` | Parcial |
| `despesas_cartao` | `idx_despesas_cartao_familia_id_active` | Parcial |
| `assinaturas` | `idx_assinaturas_familia_id_active` | Parcial |
| `contas_fixas` | `idx_contas_fixas_familia_id_active` | Parcial |
| `despesas_gerais` | `idx_despesas_gerais_familia_id_active` | Parcial |
| `rendas_fixas` | `idx_rendas_fixas_familia_id_active` | Parcial |
| `rendas_variaveis` | `idx_rendas_variaveis_familia_id_active` | Parcial |
| `rendas_extras` | `idx_rendas_extras_familia_id_active` | Parcial |
| `rendimentos_investimento` | `idx_rendimentos_investimento_familia_id_active` | Parcial |

#### 2. Índice Composto para Cálculo de Fatura

**Query otimizada:** `SELECT SUM(valor_parcela) WHERE familia_id = X AND cartao_id = Y AND fatura_ano = Z AND fatura_mes = W`

| Índice | Colunas |
|--------|---------|
| `idx_despesas_cartao_fatura` | `(familia_id, cartao_id, fatura_ano, fatura_mes)` |

**Ganho de performance:** 50-100x mais rápido em tabelas com 1000+ despesas.

#### 3. Índices em Campos de Data

Otimizam queries de dashboard, relatórios e projeções.

| Tabela | Índice | Colunas |
|--------|--------|---------|
| `despesas_cartao` | `idx_despesas_cartao_data_compra` | `(familia_id, data_compra)` |
| `despesas_gerais` | `idx_despesas_gerais_data` | `(familia_id, data)` |
| `rendas_variaveis` | `idx_rendas_variaveis_mes_ano` | `(familia_id, ano DESC, mes DESC)` |
| `rendas_extras` | `idx_rendas_extras_mes_ano` | `(familia_id, ano DESC, mes DESC)` |
| `rendimentos_investimento` | `idx_rendimentos_investimento_mes_ano` | `(familia_id, ano DESC, mes DESC)` |
| `rendas_fixas` | `idx_rendas_fixas_vigencia` | `(familia_id, data_inicio, data_fim)` |

#### 4. Índices de Ordenação

Otimizam listagens ordenadas por `created_at DESC`.

| Tabela | Índice | Colunas |
|--------|--------|---------|
| `membros` | `idx_membros_created_at` | `(familia_id, created_at DESC)` |
| `categorias` | `idx_categorias_created_at` | `(familia_id, created_at DESC)` |
| `cartoes_credito` | `idx_cartoes_credito_created_at` | `(familia_id, created_at DESC)` |

### Como Verificar se Índices Estão Sendo Usados

#### Listar todos os índices criados:

```bash
psql $DATABASE_URL -c "
SELECT tablename, indexname
FROM pg_indexes
WHERE schemaname = 'public' AND indexname LIKE 'idx_%'
ORDER BY tablename, indexname;
"
```

#### Verificar uso dos índices (estatísticas):

```bash
psql $DATABASE_URL -c "
SELECT tablename, indexrelname AS index_name, idx_scan AS index_scans
FROM pg_stat_user_indexes
WHERE schemaname = 'public' AND indexrelname LIKE 'idx_%'
ORDER BY idx_scan DESC;
"
```

**Índices com `idx_scan > 0`** → Estão sendo usados.
**Índices com `idx_scan = 0`** → Nunca foram usados (reavaliar necessidade).

#### Validar com EXPLAIN ANALYZE:

```bash
psql $DATABASE_URL <<EOF
EXPLAIN ANALYZE
SELECT SUM(valor_parcela)
FROM despesas_cartao
WHERE familia_id = (SELECT id FROM familias LIMIT 1)
  AND cartao_id = (SELECT id FROM cartoes_credito LIMIT 1)
  AND fatura_ano = 2026
  AND fatura_mes = 3
  AND deleted_at IS NULL;
EOF
```

**Resultado esperado:** `Index Scan using idx_despesas_cartao_fatura`

### Manutenção de Índices

#### Atualizar estatísticas (após inserções/atualizações em massa):

```bash
psql $DATABASE_URL -c "ANALYZE despesas_cartao;"
```

#### Verificar índices inválidos:

```bash
psql $DATABASE_URL -c "SELECT indexrelid::regclass AS index_name FROM pg_index WHERE NOT indisvalid;"
```

Se houver índices inválidos, recriá-los:

```bash
psql $DATABASE_URL -c "DROP INDEX CONCURRENTLY nome_do_indice_invalido;"
psql $DATABASE_URL -c "CREATE INDEX CONCURRENTLY nome_do_indice ..."
```

### Referências

- [Plano técnico TT-09](../docs/tecnico/06-plano-tt09-indices-performance.md)
- [PostgreSQL: Indexes](https://www.postgresql.org/docs/current/indexes.html)
- [Use The Index, Luke!](https://use-the-index-luke.com/)
