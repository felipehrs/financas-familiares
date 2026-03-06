# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Status

This project is in the **specification phase**. The only file currently present is `spec.md`, which contains the full functional specification. No implementation exists yet. All decisions about tech stack, folder structure, and tooling are open.

## Domain Overview

A family finance management system (in Portuguese Brazilian). Core domain concepts:

- **Membros** (family members) — configurable, not hardcoded to a couple
- **Categorias** — user-defined expense/income categories
- **CartaoCredito** — credit cards with closing day (`dia_fechamento`) and due day (`dia_vencimento`)
- **DespesaCartao** — credit card expenses with installments; `fatura` and `valor_parcela` are computed fields
- **Assinatura** — recurring subscriptions (active/cancelled)
- **ContaFixa** — fixed monthly bills (active/inactive)
- **DespesaGeral** — one-off general expenses (cash, debit, PIX, etc.)
- **RendaFixa** — fixed monthly income (salary); auto-included in projections
- **RendaVariavel** — variable recurring income (freelance, commissions); entered monthly
- **RendaExtra** — one-off income (bonuses, sales)
- **RendimentoInvestimento** — investment returns with `valor_distribuido` (portion injected into the budget; default 0)

## Critical Business Rules

**Credit card invoice assignment (RN01):**
- Purchase date ≤ `dia_fechamento` → invoice belongs to current month
- Purchase date > `dia_fechamento` → invoice belongs to next month

**Installment distribution (RN03):** Installment 1 goes to the calculated invoice month; each subsequent installment goes to the following month.

**Monthly balance formula (RN06):**
```
SALDO = (RendaFixa + RendaVariavel + RendaExtra + valor_distribuido_investimentos) − (DespesasGerais + Assinaturas + ContasFixas + FaturaCartoes)
```
Investment returns that are **not distributed** (`valor_distribuido = 0`) are informational only and do not affect the balance.

**Future projections (RN07, RN09):** Only `RendaFixa` (active) is projected forward. Variable income, extras, and investment returns are excluded from projections because they are unpredictable.

## Specification Reference

Full requirements are in `spec.md`:
- RF01–RF18: functional requirements (members, categories, credit cards, subscriptions, fixed bills, general expenses, income types)
- RF19–RF22: dashboard (monthly summary, category breakdown, monthly trend chart, 3-month projection)
- RN01–RN09: business rules
- Section 7: conceptual data model with all entity fields
- Section 8: dashboard wireframe showing exact layout and calculations
- Section 9: suggested tech stack (React/Vue frontend, Node/Python backend, PostgreSQL/SQLite)
- Section 10: phased roadmap (MVP → Core → Analytics → Extras)

## Tech Stack Decisions (to be made)

When implementation begins, the recommended starting point per the spec is:
- **Frontend:** React.js or Vue.js (PWA, mobile-first, 320px–1920px)
- **Storage:** SQLite (offline-first) or IndexedDB, with optional cloud sync
- **Backend (optional):** Node.js/Express or Python/FastAPI + PostgreSQL

The MVP (Fase 1) should implement: members, categories, credit cards, basic card expenses (no installments), fixed income, and a basic dashboard.
