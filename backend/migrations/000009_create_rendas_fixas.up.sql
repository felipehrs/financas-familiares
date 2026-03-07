CREATE TABLE rendas_fixas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    descricao VARCHAR(255) NOT NULL,
    membro_id UUID NOT NULL REFERENCES membros(id),
    valor NUMERIC(15,2) NOT NULL,
    dia_recebimento INTEGER NOT NULL CHECK (dia_recebimento BETWEEN 1 AND 31),
    ativa BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
