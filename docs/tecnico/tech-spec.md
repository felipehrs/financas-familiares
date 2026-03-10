# Especificação Técnica - Sistema de Gestão de Finanças Familiares

**Versão:** 1.0
**Data:** 06 de Março de 2026
**Referência:** spec.md

---

## 1. Visão Geral da Arquitetura

O sistema adota arquitetura **cliente-servidor**, com frontend PWA (Progressive Web App) e backend REST API. O armazenamento segue o modelo **offline-first**: os dados são persistidos localmente no dispositivo e sincronizados com a nuvem quando há conexão disponível.

```
┌─────────────────────────────────────────────────────┐
│                    CLIENTE (PWA)                    │
│  React + TypeScript  │  SQLite (via WASM/IndexedDB) │
└──────────────────────┬──────────────────────────────┘
                       │ HTTPS / REST
┌──────────────────────▼──────────────────────────────┐
│                   BACKEND (API)                     │
│              Go + Gin Framework                     │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│                 BANCO DE DADOS                      │
│                  PostgreSQL                         │
└─────────────────────────────────────────────────────┘
```

---

## 2. Frontend

### Linguagem e Framework
- **React 18** com **TypeScript**
  - Ecossistema maduro, excelente suporte a PWA
  - TypeScript garante segurança de tipos alinhada à experiência com Java/Go
  - Componentização facilita o isolamento de regras de negócio do domínio (ex: cálculo de fatura, parcelas)

### Principais Bibliotecas
| Propósito | Biblioteca |
|---|---|
| Roteamento | React Router v7 |
| Estado global | Zustand (leve) ou TanStack Query (para cache/sync) |
| UI Components | shadcn/ui (componentes acessíveis, customizáveis) |
| Formulários | React Hook Form + Zod (validação com schema) |
| Gráficos | Recharts |
| Formatação de datas | date-fns |
| Formatação monetária | Intl.NumberFormat (nativo) |
| Testes | Vitest + React Testing Library |

### PWA e Offline
- **Vite** como bundler (rápido, suporte nativo a PWA via `vite-plugin-pwa`)
- **Service Worker** para cache de assets e modo offline
- **IndexedDB** (via `Dexie.js`) para persistência local no browser
  - Alternativa: SQLite compilado para WASM (`wa-sqlite` ou `sql.js`) para maior fidelidade com o banco do servidor

### Estilo
- **Tailwind CSS** — utilitário, excelente para mobile-first
- Design responsivo: breakpoints de 320px (mobile) a 1920px (desktop)

---

## 3. Backend

### Linguagem e Framework
- **Go (Golang)** com framework **Gin**
  - Alinhado à experiência do desenvolvedor
  - Excelente performance para APIs REST
  - Tipagem estática facilita representação do domínio financeiro
  - Binário único de fácil deploy (sem JVM, sem runtime externo)
  - Alternativa caso prefira ecossistema maior: **Java 21 + Spring Boot 3**

### Estrutura do Projeto (Go)
```
backend/
├── cmd/
│   └── server/          # Entrypoint (main.go)
├── internal/
│   ├── domain/          # Entidades e regras de negócio puras
│   │   ├── membro.go
│   │   ├── cartao.go
│   │   ├── despesa.go
│   │   └── ...
│   ├── service/         # Casos de uso (cálculo de fatura, parcelas, saldo)
│   ├── repository/      # Interface de acesso a dados (abstrações)
│   ├── handler/         # Handlers HTTP (Gin)
│   └── middleware/      # Auth, CORS, logging
├── migrations/          # Migrações SQL versionadas
├── config/              # Configuração (env vars)
└── Makefile
```

### Principais Dependências (Go)
| Propósito | Pacote |
|---|---|
| HTTP Framework | `github.com/gin-gonic/gin` |
| ORM / Query Builder | `github.com/jackc/pgx/v5` + `github.com/jmoiron/sqlx` |
| Migrações | `github.com/golang-migrate/migrate/v4` |
| Autenticação JWT | `github.com/golang-jwt/jwt/v5` |
| Validação | `github.com/go-playground/validator/v10` |
| Configuração | `github.com/spf13/viper` |
| Testes | `testing` (stdlib) + `github.com/stretchr/testify` |

### API
- REST com JSON
- Versionamento: `/api/v1/...`
- Autenticação via **JWT Bearer Token** (header `Authorization`)
- Tokens de acesso com expiração curta (ex: 15min) + refresh token (ex: 7 dias)

---

## 4. Banco de Dados

### Produção (Cloud)
- **PostgreSQL 16**
  - Robusto, suporte a tipos avançados, transações ACID
  - Excelente para dados financeiros
  - Compatível com a maioria dos provedores cloud (Railway, Render, Supabase, AWS RDS)

### Local (Desenvolvimento)
- PostgreSQL via Docker Compose

### Migrações
- Arquivos `.sql` versionados em `migrations/`
- Ferramenta: `golang-migrate`
- Convenção: `000001_create_membros.up.sql` / `000001_create_membros.down.sql`

---

## 5. Autenticação e Segurança

- Dois usuários fixos por instância (casal), cadastrados na configuração inicial
- Senha armazenada com **bcrypt** (fator de custo ≥ 12)
- JWT para sessões stateless
- Comunicação exclusivamente via **HTTPS** (TLS obrigatório em produção)
- Dados em repouso: criptografia gerenciada pelo provedor cloud (PostgreSQL)
- CORS configurado para aceitar apenas a origem do frontend

---

## 6. Sincronização Offline

### Estratégia
- Frontend persiste operações localmente em IndexedDB com status `pendente`
- Quando a conexão é restaurada, uma fila de sincronização envia as operações ao backend
- Política de conflito: **última escrita vence** (`updated_at` como desempate)
- Todos os registros possuem campo `updated_at` (timestamp com timezone)

### Campos obrigatórios em todas as entidades para sync
```
id          UUID (gerado no cliente para evitar conflito de IDs)
created_at  TIMESTAMPTZ
updated_at  TIMESTAMPTZ
deleted_at  TIMESTAMPTZ (soft delete, para sync de exclusões)
```

---

## 7. Infraestrutura e Deploy

### Opções recomendadas (baixo custo / gratuitas para MVP)
| Componente | Opção Recomendada | Alternativa |
|---|---|---|
| Backend (Go) | **Railway** | Render, Fly.io |
| Banco de Dados | **Railway PostgreSQL** | Supabase, Render PostgreSQL |
| Frontend (PWA) | **Vercel** ou **Netlify** | GitHub Pages |

### Docker
- `Dockerfile` multi-stage para build do binário Go
- `docker-compose.yml` para desenvolvimento local (backend + PostgreSQL)

---

## 8. Tooling de Desenvolvimento

| Propósito | Ferramenta |
|---|---|
| Gerenciador de pacotes JS | pnpm |
| Linter JS/TS | ESLint + Prettier |
| Linter Go | `golangci-lint` |
| Controle de versão | Git |
| CI/CD | GitHub Actions |
| Documentação da API | Swagger / OpenAPI 3.0 (via `swaggo/swag`) |

---

## 11. Estratégia de Testes (TDD)

O desenvolvimento adota **Test-Driven Development (TDD)** como prática central. Todo código de produção deve ser precedido por testes que falham, seguindo o ciclo clássico: **Red → Green → Refactor**.

### Princípios

- Escrever o teste antes da implementação (Red)
- Implementar o mínimo necessário para o teste passar (Green)
- Refatorar mantendo os testes verdes (Refactor)
- Cobertura mínima de 80% nas camadas de domínio e serviço
- Testes de integração para endpoints HTTP críticos

### Backend (Go)

| Camada | Tipo de Teste | Ferramenta |
|---|---|---|
| Domain / Regras de negócio | Unitário | `testing` stdlib + `testify` |
| Service (casos de uso) | Unitário com mocks | `testify/mock` |
| Repository | Integração (banco real) | `testcontainers-go` + PostgreSQL |
| Handler HTTP | Integração | `net/http/httptest` + `testify` |

**Foco prioritário para TDD:**
- `CalcularFatura(dataCompra, diaFechamento)` → RN01
- `DistribuirParcelas(compra, faturaInicial)` → RN03
- `CalcularSaldoMensal(rendas, despesas)` → RN06
- `GerarProjecao(meses, rendasFixas, despesasRecorrentes)` → RN07

**Convenção de arquivos:**
```
internal/domain/cartao_test.go   # testa cartao.go
internal/service/saldo_test.go   # testa saldo.go
```

### Frontend (React + TypeScript)

| Camada | Tipo de Teste | Ferramenta |
|---|---|---|
| Funções de domínio puras | Unitário | Vitest |
| Componentes React | Componente | Vitest + React Testing Library |
| Fluxos de usuário | E2E (Fase 2+) | Playwright |

**Foco prioritário para TDD:**
- Funções de cálculo de fatura e parcelas (duplicadas no frontend para feedback imediato)
- Validações de formulário (Zod schemas)
- Lógica de formatação monetária e de datas

### Cobertura Mínima por Camada

| Camada | Cobertura Mínima |
|---|---|
| Domain (Go) | 90% |
| Service (Go) | 80% |
| Handler HTTP | 70% |
| Funções puras (TS) | 90% |
| Componentes React | 60% |

### Execução dos Testes

```bash
# Backend
make test              # todos os testes
make test-unit         # apenas unitários
make test-integration  # apenas integração (requer Docker)
make test-coverage     # relatório de cobertura

# Frontend
pnpm test              # todos os testes (watch mode)
pnpm test:run          # todos os testes (CI mode)
pnpm test:coverage     # relatório de cobertura
```

---

## 9. Considerações sobre a Escolha do Stack

### Por que Go no backend?
- Experiência prévia do desenvolvedor
- Performance nativa superior ao Node.js para cálculos intensivos (projeções, consolidações)
- Deploy simples (binário único, imagem Docker pequena ~10MB)
- Tipagem estática alinha bem com domínio financeiro (sem surpresas de tipo em cálculos)

### Por que React + TypeScript no frontend?
- Ecossistema mais maduro para PWA comparado a Vue
- TypeScript: familiaridade com tipagem estática (Java/Go background)
- Ampla disponibilidade de componentes de UI de qualidade (shadcn/ui, Radix UI)

### Alternativas descartadas
- **Flutter/React Native:** desnecessário para MVP; PWA atende mobile com menor complexidade
- **Java + Spring Boot:** viável, mas Go oferece deploy mais simples e performance superior para este porte
- **Next.js / SSR:** desnecessário; app é totalmente client-side (dados privados, sem SEO)

---

## 10. Próximos Passos

1. Validar e aprovar esta especificação técnica
2. Definir estrutura de repositório (monorepo ou repos separados)
3. Criar estrutura inicial do projeto (scaffolding)
4. Implementar MVP conforme Fase 1 do roadmap (spec.md §10)
