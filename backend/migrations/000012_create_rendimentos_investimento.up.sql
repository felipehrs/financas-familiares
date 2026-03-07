CREATE TABLE rendimentos_investimento (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    descricao VARCHAR(255) NOT NULL,
    membro_id UUID NOT NULL REFERENCES membros(id),
    data DATE NOT NULL,
    valor NUMERIC(15,2) NOT NULL,
    valor_distribuido NUMERIC(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT valor_distribuido_valido CHECK (valor_distribuido >= 0 AND valor_distribuido <= valor)
);
