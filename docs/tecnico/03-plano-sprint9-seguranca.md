# Plano de Execução — Sprint 9: Segurança e Isolamento de Dados

## Índice
1. [Visão Geral e Sequência de Execução](#1-visão-geral-e-sequência-de-execução)
2. [TT-07: Isolamento de Dados por Usuário](#2-tt-07-isolamento-de-dados-por-usuário)
3. [TT-08: CORS Restrito e Rate Limiting no Login](#3-tt-08-cors-restrito-e-rate-limiting-no-login)
4. [TT-09: Índices de Performance](#4-tt-09-índices-de-performance)
5. [Critério de Conclusão](#5-critério-de-conclusão)
6. [Arquivos Modificados por Item](#6-arquivos-modificados-por-item)

---

## 1. Visão Geral e Sequência de Execução

### Por que TT-07 → TT-08 → TT-09?

**TT-07 primeiro** porque é o problema de segurança mais crítico: sem isolamento de dados, qualquer usuário autenticado pode ler e modificar dados de outro usuário. Além disso, o TT-09 depende diretamente do TT-07: os índices em `usuario_id` só fazem sentido depois que a coluna existir no banco e nas queries. Executar TT-09 antes tornaria parte do trabalho redundante.

**TT-08 segundo** porque é independente do banco de dados (opera na camada de configuração do Gin e de middleware HTTP). Pode ser feito em paralelo com TT-07 em teoria, mas sequencialmente é mais simples e seguro. Não gera retrabalho.

**TT-09 por último** porque depende de TT-07 para os índices em `usuario_id` e complementa o trabalho já feito, sem risco de conflito com as migrations anteriores.

---

## 2. TT-07: Isolamento de Dados por Usuário

### 2.1 Análise do Estado Atual

Ao ler o código, identificou-se:

- O middleware em `backend/internal/middleware/auth.go` define o `userID` no contexto Gin com a chave `"userID"` (linha 41: `c.Set("userID", usuarioID)`).
- Nenhum handler extrai `userID` do contexto — a chave existe mas nunca é usada.
- Todos os repositórios operam sem filtro de `usuario_id`.
- A tabela `refresh_tokens` já possui `usuario_id UUID NOT NULL REFERENCES usuarios(id)`.
- O `SeedCategorias` insere categorias sem `usuario_id`, o que não funcionará após adicionar `NOT NULL`.

### 2.2 Estratégia de Migration

**Decisão: uma migration única (`000017`) com todos os `ALTER TABLE`.**

Justificativas:
- Menos arquivos para gerenciar.
- As alterações são atômicas: ou todas aplicam, ou nenhuma.
- O rollback (`down`) desfaz tudo de uma vez.
- Não há risco de estado intermediário inconsistente entre tabelas.

A numeração `000017` é a próxima disponível (a última é `000016`). O enunciado mencionava reservar `000017` para `lancamentos_despesas` em sprint futura — neste plano, `000017` será usada para `usuario_id`, e `lancamentos_despesas` usará `000018` quando chegar.

### 2.3 Migration 000017 — Adicionar `usuario_id` em Todas as Tabelas

#### Arquivo: `backend/migrations/000017_add_usuario_id_to_all_tables.up.sql`

```sql
-- Sprint 9 / TT-07: Adicionar usuario_id em todas as tabelas de dados
-- Estratégia para banco de desenvolvimento com dados existentes:
--   1. Adicionar coluna como nullable
--   2. Popular com o primeiro usuário (seed de dev)
--   3. Adicionar constraint NOT NULL e FK
-- Em banco limpo (produção), as etapas 1 e 2 são transparentes: nenhum registro existe.

-- =====================================================================
-- MEMBROS
-- =====================================================================
ALTER TABLE membros
    ADD COLUMN IF NOT EXISTS usuario_id UUID;

UPDATE membros
SET usuario_id = (SELECT id FROM usuarios ORDER BY created_at ASC LIMIT 1)
WHERE usuario_id IS NULL;

ALTER TABLE membros
    ALTER COLUMN usuario_id SET NOT NULL,
    ADD CONSTRAINT membros_usuario_id_fkey
        FOREIGN KEY (usuario_id) REFERENCES usuarios(id);

-- =====================================================================
-- CATEGORIAS
-- =====================================================================
ALTER TABLE categorias
    ADD COLUMN IF NOT EXISTS usuario_id UUID;

UPDATE categorias
SET usuario_id = (SELECT id FROM usuarios ORDER BY created_at ASC LIMIT 1)
WHERE usuario_id IS NULL;

ALTER TABLE categorias
    ALTER COLUMN usuario_id SET NOT NULL,
    ADD CONSTRAINT categorias_usuario_id_fkey
        FOREIGN KEY (usuario_id) REFERENCES usuarios(id);

-- =====================================================================
-- CARTOES_CREDITO
-- =====================================================================
ALTER TABLE cartoes_credito
    ADD COLUMN IF NOT EXISTS usuario_id UUID;

UPDATE cartoes_credito
SET usuario_id = (SELECT id FROM usuarios ORDER BY created_at ASC LIMIT 1)
WHERE usuario_id IS NULL;

ALTER TABLE cartoes_credito
    ALTER COLUMN usuario_id SET NOT NULL,
    ADD CONSTRAINT cartoes_credito_usuario_id_fkey
        FOREIGN KEY (usuario_id) REFERENCES usuarios(id);

-- =====================================================================
-- DESPESAS_CARTAO
-- =====================================================================
ALTER TABLE despesas_cartao
    ADD COLUMN IF NOT EXISTS usuario_id UUID;

UPDATE despesas_cartao
SET usuario_id = (SELECT id FROM usuarios ORDER BY created_at ASC LIMIT 1)
WHERE usuario_id IS NULL;

ALTER TABLE despesas_cartao
    ALTER COLUMN usuario_id SET NOT NULL,
    ADD CONSTRAINT despesas_cartao_usuario_id_fkey
        FOREIGN KEY (usuario_id) REFERENCES usuarios(id);

-- =====================================================================
-- ASSINATURAS
-- =====================================================================
ALTER TABLE assinaturas
    ADD COLUMN IF NOT EXISTS usuario_id UUID;

UPDATE assinaturas
SET usuario_id = (SELECT id FROM usuarios ORDER BY created_at ASC LIMIT 1)
WHERE usuario_id IS NULL;

ALTER TABLE assinaturas
    ALTER COLUMN usuario_id SET NOT NULL,
    ADD CONSTRAINT assinaturas_usuario_id_fkey
        FOREIGN KEY (usuario_id) REFERENCES usuarios(id);

-- =====================================================================
-- CONTAS_FIXAS
-- =====================================================================
ALTER TABLE contas_fixas
    ADD COLUMN IF NOT EXISTS usuario_id UUID;

UPDATE contas_fixas
SET usuario_id = (SELECT id FROM usuarios ORDER BY created_at ASC LIMIT 1)
WHERE usuario_id IS NULL;

ALTER TABLE contas_fixas
    ALTER COLUMN usuario_id SET NOT NULL,
    ADD CONSTRAINT contas_fixas_usuario_id_fkey
        FOREIGN KEY (usuario_id) REFERENCES usuarios(id);

-- =====================================================================
-- DESPESAS_GERAIS
-- =====================================================================
ALTER TABLE despesas_gerais
    ADD COLUMN IF NOT EXISTS usuario_id UUID;

UPDATE despesas_gerais
SET usuario_id = (SELECT id FROM usuarios ORDER BY created_at ASC LIMIT 1)
WHERE usuario_id IS NULL;

ALTER TABLE despesas_gerais
    ALTER COLUMN usuario_id SET NOT NULL,
    ADD CONSTRAINT despesas_gerais_usuario_id_fkey
        FOREIGN KEY (usuario_id) REFERENCES usuarios(id);

-- =====================================================================
-- RENDAS_FIXAS
-- =====================================================================
ALTER TABLE rendas_fixas
    ADD COLUMN IF NOT EXISTS usuario_id UUID;

UPDATE rendas_fixas
SET usuario_id = (SELECT id FROM usuarios ORDER BY created_at ASC LIMIT 1)
WHERE usuario_id IS NULL;

ALTER TABLE rendas_fixas
    ALTER COLUMN usuario_id SET NOT NULL,
    ADD CONSTRAINT rendas_fixas_usuario_id_fkey
        FOREIGN KEY (usuario_id) REFERENCES usuarios(id);

-- =====================================================================
-- RENDAS_VARIAVEIS
-- =====================================================================
ALTER TABLE rendas_variaveis
    ADD COLUMN IF NOT EXISTS usuario_id UUID;

UPDATE rendas_variaveis
SET usuario_id = (SELECT id FROM usuarios ORDER BY created_at ASC LIMIT 1)
WHERE usuario_id IS NULL;

ALTER TABLE rendas_variaveis
    ALTER COLUMN usuario_id SET NOT NULL,
    ADD CONSTRAINT rendas_variaveis_usuario_id_fkey
        FOREIGN KEY (usuario_id) REFERENCES usuarios(id);

-- =====================================================================
-- RENDAS_EXTRAS
-- =====================================================================
ALTER TABLE rendas_extras
    ADD COLUMN IF NOT EXISTS usuario_id UUID;

UPDATE rendas_extras
SET usuario_id = (SELECT id FROM usuarios ORDER BY created_at ASC LIMIT 1)
WHERE usuario_id IS NULL;

ALTER TABLE rendas_extras
    ALTER COLUMN usuario_id SET NOT NULL,
    ADD CONSTRAINT rendas_extras_usuario_id_fkey
        FOREIGN KEY (usuario_id) REFERENCES usuarios(id);

-- =====================================================================
-- RENDIMENTOS_INVESTIMENTO
-- =====================================================================
ALTER TABLE rendimentos_investimento
    ADD COLUMN IF NOT EXISTS usuario_id UUID;

UPDATE rendimentos_investimento
SET usuario_id = (SELECT id FROM usuarios ORDER BY created_at ASC LIMIT 1)
WHERE usuario_id IS NULL;

ALTER TABLE rendimentos_investimento
    ALTER COLUMN usuario_id SET NOT NULL,
    ADD CONSTRAINT rendimentos_investimento_usuario_id_fkey
        FOREIGN KEY (usuario_id) REFERENCES usuarios(id);
```

> **Nota sobre banco de produção (Railway):** Se o banco de produção não tiver dados reais de usuários ainda (apenas o seed de dev foi executado), a lógica do `UPDATE … WHERE usuario_id IS NULL` é segura. Se o banco de produção não tiver nenhum usuário criado antes desta migration, o `UPDATE` não afetará nenhuma linha e o `ALTER COLUMN … SET NOT NULL` aplicará sem problema (zero linhas = nenhuma viola a constraint).

#### Arquivo: `backend/migrations/000017_add_usuario_id_to_all_tables.down.sql`

```sql
-- Rollback: remover usuario_id e constraints de todas as tabelas

ALTER TABLE rendimentos_investimento
    DROP CONSTRAINT IF EXISTS rendimentos_investimento_usuario_id_fkey,
    DROP COLUMN IF EXISTS usuario_id;

ALTER TABLE rendas_extras
    DROP CONSTRAINT IF EXISTS rendas_extras_usuario_id_fkey,
    DROP COLUMN IF EXISTS usuario_id;

ALTER TABLE rendas_variaveis
    DROP CONSTRAINT IF EXISTS rendas_variaveis_usuario_id_fkey,
    DROP COLUMN IF EXISTS usuario_id;

ALTER TABLE rendas_fixas
    DROP CONSTRAINT IF EXISTS rendas_fixas_usuario_id_fkey,
    DROP COLUMN IF EXISTS usuario_id;

ALTER TABLE despesas_gerais
    DROP CONSTRAINT IF EXISTS despesas_gerais_usuario_id_fkey,
    DROP COLUMN IF EXISTS usuario_id;

ALTER TABLE contas_fixas
    DROP CONSTRAINT IF EXISTS contas_fixas_usuario_id_fkey,
    DROP COLUMN IF EXISTS usuario_id;

ALTER TABLE assinaturas
    DROP CONSTRAINT IF EXISTS assinaturas_usuario_id_fkey,
    DROP COLUMN IF EXISTS usuario_id;

ALTER TABLE despesas_cartao
    DROP CONSTRAINT IF EXISTS despesas_cartao_usuario_id_fkey,
    DROP COLUMN IF EXISTS usuario_id;

ALTER TABLE cartoes_credito
    DROP CONSTRAINT IF EXISTS cartoes_credito_usuario_id_fkey,
    DROP COLUMN IF EXISTS usuario_id;

ALTER TABLE categorias
    DROP CONSTRAINT IF EXISTS categorias_usuario_id_fkey,
    DROP COLUMN IF EXISTS usuario_id;

ALTER TABLE membros
    DROP CONSTRAINT IF EXISTS membros_usuario_id_fkey,
    DROP COLUMN IF EXISTS usuario_id;
```

### 2.4 Ajuste do Middleware: Chave Correta no Contexto

O middleware atual usa a chave `"userID"` (linha 41 de `auth.go`). A convenção escolhida para o projeto será `"userID"` (sem underscore), mas é preciso padronizar. O middleware já está correto; o que falta é os handlers lerem essa chave.

Convenção adotada neste plano: **manter a chave `"userID"`** (como está no middleware) e usar `c.MustGet("userID").(string)` nos handlers.

> Alternativa `c.GetString("userID")` também funciona, mas retorna string vazia sem indicar ausência — `MustGet` com cast é mais explícito e panica se o middleware não tiver rodado (o que nunca deve ocorrer em rotas protegidas).

### 2.5 Refatoração dos Repositórios

Para cada repositório, o padrão de refatoração segue a mesma lógica. Abaixo o mapeamento completo de **todos os métodos** que precisam receber `usuarioID string` como parâmetro:

#### 2.5.1 Interfaces a atualizar (em `service/*.go`)

Cada arquivo de service declara uma interface de repositório. Essas interfaces precisam ser atualizadas simultaneamente com as implementações, ou o compilador Go rejeitará a build.

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
| Interfaces `*ForDashboard` em `dashboard_service.go` | `ListarVigentesPorMes`, `ListarPorFaturaGlobal`, `ListarPorMes` (x3), `ListarAtivas` (x2), `ListarPorMes` (despesa geral) |

**Atenção especial ao `DashboardService`:** ele recebe múltiplos repositórios como interfaces independentes (`RendaFixaRepositoryForDashboard`, `DespesaCartaoRepositoryForDashboard`, etc.). Todas essas interfaces mini também precisam ter `usuarioID` nos métodos que fazem query.

#### 2.5.2 Exemplo completo do padrão de refatoração (MembroRepository)

**Interface em `membro_service.go` — antes:**
```go
type MembroRepository interface {
    Criar(membro *domain.Membro) (*domain.Membro, error)
    BuscarPorID(id string) (*domain.Membro, error)
    Listar() ([]*domain.Membro, error)
    Atualizar(membro *domain.Membro) (*domain.Membro, error)
    Inativar(id string) error
}
```

**Interface em `membro_service.go` — depois:**
```go
type MembroRepository interface {
    Criar(usuarioID string, membro *domain.Membro) (*domain.Membro, error)
    BuscarPorID(usuarioID, id string) (*domain.Membro, error)
    Listar(usuarioID string) ([]*domain.Membro, error)
    Atualizar(usuarioID string, membro *domain.Membro) (*domain.Membro, error)
    Inativar(usuarioID, id string) error
}
```

**MembroService — antes:**
```go
func (s *MembroService) Listar() ([]*domain.Membro, error) {
    return s.repo.Listar()
}
```

**MembroService — depois:**
```go
func (s *MembroService) Listar(usuarioID string) ([]*domain.Membro, error) {
    return s.repo.Listar(usuarioID)
}
```

**Interface pública em `membro_service.go` — depois:**
```go
type MembroServiceInterface interface {
    Criar(usuarioID, nome, relacionamento string) (*domain.Membro, error)
    BuscarPorID(usuarioID, id string) (*domain.Membro, error)
    Listar(usuarioID string) ([]*domain.Membro, error)
    Atualizar(usuarioID, id, nome, relacionamento string, ativo bool) (*domain.Membro, error)
    Inativar(usuarioID, id string) error
}
```

**MembroRepository (implementação) — método Listar — depois:**
```go
func (r *MembroRepository) Listar(usuarioID string) ([]*domain.Membro, error) {
    var rows []membroRow
    err := r.db.Select(&rows, `
        SELECT id, nome, relacionamento, ativo
        FROM membros
        WHERE usuario_id = $1 AND deleted_at IS NULL
        ORDER BY nome ASC
    `, usuarioID)
    // ...
}
```

**MembroRepository — método Criar — depois:**
```go
func (r *MembroRepository) Criar(usuarioID string, membro *domain.Membro) (*domain.Membro, error) {
    var row membroRow
    err := r.db.QueryRowx(`
        INSERT INTO membros (usuario_id, nome, relacionamento, ativo)
        VALUES ($1, $2, $3, $4)
        RETURNING id, nome, relacionamento, ativo
    `, usuarioID, membro.Nome, membro.Relacionamento, membro.Ativo).StructScan(&row)
    // ...
}
```

**MembroRepository — método BuscarPorID — depois:**
```go
func (r *MembroRepository) BuscarPorID(usuarioID, id string) (*domain.Membro, error) {
    var row membroRow
    err := r.db.QueryRowx(`
        SELECT id, nome, relacionamento, ativo
        FROM membros
        WHERE id = $1 AND usuario_id = $2 AND deleted_at IS NULL
    `, id, usuarioID).StructScan(&row)
    // ...
}
```

> Adicionar `usuario_id = $2` no `BuscarPorID` garante que um usuário não consiga acessar um recurso de outro usuário mesmo que adivinhe o UUID.

#### 2.5.3 Casos especiais

**`despesas_cartao` — ListarPorCartao e ListarPorFatura:**
O `cartao_id` é de propriedade do usuário, portanto o filtro pode ser feito indiretamente via JOIN, ou diretamente adicionando `usuario_id` na tabela. O plano adiciona `usuario_id` diretamente na tabela `despesas_cartao` (já previsto na migration), pois é mais simples e eficiente:

```sql
WHERE usuario_id = $1 AND cartao_id = $2 AND deleted_at IS NULL
```

**`categorias` — Excluir:** A query atual verifica vínculos em `despesas_cartao`, `assinaturas`, `contas_fixas`, `despesas_gerais`. Com isolamento, o filtro de vínculos também deve considerar `usuario_id`:

```sql
SELECT COUNT(*) FROM (
    SELECT id FROM despesas_cartao
        WHERE categoria_id = $1 AND usuario_id = $2 AND deleted_at IS NULL
    UNION ALL
    SELECT id FROM assinaturas
        WHERE categoria_id = $1 AND usuario_id = $2 AND deleted_at IS NULL
    UNION ALL
    SELECT id FROM contas_fixas
        WHERE categoria_id = $1 AND usuario_id = $2 AND deleted_at IS NULL
    UNION ALL
    SELECT id FROM despesas_gerais
        WHERE categoria_id = $1 AND usuario_id = $2 AND deleted_at IS NULL
) AS vinculos
```

**`dashboard_repository.go` — DespesasPorCategoria:** Esta query faz uma union de várias tabelas. Após TT-07, cada sub-query precisa filtrar por `usuario_id`. A assinatura muda para `DespesasPorCategoria(usuarioID string, mes, ano int)`.

**`assinaturas` e `contas_fixas` — ListarAtivas:** Usado pelo `DashboardService`. A assinatura muda para `ListarAtivas(usuarioID string)`.

**`renda_fixa` — ListarVigentesPorMes:** Usado pelo `DashboardService`. A assinatura muda para `ListarVigentesPorMes(usuarioID string, mes, ano int)`.

### 2.6 Refatoração dos Handlers

Cada handler precisa extrair `userID` do contexto e passá-lo para o service. A chave no contexto é `"userID"` (conforme o middleware atual).

#### Padrão de extração nos handlers

```go
// No início de cada handler que requer isolamento:
usuarioID, ok := c.Get("userID")
if !ok {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
    return
}
uid := usuarioID.(string)
```

Alternativamente, pode-se criar uma função auxiliar no pacote `handler`:

```go
// handler/helpers.go (novo arquivo no mesmo pacote)
func getUsuarioID(c *gin.Context) (string, bool) {
    v, ok := c.Get("userID")
    if !ok {
        return "", false
    }
    uid, ok := v.(string)
    return uid, ok
}
```

#### Lista completa de handlers e métodos que precisam extrair `usuarioID`

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

**Interface do handler (`MembroServiceInterface` redeclarada no handler)** também precisa ser atualizada para refletir os novos parâmetros, caso contrário o compilador rejeitará a implementação.

### 2.7 Ajuste do DashboardService

O `DashboardService` chama vários repositórios. Após a refatoração, todos os métodos que fazem query precisam receber `usuarioID`. O service passa esse parâmetro para cada repositório.

O `DashboardService.ResumoMensal` muda de:
```go
func (s *DashboardService) ResumoMensal(mes, ano int) (*ResumoMensal, error)
```
para:
```go
func (s *DashboardService) ResumoMensal(usuarioID string, mes, ano int) (*ResumoMensal, error)
```

O mesmo se aplica a `EvolucaoMensal`, `ProjecaoProximosMeses` e `DespesasPorCategoria`.

### 2.8 Ajuste do Seed de Categorias

**Problema:** Após a migration, `categorias.usuario_id` será `NOT NULL`. O `SeedCategorias` atual insere sem `usuario_id`, o que causará erro.

**Solução:** Refatorar `SeedCategorias` para receber o `usuarioID` e inserir por usuário. O `ON CONFLICT DO NOTHING` precisará de uma constraint única composta `(nome, usuario_id)` em vez da atual `(nome)` global.

A migration `000014_unique_categoria_nome.up.sql` provavelmente criou `UNIQUE(nome)`. Essa constraint precisará ser substituída.

**Migration adicional para a constraint unique (incluir no 000017 ou como 000018):**

Incluir no final do `000017.up.sql`:
```sql
-- Remover unique global de nome (será por usuario)
ALTER TABLE categorias
    DROP CONSTRAINT IF EXISTS categorias_nome_key;

-- Nova constraint: nome único por usuário
ALTER TABLE categorias
    ADD CONSTRAINT categorias_nome_usuario_unique UNIQUE (usuario_id, nome);
```

E no `000017.down.sql`, antes das remoções de `usuario_id`:
```sql
ALTER TABLE categorias
    DROP CONSTRAINT IF EXISTS categorias_nome_usuario_unique;

ALTER TABLE categorias
    ADD CONSTRAINT categorias_nome_key UNIQUE (nome);
```

**Novo `SeedCategorias` em `seed.go`:**
```go
func SeedCategorias(db *sqlx.DB, usuarioID string) error {
    categorias := []string{
        "Alimentação", "Transporte", "Lazer", "Saúde",
        "Educação", "Moradia", "Vestuário", "Outros",
    }
    for _, nome := range categorias {
        _, err := db.Exec(`
            INSERT INTO categorias (nome, usuario_id)
            VALUES ($1, $2)
            ON CONFLICT (usuario_id, nome) DO NOTHING
        `, nome, usuarioID)
        if err != nil {
            return fmt.Errorf("seed: erro ao inserir categoria %s: %w", nome, err)
        }
    }
    log.Printf("seed: categorias verificadas/inseridas para usuario %s", usuarioID)
    return nil
}
```

**Integração em `SeedUsuarios`:** Após criar (ou verificar) cada usuário, buscar o `id` e chamar `SeedCategorias`:
```go
// Após criar o usuário, buscar o id dele:
var usuarioID string
err = db.Get(&usuarioID, `SELECT id FROM usuarios WHERE email = $1`, u.Email)
if err != nil {
    return fmt.Errorf("seed: erro ao buscar id do usuario %s: %w", u.Email, err)
}
if err := SeedCategorias(db, usuarioID); err != nil {
    return err
}
```

**No `main.go`:** remover a chamada separada a `SeedCategorias(db)` (agora chamado dentro de `SeedUsuarios`).

### 2.9 Teste de Isolamento

O teste de isolamento verifica que usuário A não vê dados do usuário B. Este é um teste de integração que requer um banco de dados real (PostgreSQL). O padrão existente no projeto usa mocks para testes unitários. Para o isolamento, o teste deve ser um teste de integração ou um teste de repositório com banco de teste.

**Estrutura do teste de isolamento em `backend/internal/repository/isolamento_test.go`:**

```go
//go:build integration
// +build integration

package repository_test

import (
    "testing"
    "github.com/jmoiron/sqlx"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    // imports do projeto
)

// TestIsolamentoMembros verifica que o usuário A não vê membros do usuário B.
func TestIsolamentoMembros(t *testing.T) {
    db := conectarBancoDeTeste(t) // helper que conecta em DATABASE_URL de teste
    defer limparBanco(t, db)

    // Criar usuário A
    uidA := criarUsuario(t, db, "a@teste.com")
    // Criar usuário B
    uidB := criarUsuario(t, db, "b@teste.com")

    repo := repository.NewMembroRepository(db)

    // Usuário A cria um membro
    _, err := repo.Criar(uidA, &domain.Membro{Nome: "Ana", Relacionamento: "cônjuge"})
    require.NoError(t, err)

    // Usuário B cria um membro
    _, err = repo.Criar(uidB, &domain.Membro{Nome: "Bruno", Relacionamento: "filho"})
    require.NoError(t, err)

    // Usuário A só vê seus membros
    membrosA, err := repo.Listar(uidA)
    require.NoError(t, err)
    assert.Len(t, membrosA, 1)
    assert.Equal(t, "Ana", membrosA[0].Nome)

    // Usuário B só vê seus membros
    membrosB, err := repo.Listar(uidB)
    require.NoError(t, err)
    assert.Len(t, membrosB, 1)
    assert.Equal(t, "Bruno", membrosB[0].Nome)
}

// TestIsolamentoBuscarPorID verifica que usuário A não acessa registro do usuário B por ID.
func TestIsolamentoBuscarPorID(t *testing.T) {
    db := conectarBancoDeTeste(t)
    defer limparBanco(t, db)

    uidA := criarUsuario(t, db, "a@teste.com")
    uidB := criarUsuario(t, db, "b@teste.com")

    repo := repository.NewMembroRepository(db)

    membroB, err := repo.Criar(uidB, &domain.Membro{Nome: "Bruno"})
    require.NoError(t, err)

    // Usuário A tenta acessar membro do usuário B pelo ID
    _, err = repo.BuscarPorID(uidA, membroB.ID)
    assert.ErrorIs(t, err, domain.ErrMembroNaoEncontrado,
        "usuário A não deve enxergar membro do usuário B")
}
```

**Tag de build `integration`:** O teste só roda com `go test -tags=integration ./...`, não durante o CI normal. Isso evita exigir banco de dados em todo pull request.

### 2.10 Ajuste dos Testes Unitários Existentes

Os mocks nos testes de handler (ex: `MockMembroService` em `membro_handler_test.go`) precisam ter suas assinaturas atualizadas para refletir os novos parâmetros com `usuarioID`. Além disso, o `setupMembroRouter` precisará injetar `"userID"` no contexto Gin antes de chamar os handlers.

**Padrão para injetar `userID` nos testes de handler:**
```go
func setupMembroRouter(svc handler.MembroServiceInterface) *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    // Middleware que injeta userID fake para testes
    r.Use(func(c *gin.Context) {
        c.Set("userID", "usuario-teste-uuid")
        c.Next()
    })
    h := handler.NewMembroHandler(svc)
    // ... registrar rotas
}
```

---

## 3. TT-08: CORS Restrito e Rate Limiting no Login

### 3.1 CORS via Variável de Ambiente

**Problema atual:** `AllowAllOrigins: true` em `main.go` (linha 116) — qualquer origem pode fazer requisições.

**Solução:** Ler `ALLOWED_ORIGINS` do ambiente e construir a whitelist.

#### Passo 1: Adicionar `AllowedOrigins` ao `Config`

Em `backend/config/config.go`, adicionar um novo campo:
```go
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
        trimmed := strings.TrimSpace(o)
        if trimmed != "" {
            origins = append(origins, trimmed)
        }
    }
}
// Fallback para desenvolvimento local se não configurado
if len(origins) == 0 {
    origins = []string{"http://localhost:5173"}
}

cfg := &Config{
    // ...
    Server: ServerConfig{
        Port:           v.GetString("SERVER_PORT"),
        AllowedOrigins: origins,
    },
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

**Railway (backend):**
```
ALLOWED_ORIGINS=https://financas-familiares.vercel.app
```

**Vercel:** Não precisa de configuração de CORS (só o backend precisa).

**Arquivo `.env` local (dev):**
```
ALLOWED_ORIGINS=http://localhost:5173
```

### 3.2 Rate Limiting no Login

**Decisão de implementação:** Usar `golang.org/x/time/rate` — já está disponível como dependência transitiva de `golang.org/x/crypto` (que já está no `go.mod`). Sem dependência nova.

**Estratégia:** Um `sync.Map` que mapeia IP → `*rate.Limiter`. Cada IP tem um limiter de 10 requisições por minuto (token bucket: taxa de 10/60 tokens por segundo, burst de 10).

#### Implementação do middleware de rate limiting

**Novo arquivo: `backend/internal/middleware/ratelimit.go`**

```go
package middleware

import (
    "net/http"
    "sync"

    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

// loginLimiters armazena um rate.Limiter por IP de cliente.
var loginLimiters sync.Map

// getLoginLimiter retorna o limiter existente para o IP ou cria um novo.
// Configuração: 10 requisições por minuto por IP (burst de 10).
func getLoginLimiter(ip string) *rate.Limiter {
    val, _ := loginLimiters.LoadOrStore(ip, rate.NewLimiter(rate.Limit(10.0/60.0), 10))
    return val.(*rate.Limiter)
}

// LoginRateLimiter é um middleware Gin que limita requisições de login por IP.
// Retorna 429 com mensagem em português se o limite for excedido.
func LoginRateLimiter() gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        limiter := getLoginLimiter(ip)
        if !limiter.Allow() {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "erro": "muitas tentativas, aguarde 1 minuto",
            })
            return
        }
        c.Next()
    }
}
```

#### Integração em `main.go`

```go
// Na configuração das rotas de auth:
authGroup := v1.Group("/auth")
{
    authGroup.POST("/login", middleware.LoginRateLimiter(), authHandler.Login)
    authGroup.POST("/refresh", authHandler.Refresh)
}
```

#### Considerações sobre o `sync.Map`

- O `sync.Map` nunca é limpo: com o tempo, acumula IPs. Para produção futura, considerar limpeza periódica ou usar uma biblioteca com TTL. Para o escopo desta sprint (sistema de uso familiar, baixo volume), é aceitável.
- `c.ClientIP()` do Gin já respeita `X-Forwarded-For` quando o servidor está atrás de proxy (Railway usa proxy). Verificar se Railway injeta o IP real nesse header — geralmente sim.

#### Como testar o 429

```bash
# Usando curl em loop (11 requisições consecutivas devem resultar em 429 na décima-primeira):
for i in $(seq 1 11); do
  curl -s -o /dev/null -w "%{http_code}\n" \
    -X POST https://seu-backend.railway.app/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"teste@teste.com","senha":"errada"}'
done
```

As primeiras 10 devem retornar 401 (credenciais inválidas), a 11ª deve retornar 429.

**Teste unitário do middleware:**

```go
// backend/internal/middleware/ratelimit_test.go
func TestLoginRateLimiter_Permite10Requisicoes(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.POST("/login", LoginRateLimiter(), func(c *gin.Context) {
        c.Status(http.StatusOK)
    })

    for i := 0; i < 10; i++ {
        req := httptest.NewRequest(http.MethodPost, "/login", nil)
        w := httptest.NewRecorder()
        r.ServeHTTP(w, req)
        assert.Equal(t, http.StatusOK, w.Code, "requisição %d deveria passar", i+1)
    }
}

func TestLoginRateLimiter_Bloqueia11aRequisicao(t *testing.T) {
    // Usar IP diferente para isolar do teste anterior (sync.Map é global)
    // Alternativa: tornar loginLimiters injetável para testes — aceitável no escopo desta sprint
    // deixar como integration test ou usar IP único via header
}
```

> **Nota:** O `sync.Map` global dificulta o isolamento entre testes unitários do rate limiter. A solução mais simples é marcar o teste do bloqueio como `// +build integration` ou usar um endereço IP único por teste via `X-Forwarded-For` no request.

### 3.3 Variáveis a Configurar no Railway e Vercel

**Railway (backend) — via dashboard ou CLI:**
```
ALLOWED_ORIGINS=https://seu-dominio.vercel.app
```

**Vercel — não requer mudanças de variável para CORS.** Apenas garantir que a URL do backend (`VITE_API_BASE_URL`) já está correta.

---

## 4. TT-09: Índices de Performance

### 4.1 Estratégia

- Usar `CREATE INDEX CONCURRENTLY` para todos os índices: não bloqueia leitura/escrita durante a criação, essencial para produção.
- `CONCURRENTLY` não pode ser usado dentro de uma transação explícita. O `golang-migrate` executa cada arquivo `.up.sql` em uma transação por padrão. **Solução:** usar a diretiva especial `-- +migrate no-transaction` no topo do arquivo, ou criar os índices em migrations separadas.

A biblioteca `golang-migrate` não suporta `-- +migrate no-transaction` nativamente (essa é uma convenção do `sql-migrate`). Para `golang-migrate` com `iofs`, a solução é:

**Opção A (recomendada):** Não usar `CONCURRENTLY` — aceitar o bloqueio breve. Para um banco pequeno (sistema familiar), o bloqueio será milissegundos.

**Opção B:** Rodar os índices manualmente no banco de produção fora do processo de migration.

Este plano usa **Opção A** — sem `CONCURRENTLY` nas migrations, mas com comentário explicando o motivo. Para um sistema de finanças familiar com volume baixo, é a escolha pragmática correta.

### 4.2 Migration 000018 — Índices de Performance

**Arquivo: `backend/migrations/000018_add_performance_indexes.up.sql`**

```sql
-- Sprint 9 / TT-09: Índices de performance
-- Nota: não usamos CONCURRENTLY pois golang-migrate executa em transação.
-- Para banco de produção com alto volume, executar manualmente com CONCURRENTLY.

-- =====================================================================
-- ÍNDICES EM deleted_at (partial index: WHERE deleted_at IS NULL)
-- Acelera todas as queries com WHERE deleted_at IS NULL
-- =====================================================================
CREATE INDEX IF NOT EXISTS idx_membros_deleted_at
    ON membros (deleted_at)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_categorias_deleted_at
    ON categorias (deleted_at)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_cartoes_credito_deleted_at
    ON cartoes_credito (deleted_at)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_despesas_cartao_deleted_at
    ON despesas_cartao (deleted_at)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_assinaturas_deleted_at
    ON assinaturas (deleted_at)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_contas_fixas_deleted_at
    ON contas_fixas (deleted_at)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_despesas_gerais_deleted_at
    ON despesas_gerais (deleted_at)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_rendas_fixas_deleted_at
    ON rendas_fixas (deleted_at)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_rendas_variaveis_deleted_at
    ON rendas_variaveis (deleted_at)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_rendas_extras_deleted_at
    ON rendas_extras (deleted_at)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_rendimentos_investimento_deleted_at
    ON rendimentos_investimento (deleted_at)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_deleted_at
    ON refresh_tokens (deleted_at)
    WHERE deleted_at IS NULL;

-- =====================================================================
-- ÍNDICES EM usuario_id (depende do TT-07 — migration 000017)
-- Acelera todas as queries WHERE usuario_id = $1
-- =====================================================================
CREATE INDEX IF NOT EXISTS idx_membros_usuario_id
    ON membros (usuario_id);

CREATE INDEX IF NOT EXISTS idx_categorias_usuario_id
    ON categorias (usuario_id);

CREATE INDEX IF NOT EXISTS idx_cartoes_credito_usuario_id
    ON cartoes_credito (usuario_id);

CREATE INDEX IF NOT EXISTS idx_despesas_cartao_usuario_id
    ON despesas_cartao (usuario_id);

CREATE INDEX IF NOT EXISTS idx_assinaturas_usuario_id
    ON assinaturas (usuario_id);

CREATE INDEX IF NOT EXISTS idx_contas_fixas_usuario_id
    ON contas_fixas (usuario_id);

CREATE INDEX IF NOT EXISTS idx_despesas_gerais_usuario_id
    ON despesas_gerais (usuario_id);

CREATE INDEX IF NOT EXISTS idx_rendas_fixas_usuario_id
    ON rendas_fixas (usuario_id);

CREATE INDEX IF NOT EXISTS idx_rendas_variaveis_usuario_id
    ON rendas_variaveis (usuario_id);

CREATE INDEX IF NOT EXISTS idx_rendas_extras_usuario_id
    ON rendas_extras (usuario_id);

CREATE INDEX IF NOT EXISTS idx_rendimentos_investimento_usuario_id
    ON rendimentos_investimento (usuario_id);

-- =====================================================================
-- ÍNDICE COMPOSTO em despesas_cartao: (usuario_id, fatura_ano, fatura_mes)
-- Acelera queries de dashboard e projeção por mês/ano filtradas por usuário
-- =====================================================================
CREATE INDEX IF NOT EXISTS idx_despesas_cartao_usuario_fatura
    ON despesas_cartao (usuario_id, fatura_ano, fatura_mes);

-- =====================================================================
-- ÍNDICES EM refresh_tokens
-- =====================================================================
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_usuario_id
    ON refresh_tokens (usuario_id);

-- token_hash já tem UNIQUE constraint (que cria índice automaticamente).
-- Criar índice explícito seria redundante. Verificar se o UNIQUE já serve.
-- Se necessário para partial index (excluindo revogados):
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash_ativo
    ON refresh_tokens (token_hash)
    WHERE revogado = FALSE AND deleted_at IS NULL;
```

**Arquivo: `backend/migrations/000018_add_performance_indexes.down.sql`**

```sql
DROP INDEX IF EXISTS idx_refresh_tokens_token_hash_ativo;
DROP INDEX IF EXISTS idx_refresh_tokens_usuario_id;
DROP INDEX IF EXISTS idx_despesas_cartao_usuario_fatura;
DROP INDEX IF EXISTS idx_rendimentos_investimento_usuario_id;
DROP INDEX IF EXISTS idx_rendas_extras_usuario_id;
DROP INDEX IF EXISTS idx_rendas_variaveis_usuario_id;
DROP INDEX IF EXISTS idx_rendas_fixas_usuario_id;
DROP INDEX IF EXISTS idx_despesas_gerais_usuario_id;
DROP INDEX IF EXISTS idx_contas_fixas_usuario_id;
DROP INDEX IF EXISTS idx_assinaturas_usuario_id;
DROP INDEX IF EXISTS idx_despesas_cartao_usuario_id;
DROP INDEX IF EXISTS idx_cartoes_credito_usuario_id;
DROP INDEX IF EXISTS idx_categorias_usuario_id;
DROP INDEX IF EXISTS idx_membros_usuario_id;
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
```

### 4.3 Observação sobre Partial Index

O partial index `WHERE deleted_at IS NULL` é mais eficiente do que um índice simples em `deleted_at` porque:
- O PostgreSQL mantém apenas as linhas não-deletadas no índice (muito menor).
- As queries com `WHERE deleted_at IS NULL` usam o índice diretamente.
- Linhas soft-deleted (maioria com `deleted_at IS NOT NULL`) não ocupam espaço no índice.

Este é o padrão recomendado para soft delete com PostgreSQL.

---

## 5. Critério de Conclusão

### TT-07 — Isolamento de dados por usuário

- [ ] Migration `000017_add_usuario_id_to_all_tables.up.sql` aplicada sem erros em banco limpo e em banco com dados existentes
- [ ] Migration `000017` down funciona corretamente
- [ ] Todos os repositórios (11 + dashboard) filtram por `usuario_id` em todos os métodos de query
- [ ] Todas as interfaces de repositório em `service/` atualizadas
- [ ] Todos os services passam `usuarioID` para o repositório
- [ ] Todas as interfaces públicas de service (usadas pelos handlers) atualizadas
- [ ] Todos os handlers extraem `"userID"` do contexto antes de chamar o service
- [ ] `SeedCategorias` recebe `usuarioID` e insere com `usuario_id`
- [ ] `SeedUsuarios` chama `SeedCategorias` após criar cada usuário
- [ ] `main.go` não chama mais `SeedCategorias` separadamente
- [ ] Constraint `UNIQUE(nome)` em categorias substituída por `UNIQUE(usuario_id, nome)`
- [ ] Todos os testes unitários existentes compilam e passam (mocks atualizados)
- [ ] Teste de integração `TestIsolamentoMembros` passa com banco de teste
- [ ] Teste de integração `TestIsolamentoBuscarPorID` passa
- [ ] `go build ./...` sem erros

### TT-08 — CORS + Rate Limiting

- [ ] `ALLOWED_ORIGINS` lido do ambiente em `config.go`
- [ ] CORS configurado com whitelist no `main.go` (sem `AllowAllOrigins: true`)
- [ ] Teste manual: requisição de origem não listada recebe `403` ou cabeçalho CORS ausente
- [ ] Frontend em produção (Vercel) funciona normalmente com a nova configuração
- [ ] Middleware `LoginRateLimiter` implementado em `middleware/ratelimit.go`
- [ ] Rate limiter aplicado apenas em `POST /api/v1/auth/login`
- [ ] Após 10 tentativas seguidas do mesmo IP, a 11ª retorna `429` com `{"erro": "muitas tentativas, aguarde 1 minuto"}`
- [ ] `ALLOWED_ORIGINS` configurado no Railway com o domínio real do Vercel
- [ ] `go build ./...` sem erros
- [ ] Testes do middleware de rate limiting passam

### TT-09 — Índices de performance

- [ ] Migration `000018_add_performance_indexes.up.sql` aplicada sem erros
- [ ] Migration `000018` down funciona (todos os índices removidos)
- [ ] Verificar via `psql`: `\d membros` mostra os índices criados
- [ ] Verificar via `EXPLAIN ANALYZE`: query `SELECT ... FROM membros WHERE usuario_id = $1 AND deleted_at IS NULL` usa índice
- [ ] Índice composto `idx_despesas_cartao_usuario_fatura` existe em `despesas_cartao`
- [ ] Índices `idx_refresh_tokens_usuario_id` e `idx_refresh_tokens_token_hash_ativo` existem
- [ ] `go build ./...` sem erros (TT-09 é apenas SQL, não modifica Go)

---

## 6. Arquivos Modificados por Item

### TT-07

| Arquivo | Tipo de Mudança |
|---|---|
| `backend/migrations/000017_add_usuario_id_to_all_tables.up.sql` | **Novo** |
| `backend/migrations/000017_add_usuario_id_to_all_tables.down.sql` | **Novo** |
| `backend/internal/repository/membro_repository.go` | Modificado: adicionar `usuarioID` em todos os métodos |
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
| `backend/internal/repository/dashboard_repository.go` | Modificado: adicionar `usuarioID` em `DespesasPorCategoria` |
| `backend/internal/repository/seed.go` | Modificado: `SeedCategorias` recebe `usuarioID`; `SeedUsuarios` chama seed de categorias |
| `backend/internal/service/membro_service.go` | Modificado: interfaces e métodos com `usuarioID` |
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
| `backend/internal/service/dashboard_service.go` | Modificado: todos os métodos públicos e interfaces `*ForDashboard` |
| `backend/internal/service/renda_historico_service.go` | Modificado: idem |
| `backend/internal/handler/membro_handler.go` | Modificado: extrair `userID`; atualizar interface local |
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
| `backend/internal/handler/*_handler_test.go` (todos) | Modificado: mocks e `setupRouter` com injeção de `userID` |
| `backend/internal/repository/isolamento_test.go` | **Novo**: teste de integração de isolamento |
| `backend/cmd/server/main.go` | Modificado: remover chamada separada a `SeedCategorias` |

### TT-08

| Arquivo | Tipo de Mudança |
|---|---|
| `backend/config/config.go` | Modificado: adicionar `AllowedOrigins []string` em `ServerConfig` |
| `backend/internal/middleware/ratelimit.go` | **Novo** |
| `backend/internal/middleware/ratelimit_test.go` | **Novo** |
| `backend/cmd/server/main.go` | Modificado: substituir CORS; aplicar `LoginRateLimiter` na rota de login |

### TT-09

| Arquivo | Tipo de Mudança |
|---|---|
| `backend/migrations/000018_add_performance_indexes.up.sql` | **Novo** |
| `backend/migrations/000018_add_performance_indexes.down.sql` | **Novo** |

---

### Critical Files for Implementation

- `//wsl.localhost/Arch/home/felipehrs/workspace/pessoal/financas-familiares/backend/cmd/server/main.go` - Ponto central de composição: configura CORS, registra rotas, chama seeds; precisa de 3 mudanças (CORS, rate limit, seed)
- `//wsl.localhost/Arch/home/felipehrs/workspace/pessoal/financas-familiares/backend/internal/middleware/auth.go` - Define a chave `"userID"` no contexto Gin; todos os handlers dependem desta chave para extrair o usuário autenticado
- `//wsl.localhost/Arch/home/felipehrs/workspace/pessoal/financas-familiares/backend/internal/repository/seed.go` - Precisa de refatoração crítica: `SeedCategorias` deve receber `usuarioID` e ser chamada dentro de `SeedUsuarios` para cada usuário
- `//wsl.localhost/Arch/home/felipehrs/workspace/pessoal/financas-familiares/backend/internal/service/dashboard_service.go` - Arquivo mais complexo do TT-07: declara múltiplas interfaces `*ForDashboard` e métodos que precisam propagar `usuarioID` por toda a cadeia de chamadas
- `//wsl.localhost/Arch/home/felipehrs/workspace/pessoal/financas-familiares/backend/internal/repository/membro_repository.go` - Implementação de referência para o padrão de refatoração dos repositórios; todos os outros 10 repositórios seguem a mesma estrutura
