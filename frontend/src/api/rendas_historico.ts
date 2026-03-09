import type { HistoricoRendas, FiltroHistoricoRendas } from '@/types/renda_historico'

const API_BASE = import.meta.env.VITE_API_URL ?? ''

async function handleResponse<T>(response: Response): Promise<T> {
  if (response.status === 401) {
    throw new Error('Não autorizado')
  }

  if (!response.ok) {
    let message = response.statusText
    try {
      const body = await response.json() as { error?: string }
      if (body.error) {
        message = body.error
      }
    } catch {
      // mantém statusText se body não for JSON válido
    }
    throw new Error(message)
  }

  if (response.status === 204) {
    return undefined as unknown as T
  }

  return response.json() as Promise<T>
}

export async function buscarHistoricoRendas(
  token: string,
  filtro: FiltroHistoricoRendas = {},
): Promise<HistoricoRendas> {
  const params = new URLSearchParams()
  if (filtro.tipo) params.set('tipo', filtro.tipo)
  if (filtro.membro_id) params.set('membro_id', filtro.membro_id)
  if (filtro.mes !== undefined && filtro.mes !== 0) params.set('mes', String(filtro.mes))
  if (filtro.ano !== undefined && filtro.ano !== 0) params.set('ano', String(filtro.ano))
  const query = params.toString() ? `?${params.toString()}` : ''
  const response = await fetch(`${API_BASE}/api/v1/rendas/historico${query}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<HistoricoRendas>(response)
}
