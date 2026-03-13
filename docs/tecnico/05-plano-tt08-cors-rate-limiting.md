# Plano de Execução — TT-08: CORS Restrito e Rate Limiting

> **Sprint 9** | **Criado em:** 12/03/2026
> **Pré-requisito:** TT-07 concluído (infraestrutura de isolamento por familia_id)

---

## Índice

1. [Contexto e Motivação](#1-contexto-e-motivação)
2. [Objetivos do TT-08](#2-objetivos-do-tt-08)
3. [Estratégia de Implementação](#3-estratégia-de-implementação)
4. [Parte 1: CORS Restrito](#4-parte-1-cors-restrito)
5. [Parte 2: Rate Limiting no Login](#5-parte-2-rate-limiting-no-login)
6. [Testes e Validação](#6-testes-e-validação)
7. [Critério de Conclusão](#7-critério-de-conclusão)
8. [Checklist de Implementação](#8-checklist-de-implementação)

---

## 1. Contexto e Motivação

### Problema Atual

Conforme identificado na avaliação pós-Sprint 8 ([docs/tecnico/02-avaliacao-e-multi-tenant.md](./02-avaliacao-e-multi-tenant.md)):

**1. CORS Muito Aberto (Alto Risco)**

Em `backend/cmd/server/main.go:113-119`:

```go
r.Use(cors.New(cors.Config{
    AllowAllOrigins:  true,  // ⚠️ Aceita requisições de qualquer origem
    AllowCredentials: false,
}))
```

**Impacto:**
- Qualquer site malicioso pode fazer requisições HTTP para a API
- Mesmo com `AllowCredentials: false`, JWTs no header `Authorization` ainda estão expostos
- Um atacante pode criar um site que faz chamadas à API em nome do usuário autenticado

**2. Falta de Rate Limiting**

Não há middleware de rate limiting em nenhuma rota.

**Impacto:**
- Atacante pode fazer força bruta em `POST /api/v1/auth/login` sem limitação de tentativas
- Possível DoS (Denial of Service) em endpoints públicos
- Nenhuma proteção contra automação maliciosa

---

## 2. Objetivos do TT-08

1. **Restringir CORS** a origens específicas (whitelist de domínios permitidos)
2. **Implementar rate limiting** no endpoint de login com limite de tentativas por IP
3. **Retornar HTTP 429** (Too Many Requests) quando o limite for excedido
4. **Configurar via variáveis de ambiente** para flexibilidade entre dev/prod
5. **Adicionar testes** para validar o comportamento

---

## 3. Estratégia de Implementação

### Ordem de Execução

1. **CORS restrito** (independente, baixo risco)
2. **Rate limiting** (requer nova dependência, novo middleware)
3. **Testes de integração** (validação manual + automatizada)
4. **Documentação** (README atualizado com novas env vars)

### Dependências Novas

**Para rate limiting, vamos usar `github.com/ulule/limiter/v3`:**

```bash
cd backend
go get github.com/ulule/limiter/v3
go get github.com/ulule/limiter/v3/drivers/store/memory
```

**Justificativa:**
- Biblioteca madura e bem mantida
- Suporte a múltiplos stores (memory, redis, etc.)
- API simples e idempotente
- Usado por projetos Go de grande escala

---

## 4. Parte 1: CORS Restrito

### 4.1 Adicionar Variáveis de Ambiente

**Arquivo:** `backend/.env` e `backend/.env.example`

```env
# CORS: lista de origens permitidas, separadas por vírgula
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
```

Em produção (Railway):

```env
CORS_ALLOWED_ORIGINS=https://financas-familiares.vercel.app
```

### 4.2 Atualizar `config/config.go`

Adicionar campo `CORS` na struct `Config`:

```go
type Config struct {
    Database struct {
        URL string
    }
    JWT struct {
        Secret string
    }
    Seed struct {
        User1Email    string
        User1Password string
        User2Email    string
        User2Password string
    }
    CORS struct {
        AllowedOrigins []string  // Nova propriedade
    }
    Server struct {
        Port string
    }
}
```

**Leitura da env var:**

```go
func Load() (*Config, error) {
    // ... código existente ...

    // Parse CORS allowed origins
    originsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
    if originsStr == "" {
        originsStr = "http://localhost:5173"  // Default para dev
    }
    cfg.CORS.AllowedOrigins = strings.Split(originsStr, ",")
    for i := range cfg.CORS.AllowedOrigins {
        cfg.CORS.AllowedOrigins[i] = strings.TrimSpace(cfg.CORS.AllowedOrigins[i])
    }

    return &cfg, nil
}
```

### 4.3 Atualizar `main.go`

**Antes:**

```go
r.Use(cors.New(cors.Config{
    AllowAllOrigins:  true,
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: false,
}))
```

**Depois:**

```go
r.Use(cors.New(cors.Config{
    AllowOrigins:     cfg.CORS.AllowedOrigins,  // ✅ Whitelist específica
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: false,
}))
```

### 4.4 Teste Manual de CORS

**Validação positiva (origem permitida):**

```bash
curl -X OPTIONS http://localhost:8080/api/v1/membros \
  -H "Origin: http://localhost:5173" \
  -H "Access-Control-Request-Method: GET" \
  -v
```

Deve retornar:
```
< Access-Control-Allow-Origin: http://localhost:5173
```

**Validação negativa (origem bloqueada):**

```bash
curl -X OPTIONS http://localhost:8080/api/v1/membros \
  -H "Origin: https://evil.com" \
  -H "Access-Control-Request-Method: GET" \
  -v
```

Deve retornar sem o header `Access-Control-Allow-Origin`.

---

## 5. Parte 2: Rate Limiting no Login

### 5.1 Criar Middleware de Rate Limiting

**Arquivo:** `backend/internal/middleware/rate_limit.go`

```go
package middleware

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/ulule/limiter/v3"
    mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
    "github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimitMiddleware cria um middleware de rate limiting por IP
// rate: formato "X-Y" onde X é o número de requests e Y é o período (s, m, h)
// Exemplo: "5-M" = 5 requests por minuto
func RateLimitMiddleware(rate string) gin.HandlerFunc {
    // Parse do rate (ex: "5-M" = 5 requests por minuto)
    parsedRate, err := limiter.NewRateFromFormatted(rate)
    if err != nil {
        // Fallback seguro: 5 tentativas por minuto
        parsedRate = limiter.Rate{
            Period: 1 * time.Minute,
            Limit:  5,
        }
    }

    // Store em memória (para prod, considerar Redis)
    store := memory.NewStore()

    // Criar limiter
    lim := limiter.New(store, parsedRate)

    // Middleware do Gin
    middleware := mgin.NewMiddleware(lim)

    return func(c *gin.Context) {
        middleware(c)

        // Se o limiter bloqueou, retornar erro customizado
        if c.Writer.Header().Get("X-RateLimit-Remaining") == "0" {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error": "Muitas tentativas de login. Tente novamente mais tarde.",
            })
            return
        }
    }
}
```

### 5.2 Aplicar Rate Limiting no Login

**Arquivo:** `backend/cmd/server/main.go`

**Antes:**

```go
authGroup := v1.Group("/auth")
{
    authGroup.POST("/login", authHandler.Login)
    authGroup.POST("/refresh", authHandler.Refresh)
}
```

**Depois:**

```go
authGroup := v1.Group("/auth")
{
    // Rate limit: 5 tentativas por minuto por IP
    authGroup.POST("/login", middleware.RateLimitMiddleware("5-M"), authHandler.Login)
    authGroup.POST("/refresh", authHandler.Refresh)
}
```

### 5.3 Configuração via Env Var (Opcional)

Para tornar o rate limite configurável:

**Em `config/config.go`:**

```go
type Config struct {
    // ... campos existentes ...
    RateLimit struct {
        LoginRate string  // Ex: "5-M" (5 por minuto)
    }
}

func Load() (*Config, error) {
    // ... código existente ...

    cfg.RateLimit.LoginRate = getEnv("RATE_LIMIT_LOGIN", "5-M")

    return &cfg, nil
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}
```

**Uso em `main.go`:**

```go
authGroup.POST("/login",
    middleware.RateLimitMiddleware(cfg.RateLimit.LoginRate),
    authHandler.Login)
```

### 5.4 Teste Manual de Rate Limiting

**Script de teste (bash):**

```bash
#!/bin/bash
# Fazer 6 requests consecutivos ao endpoint de login

for i in {1..6}; do
  echo "Request $i:"
  curl -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"user1@test.com","password":"wrongpass"}' \
    -w "\nStatus: %{http_code}\n\n" \
    -s
  sleep 0.5
done
```

**Resultado esperado:**
- Requests 1-5: `HTTP 401 Unauthorized` (credenciais inválidas)
- Request 6: `HTTP 429 Too Many Requests` com mensagem de erro

**Headers esperados:**

```
X-RateLimit-Limit: 5
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1710259200
```

---

## 6. Testes e Validação

### 6.1 Testes Manuais

**Checklist de validação:**

- [ ] CORS permite `http://localhost:5173` em dev
- [ ] CORS bloqueia origens não listadas
- [ ] Rate limiting bloqueia após 5 tentativas
- [ ] Rate limiting reseta após 1 minuto
- [ ] Headers `X-RateLimit-*` presentes nas respostas
- [ ] Mensagem de erro clara no HTTP 429

### 6.2 Testes Automatizados (Opcional para TT-08)

**Arquivo:** `backend/internal/middleware/rate_limit_test.go`

```go
package middleware_test

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/felipehrs/financas-familiares/backend/internal/middleware"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware(t *testing.T) {
    gin.SetMode(gin.TestMode)

    r := gin.New()
    r.POST("/test", middleware.RateLimitMiddleware("3-M"), func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "ok"})
    })

    // Primeiras 3 requests devem passar
    for i := 0; i < 3; i++ {
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("POST", "/test", nil)
        req.RemoteAddr = "192.168.1.1:1234"  // Simular mesmo IP
        r.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
    }

    // 4ª request deve ser bloqueada
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/test", nil)
    req.RemoteAddr = "192.168.1.1:1234"
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusTooManyRequests, w.Code, "4th request should be rate limited")
}
```

---

## 7. Critério de Conclusão

### Validação Funcional

- [x] CORS configurado com whitelist de origens via `CORS_ALLOWED_ORIGINS`
- [x] `AllowAllOrigins: true` removido de `main.go`
- [x] Rate limiting ativo em `POST /api/v1/auth/login`
- [x] Limite configurável via `RATE_LIMIT_LOGIN` (default: "5-M")
- [x] HTTP 429 retornado após exceder limite
- [x] Headers `X-RateLimit-*` presentes nas respostas

### Validação Técnica

- [x] `go build ./...` — sem erros
- [x] `go test ./...` — todos os testes passando
- [x] Testes manuais de CORS (positivo + negativo) OK
- [x] Testes manuais de rate limiting (5 + 1) OK

### Documentação

- [x] `.env.example` atualizado com `CORS_ALLOWED_ORIGINS` e `RATE_LIMIT_LOGIN`
- [x] `backend/README.md` atualizado com seção de segurança explicando as novas env vars
- [x] `docs/produto/sprints.md` marcado TT-08 como ✅

---

## 8. Checklist de Implementação

### Fase 1: CORS Restrito

- [ ] Adicionar `CORS_ALLOWED_ORIGINS` em `.env` e `.env.example`
- [ ] Atualizar `config/config.go` com campo `CORS.AllowedOrigins`
- [ ] Implementar parsing da env var (split por vírgula)
- [ ] Substituir `AllowAllOrigins: true` por `AllowOrigins: cfg.CORS.AllowedOrigins` em `main.go`
- [ ] Testar CORS manualmente (curl com diferentes origens)
- [ ] Commit: `feat(backend): restringir CORS a origens específicas — TT-08 (1/2)`

### Fase 2: Rate Limiting

- [ ] Instalar dependência: `go get github.com/ulule/limiter/v3`
- [ ] Criar `backend/internal/middleware/rate_limit.go`
- [ ] Implementar `RateLimitMiddleware(rate string) gin.HandlerFunc`
- [ ] Adicionar `RATE_LIMIT_LOGIN` em `.env` e `.env.example` (default: "5-M")
- [ ] Atualizar `config/config.go` com campo `RateLimit.LoginRate`
- [ ] Aplicar middleware em `authGroup.POST("/login", ...)` no `main.go`
- [ ] Testar rate limiting manualmente (script bash com 6 requests)
- [ ] Validar headers `X-RateLimit-*` e mensagem de erro HTTP 429
- [ ] Commit: `feat(backend): adicionar rate limiting no login — TT-08 (2/2)`

### Fase 3: Documentação e Finalização

- [ ] Atualizar `backend/README.md` com seção "Segurança" explicando CORS e Rate Limiting
- [ ] Verificar `go build ./...` e `go test ./...`
- [ ] Marcar TT-08 como ✅ em `docs/produto/sprints.md`
- [ ] Commit: `docs: atualizar README com configurações de segurança — TT-08 ✅`
- [ ] Push e criar PR (ou merge direto se em branch de sprint)

---

## Arquivos Modificados (Estimativa)

| Arquivo | Tipo de Mudança |
|---------|-----------------|
| `backend/.env.example` | Adicionar `CORS_ALLOWED_ORIGINS` e `RATE_LIMIT_LOGIN` |
| `backend/config/config.go` | Adicionar campos `CORS` e `RateLimit` + parsing |
| `backend/cmd/server/main.go` | Substituir `AllowAllOrigins` + aplicar rate limit no login |
| `backend/internal/middleware/rate_limit.go` | **NOVO**: implementar middleware de rate limiting |
| `backend/go.mod` | Adicionar `github.com/ulule/limiter/v3` |
| `backend/go.sum` | Hashes das novas dependências |
| `backend/README.md` | Adicionar seção "Segurança" |
| `docs/produto/sprints.md` | Marcar TT-08 como ✅ |
| `docs/tecnico/05-plano-tt08-cors-rate-limiting.md` | **ESTE ARQUIVO** |

**Total estimado:** 8 arquivos (1 novo, 7 modificados)

---

## Considerações de Produção

### CORS

Em produção (Railway), configurar a env var:

```env
CORS_ALLOWED_ORIGINS=https://financas-familiares.vercel.app
```

Se o frontend tiver múltiplos domínios (ex: preview branches do Vercel), adicionar todos:

```env
CORS_ALLOWED_ORIGINS=https://financas-familiares.vercel.app,https://financas-familiares-preview.vercel.app
```

### Rate Limiting

**Memory store** (atual) funciona para deploy single-instance (Railway).

**Para escalar horizontalmente** (múltiplas instâncias), migrar para Redis:

```go
import "github.com/ulule/limiter/v3/drivers/store/redis"

redisClient := redis.NewClient(&redis.Options{
    Addr: cfg.Redis.URL,
})
store := redis.NewStore(redisClient)
```

Mas isso pode ficar para uma otimização futura (fora do escopo do TT-08).

---

## Referências

- [Documentação gin-contrib/cors](https://github.com/gin-contrib/cors)
- [Documentação ulule/limiter](https://github.com/ulule/limiter)
- [OWASP: Rate Limiting](https://cheatsheetseries.owasp.org/cheatsheets/Denial_of_Service_Cheat_Sheet.html#rate-limiting)
- [MDN: CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
