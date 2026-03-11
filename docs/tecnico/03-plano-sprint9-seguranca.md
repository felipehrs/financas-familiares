# Plano de Execução — Sprint 9: Segurança e Isolamento de Dados

## Índice
1. [Visão Geral e Sequência de Execução](#1-visão-geral-e-sequência-de-execução)
2. [TT-07: Isolamento de Dados por Família](#2-tt-07-isolamento-de-dados-por-família)
3. [TT-08: CORS Restrito e Rate Limiting no Login](#3-tt-08-cors-restrito-e-rate-limiting-no-login)
4. [TT-09: Índices de Performance](#4-tt-09-índices-de-performance)
5. [Critério de Conclusão](#5-critério-de-conclusão)
6. [Arquivos Modificados por Item](#6-arquivos-modificados-por-item)

---

## 1. Visão Geral e Sequência de Execução

### Por que TT-07 → TT-08 → TT-09?

**TT-07 primeiro** porque é o problema de segurança mais crítico: sem isolamento, qualquer usuário autenticado lê e modifica dados de qualquer outro. Além disso, o TT-09 depende do TT-07: os índices em `familia_id` só fazem sentido depois que a coluna existir.

**TT-08 segundo** porque é independente do banco de dados (opera na camada de middleware HTTP). Pode ser feito em paralelo com TT-07, mas sequencialmente é mais simples.

**TT-09 por último** porque os índices em `familia_id` dependem da migration do TT-07 estar aplicada.

### Por que `familia_id` e não `usuario_id`?

O domínio do produto é **finanças familiares** — múltiplos usuários da mesma família devem ver os mesmos dados (despesas do cônjuge, cartões compartilhados, etc.). Isolar por `usuario_id` quebraria essa colaboração. A escolha correta é isolar por `familia_id`:

- Dados pertencem à **família**, não ao usuário individual
- Todos os membros da família veem e editam os mesmos registros
- Famílias diferentes são completamente isoladas entre si
- O `usuario_id` continua existindo para autenticação (JWT), mas não para filtrar dados

---

## 2. TT-07: Isolamento de Dados por Família

### 2.1 Análise do Estado Atual

- O middleware em `backend/internal/middleware/auth.go` define `"userID"` no contexto Gin — mas nenhum handler usa essa chave.
- Todos os repositórios operam sem filtro de tenant.
- A tabela `refresh_tokens` já possui `usuario_id` — **esta coluna permanece**, pois tokens são por usuário, não por família.
- O `SeedCategorias` insere categorias globais sem vínculo de família.
- O JWT atual contém apenas `sub` (usuarioID) — precisa passar a incluir `familia_id`.

### 2.2 Estratégia de Migration

**Decisão: uma migration única (`000017`) que:**
1. Cria as tabelas `familias` e `familia_usuarios`
2. Adiciona `familia_id` em todas as tabelas de dados
3. Popula `familia_id` nos registros existentes (seed de dev)
4. Aplica `NOT NULL` e FKs
5. Substitui a constraint `UNIQUE(nome)` de categorias por `UNIQUE(familia_id, nome)`

Uma migration única garante atomicidade: ou tudo aplica, ou nada.

### 2.3 Migration 000017

#### `backend/migrations/000017_add_familia_isolation.up.sql`

```sql
-- Sprint 9 / TT-07: Isolamento de dados por família
-- Estratégia para banco com dados existentes:
--   1. Criar tabelas familias e familia_usuarios
--   2. Inserir uma família padrão vinculada ao primeiro usuário
--   3. Adicionar familia_id nas tabelas de dados como nullable
--   4. Popular familia_id com a família criada
--   5. Aplicar NOT NULL e FK
-- Em banco limpo (produção), a lógica de UPDATE não afeta nenhuma linha.

-- =====================================================================
-- NOVAS TABELAS DE TENANT
-- =====================================================================
CREATE TABLE IF NOT EXISTS familias (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nome       VARCHAR(100) NOT NULL,
    owner_id   UUID NOT NULL REFERENCES usuarios(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS familia_usuarios (
    familia_id UUID NOT NULL REFERENCES familias(id),
    usuario_id UUID NOT NULL REFERENCES usuarios(id),
    role       VARCHAR(20) NOT NULL DEFAULT 'membro'
                CHECK (role IN ('admin', 'membro')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (familia_id, usuario_id)
);

-- Criar família padrão para o primeiro usuário (seed de dev)
-- Em banco limpo de produção, este bloco não cria nada (sem usuários = sem família)
DO $$
DECLARE
    v_usuario_id UUID;
    v_familia_id UUID;
BEGIN
    SELECT id INTO v_usuario_id FROM usuarios ORDER BY created_at ASC LIMIT 1;
    IF v_usuario_id IS NOT NULL THEN
        INSERT INTO familias (nome, owner_id)
        VALUES ('Família Principal', v_usuario_id)
        RETURNING id INTO v_familia_id;

        -- Vincular todos os usuários existentes à família
        INSERT INTO familia_usuarios (familia_id, usuario_id, role)
        SELECT v_familia_id, id, CASE WHEN id = v_usuario_id THEN 'admin' ELSE 'membro' END
        FROM usuarios
        WHERE deleted_at IS NULL
        ON CONFLICT DO NOTHING;
    END IF;
END $$;

-- =====================================================================
-- MEMBROS
-- =====================================================================
ALTER TABLE membros ADD COLUMN IF NOT EXISTS familia_id UUID;

UPDATE membros
SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1)
WHERE familia_id IS NULL;

ALTER TABLE membros
    ALTER COLUMN familia_id SET NOT NULL,
    ADD CONSTRAINT membros_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

-- =====================================================================
-- CATEGORIAS
-- =====================================================================
ALTER TABLE categorias ADD COLUMN IF NOT EXISTS familia_id UUID;

UPDATE categorias
SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1)
WHERE familia_id IS NULL;

ALTER TABLE categorias
    ALTER COLUMN familia_id SET NOT NULL,
    ADD CONSTRAINT categorias_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

-- Substituir UNIQUE(nome) global por UNIQUE(familia_id, nome)
ALTER TABLE categorias DROP CONSTRAINT IF EXISTS categorias_nome_key;
ALTER TABLE categorias ADD CONSTRAINT categorias_nome_familia_unique UNIQUE (familia_id, nome);

-- =====================================================================
-- CARTOES_CREDITO
-- =====================================================================
ALTER TABLE cartoes_credito ADD COLUMN IF NOT EXISTS familia_id UUID;

UPDATE cartoes_credito
SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1)
WHERE familia_id IS NULL;

ALTER TABLE cartoes_credito
    ALTER COLUMN familia_id SET NOT NULL,
    ADD CONSTRAINT cartoes_credito_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

-- =====================================================================
-- DESPESAS_CARTAO
-- =====================================================================
ALTER TABLE despesas_cartao ADD COLUMN IF NOT EXISTS familia_id UUID;

UPDATE despesas_cartao
SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1)
WHERE familia_id IS NULL;

ALTER TABLE despesas_cartao
    ALTER COLUMN familia_id SET NOT NULL,
    ADD CONSTRAINT despesas_cartao_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

-- =====================================================================
-- ASSINATURAS
-- =====================================================================
ALTER TABLE assinaturas ADD COLUMN IF NOT EXISTS familia_id UUID;

UPDATE assinaturas
SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1)
WHERE familia_id IS NULL;

ALTER TABLE assinaturas
    ALTER COLUMN familia_id SET NOT NULL,
    ADD CONSTRAINT assinaturas_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

-- =====================================================================
-- CONTAS_FIXAS
-- =====================================================================
ALTER TABLE contas_fixas ADD COLUMN IF NOT EXISTS familia_id UUID;

UPDATE contas_fixas
SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1)
WHERE familia_id IS NULL;

ALTER TABLE contas_fixas
    ALTER COLUMN familia_id SET NOT NULL,
    ADD CONSTRAINT contas_fixas_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

-- =====================================================================
-- DESPESAS_GERAIS
-- =====================================================================
ALTER TABLE despesas_gerais ADD COLUMN IF NOT EXISTS familia_id UUID;

UPDATE despesas_gerais
SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1)
WHERE familia_id IS NULL;

ALTER TABLE despesas_gerais
    ALTER COLUMN familia_id SET NOT NULL,
    ADD CONSTRAINT despesas_gerais_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

-- =====================================================================
-- RENDAS_FIXAS
-- =====================================================================
ALTER TABLE rendas_fixas ADD COLUMN IF NOT EXISTS familia_id UUID;

UPDATE rendas_fixas
SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1)
WHERE familia_id IS NULL;

ALTER TABLE rendas_fixas
    ALTER COLUMN familia_id SET NOT NULL,
    ADD CONSTRAINT rendas_fixas_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

-- =====================================================================
-- RENDAS_VARIAVEIS
-- =====================================================================
ALTER TABLE rendas_variaveis ADD COLUMN IF NOT EXISTS familia_id UUID;

UPDATE rendas_variaveis
SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1)
WHERE familia_id IS NULL;

ALTER TABLE rendas_variaveis
    ALTER COLUMN familia_id SET NOT NULL,
    ADD CONSTRAINT rendas_variaveis_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

-- =====================================================================
-- RENDAS_EXTRAS
-- =====================================================================
ALTER TABLE rendas_extras ADD COLUMN IF NOT EXISTS familia_id UUID;

UPDATE rendas_extras
SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1)
WHERE familia_id IS NULL;

ALTER TABLE rendas_extras
    ALTER COLUMN familia_id SET NOT NULL,
    ADD CONSTRAINT rendas_extras_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

-- =====================================================================
-- RENDIMENTOS_INVESTIMENTO
-- =====================================================================
ALTER TABLE rendimentos_investimento ADD COLUMN IF NOT EXISTS familia_id UUID;

UPDATE rendimentos_investimento
SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1)
WHERE familia_id IS NULL;

ALTER TABLE rendimentos_investimento
    ALTER COLUMN familia_id SET NOT NULL,
    ADD CONSTRAINT rendimentos_investimento_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);
```

> **Nota sobre produção (Railway):** Se o banco não tiver usuários cadastrados, o bloco `DO $$` não executa nada — `v_usuario_id` será NULL. Os `UPDATE ... WHERE familia_id IS NULL` também não afetarão nenhuma linha. A migration aplica limpa.

#### `backend/migrations/000017_add_familia_isolation.down.sql`

```sql
-- Rollback: remover familia_id de todas as tabelas e apagar tabelas de tenant

ALTER TABLE rendimentos_investimento
    DROP CONSTRAINT IF EXISTS rendimentos_investimento_familia_id_fkey,
    DROP COLUMN IF EXISTS familia_id;

ALTER TABLE rendas_extras
    DROP CONSTRAINT IF EXISTS rendas_extras_familia_id_fkey,
    DROP COLUMN IF EXISTS familia_id;

ALTER TABLE rendas_variaveis
    DROP CONSTRAINT IF EXISTS rendas_variaveis_familia_id_fkey,
    DROP COLUMN IF EXISTS familia_id;

ALTER TABLE rendas_fixas
    DROP CONSTRAINT IF EXISTS rendas_fixas_familia_id_fkey,
    DROP COLUMN IF EXISTS familia_id;

ALTER TABLE despesas_gerais
    DROP CONSTRAINT IF EXISTS despesas_gerais_familia_id_fkey,
    DROP COLUMN IF EXISTS familia_id;

ALTER TABLE contas_fixas
    DROP CONSTRAINT IF EXISTS contas_fixas_familia_id_fkey,
    DROP COLUMN IF EXISTS familia_id;

ALTER TABLE assinaturas
    DROP CONSTRAINT IF EXISTS assinaturas_familia_id_fkey,
    DROP COLUMN IF EXISTS familia_id;

ALTER TABLE despesas_cartao
    DROP CONSTRAINT IF EXISTS despesas_cartao_familia_id_fkey,
    DROP COLUMN IF EXISTS familia_id;

ALTER TABLE cartoes_credito
    DROP CONSTRAINT IF EXISTS cartoes_credito_familia_id_fkey,
    DROP COLUMN IF EXISTS familia_id;

ALTER TABLE categorias
    DROP CONSTRAINT IF EXISTS categorias_nome_familia_unique,
    DROP CONSTRAINT IF EXISTS categorias_familia_id_fkey,
    DROP COLUMN IF EXISTS familia_id;

ALTER TABLE categorias ADD CONSTRAINT categorias_nome_key UNIQUE (nome);

ALTER TABLE membros
    DROP CONSTRAINT IF EXISTS membros_familia_id_fkey,
    DROP COLUMN IF EXISTS familia_id;

DROP TABLE IF EXISTS familia_usuarios;
DROP TABLE IF EXISTS familias;
```

### 2.4 Ajuste do JWT e Middleware

O JWT atual contém apenas `sub` (usuarioID). Para suportar isolamento por família, precisa incluir `familia_id`.

#### Auth Service — Claims do JWT

```go
// backend/internal/domain/auth.go (ou onde Claims estiver definido)
type Claims struct {
    UsuarioID string `json:"sub"`
    FamiliaID string `json:"familia_id"`  // NOVO
    jwt.RegisteredClaims
}
```

#### Auth Service — Geração do Token no Login

Ao fazer login, buscar a família do usuário e incluir no token:

```go
// No AuthService.Login, após validar credenciais:
familiaID, err := s.familiaRepo.BuscarFamiliaPorUsuario(usuarioID)
if err != nil {
    return nil, domain.ErrFamiliaNaoEncontrada
}

claims := domain.Claims{
    UsuarioID: usuarioID,
    FamiliaID: familiaID,
    RegisteredClaims: jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
    },
}
```

#### Middleware — Disponibilizar `familiaID` no Contexto

```go
// backend/internal/middleware/auth.go
// Após validar o token, extrair FamiliaID das claims:
c.Set("userID", claims.UsuarioID)    // mantido para compatibilidade
c.Set("familiaID", claims.FamiliaID) // NOVO
```

#### Handlers — Extração do Contexto

Convenção adotada: chave `"familiaID"` no contexto Gin.

```go
// Padrão em todos os handlers que acessam dados:
familiaID, ok := c.Get("familiaID")
if !ok {
    c.JSON(http.StatusUnauthorized, gin.H{"erro": "não autenticado"})
    return
}
fid := familiaID.(string)
```

Função auxiliar recomendada em `handler/helpers.go`:

```go
func getFamiliaID(c *gin.Context) (string, bool) {
    v, ok := c.Get("familiaID")
    if !ok {
        return "", false
    }
    fid, ok := v.(string)
    return fid, ok
}
```

### 2.5 Refatoração dos Repositórios

O padrão de refatoração é idêntico ao anterior, mas o parâmetro muda de `usuarioID` para `familiaID` — e a semântica muda: agora todos os membros da família veem os mesmos dados.

#### 2.5.1 Interfaces a atualizar (em `service/*.go`)

| Interface (em service/) | Métodos que mudam |
|---|---|
| `MembroRepository` | `Criar`, `BuscarPorID`, `Listar`, `Atualizar`, `Inativar` |
| `CategoriaRepository` | `Criar`, `Listar`, `BuscarPorID`, `Atualizar`, `Excluir` |
| `CartaoCreditoRepository` | `Criar`, `BuscarPorID`, `Listar`, `Atualizar`, `Inativar` |
| `DespesaCartaoRepository` | `Criar`, `ListarPorCartao`, `ListarPorFatura`, `BuscarPorID`, `ListarPorFaturaGlobal`, `Excluir` |
| `AssinaturaRepository` | `Criar`, `BuscarPorID`, `Listar`, `Atualizar`, `ListarAtivas`, `AlterarStatus` |
| `ContaFixaRepository` | `Criar`, `BuscarPorID`, `Listar`, `Atualizar`, `ListarAtivas`, `AlterarAtivo` |
| `DespesaGeralRepository` | `Criar`, `BuscarPorID`, `Listar`, `ListarPorMes`, `Atualizar`, `Excluir` |
| `RendaFixaRepository` | `Criar`, `BuscarPorID`, `Listar`, `ListarAtivas`, `ListarVigentesPorMes`, `Atualizar`, `Inativar` |
| `RendaVariavelRepository` | `Criar`, `BuscarPorID`, `Listar`, `ListarPorMes`, `Atualizar`, `Excluir` |
| `RendaExtraRepository` | `Criar`, `BuscarPorID`, `Listar`, `ListarPorMes`, `Atualizar`, `Excluir` |
| `RendimentoInvestimentoRepository` | `Criar`, `BuscarPorID`, `Listar`, `ListarPorMes`, `Atualizar`, `Excluir` |
| `DashboardRepository` | `DespesasPorCategoria` |
| Interfaces `*ForDashboard` | `ListarVigentesPorMes`, `ListarPorFaturaGlobal`, `ListarAtivas` (x2), `ListarPorMes` (x3) |

#### 2.5.2 Exemplo completo (MembroRepository)

**Interface — depois:**
```go
type MembroRepository interface {
    Criar(familiaID string, membro *domain.Membro) (*domain.Membro, error)
    BuscarPorID(familiaID, id string) (*domain.Membro, error)
    Listar(familiaID string) ([]*domain.Membro, error)
    Atualizar(familiaID string, membro *domain.Membro) (*domain.Membro, error)
    Inativar(familiaID, id string) error
}
```

**MembroService — depois:**
```go
func (s *MembroService) Listar(familiaID string) ([]*domain.Membro, error) {
    return s.repo.Listar(familiaID)
}
```

**MembroRepository.Listar — depois:**
```go
func (r *MembroRepository) Listar(familiaID string) ([]*domain.Membro, error) {
    var rows []membroRow
    err := r.db.Select(&rows, `
        SELECT id, nome, relacionamento, ativo
        FROM membros
        WHERE familia_id = $1 AND deleted_at IS NULL
        ORDER BY nome ASC
    `, familiaID)
    // ...
}
```

**MembroRepository.Criar — depois:**
```go
func (r *MembroRepository) Criar(familiaID string, membro *domain.Membro) (*domain.Membro, error) {
    var row membroRow
    err := r.db.QueryRowx(`
        INSERT INTO membros (familia_id, nome, relacionamento, ativo)
        VALUES ($1, $2, $3, $4)
        RETURNING id, nome, relacionamento, ativo
    `, familiaID, membro.Nome, membro.Relacionamento, membro.Ativo).StructScan(&row)
    // ...
}
```

**MembroRepository.BuscarPorID — depois:**
```go
func (r *MembroRepository) BuscarPorID(familiaID, id string) (*domain.Membro, error) {
    var row membroRow
    err := r.db.QueryRowx(`
        SELECT id, nome, relacionamento, ativo
        FROM membros
        WHERE id = $1 AND familia_id = $2 AND deleted_at IS NULL
    `, id, familiaID).StructScan(&row)
    // ...
}
```

> Adicionar `familia_id = $2` no `BuscarPorID` garante que um usuário de outra família não acessa o registro mesmo que adivinhe o UUID.

#### 2.5.3 Casos especiais

**`categorias` — Excluir:** verificação de vínculos também filtra por `familia_id`:
```sql
SELECT COUNT(*) FROM (
    SELECT id FROM despesas_cartao
        WHERE categoria_id = $1 AND familia_id = $2 AND deleted_at IS NULL
    UNION ALL
    SELECT id FROM assinaturas
        WHERE categoria_id = $1 AND familia_id = $2 AND deleted_at IS NULL
    UNION ALL
    SELECT id FROM contas_fixas
        WHERE categoria_id = $1 AND familia_id = $2 AND deleted_at IS NULL
    UNION ALL
    SELECT id FROM despesas_gerais
        WHERE categoria_id = $1 AND familia_id = $2 AND deleted_at IS NULL
) AS vinculos
```

**`dashboard_repository.go` — DespesasPorCategoria:** assinatura muda para `DespesasPorCategoria(familiaID string, mes, ano int)`.

### 2.6 Refatoração dos Handlers

#### Lista completa de handlers

| Handler | Métodos |
|---|---|
| `membro_handler.go` | `Listar`, `Criar`, `BuscarPorID`, `Atualizar`, `Inativar` |
| `categoria_handler.go` | `Listar`, `Criar`, `Atualizar`, `Excluir` |
| `cartao_credito_handler.go` | `Listar`, `Criar`, `BuscarPorID`, `Atualizar`, `Inativar` |
| `despesa_cartao_handler.go` | `ListarPorCartao`, `ListarPorFatura`, `Criar`, `Excluir` |
| `assinatura_handler.go` | `Listar`, `Criar`, `BuscarPorID`, `Atualizar`, `AlterarStatus` |
| `conta_fixa_handler.go` | `Listar`, `Criar`, `BuscarPorID`, `Atualizar`, `AlterarAtivo` |
| `despesa_geral_handler.go` | `Listar`, `Criar`, `BuscarPorID`, `Atualizar`, `Excluir` |
| `renda_fixa_handler.go` | `Listar`, `Criar`, `BuscarPorID`, `Atualizar`, `Inativar`, `ListarVigentesPorMes` |
| `renda_variavel_handler.go` | `Listar`, `Criar`, `BuscarPorID`, `Atualizar`, `Excluir` |
| `renda_extra_handler.go` | `Listar`, `Criar`, `BuscarPorID`, `Atualizar`, `Excluir` |
| `rendimento_investimento_handler.go` | `Listar`, `Criar`, `BuscarPorID`, `Atualizar`, `Excluir` |
| `dashboard_handler.go` | `ResumoMensal`, `DespesasPorCategoria`, `EvolucaoMensal`, `Projecao` |
| `renda_historico_handler.go` | `Historico` |

### 2.7 Ajuste do DashboardService

```go
// Antes:
func (s *DashboardService) ResumoMensal(mes, ano int) (*ResumoMensal, error)

// Depois:
func (s *DashboardService) ResumoMensal(familiaID string, mes, ano int) (*ResumoMensal, error)
```

O mesmo se aplica a `EvolucaoMensal`, `ProjecaoProximosMeses` e `DespesasPorCategoria`.

### 2.8 Ajuste do Seed

**Novo repositório: `familia_repository.go`** com método `BuscarFamiliaPorUsuario(usuarioID string) (string, error)` — usado pelo AuthService no login.

**Refatorar `SeedCategorias`:**
```go
func SeedCategorias(db *sqlx.DB, familiaID string) error {
    categorias := []string{
        "Alimentação", "Transporte", "Lazer", "Saúde",
        "Educação", "Moradia", "Vestuário", "Outros",
    }
    for _, nome := range categorias {
        _, err := db.Exec(`
            INSERT INTO categorias (nome, familia_id)
            VALUES ($1, $2)
            ON CONFLICT (familia_id, nome) DO NOTHING
        `, nome, familiaID)
        if err != nil {
            return fmt.Errorf("seed: erro ao inserir categoria %s: %w", nome, err)
        }
    }
    return nil
}
```

**Refatorar `SeedUsuarios`:** após criar cada usuário, criar sua família padrão e chamar `SeedCategorias`:
```go
// Criar família para o usuário
var familiaID string
err = db.QueryRow(`
    INSERT INTO familias (nome, owner_id) VALUES ($1, $2) RETURNING id
`, "Família Principal", usuarioID).Scan(&familiaID)
if err != nil {
    return fmt.Errorf("seed: erro ao criar família para %s: %w", u.Email, err)
}

// Vincular usuário à família como admin
_, err = db.Exec(`
    INSERT INTO familia_usuarios (familia_id, usuario_id, role)
    VALUES ($1, $2, 'admin')
    ON CONFLICT DO NOTHING
`, familiaID, usuarioID)

// Seed de categorias da família
if err := SeedCategorias(db, familiaID); err != nil {
    return err
}
```

**Importante:** os dois usuários de seed de dev devem ser vinculados à **mesma família** para testar a colaboração entre membros. Ajustar o seed para que o segundo usuário entre como `membro` da família do primeiro.

### 2.9 Teste de Isolamento

```go
//go:build integration
// +build integration

package repository_test

// TestIsolamentoFamilias verifica que a família A não vê dados da família B.
func TestIsolamentoFamilias(t *testing.T) {
    db := conectarBancoDeTeste(t)
    defer limparBanco(t, db)

    // Criar duas famílias independentes
    fidA := criarFamilia(t, db, "Família A")
    fidB := criarFamilia(t, db, "Família B")

    repo := repository.NewMembroRepository(db)

    // Família A cria um membro
    _, err := repo.Criar(fidA, &domain.Membro{Nome: "Ana"})
    require.NoError(t, err)

    // Família B cria um membro
    _, err = repo.Criar(fidB, &domain.Membro{Nome: "Bruno"})
    require.NoError(t, err)

    // Família A só vê seus dados
    membrosA, err := repo.Listar(fidA)
    require.NoError(t, err)
    assert.Len(t, membrosA, 1)
    assert.Equal(t, "Ana", membrosA[0].Nome)

    // Família B só vê seus dados
    membrosB, err := repo.Listar(fidB)
    require.NoError(t, err)
    assert.Len(t, membrosB, 1)
    assert.Equal(t, "Bruno", membrosB[0].Nome)
}

// TestIsolamentoBuscarPorID verifica que família A não acessa registro da família B por ID.
func TestIsolamentoBuscarPorID(t *testing.T) {
    db := conectarBancoDeTeste(t)
    defer limparBanco(t, db)

    fidA := criarFamilia(t, db, "Família A")
    fidB := criarFamilia(t, db, "Família B")

    repo := repository.NewMembroRepository(db)

    membroB, err := repo.Criar(fidB, &domain.Membro{Nome: "Bruno"})
    require.NoError(t, err)

    // Família A tenta acessar membro da família B
    _, err = repo.BuscarPorID(fidA, membroB.ID)
    assert.ErrorIs(t, err, domain.ErrMembroNaoEncontrado,
        "família A não deve enxergar membro da família B")
}
```

### 2.10 Ajuste dos Testes Unitários

**Injetar `familiaID` no contexto nos testes de handler:**
```go
func setupMembroRouter(svc handler.MembroServiceInterface) *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.Use(func(c *gin.Context) {
        c.Set("userID", "usuario-teste-uuid")
        c.Set("familiaID", "familia-teste-uuid") // NOVO
        c.Next()
    })
    h := handler.NewMembroHandler(svc)
    // ... registrar rotas
}
```

---

## 3. TT-08: CORS Restrito e Rate Limiting no Login

### 3.1 CORS via Variável de Ambiente

**Problema atual:** `AllowAllOrigins: true` em `main.go` — qualquer origem pode fazer requisições.

#### Passo 1: Adicionar `AllowedOrigins` ao `Config`

```go
// backend/config/config.go
type ServerConfig struct {
    Port           string   // SERVER_PORT, default "8080"
    AllowedOrigins []string // ALLOWED_ORIGINS, separado por vírgula
}
```

Na função `Load()`:
```go
originsRaw := v.GetString("ALLOWED_ORIGINS")
var origins []string
if originsRaw != "" {
    for _, o := range strings.Split(originsRaw, ",") {
        if trimmed := strings.TrimSpace(o); trimmed != "" {
            origins = append(origins, trimmed)
        }
    }
}
if len(origins) == 0 {
    origins = []string{"http://localhost:5173"} // fallback dev
}
```

#### Passo 2: Substituir CORS em `main.go`

```go
r.Use(cors.New(cors.Config{
    AllowOrigins:     cfg.Server.AllowedOrigins,
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: false,
}))
```

#### Passo 3: Variáveis de ambiente

**Railway:**
```
ALLOWED_ORIGINS=https://financas-familiares.vercel.app
```

**`.env` local:**
```
ALLOWED_ORIGINS=http://localhost:5173
```

### 3.2 Rate Limiting no Login

**Implementação:** `golang.org/x/time/rate` (dependência transitiva — sem pacote novo).

**Novo arquivo: `backend/internal/middleware/ratelimit.go`**

```go
package middleware

import (
    "net/http"
    "sync"

    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

var loginLimiters sync.Map

func getLoginLimiter(ip string) *rate.Limiter {
    val, _ := loginLimiters.LoadOrStore(ip, rate.NewLimiter(rate.Limit(10.0/60.0), 10))
    return val.(*rate.Limiter)
}

// LoginRateLimiter limita a 10 req/min por IP em POST /auth/login.
func LoginRateLimiter() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !getLoginLimiter(c.ClientIP()).Allow() {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "erro": "muitas tentativas, aguarde 1 minuto",
            })
            return
        }
        c.Next()
    }
}
```

**Integração em `main.go`:**
```go
authGroup.POST("/login", middleware.LoginRateLimiter(), authHandler.Login)
```

**Como testar o 429:**
```bash
for i in $(seq 1 11); do
  curl -s -o /dev/null -w "%{http_code}\n" \
    -X POST https://seu-backend.railway.app/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"teste@teste.com","senha":"errada"}'
done
# Primeiras 10 → 401; 11ª → 429
```

---

## 4. TT-09: Índices de Performance

### 4.1 Estratégia

Não usar `CONCURRENTLY` nas migrations (golang-migrate executa em transação). Para banco de volume baixo (sistema familiar), o bloqueio é de milissegundos — aceitável.

### 4.2 Migration 000018

**`backend/migrations/000018_add_performance_indexes.up.sql`**

```sql
-- Sprint 9 / TT-09: Índices de performance

-- =====================================================================
-- PARTIAL INDEXES em deleted_at (WHERE deleted_at IS NULL)
-- =====================================================================
CREATE INDEX IF NOT EXISTS idx_familias_deleted_at
    ON familias (deleted_at) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_membros_deleted_at
    ON membros (deleted_at) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_categorias_deleted_at
    ON categorias (deleted_at) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_cartoes_credito_deleted_at
    ON cartoes_credito (deleted_at) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_despesas_cartao_deleted_at
    ON despesas_cartao (deleted_at) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_assinaturas_deleted_at
    ON assinaturas (deleted_at) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_contas_fixas_deleted_at
    ON contas_fixas (deleted_at) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_despesas_gerais_deleted_at
    ON despesas_gerais (deleted_at) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_rendas_fixas_deleted_at
    ON rendas_fixas (deleted_at) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_rendas_variaveis_deleted_at
    ON rendas_variaveis (deleted_at) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_rendas_extras_deleted_at
    ON rendas_extras (deleted_at) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_rendimentos_investimento_deleted_at
    ON rendimentos_investimento (deleted_at) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_deleted_at
    ON refresh_tokens (deleted_at) WHERE deleted_at IS NULL;

-- =====================================================================
-- ÍNDICES em familia_id (depende do TT-07 — migration 000017)
-- =====================================================================
CREATE INDEX IF NOT EXISTS idx_membros_familia_id
    ON membros (familia_id);

CREATE INDEX IF NOT EXISTS idx_categorias_familia_id
    ON categorias (familia_id);

CREATE INDEX IF NOT EXISTS idx_cartoes_credito_familia_id
    ON cartoes_credito (familia_id);

CREATE INDEX IF NOT EXISTS idx_despesas_cartao_familia_id
    ON despesas_cartao (familia_id);

CREATE INDEX IF NOT EXISTS idx_assinaturas_familia_id
    ON assinaturas (familia_id);

CREATE INDEX IF NOT EXISTS idx_contas_fixas_familia_id
    ON contas_fixas (familia_id);

CREATE INDEX IF NOT EXISTS idx_despesas_gerais_familia_id
    ON despesas_gerais (familia_id);

CREATE INDEX IF NOT EXISTS idx_rendas_fixas_familia_id
    ON rendas_fixas (familia_id);

CREATE INDEX IF NOT EXISTS idx_rendas_variaveis_familia_id
    ON rendas_variaveis (familia_id);

CREATE INDEX IF NOT EXISTS idx_rendas_extras_familia_id
    ON rendas_extras (familia_id);

CREATE INDEX IF NOT EXISTS idx_rendimentos_investimento_familia_id
    ON rendimentos_investimento (familia_id);

-- =====================================================================
-- ÍNDICE COMPOSTO em despesas_cartao: (familia_id, fatura_ano, fatura_mes)
-- =====================================================================
CREATE INDEX IF NOT EXISTS idx_despesas_cartao_familia_fatura
    ON despesas_cartao (familia_id, fatura_ano, fatura_mes);

-- =====================================================================
-- ÍNDICES em familia_usuarios
-- =====================================================================
CREATE INDEX IF NOT EXISTS idx_familia_usuarios_usuario_id
    ON familia_usuarios (usuario_id);

-- =====================================================================
-- ÍNDICES em refresh_tokens (usuario_id mantido — tokens são por usuário)
-- =====================================================================
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_usuario_id
    ON refresh_tokens (usuario_id);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash_ativo
    ON refresh_tokens (token_hash)
    WHERE revogado = FALSE AND deleted_at IS NULL;
```

**`backend/migrations/000018_add_performance_indexes.down.sql`**

```sql
DROP INDEX IF EXISTS idx_refresh_tokens_token_hash_ativo;
DROP INDEX IF EXISTS idx_refresh_tokens_usuario_id;
DROP INDEX IF EXISTS idx_familia_usuarios_usuario_id;
DROP INDEX IF EXISTS idx_despesas_cartao_familia_fatura;
DROP INDEX IF EXISTS idx_rendimentos_investimento_familia_id;
DROP INDEX IF EXISTS idx_rendas_extras_familia_id;
DROP INDEX IF EXISTS idx_rendas_variaveis_familia_id;
DROP INDEX IF EXISTS idx_rendas_fixas_familia_id;
DROP INDEX IF EXISTS idx_despesas_gerais_familia_id;
DROP INDEX IF EXISTS idx_contas_fixas_familia_id;
DROP INDEX IF EXISTS idx_assinaturas_familia_id;
DROP INDEX IF EXISTS idx_despesas_cartao_familia_id;
DROP INDEX IF EXISTS idx_cartoes_credito_familia_id;
DROP INDEX IF EXISTS idx_categorias_familia_id;
DROP INDEX IF EXISTS idx_membros_familia_id;
DROP INDEX IF EXISTS idx_refresh_tokens_deleted_at;
DROP INDEX IF EXISTS idx_rendimentos_investimento_deleted_at;
DROP INDEX IF EXISTS idx_rendas_extras_deleted_at;
DROP INDEX IF EXISTS idx_rendas_variaveis_deleted_at;
DROP INDEX IF EXISTS idx_rendas_fixas_deleted_at;
DROP INDEX IF EXISTS idx_despesas_gerais_deleted_at;
DROP INDEX IF EXISTS idx_contas_fixas_deleted_at;
DROP INDEX IF EXISTS idx_assinaturas_deleted_at;
DROP INDEX IF EXISTS idx_despesas_cartao_deleted_at;
DROP INDEX IF EXISTS idx_cartoes_credito_deleted_at;
DROP INDEX IF EXISTS idx_categorias_deleted_at;
DROP INDEX IF EXISTS idx_membros_deleted_at;
DROP INDEX IF EXISTS idx_familias_deleted_at;
```

> **Nota:** `refresh_tokens.usuario_id` mantém índice por `usuario_id` (não `familia_id`) pois tokens são por usuário individual, não por família.

---

## 5. Critério de Conclusão

### TT-07 — Isolamento de dados por família

- [ ] Migration `000017` aplicada sem erros em banco limpo e em banco com dados existentes
- [ ] Migration `000017` down funciona corretamente (apaga `familias`, `familia_usuarios` e remove colunas)
- [ ] Tabelas `familias` e `familia_usuarios` criadas corretamente
- [ ] JWT inclui `familia_id` nas claims
- [ ] Middleware extrai `familia_id` do token e disponibiliza como `"familiaID"` no contexto Gin
- [ ] `FamiliaRepository.BuscarFamiliaPorUsuario` implementado e usado pelo AuthService
- [ ] Todos os repositórios (11 + dashboard) filtram por `familia_id`
- [ ] Todas as interfaces de repositório em `service/` atualizadas
- [ ] Todos os services passam `familiaID` para o repositório
- [ ] Todos os handlers extraem `"familiaID"` do contexto
- [ ] `SeedCategorias` recebe `familiaID` e insere com `familia_id`
- [ ] `SeedUsuarios` cria família e chama `SeedCategorias` para cada usuário de seed
- [ ] Os dois usuários de seed compartilham a mesma família (testam colaboração)
- [ ] Constraint `UNIQUE(familia_id, nome)` em categorias substituiu `UNIQUE(nome)` global
- [ ] Todos os testes unitários compilam e passam (mocks atualizados com `familiaID`)
- [ ] Teste de integração `TestIsolamentoFamilias` passa
- [ ] Teste de integração `TestIsolamentoBuscarPorID` passa
- [ ] `go build ./...` sem erros

### TT-08 — CORS + Rate Limiting

- [ ] `ALLOWED_ORIGINS` lido do ambiente em `config.go`
- [ ] CORS com whitelist no `main.go` (sem `AllowAllOrigins: true`)
- [ ] Requisição de origem não listada recebe resposta sem cabeçalho CORS
- [ ] Frontend em produção (Vercel) funciona normalmente
- [ ] Middleware `LoginRateLimiter` implementado em `middleware/ratelimit.go`
- [ ] Rate limiter aplicado apenas em `POST /api/v1/auth/login`
- [ ] 11ª requisição consecutiva do mesmo IP retorna `429` com `{"erro": "muitas tentativas, aguarde 1 minuto"}`
- [ ] `ALLOWED_ORIGINS` configurado no Railway com o domínio real do Vercel
- [ ] `go build ./...` sem erros

### TT-09 — Índices de performance

- [ ] Migration `000018` aplicada sem erros
- [ ] Migration `000018` down funciona
- [ ] `\d membros` no psql mostra `idx_membros_familia_id` e `idx_membros_deleted_at`
- [ ] `EXPLAIN ANALYZE` na query do dashboard usa índice composto `idx_despesas_cartao_familia_fatura`
- [ ] Índices de `familia_usuarios` e `refresh_tokens` criados
- [ ] `go build ./...` sem erros (TT-09 é apenas SQL)

---

## 6. Arquivos Modificados por Item

### TT-07

| Arquivo | Tipo de Mudança |
|---|---|
| `backend/migrations/000017_add_familia_isolation.up.sql` | **Novo** |
| `backend/migrations/000017_add_familia_isolation.down.sql` | **Novo** |
| `backend/internal/domain/auth.go` | Modificado: adicionar `FamiliaID` nas Claims do JWT |
| `backend/internal/domain/familia.go` | **Novo**: struct `Familia`, `FamiliaUsuario`, erros de domínio |
| `backend/internal/repository/familia_repository.go` | **Novo**: `BuscarFamiliaPorUsuario`, `Criar`, `VincularUsuario` |
| `backend/internal/middleware/auth.go` | Modificado: extrair `familia_id` do token e setar `"familiaID"` no contexto |
| `backend/internal/service/auth_service.go` | Modificado: buscar `familiaID` ao fazer login e incluir no JWT |
| `backend/internal/handler/helpers.go` | **Novo**: funções `getFamiliaID`, `getUsuarioID` |
| `backend/internal/repository/membro_repository.go` | Modificado: `familiaID` em todos os métodos |
| `backend/internal/repository/categoria_repository.go` | Modificado: idem |
| `backend/internal/repository/cartao_credito_repository.go` | Modificado: idem |
| `backend/internal/repository/despesa_cartao_repository.go` | Modificado: idem |
| `backend/internal/repository/assinatura_repository.go` | Modificado: idem |
| `backend/internal/repository/conta_fixa_repository.go` | Modificado: idem |
| `backend/internal/repository/despesa_geral_repository.go` | Modificado: idem |
| `backend/internal/repository/renda_fixa_repository.go` | Modificado: idem |
| `backend/internal/repository/renda_variavel_repository.go` | Modificado: idem |
| `backend/internal/repository/renda_extra_repository.go` | Modificado: idem |
| `backend/internal/repository/rendimento_investimento_repository.go` | Modificado: idem |
| `backend/internal/repository/dashboard_repository.go` | Modificado: `familiaID` em `DespesasPorCategoria` |
| `backend/internal/repository/seed.go` | Modificado: criar família no seed; `SeedCategorias` recebe `familiaID` |
| `backend/internal/service/membro_service.go` | Modificado: interfaces e métodos com `familiaID` |
| `backend/internal/service/categoria_service.go` | Modificado: idem |
| `backend/internal/service/cartao_credito_service.go` | Modificado: idem |
| `backend/internal/service/despesa_cartao_service.go` | Modificado: idem |
| `backend/internal/service/assinatura_service.go` | Modificado: idem |
| `backend/internal/service/conta_fixa_service.go` | Modificado: idem |
| `backend/internal/service/despesa_geral_service.go` | Modificado: idem |
| `backend/internal/service/renda_fixa_service.go` | Modificado: idem |
| `backend/internal/service/renda_variavel_service.go` | Modificado: idem |
| `backend/internal/service/renda_extra_service.go` | Modificado: idem |
| `backend/internal/service/rendimento_investimento_service.go` | Modificado: idem |
| `backend/internal/service/dashboard_service.go` | Modificado: `familiaID` em todos os métodos e interfaces `*ForDashboard` |
| `backend/internal/service/renda_historico_service.go` | Modificado: idem |
| `backend/internal/handler/membro_handler.go` | Modificado: extrair `familiaID`; atualizar interface local |
| `backend/internal/handler/categoria_handler.go` | Modificado: idem |
| `backend/internal/handler/cartao_credito_handler.go` | Modificado: idem |
| `backend/internal/handler/despesa_cartao_handler.go` | Modificado: idem |
| `backend/internal/handler/assinatura_handler.go` | Modificado: idem |
| `backend/internal/handler/conta_fixa_handler.go` | Modificado: idem |
| `backend/internal/handler/despesa_geral_handler.go` | Modificado: idem |
| `backend/internal/handler/renda_fixa_handler.go` | Modificado: idem |
| `backend/internal/handler/renda_variavel_handler.go` | Modificado: idem |
| `backend/internal/handler/renda_extra_handler.go` | Modificado: idem |
| `backend/internal/handler/rendimento_investimento_handler.go` | Modificado: idem |
| `backend/internal/handler/dashboard_handler.go` | Modificado: idem |
| `backend/internal/handler/renda_historico_handler.go` | Modificado: idem |
| `backend/internal/handler/*_handler_test.go` (todos) | Modificado: mocks e `setupRouter` com injeção de `familiaID` |
| `backend/internal/repository/isolamento_test.go` | **Novo**: teste de integração de isolamento por família |
| `backend/cmd/server/main.go` | Modificado: remover chamada separada a `SeedCategorias` |

### TT-08

| Arquivo | Tipo de Mudança |
|---|---|
| `backend/config/config.go` | Modificado: `AllowedOrigins []string` em `ServerConfig` |
| `backend/internal/middleware/ratelimit.go` | **Novo** |
| `backend/internal/middleware/ratelimit_test.go` | **Novo** |
| `backend/cmd/server/main.go` | Modificado: CORS com whitelist; `LoginRateLimiter` na rota de login |

### TT-09

| Arquivo | Tipo de Mudança |
|---|---|
| `backend/migrations/000018_add_performance_indexes.up.sql` | **Novo** |
| `backend/migrations/000018_add_performance_indexes.down.sql` | **Novo** |

---

### Arquivos Críticos para Implementação

- `backend/cmd/server/main.go` — composição central: CORS, rotas, seeds; 3 mudanças
- `backend/internal/middleware/auth.go` — extrai `familia_id` do JWT e seta `"familiaID"` no contexto; todos os handlers dependem disso
- `backend/internal/service/auth_service.go` — busca `familiaID` no login e inclui no token; ponto de entrada do isolamento
- `backend/internal/repository/seed.go` — criar família no seed e chamar `SeedCategorias` por família
- `backend/internal/service/dashboard_service.go` — arquivo mais complexo: múltiplas interfaces `*ForDashboard` e propagação de `familiaID`
