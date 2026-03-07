CREATE TABLE despesas_gerais (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    membro_id UUID NOT NULL REFERENCES membros(id),
    categoria_id UUID REFERENCES categorias(id),
    descricao VARCHAR(500) NOT NULL,
    data DATE NOT NULL,
    valor NUMERIC(15,2) NOT NULL,
    forma_pagamento VARCHAR(50) NOT NULL CHECK (forma_pagamento IN ('dinheiro', 'debito', 'pix')),
    observacoes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
