import type { CategoriaDespesa } from '@/types/dashboard'

function formatarMoeda(valor: number) {
  return valor.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })
}

interface Props {
  categorias: CategoriaDespesa[]
}

export function GraficoCategoriaDespesas({ categorias }: Props) {
  if (categorias.length === 0) {
    return <p className="text-muted-foreground text-sm">Nenhuma despesa no período.</p>
  }

  return (
    <ul className="space-y-3">
      {categorias.map((cat) => (
        <li key={cat.nome}>
          <div className="flex justify-between text-sm mb-1">
            <span>{cat.nome}</span>
            <span className="text-muted-foreground">
              {formatarMoeda(cat.total)} ({cat.percentual.toFixed(1)}%)
            </span>
          </div>
          <div className="h-2 bg-muted rounded-full overflow-hidden">
            <div
              className="h-full bg-primary rounded-full"
              style={{ width: `${cat.percentual}%` }}
              role="progressbar"
              aria-valuenow={cat.percentual}
              aria-valuemin={0}
              aria-valuemax={100}
              aria-label={cat.nome}
            />
          </div>
        </li>
      ))}
    </ul>
  )
}
