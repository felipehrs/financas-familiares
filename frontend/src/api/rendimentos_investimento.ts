import type { RendimentoInvestimento, CriarRendimentoInvestimentoRequest, AtualizarRendimentoInvestimentoRequest } from '@/types/rendimento_investimento'

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

export async function listarRendimentos(
  token: string,
  mes?: number,
  ano?: number,
): Promise<RendimentoInvestimento[]> {
  const params = new URLSearchParams()
  if (mes !== undefined) params.set('mes', String(mes))
  if (ano !== undefined) params.set('ano', String(ano))
  const query = params.toString() ? `?${params.toString()}` : ''
  const response = await fetch(`${API_BASE}/api/v1/rendimentos-investimento${query}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<RendimentoInvestimento[]>(response)
}

export async function criarRendimento(
  token: string,
  data: CriarRendimentoInvestimentoRequest,
): Promise<RendimentoInvestimento> {
  const response = await fetch(`${API_BASE}/api/v1/rendimentos-investimento`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<RendimentoInvestimento>(response)
}

export async function atualizarRendimento(
  token: string,
  id: string,
  data: AtualizarRendimentoInvestimentoRequest,
): Promise<RendimentoInvestimento> {
  const response = await fetch(`${API_BASE}/api/v1/rendimentos-investimento/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<RendimentoInvestimento>(response)
}

export async function excluirRendimento(token: string, id: string): Promise<void> {
  const response = await fetch(`${API_BASE}/api/v1/rendimentos-investimento/${id}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<void>(response)
}
