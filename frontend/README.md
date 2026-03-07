# Frontend — Finanças Familiares

> Interface web PWA do sistema de gestão financeira familiar.

![React](https://img.shields.io/badge/React-19-61DAFB?style=flat-square&logo=react&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-5.9-3178C6?style=flat-square&logo=typescript&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-7-646CFF?style=flat-square&logo=vite&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-4-06B6D4?style=flat-square&logo=tailwindcss&logoColor=white)
![Vitest](https://img.shields.io/badge/Vitest-4-6E9F18?style=flat-square&logo=vitest&logoColor=white)

---

## Stack

| Propósito | Biblioteca |
|---|---|
| Framework | React 19 + TypeScript |
| Bundler | Vite 7 |
| Estilo | Tailwind CSS 4 + shadcn/ui |
| Roteamento | React Router v7 |
| Formulários | React Hook Form + Zod |
| Testes | Vitest + React Testing Library |
| Gerenciador de pacotes | pnpm |

---

## Estrutura

```
src/
├── api/           # Clientes HTTP para cada recurso da API
├── components/
│   ├── ui/        # Componentes shadcn/ui (Button, Card, Input, Label, Form)
│   └── ProtectedRoute.tsx
├── hooks/         # Hooks customizados (useAuth)
├── pages/         # Páginas da aplicação
├── store/         # Estado global (AuthContext)
├── test/          # Setup do Vitest
└── types/         # Tipos TypeScript compartilhados
```

---

## Pré-requisitos

- Node.js 20+
- pnpm 9+

---

## Instalação

```bash
pnpm install
```

---

## Variáveis de ambiente

Copie `.env.example` para `.env.local` e ajuste:

```bash
cp .env.example .env.local
```

| Variável | Descrição | Padrão |
|---|---|---|
| `VITE_API_URL` | URL base da API backend | `http://localhost:8080` |

---

## Comandos

```bash
# Desenvolvimento (hot reload)
pnpm dev

# Build de produção
pnpm build

# Preview do build
pnpm preview

# Testes (watch mode)
pnpm test

# Testes (CI — sem watch)
pnpm test:run

# Cobertura de testes
pnpm test:coverage

# Lint
pnpm lint

# Formatação
pnpm format
```

---

## Testes

O projeto adota **TDD** — testes são escritos antes da implementação.

```bash
pnpm test:run
```

| Camada | Ferramenta | Cobertura mínima |
|---|---|---|
| Funções puras / hooks | Vitest | 90% |
| Componentes React | Vitest + RTL | 60% |

---

## Autenticação

- Access token armazenado **em memória** (mais seguro contra XSS)
- Refresh token no `localStorage`
- Renovação automática do token ao expirar ou receber 401
