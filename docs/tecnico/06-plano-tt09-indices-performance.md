# Plano de Execução — TT-09: Índices de Performance

> **Sprint 9** | **Criado em:** 12/03/2026
> **Pré-requisito:** TT-07 e TT-08 concluídos (isolamento por família + segurança)

---

## Índice

1. [Contexto e Motivação](#1-contexto-e-motivação)
2. [Objetivos do TT-09](#2-objetivos-do-tt-09)
3. [Estratégia de Implementação](#3-estratégia-de-implementação)
4. [Análise de Índices Necessários](#4-análise-de-índices-necessários)
5. [Criação da Migration](#5-criação-da-migration)
6. [Testes e Validação](#6-testes-e-validação)
7. [Critério de Conclusão](#7-critério-de-conclusão)
8. [Checklist de Implementação](#8-checklist-de-implementação)

---

## 1. Contexto e Motivação

### Problema Atual

Conforme identificado na avaliação pós-Sprint 8 ([docs/tecnico/02-avaliacao-e-multi-tenant.md](./02-avaliacao-e-multi-tenant.md)) e após a implementação do TT-07 (isolamento por `familia_id`):

**1. Ausência de Índices em Colunas Frequentemente Filtradas**

Todas as queries do sistema agora filtram por `familia_id` (após TT-07), mas não há índices específicos para essa coluna. Além disso, o soft delete (`deleted_at IS NULL`) é aplicado em praticamente todas as consultas, mas também não possui índices.

**Impacto:**
- Queries lentas em tabelas com muitos registros (full table scan)
- Performance degrada conforme o volume de dados cresce
- Dashboard e listagens podem ficar lentos com milhares de lançamentos

**2. Falta de Índice Composto em `despesas_cartao`**

A query de cálculo de fatura filtra simultaneamente por:
- `familia_id`
- `cartao_id`
- `fatura_mes`
- `fatura_ano`
- `deleted_at IS NULL`

Sem um índice composto, o PostgreSQL precisa fazer scans completos ou usar índices sub-ótimos.

**3. Ausência de Índices em Campos de Data**

Queries de dashboard e relatórios filtram por:
- `data_compra` (despesas_cartao, despesas_gerais)
- `mes` e `ano` (rendas_variaveis, rendas_extras)
- `data_inicio` e `data_fim` (rendas_fixas — após Sprint 4.5)

---

## 2. Objetivos do TT-09

1. **Criar índices em `deleted_at`** em todas as tabelas que usam soft delete
2. **Criar índices em `familia_id`** em todas as tabelas de domínio (após TT-07)
3. **Criar índice composto em `despesas_cartao`** para otimizar cálculo de fatura
4. **Criar índices em campos de data** usados em filtros e ordenações
5. **Validar impacto dos índices** com queries de exemplo e `EXPLAIN ANALYZE`

---

## 3. Estratégia de Implementação

### Ordem de Execução

1. **Levantamento completo** de tabelas e colunas que precisam de índices
2. **Criação de migration** única (`000018_add_performance_indexes.up.sql`)
3. **Implementação dos índices** em ordem lógica (simples → compostos)
4. **Testes com EXPLAIN ANALYZE** para validar uso dos índices
5. **Documentação** atualizada com os índices criados

### Princípios de Design de Índices

**Índices Simples:**
- Colunas frequentemente usadas em `WHERE` isoladamente
- Colunas usadas em `ORDER BY`
- Colunas de chaves estrangeiras (quando não há índice automático)

**Índices Compostos:**
- Colunas usadas **juntas** em `WHERE` da mesma query
- Ordem: seletividade decrescente (mais específico primeiro)
- Exemplo: `(familia_id, cartao_id, fatura_ano, fatura_mes)` para fatura de cartão

**Índices Parciais:**
- Quando queremos indexar apenas registros ativos: `WHERE deleted_at IS NULL`
- Reduz tamanho do índice e melhora performance

---

## 4. Análise de Índices Necessários

### 4.1 Tabelas com Soft Delete (deleted_at)

Todas as queries filtram por `deleted_at IS NULL`. Usar **índice parcial** para otimizar:

| Tabela | Índice |
|--------|--------|
| `familias` | `idx_familias_deleted_at` |
| `membros` | `idx_membros_deleted_at` |
| `categorias` | `idx_categorias_deleted_at` |
| `cartoes_credito` | `idx_cartoes_credito_deleted_at` |
| `despesas_cartao` | `idx_despesas_cartao_deleted_at` |
| `assinaturas` | `idx_assinaturas_deleted_at` |
| `contas_fixas` | `idx_contas_fixas_deleted_at` |
| `despesas_gerais` | `idx_despesas_gerais_deleted_at` |
| `rendas_fixas` | `idx_rendas_fixas_deleted_at` |
| `rendas_variaveis` | `idx_rendas_variaveis_deleted_at` |
| `rendas_extras` | `idx_rendas_extras_deleted_at` |
| `rendimentos_investimento` | `idx_rendimentos_investimento_deleted_at` |
| `refresh_tokens` | `idx_refresh_tokens_deleted_at` |

**Observação:** `usuarios` também tem `deleted_at`, mas é consultado raramente (apenas no seed/login). Pode receber índice por completude.

---

### 4.2 Tabelas com `familia_id` (Isolamento Multi-Tenant)

Todas as queries de domínio filtram por `familia_id`. Criar índice simples:

| Tabela | Índice |
|--------|--------|
| `membros` | `idx_membros_familia_id` |
| `categorias` | `idx_categorias_familia_id` |
| `cartoes_credito` | `idx_cartoes_credito_familia_id` |
| `despesas_cartao` | `idx_despesas_cartao_familia_id` |
| `assinaturas` | `idx_assinaturas_familia_id` |
| `contas_fixas` | `idx_contas_fixas_familia_id` |
| `despesas_gerais` | `idx_despesas_gerais_familia_id` |
| `rendas_fixas` | `idx_rendas_fixas_familia_id` |
| `rendas_variaveis` | `idx_rendas_variaveis_familia_id` |
| `rendas_extras` | `idx_rendas_extras_familia_id` |
| `rendimentos_investimento` | `idx_rendimentos_investimento_familia_id` |

**Observação:** Como todas as queries já filtram por `deleted_at IS NULL` **e** `familia_id`, podemos considerar **índices compostos** em vez de índices simples. Ver seção 4.5.

---

### 4.3 Índice Composto em `despesas_cartao` (Cálculo de Fatura)

**Query típica (ver `backend/internal/repository/despesa_cartao_repository.go`):**

```sql
SELECT SUM(valor_parcela) as total_fatura
FROM despesas_cartao
WHERE familia_id = $1
  AND cartao_id = $2
  AND fatura_ano = $3
  AND fatura_mes = $4
  AND deleted_at IS NULL;
```

**Índice composto ideal:**

```sql
CREATE INDEX idx_despesas_cartao_fatura
ON despesas_cartao (familia_id, cartao_id, fatura_ano, fatura_mes)
WHERE deleted_at IS NULL;
```

**Justificativa da ordem:**
1. `familia_id` — máxima seletividade (isola tenant)
2. `cartao_id` — alta seletividade (cada família tem ~2-5 cartões)
3. `fatura_ano` — média seletividade
4. `fatura_mes` — baixa seletividade (1-12)

---

### 4.4 Índices em Campos de Data

**Tabelas que filtram/ordenam por data:**

| Tabela | Coluna(s) | Uso |
|--------|-----------|-----|
| `despesas_cartao` | `data_compra` | Listagens filtradas por período, cálculo de fatura |
| `despesas_gerais` | `data` | Listagens e relatórios por período |
| `rendas_variaveis` | `mes`, `ano` | Projeções e histórico mensal |
| `rendas_extras` | `mes`, `ano` | Histórico de rendas extras |
| `rendimentos_investimento` | `mes`, `ano` | Histórico de rendimentos |
| `rendas_fixas` | `data_inicio`, `data_fim` | Filtro de vigência (RN10) |

**Índices recomendados:**

```sql
-- Despesas com data de compra (queries de dashboard por período)
CREATE INDEX idx_despesas_cartao_data_compra
ON despesas_cartao (familia_id, data_compra)
WHERE deleted_at IS NULL;

CREATE INDEX idx_despesas_gerais_data
ON despesas_gerais (familia_id, data)
WHERE deleted_at IS NULL;

-- Rendas mensais (projeções e histórico)
CREATE INDEX idx_rendas_variaveis_mes_ano
ON rendas_variaveis (familia_id, ano, mes)
WHERE deleted_at IS NULL;

CREATE INDEX idx_rendas_extras_mes_ano
ON rendas_extras (familia_id, ano, mes)
WHERE deleted_at IS NULL;

CREATE INDEX idx_rendimentos_investimento_mes_ano
ON rendimentos_investimento (familia_id, ano, mes)
WHERE deleted_at IS NULL;

-- Vigência de rendas fixas (Sprint 4.5 — US-22)
CREATE INDEX idx_rendas_fixas_vigencia
ON rendas_fixas (familia_id, data_inicio, data_fim)
WHERE deleted_at IS NULL;
```

---

### 4.5 Índices Compostos `(familia_id, deleted_at)` vs. Índices Parciais

**Opção A: Índice Simples em `familia_id` + Filtro em Memória**

```sql
CREATE INDEX idx_membros_familia_id ON membros (familia_id);
-- Query: WHERE familia_id = $1 AND deleted_at IS NULL
-- PostgreSQL usa idx_membros_familia_id e filtra deleted_at em memória
```

**Opção B: Índice Parcial em `familia_id` (apenas registros ativos)**

```sql
CREATE INDEX idx_membros_familia_id_active
ON membros (familia_id)
WHERE deleted_at IS NULL;
-- Query: WHERE familia_id = $1 AND deleted_at IS NULL
-- PostgreSQL usa índice menor e mais eficiente
```

**Decisão:** Usar **Opção B (índices parciais)** para todas as tabelas de domínio.

**Justificativa:**
- Registros deletados (`deleted_at IS NOT NULL`) nunca são consultados
- Índice parcial é **menor** (menos I/O, mais cache-friendly)
- Melhor performance em queries que sempre filtram por `deleted_at IS NULL`

---

### 4.6 Índices em Chaves Estrangeiras

PostgreSQL **não cria índices automáticos** em foreign keys (diferente de outros DBs). Verificar se já existem:

| Tabela | FK | Índice Necessário? |
|--------|----|--------------------|
| `membros` | `familia_id` | ✅ Coberto por `idx_membros_familia_id_active` |
| `categorias` | `familia_id` | ✅ Coberto por `idx_categorias_familia_id_active` |
| `cartoes_credito` | `familia_id`, `membro_id` | ✅ `familia_id` coberto; `membro_id` raramente filtrado isoladamente |
| `despesas_cartao` | `familia_id`, `cartao_id`, `categoria_id` | ✅ Coberto por índice composto de fatura |
| `assinaturas` | `familia_id`, `membro_id`, `categoria_id` | ✅ `familia_id` coberto; outros raramente filtrados |
| `contas_fixas` | `familia_id`, `membro_id`, `categoria_id` | ✅ `familia_id` coberto |
| `despesas_gerais` | `familia_id`, `membro_id`, `categoria_id`, `cartao_id` | ✅ `familia_id` coberto; avaliar `cartao_id` |
| `rendas_fixas` | `familia_id`, `membro_id` | ✅ `familia_id` coberto |
| `rendas_variaveis` | `familia_id`, `membro_id` | ✅ `familia_id` coberto |
| `rendas_extras` | `familia_id`, `membro_id` | ✅ `familia_id` coberto |
| `rendimentos_investimento` | `familia_id`, `membro_id` | ✅ `familia_id` coberto |

**Decisão:** Não criar índices isolados em `membro_id`, `categoria_id` por enquanto. Reavaliar em futuras sprints se surgirem queries específicas.

---

### 4.7 Resumo de Índices a Criar

**Total: ~30 índices distribuídos em:**

1. **Índices parciais em `familia_id`** (11 tabelas)
2. **Índice composto em `despesas_cartao`** para cálculo de fatura
3. **Índices compostos em campos de data** (6 tabelas)

---

## 5. Criação da Migration

### 5.1 Estrutura da Migration

**Arquivo:** `backend/migrations/000018_add_performance_indexes.up.sql`

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

**Observações:**
- **`CONCURRENTLY`**: Permite criar índices sem bloquear a tabela (importante em produção)
- **`IF NOT EXISTS`**: Evita erros se a migration for rodada múltiplas vezes
- **Índices parciais (`WHERE deleted_at IS NULL`)**: Reduzem tamanho e melhoram performance

---

### 5.2 Migration de Rollback

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

---

## 6. Testes e Validação

### 6.1 Aplicar a Migration

```bash
cd backend
make migrate-up
```

**Verificar que a migration foi aplicada:**

```sql
SELECT version, dirty FROM schema_migrations ORDER BY version DESC LIMIT 1;
-- Deve retornar: version = 18, dirty = false
```

---

### 6.2 Listar Índices Criados

**Query para verificar todos os índices do schema:**

```sql
SELECT
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
  AND indexname LIKE 'idx_%'
ORDER BY tablename, indexname;
```

**Verificar tamanho dos índices:**

```sql
SELECT
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) as index_size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND indexrelname LIKE 'idx_%'
ORDER BY pg_relation_size(indexrelid) DESC;
```

---

### 6.3 Validar Uso dos Índices com EXPLAIN ANALYZE

**Teste 1: Query de listagem de membros por família**

```sql
EXPLAIN ANALYZE
SELECT *
FROM membros
WHERE familia_id = 'uuid-da-familia'
  AND deleted_at IS NULL
ORDER BY created_at DESC;
```

**Resultado esperado:**

```
Index Scan using idx_membros_created_at on membros
  Filter: familia_id = 'uuid-da-familia'
```

ou

```
Bitmap Index Scan on idx_membros_familia_id_active
```

**✅ Sucesso:** Se aparecer `Index Scan` ou `Bitmap Index Scan` usando um dos índices criados.
**❌ Falha:** Se aparecer `Seq Scan` (full table scan).

---

**Teste 2: Query de cálculo de fatura**

```sql
EXPLAIN ANALYZE
SELECT SUM(valor_parcela) as total_fatura
FROM despesas_cartao
WHERE familia_id = 'uuid-da-familia'
  AND cartao_id = 'uuid-do-cartao'
  AND fatura_ano = 2026
  AND fatura_mes = 3
  AND deleted_at IS NULL;
```

**Resultado esperado:**

```
Index Scan using idx_despesas_cartao_fatura on despesas_cartao
  Index Cond: (familia_id = 'uuid' AND cartao_id = 'uuid' AND fatura_ano = 2026 AND fatura_mes = 3)
  Filter: deleted_at IS NULL
```

**✅ Sucesso:** Usa `idx_despesas_cartao_fatura` (índice composto).
**❌ Falha:** Usa `Seq Scan` ou índice sub-ótimo.

---

**Teste 3: Query de despesas gerais por período**

```sql
EXPLAIN ANALYZE
SELECT *
FROM despesas_gerais
WHERE familia_id = 'uuid-da-familia'
  AND data BETWEEN '2026-01-01' AND '2026-03-31'
  AND deleted_at IS NULL
ORDER BY data DESC;
```

**Resultado esperado:**

```
Index Scan using idx_despesas_gerais_data on despesas_gerais
  Index Cond: (familia_id = 'uuid' AND data >= '2026-01-01' AND data <= '2026-03-31')
```

---

**Teste 4: Query de projeção de rendas variáveis (últimos 3 meses)**

```sql
EXPLAIN ANALYZE
SELECT *
FROM rendas_variaveis
WHERE familia_id = 'uuid-da-familia'
  AND (ano = 2026 AND mes >= 1) OR (ano = 2025 AND mes >= 10)
  AND deleted_at IS NULL
ORDER BY ano DESC, mes DESC;
```

**Resultado esperado:**

```
Index Scan using idx_rendas_variaveis_mes_ano on rendas_variaveis
  Index Cond: (familia_id = 'uuid')
  Filter: (ano = 2026 AND mes >= 1) OR (ano = 2025 AND mes >= 10)
  Order By: ano DESC, mes DESC
```

---

### 6.4 Benchmark de Performance (Opcional)

**Criar dados de teste (1000 despesas de cartão):**

```sql
DO $$
DECLARE
    v_familia_id UUID := (SELECT id FROM familias LIMIT 1);
    v_cartao_id UUID := (SELECT id FROM cartoes_credito LIMIT 1);
    v_categoria_id UUID := (SELECT id FROM categorias LIMIT 1);
    i INTEGER;
BEGIN
    FOR i IN 1..1000 LOOP
        INSERT INTO despesas_cartao (
            familia_id, cartao_id, categoria_id, descricao,
            data_compra, valor_total, numero_parcelas, valor_parcela,
            fatura_mes, fatura_ano
        ) VALUES (
            v_familia_id, v_cartao_id, v_categoria_id,
            'Despesa de teste ' || i,
            CURRENT_DATE - (i || ' days')::INTERVAL,
            100.00, 1, 100.00,
            EXTRACT(MONTH FROM CURRENT_DATE - (i || ' days')::INTERVAL)::INTEGER,
            EXTRACT(YEAR FROM CURRENT_DATE - (i || ' days')::INTERVAL)::INTEGER
        );
    END LOOP;
END $$;
```

**Medir tempo de execução:**

```sql
-- Desabilitar cache para teste real
SET enable_seqscan = off;

-- Query de fatura (deve usar índice)
EXPLAIN (ANALYZE, BUFFERS)
SELECT SUM(valor_parcela)
FROM despesas_cartao
WHERE familia_id = (SELECT id FROM familias LIMIT 1)
  AND cartao_id = (SELECT id FROM cartoes_credito LIMIT 1)
  AND fatura_ano = 2026
  AND fatura_mes = 3
  AND deleted_at IS NULL;
```

**Comparar:**
- **Execution Time** antes e depois da migration
- **Buffers** (número de blocos lidos — deve diminuir significativamente)

---

## 7. Critério de Conclusão

### Validação Funcional

- [ ] Migration `000018_add_performance_indexes.up.sql` criada e aplicada com sucesso
- [ ] Todos os 25+ índices criados sem erros
- [ ] `make migrate-up` executa sem erros em ambiente de desenvolvimento
- [ ] `make migrate-down` remove todos os índices corretamente (rollback funcional)

### Validação de Performance

- [ ] `EXPLAIN ANALYZE` em query de listagem de membros usa `idx_membros_familia_id_active` ou `idx_membros_created_at`
- [ ] `EXPLAIN ANALYZE` em query de cálculo de fatura usa `idx_despesas_cartao_fatura`
- [ ] `EXPLAIN ANALYZE` em query de despesas por data usa `idx_despesas_gerais_data`
- [ ] `EXPLAIN ANALYZE` em query de rendas mensais usa `idx_rendas_variaveis_mes_ano`
- [ ] Nenhuma query crítica usa `Seq Scan` em tabelas indexadas

### Validação Técnica

- [ ] `go build ./...` — sem erros
- [ ] `go test ./...` — todos os testes passando (índices não devem quebrar testes)
- [ ] Verificação de tamanho dos índices (`pg_relation_size`) — valores razoáveis

### Documentação

- [ ] `docs/tecnico/06-plano-tt09-indices-performance.md` — este arquivo criado
- [ ] `backend/README.md` atualizado com seção "Índices de Performance"
- [ ] `docs/produto/sprints.md` — TT-09 marcado como ✅

---

## 8. Checklist de Implementação

### Fase 1: Criação da Migration

- [ ] Criar arquivo `backend/migrations/000018_add_performance_indexes.up.sql`
- [ ] Implementar **Parte 1:** Índices em `familia_id` (11 tabelas + familias + refresh_tokens)
- [ ] Implementar **Parte 2:** Índice composto em `despesas_cartao` para fatura
- [ ] Implementar **Parte 3:** Índices em campos de data (6 tabelas)
- [ ] Implementar **Parte 4:** Índices de ordenação (`created_at`) em 3 tabelas principais
- [ ] Criar arquivo `backend/migrations/000018_add_performance_indexes.down.sql` com rollback
- [ ] Validar sintaxe SQL (copiar no psql e executar `\i 000018_add_performance_indexes.up.sql`)

### Fase 2: Aplicação e Testes Locais

- [ ] Rodar `make migrate-up` no ambiente de desenvolvimento
- [ ] Verificar que todos os 25+ índices foram criados (`SELECT * FROM pg_indexes WHERE indexname LIKE 'idx_%'`)
- [ ] Executar `EXPLAIN ANALYZE` em 4 queries críticas (ver seção 6.3)
- [ ] Validar que todas usam os índices criados (não `Seq Scan`)
- [ ] Rodar `go test ./...` — garantir que nenhum teste quebrou

### Fase 3: Validação de Rollback

- [ ] Rodar `make migrate-down` uma vez (deve remover migration 000018)
- [ ] Verificar que todos os índices foram removidos
- [ ] Rodar `make migrate-up` novamente (reaplicar migration)
- [ ] Confirmar que os índices foram recriados

### Fase 4: Documentação

- [ ] Atualizar `backend/README.md` com seção "Índices de Performance"
- [ ] Listar os principais índices e seu propósito
- [ ] Adicionar exemplo de como verificar uso de índices com `EXPLAIN ANALYZE`
- [ ] Commit: `feat(backend): adicionar índices de performance — TT-09`

### Fase 5: Deploy e Validação em Produção

- [ ] Push da migration para o repositório
- [ ] Aplicar migration em produção (Railway) via CI/CD ou manual
- [ ] Executar `EXPLAIN ANALYZE` em produção (via Railway CLI ou pgAdmin)
- [ ] Monitorar performance do dashboard após deploy
- [ ] Marcar TT-09 como ✅ em `docs/produto/sprints.md`
- [ ] Commit: `docs: finalizar TT-09 com validação de índices — TT-09 ✅`

---

## Arquivos Modificados (Estimativa)

| Arquivo | Tipo de Mudança |
|---------|-----------------|
| `backend/migrations/000018_add_performance_indexes.up.sql` | **NOVO**: criação de ~25 índices |
| `backend/migrations/000018_add_performance_indexes.down.sql` | **NOVO**: rollback dos índices |
| `backend/README.md` | Adicionar seção "Índices de Performance" |
| `docs/produto/sprints.md` | Marcar TT-09 como ✅ |
| `docs/tecnico/06-plano-tt09-indices-performance.md` | **ESTE ARQUIVO** |

**Total estimado:** 5 arquivos (3 novos, 2 modificados)

---

## Considerações de Produção

### Aplicação de Índices em Produção

**Problema:** Criar índices em produção pode bloquear a tabela por alguns segundos/minutos (dependendo do tamanho).

**Solução:** Usar `CREATE INDEX CONCURRENTLY` (já incluído na migration).

**Cuidados:**
- `CONCURRENTLY` não funciona dentro de transações (migrations devem ser executadas sem `BEGIN/COMMIT` para esses comandos)
- Se a criação falhar, o índice fica em estado `INVALID` — deve ser dropado e recriado

**Validar se algum índice ficou inválido:**

```sql
SELECT indexrelid::regclass AS index_name
FROM pg_index
WHERE NOT indisvalid;
```

Se houver algum, dropar e recriar:

```sql
DROP INDEX CONCURRENTLY nome_do_indice_invalido;
CREATE INDEX CONCURRENTLY nome_do_indice ...;
```

---

### Monitoramento de Performance Pós-Deploy

**Queries para monitorar uso de índices:**

**1. Índices nunca usados (candidatos a remoção futura):**

```sql
SELECT
    schemaname,
    tablename,
    indexrelname AS index_name,
    idx_scan AS index_scans
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND idx_scan = 0
  AND indexrelname LIKE 'idx_%'
ORDER BY tablename, indexrelname;
```

**2. Índices mais usados:**

```sql
SELECT
    schemaname,
    tablename,
    indexrelname AS index_name,
    idx_scan AS index_scans,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND indexrelname LIKE 'idx_%'
ORDER BY idx_scan DESC
LIMIT 20;
```

**3. Tabelas com mais Seq Scans (sem usar índices):**

```sql
SELECT
    schemaname,
    tablename,
    seq_scan AS sequential_scans,
    seq_tup_read AS rows_read_sequentially,
    idx_scan AS index_scans,
    n_live_tup AS estimated_rows
FROM pg_stat_user_tables
WHERE schemaname = 'public'
  AND seq_scan > 0
ORDER BY seq_scan DESC;
```

**Ação:** Se após 1 semana em produção alguns índices tiverem `idx_scan = 0`, considerar removê-los em uma sprint futura.

---

### Manutenção de Índices

**VACUUM e ANALYZE:**

PostgreSQL atualiza automaticamente as estatísticas dos índices, mas em tabelas com muitos `INSERT/UPDATE/DELETE`, pode ser necessário executar manualmente:

```sql
VACUUM ANALYZE despesas_cartao;
VACUUM ANALYZE despesas_gerais;
```

**REINDEX (raramente necessário):**

Se um índice ficar corrompido ou muito fragmentado:

```sql
REINDEX INDEX CONCURRENTLY idx_despesas_cartao_fatura;
```

---

## Possíveis Melhorias Futuras (Fora do Escopo do TT-09)

1. **Índices em `membro_id`** nas tabelas de domínio (se surgirem queries "todas as despesas do membro X")
2. **Índices GIN em campos JSONB** (se adicionarmos campos de metadados no futuro)
3. **Índices de texto completo (GIN/GIST)** em `descricao` para busca textual
4. **Particionamento de tabelas** (ex: `despesas_cartao` particionada por ano) se o volume crescer muito
5. **Índices em `categoria_id`** se houver relatórios "todas as despesas da categoria X independente da família"

---

## Referências

- [PostgreSQL Documentation: Indexes](https://www.postgresql.org/docs/current/indexes.html)
- [PostgreSQL: Partial Indexes](https://www.postgresql.org/docs/current/indexes-partial.html)
- [PostgreSQL: EXPLAIN](https://www.postgresql.org/docs/current/using-explain.html)
- [PostgreSQL: CREATE INDEX CONCURRENTLY](https://www.postgresql.org/docs/current/sql-createindex.html#SQL-CREATEINDEX-CONCURRENTLY)
- [Use The Index, Luke!](https://use-the-index-luke.com/) — guia completo sobre índices SQL
- [PostgreSQL Wiki: Index Maintenance](https://wiki.postgresql.org/wiki/Index_Maintenance)

---

**Fim do Plano de Execução TT-09** 🚀
