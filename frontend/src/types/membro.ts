export interface Membro {
  id: string
  nome: string
  relacionamento: string
  ativo: boolean
}

export interface CriarMembroRequest {
  nome: string
  relacionamento?: string
}

export interface AtualizarMembroRequest {
  nome: string
  relacionamento?: string
  ativo: boolean
}
