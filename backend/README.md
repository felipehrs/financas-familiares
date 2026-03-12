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
