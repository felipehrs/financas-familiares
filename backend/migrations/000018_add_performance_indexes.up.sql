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
-- Ordem: familia_id, ano_referencia DESC, mes_referencia DESC (para queries de "últimos X meses")
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_rendas_variaveis_mes_ano
ON rendas_variaveis (familia_id, ano_referencia DESC, mes_referencia DESC)
WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_rendas_extras_mes_ano
ON rendas_extras (familia_id, EXTRACT(YEAR FROM data_recebimento)::INTEGER DESC, EXTRACT(MONTH FROM data_recebimento)::INTEGER DESC)
WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_rendimentos_investimento_mes_ano
ON rendimentos_investimento (familia_id, EXTRACT(YEAR FROM data)::INTEGER DESC, EXTRACT(MONTH FROM data)::INTEGER DESC)
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
