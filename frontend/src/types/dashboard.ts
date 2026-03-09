export interface PontoEvolucao {
  mes: number
  ano: number
  total_rendas: number
  total_despesas: number
  saldo: number
}

export interface CategoriaDespesa {
  nome: string
  total: number
  percentual: number
}

export interface ResumoCategorias {
  mes: number
  ano: number
  total_despesas: number
  categorias: CategoriaDespesa[]
}

export interface ResumoMensal {
  mes: number
  ano: number
  // Rendas operacionais (entram no saldo — RN06)
  total_renda_fixa: number
  total_renda_variavel: number
  total_renda_extra: number
  total_rendimento_distribuido: number
  total_rendas_operacionais: number
  // Rendimentos informativos (não entram no saldo — RN08)
  total_rendimento_investimento: number
  // Despesas
  total_fatura_cartoes: number
  total_assinaturas: number
  total_contas_fixas: number
  total_despesas_gerais: number
  total_despesas: number
  // Resultado
  saldo: number
}
