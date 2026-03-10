import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, renderHook, screen, waitFor, act } from '@testing-library/react'
import React from 'react'
import { AuthProvider, useAuthContext } from './authStore'
import * as authApi from '@/api/auth'

// Helper para consumir o contexto em testes
function AuthConsumer({ onMount }: { onMount?: (ctx: ReturnType<typeof useAuthContext>) => void }) {
  const ctx = useAuthContext()
  React.useEffect(() => {
    onMount?.(ctx)
  })
  return (
    <div>
      <span data-testid="isAuthenticated">{String(ctx.isAuthenticated)}</span>
      <span data-testid="isRestoringSession">{String(ctx.isRestoringSession)}</span>
      <span data-testid="accessToken">{ctx.accessToken ?? 'null'}</span>
    </div>
  )
}

function renderWithProvider(ui: React.ReactNode = <AuthConsumer />) {
  return render(<AuthProvider>{ui}</AuthProvider>)
}

const REFRESH_TOKEN_KEY = 'refresh_token'
const TOKEN_EXPIRY_KEY = 'token_expiry'

const FUTURE_EXPIRY = String(Date.now() + 10 * 60 * 1000) // 10 min no futuro

const mockRefreshResponse = {
  access_token: 'new-access-token',
  token_type: 'Bearer',
  expires_in: 900,
}

describe('AuthProvider — restauração de sessão ao montar (useEffect)', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.restoreAllMocks()
  })

  afterEach(() => {
    localStorage.clear()
  })

  it('isRestoringSession deve ser true ao montar (antes da restauração terminar)', async () => {
    // Mantém a promise pendente para capturar o estado inicial
    let resolveRefresh!: (v: typeof mockRefreshResponse) => void
    vi.spyOn(authApi, 'refreshToken').mockReturnValue(
      new Promise((res) => { resolveRefresh = res }),
    )

    localStorage.setItem(REFRESH_TOKEN_KEY, 'valid-token')
    localStorage.setItem(TOKEN_EXPIRY_KEY, FUTURE_EXPIRY)

    const values: boolean[] = []
    function Spy() {
      const ctx = useAuthContext()
      values.push(ctx.isRestoringSession)
      return null
    }

    render(<AuthProvider><Spy /></AuthProvider>)

    // Primeiro render: deve estar restaurando
    expect(values[0]).toBe(true)

    // Resolve para limpar
    act(() => { resolveRefresh(mockRefreshResponse) })
  })

  it('deve chamar apiRefreshToken se houver refresh_token válido no localStorage', async () => {
    const spy = vi.spyOn(authApi, 'refreshToken').mockResolvedValue(mockRefreshResponse)

    localStorage.setItem(REFRESH_TOKEN_KEY, 'valid-token')
    localStorage.setItem(TOKEN_EXPIRY_KEY, FUTURE_EXPIRY)

    renderWithProvider()

    await waitFor(() => {
      expect(spy).toHaveBeenCalledWith('valid-token')
    })
  })

  it('deve definir accessToken após refresh bem-sucedido', async () => {
    vi.spyOn(authApi, 'refreshToken').mockResolvedValue(mockRefreshResponse)

    localStorage.setItem(REFRESH_TOKEN_KEY, 'valid-token')
    localStorage.setItem(TOKEN_EXPIRY_KEY, FUTURE_EXPIRY)

    renderWithProvider()

    await waitFor(() => {
      expect(screen.getByTestId('accessToken').textContent).toBe('new-access-token')
    })
  })

  it('deve definir isRestoringSession = false após restauração bem-sucedida', async () => {
    vi.spyOn(authApi, 'refreshToken').mockResolvedValue(mockRefreshResponse)

    localStorage.setItem(REFRESH_TOKEN_KEY, 'valid-token')
    localStorage.setItem(TOKEN_EXPIRY_KEY, FUTURE_EXPIRY)

    renderWithProvider()

    await waitFor(() => {
      expect(screen.getByTestId('isRestoringSession').textContent).toBe('false')
    })
  })

  it('deve definir isRestoringSession = false após falha na restauração', async () => {
    vi.spyOn(authApi, 'refreshToken').mockRejectedValue(new Error('Network error'))

    localStorage.setItem(REFRESH_TOKEN_KEY, 'valid-token')
    localStorage.setItem(TOKEN_EXPIRY_KEY, FUTURE_EXPIRY)

    renderWithProvider()

    await waitFor(() => {
      expect(screen.getByTestId('isRestoringSession').textContent).toBe('false')
    })
  })

  it('NÃO deve limpar localStorage em erro de rede (não-401)', async () => {
    vi.spyOn(authApi, 'refreshToken').mockRejectedValue(new Error('Network error'))

    localStorage.setItem(REFRESH_TOKEN_KEY, 'valid-token')
    localStorage.setItem(TOKEN_EXPIRY_KEY, FUTURE_EXPIRY)

    renderWithProvider()

    await waitFor(() => {
      expect(screen.getByTestId('isRestoringSession').textContent).toBe('false')
    })

    expect(localStorage.getItem(REFRESH_TOKEN_KEY)).toBe('valid-token')
    expect(localStorage.getItem(TOKEN_EXPIRY_KEY)).toBe(FUTURE_EXPIRY)
  })

  it('deve limpar localStorage e fazer logout em erro 401', async () => {
    const unauthorized = new Response(null, { status: 401 })
    vi.spyOn(authApi, 'refreshToken').mockRejectedValue(unauthorized)

    localStorage.setItem(REFRESH_TOKEN_KEY, 'valid-token')
    localStorage.setItem(TOKEN_EXPIRY_KEY, FUTURE_EXPIRY)

    renderWithProvider()

    await waitFor(() => {
      expect(screen.getByTestId('isRestoringSession').textContent).toBe('false')
    })

    expect(localStorage.getItem(REFRESH_TOKEN_KEY)).toBeNull()
    expect(localStorage.getItem(TOKEN_EXPIRY_KEY)).toBeNull()
    expect(screen.getByTestId('isAuthenticated').textContent).toBe('false')
  })

  it('deve restaurar sessão quando access token expirou mas refresh token ainda é válido', async () => {
    vi.spyOn(authApi, 'refreshToken').mockResolvedValue(mockRefreshResponse)

    localStorage.setItem(REFRESH_TOKEN_KEY, 'valid-refresh-token')
    // TOKEN_EXPIRY_KEY no passado (access token expirou há 20min)
    localStorage.setItem(TOKEN_EXPIRY_KEY, String(Date.now() - 20 * 60 * 1000))

    renderWithProvider()

    await waitFor(() => {
      expect(screen.getByTestId('isAuthenticated').textContent).toBe('true')
      expect(screen.getByTestId('accessToken').textContent).toBe('new-access-token')
    })

    expect(localStorage.getItem(REFRESH_TOKEN_KEY)).toBe('valid-refresh-token')
  })

  it('deve fazer logout quando refresh token expirou no servidor (401)', async () => {
    const expired = new Response(null, { status: 401 })
    vi.spyOn(authApi, 'refreshToken').mockRejectedValue(expired)

    localStorage.setItem(REFRESH_TOKEN_KEY, 'expired-refresh-token')
    localStorage.setItem(TOKEN_EXPIRY_KEY, String(Date.now() - 20 * 60 * 1000))

    renderWithProvider()

    await waitFor(() => {
      expect(screen.getByTestId('isAuthenticated').textContent).toBe('false')
      expect(screen.getByTestId('isRestoringSession').textContent).toBe('false')
    })

    expect(localStorage.getItem(REFRESH_TOKEN_KEY)).toBeNull()
  })

  it('deve definir isRestoringSession = false se não há refresh_token', async () => {
    // Sem nada no localStorage

    renderWithProvider()

    await waitFor(() => {
      expect(screen.getByTestId('isRestoringSession').textContent).toBe('false')
    })
  })
})

describe('AuthProvider — refreshIfNeeded (renovação periódica)', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.restoreAllMocks()
  })

  afterEach(() => {
    localStorage.clear()
  })

  it('NÃO deve fazer logout em erro de rede (não-401)', async () => {
    vi.spyOn(authApi, 'refreshToken').mockRejectedValue(new Error('Network error'))

    localStorage.setItem(REFRESH_TOKEN_KEY, 'valid-token')
    localStorage.setItem(TOKEN_EXPIRY_KEY, FUTURE_EXPIRY)

    const { result } = renderHook(() => useAuthContext(), {
      wrapper: AuthProvider,
    })

    // Espera restauração terminar (primeira chamada, que vai falhar com erro de rede)
    await waitFor(() => {
      expect(result.current.isRestoringSession).toBe(false)
    })

    // tokens devem ainda estar no localStorage (erro de rede não limpa)
    expect(localStorage.getItem(REFRESH_TOKEN_KEY)).toBe('valid-token')

    // Chama refreshIfNeeded novamente com erro de rede
    vi.spyOn(authApi, 'refreshToken').mockRejectedValue(new Error('Network error'))
    await act(async () => {
      await result.current.refreshIfNeeded()
    })

    // Tokens não devem ter sido removidos
    expect(localStorage.getItem(REFRESH_TOKEN_KEY)).toBe('valid-token')
  })

  it('deve fazer logout em erro 401 no refreshIfNeeded', async () => {
    // Arrange: restauração falha com erro de rede (accessToken fica null, mas tokens no localStorage)
    // Depois refreshIfNeeded é chamado e recebe 401
    vi.spyOn(authApi, 'refreshToken').mockRejectedValueOnce(new Error('Network error'))

    localStorage.setItem(REFRESH_TOKEN_KEY, 'valid-token')
    localStorage.setItem(TOKEN_EXPIRY_KEY, FUTURE_EXPIRY)

    const { result } = renderHook(() => useAuthContext(), {
      wrapper: AuthProvider,
    })

    // Espera restauração terminar (falha de rede: accessToken null, tokens preservados)
    await waitFor(() => {
      expect(result.current.isRestoringSession).toBe(false)
    })
    expect(result.current.isAuthenticated).toBe(false)
    expect(localStorage.getItem(REFRESH_TOKEN_KEY)).toBe('valid-token')

    // Segunda chamada: refreshIfNeeded com 401
    // Como accessToken === null, refreshIfNeeded vai chamar a API
    const unauthorized = new Response(null, { status: 401 })
    vi.spyOn(authApi, 'refreshToken').mockRejectedValue(unauthorized)

    await act(async () => {
      await result.current.refreshIfNeeded()
    })

    expect(localStorage.getItem(REFRESH_TOKEN_KEY)).toBeNull()
    expect(result.current.isAuthenticated).toBe(false)
  })
})
