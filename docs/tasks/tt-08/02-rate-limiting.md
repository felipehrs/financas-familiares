# Tarefa TT-08.2: Implementar Rate Limiting no Login (TDD)

**Sprint:** 9
**Item:** TT-08
**Fase:** 2 de 3
**Estimativa:** 60 min (incluindo TDD)
**Dependências:** TT-08.1 (CORS Restrito)
**Abordagem:** Test-Driven Development (Red → Green → Refactor)

---

## Objetivo

Implementar middleware de rate limiting para proteger o endpoint `POST /api/v1/auth/login` contra ataques de força bruta, limitando tentativas de login por IP.

**Por quê?** Atualmente não há limitação de tentativas de login. Um atacante pode fazer força bruta sem restrições.

---

## Checklist de Implementação

### 1. Instalar Dependência

**Biblioteca:** `github.com/ulule/limiter/v3`

```bash
cd backend
go get github.com/ulule/limiter/v3
go get github.com/ulule/limiter/v3/drivers/store/memory
go get github.com/ulule/limiter/v3/drivers/middleware/gin
```

- [ ] Dependências instaladas
- [ ] `go.mod` e `go.sum` atualizados

---

### 2. Escrever Teste de Integração (TDD — Red Phase)

**IMPORTANTE:** Seguindo TDD, vamos escrever o teste ANTES do middleware.

**Criar arquivo:** `backend/internal/middleware/rate_limit_test.go`

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

	// Setup: criar rota com rate limit de 3 requests por minuto
	r := gin.New()
	r.POST("/test", middleware.RateLimitMiddleware("3-M"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	t.Run("permite requests dentro do limite", func(t *testing.T) {
		// Primeiras 3 requests devem passar
		for i := 0; i < 3; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/test", nil)
			req.RemoteAddr = "192.168.1.1:1234" // Simular mesmo IP
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)

			// Verificar headers de rate limiting
			assert.NotEmpty(t, w.Header().Get("X-RateLimit-Limit"))
			assert.NotEmpty(t, w.Header().Get("X-RateLimit-Remaining"))
		}
	})

	t.Run("bloqueia 4ª request (excede limite)", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		r.ServeHTTP(w, req)

		// 4ª request deve ser bloqueada com HTTP 429
		assert.Equal(t, http.StatusTooManyRequests, w.Code, "4th request should be rate limited")

		// Header de remaining deve ser 0
		assert.Equal(t, "0", w.Header().Get("X-RateLimit-Remaining"))
	})

	t.Run("IPs diferentes têm limites independentes", func(t *testing.T) {
		// Novo IP deve ter seu próprio limite
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "192.168.1.2:5678" // IP diferente
		r.ServeHTTP(w, req)

		// Deve passar (primeiro request deste IP)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
```

**Executar o teste (deve FALHAR — Red phase):**

```bash
cd backend
go test ./internal/middleware -v
```

**Resultado esperado:**

```
# Erro de compilação: RateLimitMiddleware não existe
./rate_limit_test.go:XX:XX: undefined: middleware.RateLimitMiddleware
FAIL    github.com/felipehrs/financas-familiares/backend/internal/middleware [build failed]
```

✅ **Teste falhando = Red phase concluída**

- [ ] Arquivo `rate_limit_test.go` criado
- [ ] Teste executado e FALHANDO (Red phase ✅)

---

### 3. Criar Middleware de Rate Limiting (Green Phase)

**Criar arquivo:** `backend/internal/middleware/rate_limit.go`

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
// rate: formato "X-Y" onde X é o número de requests e Y é o período (S, M, H, D)
// Exemplos:
//   - "5-M" = 5 requests por minuto
//   - "10-S" = 10 requests por segundo
//   - "100-H" = 100 requests por hora
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

	// Store em memória
	// Para produção com múltiplas instâncias, considerar Redis
	store := memory.NewStore()

	// Criar limiter
	lim := limiter.New(store, parsedRate)

	// Middleware do Gin
	middleware := mgin.NewMiddleware(lim)

	return func(c *gin.Context) {
		// Aplicar o middleware do limiter
		middleware(c)

		// Verificar se o limite foi atingido
		// O middleware do limiter já seta os headers X-RateLimit-*
		// e aborta com status 429 automaticamente
	}
}
```

**Pontos importantes:**
- O middleware usa o IP do cliente (`c.ClientIP()`) automaticamente
- Headers `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset` são adicionados automaticamente
- HTTP 429 é retornado automaticamente quando o limite é excedido

- [ ] Arquivo `rate_limit.go` criado
- [ ] Função `RateLimitMiddleware` implementada
- [ ] Imports corretos

**Executar o teste novamente (deve PASSAR — Green phase):**

```bash
cd backend
go test ./internal/middleware -v
```

**Resultado esperado:**

```
=== RUN   TestRateLimitMiddleware
=== RUN   TestRateLimitMiddleware/permite_requests_dentro_do_limite
=== RUN   TestRateLimitMiddleware/bloqueia_4ª_request_(excede_limite)
=== RUN   TestRateLimitMiddleware/IPs_diferentes_têm_limites_independentes
--- PASS: TestRateLimitMiddleware (0.XXs)
    --- PASS: TestRateLimitMiddleware/permite_requests_dentro_do_limite (0.XXs)
    --- PASS: TestRateLimitMiddleware/bloqueia_4ª_request_(excede_limite) (0.XXs)
    --- PASS: TestRateLimitMiddleware/IPs_diferentes_têm_limites_independentes (0.XXs)
PASS
ok      github.com/felipehrs/financas-familiares/backend/internal/middleware   0.XXXs
```

✅ **Teste passando = Green phase concluída**

- [ ] Teste passando (Green phase ✅)

---

### 4. Refatorar (Refactor Phase)

Revisar o middleware:
- Código está limpo e legível?
- Comentários adequados?
- Tratamento de erro do parse está adequado?

Para este caso, o código está bem estruturado. **Refactor opcional.**

- [ ] Code review interno (refactoring se necessário)

---

### 5. Adicionar Variável de Ambiente (Opcional)

Para tornar o rate limite configurável:

**Arquivo:** `backend/.env`

```env
# Rate Limiting: formato "X-Y" (X requests por Y período: S, M, H, D)
RATE_LIMIT_LOGIN=5-M
```

**Arquivo:** `backend/.env.example`

```env
# Rate Limiting: formato "X-Y" (X requests por Y período: S, M, H, D)
# Exemplo: "5-M" = 5 tentativas por minuto
RATE_LIMIT_LOGIN=5-M
```

- [ ] `.env` atualizado
- [ ] `.env.example` atualizado

---

### 4. Atualizar Configuração (Opcional)

**Arquivo:** `backend/config/config.go`

**Passo 4.1:** Adicionar campo `RateLimit` na struct

```go
type Config struct {
    // ... campos existentes ...
    RateLimit struct {
        LoginRate string  // Ex: "5-M" (5 por minuto)
    }
}
```

**Passo 4.2:** Implementar parsing

Na função `Load()`, adicionar:

```go
// Parse rate limiting config
cfg.RateLimit.LoginRate = os.Getenv("RATE_LIMIT_LOGIN")
if cfg.RateLimit.LoginRate == "" {
    cfg.RateLimit.LoginRate = "5-M"  // Default: 5 por minuto
}
```

- [ ] Campo `RateLimit` adicionado à struct `Config`
- [ ] Parsing implementado em `Load()`

---

### 5. Aplicar Middleware no Endpoint de Login

**Arquivo:** `backend/cmd/server/main.go`

**ANTES (linhas ~132-136):**

```go
authGroup := v1.Group("/auth")
{
    authGroup.POST("/login", authHandler.Login)
    authGroup.POST("/refresh", authHandler.Refresh)
}
```

**DEPOIS (com rate limiting configurável):**

```go
authGroup := v1.Group("/auth")
{
    // Rate limit no login: proteger contra força bruta
    authGroup.POST("/login",
        middleware.RateLimitMiddleware(cfg.RateLimit.LoginRate),
        authHandler.Login)

    // Refresh token sem rate limit (já é protegido pelo token válido)
    authGroup.POST("/refresh", authHandler.Refresh)
}
```

**OU (se não implementou configuração opcional):**

```go
authGroup := v1.Group("/auth")
{
    // Rate limit: 5 tentativas por minuto por IP
    authGroup.POST("/login",
        middleware.RateLimitMiddleware("5-M"),
        authHandler.Login)

    authGroup.POST("/refresh", authHandler.Refresh)
}
```

- [ ] Middleware aplicado em `POST /auth/login`
- [ ] `/auth/refresh` permanece sem rate limiting

---

### 6. Verificar Compilação

```bash
cd backend
go build ./...
```

- [ ] `go build ./...` executado com sucesso

---

### 7. Testar Rate Limiting Manualmente

**Pré-requisito:** Servidor rodando (`go run cmd/server/main.go`)

**Teste: Fazer 6 requisições consecutivas**

Criar um script bash `test-rate-limit.sh`:

```bash
#!/bin/bash
# Testar rate limiting (5 por minuto = 6ª request deve falhar)

echo "=== Teste de Rate Limiting - POST /api/v1/auth/login ==="
echo

for i in {1..6}; do
  echo "Request $i:"

  response=$(curl -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"user1@test.com","password":"wrongpass"}' \
    -w "\nHTTP Status: %{http_code}\n" \
    -s -i)

  # Extrair headers de rate limit
  echo "$response" | grep -i "x-ratelimit"
  echo "$response" | grep "HTTP Status"

  # Extrair mensagem de erro se houver
  echo "$response" | grep -o '"error":"[^"]*"' || echo "(sem erro de rate limit)"

  echo "---"
  sleep 0.5
done

echo
echo "✅ Resultado esperado:"
echo "  - Requests 1-5: HTTP 401 (credenciais inválidas)"
echo "  - Request 6: HTTP 429 (Too Many Requests)"
```

Executar:

```bash
chmod +x test-rate-limit.sh
./test-rate-limit.sh
```

**Resultado esperado:**

```
Request 1:
X-RateLimit-Limit: 5
X-RateLimit-Remaining: 4
HTTP Status: 401
---

Request 2:
X-RateLimit-Limit: 5
X-RateLimit-Remaining: 3
HTTP Status: 401
---

... (requests 3-5 similares) ...

Request 6:
X-RateLimit-Limit: 5
X-RateLimit-Remaining: 0
HTTP Status: 429
---
```

---

**Teste manual simplificado (curl direto):**

```bash
# Request 1-5 devem retornar HTTP 401
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"wrong"}' \
  -v

# Repetir 5 vezes...

# Request 6 deve retornar HTTP 429
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"wrong"}' \
  -v
```

- [ ] Requests 1-5: HTTP 401 (credenciais inválidas)
- [ ] Request 6: HTTP 429 (Too Many Requests)
- [ ] Headers `X-RateLimit-*` presentes em todas as respostas
- [ ] Após ~60 segundos, rate limit reseta e permite novas tentativas

---

### 8. Testar Rate Limiting com Credenciais Corretas

**Validação adicional:** Rate limiting deve bloquear mesmo com credenciais válidas (proteção contra force brute em geral).

```bash
# Fazer 6 logins com credenciais CORRETAS
for i in {1..6}; do
  curl -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"user1@test.com","password":"senha123"}' \
    -w "\nStatus: %{http_code}\n"
  sleep 0.5
done
```

**Resultado esperado:**
- Requests 1-5: HTTP 200 (login OK)
- Request 6: HTTP 429 (bloqueado mesmo com senha correta)

- [ ] Rate limiting aplica mesmo para credenciais corretas

---

## Critério de Conclusão (TDD)

- [x] **Red Phase:** Teste `TestRateLimitMiddleware` criado e FALHANDO
- [x] **Green Phase:** Middleware implementado e teste PASSANDO (3 cenários)
- [x] **Refactor Phase:** Code review realizado (opcional neste caso)
- [x] Dependência `github.com/ulule/limiter/v3` instalada
- [x] Middleware `rate_limit.go` criado e implementado
- [x] `go test ./internal/middleware` executado com sucesso (3 subtestes)
- [x] Env var `RATE_LIMIT_LOGIN` adicionada (se implementou config opcional)
- [x] Configuração em `config.go` atualizada (se implementou config opcional)
- [x] Middleware aplicado em `POST /auth/login`
- [x] `go build ./...` executado com sucesso
- [x] Teste manual: 5 requests OK, 6ª bloqueada com HTTP 429
- [x] Headers `X-RateLimit-*` presentes nas respostas

---

## Commit Sugerido

```bash
git add backend/internal/middleware/rate_limit.go \
        backend/internal/middleware/rate_limit_test.go \
        backend/.env.example \
        backend/config/config.go \
        backend/cmd/server/main.go \
        backend/go.mod \
        backend/go.sum

git commit -m "feat(backend): adicionar rate limiting no login — TT-08 (2/3)

Implementa rate limiting (5 tentativas/minuto) no endpoint de autenticação
para proteger contra ataques de força bruta.

**Mudanças:**
- Criar middleware RateLimitMiddleware usando ulule/limiter
- Adicionar env var RATE_LIMIT_LOGIN (default: 5-M)
- Aplicar middleware em POST /auth/login
- Headers X-RateLimit-* incluídos nas respostas

**Segurança:**
- Bloqueia força bruta após 5 tentativas/minuto por IP
- Retorna HTTP 429 quando limite excedido
- Rate limit reseta após 60 segundos

**Testes (TDD):**
- TestRateLimitMiddleware: 3 cenários
  1. Permite requests dentro do limite
  2. Bloqueia 4ª request (excede limite)
  3. IPs diferentes têm limites independentes
- Teste manual: 5 OK / 6ª bloqueada ✅

**Dependências:**
- github.com/ulule/limiter/v3

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Próximos Passos

Após conclusão desta tarefa, prosseguir para:
- **[03-testes-e-docs.md](./03-testes-e-docs.md)** — Documentação e finalização do TT-08

---

## Referências

- [Plano técnico TT-08](../../tecnico/05-plano-tt08-cors-rate-limiting.md)
- [ulule/limiter docs](https://github.com/ulule/limiter)
- [OWASP: Brute Force Protection](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html#protect-against-automated-attacks)
