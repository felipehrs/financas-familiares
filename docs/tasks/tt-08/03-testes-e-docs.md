# Tarefa TT-08.3: Validação Final e Documentação

**Sprint:** 9
**Item:** TT-08
**Fase:** 3 de 3
**Estimativa:** 30 min
**Dependências:** TT-08.1 (CORS) + TT-08.2 (Rate Limiting)
**Nota:** As tarefas 1 e 2 já incluem testes automatizados (TDD). Esta fase foca em validação end-to-end e documentação.

---

## Objetivo

Validar a implementação completa do TT-08 (CORS + Rate Limiting) end-to-end, atualizar documentação do backend e marcar o item como concluído na Sprint 9.

**Por quê?** Garantir que os testes automatizados refletem o comportamento real do sistema em execução e que a documentação está completa para deploy em produção.

---

## Checklist de Implementação

### 1. Verificação Final de Build e Testes Automatizados

**Executar todos os testes (incluindo os TDD das tarefas anteriores):**

```bash
cd backend

# Compilação
go build ./...

# Testes automatizados (config + middleware)
go test ./...  -v
```

**Resultado esperado:**

```
=== RUN   TestLoad_CORSAllowedOrigins
--- PASS: TestLoad_CORSAllowedOrigins (0.XXs)
=== RUN   TestRateLimitMiddleware
--- PASS: TestRateLimitMiddleware (0.XXs)
... (outros testes) ...

PASS
ok      github.com/felipehrs/financas-familiares/backend/config              0.XXs
ok      github.com/felipehrs/financas-familiares/backend/internal/middleware 0.XXs
ok      github.com/felipehrs/financas-familiares/backend/internal/handler    (cached)
ok      github.com/felipehrs/financas-familiares/backend/internal/service    (cached)
```

- ✅ `go build ./...` — sem erros
- ✅ `go test ./...` — todos os testes passando (incluindo TDD do TT-08)

- [ ] Build executado com sucesso
- [ ] Testes automatizados executados com sucesso (4+ cenários CORS, 3+ cenários rate limit)

---

### 2. Teste de Integração (Manual)

**Pré-requisito:** Servidor rodando localmente

```bash
cd backend
go run cmd/server/main.go
```

#### Teste 2.1: CORS Permitido

```bash
curl -X OPTIONS http://localhost:8080/api/v1/membros \
  -H "Origin: http://localhost:5173" \
  -H "Access-Control-Request-Method: GET" \
  -v 2>&1 | grep -i "access-control"
```

**Resultado esperado:**

```
< Access-Control-Allow-Origin: http://localhost:5173
< Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
```

- [ ] CORS permitido para origem autorizada

---

#### Teste 2.2: CORS Bloqueado

```bash
curl -X OPTIONS http://localhost:8080/api/v1/membros \
  -H "Origin: https://malicious-site.com" \
  -H "Access-Control-Request-Method: GET" \
  -v 2>&1 | grep -i "access-control"
```

**Resultado esperado:**

```
(nenhum output — header CORS ausente)
```

- [ ] CORS bloqueado para origem não autorizada

---

#### Teste 2.3: Rate Limiting (6 tentativas de login)

```bash
# Script de teste rápido
for i in {1..6}; do
  echo "Request $i:"
  curl -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@test.com","password":"wrong"}' \
    -w " | Status: %{http_code}\n" \
    -s
  sleep 0.5
done
```

**Resultado esperado:**

```
Request 1: ... | Status: 401
Request 2: ... | Status: 401
Request 3: ... | Status: 401
Request 4: ... | Status: 401
Request 5: ... | Status: 401
Request 6: ... | Status: 429  ← Bloqueado!
```

- [ ] Primeiras 5 tentativas permitidas (HTTP 401)
- [ ] 6ª tentativa bloqueada (HTTP 429)

---

#### Teste 2.4: Headers de Rate Limiting

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"wrong"}' \
  -v 2>&1 | grep -i "x-ratelimit"
```

**Resultado esperado:**

```
< X-RateLimit-Limit: 5
< X-RateLimit-Remaining: 4  (ou 3, 2, 1, 0 dependendo de quantas requests já foram feitas)
< X-RateLimit-Reset: 1710259200  (timestamp Unix)
```

- [ ] Headers `X-RateLimit-*` presentes

---

### 3. Atualizar `backend/README.md`

**Localização:** `backend/README.md`

Adicionar uma nova seção após a seção de "Configuração" (ou criar se não existir):

```markdown
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
```

- [ ] Seção "Segurança" adicionada ao README
- [ ] CORS documentado
- [ ] Rate Limiting documentado

---

### 4. Atualizar `docs/produto/sprints.md`

**Arquivo:** `docs/produto/sprints.md`

**Encontrar a seção "Progresso TT-08" e atualizar todos os checkboxes:**

```markdown
### Progresso TT-08 (concluído em XX/03/2026)

> **Contexto:** Implementação de segurança adicional: restrição de CORS e rate limiting no login.
> **Plano detalhado:** [docs/tecnico/05-plano-tt08-cors-rate-limiting.md](../tecnico/05-plano-tt08-cors-rate-limiting.md)

| Fase | Status | Detalhe |
|------|--------|---------|
| **Fase 1: CORS Restrito** | ✅ | |
| Adicionar env vars `CORS_ALLOWED_ORIGINS` | ✅ | Em `.env` e `.env.example` |
| Atualizar `config/config.go` | ✅ | Campo `CORS.AllowedOrigins` + parsing |
| Substituir `AllowAllOrigins` em `main.go` | ✅ | Usar whitelist de origens |
| Testar CORS manualmente | ✅ | curl com origem permitida + bloqueada |
| **Fase 2: Rate Limiting** | ✅ | |
| Instalar `github.com/ulule/limiter/v3` | ✅ | `go get` |
| Criar `middleware/rate_limit.go` | ✅ | Middleware genérico |
| Adicionar env var `RATE_LIMIT_LOGIN` | ✅ | Default: "5-M" (5 por minuto) |
| Aplicar middleware no `/auth/login` | ✅ | Em `main.go` |
| Testar rate limiting manualmente | ✅ | Script bash com 6 requests |
| Validar HTTP 429 e headers | ✅ | `X-RateLimit-*` presentes |
| **Fase 3: Docs e Finalização** | ✅ | |
| Atualizar `backend/README.md` | ✅ | Seção "Segurança" |
| Verificação final | ✅ | `go build ./...` e `go test ./...` |
```

**Atualizar também a tabela de itens da Sprint 9:**

```markdown
| Seq | Status | ID | Descrição |
|-----|--------|-----|-----------|
| 1 | ✅ | TT-07 | Infraestrutura Multi-Tenant: tabelas `familias`, colunas `familia_id`, alteração JWT e Middleware |
| 2 | ✅ | TT-10 | Refatoração para Isolamento: filtrar 15+ repositórios e handlers por `familia_id` |
| 3 | ✅ | TT-08 | Segurança: restringir CORS a origens específicas; rate limiting em `POST /auth/login` |
| 4 | 🔲 | TT-09 | Índices de performance: `deleted_at` em todas as tabelas, `familia_id` em todas as tabelas, índice composto em `despesas_cartao` |
```

- [ ] Progresso TT-08 atualizado com todos os ✅
- [ ] Status do TT-08 alterado de 🔄 para ✅
- [ ] Data de conclusão adicionada

---

### 5. Atualizar `.env.example` (Verificação Final)

Garantir que o arquivo `.env.example` está completo:

```env
# Database
DATABASE_URL=postgresql://user:password@localhost:5432/financas_db?sslmode=disable

# JWT
JWT_SECRET=your-super-secret-key-change-in-production

# Seed (opcional, apenas para dev)
SEED_USER1_EMAIL=user1@test.com
SEED_USER1_PASSWORD=senha123
SEED_USER2_EMAIL=user2@test.com
SEED_USER2_PASSWORD=senha456

# CORS: lista de origens permitidas, separadas por vírgula
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

# Rate Limiting: formato "X-Y" (X requests por Y período: S, M, H, D)
# Exemplo: "5-M" = 5 tentativas por minuto
RATE_LIMIT_LOGIN=5-M

# Server
PORT=8080
```

- [ ] `.env.example` contém todas as env vars necessárias
- [ ] Comentários explicativos presentes

---

## Critério de Conclusão

**Testes Automatizados (TDD - concluídos nas tarefas anteriores):**
- [x] `TestLoad_CORSAllowedOrigins`: 4 cenários passando
- [x] `TestRateLimitMiddleware`: 3 cenários passando

**Validação End-to-End:**
- [x] `go build ./...` executado com sucesso
- [x] `go test ./...` executado com sucesso (todos os testes)
- [x] Testes manuais de CORS validados (permitido + bloqueado)
- [x] Testes manuais de rate limiting validados (5 OK + 6ª bloqueada)
- [x] Headers `X-RateLimit-*` confirmados

**Documentação:**
- [x] `backend/README.md` atualizado com seção "Segurança"
- [x] `docs/produto/sprints.md` marcado TT-08 como ✅
- [x] `.env.example` completo e documentado

---

## Commit Sugerido

```bash
git add backend/README.md \
        backend/.env.example \
        docs/produto/sprints.md

git commit -m "docs: atualizar README e sprints com conclusão TT-08 ✅

Finaliza implementação de segurança adicional (CORS + rate limiting).

**Documentação:**
- Adicionar seção 'Segurança' no backend/README.md
- Explicar configuração de CORS_ALLOWED_ORIGINS
- Explicar configuração de RATE_LIMIT_LOGIN
- Descrever comportamento e headers de rate limiting
- Nota sobre Redis para escala horizontal

**Sprint 9:**
- Marcar TT-08 como ✅ concluído
- Atualizar checklist de progresso (todas as fases ✅)
- Adicionar data de conclusão

**Validação:**
- go build ./... ✅
- go test ./... ✅
- CORS restrito: testado ✅
- Rate limiting: testado ✅

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Resultado Final do TT-08

Ao concluir esta tarefa, o TT-08 está completo com:

✅ **CORS Restrito**
- Whitelist de origens via `CORS_ALLOWED_ORIGINS`
- Bloqueio de sites não autorizados

✅ **Rate Limiting no Login**
- Limite de 5 tentativas por minuto por IP
- HTTP 429 após exceder limite
- Headers `X-RateLimit-*` informativos

✅ **Documentação**
- README atualizado com seção de segurança
- Sprints marcado como concluído
- `.env.example` completo

---

## Próximos Passos (Sprint 9)

Após conclusão do TT-08, prosseguir para:
- **TT-09: Índices de Performance** — Criar índices em `deleted_at`, `familia_id` e índice composto em `despesas_cartao`

---

## Referências

- [Plano técnico TT-08](../../tecnico/05-plano-tt08-cors-rate-limiting.md)
- [Tarefa TT-08.1](./01-cors-restrito.md)
- [Tarefa TT-08.2](./02-rate-limiting.md)
