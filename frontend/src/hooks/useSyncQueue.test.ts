import { renderHook } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

vi.mock('@/hooks/useAuth', () => ({
  useAuth: vi.fn(() => ({ accessToken: 'token-teste' })),
}))

vi.mock('@/lib/syncQueue', () => ({
  processarFila: vi.fn().mockResolvedValue(undefined),
  configurarListenerOnline: vi.fn(),
}))

import { useAuth } from '@/hooks/useAuth'
import { processarFila, configurarListenerOnline } from '@/lib/syncQueue'
import { useSyncQueue } from './useSyncQueue'

describe('useSyncQueue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    Object.defineProperty(navigator, 'onLine', {
      value: true,
      writable: true,
      configurable: true,
    })
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('chama configurarListenerOnline ao montar quando há token', () => {
    renderHook(() => useSyncQueue())
    expect(configurarListenerOnline).toHaveBeenCalledOnce()
  })

  it('passa um getter de token para configurarListenerOnline', () => {
    renderHook(() => useSyncQueue())
    const getter = vi.mocked(configurarListenerOnline).mock.calls[0][0]
    expect(typeof getter).toBe('function')
    expect(getter()).toBe('token-teste')
  })

  it('chama processarFila ao montar quando navigator.onLine é true', () => {
    renderHook(() => useSyncQueue())
    expect(processarFila).toHaveBeenCalledWith('token-teste')
  })

  it('não chama processarFila ao montar quando navigator.onLine é false', () => {
    Object.defineProperty(navigator, 'onLine', {
      value: false,
      writable: true,
      configurable: true,
    })
    renderHook(() => useSyncQueue())
    expect(processarFila).not.toHaveBeenCalled()
  })

  it('chama configurarListenerOnline mesmo quando navigator.onLine é false', () => {
    Object.defineProperty(navigator, 'onLine', {
      value: false,
      writable: true,
      configurable: true,
    })
    renderHook(() => useSyncQueue())
    expect(configurarListenerOnline).toHaveBeenCalledOnce()
  })

  it('não chama configurarListenerOnline nem processarFila quando accessToken é null', () => {
    vi.mocked(useAuth).mockReturnValue({
      accessToken: null,
      isAuthenticated: false,
      isRestoringSession: false,
      login: vi.fn(),
      logout: vi.fn(),
      refreshIfNeeded: vi.fn(),
    })
    renderHook(() => useSyncQueue())
    expect(configurarListenerOnline).not.toHaveBeenCalled()
    expect(processarFila).not.toHaveBeenCalled()
  })
})
