# Sprints — Finanças Familiares

**Referência:** stories.md | spec.md | tech-spec.md
**Atualizado em:** 11/03/2026 (Sprint 9 TT-07 em andamento — progresso detalhado abaixo)

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
| ✅ | US-15 | Registrar renda extra / pontual |
| ✅ | US-16 | Registrar rendimento de investimento com valor_distribuido (RN09) |

**Critério de conclusão:** Todos os tipos de lançamento financeiro funcionando. RN09 coberto por testes.

---

## Sprint 6 — Dashboard Completo e Offline

**Objetivo:** Completar o dashboard com todas as seções e implementar a infraestrutura offline (PWA + IndexedDB).

**Itens:**

| Status | ID | Descrição |
|--------|-----|-----------|
| ✅ | US-17 | Dashboard completo: rendas operacionais + rendimentos de investimentos + despesas + saldo (RN06, RN08) |
| ✅ | US-18 | Despesas por categoria: gráfico de pizza/barras com percentuais (RN04) |
| ✅ | TT-04 | Infraestrutura offline: Dexie.js + fila de sincronização + política last-write-wins |
| ✅ | TT-05 | PWA: Service Worker + manifest.json + funcionamento offline completo |

**Critério de conclusão:** Dashboard completo conforme wireframe da spec.md §8. App funciona offline e sincroniza ao reconectar.

---

## Sprint 7 — Analytics e Projeções

**Objetivo:** Implementar os gráficos de evolução mensal e as projeções dos próximos 3 meses.

**Itens:**

| Status | ID | Descrição |
|--------|-----|-----------|
| ✅ | US-19 | Gráfico de evolução mensal (Recharts): rendas, despesas e saldo dos últimos 12 meses |
| ✅ | US-20 | Projeção dos próximos 3 meses: rendas fixas vs. despesas recorrentes + parcelas (RN07) |
| ✅ | US-21 | Histórico e resumo de rendas filtráveis por tipo, membro e período |

**Critério de conclusão:** Projeção correta conforme RN07 e RN09, coberta por testes. Gráficos renderizando com dados reais.

---

## Sprint 8 — Qualidade e Deploy

**Objetivo:** CI/CD configurado, cobertura de testes validada e primeiro deploy em produção.

**Itens:**

| Status | ID | Descrição |
|--------|-----|-----------|
| ✅ | TT-06 | CI/CD com GitHub Actions: build, lint e testes obrigatórios; deploy automático |
| ✅ | — | Validar cobertura mínima por camada (tech-spec.md §11) — backend: handler+service testados; frontend: 176 testes em 22 suítes |
| ✅ | — | Configs de deploy criadas: `railway.toml` (backend) + `vercel.json` (frontend) + `backend/Dockerfile` multi-stage |
| ✅ | — | Deploy efetivo em produção + smoke tests (requer contas Railway e Vercel) |

**Critério de conclusão:** Pipeline verde, app acessível em produção, smoke tests passando.

---

## Sprint 9 — Segurança e Isolamento de Dados

**Objetivo:** Corrigir problemas críticos de segurança identificados na avaliação pós-Sprint 8: isolamento de dados por usuário (atualmente todos os usuários veem os dados uns dos outros), restrição de CORS e rate limiting no login, e índices de performance no banco.

**Itens:**

| Seq | Status | ID | Descrição |
|-----|--------|-----|-----------|
| 1 | 🔄 | TT-07 | Infraestrutura Multi-Tenant: tabelas `familias`, colunas `familia_id`, alteração JWT e Middleware |
| 2 | 🔲 | TT-10 | Refatoração para Isolamento: filtrar 15+ repositórios e handlers por `familia_id` |
| 3 | 🔲 | TT-08 | Segurança: restringir CORS a origens específicas; rate limiting em `POST /auth/login` |
| 4 | 🔲 | TT-09 | Índices de performance: `deleted_at` em todas as tabelas, `familia_id` em todas as tabelas, índice composto em `despesas_cartao` |

**Critério de conclusão:** TT-07 validado com teste de isolamento (usuário A não vê dados do B). CORS restrito a domínios específicos e rate limiting ativo no login com resposta 429. Índices criados e verificados em todas as tabelas afetadas.

### Progresso TT-07 (atualizado em 11/03/2026)

> **Contexto:** O TT-07 abrange toda a refatoração backend para isolamento por família. As subtarefas abaixo cobrem o trabalho já feito (marcado ✅) e o que ainda falta (🔲).

| Subtarefa | Status | Detalhe |
|-----------|--------|---------|
| Migrations 000017 up/down | ✅ | Criadas: tabelas `familias`, `familia_usuarios`; colunas `familia_id` em todas as tabelas de dados |
| `domain/familia.go` | ✅ | Criado |
| `repository/familia_repository.go` | ✅ | Criado (`BuscarFamiliaPorUsuario`) |
| `middleware/auth.go` | ✅ | `ValidateAccessToken` retorna 3 valores; seta `"familiaID"` no contexto |
| `handler/helpers.go` | ✅ | `getFamiliaID(c)` criado |
| `service/auth_service.go` | ✅ | Interface `FamiliaRepository` adicionada; `familiaID` no JWT |
| 11 domain services (interfaces + implementações) | ✅ | `familiaID string` como 1º param em todas as interfaces e chamadas repo |
| `service/rendimento_investimento_service_test.go` | ✅ | Mocks e chamadas atualizados com `familiaID` |
| `service/dashboard_service.go` | ✅ | Interfaces e métodos públicos atualizados com `familiaID` |
| `service/dashboard_service_test.go` | ✅ | Mocks e chamadas atualizados com `familiaID` |
| `service/renda_historico_service.go` | 🔲 | Interfaces e `BuscarHistorico` precisam de `familiaID` |
| `service/renda_historico_service_test.go` | 🔲 | Mocks e chamadas precisam de `familiaID` |
| 13 repositories (queries SQL) | 🔲 | Cada método precisa de `familiaID` param + `AND familia_id = $N` nas queries |
| 13 handlers (interfaces locais + chamadas) | 🔲 | Interfaces locais e cada chamada de service precisam de `familiaID` |
| 13 handler tests (mocks + setupRouter) | 🔲 | Mocks e `setupRouter` precisam injetar `familiaID` no contexto |
| `repository/seed.go` | 🔲 | `SeedCategorias(familiaID)` + lógica de família no `SeedUsuarios` |
| `cmd/server/main.go` | 🔲 | Criar `familiaRepo`, passar para `authService`; remover chamada separada de `SeedCategorias` |
| `go build ./...` limpo | 🔲 | Zero erros de compilação |
| `go test ./...` verde | 🔲 | Zero falhas de teste |

**Pré-requisito:** Sprints 1–8 concluídas.

---

## Sprint 10 — Melhorias de Usabilidade

**Objetivo:** Implementar melhorias de UX e de modelagem decididas após a conclusão das sprints 1–8: projeção com pior cenário histórico, contas fixas com valor variável e reajustes, suporte a moeda estrangeira em assinaturas, vinculação de despesa geral a cartão de crédito, e tela centralizada de despesas com formulário dinâmico por tipo.

**Itens:**

| Seq | Status | ID | Descrição |
|-----|--------|-----|-----------|
| 1 | 🔲 | US-23 | Assinaturas em moeda estrangeira: campo `moeda`, busca de cotação via API de câmbio na data de fechamento, cotação manual pelo usuário (RF08a, RN11) |
| 2 | 🔲 | US-24 | Despesa geral com cartão de crédito: campo `cartao_id` quando forma de pagamento = cartão; entrada na fatura via RN01 (RF12 revisado, RN12) |
| 3 | 🔲 | US-26 | Contas fixas com valor variável: campo `tipo_valor` (fixo/variável); lançamento mensal do valor real; estimativa automática = máximo dos últimos 3 meses (RF10a, RN14) |
| 4 | 🔲 | US-27 | Reajuste de contas fixas: registro de novo valor com data de vigência; histórico de reajustes; valor correto aplicado automaticamente por mês (RF10b, RN15) |
| 5 | 🔲 | US-28 | Projeção com pior cenário: contas fixas e despesas gerais usam o maior valor dos últimos 3 meses como estimativa (RF22 revisado, RN07 revisado) — depende de US-26/US-27 |
| 6 | 🔲 | US-25 | Tela centralizada de despesas: lista unificada com badges por tipo abaixo do dashboard; botão "+ Nova Despesa" com seletor de tipo e formulário dinâmico (RF12a, RN13) |

**Critério de conclusão:** Projeção conservadora usando pior cenário histórico (US-28). Estimativa de conta variável usando max dos últimos 3 meses; valor real substitui estimativa quando lançado. Reajuste aplicado corretamente na linha do tempo. Cotação de moeda estrangeira automática com override manual. Lista centralizada com badges e formulário dinâmico. Testes unitários ≥ 90% cobrindo RN07 revisado, RN11–RN15.

**Pré-requisito:** Sprint 9 concluída.

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
| 9 | Segurança e isolamento | TT-07, TT-10, TT-08, TT-09 |
| 10 | Melhorias de usabilidade | US-28, US-23, US-24, US-25, US-26, US-27 |
