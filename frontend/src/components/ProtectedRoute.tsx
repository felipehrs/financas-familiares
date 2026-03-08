import { Navigate, Outlet } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { SyncQueueInitializer } from '@/components/SyncQueueInitializer'

export function ProtectedRoute() {
  const { isAuthenticated } = useAuth()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  return (
    <SyncQueueInitializer>
      <Outlet />
    </SyncQueueInitializer>
  )
}
