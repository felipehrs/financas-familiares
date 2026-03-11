# Análise: Estrutura de Dados para Despesas

**Data:** 2026-03-10
**Contexto:** O sistema possui 4 tabelas separadas de despesas (`despesas_cartao`, `assinaturas`, `contas_fixas`, `despesas_gerais`). O problema central é que **alterar o valor de uma assinatura ou conta fixa modifica retroativamente todos os meses passados e futuros**, pois não há separação entre o "cadastro" da despesa recorrente e sua "ocorrência" em uma competência.

---

## Problema Identificado

### Comportamento Atual

| Tipo | Como funciona hoje | Problema |
|---|---|---|
| `assinaturas` | Registro único com `valor`. O dashboard usa esse valor para qualquer mês. | Editar o valor da Netflix hoje altera o histórico de janeiro. |
| `contas_fixas` | Registro único com `valor`. Mesma lógica. | Editar a conta de luz altera todos os meses passados. |
| `despesas_cartao` | Cada parcela é um registro separado com `valor_parcela` fixo. | **Sem problema** — já é imutável por design. |
| `despesas_gerais` | Cada despesa é um registro pontual com data. | **Sem problema** — já é imutável por design. |

### Raiz do Problema

Despesas recorrentes (`assinaturas`, `contas_fixas`) são modeladas como **cadastros mutáveis**, não como **ocorrências históricas**. O sistema não distingue:
- O *contrato* da assinatura (Netflix existe desde jan/2024)
- O *lançamento* da assinatura em fev/2026 com o valor vigente naquele mês

---

## Solução 1 — Vigência com SCD Tipo 2 (Slowly Changing Dimension)

### Conceito

Inspirada na migration já existente para `rendas_fixas` (`data_inicio`, `data_fim`), esta solução adiciona **período de vigência** a `assinaturas` e `contas_fixas`. Ao editar o valor, não se atualiza o registro — cria-se um novo com `data_inicio` = data da alteração e encerra-se o anterior com `data_fim`.

### Mudanças no Banco

```sql
-- Assinaturas: adicionar vigência
ALTER TABLE assinaturas
    ADD COLUMN data_inicio DATE NOT NULL DEFAULT CURRENT_DATE,
    ADD COLUMN data_fim    DATE;

-- Contas fixas: mesma abordagem (substitui o campo `ativo` simples)
ALTER TABLE contas_fixas
    ADD COLUMN data_inicio DATE NOT NULL DEFAULT CURRENT_DATE,
    ADD COLUMN data_fim    DATE;
```

### Fluxo de Edição de Valor

```
Hoje (2026-03-10): Netflix custa R$ 45,90 (registro A, data_inicio=2024-01-01, data_fim=NULL)

Usuário edita para R$ 55,90:
  1. UPDATE assinaturas SET data_fim = '2026-03-10' WHERE id = A
  2. INSERT INTO assinaturas (..., valor=55.90, data_inicio='2026-03-11', data_fim=NULL)
```

### Query para um Mês Específico

```sql
-- Assinaturas vigentes em fev/2026
SELECT * FROM assinaturas
WHERE status = 'ativa'
  AND data_inicio <= '2026-02-28'
  AND (data_fim IS NULL OR data_fim >= '2026-02-01')
  AND deleted_at IS NULL;
```

### Prós e Contras

| Prós | Contras |
|---|---|
| Mínima reestruturação — mesmas tabelas | Queries ficam mais complexas (filtro de vigência em todo lugar) |
| Histórico completo de valores | Risco de inconsistência se a lógica de encerramento falhar |
| Mesma abordagem já usada em `rendas_fixas` (consistência) | UX precisa comunicar ao usuário que a alteração é "a partir de hoje" |
| Sem duplicação de dados | Dificuldade em fazer ajustes retroativos (ex: "o valor correto desde janeiro era X") |
| Projeções futuras usam o registro com `data_fim IS NULL` (simples) | |

### Esforço de Implementação

- **Baixo a médio** — 2 migrations, refatorar `assinatura_repository.go` e `conta_fixa_repository.go`, ajustar o dashboard.
- Reutiliza o padrão já implementado em `renda_fixa`.

---

## Solução 2 — Tabela de Lançamentos (Ledger/Journal)

### Conceito

Separar completamente o **cadastro** da despesa recorrente (o "contrato") da sua **ocorrência mensal** (o "lançamento"). Cria-se uma tabela `lancamentos_despesas` que registra o valor real cobrado em cada competência. O cadastro (`assinaturas`, `contas_fixas`) torna-se apenas um template de referência.

### Nova Estrutura

```sql
CREATE TABLE lancamentos_despesas (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tipo            VARCHAR(30) NOT NULL, -- 'assinatura' | 'conta_fixa' | 'despesa_geral' | 'parcela_cartao'
    referencia_id   UUID NOT NULL,        -- FK para a tabela de origem (assinatura, conta_fixa, etc.)
    mes_referencia  INTEGER NOT NULL CHECK (mes_referencia BETWEEN 1 AND 12),
    ano_referencia  INTEGER NOT NULL,
    valor           NUMERIC(15,2) NOT NULL, -- valor REAL cobrado naquele mês (snapshot)
    descricao       VARCHAR(500),
    categoria_id    UUID REFERENCES categorias(id),
    membro_id       UUID NOT NULL REFERENCES membros(id),
    data_efetiva    DATE,                 -- data real do débito/cobrança
    observacoes     TEXT,
    confirmado      BOOLEAN NOT NULL DEFAULT FALSE, -- lançamento confirmado vs. projetado
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    UNIQUE (tipo, referencia_id, mes_referencia, ano_referencia)
);
```

### Fluxo de Operação

```
Cadastro: usuário cria assinatura Netflix = R$ 45,90 (apenas o template)

Todo mês, o sistema (ou o usuário) "confirma" o lançamento:
  → INSERT INTO lancamentos_despesas (tipo='assinatura', referencia_id=<netflix_id>,
      mes_referencia=3, ano_referencia=2026, valor=45.90, confirmado=true)

Se o valor foi diferente em março (ex: cobrança em dólar):
  → O lançamento pode ter valor=48.30 sem alterar o cadastro da assinatura

Se a assinatura mudar de valor a partir de abril:
  → Atualiza apenas o template (assinaturas.valor = 55.90)
  → Lançamentos passados permanecem com os valores históricos (intocados)
```

### Modos de Geração de Lançamentos

1. **Auto-geração**: ao abrir o dashboard de um mês, gera automaticamente os lançamentos das recorrentes com `confirmado=false` (projeção).
2. **Confirmação manual**: usuário confirma/ajusta o valor real ao final do mês.
3. **Retroativo**: possível criar/editar lançamentos de meses passados sem afetar outros.

### Prós e Contras

| Prós | Contras |
|---|---|
| **Isolamento total** entre histórico e configuração atual | Maior reestruturação — novo modelo de dados, novo repositório |
| Valor exato de cada mês preservado para sempre | UX mais complexa: usuário precisa "confirmar" lançamentos |
| Suporte natural a ajustes retroativos pontuais | Risco de lançamentos não confirmados em meses esquecidos |
| Dashboard sempre lê de uma única tabela (mais rápido e simples) | Migração dos dados existentes requer script cuidadoso |
| Fundação para funcionalidades futuras (conciliação bancária, alertas de cobrança) | |
| `despesas_gerais` e `despesas_cartao` podem migrar para o mesmo modelo | |

### Esforço de Implementação

- **Alto** — nova tabela central, novo repositório, refatorar dashboard, lógica de auto-geração/confirmação, migração dos dados históricos.
- Representa uma mudança arquitetural significativa.

---

## Comparativo Direto

| Critério | Solução 1 (Vigência SCD) | Solução 2 (Lançamentos) |
|---|---|---|
| **Integridade histórica** | Parcial — valor correto por período, não por mês exato | Total — valor exato capturado por competência |
| **Ajuste retroativo pontual** | Difícil (precisaria alterar os períodos) | Simples (editar o lançamento do mês) |
| **Complexidade de implementação** | Baixa a média | Alta |
| **Complexidade de queries** | Média (filtro de vigência em todo lugar) | Baixa (uma tabela centralizada) |
| **Performance do dashboard** | Mantém múltiplos JOINs | Potencial para uma única query |
| **UX para o usuário** | Transparente (igual ao atual) | Requer ação de confirmação mensal |
| **Projeções futuras** | Usa o registro vigente atual | Usa o template + gera lançamentos futuros automaticamente |
| **Risco de migração** | Baixo | Médio |

---

## Recomendação

### Curto prazo (Sprint atual ou próxima): **Solução 1**

A Solução 1 resolve o problema de imutabilidade histórica com risco mínimo, seguindo o padrão já estabelecido em `rendas_fixas`. Entrega o valor imediato sem paralisar o roadmap atual (Sprint 9 de segurança está pendente).

**Passos sugeridos:**
1. Migration para adicionar `data_inicio`/`data_fim` em `assinaturas` e `contas_fixas`.
2. Refatorar repositories para filtrar por vigência.
3. Ajustar o service de edição para criar novo registro em vez de atualizar o valor.

### Médio/longo prazo: **Solução 2 (ou híbrido)**

Quando o produto estiver maduro e houver demanda por conciliação real vs. projetado (ex: "a Netflix cobrou R$ 48 em março porque estava em dólar"), a Solução 2 entrega o modelo correto. Pode ser implementada de forma incremental:
1. Primeiro para `assinaturas` e `contas_fixas`.
2. Depois `despesas_gerais` pode ser absorvida como um tipo de lançamento pontual.
3. As parcelas de cartão (`despesas_cartao`) já são imutáveis por design e podem se manter separadas ou migrar também.

### O que NÃO fazer

- **Não ignorar o problema**: ao escalar para múltiplos usuários (Sprint 9+), relatórios financeiros com valores históricos incorretos corroem a confiança no produto.
- **Não misturar as duas soluções** sem uma estratégia clara: adicionar vigência E lançamentos duplica a complexidade.
