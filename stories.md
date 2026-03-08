# Histórias de Usuário — Finanças Familiares

**Referência:** spec.md v1.1 | tech-spec.md v1.0
**Atualizado em:** 06/03/2026

---

## Convenções

- **Formato:** Como [ator], quero [ação], para que [objetivo].
- **Ator padrão:** Usuário (membro autenticado da família)
- **Prioridade:** 🔴 MVP | 🟡 Core | 🟢 Analytics | ⚪ Extras
- **Critérios de aceite (AC):** derivados diretamente de RF/RN da spec.md

---

## FASE 1 — MVP

### EPIC: Autenticação

---

#### US-01 — Login no sistema 🔴

**Como** usuário, **quero** me autenticar com e-mail e senha, **para que** meus dados financeiros sejam protegidos.

**Critérios de aceite:**
- AC01: O sistema aceita e-mail e senha válidos e emite um JWT com expiração de 15 minutos + refresh token de 7 dias.
- AC02: Credenciais inválidas retornam mensagem de erro sem revelar qual campo está incorreto.
- AC03: Após expiração do token de acesso, o refresh token permite renovar a sessão sem novo login.
- AC04: O sistema suporta exatamente 2 usuários por instância (casal).

---

### EPIC: Membros da Família

---

#### US-02 — Cadastrar membro da família 🔴

**Como** usuário, **quero** cadastrar os membros da família, **para que** eu possa vincular rendas, despesas e cartões a cada pessoa.

**Critérios de aceite (RF01):**
- AC01: É possível criar um membro informando nome (obrigatório) e relacionamento (opcional).
- AC02: O membro é criado com status ativo por padrão.
- AC03: É possível editar nome, relacionamento e status de um membro existente.
- AC04: Não é possível excluir um membro que possui vínculos ativos (cartões, despesas, rendas). O sistema exibe mensagem explicativa.
- AC05: É possível inativar um membro, fazendo com que seus registros sejam excluídos dos cálculos futuros sem apagar o histórico.

---

### EPIC: Categorias

---

#### US-03 — Gerenciar categorias 🔴

**Como** usuário, **quero** criar e gerenciar categorias de despesas e rendas, **para que** eu possa organizar e analisar meus gastos por tipo.

**Critérios de aceite (RF02):**
- AC01: O sistema oferece categorias padrão na criação da conta: Alimentação, Transporte, Lazer, Saúde, Educação, Moradia, Vestuário, Outros.
- AC02: É possível adicionar novas categorias informando apenas o nome.
- AC03: É possível editar o nome de uma categoria existente.
- AC04: É possível excluir uma categoria que não possui registros vinculados.
- AC05: Tentativa de excluir categoria com vínculos exibe mensagem de erro e impede a exclusão.

---

### EPIC: Cartões de Crédito

---

#### US-04 — Cadastrar cartão de crédito 🔴

**Como** usuário, **quero** cadastrar meus cartões de crédito com dia de fechamento e vencimento, **para que** o sistema possa calcular automaticamente a qual fatura cada compra pertence.

**Critérios de aceite (RF03):**
- AC01: É possível cadastrar um cartão informando: nome/identificação, membro responsável, dia de fechamento e dia de vencimento.
- AC02: Limite do cartão é opcional.
- AC03: É possível cadastrar múltiplos cartões.
- AC04: É possível editar qualquer campo do cartão.
- AC05: Cartão pode ser inativado sem apagar o histórico de despesas.

---

#### US-05 — Lançar despesa no cartão (à vista) 🔴

**Como** usuário, **quero** registrar uma compra à vista no cartão de crédito, **para que** o valor apareça na fatura correta do mês.

**Critérios de aceite (RF04, RF05, RF06 — RN01, RN02):**
- AC01: É possível lançar uma despesa informando: data da compra, cartão, descrição, categoria, valor total. O número de parcelas padrão é 1.
- AC02: O sistema calcula automaticamente a fatura com base na data da compra e no dia de fechamento do cartão (RN01):
  - Data ≤ dia_fechamento → fatura do mês corrente.
  - Data > dia_fechamento → fatura do mês seguinte.
- AC03: A fatura é exibida no formato "MMM/AA" (ex: "MAR/26").
- AC04: O valor da parcela é calculado como valor_total / numero_parcelas (RN02).
- AC05: A despesa lançada aparece imediatamente na listagem de faturas do cartão.

---

### EPIC: Rendas

---

#### US-06 — Cadastrar renda fixa (salário) 🔴

**Como** usuário, **quero** cadastrar minha renda fixa mensal, **para que** ela seja incluída automaticamente nos cálculos de saldo de todos os meses.

**Critérios de aceite (RF14):**
- AC01: É possível cadastrar uma renda fixa informando: descrição, membro, valor mensal, dia de recebimento.
- AC02: Renda fixa ativa é incluída automaticamente no cálculo do saldo de qualquer mês (atual e projeções futuras).
- AC03: É possível inativar uma renda fixa sem excluir o histórico; renda inativa não entra nas projeções futuras.
- AC04: É possível editar valor, descrição e demais campos a qualquer momento.

---

### EPIC: Dashboard

---

#### US-07 — Visualizar resumo financeiro do mês 🔴

**Como** usuário, **quero** ver um resumo financeiro do mês atual, **para que** eu saiba rapidamente o saldo da família.

**Critérios de aceite (RF19 — RN06):**
- AC01: O dashboard exibe, para o mês selecionado:
  - Total de rendas fixas (salários ativos)
  - **TOTAL RENDAS OPERACIONAIS**
  - Fatura total dos cartões do mês
  - **TOTAL DE DESPESAS**
  - **SALDO DO MÊS = TOTAL RENDAS OPERACIONAIS − TOTAL DESPESAS**
- AC02: O saldo é exibido com indicação visual: verde (positivo), vermelho (negativo).
- AC03: É possível selecionar o mês/ano para visualizar dados históricos.
- AC04: É possível alternar entre visão consolidada (todos os membros) e visão individual por membro (RN08).

---

## FASE 2 — Core

### EPIC: Cartões de Crédito (avançado)

---

#### US-08 — Lançar compra parcelada no cartão 🟡

**Como** usuário, **quero** registrar uma compra parcelada no cartão, **para que** cada parcela apareça na fatura do mês correto automaticamente.

**Critérios de aceite (RF04, RF05, RF06 — RN01, RN02, RN03):**
- AC01: É possível informar o número de parcelas ao lançar uma despesa.
- AC02: O valor de cada parcela é calculado como valor_total / numero_parcelas (RN02).
- AC03: A parcela 1 é alocada na fatura calculada pela regra RN01.
- AC04: Cada parcela subsequente é alocada no mês seguinte à parcela anterior (RN03).
- AC05: Todas as parcelas futuras aparecem na projeção dos meses correspondentes.

---

#### US-09 — Visualizar fatura de cartão por mês 🟡

**Como** usuário, **quero** ver todas as despesas de um cartão filtradas por fatura (mês/ano), **para que** eu saiba o total exato que vou pagar naquela fatura.

**Critérios de aceite (RF06, RF07):**
- AC01: É possível filtrar despesas por cartão e por fatura (mês/ano).
- AC02: O total da fatura selecionada é exibido.
- AC03: Cada linha da listagem exibe: data da compra, descrição, valor da parcela, indicação "X/Y" (parcela atual/total).
- AC04: O saldo disponível do cartão é exibido quando o limite está configurado.
- AC05: Alerta visual é exibido quando o uso ultrapassa 80% do limite configurado (RF07).

---

### EPIC: Assinaturas

---

#### US-10 — Cadastrar assinatura recorrente 🟡

**Como** usuário, **quero** cadastrar assinaturas de serviços, **para que** elas sejam incluídas automaticamente nos cálculos mensais sem precisar lançar todo mês.

**Critérios de aceite (RF08, RF09 — RN05):**
- AC01: É possível cadastrar uma assinatura informando: nome do serviço, membro, categoria, valor mensal, dia da cobrança, forma de pagamento.
- AC02: Assinatura ativa é automaticamente incluída no total de despesas de todos os meses.
- AC03: É possível pausar ou cancelar uma assinatura; registros históricos são mantidos (RN05).
- AC04: Assinatura cancelada não aparece nas projeções futuras.

---

### EPIC: Contas Fixas

---

#### US-11 — Cadastrar conta fixa mensal 🟡

**Como** usuário, **quero** cadastrar contas fixas como condomínio e financiamento, **para que** elas sejam consideradas automaticamente nos cálculos mensais.

**Critérios de aceite (RF10, RF11 — RN05):**
- AC01: É possível cadastrar uma conta fixa informando: descrição, membro, categoria, valor, dia do vencimento, forma de pagamento.
- AC02: Conta fixa ativa é incluída automaticamente no total de despesas de todos os meses.
- AC03: É possível alterar o valor da conta fixa (para reajustes); a alteração vale a partir do mês seguinte ao da edição.
- AC04: Conta inativada é removida das projeções futuras, mas mantida no histórico.

---

### EPIC: Despesas Gerais

---

#### US-12 — Lançar despesa geral (à vista, débito, PIX) 🟡

**Como** usuário, **quero** registrar despesas pagas em dinheiro, débito ou PIX, **para que** elas apareçam no saldo do mês correto.

**Critérios de aceite (RF12):**
- AC01: É possível lançar uma despesa geral informando: data, membro, descrição, categoria, valor, forma de pagamento (dinheiro, débito, PIX).
- AC02: A despesa aparece no total de despesas do mês correspondente à data informada.
- AC03: Campo de observações é opcional.

---

#### US-13 — Filtrar e visualizar histórico de despesas gerais 🟡

**Como** usuário, **quero** filtrar o histórico de despesas gerais, **para que** eu encontre rapidamente um lançamento específico ou analise gastos de um período.

**Critérios de aceite (RF13):**
- AC01: É possível filtrar por: período (mês/ano), categoria, forma de pagamento, membro.
- AC02: Os filtros podem ser combinados.
- AC03: O total das despesas filtradas é exibido.
- AC04: É possível exportar o resultado em CSV.

---

### EPIC: Rendas (avançado)

---

#### US-14 — Registrar renda variável do mês 🟡

**Como** usuário, **quero** lançar a renda variável (freelance, comissão) de cada mês, **para que** ela seja considerada no saldo do mês em que foi recebida.

**Critérios de aceite (RF15 — RN07):**
- AC01: É possível registrar uma renda variável informando: descrição, membro, mês/ano de referência, valor recebido, data de recebimento.
- AC02: Cada lançamento é um registro individual por mês.
- AC03: Renda variável entra no saldo do mês referenciado.
- AC04: Renda variável **não** é projetada automaticamente para meses futuros (RN07).

---

#### US-15 — Registrar renda extra / pontual 🟡

**Como** usuário, **quero** registrar recebimentos pontuais como bônus ou venda de item, **para que** eles sejam contabilizados no saldo do mês correto.

**Critérios de aceite (RF16 — RN07):**
- AC01: É possível registrar uma renda extra informando: descrição, membro, data de recebimento, valor.
- AC02: Renda extra entra no saldo do mês da data de recebimento.
- AC03: Renda extra **não** é projetada para meses futuros (RN07).

---

#### US-16 — Registrar rendimento de investimento 🟡

**Como** usuário, **quero** registrar dividendos e rendimentos de CDB/FIIs, **para que** eu tenha visibilidade do total recebido e possa decidir quanto injetar no orçamento familiar.

**Critérios de aceite (RF17, RF17a — RN09):**
- AC01: É possível registrar um rendimento informando: descrição, membro, data, valor total recebido.
- AC02: O campo `valor_distribuido` é opcional e tem padrão R$ 0,00.
- AC03: `valor_distribuido` não pode ser maior que o `valor` do rendimento.
- AC04: O dashboard exibe em seção separada: total recebido, total distribuído, total reinvestido (RN09).
- AC05: Somente o `valor_distribuido` entra no cálculo do saldo do mês (RN06, RN09).
- AC06: É possível editar o `valor_distribuido` de qualquer rendimento do mês corrente.

---

### EPIC: Dashboard (completo)

---

#### US-17 — Visualizar resumo completo do mês 🟡

**Como** usuário, **quero** ver o dashboard completo com todas as fontes de renda e despesa, **para que** eu tenha uma visão clara do saldo real da família.

**Critérios de aceite (RF19 — RN06, RN08):**
- AC01: Dashboard exibe todas as seções do mês selecionado:
  - **Rendas Operacionais:** fixas + variáveis lançadas + extras
  - **Rendimentos de Investimentos:** total recebido / total distribuído / total reinvestido
  - **Despesas:** gerais + assinaturas + contas fixas + fatura dos cartões
  - **SALDO = (Rendas Operacionais + Total Distribuído) − Total Despesas**
- AC02: Saldo positivo em verde; negativo em vermelho.
- AC03: Visão consolidada familiar e por membro (RN08).
- AC04: Seletor de mês/ano para navegar entre períodos.

---

#### US-18 — Visualizar despesas por categoria 🟡

**Como** usuário, **quero** ver um gráfico de despesas por categoria do mês, **para que** eu identifique onde estou gastando mais.

**Critérios de aceite (RF20 — RN04):**
- AC01: Exibe valor e percentual de cada categoria sobre o total de despesas.
- AC02: Considera todas as fontes: despesas gerais, parcelas de cartão, assinaturas e contas fixas (RN04).
- AC03: É possível filtrar por membro ou exibir consolidado familiar.
- AC04: Gráfico de pizza ou barras horizontais.

---

## FASE 3 — Analytics

### EPIC: Relatórios e Gráficos

---

#### US-19 — Visualizar evolução mensal (gráfico de linha) 🟢

**Como** usuário, **quero** ver um gráfico com a evolução de rendas, despesas e saldo mês a mês, **para que** eu identifique tendências financeiras da família.

**Critérios de aceite (RF21):**
- AC01: Gráfico de linha exibindo os últimos 12 meses (ou período configurável) com três séries: Rendas, Despesas e Saldo.
- AC02: Ao passar o cursor sobre um ponto, exibe os valores detalhados do mês.
- AC03: Comparativo do mês atual vs. mês anterior é exibido em destaque.

---

#### US-20 — Visualizar projeção dos próximos 3 meses 🟢

**Como** usuário, **quero** ver uma projeção de rendas e despesas dos próximos 3 meses, **para que** eu possa planejar gastos futuros com antecedência.

**Critérios de aceite (RF22 — RN07, RN09):**
- AC01: A projeção exibe para cada um dos 3 meses seguintes:
  - Rendas fixas ativas (única renda projetada automaticamente — RN07)
  - Parcelas de cartão a vencer no mês
  - Total de assinaturas ativas
  - Total de contas fixas ativas
  - **Total estimado de despesas**
  - **Saldo estimado** (rendas fixas − despesas estimadas)
- AC02: Rendas variáveis, extras e rendimentos de investimentos **não** entram nas projeções (RN07, RN09).
- AC03: O saldo estimado é exibido com indicação visual verde/vermelho.

---

#### US-21 — Visualizar histórico e resumo de rendas 🟢

**Como** usuário, **quero** visualizar e filtrar todas as rendas por tipo, membro e período, **para que** eu tenha controle do histórico de recebimentos da família.

**Critérios de aceite (RF18):**
- AC01: Listagem de rendas filtráveis por tipo (fixa, variável, extra, investimento), membro e mês/ano.
- AC02: Total consolidado por membro e da família é exibido.
- AC03: Filtros podem ser combinados.

---

## MELHORIAS — Revisões pós-MVP

---

#### US-22 — Definir período de vigência da renda fixa 🟡

**Como** usuário, **quero** informar uma data de início (obrigatória) e uma data de fim (opcional) ao cadastrar uma renda fixa, **para que** o sistema inclua essa renda somente nos meses em que ela realmente vigora, sem precisar inativá-la manualmente.

**Critérios de aceite (RF14 — RN10):**
- AC01: O formulário de criação e edição de renda fixa exibe o campo "Data de início" (obrigatório) e "Data de fim" (opcional).
- AC02: O sistema usa apenas o mês e o ano de `data_inicio` e `data_fim` para determinar a vigência; o dia exato é ignorado no cálculo (RN10).
- AC03: Uma renda fixa ativa é incluída no cálculo e nas projeções de um mês M somente se M >= mês/ano de `data_inicio` E (data_fim é nula OU M <= mês/ano de `data_fim`).
- AC04: Quando `data_fim` é nula, a renda é tratada como sem fim previsto — comportamento idêntico ao anterior à melhoria.
- AC05: O dashboard e as projeções respeitam automaticamente o período de vigência, sem nenhuma ação manual do usuário.
- AC06: Rendas com `data_fim` no passado continuam visíveis no histórico, mas não entram em cálculos de meses futuros.
- AC07: A migration que adiciona `data_inicio` define como valor padrão `created_at` dos registros existentes (retrocompatibilidade); `data_fim` é nullable sem padrão.

---

## Tarefas Técnicas (sem valor direto para o usuário, mas necessárias)

---

#### TT-01 — Setup do projeto (scaffolding) 🔴

- Criar estrutura de repositório (monorepo ou repos separados: `frontend/` e `backend/`)
- Configurar frontend: React 18 + TypeScript + Vite + Tailwind CSS + shadcn/ui
- Configurar backend: Go + Gin, estrutura `cmd/`, `internal/domain/`, `internal/service/`, `internal/handler/`, `internal/repository/`
- Configurar linters: ESLint + Prettier (JS/TS), golangci-lint (Go)
- Configurar Makefile com targets de build, test e run

---

#### TT-02 — Banco de dados e migrações 🔴

- Configurar PostgreSQL via Docker Compose para desenvolvimento local
- Criar migrações iniciais para todas as entidades do modelo de dados (spec.md §7):
  `membros`, `categorias`, `cartoes_credito`, `despesas_cartao`, `assinaturas`, `contas_fixas`, `despesas_gerais`, `rendas_fixas`, `rendas_variaveis`, `rendas_extras`, `rendimentos_investimento`
- Todos os registros devem incluir: `id` (UUID), `created_at`, `updated_at`, `deleted_at` (soft delete para sync)
- Configurar ferramenta de migração: `golang-migrate`

---

#### TT-03 — Autenticação JWT 🔴

- Implementar endpoint `POST /api/v1/auth/login`
- Implementar endpoint `POST /api/v1/auth/refresh`
- Middleware de autenticação para rotas protegidas
- Senhas armazenadas com bcrypt (custo ≥ 12)
- Seed de 2 usuários iniciais via configuração (RNF04)

---

#### TT-04 — Infraestrutura offline (IndexedDB + sync) 🟡

- Configurar Dexie.js para persistência local no frontend
- Implementar fila de sincronização: operações locais com status `pendente` → envio ao backend quando online
- Política de conflito: última escrita vence por `updated_at`
- Soft delete: campo `deleted_at` sincronizado entre cliente e servidor (RNF03, RNF05)

---

#### TT-05 — PWA (Service Worker + manifest) 🟡

- Configurar `vite-plugin-pwa` com Service Worker para cache de assets
- Criar `manifest.json` (ícones, nome, display standalone)
- Garantir funcionamento offline completo para lançamento e consulta (RNF02, RNF05)

---

#### TT-06 — CI/CD com GitHub Actions 🟢

- Pipeline de build e testes para frontend (Vitest) e backend (Go test)
- Deploy automático: frontend → Vercel/Netlify; backend → Railway
- Lint obrigatório no CI

---

## Resumo por Fase

| Fase | Stories | Tarefas Técnicas |
|------|---------|-----------------|
| 🔴 MVP | US-01 a US-07 | TT-01, TT-02, TT-03 |
| 🟡 Core | US-08 a US-18 | TT-04, TT-05 |
| 🟢 Analytics | US-19 a US-21 | TT-06 |
| 🟡 Melhoria | US-22 | — |

---

**Total: 22 histórias de usuário + 6 tarefas técnicas**
