import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { GraficoCategoriaDespesas } from './GraficoCategoriaDespesas'
import type { CategoriaDespesa } from '@/types/dashboard'

const categoriasFixture: CategoriaDespesa[] = [
  { nome: 'Alimentação', total: 650, percentual: 65 },
  { nome: 'Transporte', total: 350, percentual: 35 },
]

describe('GraficoCategoriaDespesas', () => {
  it('renderiza uma barra para cada categoria', () => {
    render(<GraficoCategoriaDespesas categorias={categoriasFixture} />)
    expect(screen.getAllByRole('progressbar')).toHaveLength(2)
  })

  it('exibe nome, valor formatado e percentual de cada categoria', () => {
    render(<GraficoCategoriaDespesas categorias={categoriasFixture} />)
    expect(screen.getByText('Alimentação')).toBeInTheDocument()
    expect(screen.getByText('Transporte')).toBeInTheDocument()
    expect(screen.getByText(/65/)).toBeInTheDocument()
    expect(screen.getByText(/35/)).toBeInTheDocument()
  })

  it('a progressbar tem aria-valuenow correto', () => {
    render(<GraficoCategoriaDespesas categorias={categoriasFixture} />)
    const bars = screen.getAllByRole('progressbar')
    expect(bars[0]).toHaveAttribute('aria-valuenow', '65')
    expect(bars[1]).toHaveAttribute('aria-valuenow', '35')
  })

  it('exibe mensagem quando categorias está vazio', () => {
    render(<GraficoCategoriaDespesas categorias={[]} />)
    expect(screen.getByText(/nenhuma despesa/i)).toBeInTheDocument()
  })
})
