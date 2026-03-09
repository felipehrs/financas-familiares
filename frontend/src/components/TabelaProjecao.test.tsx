import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { TabelaProjecao } from './TabelaProjecao'
import type { MesProjecao } from '@/types/dashboard'

const fixture: MesProjecao[] = [
  { mes: 4, ano: 2026, total_rendas: 5000, total_cartoes: 800, total_assinaturas: 150, total_contas_fixas: 300, total_despesas: 1250, saldo_estimado: 3750 },
  { mes: 5, ano: 2026, total_rendas: 5000, total_cartoes: 700, total_assinaturas: 150, total_contas_fixas: 300, total_despesas: 1150, saldo_estimado: 3850 },
  { mes: 6, ano: 2026, total_rendas: 5000, total_cartoes: 600, total_assinaturas: 150, total_contas_fixas: 300, total_despesas: 1050, saldo_estimado: -500 },
]

describe('TabelaProjecao', () => {
  it('renderiza os 3 meses com labels corretos', () => {
    render(<TabelaProjecao projecao={fixture} />)
    expect(screen.getByText('ABR/26')).toBeInTheDocument()
    expect(screen.getByText('MAI/26')).toBeInTheDocument()
    expect(screen.getByText('JUN/26')).toBeInTheDocument()
  })

  it('exibe os valores de rendas formatados', () => {
    render(<TabelaProjecao projecao={fixture} />)
    expect(screen.getAllByText(/5\.000/).length).toBeGreaterThan(0)
  })

  it('exibe os subtotais de cartões, assinaturas e contas fixas', () => {
    render(<TabelaProjecao projecao={fixture} />)
    expect(screen.getByText(/800/)).toBeInTheDocument()
    expect(screen.getAllByText(/150/).length).toBeGreaterThan(0)
    expect(screen.getAllByText(/300/).length).toBeGreaterThan(0)
  })

  it('saldo positivo tem classe de cor verde', () => {
    render(<TabelaProjecao projecao={fixture} />)
    const saldos = screen.getAllByText(/3\.750|3\.850/)
    saldos.forEach(el => {
      expect(el.className).toMatch(/green/)
    })
  })

  it('saldo negativo tem classe de cor vermelha', () => {
    render(<TabelaProjecao projecao={fixture} />)
    // JUN/26 tem saldo_estimado: -500 → formatarMoeda produz algo como "-R$ 500,00"
    const saldoNegativo = screen.getByText(/500,00/)
    expect(saldoNegativo.className).toMatch(/red/)
  })

  it('exibe mensagem quando projecao está vazia', () => {
    render(<TabelaProjecao projecao={[]} />)
    expect(screen.getByText(/nenhum dado/i)).toBeInTheDocument()
  })
})
