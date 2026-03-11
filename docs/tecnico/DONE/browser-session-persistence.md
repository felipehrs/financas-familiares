# Browser: Sessão Perdida ao Fechar e Reabrir Aba

> Criado em: 10/03/2026

## Problema

Ao fechar uma aba do navegador e reabri-la após 15+ minutos, o usuário precisa fazer login novamente — mesmo com o `refresh_token` ainda válido (7 dias).

---

## Causa Raiz

**Confusão semântica entre `TOKEN_EXPIRY_KEY` (access token, 15 min) e a validade real do refresh token (7 dias).**

Em `frontend/src/store/authStore.tsx`, a função `restoreSession()` faz a seguinte verificação:

```typescript
// BUGADO — authStore.tsx, linhas 73-90
if (storedRefreshToken && storedExpiry > Date.now()) {
  // tenta refresh
} else if (storedRefreshToken && storedExpiry <= Date.now()) {
  logout() // ← BUG: interpreta "access token expirado" como "refresh token expirado"
}
```

`TOKEN_EXPIRY_KEY` armazena a expiração do **access token** (15 minutos). Quando o usuário fecha a aba e reabre depois de 15 minutos, o access token já expirou, mas o código interpreta isso como se o **refresh token** também tivesse expirado — e chama `logout()`, apagando o refresh token do localStorage.

Na próxima reabertura não há mais token, obrigando novo login.

---

## Diferença em relação ao bug do PWA (já corrigido)

O fix anterior (`91203db5`) corrigiu:
- Race condition: `isRestoringSession` aguarda restore
- `logout()` em qualquer erro de rede → agora só em 401

**Não corrigiu** a condição de expiração porque:
- **PWA reabre rapidamente** (minutos) → `TOKEN_EXPIRY_KEY` ainda está no futuro na maioria dos casos
- **Navegador desktop reabre horas depois** → `TOKEN_EXPIRY_KEY` já expirou → linha 87 dispara logout

---

## Timeline do Bug

```
T=0:    Login
        access_token válido por 15min
        TOKEN_EXPIRY_KEY = Date.now() + 900_000  ← guarda expiry do access token
        refresh_token válido por 7 dias

T=20:   Usuário FECHA a aba
        React state é destruído (accessToken = null)
        localStorage PERSISTE com os dois tokens

T=40:   Usuário REABRE a aba
        restoreSession() executa:
          storedExpiry = TOKEN_EXPIRY_KEY = (T=0 + 900s) → já no passado!
          Condição: storedRefreshToken && storedExpiry <= Date.now()
          → TRUE → logout() chamado!
          → localStorage.removeItem('refresh_token') ← apaga token válido!

T=41:   Próxima reabertura: refresh_token = null → login obrigatório
        Refresh token ainda válido no servidor por 6+ dias, mas foi destruído localmente
```

---

## Solução

Remover a verificação de `TOKEN_EXPIRY_KEY` em `restoreSession()`. O backend é a fonte da verdade sobre a validade do refresh token.

**Lógica corrigida:**
```typescript
// Se há refresh_token, tenta usá-lo — o backend diz se é válido (200) ou não (401)
if (storedRefreshToken) {
  try {
    const response = await apiRefreshToken(storedRefreshToken)
    // sucesso → restaura sessão
  } catch (err) {
    if (isUnauthorizedError(err)) {
      logout() // 401 = token revogado ou expirado no servidor
    }
    // erro de rede → mantém tokens, isAuthenticated fica false
  }
}
```

Vantagens:
- Simples: remove lógica confusa de expiração
- Correto: o servidor mantém `expires_at` no banco; um token de 7 dias é rejeitado com 401
- Resiliente: erros de rede não destroem tokens

---

## Plano de Implementação

### Arquivo 1: `frontend/src/store/authStore.tsx`

**Reescrever `restoreSession()`** para não verificar `TOKEN_EXPIRY_KEY`:

```typescript
async function restoreSession() {
  const storedRefreshToken = localStorage.getItem(REFRESH_TOKEN_KEY)

  if (storedRefreshToken) {
    try {
      const response = await apiRefreshToken(storedRefreshToken)
      const expiry = Date.now() + response.expires_in * 1000
      setAccessToken(response.access_token)
      setTokenExpiry(expiry)
      localStorage.setItem(TOKEN_EXPIRY_KEY, String(expiry))
    } catch (err) {
      if (isUnauthorizedError(err)) {
        logout() // token revogado/expirado no servidor
      }
      // erro de rede: mantém tokens para retry posterior
    }
  }

  setIsRestoringSession(false)
}
```

> `TOKEN_EXPIRY_KEY` permanece útil para `refreshIfNeeded()` (renovação periódica enquanto o app está aberto — detectar quando o access token está prestes a expirar). Apenas deixa de ser usado em `restoreSession()`.

---

### Arquivo 2: `frontend/src/store/authStore.test.tsx`

**Adicionar 2 novos testes** cobrindo o caso edge:

**Teste 1** — access token expirado, refresh token válido (o bug):
```typescript
it('deve restaurar sessão quando access token expirou mas refresh token ainda é válido', async () => {
  vi.spyOn(authApi, 'refreshToken').mockResolvedValue(mockRefreshResponse)

  localStorage.setItem(REFRESH_TOKEN_KEY, 'valid-refresh-token')
  // TOKEN_EXPIRY_KEY aponta para 20 minutos atrás (access token expirado)
  localStorage.setItem(TOKEN_EXPIRY_KEY, String(Date.now() - 20 * 60 * 1000))

  renderWithProvider()

  await waitFor(() => {
    expect(screen.getByTestId('isAuthenticated').textContent).toBe('true')
    expect(screen.getByTestId('accessToken').textContent).toBe('new-access-token')
  })

  expect(localStorage.getItem(REFRESH_TOKEN_KEY)).toBe('valid-refresh-token')
})
```

**Teste 2** — refresh token realmente expirado (servidor retorna 401):
```typescript
it('deve fazer logout quando refresh token expirou no servidor (401)', async () => {
  const expired = new Response(null, { status: 401 })
  vi.spyOn(authApi, 'refreshToken').mockRejectedValue(expired)

  localStorage.setItem(REFRESH_TOKEN_KEY, 'expired-refresh-token')
  // TOKEN_EXPIRY_KEY no passado também (access token expirado)
  localStorage.setItem(TOKEN_EXPIRY_KEY, String(Date.now() - 20 * 60 * 1000))

  renderWithProvider()

  await waitFor(() => {
    expect(screen.getByTestId('isAuthenticated').textContent).toBe('false')
    expect(screen.getByTestId('isRestoringSession').textContent).toBe('false')
  })

  // Tokens apagados por 401
  expect(localStorage.getItem(REFRESH_TOKEN_KEY)).toBeNull()
})
```

---

## Critério de Aceite

- [ ] Fechar e reabrir aba após 20 minutos não exige novo login
- [ ] Fechar e reabrir aba após 6 dias não exige novo login (refresh token ainda válido)
- [ ] Fechar e reabrir aba após 8 dias exige login (refresh token realmente expirou → 401 do servidor)
- [ ] Erro de rede ao reabrir não apaga tokens
- [ ] Testes novos passando

---

## Casos Edge Cobertos pela Solução

| Cenário | Comportamento esperado | Resultado |
|---------|----------------------|-----------|
| Reabre após 20min (access expirado, refresh válido) | Restaura sessão normalmente | ✅ Correto |
| Reabre após 8 dias (refresh expirado no servidor) | Backend retorna 401 → logout() | ✅ Correto |
| Token revogado (outra sessão usou rotation) | Backend retorna 401 → logout() | ✅ Correto |
| Erro de rede ao restaurar | Mantém tokens, isAuthenticated = false | ✅ Correto |
| Clock skew (relógio do cliente errado) | Backend decide, não o cliente | ✅ Correto |
| Reabre com aba duplicada (2 tabs) | Cada aba restaura independentemente | ✅ Correto |
