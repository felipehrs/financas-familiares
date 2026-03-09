import type { MesProjecao } from '@/types/dashboard'

const MESES = ['JAN', 'FEV', 'MAR', 'ABR', 'MAI', 'JUN', 'JUL', 'AGO', 'SET', 'OUT', 'NOV', 'DEZ']

const formatarMoeda = (v: number) =>
  v.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })

interface Props {
  projecao: MesProjecao[]
}

export function TabelaProjecao({ projecao }: Props) {
  if (projecao.length === 0) {
    return <p className="text-muted-foreground text-sm">Nenhum dado disponível.</p>
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b text-muted-foreground">
            <th className="text-left py-2 font-medium">Mês</th>
            <th className="text-right py-2 font-medium">Rendas Fixas</th>
            <th className="text-right py-2 font-medium">Cartões</th>
            <th className="text-right py-2 font-medium">Assinaturas</th>
            <th className="text-right py-2 font-medium">Contas Fixas</th>
            <th className="text-right py-2 font-medium">Total Despesas</th>
            <th className="text-right py-2 font-medium">Saldo Estimado</th>
          </tr>
        </thead>
        <tbody>
          {projecao.map((p) => {
            const label = `${MESES[p.mes - 1]}/${String(p.ano).slice(2)}`
            const saldoPositivo = p.saldo_estimado >= 0
            return (
              <tr key={label} className="border-b last:border-0">
                <td className="py-2 font-medium">{label}</td>
                <td className="py-2 text-right text-green-600 dark:text-green-400">
                  {formatarMoeda(p.total_rendas)}
                </td>
                <td className="py-2 text-right">{formatarMoeda(p.total_cartoes)}</td>
                <td className="py-2 text-right">{formatarMoeda(p.total_assinaturas)}</td>
                <td className="py-2 text-right">{formatarMoeda(p.total_contas_fixas)}</td>
                <td className="py-2 text-right text-red-600 dark:text-red-400">
                  {formatarMoeda(p.total_despesas)}
                </td>
                <td
                  className={`py-2 text-right font-semibold ${saldoPositivo ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}`}
                >
                  {formatarMoeda(p.saldo_estimado)}
                </td>
              </tr>
            )
          })}
        </tbody>
      </table>
    </div>
  )
}
