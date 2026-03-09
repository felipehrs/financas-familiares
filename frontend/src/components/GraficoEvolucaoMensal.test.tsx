import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect } from 'vitest'
import { GraficoEvolucaoMensal } from './GraficoEvolucaoMensal'
import type { PontoEvolucao } from '@/types/dashboard'

// 12 meses de ABR/2025 a MAR/2026
const fixture: PontoEvolucao[] = Array.from({ length: 12 }, (_, i) => {
  const data = new Date(2025, 3 + i, 1) // ABR/2025 + i meses
  return {
    mes: data.getMonth() + 1,
    ano: data.getFullYear(),
    total_rendas: 5000,
    total_despesas: 3000,
    saldo: 2000,
  }
})

describe('GraficoEvolucaoMensal', () => {
  it('exibe mensagem quando pontos está vazio', () => {
    render(<GraficoEvolucaoMensal pontos={[]} />)
    expect(screen.getByText(/nenhum dado/i)).toBeInTheDocument()
  })

  it('renderiza o SVG com role img', () => {
    render(<GraficoEvolucaoMensal pontos={fixture} />)
    expect(screen.getByRole('img')).toBeInTheDocument()
  })

  it('renderiza 3 polylines (rendas, despesas, saldo)', () => {
    const { container } = render(<GraficoEvolucaoMensal pontos={fixture} />)
    const polylines = container.querySelectorAll('polyline')
    expect(polylines).toHaveLength(3)
  })

  it('exibe label do primeiro e último mês', () => {
    render(<GraficoEvolucaoMensal pontos={fixture} />)
    // fixture começa em ABR/2025, termina em MAR/2026
    expect(screen.getAllByText('Abr').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Mar').length).toBeGreaterThan(0)
  })

  it('exibe legenda com Rendas, Despesas e Saldo', () => {
    render(<GraficoEvolucaoMensal pontos={fixture} />)
    expect(screen.getByText('Rendas')).toBeInTheDocument()
    expect(screen.getByText('Despesas')).toBeInTheDocument()
    expect(screen.getByText('Saldo')).toBeInTheDocument()
  })

  it('mostra tooltip ao hover em um círculo', async () => {
    const user = userEvent.setup()
    render(<GraficoEvolucaoMensal pontos={fixture} />)
    const circles = document.querySelectorAll('circle')
    await user.hover(circles[0])
    // Tooltip deve mostrar os valores do primeiro ponto
    expect(screen.getByText(/5\.000/)).toBeInTheDocument()
  })
})
