-- Adiciona suporte a compras parceladas (US-08 / RN03).
-- compra_id agrupa todas as parcelas de uma mesma compra.
-- parcela_numero identifica a posição da parcela (1, 2, 3, ...).

ALTER TABLE despesas_cartao
    ADD COLUMN compra_id     UUID NOT NULL DEFAULT gen_random_uuid(),
    ADD COLUMN parcela_numero INTEGER NOT NULL DEFAULT 1;

-- Registros existentes (à vista) já têm numero_parcelas = 1, então
-- compra_id e parcela_numero 1 estão corretos por padrão.
