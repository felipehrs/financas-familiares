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
