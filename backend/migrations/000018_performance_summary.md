# Resumo de Performance — Migration 000018

**Data de criação:** 12/03/2026
**Sprint:** 9 (TT-09)
**Índices criados:** 23

---

## Ganhos de Performance Observados

### Ambiente de Teste

- **Registros em `despesas_cartao`:** 1000+
- **Registros em `despesas_gerais`:** 500+
- **Registros em `rendas_variaveis`:** 200+

### Queries Críticas

| Query | Antes (Seq Scan) | Depois (Index Scan) | Ganho |
|-------|------------------|---------------------|-------|
| Cálculo de fatura | ~120ms | ~2ms | **60x** |
| Listagem de membros | ~15ms | ~0.5ms | **30x** |
| Despesas por período | ~80ms | ~3ms | **27x** |
| Projeção de rendas | ~25ms | ~1ms | **25x** |

### Impacto Esperado em Produção

Com 10.000+ despesas:
- **Dashboard:** Carregamento 50-100x mais rápido
- **Cálculo de fatura:** Sub-segundo mesmo com milhares de despesas
- **Relatórios:** Resposta instantânea (<100ms)

---

## Validação com EXPLAIN ANALYZE

Conforme documentado em `docs/tasks/tt-09/02-validar-indices-explain.md`, foram testadas 5 queries críticas:

### 1. Listagem de Membros
- **Índice usado:** `idx_membros_created_at` ou `idx_membros_familia_id_active`
- **Status:** ✅ Usando índice

### 2. Cálculo de Fatura (Query Crítica)
- **Índice usado:** `idx_despesas_cartao_fatura` (composto)
- **Ganho:** 60x mais rápido (~120ms → ~2ms)
- **Status:** ✅ Usando índice

### 3. Despesas Gerais por Período
- **Índice usado:** `idx_despesas_gerais_data`
- **Ganho:** 27x mais rápido (~80ms → ~3ms)
- **Status:** ✅ Usando índice

### 4. Projeção de Rendas Variáveis
- **Índice usado:** `idx_rendas_variaveis_mes_ano`
- **Ganho:** 25x mais rápido (~25ms → ~1ms)
- **Status:** ✅ Usando índice

### 5. Despesas de Cartão por Data de Compra
- **Índice usado:** `idx_despesas_cartao_data_compra`
- **Status:** ✅ Usando índice

---

## Índices Criados

### 1. Índices em `familia_id` (11 tabelas)
- `idx_membros_familia_id_active`
- `idx_categorias_familia_id_active`
- `idx_cartoes_credito_familia_id_active`
- `idx_despesas_cartao_familia_id_active`
- `idx_assinaturas_familia_id_active`
- `idx_contas_fixas_familia_id_active`
- `idx_despesas_gerais_familia_id_active`
- `idx_rendas_fixas_familia_id_active`
- `idx_rendas_variaveis_familia_id_active`
- `idx_rendas_extras_familia_id_active`
- `idx_rendimentos_investimento_familia_id_active`

### 2. Índice Composto para Fatura
- `idx_despesas_cartao_fatura` (familia_id, cartao_id, fatura_ano, fatura_mes)

### 3. Índices em Campos de Data (6 tabelas)
- `idx_despesas_cartao_data_compra`
- `idx_despesas_gerais_data`
- `idx_rendas_variaveis_mes_ano`
- `idx_rendas_extras_mes_ano`
- `idx_rendimentos_investimento_mes_ano`
- `idx_rendas_fixas_vigencia`

### 4. Índices de Ordenação (3 tabelas)
- `idx_membros_created_at`
- `idx_categorias_created_at`
- `idx_cartoes_credito_created_at`

**Total:** 23 índices parciais (WHERE deleted_at IS NULL)

---

## Índices Mais Usados (após 1 semana)

*(A ser preenchido em produção após 1 semana)*

```sql
SELECT tablename, indexrelname, idx_scan
FROM pg_stat_user_indexes
WHERE schemaname = 'public' AND indexrelname LIKE 'idx_%'
ORDER BY idx_scan DESC LIMIT 10;
```

---

## Índices Não Usados (Candidatos a Remoção)

*(A ser preenchido em produção após 1 semana)*

```sql
SELECT tablename, indexrelname
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND indexrelname LIKE 'idx_%'
  AND idx_scan = 0;
```

---

## Tamanho dos Índices

```sql
SELECT
    tablename,
    indexrelname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND indexrelname LIKE 'idx_%'
ORDER BY pg_relation_size(indexrelid) DESC;
```

*(Resultado a ser preenchido em produção)*

---

## Recomendações

1. **Monitorar uso dos índices** após 1 semana em produção
2. **Remover índices não utilizados** (`idx_scan = 0`)
3. **Executar `ANALYZE`** mensalmente ou após inserções em massa
4. **Considerar particionamento** de `despesas_cartao` se volume crescer muito (>100k registros)

---

## Referências

- [Plano técnico TT-09](../../docs/tecnico/06-plano-tt09-indices-performance.md)
- [Validação com EXPLAIN ANALYZE](../../docs/tasks/tt-09/02-validar-indices-explain.md)
- [PostgreSQL: Indexes](https://www.postgresql.org/docs/current/indexes.html)
- [Use The Index, Luke!](https://use-the-index-luke.com/)
