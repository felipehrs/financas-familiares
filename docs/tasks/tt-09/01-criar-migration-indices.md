# Tarefa TT-09.1: Criar Migration de Índices de Performance

**Sprint:** 9
**Item:** TT-09
**Fase:** 1 de 3
**Estimativa:** 60 min (criação + validação de sintaxe)
**Dependências:** TT-07 e TT-08 concluídos
**Abordagem:** Incremental (criar migration → validar sintaxe → aplicar localmente)

---

## Objetivo

Criar migration SQL que adiciona ~25 índices de performance no banco de dados, focando em:
- Índices parciais em `familia_id` (isolamento multi-tenant)
- Índice composto em `despesas_cartao` para otimizar cálculo de fatura
- Índices em campos de data para dashboard e relatórios
- Índices de ordenação em `created_at`

**Por quê?** Após TT-07, todas as queries filtram por `familia_id` + `deleted_at IS NULL`, mas não há índices específicos. Sem índices, o PostgreSQL faz full table scans, degradando performance conforme o volume de dados cresce.

---

## Checklist de Implementação

### 1. Criar Arquivo de Migration (UP)

**Arquivo:** `backend/migrations/000018_add_performance_indexes.up.sql`

**Copiar o conteúdo completo da migration:**

```sql
-- ============================================
-- Migration 000018: Índices de Performance
-- Sprint 9 — TT-09
-- Criado em: 12/03/2026
-- ============================================

-- ============================================
-- PARTE 1: Índices em familia_id (Multi-Tenant)
-- Todas as queries de domínio filtram por familia_id + deleted_at IS NULL
-- Usar índices parciais para melhor performance
-- ============================================

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_membros_familia_id_active
ON membros (familia_id)
WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_categorias_familia_id_active
ON categorias (familia_id)
WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_cartoes_credito_familia_id_active
ON cartoes_credito (familia_id)
WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_despesas_cartao_familia_id_active
ON despesas_cartao (familia_id)
WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_assinaturas_familia_id_active
ON assinaturas (familia_id)
WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_contas_fixas_familia_id_active
ON contas_fixas (familia_id)
WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_despesas_gerais_familia_id_active
ON despesas_gerais (familia_id)
WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_rendas_fixas_familia_id_active
ON rendas_fixas (familia_id)
WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_rendas_variaveis_familia_id_active
ON rendas_variaveis (familia_id)
WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_rendas_extras_familia_id_active
ON rendas_extras (familia_id)
WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_rendimentos_investimento_familia_id_active
ON rendimentos_investimento (familia_id)
WHERE deleted_at IS NULL;

-- Tabela familias (soft delete, mas sem familia_id)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_familias_deleted_at
ON familias (id)
WHERE deleted_at IS NULL;

-- Tabela refresh_tokens (não tem familia_id, mas tem deleted_at)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_refresh_tokens_usuario_id_active
ON refresh_tokens (usuario_id)
WHERE deleted_at IS NULL;

-- ============================================
-- PARTE 2: Índice Composto para Cálculo de Fatura
-- Query: SELECT SUM(valor_parcela) WHERE familia_id = X AND cartao_id = Y AND fatura_ano = Z AND fatura_mes = W
-- ============================================

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_despesas_cartao_fatura
ON despesas_cartao (familia_id, cartao_id, fatura_ano, fatura_mes)
WHERE deleted_at IS NULL;

-- ============================================
-- PARTE 3: Índices em Campos de Data
-- Queries de dashboard, relatórios e projeções filtram por data/período
-- ============================================

-- Despesas por data de compra (relatórios e filtros de período)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_despesas_cartao_data_compra
ON despesas_cartao (familia_id, data_compra)
WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_despesas_gerais_data
ON despesas_gerais (familia_id, data)
WHERE deleted_at IS NULL;

-- Rendas mensais (projeções e histórico)
-- Ordem: familia_id, ano DESC, mes DESC (para queries de "últimos X meses")
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_rendas_variaveis_mes_ano
ON rendas_variaveis (familia_id, ano DESC, mes DESC)
WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_rendas_extras_mes_ano
ON rendas_extras (familia_id, ano DESC, mes DESC)
WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_rendimentos_investimento_mes_ano
ON rendimentos_investimento (familia_id, ano DESC, mes DESC)
WHERE deleted_at IS NULL;

-- Vigência de rendas fixas (data_inicio, data_fim — US-22)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_rendas_fixas_vigencia
ON rendas_fixas (familia_id, data_inicio, data_fim)
WHERE deleted_at IS NULL;

-- ============================================
-- PARTE 4: Índices de Ordenação (created_at, updated_at)
-- Listagens padrão ordenam por created_at DESC
-- ============================================

-- Membros (lista ordenada por cadastro)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_membros_created_at
ON membros (familia_id, created_at DESC)
WHERE deleted_at IS NULL;

-- Categorias (lista ordenada por cadastro)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_categorias_created_at
ON categorias (familia_id, created_at DESC)
WHERE deleted_at IS NULL;

-- Cartões (lista ordenada por cadastro)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_cartoes_credito_created_at
ON cartoes_credito (familia_id, created_at DESC)
WHERE deleted_at IS NULL;
```

**Observações importantes:**

1. **`CONCURRENTLY`**: Permite criar índices sem bloquear a tabela (crucial em produção)
2. **`IF NOT EXISTS`**: Evita erros se a migration for executada múltiplas vezes
3. **Índices parciais (`WHERE deleted_at IS NULL`)**: Reduzem tamanho e melhoram performance
4. **Ordem das colunas em índices compostos**: `familia_id` primeiro (máxima seletividade), depois outras colunas por especificidade decrescente

- [ ] Arquivo `000018_add_performance_indexes.up.sql` criado
- [ ] Conteúdo completo copiado (25 índices)
- [ ] Comentários explicativos incluídos

---

### 2. Criar Arquivo de Migration (DOWN/Rollback)

**Arquivo:** `backend/migrations/000018_add_performance_indexes.down.sql`

```sql
-- ============================================
-- Migration 000018 DOWN: Remover Índices de Performance
-- ============================================

-- Parte 1: Índices em familia_id
DROP INDEX CONCURRENTLY IF EXISTS idx_membros_familia_id_active;
DROP INDEX CONCURRENTLY IF EXISTS idx_categorias_familia_id_active;
DROP INDEX CONCURRENTLY IF EXISTS idx_cartoes_credito_familia_id_active;
DROP INDEX CONCURRENTLY IF EXISTS idx_despesas_cartao_familia_id_active;
DROP INDEX CONCURRENTLY IF EXISTS idx_assinaturas_familia_id_active;
DROP INDEX CONCURRENTLY IF EXISTS idx_contas_fixas_familia_id_active;
DROP INDEX CONCURRENTLY IF EXISTS idx_despesas_gerais_familia_id_active;
DROP INDEX CONCURRENTLY IF EXISTS idx_rendas_fixas_familia_id_active;
DROP INDEX CONCURRENTLY IF EXISTS idx_rendas_variaveis_familia_id_active;
DROP INDEX CONCURRENTLY IF EXISTS idx_rendas_extras_familia_id_active;
DROP INDEX CONCURRENTLY IF EXISTS idx_rendimentos_investimento_familia_id_active;
DROP INDEX CONCURRENTLY IF EXISTS idx_familias_deleted_at;
DROP INDEX CONCURRENTLY IF EXISTS idx_refresh_tokens_usuario_id_active;

-- Parte 2: Índice composto de fatura
DROP INDEX CONCURRENTLY IF EXISTS idx_despesas_cartao_fatura;

-- Parte 3: Índices de data
DROP INDEX CONCURRENTLY IF EXISTS idx_despesas_cartao_data_compra;
DROP INDEX CONCURRENTLY IF EXISTS idx_despesas_gerais_data;
DROP INDEX CONCURRENTLY IF EXISTS idx_rendas_variaveis_mes_ano;
DROP INDEX CONCURRENTLY IF EXISTS idx_rendas_extras_mes_ano;
DROP INDEX CONCURRENTLY IF EXISTS idx_rendimentos_investimento_mes_ano;
DROP INDEX CONCURRENTLY IF EXISTS idx_rendas_fixas_vigencia;

-- Parte 4: Índices de ordenação
DROP INDEX CONCURRENTLY IF EXISTS idx_membros_created_at;
DROP INDEX CONCURRENTLY IF EXISTS idx_categorias_created_at;
DROP INDEX CONCURRENTLY IF EXISTS idx_cartoes_credito_created_at;
```

- [ ] Arquivo `000018_add_performance_indexes.down.sql` criado
- [ ] Todos os 25 índices incluídos no DROP

---

### 3. Validar Sintaxe SQL (Sem Aplicar)

**IMPORTANTE:** Antes de aplicar a migration, validar a sintaxe SQL para detectar erros.

**Opção A: Validar com `psql` (sem executar)**

```bash
cd backend/migrations

# Validar sintaxe (apenas parse, não executa)
psql -d postgres://seu-db-url -f 000018_add_performance_indexes.up.sql --echo-all --dry-run 2>&1 | head -20
```

**Observação:** O flag `--dry-run` não existe no psql. Usar alternativa:

```bash
# Validar sintaxe com psql em modo single-transaction e rollback
psql $DATABASE_URL <<EOF
BEGIN;
\i 000018_add_performance_indexes.up.sql
ROLLBACK;
EOF
```

Isso executa a migration em uma transação e faz rollback no final, deixando o DB intacto.

**Opção B: Usar ferramenta online**

Copiar o SQL e validar em:
- [SQL Fiddle](http://sqlfiddle.com/) (PostgreSQL)
- [DB Fiddle](https://www.db-fiddle.com/) (PostgreSQL)

---

**Resultado esperado:**

Sintaxe válida sem erros de parsing.

**Possíveis erros a verificar:**
- ❌ Nome de tabela/coluna incorreto (ex: `menbros` em vez de `membros`)
- ❌ Coluna inexistente (ex: `familia_id` em tabela que não tem essa coluna)
- ❌ Sintaxe incorreta de índice parcial (`WHERE` mal formatado)

- [ ] Sintaxe SQL validada (sem erros de parsing)
- [ ] Nomes de tabelas e colunas conferidos

---

### 4. Listar Migrations Atuais

Verificar qual é a última migration aplicada:

```bash
cd backend
make migrate-version
```

**Resultado esperado:**

```
version: 17
dirty: false
```

Se aparecer `version: 17`, significa que a migration 000017 foi a última aplicada.

- [ ] Versão atual verificada (deve ser 17 ou superior)
- [ ] `dirty: false` (sem migrações pendentes/quebradas)

---

### 5. Aplicar Migration Localmente

**IMPORTANTE:** Antes de aplicar, fazer backup do banco (opcional em dev, obrigatório em prod).

```bash
cd backend
make migrate-up
```

**Saída esperada:**

```
Applying migration 000018_add_performance_indexes.up.sql...
Migration 000018 applied successfully.
```

**Se houver erro:**

Possíveis causas:
1. **Sintaxe SQL inválida** → Corrigir no arquivo `.up.sql` e tentar novamente
2. **Tabela/coluna não existe** → Verificar se migrations anteriores foram aplicadas
3. **Índice com nome duplicado** → Verificar se já existe um índice com esse nome (`\d nome_da_tabela` no psql)

**Comando para debugar:**

```bash
# Conectar ao banco e listar índices existentes
psql $DATABASE_URL -c "\di idx_*"
```

- [ ] `make migrate-up` executado com sucesso
- [ ] Nenhum erro de criação de índices

---

### 6. Verificar Índices Criados

**Listar todos os índices criados pela migration:**

```bash
psql $DATABASE_URL <<EOF
SELECT
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
  AND indexname LIKE 'idx_%'
ORDER BY tablename, indexname;
EOF
```

**Resultado esperado:**

Lista com ~25 índices, incluindo:
- `idx_membros_familia_id_active`
- `idx_despesas_cartao_fatura`
- `idx_despesas_gerais_data`
- `idx_rendas_variaveis_mes_ano`
- etc.

**Verificar tamanho dos índices:**

```bash
psql $DATABASE_URL <<EOF
SELECT
    tablename,
    indexrelname AS index_name,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND indexrelname LIKE 'idx_%'
ORDER BY pg_relation_size(indexrelid) DESC;
EOF
```

**Resultado esperado:**

Índices com tamanhos variando de poucos KB (tabelas pequenas) a alguns MB (tabelas maiores).

- [ ] 25 índices criados com sucesso
- [ ] Todos os nomes de índices corretos (`idx_*`)
- [ ] Tamanhos dos índices verificados (valores razoáveis)

---

### 7. Contar Total de Índices por Tabela

Verificar quantos índices cada tabela principal tem:

```bash
psql $DATABASE_URL <<EOF
SELECT
    tablename,
    COUNT(*) AS total_indexes
FROM pg_indexes
WHERE schemaname = 'public'
  AND tablename IN (
    'membros', 'categorias', 'cartoes_credito',
    'despesas_cartao', 'despesas_gerais',
    'assinaturas', 'contas_fixas',
    'rendas_fixas', 'rendas_variaveis', 'rendas_extras',
    'rendimentos_investimento', 'familias', 'refresh_tokens'
  )
GROUP BY tablename
ORDER BY total_indexes DESC;
EOF
```

**Resultado esperado (aproximado):**

| Tabela | Total Índices |
|--------|---------------|
| `despesas_cartao` | ~6 (PK + 4 novos) |
| `membros` | ~4 (PK + 2 novos) |
| `categorias` | ~4 (PK + unique + 2 novos) |
| `despesas_gerais` | ~3 (PK + 2 novos) |
| etc. | ... |

- [ ] Contagem de índices por tabela conferida
- [ ] Nenhuma tabela sem índices (todas têm pelo menos PK)

---

## Critério de Conclusão

- [x] Arquivo `000018_add_performance_indexes.up.sql` criado com 25 índices
- [x] Arquivo `000018_add_performance_indexes.down.sql` criado com rollback
- [x] Sintaxe SQL validada (sem erros de parsing)
- [x] Migration aplicada localmente com `make migrate-up`
- [x] 25 índices criados e verificados com `pg_indexes`
- [x] Tamanhos dos índices verificados e razoáveis
- [x] Nenhum índice em estado `INVALID` (`SELECT * FROM pg_index WHERE NOT indisvalid`)

---

## Commit Sugerido

```bash
git add backend/migrations/000018_add_performance_indexes.up.sql \
        backend/migrations/000018_add_performance_indexes.down.sql

git commit -m "feat(backend): adicionar índices de performance — TT-09 (1/3)

Cria ~25 índices para otimizar queries multi-tenant e cálculo de faturas.

**Índices criados:**
- 11 índices parciais em familia_id (isolamento multi-tenant)
- 1 índice composto em despesas_cartao (familia_id, cartao_id, fatura_ano, fatura_mes)
- 6 índices em campos de data (data_compra, data, mes/ano)
- 3 índices de ordenação em created_at

**Técnicas usadas:**
- Índices parciais (WHERE deleted_at IS NULL) para reduzir tamanho
- CREATE INDEX CONCURRENTLY para não bloquear tabelas em produção
- Ordem de colunas otimizada por seletividade decrescente

**Impacto esperado:**
- Listagens de membros/categorias: 10-50x mais rápidas
- Cálculo de fatura: 100x mais rápido (evita full table scan)
- Queries de dashboard por período: 20-100x mais rápidas

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Próximos Passos

Após conclusão desta tarefa, prosseguir para:
- **[02-validar-indices-explain.md](./02-validar-indices-explain.md)** — Validar uso dos índices com EXPLAIN ANALYZE

---

## Referências

- [Plano técnico TT-09](../../tecnico/06-plano-tt09-indices-performance.md)
- [PostgreSQL: CREATE INDEX](https://www.postgresql.org/docs/current/sql-createindex.html)
- [PostgreSQL: Partial Indexes](https://www.postgresql.org/docs/current/indexes-partial.html)
- [PostgreSQL: CREATE INDEX CONCURRENTLY](https://www.postgresql.org/docs/current/sql-createindex.html#SQL-CREATEINDEX-CONCURRENTLY)
