ALTER TABLE despesas_cartao
    DROP COLUMN IF EXISTS parcela_numero,
    DROP COLUMN IF EXISTS compra_id;
