# Plano de Implementação: Tabela de Lançamentos (Ledger) + Tela de Pendências

**Data:** 10/03/2026
**Pré-requisito:** Sprint 9 concluída (TT-07 — `familia_id` em todas as tabelas + tabelas `familias` e `familia_usuarios`)

---

## 1. Análise de Impacto

### 1.1 O Problema Atual

`assinaturas` e `contas_fixas` são cadastros mutáveis usados diretamente pelo dashboard como se fossem lançamentos do mês. O campo `valor` representa tanto o valor "do template" quanto o valor "cobrado este mês" — são a mesma coisa. Isso cria dois problemas:

1. **Corrupção de histórico:** editar o valor da Netflix hoje retroage para todos os meses passados no dashboard/evolução mensal (o método `ResumoMensal` chama `ListarAtivas()` e soma `a.Valor` sem nenhum filtro de mês).
2. **Ausência de confirmação:** o sistema assume que toda assinatura ativa foi cobrada todo mês, sem que o usuário confirme o lançamento real.

### 1.2 O Que Quebra com a Mudança

| Componente | Impacto | Gravidade |
|---|---|---|
| `DashboardService.ResumoMensal` | Precisa buscar `lancamentos_despesas` ao invés de `ListarAtivas()` para meses passados/atual | Alta |
| `DashboardService.ProjecaoProximosMeses` | Continua usando o valor do template (correto para projeção futura) | Nenhum |
| `DashboardService.EvolucaoMensal` | Chama `ResumoMensal` para meses passados — agora usará lançamentos reais | Alta |
| `DashboardRepository.DespesasPorCategoria` | A UNION ALL que pega `assinaturas` e `contas_fixas` precisa virar JOIN em `lancamentos_despesas` | Alta |
| `AssinaturaRepository.ListarAtivas` | Permanece, mas não é mais usada pelo dashboard para meses passados | Baixa |
| `ContaFixaRepository.ListarAtivas` | Idem | Baixa |
| Frontend `AssinaturasPage` | Sem mudança no CRUD; ganha botão "Lançar este mês" | Média |
| Frontend `ContasFixasPage` | Idem | Média |
| Frontend `DashboardPage` | Sem mudança de tela, mas dados virão de fontes diferentes | Baixa |

### 1.3 O Que NÃO Muda

- CRUD de `assinaturas` e `contas_fixas` (create/update/delete continuam iguais)
- Projeção de meses futuros (continua usando valor do template via `ListarAtivas`)
- Despesas de cartão (`despesas_cartao` já é imutável por design)
- Despesas gerais (já são pontuais)
- Todo o módulo de rendas

---

## 2. Modelo de Dados

### 2.1 Nova Tabela: `lancamentos_despesas`

```sql
CREATE TABLE lancamentos_despesas (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    familia_id    UUID NOT NULL REFERENCES familias(id),

    -- Identificação do tipo e template de origem
    tipo          VARCHAR(20) NOT NULL
                  CHECK (tipo IN ('assinatura', 'conta_fixa')),
    referencia_id UUID NOT NULL,
    -- referencia_id aponta para assinaturas.id ou contas_fixas.id
    -- FK omitida intencionalmente: tipo é polimórfico e template pode ser
    -- soft-deleted depois do lançamento existir

    -- Competência: mês/ano a que se refere este lançamento
    competencia_mes INTEGER NOT NULL CHECK (competencia_mes BETWEEN 1 AND 12),
    competencia_ano INTEGER NOT NULL CHECK (competencia_ano >= 2020),

    -- Valor real cobrado (pode diferir do template, ex: dólar variou)
    valor         NUMERIC(15,2) NOT NULL CHECK (valor > 0),

    -- Categoria desnormalizada do template no momento do lançamento
    categoria_id  UUID REFERENCES categorias(id),

    -- Marcação de "não se aplica este mês" (assinatura pausada temporariamente)
    nao_se_aplica BOOLEAN NOT NULL DEFAULT FALSE,

    -- Metadados
    observacao    VARCHAR(500),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,

    -- Garante unicidade: um único lançamento por template por competência
    CONSTRAINT uq_lancamento_competencia
        UNIQUE (referencia_id, competencia_mes, competencia_ano)
);
```

> **Por que `categoria_id` desnormalizado?** Simplifica a query `DespesasPorCategoria` no dashboard — sem CASE WHEN para resolver categoria via JOIN polimórfico na tabela de origem. O valor é copiado do template no momento do lançamento.

### 2.2 Índices

```sql
CREATE INDEX idx_lancamentos_familia
    ON lancamentos_despesas(familia_id)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_lancamentos_competencia
    ON lancamentos_despesas(familia_id, competencia_ano, competencia_mes)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_lancamentos_referencia
    ON lancamentos_despesas(referencia_id, competencia_mes, competencia_ano)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_lancamentos_deleted_at
    ON lancamentos_despesas(deleted_at);
```

### 2.3 Domínio Go

```go
// backend/internal/domain/lancamento_despesa.go

type LancamentoDespesa struct {
    ID             string
    FamiliaID      string
    Tipo           string  // "assinatura" | "conta_fixa"
    ReferenciaID   string
    CompetenciaMes int
    CompetenciaAno int
    Valor          float64
    CategoriaID    *string
    NaoSeAplica    bool
    Observacao     *string
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

var (
    ErrLancamentoDespesaNaoEncontrado        = errors.New("lançamento não encontrado")
    ErrLancamentoDespesaJaExiste             = errors.New("já existe lançamento para esta competência")
    ErrLancamentoDespesaTipoInvalido         = errors.New("tipo deve ser: assinatura ou conta_fixa")
    ErrLancamentoDespesaValorInvalido        = errors.New("valor deve ser maior que zero")
    ErrLancamentoDespesaCompetenciaInvalida  = errors.New("competência inválida")
    ErrLancamentoDespesaReferenciaObrigatoria = errors.New("referencia_id é obrigatório")
)
```

---

## 3. Estratégia de Migração de Dados

### 3.1 Sem Migração Retroativa

**Decisão:** não gerar lançamentos automáticos para meses passados. O sistema adota uma **data de corte** (`data_implantacao_ledger`).

- **Antes da data de corte:** dashboard usa o valor atual do template (comportamento atual — aceito como trade-off)
- **A partir da data de corte:** dashboard usa apenas lançamentos confirmados

**Implementação da data de corte:** constante no código ou variável de ambiente `LEDGER_DATA_CORTE=2026-03` definida no Railway e Vercel.

### 3.2 Migration SQL

**Número:** `000018_create_lancamentos_despesas` (000017 é reservada para a migration de `familia_id` do TT-07)

- `up.sql`: DDL completo da seção 2.1 + índices da seção 2.2
- `down.sql`: `DROP TABLE IF EXISTS lancamentos_despesas;`

Nenhuma alteração nas tabelas existentes.

---

## 4. Fases de Implementação

### Fase A — Backend: Infraestrutura do Ledger

**Objetivo:** Migration, domain, repository, service, handler e rotas. Refatorar o dashboard para usar lançamentos reais.

#### Tasks de Backend

| ID | Tarefa |
|---|---|
| **A-B01** | Migration `000017`: criar tabela + índices |
| **A-B02** | `domain/lancamento_despesa.go`: struct, constantes de tipo, erros de domínio |
| **A-B03** | `repository/lancamento_despesa_repository.go`: `Criar`, `BuscarPorID`, `BuscarPorReferenciaECompetencia`, `ListarPorCompetencia`, `ListarPorReferencia`, `Atualizar`, `MarcarNaoSeAplica`, `Deletar` |
| **A-B04** | `service/lancamento_despesa_service.go`: validações, verificação de duplicidade, interface de repositório |
| **A-B05** | `handler/lancamento_despesa_handler.go`: endpoints CRUD + marcar não se aplica, mapeamento de erros HTTP |
| **A-B06** | Registrar rotas no servidor |
| **A-B07** | Refatorar `DashboardService.ResumoMensal`: lógica de data de corte — usa lançamentos para meses ≥ corte, template para meses anteriores |
| **A-B08** | Refatorar `DashboardRepository.DespesasPorCategoria`: UNION ALL passa a incluir `lancamentos_despesas` (onde `nao_se_aplica=FALSE`) em vez de `assinaturas`/`contas_fixas` direto |
| **A-B09** | `service/lancamento_despesa_service_test.go`: criação, duplicidade, nao_se_aplica, validações |
| **A-B10** | `handler/lancamento_despesa_handler_test.go`: todos os endpoints |
| **A-B11** | Atualizar `dashboard_service_test.go`: novos mocks de `LancamentoDespesaRepository`, cenários com/sem lançamentos |

#### Critério de Pronto (Fase A)

- Migration aplicada localmente e na produção sem erros
- `make test` verde, cobertura ≥ 90% nas novas camadas
- `GET /api/v1/lancamentos-despesas?mes=3&ano=2026` retorna lista correta
- `POST /api/v1/lancamentos-despesas` cria lançamento, responde 201
- `POST` com mesma `(referencia_id, mes, ano)` retorna 409
- Dashboard usa lançamentos reais para meses ≥ data de corte; não quebra para meses anteriores

---

### Fase B — Frontend: Lançamento nas Páginas Existentes

**Objetivo:** Adicionar em `AssinaturasPage` e `ContasFixasPage` a capacidade de lançar o mês atual para cada item, com badge de status e modal de confirmação.

#### Tasks de Backend

| ID | Tarefa |
|---|---|
| **B-B01** | `GET /api/v1/lancamentos-despesas/pendencias?mes=M&ano=A`: retorna assinaturas/contas sem lançamento + faturas de cartão do mês (ver estrutura abaixo) |

**Estrutura de resposta do endpoint de pendências:**

```json
{
  "mes": 3,
  "ano": 2026,
  "assinaturas_pendentes": [
    {
      "assinatura_id": "uuid",
      "nome": "Netflix",
      "valor_template": 48.90,
      "dia_cobranca": 15,
      "membro_nome": "Felipe"
    }
  ],
  "contas_fixas_pendentes": [
    {
      "conta_fixa_id": "uuid",
      "descricao": "Conta de Luz",
      "valor_template": 180.00,
      "dia_vencimento": 10,
      "membro_nome": "Ana"
    }
  ],
  "faturas_cartao": [
    {
      "cartao_id": "uuid",
      "cartao_nome": "Nubank",
      "dia_vencimento": 20,
      "total_fatura": 1250.00,
      "num_parcelas": 8
    }
  ]
}
```

**Regras de negócio do endpoint:**
1. Busca `assinaturas` com `status='ativa'` do usuário → verifica se existe lançamento para `(id, mes, ano)` → inclui as que não têm
2. Repete para `contas_fixas` com `ativa=true`
3. Para faturas: agrupa `despesas_cartao` por `cartao_id` onde `fatura_mes=mes AND fatura_ano=ano`, join com `cartoes_credito` para `dia_vencimento`

#### Tasks de Frontend

| ID | Tarefa |
|---|---|
| **B-F01** | `src/types/lancamento_despesa.ts`: types `LancamentoDespesa`, `CriarLancamentoRequest`, `PendenciasMes` |
| **B-F02** | `src/api/lancamentos_despesas.ts`: `criarLancamento`, `listarPorCompetencia`, `atualizarLancamento`, `marcarNaoSeAplica`, `deletarLancamento`, `buscarPendencias` |
| **B-F03** | `AssinaturasPage`: coluna/badge mostrando status do lançamento do mês atual; botão "Lançar" |
| **B-F04** | `ContasFixasPage`: idem |
| **B-F05** | `components/ModalLancamento.tsx`: modal reutilizável (assinatura ou conta fixa), campos valor + observação, checkbox "não se aplica este mês" |
| **B-F06** | Testes: `AssinaturasPage.test.tsx`, `ContasFixasPage.test.tsx`, `ModalLancamento.test.tsx` |

#### Critério de Pronto (Fase B)

- Badge de status visível em cada assinatura/conta fixa (lançado / pendente / não se aplica)
- Fluxo completo: abrir modal → preencher valor → salvar → badge muda
- Fluxo "não se aplica": checkbox → confirmar → badge muda
- Dashboard atualiza ao retornar após lançar
- `pnpm test:run` verde

---

### Fase C — Tela de Pendências

**Objetivo:** Nova página `/pendencias` com checklist mensal completo: assinaturas pendentes, contas fixas pendentes e faturas de cartão informativas.

#### Tasks de Frontend

| ID | Tarefa |
|---|---|
| **C-F01** | `pages/PendenciasPage.tsx`: página principal com seletor de mês/ano |
| **C-F02** | `components/CardPendenciaAssinatura.tsx`: card com nome, valor esperado, dia de cobrança, membro, botões "Não se aplica" e "Confirmar" |
| **C-F03** | `components/CardPendenciaContaFixa.tsx`: idem para contas fixas |
| **C-F04** | `components/CardPendenciaFatura.tsx`: card read-only com cartão, total da fatura, dia de vencimento, badge "Automático" |
| **C-F05** | `components/SecaoPendencias.tsx`: seção genérica com título, badge de contagem, lista de cards |
| **C-F06** | Estado de "tudo em dia": mensagem de parabéns + ícone check quando lista vazia |
| **C-F07** | Rota `/pendencias` em `App.tsx` |
| **C-F08** | Entrada "Pendências" em `AppLayout.tsx` com badge de contagem (mês atual) |
| **C-F09** | `pages/PendenciasPage.test.tsx`: estados carregando, com pendências, sem pendências, erro |

#### Critério de Pronto (Fase C)

- Tela `/pendencias` acessível pelo menu com badge de contagem
- Confirmar assinatura: modal com valor pré-preenchido → salvar → item some da lista
- Marcar "não se aplica": item some da lista (não contabilizado no dashboard)
- Faturas de cartão listadas como informativas (sem ação de lançar)
- Lista vazia: "Tudo em dia para março de 2026 ✓"
- Navegação entre meses funciona
- `pnpm test:run` verde

---

## 5. Endpoints REST

```
POST   /api/v1/lancamentos-despesas
GET    /api/v1/lancamentos-despesas?mes=3&ano=2026
GET    /api/v1/lancamentos-despesas/:id
PUT    /api/v1/lancamentos-despesas/:id
DELETE /api/v1/lancamentos-despesas/:id
PATCH  /api/v1/lancamentos-despesas/:id/nao-se-aplica
GET    /api/v1/lancamentos-despesas/pendencias?mes=3&ano=2026
```

**POST — Body:**
```json
{
  "tipo": "assinatura",
  "referencia_id": "uuid-da-assinatura",
  "competencia_mes": 3,
  "competencia_ano": 2026,
  "valor": 48.90,
  "observacao": "cobrança em dólar — cotação 5.78"
}
```

**Erros específicos:**
- `409 Conflict` — já existe lançamento para `(referencia_id, mes, ano)`
- `404 Not Found` — `referencia_id` não existe ou não pertence ao usuário
- `400 Bad Request` — validações de campo

**PATCH /:id/nao-se-aplica — Body:**
```json
{ "nao_se_aplica": true }
```

---

## 6. Layout da Tela de Pendências

```
┌─────────────────────────────────────────────┐
│  Pendências de março 2026          [◀] [▶]  │
│  3 itens pendentes                           │
├─────────────────────────────────────────────┤
│  ASSINATURAS (2 pendentes)                  │
│  ┌──────────────────────────────────────┐   │
│  │ Netflix                              │   │
│  │ Valor esperado: R$ 48,90 · Dia 15   │   │
│  │ Responsável: Felipe                  │   │
│  │         [Não se aplica] [Confirmar]  │   │
│  └──────────────────────────────────────┘   │
│  ┌──────────────────────────────────────┐   │
│  │ Spotify                              │   │
│  │ Valor esperado: R$ 21,90 · Dia 20   │   │
│  │ Responsável: Ana                     │   │
│  │         [Não se aplica] [Confirmar]  │   │
│  └──────────────────────────────────────┘   │
├─────────────────────────────────────────────┤
│  CONTAS FIXAS (1 pendente)                  │
│  ┌──────────────────────────────────────┐   │
│  │ Conta de Luz                         │   │
│  │ Valor esperado: R$ 180,00 · Dia 10  │   │
│  │         [Não se aplica] [Confirmar]  │   │
│  └──────────────────────────────────────┘   │
├─────────────────────────────────────────────┤
│  FATURAS DE CARTÃO (informativo)            │
│  ┌──────────────────────────────────────┐   │
│  │ Nubank · Vence dia 20               │   │
│  │ R$ 1.250,00 (8 parcelas)            │   │
│  │                        ✓ Automático  │   │
│  └──────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
```

**Modal de Confirmação:**

```
┌─────────────────────────────────┐
│  Lançar Netflix — março 2026    │
│                                 │
│  Valor cobrado (R$)             │
│  [   48,90                  ]   │
│                                 │
│  Observação (opcional)          │
│  [                          ]   │
│                                 │
│  [ ] Não se aplica este mês     │
│                                 │
│         [Cancelar] [Confirmar]  │
└─────────────────────────────────┘
```

---

## 7. Mudanças no DashboardService

### 7.1 Lógica de Fonte de Dados por Mês

```go
// ResumoMensal(mes, ano):
temLancamentos, _ := lancamentosRepo.ExisteAlgumLancamento(mes, ano, familiaID)
ehMesApartirDaCorte := competenciaApartirDaCorte(mes, ano) // env LEDGER_DATA_CORTE

if temLancamentos || ehMesApartirDaCorte {
    // Usa lançamentos reais (nao_se_aplica=false)
    lancamentos, _ := lancamentosRepo.ListarPorCompetencia(mes, ano, familiaID)
    totalAssinaturas = somarPorTipo(lancamentos, "assinatura")
    totalContasFixas = somarPorTipo(lancamentos, "conta_fixa")
} else {
    // Compatibilidade retroativa: usa valor do template
    assinaturas, _ = assinaturaRepo.ListarAtivas(familiaID)
    contasFixas, _ = contaFixaRepo.ListarAtivas(familiaID)
    // soma como hoje
}
```

### 7.2 Query `DespesasPorCategoria` Atualizada

A UNION ALL existente substitui os blocos de `assinaturas`/`contas_fixas` por:

```sql
UNION ALL
SELECT ld.categoria_id, SUM(ld.valor) as total
FROM lancamentos_despesas ld
WHERE ld.competencia_mes = $1
  AND ld.competencia_ano = $2
  AND ld.familia_id      = $3
  AND ld.nao_se_aplica   = FALSE
  AND ld.deleted_at IS NULL
GROUP BY ld.categoria_id
```

### 7.3 Aviso no Dashboard Quando Há Pendências

Quando o mês visualizado tem pendências (assinaturas/contas sem lançamento), o dashboard exibe um banner:

> ⚠️ Este mês tem 3 despesas pendentes de lançamento. Os totais abaixo podem estar incompletos.

O banner inclui link para a tela de pendências.

---

## 8. Sequenciamento e Dependências

```
Sprint 9 (TT-07 — familia_id em todas as tabelas + tabelas familias/familia_usuarios)
        │
        ▼
   Fase A — Backend Ledger
   (migration + domain + repo + service + handler + refatorar dashboard)
        │
        ├──► Fase B — Lançamento nas Páginas Existentes
        │    (AssinaturasPage + ContasFixasPage + ModalLancamento)
        │
        └──► Fase C — Tela de Pendências
             (PendenciasPage + cards + menu badge)
```

Fase B e Fase C podem rodar **em paralelo** após a Fase A estar completa.

---

## 9. Riscos e Mitigações

| Risco | Probabilidade | Impacto | Mitigação |
|---|---|---|---|
| Meses sem lançamento mostram R$ 0 | Alta | Alto | Banner contextual no dashboard: "X pendências não lançadas — valores incompletos" |
| Usuário esquecer de lançar | Média | Alto | Badge de contagem no menu; banner no dashboard |
| Evolução mensal incorreta antes da data de corte | Alta | Médio | Trade-off documentado; exibir nota na UI para meses históricos |
| Constraint UNIQUE rejeitar tentativa dupla | Baixa | Baixo | Backend 409 com mensagem clara; frontend verifica existência antes do POST |
| Template soft-deleted com lançamentos existentes | Baixa | Baixo | FK polimórfica omitida; lançamentos sobrevivem ao delete do template via JOIN LEFT |
| Sprint 9 não concluída | Alta (pré-req) | Crítico | Bloquear Fase A até TT-07 estar em ✅ (tabela `familias` precisa existir para a FK de `lancamentos_despesas`) |

---

## 10. Arquivos por Fase

### Fase A — Novos

```
backend/migrations/000018_create_lancamentos_despesas.up.sql
backend/migrations/000018_create_lancamentos_despesas.down.sql
backend/internal/domain/lancamento_despesa.go
backend/internal/repository/lancamento_despesa_repository.go
backend/internal/service/lancamento_despesa_service.go
backend/internal/service/lancamento_despesa_service_test.go
backend/internal/handler/lancamento_despesa_handler.go
backend/internal/handler/lancamento_despesa_handler_test.go
```

### Fase A — Modificados

```
backend/internal/service/dashboard_service.go
backend/internal/service/dashboard_service_test.go
backend/internal/repository/dashboard_repository.go
backend/server.go  (ou main.go — registro de rotas)
```

### Fase B — Novos

```
frontend/src/types/lancamento_despesa.ts
frontend/src/api/lancamentos_despesas.ts
frontend/src/components/ModalLancamento.tsx
frontend/src/components/ModalLancamento.test.tsx
```

### Fase B — Modificados

```
frontend/src/pages/AssinaturasPage.tsx
frontend/src/pages/AssinaturasPage.test.tsx
frontend/src/pages/ContasFixasPage.tsx
frontend/src/pages/ContasFixasPage.test.tsx
```

### Fase C — Novos

```
frontend/src/pages/PendenciasPage.tsx
frontend/src/pages/PendenciasPage.test.tsx
frontend/src/components/CardPendenciaAssinatura.tsx
frontend/src/components/CardPendenciaContaFixa.tsx
frontend/src/components/CardPendenciaFatura.tsx
frontend/src/components/SecaoPendencias.tsx
```

### Fase C — Modificados

```
frontend/src/App.tsx
frontend/src/components/AppLayout.tsx
```
