import { useEffect } from 'react'
import { useAuth } from '@/hooks/useAuth'
import { configurarListenerOnline, processarFila } from '@/lib/syncQueue'

export function useSyncQueue() {
  const { accessToken } = useAuth()

  useEffect(() => {
    if (!accessToken) return

    configurarListenerOnline(() => accessToken)

    // Se já estiver online ao montar, processa fila residual
    if (navigator.onLine) {
      void processarFila(accessToken)
    }
  }, [accessToken])
}
