CREATE TABLE rendas_variaveis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    descricao VARCHAR(255) NOT NULL,
    membro_id UUID NOT NULL REFERENCES membros(id),
    mes_referencia INTEGER NOT NULL CHECK (mes_referencia BETWEEN 1 AND 12),
    ano_referencia INTEGER NOT NULL,
    valor NUMERIC(15,2) NOT NULL,
    data_recebimento DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
