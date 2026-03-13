# Tarefa TT-08.1: Implementar CORS Restrito (TDD)

**Sprint:** 9
**Item:** TT-08
**Fase:** 1 de 3
**Estimativa:** 45 min (incluindo TDD)
**Dependências:** Nenhuma
**Abordagem:** Test-Driven Development (Red → Green → Refactor)

---

## Objetivo

Substituir a configuração atual de CORS (`AllowAllOrigins: true`) por uma whitelist de origens permitidas, configurável via variável de ambiente.

**Por quê?** Atualmente qualquer site malicioso pode fazer requisições HTTP para a API. Mesmo com `AllowCredentials: false`, JWTs no header `Authorization` ainda estão expostos.

---

## Checklist de Implementação

### 1. Adicionar Variáveis de Ambiente

**Arquivo:** `backend/.env`

Adicionar no final do arquivo:

```env
# CORS: lista de origens permitidas, separadas por vírgula
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
```

**Arquivo:** `backend/.env.example`

Adicionar a mesma linha para documentar a configuração:

```env
# CORS: lista de origens permitidas, separadas por vírgula
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
```

- [ ] `.env` atualizado
- [ ] `.env.example` atualizado

---

### 2. Atualizar Configuração (`config/config.go`)

**Localização:** `backend/config/config.go`

**Passo 2.1:** Adicionar campo `CORS` na struct `Config`

Encontrar a struct `Config` e adicionar:

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
        AllowedOrigins []string  // ✅ Nova propriedade
    }
    Server struct {
        Port string
    }
}
```

**Passo 2.2:** Implementar parsing da env var

Na função `Load()`, adicionar (antes do `return &cfg, nil`):

```go
// Parse CORS allowed origins
originsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
if originsStr == "" {
    // Default para desenvolvimento local
    originsStr = "http://localhost:5173"
}

// Split por vírgula e trim espaços
cfg.CORS.AllowedOrigins = strings.Split(originsStr, ",")
for i := range cfg.CORS.AllowedOrigins {
    cfg.CORS.AllowedOrigins[i] = strings.TrimSpace(cfg.CORS.AllowedOrigins[i])
}
```

**IMPORTANTE:** Verificar se o import `strings` já existe. Se não, adicionar no topo:

```go
import (
    // ... outros imports ...
    "strings"
)
```

- [ ] Struct `Config` atualizada com campo `CORS`
- [ ] Parsing implementado na função `Load()`
- [ ] Import `strings` adicionado (se necessário)

---

### 3. Atualizar `main.go`

**Localização:** `backend/cmd/server/main.go`

**Encontrar o bloco CORS atual (linhas ~113-119):**

```go
r.Use(cors.New(cors.Config{
    AllowAllOrigins:  true,  // ❌ REMOVER esta linha
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: false,
}))
```

**Substituir por:**

```go
r.Use(cors.New(cors.Config{
    AllowOrigins:     cfg.CORS.AllowedOrigins,  // ✅ Whitelist específica
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: false,
}))
```

**Mudanças:**
- `AllowAllOrigins: true` → **REMOVIDO**
- **ADICIONADO:** `AllowOrigins: cfg.CORS.AllowedOrigins`

- [ ] `AllowAllOrigins` removido
- [ ] `AllowOrigins` adicionado com referência a `cfg.CORS.AllowedOrigins`

---

### 4. Escrever Teste (TDD — Red Phase)

**IMPORTANTE:** Este projeto segue **TDD**. Vamos escrever o teste ANTES da implementação.

**Criar arquivo:** `backend/config/config_test.go`

```go
package config_test

import (
	"os"
	"testing"

	"github.com/felipehrs/financas-familiares/backend/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad_CORSAllowedOrigins(t *testing.T) {
	tests := []struct {
		name            string
		envValue        string
		expectedOrigins []string
	}{
		{
			name:            "múltiplas origens separadas por vírgula",
			envValue:        "http://localhost:5173,http://localhost:3000,https://app.vercel.app",
			expectedOrigins: []string{"http://localhost:5173", "http://localhost:3000", "https://app.vercel.app"},
		},
		{
			name:            "origem única",
			envValue:        "http://localhost:5173",
			expectedOrigins: []string{"http://localhost:5173"},
		},
		{
			name:            "origens com espaços extras (deve trimmar)",
			envValue:        "http://localhost:5173 , http://localhost:3000 , https://app.vercel.app",
			expectedOrigins: []string{"http://localhost:5173", "http://localhost:3000", "https://app.vercel.app"},
		},
		{
			name:            "env var vazia (deve usar default)",
			envValue:        "",
			expectedOrigins: []string{"http://localhost:5173"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup: configurar env vars
			if tt.envValue != "" {
				os.Setenv("CORS_ALLOWED_ORIGINS", tt.envValue)
			} else {
				os.Unsetenv("CORS_ALLOWED_ORIGINS")
			}
			defer os.Unsetenv("CORS_ALLOWED_ORIGINS")

			// Env vars obrigatórias para Load() não falhar
			os.Setenv("DATABASE_URL", "postgres://test")
			os.Setenv("JWT_SECRET", "test-secret")
			defer os.Unsetenv("DATABASE_URL")
			defer os.Unsetenv("JWT_SECRET")

			// Execute
			cfg, err := config.Load()

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOrigins, cfg.CORS.AllowedOrigins)
		})
	}
}
```

**Executar o teste (deve FALHAR — Red phase):**

```bash
cd backend
go test ./config -v
```

**Resultado esperado:**

```
# Erro de compilação: cfg.CORS field não existe
./config_test.go:XX:XX: cfg.CORS undefined (type *config.Config has no field or method CORS)
FAIL    github.com/felipehrs/financas-familiares/backend/config [build failed]
```

✅ **Teste falhando = Red phase concluída**

- [ ] Arquivo `config/config_test.go` criado
- [ ] Teste executado e FALHANDO (Red phase ✅)

---

### 5. Implementar Código (Green Phase)

Agora implementar o código mínimo para fazer o teste passar.

**Já descrito anteriormente:**
- Adicionar campo `CORS` na struct `Config`
- Implementar parsing em `Load()`
- Adicionar import `strings`

**Executar o teste novamente (deve PASSAR — Green phase):**

```bash
cd backend
go test ./config -v
```

**Resultado esperado:**

```
=== RUN   TestLoad_CORSAllowedOrigins
=== RUN   TestLoad_CORSAllowedOrigins/múltiplas_origens_separadas_por_vírgula
=== RUN   TestLoad_CORSAllowedOrigins/origem_única
=== RUN   TestLoad_CORSAllowedOrigins/origens_com_espaços_extras_(deve_trimmar)
=== RUN   TestLoad_CORSAllowedOrigins/env_var_vazia_(deve_usar_default)
--- PASS: TestLoad_CORSAllowedOrigins (0.00s)
    --- PASS: TestLoad_CORSAllowedOrigins/múltiplas_origens_separadas_por_vírgula (0.00s)
    --- PASS: TestLoad_CORSAllowedOrigins/origem_única (0.00s)
    --- PASS: TestLoad_CORSAllowedOrigins/origens_com_espaços_extras_(deve_trimmar) (0.00s)
    --- PASS: TestLoad_CORSAllowedOrigins/env_var_vazia_(deve_usar_default) (0.00s)
PASS
ok      github.com/felipehrs/financas-familiares/backend/config        0.XXXs
```

✅ **Teste passando = Green phase concluída**

- [ ] Código implementado em `config/config.go`
- [ ] Teste passando (Green phase ✅)

---

### 6. Refatorar (Refactor Phase)

Revisar o código e melhorar se houver oportunidade:
- Extrair lógica de parsing para função auxiliar?
- Adicionar comentários GoDoc?
- Validar formato das URLs?

Para este caso, o código já está simples. **Refactor opcional.**

- [ ] Code review interno (refactoring se necessário)

---

### 7. Verificar Compilação Geral

```bash
cd backend
go build ./...
```

- [ ] `go build ./...` executado com sucesso

---

### 8. Testar Manualmente (CORS)

**Teste 1: Origem PERMITIDA (deve retornar header CORS)**

```bash
curl -X OPTIONS http://localhost:8080/api/v1/membros \
  -H "Origin: http://localhost:5173" \
  -H "Access-Control-Request-Method: GET" \
  -v
```

**Resultado esperado:**

```
< HTTP/1.1 204 No Content
< Access-Control-Allow-Origin: http://localhost:5173
< Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
```

---

**Teste 2: Origem BLOQUEADA (NÃO deve retornar header CORS)**

```bash
curl -X OPTIONS http://localhost:8080/api/v1/membros \
  -H "Origin: https://evil.com" \
  -H "Access-Control-Request-Method: GET" \
  -v
```

**Resultado esperado:**

```
< HTTP/1.1 204 No Content
(SEM o header Access-Control-Allow-Origin)
```

O servidor responde, mas **não inclui** o header `Access-Control-Allow-Origin`, o que faz o navegador bloquear a resposta.

- [ ] Teste com origem permitida: ✅ header CORS presente
- [ ] Teste com origem bloqueada: ❌ header CORS ausente

---

## Critério de Conclusão (TDD)

- [x] **Red Phase:** Teste `TestLoad_CORSAllowedOrigins` criado e FALHANDO
- [x] **Green Phase:** Código implementado e teste PASSANDO
- [x] **Refactor Phase:** Code review realizado (opcional neste caso)
- [x] Env vars `CORS_ALLOWED_ORIGINS` adicionadas em `.env` e `.env.example`
- [x] `config/config.go` atualizado com campo `CORS.AllowedOrigins` e parsing
- [x] `main.go` usando `AllowOrigins` com whitelist (sem `AllowAllOrigins`)
- [x] `go test ./config` executado com sucesso (4 cenários de teste)
- [x] `go build ./...` executado com sucesso
- [x] Testes manuais de CORS validados (positivo + negativo)

---

## Commit Sugerido

```bash
git add backend/.env.example \
        backend/config/config.go \
        backend/config/config_test.go \
        backend/cmd/server/main.go

git commit -m "feat(backend): restringir CORS a origens específicas — TT-08 (1/3)

Substitui AllowAllOrigins por whitelist configurável via CORS_ALLOWED_ORIGINS.

**Mudanças:**
- Adicionar env var CORS_ALLOWED_ORIGINS (default: localhost:5173)
- Atualizar config.go com campo CORS.AllowedOrigins + parsing
- Substituir AllowAllOrigins: true por AllowOrigins: cfg.CORS.AllowedOrigins

**Segurança:**
- Bloqueia requisições de origens não autorizadas
- Mitiga ataques CSRF de sites maliciosos

**Testes (TDD):**
- TestLoad_CORSAllowedOrigins: 4 cenários (múltiplas origens, única, com espaços, default)
- Teste manual: origem permitida ✅ / bloqueada ❌

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Próximos Passos

Após conclusão desta tarefa, prosseguir para:
- **[02-rate-limiting.md](./02-rate-limiting.md)** — Implementar rate limiting no login

---

## Referências

- [Plano técnico TT-08](../../tecnico/05-plano-tt08-cors-rate-limiting.md)
- [MDN: CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [gin-contrib/cors](https://github.com/gin-contrib/cors)
