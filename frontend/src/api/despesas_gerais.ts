import type { DespesaGeral, CriarDespesaGeralRequest, AtualizarDespesaGeralRequest } from '@/types/despesa_geral'

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

export async function listarDespesasGerais(
  token: string,
  mes?: number,
  ano?: number,
): Promise<DespesaGeral[]> {
  const params = new URLSearchParams()
  if (mes !== undefined) params.set('mes', String(mes))
  if (ano !== undefined) params.set('ano', String(ano))
  const query = params.toString() ? `?${params.toString()}` : ''
  const response = await fetch(`${API_BASE}/api/v1/despesas-gerais${query}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<DespesaGeral[]>(response)
}

export async function criarDespesaGeral(
  token: string,
  data: CriarDespesaGeralRequest,
): Promise<DespesaGeral> {
  const response = await fetch(`${API_BASE}/api/v1/despesas-gerais`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<DespesaGeral>(response)
}

export async function atualizarDespesaGeral(
  token: string,
  id: string,
  data: AtualizarDespesaGeralRequest,
): Promise<DespesaGeral> {
  const response = await fetch(`${API_BASE}/api/v1/despesas-gerais/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<DespesaGeral>(response)
}

export async function excluirDespesaGeral(token: string, id: string): Promise<void> {
  const response = await fetch(`${API_BASE}/api/v1/despesas-gerais/${id}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<void>(response)
}
