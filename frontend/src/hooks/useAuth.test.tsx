import { renderHook, act } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useAuth } from './useAuth'
import { AuthProvider } from '@/store/authStore'

// Mock da API de autenticação
vi.mock('@/api/auth', () => ({
  login: vi.fn(),
  refreshToken: vi.fn(),
}))

import { login as mockLogin } from '@/api/auth'

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <AuthProvider>{children}</AuthProvider>
)

describe('useAuth', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('isAuthenticated é false inicialmente', () => {
    const { result } = renderHook(() => useAuth(), { wrapper })
    expect(result.current.isAuthenticated).toBe(false)
    expect(result.current.accessToken).toBeNull()
  })

  it('isAuthenticated é true após login bem-sucedido', async () => {
    vi.mocked(mockLogin).mockResolvedValueOnce({
      access_token: 'access-token-123',
      refresh_token: 'refresh-token-456',
      token_type: 'Bearer',
      expires_in: 3600,
    })

    const { result } = renderHook(() => useAuth(), { wrapper })

    await act(async () => {
      await result.current.login('user@example.com', 'senha123')
    })

    expect(result.current.isAuthenticated).toBe(true)
    expect(result.current.accessToken).toBe('access-token-123')
  })

  it('logout limpa o estado', async () => {
    vi.mocked(mockLogin).mockResolvedValueOnce({
      access_token: 'access-token-123',
      refresh_token: 'refresh-token-456',
      token_type: 'Bearer',
      expires_in: 3600,
    })

    const { result } = renderHook(() => useAuth(), { wrapper })

    await act(async () => {
      await result.current.login('user@example.com', 'senha123')
    })

    expect(result.current.isAuthenticated).toBe(true)

    act(() => {
      result.current.logout()
    })

    expect(result.current.isAuthenticated).toBe(false)
    expect(result.current.accessToken).toBeNull()
  })

  it('login salva refresh_token no localStorage', async () => {
    vi.mocked(mockLogin).mockResolvedValueOnce({
      access_token: 'access-token-123',
      refresh_token: 'refresh-token-456',
      token_type: 'Bearer',
      expires_in: 3600,
    })

    const { result } = renderHook(() => useAuth(), { wrapper })

    await act(async () => {
      await result.current.login('user@example.com', 'senha123')
    })

    expect(localStorage.getItem('refresh_token')).toBe('refresh-token-456')
  })

  it('logout remove refresh_token do localStorage', async () => {
    vi.mocked(mockLogin).mockResolvedValueOnce({
      access_token: 'access-token-123',
      refresh_token: 'refresh-token-456',
      token_type: 'Bearer',
      expires_in: 3600,
    })

    const { result } = renderHook(() => useAuth(), { wrapper })

    await act(async () => {
      await result.current.login('user@example.com', 'senha123')
    })

    act(() => {
      result.current.logout()
    })

    expect(localStorage.getItem('refresh_token')).toBeNull()
  })
})
