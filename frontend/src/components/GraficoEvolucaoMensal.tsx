import { useState } from 'react'
import type { PontoEvolucao } from '@/types/dashboard'

const MESES = ['Jan', 'Fev', 'Mar', 'Abr', 'Mai', 'Jun', 'Jul', 'Ago', 'Set', 'Out', 'Nov', 'Dez']

const formatarMoeda = (v: number) =>
  v.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })

const formatarEixoY = (v: number) =>
  Math.abs(v) >= 1000
    ? `R$${(v / 1000).toFixed(0)}k`
    : `R$${v.toFixed(0)}`

interface Props {
  pontos: PontoEvolucao[]
}

export function GraficoEvolucaoMensal({ pontos }: Props) {
  const [hovered, setHovered] = useState<number | null>(null)

  if (pontos.length === 0) {
    return <p className="text-muted-foreground text-sm">Nenhum dado para exibir.</p>
  }

  const PAD = { top: 20, right: 20, bottom: 40, left: 60 }
  const W = 600
  const H = 220
  const uw = W - PAD.left - PAD.right
  const uh = H - PAD.top - PAD.bottom

  const maxVal = Math.max(...pontos.flatMap(p => [p.total_rendas, p.total_despesas, 0]))
  const minVal = Math.min(...pontos.map(p => p.saldo), 0)
  const range = maxVal - minVal || 1

  const cx = (i: number) => PAD.left + (i / (pontos.length - 1)) * uw
  const cy = (v: number) => PAD.top + (1 - (v - minVal) / range) * uh

  const polyline = (getter: (p: PontoEvolucao) => number) =>
    pontos.map((p, i) => `${cx(i)},${cy(getter(p))}`).join(' ')

  const yMarks = Array.from({ length: 5 }, (_, i) => minVal + (range / 4) * i).reverse()

  return (
    <div>
      <svg
        viewBox={`0 0 ${W} ${H}`}
        width="100%"
        aria-label="Gráfico de evolução mensal"
        role="img"
      >
        {/* Grid horizontal */}
        {yMarks.map((v, i) => (
          <g key={i}>
            <line
              x1={PAD.left} y1={cy(v)}
              x2={W - PAD.right} y2={cy(v)}
              stroke="currentColor" strokeOpacity={0.1} strokeWidth={1}
            />
            <text
              x={PAD.left - 6} y={cy(v)}
              textAnchor="end" dominantBaseline="middle"
              fontSize={10} fill="currentColor" opacity={0.5}
            >
              {formatarEixoY(v)}
            </text>
          </g>
        ))}

        {/* Linha do zero se saldo pode ser negativo */}
        {minVal < 0 && (
          <line
            x1={PAD.left} y1={cy(0)}
            x2={W - PAD.right} y2={cy(0)}
            stroke="currentColor" strokeOpacity={0.3} strokeWidth={1} strokeDasharray="4 2"
          />
        )}

        {/* Polylines */}
        <polyline points={polyline(p => p.total_rendas)}
          fill="none" stroke="#22c55e" strokeWidth={2} strokeLinejoin="round" />
        <polyline points={polyline(p => p.total_despesas)}
          fill="none" stroke="#ef4444" strokeWidth={2} strokeLinejoin="round" />
        <polyline points={polyline(p => p.saldo)}
          fill="none" stroke="#3b82f6" strokeWidth={2} strokeLinejoin="round" strokeDasharray="5 3" />

        {/* Pontos e labels do eixo X */}
        {pontos.map((p, i) => (
          <g key={i}>
            {/* Label mês */}
            <text
              x={cx(i)} y={H - PAD.bottom + 14}
              textAnchor="middle" fontSize={10}
              fill="currentColor" opacity={0.6}
            >
              {MESES[p.mes - 1]}
            </text>

            {/* Círculos interativos (saldo) */}
            <circle
              cx={cx(i)} cy={cy(p.saldo)}
              r={hovered === i ? 5 : 3}
              fill="#3b82f6"
              onMouseEnter={() => setHovered(i)}
              onMouseLeave={() => setHovered(null)}
              style={{ cursor: 'pointer' }}
              aria-label={`${MESES[p.mes - 1]} ${p.ano}`}
            />
          </g>
        ))}

        {/* Tooltip */}
        {hovered !== null && (() => {
          const p = pontos[hovered]
          const tx = Math.min(cx(hovered), W - PAD.right - 130)
          const ty = Math.max(cy(p.saldo) - 70, PAD.top)
          return (
            <g>
              <rect x={tx} y={ty} width={130} height={64} rx={4}
                fill="var(--card)" stroke="var(--border)" strokeWidth={1} opacity={0.95} />
              <text x={tx + 8} y={ty + 14} fontSize={10} fontWeight="600" fill="var(--card-foreground)">
                {MESES[p.mes - 1]}/{p.ano}
              </text>
              <text x={tx + 8} y={ty + 28} fontSize={10} fill="#22c55e">
                Rendas: {formatarMoeda(p.total_rendas)}
              </text>
              <text x={tx + 8} y={ty + 42} fontSize={10} fill="#ef4444">
                Despesas: {formatarMoeda(p.total_despesas)}
              </text>
              <text x={tx + 8} y={ty + 56} fontSize={10} fill="#3b82f6">
                Saldo: {formatarMoeda(p.saldo)}
              </text>
            </g>
          )
        })()}
      </svg>

      {/* Legenda */}
      <div className="flex gap-4 justify-center mt-2 text-xs text-muted-foreground">
        <span className="flex items-center gap-1">
          <span className="inline-block w-4 h-0.5 bg-green-500 rounded" />
          Rendas
        </span>
        <span className="flex items-center gap-1">
          <span className="inline-block w-4 h-0.5 bg-red-500 rounded" />
          Despesas
        </span>
        <span className="flex items-center gap-1">
          <span className="inline-block w-4 h-0.5 bg-blue-500 rounded" style={{ backgroundImage: 'repeating-linear-gradient(90deg, #3b82f6 0,#3b82f6 5px,transparent 5px,transparent 8px)' }} />
          Saldo
        </span>
      </div>
    </div>
  )
}
