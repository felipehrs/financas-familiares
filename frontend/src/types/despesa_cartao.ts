export interface DespesaCartao {
  id: string
  compra_id: string
  cartao_id: string
  categoria_id: string | null
  descricao: string
  data_compra: string       // "YYYY-MM-DD"
  valor_total: number
  numero_parcelas: number
  parcela_numero: number
  valor_parcela: number
  fatura_mes: number
  fatura_ano: number
  fatura: string            // "MAR/26" — calculado/retornado pelo backend
}

export interface CriarDespesaCartaoRequest {
  descricao: string
  categoria_id?: string
  data_compra: string       // "YYYY-MM-DD"
  valor_total: number
  numero_parcelas: number   // mínimo 1; 1 = à vista
}
