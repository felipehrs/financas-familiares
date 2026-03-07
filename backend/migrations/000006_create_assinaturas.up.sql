CREATE TABLE assinaturas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nome VARCHAR(255) NOT NULL,
    membro_id UUID NOT NULL REFERENCES membros(id),
    categoria_id UUID REFERENCES categorias(id),
    valor NUMERIC(15,2) NOT NULL,
    dia_cobranca INTEGER NOT NULL CHECK (dia_cobranca BETWEEN 1 AND 31),
    forma_pagamento VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ativa' CHECK (status IN ('ativa', 'pausada', 'cancelada')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
