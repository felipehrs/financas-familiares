# Sprints — Finanças Familiares

**Referência:** stories.md | spec.md | tech-spec.md
**Atualizado em:** 07/03/2026 (US-04 concluída)

---

## Convenções

- **Status:** 🔲 Pendente | 🔄 Em andamento | ✅ Concluído | ⛔ Bloqueado
- **Itens:** cada linha referencia uma US ou TT do stories.md
- Atualizar o status dos itens conforme o desenvolvimento avança
- Registrar impedimentos na seção de cada sprint quando ocorrerem

---

## Sprint 1 — Fundação do Projeto

**Objetivo:** Ter o projeto scaffoldado, banco de dados funcionando localmente e autenticação implementada com TDD. Ao final desta sprint, é possível fazer login e o backend está pronto para receber as primeiras features.

**Itens:**

| Status | ID | Descrição |
|--------|-----|-----------|
| ✅ | TT-01 | Setup do projeto: monorepo, frontend (React + Vite + TS + Tailwind + shadcn/ui), backend (Go + Gin), linters, Makefile |
| ✅ | TT-02 | Banco de dados: Docker Compose + PostgreSQL + golang-migrate + migrações de todas as entidades |
| ✅ | TT-03 | Autenticação JWT: `POST /auth/login`, `POST /auth/refresh`, middleware, bcrypt, seed de 2 usuários |
| ✅ | US-01 | Login no sistema (frontend): tela de login, integração com JWT, refresh automático de token |

**Critério de conclusão:** É possível subir o ambiente com `docker compose up`, rodar `make test` no backend (verde), rodar `pnpm test:run` no frontend (verde) e fazer login com os dois usuários seed.

---

## Sprint 2 — Cadastros Base (MVP)

**Objetivo:** Implementar os cadastros fundamentais que são pré-requisito para todos os outros módulos: membros, categorias e cartões de crédito.

**Itens:**

| Status | ID | Descrição |
|--------|-----|-----------|
| ✅ | US-02 | Cadastrar/editar/inativar membros da família (backend + frontend) |
| ✅ | US-03 | Gerenciar categorias com sugestões padrão (backend + frontend) |
| ✅ | US-04 | Cadastrar/editar/inativar cartões de crédito (backend + frontend) |

**Critério de conclusão:** CRUD completo de membros, categorias e cartões funcionando com testes passando. Categorias padrão criadas no seed.

---

## Sprint 3 — Lançamentos e Dashboard MVP

**Objetivo:** Implementar o lançamento de despesas no cartão (à vista) e a renda fixa, culminando no dashboard básico com saldo do mês. Ao final desta sprint, o MVP está funcional.

**Itens:**

| Status | ID | Descrição |
|--------|-----|-----------|
| ✅ | US-05 | Lançar despesa no cartão (à vista) com cálculo automático de fatura (RN01, RN02) |
| ✅ | US-06 | Cadastrar renda fixa (salário) com inclusão automática nas projeções |
| ✅ | US-07 | Dashboard básico: total rendas fixas, fatura dos cartões, saldo do mês, visão por membro |

**Critério de conclusão:** Dashboard exibe saldo correto calculado com RN06. Regras RN01 e RN02 cobertas por testes unitários ≥ 90%. MVP utilizável de ponta a ponta.

---

## Sprint 4 — Cartões Avançados e Recorrências

**Objetivo:** Adicionar parcelamento de compras no cartão e os módulos de assinaturas e contas fixas.

**Itens:**

| Status | ID | Descrição |
|--------|-----|-----------|
| ✅ | US-08 | Lançar compra parcelada no cartão com distribuição automática por mês (RN03) |
| ✅ | US-09 | Visualizar fatura de cartão por mês com saldo disponível e alerta de limite |
| ✅ | US-10 | Cadastrar e gerenciar assinaturas recorrentes |
| ✅ | US-11 | Cadastrar e gerenciar contas fixas mensais |

**Critério de conclusão:** Regra RN03 coberta por testes ≥ 90%. Assinaturas e contas fixas entram automaticamente no saldo do mês.

---

## Sprint 4.5 — Melhoria: Vigência de Renda Fixa

**Objetivo:** Enriquecer o cadastro de renda fixa com período de vigência (data de início obrigatória e data de fim opcional), para que o sistema inclua cada renda somente nos meses em que ela realmente vigora — sem necessidade de inativação manual.

**Itens:**

| Status | ID | Descrição |
|--------|-----|-----------|
| ✅ | US-22 | Adicionar data de início (obrigatória) e data de fim (opcional) à renda fixa; cálculos e projeções respeitam o período (RN10) |

**Critério de conclusão:** Migration aplicada com retrocompatibilidade. Dashboard e projeções consideram `data_inicio` e `data_fim` corretamente. RN10 coberta por testes unitários ≥ 90%.

---

## Sprint 5 — Despesas Gerais e Rendas Completas

**Objetivo:** Completar os módulos de despesas gerais e todos os tipos de renda (variável, extra, investimentos).

**Itens:**

| Status | ID | Descrição |
|--------|-----|-----------|
| ✅ | US-12 | Lançar despesa geral (dinheiro, débito, PIX) |
| ✅ | US-13 | Filtrar e visualizar histórico de despesas gerais com exportação CSV |
| ✅ | US-14 | Registrar renda variável do mês |
| 🔲 | US-15 | Registrar renda extra / pontual |
| 🔲 | US-16 | Registrar rendimento de investimento com valor_distribuido (RN09) |

**Critério de conclusão:** Todos os tipos de lançamento financeiro funcionando. RN09 coberto por testes.

---

## Sprint 6 — Dashboard Completo e Offline

**Objetivo:** Completar o dashboard com todas as seções e implementar a infraestrutura offline (PWA + IndexedDB).

**Itens:**

| Status | ID | Descrição |
|--------|-----|-----------|
| 🔲 | US-17 | Dashboard completo: rendas operacionais + rendimentos de investimentos + despesas + saldo (RN06, RN08) |
| 🔲 | US-18 | Despesas por categoria: gráfico de pizza/barras com percentuais (RN04) |
| 🔲 | TT-04 | Infraestrutura offline: Dexie.js + fila de sincronização + política last-write-wins |
| 🔲 | TT-05 | PWA: Service Worker + manifest.json + funcionamento offline completo |

**Critério de conclusão:** Dashboard completo conforme wireframe da spec.md §8. App funciona offline e sincroniza ao reconectar.

---

## Sprint 7 — Analytics e Projeções

**Objetivo:** Implementar os gráficos de evolução mensal e as projeções dos próximos 3 meses.

**Itens:**

| Status | ID | Descrição |
|--------|-----|-----------|
| 🔲 | US-19 | Gráfico de evolução mensal (Recharts): rendas, despesas e saldo dos últimos 12 meses |
| 🔲 | US-20 | Projeção dos próximos 3 meses: rendas fixas vs. despesas recorrentes + parcelas (RN07) |
| 🔲 | US-21 | Histórico e resumo de rendas filtráveis por tipo, membro e período |

**Critério de conclusão:** Projeção correta conforme RN07 e RN09, coberta por testes. Gráficos renderizando com dados reais.

---

## Sprint 8 — Qualidade e Deploy

**Objetivo:** CI/CD configurado, cobertura de testes validada e primeiro deploy em produção.

**Itens:**

| Status | ID | Descrição |
|--------|-----|-----------|
| 🔲 | TT-06 | CI/CD com GitHub Actions: build, lint e testes obrigatórios; deploy automático |
| 🔲 | — | Validar cobertura mínima por camada (tech-spec.md §11) |
| 🔲 | — | Deploy: frontend → Vercel/Netlify; backend → Railway; banco → Railway PostgreSQL |
| 🔲 | — | Smoke tests em produção com os dois usuários seed |

**Critério de conclusão:** Pipeline verde, app acessível em produção, smoke tests passando.

---

## Resumo Geral

| Sprint | Foco | Stories/TTs |
|--------|------|-------------|
| 1 | Fundação | TT-01, TT-02, TT-03, US-01 |
| 2 | Cadastros base | US-02, US-03, US-04 |
| 3 | MVP funcional | US-05, US-06, US-07 |
| 4 | Cartões avançados e recorrências | US-08, US-09, US-10, US-11 |
| 4.5 | Melhoria: vigência de renda fixa | US-22 |
| 5 | Despesas gerais e rendas completas | US-12 a US-16 |
| 6 | Dashboard completo e offline | US-17, US-18, TT-04, TT-05 |
| 7 | Analytics e projeções | US-19, US-20, US-21 |
| 8 | Qualidade e deploy | TT-06 + validações |
