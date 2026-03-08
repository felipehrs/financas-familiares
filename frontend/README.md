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
| Gerenciador de pacotes | npm |
| Offline / IndexedDB | Dexie.js |
| PWA | vite-plugin-pwa (Workbox) |

---

## Estrutura

```
src/
├── api/           # Clientes HTTP para cada recurso da API (chamadas diretas ao backend)
├── offline/       # Wrappers offline: network-first com fallback Dexie + sync queue
├── lib/           # db.ts (Dexie schema) e syncQueue.ts (fila de sincronização)
├── components/
│   ├── ui/        # Componentes shadcn/ui (Button, Card, Input, Label, Form)
│   ├── GraficoCategoriaDespesas.tsx
│   ├── ProtectedRoute.tsx
│   └── SyncQueueInitializer.tsx
├── hooks/         # Hooks customizados (useAuth, useTheme, useSyncQueue)
├── pages/         # Páginas da aplicação
├── store/         # Estado global (AuthContext)
├── test/          # Setup do Vitest (fake-indexeddb/auto)
└── types/         # Tipos TypeScript compartilhados
```

---

## Pré-requisitos

- Node.js 20+

---

## Instalação

```bash
npm install
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
npm run dev

# Build de produção
npm run build

# Preview do build (PWA ativo)
npm run preview

# Testes (watch mode)
npm test

# Testes (CI — sem watch)
npm run test:run

# Cobertura de testes
npm run test:coverage

# Lint
npm run lint
```

---

## Testes

O projeto adota **TDD** — testes são escritos antes da implementação.

```bash
npm run test:run
```

| Camada | Ferramenta | Cobertura mínima |
|---|---|---|
| Funções puras / hooks | Vitest | 90% |
| Componentes React | Vitest + RTL | 60% |
| Offline / sync queue | Vitest + fake-indexeddb | 90% |

---

## Offline

As páginas importam de `@/offline/*` em vez de `@/api/*`. Cada wrapper:

- **Leituras:** tenta a rede; em caso de falha, serve do cache Dexie (IndexedDB)
- **Escritas online:** chama a API e atualiza o cache local
- **Escritas offline:** salva localmente (UUID gerado no cliente) e enfileira na `sync_queue`
- **Reconexão:** o evento `window.online` dispara `processarFila`, que envia as operações pendentes ao backend em ordem cronológica (política last-write-wins)

---

## PWA

O app é instalável como PWA. O Service Worker (Workbox, modo `generateSW`) faz precache de todos os assets estáticos. Chamadas à API (`/api/*`) não são interceptadas pelo SW — tratadas pelo Dexie.

Para verificar o PWA localmente:

```bash
npm run build && npm run preview
```

Abra `http://localhost:4173` no Chrome e use DevTools > Application.

---

## Autenticação

- Access token armazenado **em memória** (mais seguro contra XSS)
- Refresh token no `localStorage`
- Renovação automática do token ao expirar ou receber 401
