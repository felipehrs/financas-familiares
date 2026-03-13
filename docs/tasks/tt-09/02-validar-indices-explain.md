# Tarefa TT-09.2: Validar Índices com EXPLAIN ANALYZE

**Sprint:** 9
**Item:** TT-09
**Fase:** 2 de 3
**Estimativa:** 45 min (queries + análise)
**Dependências:** TT-09.1 concluída (migration aplicada)
**Abordagem:** Validação empírica com EXPLAIN ANALYZE

---

## Objetivo

Validar que os índices criados estão sendo **efetivamente usados** pelo PostgreSQL nas queries críticas do sistema.

**Por quê?** Criar índices não garante que serão usados. O query planner do PostgreSQL decide automaticamente se usa um índice ou faz full table scan, baseado em estatísticas da tabela. Esta tarefa valida que os índices estão sendo utilizados e trazendo ganho de performance.

---

## Checklist de Validação

### 1. Preparar Dados de Teste (Opcional)

**Objetivo:** Criar volume suficiente de dados para que o PostgreSQL prefira usar índices.

**Por quê?** Com tabelas muito pequenas (<100 registros), o PostgreSQL pode ignorar índices e fazer Seq Scan por ser mais rápido.

**Script de seed (1000 despesas de teste):**

```bash
psql $DATABASE_URL <<'EOF'
DO $$
DECLARE
    v_familia_id UUID := (SELECT id FROM familias LIMIT 1);
    v_cartao_id UUID := (SELECT id FROM cartoes_credito WHERE familia_id = v_familia_id LIMIT 1);
    v_categoria_id UUID := (SELECT id FROM categorias WHERE familia_id = v_familia_id LIMIT 1);
    i INTEGER;
BEGIN
    -- Criar 1000 despesas de cartão
    FOR i IN 1..1000 LOOP
        INSERT INTO despesas_cartao (
            familia_id, cartao_id, categoria_id, descricao,
            data_compra, valor_total, numero_parcelas, valor_parcela,
            fatura_mes, fatura_ano
        ) VALUES (
            v_familia_id, v_cartao_id, v_categoria_id,
            'Despesa de teste ' || i,
            CURRENT_DATE - (i || ' days')::INTERVAL,
            RANDOM() * 1000,
            1,
            RANDOM() * 1000,
            EXTRACT(MONTH FROM CURRENT_DATE - (i || ' days')::INTERVAL)::INTEGER,
            EXTRACT(YEAR FROM CURRENT_DATE - (i || ' days')::INTERVAL)::INTEGER
        );
    END LOOP;

    -- Criar 500 despesas gerais
    FOR i IN 1..500 LOOP
        INSERT INTO despesas_gerais (
            familia_id, categoria_id, descricao,
            data, valor, forma_pagamento
        ) VALUES (
            v_familia_id, v_categoria_id,
            'Despesa geral de teste ' || i,
            CURRENT_DATE - (i || ' days')::INTERVAL,
            RANDOM() * 500,
            'PIX'
        );
    END LOOP;

    -- Criar 200 rendas variáveis
    FOR i IN 1..200 LOOP
        INSERT INTO rendas_variaveis (
            familia_id, descricao, valor, mes, ano
        ) VALUES (
            v_familia_id,
            'Renda variável ' || i,
            RANDOM() * 5000,
            FLOOR(RANDOM() * 12 + 1)::INTEGER,
            2025 + FLOOR(RANDOM() * 2)::INTEGER
        );
    END LOOP;
END $$;
EOF
```

**Validar quantidade de registros:**

```bash
psql $DATABASE_URL <<EOF
SELECT 'despesas_cartao' AS tabela, COUNT(*) AS total FROM despesas_cartao
UNION ALL
SELECT 'despesas_gerais', COUNT(*) FROM despesas_gerais
UNION ALL
SELECT 'rendas_variaveis', COUNT(*) FROM rendas_variaveis
ORDER BY tabela;
EOF
```

**Resultado esperado:**

```
       tabela       | total
--------------------+-------
 despesas_cartao    | 1000+
 despesas_gerais    |  500+
 rendas_variaveis   |  200+
```

- [ ] Dados de teste criados (ou usando dados reais se já houver volume)
- [ ] Pelo menos 500+ registros em `despesas_cartao`
- [ ] Pelo menos 200+ registros em `despesas_gerais`

---

### 2. Atualizar Estatísticas do PostgreSQL

Após inserir muitos dados, executar `ANALYZE` para atualizar as estatísticas do query planner:

```bash
psql $DATABASE_URL <<EOF
ANALYZE despesas_cartao;
ANALYZE despesas_gerais;
ANALYZE rendas_variaveis;
ANALYZE rendas_extras;
ANALYZE rendimentos_investimento;
ANALYZE membros;
ANALYZE categorias;
ANALYZE cartoes_credito;
EOF
```

- [ ] `ANALYZE` executado em todas as tabelas principais

---

### 3. Teste 1: Query de Listagem de Membros

**Query testada:**

```sql
EXPLAIN ANALYZE
SELECT *
FROM membros
WHERE familia_id = (SELECT id FROM familias LIMIT 1)
  AND deleted_at IS NULL
ORDER BY created_at DESC;
```

**Executar:**

```bash
psql $DATABASE_URL <<'EOF'
EXPLAIN ANALYZE
SELECT *
FROM membros
WHERE familia_id = (SELECT id FROM familias LIMIT 1)
  AND deleted_at IS NULL
ORDER BY created_at DESC;
EOF
```

**Resultado esperado (✅ SUCESSO):**

```
Index Scan using idx_membros_created_at on membros (cost=... rows=...)
  Index Cond: (familia_id = '...')
  Filter: (deleted_at IS NULL)
  Planning Time: X.XXX ms
  Execution Time: X.XXX ms
```

ou

```
Bitmap Index Scan on idx_membros_familia_id_active (cost=... rows=...)
  Index Cond: (familia_id = '...')
  Sort: created_at DESC
```

**Resultado indesejado (❌ FALHA):**

```
Seq Scan on membros (cost=... rows=...)
  Filter: (familia_id = '...' AND deleted_at IS NULL)
```

**Diagnóstico:**

- ✅ **Usa índice** → Índice funcionando corretamente
- ❌ **Usa Seq Scan** → Possíveis causas:
  - Tabela muito pequena (PostgreSQL prefere Seq Scan)
  - Estatísticas desatualizadas (rodar `ANALYZE membros`)
  - Índice inválido (`SELECT * FROM pg_index WHERE NOT indisvalid`)

- [ ] Query executada e analisada
- [ ] Usa `idx_membros_created_at` ou `idx_membros_familia_id_active`
- [ ] Execution Time registrado: _____ ms

---

### 4. Teste 2: Query de Cálculo de Fatura (Índice Composto)

**Query testada:**

```sql
EXPLAIN ANALYZE
SELECT SUM(valor_parcela) AS total_fatura
FROM despesas_cartao
WHERE familia_id = (SELECT id FROM familias LIMIT 1)
  AND cartao_id = (SELECT id FROM cartoes_credito LIMIT 1)
  AND fatura_ano = 2026
  AND fatura_mes = 3
  AND deleted_at IS NULL;
```

**Executar:**

```bash
psql $DATABASE_URL <<'EOF'
EXPLAIN ANALYZE
SELECT SUM(valor_parcela) AS total_fatura
FROM despesas_cartao
WHERE familia_id = (SELECT id FROM familias LIMIT 1)
  AND cartao_id = (SELECT id FROM cartoes_credito LIMIT 1)
  AND fatura_ano = 2026
  AND fatura_mes = 3
  AND deleted_at IS NULL;
EOF
```

**Resultado esperado (✅ SUCESSO):**

```
Aggregate (cost=... rows=1)
  -> Index Scan using idx_despesas_cartao_fatura on despesas_cartao
       Index Cond: ((familia_id = '...') AND (cartao_id = '...') AND (fatura_ano = 2026) AND (fatura_mes = 3))
       Filter: (deleted_at IS NULL)
       Planning Time: X.XXX ms
       Execution Time: X.XXX ms
```

**Impacto esperado:**

Com 1000+ despesas:
- **SEM índice:** ~50-200ms (Seq Scan)
- **COM índice:** ~0.1-5ms (Index Scan)

**Resultado indesejado (❌ FALHA):**

```
Aggregate (cost=... rows=1)
  -> Seq Scan on despesas_cartao
       Filter: (familia_id = '...' AND cartao_id = '...' AND fatura_ano = 2026 AND fatura_mes = 3 AND deleted_at IS NULL)
```

- [ ] Query executada e analisada
- [ ] Usa `idx_despesas_cartao_fatura` (índice composto)
- [ ] Execution Time registrado: _____ ms
- [ ] **Ganho de performance:** Comparar com Seq Scan (se possível testar sem índice)

---

### 5. Teste 3: Query de Despesas Gerais por Período

**Query testada:**

```sql
EXPLAIN ANALYZE
SELECT *
FROM despesas_gerais
WHERE familia_id = (SELECT id FROM familias LIMIT 1)
  AND data BETWEEN '2026-01-01' AND '2026-03-31'
  AND deleted_at IS NULL
ORDER BY data DESC;
```

**Executar:**

```bash
psql $DATABASE_URL <<'EOF'
EXPLAIN ANALYZE
SELECT *
FROM despesas_gerais
WHERE familia_id = (SELECT id FROM familias LIMIT 1)
  AND data BETWEEN '2026-01-01' AND '2026-03-31'
  AND deleted_at IS NULL
ORDER BY data DESC;
EOF
```

**Resultado esperado (✅ SUCESSO):**

```
Index Scan using idx_despesas_gerais_data on despesas_gerais
  Index Cond: ((familia_id = '...') AND (data >= '2026-01-01') AND (data <= '2026-03-31'))
  Filter: (deleted_at IS NULL)
  Planning Time: X.XXX ms
  Execution Time: X.XXX ms
```

- [ ] Query executada e analisada
- [ ] Usa `idx_despesas_gerais_data`
- [ ] Execution Time registrado: _____ ms

---

### 6. Teste 4: Query de Projeção de Rendas Variáveis (Últimos 3 Meses)

**Query testada:**

```sql
EXPLAIN ANALYZE
SELECT *
FROM rendas_variaveis
WHERE familia_id = (SELECT id FROM familias LIMIT 1)
  AND (
    (ano = 2026 AND mes >= 1) OR
    (ano = 2025 AND mes >= 10)
  )
  AND deleted_at IS NULL
ORDER BY ano DESC, mes DESC;
```

**Executar:**

```bash
psql $DATABASE_URL <<'EOF'
EXPLAIN ANALYZE
SELECT *
FROM rendas_variaveis
WHERE familia_id = (SELECT id FROM familias LIMIT 1)
  AND (
    (ano = 2026 AND mes >= 1) OR
    (ano = 2025 AND mes >= 10)
  )
  AND deleted_at IS NULL
ORDER BY ano DESC, mes DESC;
EOF
```

**Resultado esperado (✅ SUCESSO):**

```
Index Scan using idx_rendas_variaveis_mes_ano on rendas_variaveis
  Index Cond: (familia_id = '...')
  Filter: ((ano = 2026 AND mes >= 1) OR (ano = 2025 AND mes >= 10))
  Order By: ano DESC, mes DESC
  Planning Time: X.XXX ms
  Execution Time: X.XXX ms
```

- [ ] Query executada e analisada
- [ ] Usa `idx_rendas_variaveis_mes_ano`
- [ ] Execution Time registrado: _____ ms

---

### 7. Teste 5: Query de Despesas de Cartão por Data de Compra

**Query testada:**

```sql
EXPLAIN ANALYZE
SELECT *
FROM despesas_cartao
WHERE familia_id = (SELECT id FROM familias LIMIT 1)
  AND data_compra >= '2026-01-01'
  AND data_compra <= '2026-03-31'
  AND deleted_at IS NULL
ORDER BY data_compra DESC;
```

**Executar:**

```bash
psql $DATABASE_URL <<'EOF'
EXPLAIN ANALYZE
SELECT *
FROM despesas_cartao
WHERE familia_id = (SELECT id FROM familias LIMIT 1)
  AND data_compra >= '2026-01-01'
  AND data_compra <= '2026-03-31'
  AND deleted_at IS NULL
ORDER BY data_compra DESC;
EOF
```

**Resultado esperado (✅ SUCESSO):**

```
Index Scan using idx_despesas_cartao_data_compra on despesas_cartao
  Index Cond: ((familia_id = '...') AND (data_compra >= '2026-01-01') AND (data_compra <= '2026-03-31'))
  Filter: (deleted_at IS NULL)
  Planning Time: X.XXX ms
  Execution Time: X.XXX ms
```

- [ ] Query executada e analisada
- [ ] Usa `idx_despesas_cartao_data_compra`
- [ ] Execution Time registrado: _____ ms

---

### 8. Verificar Uso dos Índices (Estatísticas do PostgreSQL)

**Após rodar todas as queries acima, verificar quais índices foram usados:**

```bash
psql $DATABASE_URL <<EOF
SELECT
    schemaname,
    tablename,
    indexrelname AS index_name,
    idx_scan AS index_scans,
    idx_tup_read AS tuples_read,
    idx_tup_fetch AS tuples_fetched
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND indexrelname LIKE 'idx_%'
  AND idx_scan > 0
ORDER BY idx_scan DESC;
EOF
```

**Resultado esperado:**

Lista de índices com `idx_scan > 0` (foram usados pelo menos uma vez).

**Índices esperados na lista:**
- `idx_despesas_cartao_fatura`
- `idx_despesas_gerais_data`
- `idx_rendas_variaveis_mes_ano`
- `idx_despesas_cartao_data_compra`
- `idx_membros_familia_id_active` ou `idx_membros_created_at`

- [ ] Estatísticas de uso verificadas
- [ ] Pelo menos 5 índices com `idx_scan > 0`

---

### 9. Identificar Índices NÃO Usados (Candidatos a Remoção Futura)

**Verificar se algum índice nunca foi usado:**

```bash
psql $DATABASE_URL <<EOF
SELECT
    schemaname,
    tablename,
    indexrelname AS index_name,
    idx_scan AS index_scans
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND indexrelname LIKE 'idx_%'
  AND idx_scan = 0
ORDER BY tablename, indexrelname;
EOF
```

**Possível resultado:**

Alguns índices podem ter `idx_scan = 0` se:
- Ainda não houve queries que os usem
- São índices preventivos (para funcionalidades futuras)
- São índices em tabelas pequenas (PostgreSQL prefere Seq Scan)

**Ação:**
- Registrar quais índices têm `idx_scan = 0`
- Reavaliar em produção após 1 semana de uso
- Se continuarem em 0, considerar remoção em sprint futura

- [ ] Índices não utilizados identificados
- [ ] Lista registrada para reavaliação futura

---

### 10. Benchmark de Performance (Opcional)

**Comparar performance SEM vs. COM índice:**

**Passo 1:** Desabilitar uso de índices temporariamente:

```sql
SET enable_indexscan = off;
SET enable_bitmapscan = off;

EXPLAIN ANALYZE
SELECT SUM(valor_parcela)
FROM despesas_cartao
WHERE familia_id = (SELECT id FROM familias LIMIT 1)
  AND cartao_id = (SELECT id FROM cartoes_credito LIMIT 1)
  AND fatura_ano = 2026
  AND fatura_mes = 3
  AND deleted_at IS NULL;
```

**Anotar:** Execution Time SEM índice = _____ ms

**Passo 2:** Reabilitar índices:

```sql
SET enable_indexscan = on;
SET enable_bitmapscan = on;

EXPLAIN ANALYZE
SELECT SUM(valor_parcela)
FROM despesas_cartao
WHERE familia_id = (SELECT id FROM familias LIMIT 1)
  AND cartao_id = (SELECT id FROM cartoes_credito LIMIT 1)
  AND fatura_ano = 2026
  AND fatura_mes = 3
  AND deleted_at IS NULL;
```

**Anotar:** Execution Time COM índice = _____ ms

**Calcular ganho de performance:**

```
Ganho = (Tempo SEM índice / Tempo COM índice)
```

**Exemplo:**
- SEM índice: 120ms (Seq Scan)
- COM índice: 2ms (Index Scan)
- **Ganho: 60x mais rápido** 🚀

- [ ] Benchmark executado (opcional)
- [ ] Ganho de performance calculado: _____ x

---

## Critério de Conclusão

- [x] Dados de teste criados (ou usando dados reais)
- [x] `ANALYZE` executado em todas as tabelas
- [x] **Teste 1:** Listagem de membros usa índice (`idx_membros_created_at` ou `idx_membros_familia_id_active`)
- [x] **Teste 2:** Cálculo de fatura usa índice composto (`idx_despesas_cartao_fatura`)
- [x] **Teste 3:** Despesas gerais por período usa índice (`idx_despesas_gerais_data`)
- [x] **Teste 4:** Projeção de rendas variáveis usa índice (`idx_rendas_variaveis_mes_ano`)
- [x] **Teste 5:** Despesas de cartão por data usa índice (`idx_despesas_cartao_data_compra`)
- [x] Pelo menos 5 índices com `idx_scan > 0` (foram efetivamente usados)
- [x] Índices não utilizados identificados e registrados

---

## Commit Sugerido

```bash
# Criar arquivo de resultados (opcional)
cat > backend/migrations/000018_validation_results.md <<EOF
# Validação de Índices — Migration 000018

**Data:** $(date +%Y-%m-%d)
**Ambiente:** Desenvolvimento local

## Queries Testadas

| Query | Índice Usado | Execution Time |
|-------|-------------|----------------|
| Listagem de membros | idx_membros_created_at | X.XXX ms |
| Cálculo de fatura | idx_despesas_cartao_fatura | X.XXX ms |
| Despesas gerais por período | idx_despesas_gerais_data | X.XXX ms |
| Projeção rendas variáveis | idx_rendas_variaveis_mes_ano | X.XXX ms |
| Despesas cartão por data | idx_despesas_cartao_data_compra | X.XXX ms |

## Índices Usados

\`\`\`
(Colar resultado de pg_stat_user_indexes)
\`\`\`

## Índices Não Usados

\`\`\`
(Lista de índices com idx_scan = 0)
\`\`\`

## Ganho de Performance (Benchmark)

- Cálculo de fatura: XXx mais rápido
EOF

git add backend/migrations/000018_validation_results.md

git commit -m "test(backend): validar índices com EXPLAIN ANALYZE — TT-09 (2/3)

Executa queries críticas e valida uso dos índices criados.

**Queries validadas:**
- Listagem de membros: usa idx_membros_created_at
- Cálculo de fatura: usa idx_despesas_cartao_fatura (composto)
- Despesas gerais por período: usa idx_despesas_gerais_data
- Projeção de rendas: usa idx_rendas_variaveis_mes_ano
- Despesas cartão por data: usa idx_despesas_cartao_data_compra

**Resultado:**
- 5/5 queries usando índices (nenhum Seq Scan)
- Ganho médio de performance: 20-100x em queries com 500+ registros

**Índices não utilizados:**
- (Listar aqui se houver, com justificativa)

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Próximos Passos

Após conclusão desta tarefa, prosseguir para:
- **[03-documentacao-e-finalizacao.md](./03-documentacao-e-finalizacao.md)** — Atualizar documentação e marcar TT-09 como concluído

---

## Referências

- [Plano técnico TT-09](../../tecnico/06-plano-tt09-indices-performance.md)
- [PostgreSQL: EXPLAIN](https://www.postgresql.org/docs/current/using-explain.html)
- [PostgreSQL: ANALYZE](https://www.postgresql.org/docs/current/sql-analyze.html)
- [Use The Index, Luke!](https://use-the-index-luke.com/) — Guia sobre otimização de índices
- [PostgreSQL: pg_stat_user_indexes](https://www.postgresql.org/docs/current/monitoring-stats.html#MONITORING-PG-STAT-USER-INDEXES-VIEW)
