-- Seed de dados iniciais para o sistema de finanças familiares
-- Nota: usuários devem ser criados via API (POST /api/v1/auth/register)
-- pois as senhas são processadas pelo backend com bcrypt.
-- Este arquivo é apenas para dados de referência que não dependem de lógica de aplicação.

-- Categorias padrão
INSERT INTO categorias (nome) VALUES
    ('Alimentação'),
    ('Transporte'),
    ('Lazer'),
    ('Saúde'),
    ('Educação'),
    ('Moradia'),
    ('Vestuário'),
    ('Outros')
ON CONFLICT DO NOTHING;
