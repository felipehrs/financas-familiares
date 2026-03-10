# 💰 Finanças Familiares

> Sistema de gestão financeira familiar com controle de cartões, parcelamentos, assinaturas, contas fixas e rendas — com dashboard consolidado e projeções futuras.

![Status](https://img.shields.io/badge/status-em%20desenvolvimento-yellow?style=flat-square)
![Last Commit](https://img.shields.io/github/last-commit/felipehrs/financas-familiares?style=flat-square)
![Repo Size](https://img.shields.io/github/repo-size/felipehrs/financas-familiares?style=flat-square)
[![Backend Coverage](https://codecov.io/gh/felipehrs/financas-familiares/graph/badge.svg?flag=backend)](https://codecov.io/gh/felipehrs/financas-familiares)
[![Frontend Coverage](https://codecov.io/gh/felipehrs/financas-familiares/graph/badge.svg?flag=frontend)](https://codecov.io/gh/felipehrs/financas-familiares)

---

## Stack

**Frontend**

![React](https://img.shields.io/badge/React-19-61DAFB?style=flat-square&logo=react&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-5.9-3178C6?style=flat-square&logo=typescript&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-7-646CFF?style=flat-square&logo=vite&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-4-06B6D4?style=flat-square&logo=tailwindcss&logoColor=white)

**Backend**

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.12-00ADD8?style=flat-square&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat-square&logo=postgresql&logoColor=white)

**Infra**

![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat-square&logo=docker&logoColor=white)
![PWA](https://img.shields.io/badge/PWA-offline--first-5A0FC8?style=flat-square&logo=pwa&logoColor=white)

---

## Sobre o projeto

O sistema resolve a necessidade de visibilidade clara e consolidada das finanças de uma família, com suporte a múltiplos membros, múltiplos cartões de crédito e diferentes tipos de renda.

**Principais funcionalidades:**

- Controle de múltiplos cartões de crédito com cálculo automático de fatura e parcelamento
- Gestão de assinaturas e contas fixas com recorrência automática
- Registro de despesas gerais (dinheiro, débito, PIX)
- Múltiplas fontes de renda: fixa (salário), variável, extra e rendimentos de investimentos
- Dashboard com saldo real do mês: `Rendas − Despesas = Saldo`
- Projeção de despesas e saldo para os próximos 3 meses
- Visão individual por membro e consolidada da família
- Modo offline-first com sincronização automática em nuvem

---

## Arquitetura

```
┌─────────────────────────────────────────────────────┐
│                    CLIENTE (PWA)                    │
│     React + TypeScript      │      IndexedDB        │
└──────────────────────┬──────────────────────────────┘
                       │ HTTPS / REST
┌──────────────────────▼──────────────────────────────┐
│                   BACKEND (API)                     │
│              Go + Gin Framework                     │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│                 BANCO DE DADOS                      │
│                  PostgreSQL                         │
└─────────────────────────────────────────────────────┘
```

---

## Desenvolvimento

### Pré-requisitos

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (ou Docker Engine no Linux)
- [VS Code](https://code.visualstudio.com/) com a extensão [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

### Subindo o ambiente com Dev Container (recomendado)

O devcontainer já inclui Go, Node.js, pnpm e `go-migrate` — não é necessário instalar nada na máquina host além do Docker.

1. Abra a pasta do projeto no VS Code
2. Aceite o prompt **"Reopen in Container"** (ou use `Ctrl+Shift+P` → `Dev Containers: Reopen in Container`)
3. Aguarde o build da imagem e a instalação das dependências do frontend (feita automaticamente no `postCreateCommand`)
4. Dentro do terminal do container:

```bash
make dev
```

Isso vai rodar as migrations e subir o backend (`http://localhost:8080`) e o frontend (`http://localhost:5173`).

### Sem Dev Container (ambiente local)

Requer Go 1.25+, Node.js, pnpm e [`golang-migrate`](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) instalados na máquina.

```bash
make dev
```

---

## Roadmap

| Sprint | Foco | Status |
|--------|------|--------|
| **Sprint 1 — Fundação** | Scaffolding do monorepo, Docker + PostgreSQL + migrations, autenticação JWT, tela de login | ✅ Concluído |
| **Sprint 2 — Cadastros Base** | Membros, categorias, cartões de crédito (CRUD completo com TDD) | ✅ Concluído |
| **Sprint 3 — MVP** | Despesas à vista no cartão (RN01/RN02), renda fixa, dashboard básico com saldo | ✅ Concluído |
| **Sprint 4 — Cartões Avançados** | Compras parceladas (RN03), fatura por mês, assinaturas, contas fixas | ✅ Concluído |
| **Sprint 4.5 — Vigência de Renda Fixa** | Data de início/fim na renda fixa; proporcionalidade por dia (RN10) | ✅ Concluído |
| **Sprint 5 — Despesas e Rendas Completas** | Despesas gerais, filtros e CSV, renda variável, renda extra, rendimentos de investimento (RN09) | ✅ Concluído |
| **Sprint 6 — Dashboard Completo e Offline** | Dashboard com todas as seções (RN06/RN08), gráfico de categorias, PWA, IndexedDB offline-first | ✅ Concluído |
| **Sprint 7 — Analytics e Projeções** | Gráfico de evolução mensal (12 meses), projeção 3 meses (RN07), histórico de rendas | ✅ Concluído |
| **Sprint 8 — Qualidade e Deploy** | CI/CD GitHub Actions, cobertura validada, deploy (Vercel + Railway) | ✅ Concluído |
| **Sprint 9 — Segurança e Isolamento de Dados** | Isolamento de dados por usuário (TT-07), CORS restrito + rate limiting (TT-08), índices de performance (TT-09) | 🔲 Futuro |
| **Sprint 10 — Melhorias de Usabilidade** | Contas fixas variáveis + reajustes, assinaturas em moeda estrangeira, lista centralizada de despesas, projeção com pior cenário | 🔲 Futuro |

---

## Documentação

| Arquivo | Descrição |
|---------|-----------|
| [`docs/produto/spec.md`](docs/produto/spec.md) | Especificação funcional completa (RF01–RF22, RN01–RN15, modelo de dados, wireframes) |
| [`docs/produto/stories.md`](docs/produto/stories.md) | Histórias de usuário com critérios de aceite por fase |
| [`docs/produto/sprints.md`](docs/produto/sprints.md) | Status detalhado de cada sprint com itens e critérios de conclusão |
| [`docs/tecnico/tech-spec.md`](docs/tecnico/tech-spec.md) | Especificação técnica (stack, arquitetura, autenticação, sync offline, infra) |
| [`docs/tecnico/avaliacao-e-multi-tenant.md`](docs/tecnico/avaliacao-e-multi-tenant.md) | Avaliação pós-deploy: problemas encontrados e roadmap para multi-tenant |

---

## Regras de negócio principais

**Cálculo de fatura do cartão (RN01)**
```
compra ≤ dia_fechamento  →  fatura do mês corrente
compra > dia_fechamento  →  fatura do mês seguinte
```

**Distribuição de parcelas (RN03)**
```
parcela 1 → fatura calculada pela data da compra
parcela 2 → mês seguinte
parcela N → mês seguinte à parcela N-1
```

**Saldo mensal (RN06)**
```
SALDO = (Rendas Fixas + Rendas Variáveis + Rendas Extras + Valor Distribuído de Investimentos)
      − (Despesas Gerais + Assinaturas + Contas Fixas + Fatura dos Cartões)
```

> Rendimentos de investimentos não distribuídos são informativos e não impactam o saldo.

**Proporcionalidade de renda fixa (RN10)**
```
Mês de início → valor × (dias restantes no mês / dias no mês)
Mês de fim    → valor × (dia de fim / dias no mês)
Mês intermediário → valor cheio
```

---

## Licença

Uso pessoal e familiar. Sem licença de código aberto definida.
