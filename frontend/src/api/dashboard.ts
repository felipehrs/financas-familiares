import type { ResumoMensal, ResumoCategorias, PontoEvolucao, MesProjecao } from '@/types/dashboard'

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

export async function buscarCategoriasDespesas(
  token: string,
  mes: number,
  ano: number,
): Promise<ResumoCategorias> {
  const response = await fetch(`${API_BASE}/api/v1/dashboard/categorias?mes=${mes}&ano=${ano}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<ResumoCategorias>(response)
}

export async function buscarEvolucaoMensal(token: string): Promise<PontoEvolucao[]> {
  const response = await fetch(`${API_BASE}/api/v1/dashboard/evolucao`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<PontoEvolucao[]>(response)
}

export async function buscarResumoMensal(
  token: string,
  mes: number,
  ano: number,
): Promise<ResumoMensal> {
  const response = await fetch(`${API_BASE}/api/v1/dashboard/resumo?mes=${mes}&ano=${ano}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<ResumoMensal>(response)
}

export async function buscarProjecao(token: string): Promise<MesProjecao[]> {
  const response = await fetch(`${API_BASE}/api/v1/dashboard/projecao`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<MesProjecao[]>(response)
}
