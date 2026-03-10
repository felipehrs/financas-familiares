# PWA: Persistência de Sessão no iPhone

> Criado em: 10/03/2026

## Problema

Ao salvar o app na tela inicial do iPhone e abri-lo pelo ícone (modo standalone), o usuário precisa fazer login novamente a cada abertura, mesmo com refresh token ainda válido (7 dias).

---

## Causa Raiz

Os tokens são armazenados corretamente em `localStorage` (que persiste entre aberturas no iOS standalone). O problema está em **como a sessão é restaurada**, não em onde os tokens ficam.

Há duas falhas críticas em `frontend/src/store/authStore.tsx`:

### Falha 1 — Race condition: ProtectedRoute renderiza antes da restauração terminar

O `useEffect` dispara a restauração da sessão, mas **não aguarda a Promise** antes de retornar. O `ProtectedRoute` lê `isAuthenticated` antes de `refreshIfNeeded()` terminar, vê `accessToken === null`, e redireciona para `/login`.

```typescript
// authStore.tsx — useEffect atual (DEFEITUOSO)
useEffect(() => {
  const storedRefreshToken = localStorage.getItem(REFRESH_TOKEN_KEY)
  const storedExpiry = Number(localStorage.getItem(TOKEN_EXPIRY_KEY))

  if (storedRefreshToken && storedExpiry > Date.now()) {
    refreshIfNeeded().catch(() => logout())  // ← Não aguarda; .catch chama logout!
  } else if (storedRefreshToken && storedExpiry <= Date.now()) {
    logout()
  }
}, [])
```

### Falha 2 — logout() destrói os tokens em qualquer erro de rede

Dentro de `refreshIfNeeded()`, qualquer falha (timeout, backend lento, CORS, 5xx) chama `logout()`, que apaga `refresh_token` do `localStorage`. Na próxima abertura do app não há mais token para restaurar.

```typescript
// authStore.tsx — refreshIfNeeded() atual (DEFEITUOSO)
try {
  const response = await apiRefreshToken(storedRefreshToken)
  // ...
} catch {
  logout()  // ← Apaga localStorage por erro temporário de rede!
}
```

**Consequência combinada**: Na primeira abertura o refresh falha (ex: rede mobile instável), `logout()` limpa o `localStorage`, e da segunda abertura em diante não há tokens — obrigando novo login.

---

## Fluxo Atual vs. Fluxo Esperado

### Atual (defeituoso)
```
Abre PWA
  → AuthProvider monta, dispara refresh (sem aguardar)
  → ProtectedRoute renderiza: accessToken = null → redireciona /login  ← usuário vê login
  → (em paralelo) refresh tenta chamar API
      → SE erro de rede → logout() → limpa localStorage
      → SE sucesso (tarde demais) → usuário já está em /login
```

### Esperado
```
Abre PWA
  → AuthProvider monta com isRestoringSession = true
  → ProtectedRoute: vê isRestoringSession → mostra spinner, NÃO redireciona
  → refresh aguardado: chama API, obtém novo access_token
      → SE sucesso → isRestoringSession = false, isAuthenticated = true → vai para /dashboard
      → SE erro de rede → isRestoringSession = false, mantém tokens no localStorage
          → isAuthenticated = false → redireciona /login (mas tokens preservados para próxima vez)
      → SE 401 (token revogado/expirado real) → limpa tokens → redireciona /login
```

---

## Plano de Implementação

### Arquivo 1: `frontend/src/store/authStore.tsx`

**Mudança 1a** — Adicionar estado `isRestoringSession` inicializado como `true`:

```typescript
const [isRestoringSession, setIsRestoringSession] = useState(true)
```

**Mudança 1b** — Reescrever `useEffect` para aguardar Promise e nunca chamar `logout()` em erro temporário:

```typescript
useEffect(() => {
  async function restoreSession() {
    const storedRefreshToken = localStorage.getItem(REFRESH_TOKEN_KEY)
    const storedExpiry = Number(localStorage.getItem(TOKEN_EXPIRY_KEY))

    if (storedRefreshToken && storedExpiry > Date.now()) {
      try {
        const response = await apiRefreshToken(storedRefreshToken)
        const expiry = Date.now() + response.expires_in * 1000
        setAccessToken(response.access_token)
        setTokenExpiry(expiry)
        localStorage.setItem(TOKEN_EXPIRY_KEY, String(expiry))
      } catch (err) {
        // Diferenciar erros:
        // 401 = token inválido/revogado → limpar e forçar login
        // Outros (rede, timeout, 5xx) → NÃO limpar, deixar usuário tentar manualmente
        if (err instanceof Response && err.status === 401) {
          logout()
        }
        // Caso contrário: mantém tokens, isAuthenticated ficará false (sem access_token)
      }
    } else if (storedRefreshToken && storedExpiry <= Date.now()) {
      // Refresh token expirou de verdade
      logout()
    }

    setIsRestoringSession(false)
  }

  void restoreSession()
  // eslint-disable-next-line react-hooks/exhaustive-deps
}, [])
```

**Mudança 1c** — Remover `logout` do catch de `refreshIfNeeded()` (renovação periódica):

```typescript
const refreshIfNeeded = useCallback(async () => {
  const storedExpiry = tokenExpiry ?? Number(localStorage.getItem(TOKEN_EXPIRY_KEY))
  const storedRefreshToken = localStorage.getItem(REFRESH_TOKEN_KEY)

  if (!storedRefreshToken) return

  const isExpiringSoon = storedExpiry - Date.now() < 60_000

  if (!accessToken || isExpiringSoon) {
    try {
      const response = await apiRefreshToken(storedRefreshToken)
      const expiry = Date.now() + response.expires_in * 1000
      setAccessToken(response.access_token)
      setTokenExpiry(expiry)
      localStorage.setItem(TOKEN_EXPIRY_KEY, String(expiry))
    } catch (err) {
      // Só faz logout se 401 (token revogado). Erros de rede: ignora silenciosamente.
      if (err instanceof Response && err.status === 401) {
        logout()
      }
    }
  }
}, [accessToken, tokenExpiry, logout])
```

**Mudança 1d** — Expor `isRestoringSession` no contexto:

```typescript
// No value do Provider
{
  accessToken,
  isAuthenticated: accessToken !== null,
  isRestoringSession,
  login,
  logout,
  refreshIfNeeded,
}
```

**Mudança 1e** — Adicionar `isRestoringSession` ao tipo do contexto:

```typescript
interface AuthContextValue {
  accessToken: string | null
  isAuthenticated: boolean
  isRestoringSession: boolean
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  refreshIfNeeded: () => Promise<void>
}
```

---

### Arquivo 2: `frontend/src/components/ProtectedRoute.tsx`

Aguardar restauração antes de decidir redirecionar:

```typescript
export function ProtectedRoute() {
  const { isAuthenticated, isRestoringSession } = useAuth()

  if (isRestoringSession) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="text-muted-foreground text-sm">Carregando...</div>
      </div>
    )
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  return (
    <SyncQueueInitializer>
      <Outlet />
    </SyncQueueInitializer>
  )
}
```

---

### Arquivo 3: `frontend/src/api/auth.ts` (verificar)

Confirmar que o backend recebe o `refresh_token` como `Authorization: Bearer`:

```typescript
// Verificar em backend/internal/handler/auth_handler.go
// Se o handler lê do header → está correto
// Se o handler lê do body → mudar para:
body: JSON.stringify({ refresh_token: token })
```

> Checar `backend/internal/handler/auth_handler.go` função `Refresh()` antes de alterar.

---

## Critério de Aceite

- [ ] Abrir o PWA pelo ícone da home screen do iPhone sem fazer login novamente (com refresh token válido)
- [ ] Erro de rede ao restaurar sessão **não** apaga os tokens do localStorage
- [ ] Token revogado (401) faz logout normalmente e redireciona para /login
- [ ] Sem flash de tela de login antes de redirecionar para o dashboard
- [ ] Funciona em Safari standalone (iOS) e Chrome (Android)

---

## Notas sobre o iOS Standalone

- No modo `display: standalone` (PWA na home screen), o iOS trata o app como contexto isolado
- `localStorage` **persiste** entre aberturas — correto para armazenar tokens
- `sessionStorage` seria zerado a cada abertura — **não usar para tokens**
- O manifest.json já está configurado com `"display": "standalone"` ✓
- O Service Worker (Workbox/VitePWA) não interfere com chamadas de API autenticadas (não há cache de responses de `/auth/*`) ✓
