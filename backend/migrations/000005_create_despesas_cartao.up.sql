CREATE TABLE despesas_cartao (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cartao_id UUID NOT NULL REFERENCES cartoes_credito(id),
    categoria_id UUID REFERENCES categorias(id),
    descricao VARCHAR(500) NOT NULL,
    data_compra DATE NOT NULL,
    valor_total NUMERIC(15,2) NOT NULL,
    numero_parcelas INTEGER NOT NULL DEFAULT 1,
    valor_parcela NUMERIC(15,2) NOT NULL, -- valor_total / numero_parcelas
    fatura_mes INTEGER NOT NULL,           -- mês da primeira fatura (1-12)
    fatura_ano INTEGER NOT NULL,           -- ano da primeira fatura
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
