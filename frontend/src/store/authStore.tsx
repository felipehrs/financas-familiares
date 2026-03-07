import { createContext, useContext, useState, useCallback, useEffect } from 'react'
import { login as apiLogin, refreshToken as apiRefreshToken } from '@/api/auth'

const REFRESH_TOKEN_KEY = 'refresh_token'
const TOKEN_EXPIRY_KEY = 'token_expiry'

interface AuthContextValue {
  accessToken: string | null
  isAuthenticated: boolean
  login: (email: string, senha: string) => Promise<void>
  logout: () => void
  refreshIfNeeded: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [accessToken, setAccessToken] = useState<string | null>(null)
  const [tokenExpiry, setTokenExpiry] = useState<number | null>(null)

  const logout = useCallback(() => {
    setAccessToken(null)
    setTokenExpiry(null)
    localStorage.removeItem(REFRESH_TOKEN_KEY)
    localStorage.removeItem(TOKEN_EXPIRY_KEY)
  }, [])

  const login = useCallback(async (email: string, senha: string) => {
    const response = await apiLogin({ email, senha })
    const expiry = Date.now() + response.expires_in * 1000
    setAccessToken(response.access_token)
    setTokenExpiry(expiry)
    localStorage.setItem(REFRESH_TOKEN_KEY, response.refresh_token)
    localStorage.setItem(TOKEN_EXPIRY_KEY, String(expiry))
  }, [])

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
      } catch {
        logout()
      }
    }
  }, [accessToken, tokenExpiry, logout])

  // Tenta restaurar sessão via refresh token ao inicializar
  useEffect(() => {
    const storedRefreshToken = localStorage.getItem(REFRESH_TOKEN_KEY)
    const storedExpiry = Number(localStorage.getItem(TOKEN_EXPIRY_KEY))

    if (storedRefreshToken && storedExpiry > Date.now()) {
      refreshIfNeeded().catch(() => logout())
    } else if (storedRefreshToken && storedExpiry <= Date.now()) {
      // Token expirado, limpa tudo
      logout()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return (
    <AuthContext.Provider
      value={{
        accessToken,
        isAuthenticated: accessToken !== null,
        login,
        logout,
        refreshIfNeeded,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuthContext(): AuthContextValue {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuthContext deve ser usado dentro de AuthProvider')
  }
  return context
}
