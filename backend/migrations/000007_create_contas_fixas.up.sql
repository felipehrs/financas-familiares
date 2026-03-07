CREATE TABLE contas_fixas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    descricao VARCHAR(255) NOT NULL,
    membro_id UUID NOT NULL REFERENCES membros(id),
    categoria_id UUID REFERENCES categorias(id),
    valor NUMERIC(15,2) NOT NULL,
    dia_vencimento INTEGER NOT NULL CHECK (dia_vencimento BETWEEN 1 AND 31),
    forma_pagamento VARCHAR(50) NOT NULL,
    ativa BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
