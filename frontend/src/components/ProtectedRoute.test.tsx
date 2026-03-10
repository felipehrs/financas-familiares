import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter, Routes, Route } from 'react-router-dom'
import { ProtectedRoute } from './ProtectedRoute'

// Mock do useAuth para controlar isAuthenticated e isRestoringSession
vi.mock('@/hooks/useAuth')

// Mock do SyncQueueInitializer para evitar dependências externas
vi.mock('@/components/SyncQueueInitializer', () => ({
  SyncQueueInitializer: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}))

import { useAuth } from '@/hooks/useAuth'

const mockUseAuth = vi.mocked(useAuth)

function renderProtectedRoute() {
  return render(
    <MemoryRouter initialEntries={['/dashboard']}>
      <Routes>
        <Route path="/login" element={<div data-testid="login-page">Login</div>} />
        <Route element={<ProtectedRoute />}>
          <Route path="/dashboard" element={<div data-testid="dashboard">Dashboard</div>} />
        </Route>
      </Routes>
    </MemoryRouter>,
  )
}

describe('ProtectedRoute', () => {
  it('deve renderizar loading enquanto isRestoringSession = true', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isRestoringSession: true,
      accessToken: null,
      login: vi.fn(),
      logout: vi.fn(),
      refreshIfNeeded: vi.fn(),
    })

    renderProtectedRoute()

    expect(screen.getByText('Carregando...')).toBeTruthy()
  })

  it('NÃO deve redirecionar para /login enquanto isRestoringSession = true', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isRestoringSession: true,
      accessToken: null,
      login: vi.fn(),
      logout: vi.fn(),
      refreshIfNeeded: vi.fn(),
    })

    renderProtectedRoute()

    expect(screen.queryByTestId('login-page')).toBeNull()
  })

  it('deve redirecionar para /login quando isRestoringSession = false e isAuthenticated = false', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isRestoringSession: false,
      accessToken: null,
      login: vi.fn(),
      logout: vi.fn(),
      refreshIfNeeded: vi.fn(),
    })

    renderProtectedRoute()

    expect(screen.getByTestId('login-page')).toBeTruthy()
    expect(screen.queryByTestId('dashboard')).toBeNull()
  })

  it('deve renderizar Outlet quando isRestoringSession = false e isAuthenticated = true', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isRestoringSession: false,
      accessToken: 'some-token',
      login: vi.fn(),
      logout: vi.fn(),
      refreshIfNeeded: vi.fn(),
    })

    renderProtectedRoute()

    expect(screen.getByTestId('dashboard')).toBeTruthy()
    expect(screen.queryByTestId('login-page')).toBeNull()
  })
})
