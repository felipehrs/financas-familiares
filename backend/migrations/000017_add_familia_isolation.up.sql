CREATE TABLE IF NOT EXISTS familias (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nome       VARCHAR(100) NOT NULL,
    owner_id   UUID NOT NULL REFERENCES usuarios(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS familia_usuarios (
    familia_id UUID NOT NULL REFERENCES familias(id),
    usuario_id UUID NOT NULL REFERENCES usuarios(id),
    role       VARCHAR(20) NOT NULL DEFAULT 'membro'
                CHECK (role IN ('admin', 'membro')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (familia_id, usuario_id)
);

DO $$
DECLARE
    v_usuario_id UUID;
    v_familia_id UUID;
BEGIN
    SELECT id INTO v_usuario_id FROM usuarios ORDER BY created_at ASC LIMIT 1;
    IF v_usuario_id IS NOT NULL THEN
        INSERT INTO familias (nome, owner_id)
        VALUES ('Família Principal', v_usuario_id)
        RETURNING id INTO v_familia_id;

        INSERT INTO familia_usuarios (familia_id, usuario_id, role)
        SELECT v_familia_id, id, CASE WHEN id = v_usuario_id THEN 'admin' ELSE 'membro' END
        FROM usuarios
        WHERE deleted_at IS NULL
        ON CONFLICT DO NOTHING;
    END IF;
END $$;

ALTER TABLE membros ADD COLUMN IF NOT EXISTS familia_id UUID;
UPDATE membros SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1) WHERE familia_id IS NULL;
ALTER TABLE membros ALTER COLUMN familia_id SET NOT NULL, ADD CONSTRAINT membros_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

ALTER TABLE categorias ADD COLUMN IF NOT EXISTS familia_id UUID;
UPDATE categorias SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1) WHERE familia_id IS NULL;
ALTER TABLE categorias ALTER COLUMN familia_id SET NOT NULL, ADD CONSTRAINT categorias_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);
ALTER TABLE categorias DROP CONSTRAINT IF EXISTS categorias_nome_key;
ALTER TABLE categorias ADD CONSTRAINT categorias_nome_familia_unique UNIQUE (familia_id, nome);

ALTER TABLE cartoes_credito ADD COLUMN IF NOT EXISTS familia_id UUID;
UPDATE cartoes_credito SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1) WHERE familia_id IS NULL;
ALTER TABLE cartoes_credito ALTER COLUMN familia_id SET NOT NULL, ADD CONSTRAINT cartoes_credito_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

ALTER TABLE despesas_cartao ADD COLUMN IF NOT EXISTS familia_id UUID;
UPDATE despesas_cartao SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1) WHERE familia_id IS NULL;
ALTER TABLE despesas_cartao ALTER COLUMN familia_id SET NOT NULL, ADD CONSTRAINT despesas_cartao_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

ALTER TABLE assinaturas ADD COLUMN IF NOT EXISTS familia_id UUID;
UPDATE assinaturas SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1) WHERE familia_id IS NULL;
ALTER TABLE assinaturas ALTER COLUMN familia_id SET NOT NULL, ADD CONSTRAINT assinaturas_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

ALTER TABLE contas_fixas ADD COLUMN IF NOT EXISTS familia_id UUID;
UPDATE contas_fixas SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1) WHERE familia_id IS NULL;
ALTER TABLE contas_fixas ALTER COLUMN familia_id SET NOT NULL, ADD CONSTRAINT contas_fixas_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

ALTER TABLE despesas_gerais ADD COLUMN IF NOT EXISTS familia_id UUID;
UPDATE despesas_gerais SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1) WHERE familia_id IS NULL;
ALTER TABLE despesas_gerais ALTER COLUMN familia_id SET NOT NULL, ADD CONSTRAINT despesas_gerais_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

ALTER TABLE rendas_fixas ADD COLUMN IF NOT EXISTS familia_id UUID;
UPDATE rendas_fixas SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1) WHERE familia_id IS NULL;
ALTER TABLE rendas_fixas ALTER COLUMN familia_id SET NOT NULL, ADD CONSTRAINT rendas_fixas_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

ALTER TABLE rendas_variaveis ADD COLUMN IF NOT EXISTS familia_id UUID;
UPDATE rendas_variaveis SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1) WHERE familia_id IS NULL;
ALTER TABLE rendas_variaveis ALTER COLUMN familia_id SET NOT NULL, ADD CONSTRAINT rendas_variaveis_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

ALTER TABLE rendas_extras ADD COLUMN IF NOT EXISTS familia_id UUID;
UPDATE rendas_extras SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1) WHERE familia_id IS NULL;
ALTER TABLE rendas_extras ALTER COLUMN familia_id SET NOT NULL, ADD CONSTRAINT rendas_extras_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);

ALTER TABLE rendimentos_investimento ADD COLUMN IF NOT EXISTS familia_id UUID;
UPDATE rendimentos_investimento SET familia_id = (SELECT id FROM familias ORDER BY created_at ASC LIMIT 1) WHERE familia_id IS NULL;
ALTER TABLE rendimentos_investimento ALTER COLUMN familia_id SET NOT NULL, ADD CONSTRAINT rendimentos_investimento_familia_id_fkey FOREIGN KEY (familia_id) REFERENCES familias(id);
