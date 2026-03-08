import type { DespesaCartao, CriarDespesaCartaoRequest } from '@/types/despesa_cartao'

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

export async function listarDespesasPorCartao(token: string, cartaoId: string): Promise<DespesaCartao[]> {
  const response = await fetch(`${API_BASE}/api/v1/cartoes/${cartaoId}/despesas`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<DespesaCartao[]>(response)
}

export async function listarDespesasPorFatura(
  token: string,
  cartaoId: string,
  mes: number,
  ano: number,
): Promise<DespesaCartao[]> {
  const response = await fetch(`${API_BASE}/api/v1/cartoes/${cartaoId}/despesas/fatura?mes=${mes}&ano=${ano}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<DespesaCartao[]>(response)
}

export async function criarDespesa(
  token: string,
  cartaoId: string,
  data: CriarDespesaCartaoRequest,
): Promise<DespesaCartao[]> {
  const response = await fetch(`${API_BASE}/api/v1/cartoes/${cartaoId}/despesas`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<DespesaCartao[]>(response)
}

export async function excluirDespesa(token: string, id: string): Promise<void> {
  const response = await fetch(`${API_BASE}/api/v1/despesas/${id}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<void>(response)
}
