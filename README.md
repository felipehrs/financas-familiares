# 💰 Finanças Familiares

> Sistema de gestão financeira familiar com controle de cartões, parcelamentos, assinaturas, contas fixas e rendas — com dashboard consolidado e projeções futuras.

![Status](https://img.shields.io/badge/status-em%20desenvolvimento-yellow?style=flat-square)
![Last Commit](https://img.shields.io/github/last-commit/felipehrs/financas-familiares?style=flat-square)
![Repo Size](https://img.shields.io/github/repo-size/felipehrs/financas-familiares?style=flat-square)

---

## Stack

**Frontend**

![React](https://img.shields.io/badge/React-18-61DAFB?style=flat-square&logo=react&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?style=flat-square&logo=typescript&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-5-646CFF?style=flat-square&logo=vite&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-3-06B6D4?style=flat-square&logo=tailwindcss&logoColor=white)

**Backend**

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?style=flat-square&logo=go&logoColor=white)
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

## Roadmap

| Fase | Descrição | Status |
|------|-----------|--------|
| **Fase 1 — MVP** | Membros, categorias, cartões, despesas à vista, renda fixa, dashboard básico | 🔲 Planejado |
| **Fase 2 — Core** | Parcelamento, assinaturas, contas fixas, despesas gerais, todos os tipos de renda, dashboard completo | 🔲 Planejado |
| **Fase 3 — Analytics** | Gráficos de evolução, projeções futuras, relatórios | 🔲 Planejado |
| **Fase 4 — Extras** | Integração bancária, alertas, importação de extratos, app mobile | 🔲 Futuro |

---

## Documentação

| Arquivo | Descrição |
|---------|-----------|
| [`spec.md`](spec.md) | Especificação funcional completa (RF01–RF22, RN01–RN09, modelo de dados, wireframes) |
| [`tech-spec.md`](tech-spec.md) | Especificação técnica (stack, arquitetura, autenticação, sync offline, infra) |
| [`stories.md`](stories.md) | Histórias de usuário com critérios de aceite por fase |

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

---

## Licença

Uso pessoal e familiar. Sem licença de código aberto definida.
