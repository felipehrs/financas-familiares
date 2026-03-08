import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { listarRendimentos, criarRendimento, atualizarRendimento, excluirRendimento } from '@/offline/rendimentos_investimento'
import { listarMembros } from '@/api/membros'
import type { RendimentoInvestimento } from '@/types/rendimento_investimento'
import type { Membro } from '@/types/membro'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

// ─── Constantes ───────────────────────────────────────────────────────────────

const MESES = [
  { value: 1, label: 'Janeiro' },
  { value: 2, label: 'Fevereiro' },
  { value: 3, label: 'Março' },
  { value: 4, label: 'Abril' },
  { value: 5, label: 'Maio' },
  { value: 6, label: 'Junho' },
  { value: 7, label: 'Julho' },
  { value: 8, label: 'Agosto' },
  { value: 9, label: 'Setembro' },
  { value: 10, label: 'Outubro' },
  { value: 11, label: 'Novembro' },
  { value: 12, label: 'Dezembro' },
]

const SELECT_CLASS = 'mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2'

// ─── Schema de validação ──────────────────────────────────────────────────────

const rendimentoSchema = z.object({
  descricao: z.string().min(2, 'Descrição é obrigatória e deve ter pelo menos 2 caracteres'),
  membro_id: z.string().min(1, 'Membro é obrigatório'),
  data: z.string().min(1, 'Data é obrigatória'),
  valor: z.coerce.number({ invalid_type_error: 'Informe o valor' }).positive('Valor deve ser maior que zero'),
  valor_distribuido: z.coerce.number().min(0, 'Mínimo 0').default(0),
}).refine((data) => data.valor_distribuido <= data.valor, {
  message: 'Valor distribuído não pode ser maior que o valor total recebido',
  path: ['valor_distribuido'],
})

type RendimentoFormValues = z.infer<typeof rendimentoSchema>

// ─── Helpers ─────────────────────────────────────────────────────────────────

function formatarMoeda(valor: number): string {
  return valor.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })
}

// ─── Componente ───────────────────────────────────────────────────────────────

type FormMode = 'hidden' | 'criar' | 'editar'

export function RendimentosInvestimentoPage() {
  const { accessToken } = useAuth()

  const hoje = new Date()

  const [rendimentos, setRendimentos] = useState<RendimentoInvestimento[]>([])
  const [membros, setMembros] = useState<Membro[]>([])
  const [loading, setLoading] = useState(true)
  const [apiError, setApiError] = useState<string | null>(null)
  const [formMode, setFormMode] = useState<FormMode>('hidden')
  const [editandoId, setEditandoId] = useState<string | null>(null)

  // ─── Estado do filtro de mês/ano ──────────────────────────────────────────
  const [filtroMes, setFiltroMes] = useState(hoje.getMonth() + 1)
  const [filtroAno, setFiltroAno] = useState(hoje.getFullYear())

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<RendimentoFormValues>({
    resolver: zodResolver(rendimentoSchema),
    defaultValues: {
      descricao: '',
      membro_id: '',
      data: '',
      valor: '' as unknown as number,
      valor_distribuido: 0,
    },
  })

  // ─── Carrega dados ──────────────────────────────────────────────────────────

  async function carregarDados(mes?: number, ano?: number) {
    if (!accessToken) return
    try {
      const [listaRendimentos, listaMembros] = await Promise.all([
        listarRendimentos(accessToken, mes, ano),
        listarMembros(accessToken),
      ])
      setRendimentos(listaRendimentos)
      setMembros(listaMembros)
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro ao carregar dados')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void carregarDados(filtroMes, filtroAno)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filtroMes, filtroAno])

  // ─── Helpers ────────────────────────────────────────────────────────────────

  function nomeMembroById(membroId: string): string {
    return membros.find((m) => m.id === membroId)?.nome ?? membroId
  }

  // ─── Handlers de UI ─────────────────────────────────────────────────────────

  function abrirFormCriacao() {
    setFormMode('criar')
    setEditandoId(null)
    setApiError(null)
    reset({
      descricao: '',
      membro_id: '',
      data: '',
      valor: undefined,
      valor_distribuido: 0,
    })
  }

  function abrirFormEdicao(rendimento: RendimentoInvestimento) {
    setFormMode('editar')
    setEditandoId(rendimento.id)
    setApiError(null)
    reset({
      descricao: rendimento.descricao,
      membro_id: rendimento.membro_id,
      data: rendimento.data,
      valor: rendimento.valor,
      valor_distribuido: rendimento.valor_distribuido,
    })
  }

  function fecharForm() {
    setFormMode('hidden')
    setEditandoId(null)
    setApiError(null)
    reset({
      descricao: '',
      membro_id: '',
      data: '',
      valor: undefined,
      valor_distribuido: 0,
    })
  }

  // ─── Submit ─────────────────────────────────────────────────────────────────

  async function onSubmit(values: RendimentoFormValues) {
    if (!accessToken) return
    setApiError(null)

    try {
      if (formMode === 'criar') {
        await criarRendimento(accessToken, {
          descricao: values.descricao,
          membro_id: values.membro_id,
          data: values.data,
          valor: values.valor,
          valor_distribuido: values.valor_distribuido,
        })
      } else if (formMode === 'editar' && editandoId) {
        await atualizarRendimento(accessToken, editandoId, {
          descricao: values.descricao,
          membro_id: values.membro_id,
          data: values.data,
          valor: values.valor,
          valor_distribuido: values.valor_distribuido,
        })
      }

      fecharForm()
      setLoading(true)
      await carregarDados(filtroMes, filtroAno)
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro desconhecido')
    }
  }

  async function handleExcluir(id: string) {
    if (!accessToken) return
    setApiError(null)
    try {
      await excluirRendimento(accessToken, id)
      setLoading(true)
      await carregarDados(filtroMes, filtroAno)
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro ao excluir rendimento')
    }
  }

  // ─── Totais ───────────────────────────────────────────────────────────────

  const totalRecebido = rendimentos.reduce((acc, r) => acc + r.valor, 0)
  const totalDistribuido = rendimentos.reduce((acc, r) => acc + r.valor_distribuido, 0)
  const totalReinvestido = totalRecebido - totalDistribuido

  // ─── Render ──────────────────────────────────────────────────────────────────

  return (
    <div className="p-6 max-w-2xl mx-auto">
      <div className="mb-4">
        <Link to="/" className="text-blue-600 dark:text-blue-400 hover:underline text-sm">
          ← Voltar para o dashboard
        </Link>
      </div>

      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Rendimentos de Investimento</h1>
        {formMode === 'hidden' && (
          <Button onClick={abrirFormCriacao}>Novo rendimento</Button>
        )}
      </div>

      {/* Erro global de API */}
      {apiError && (
        <p className="mb-4 text-sm text-red-600 dark:text-red-400" role="alert">
          {apiError}
        </p>
      )}

      {/* Filtro de mês/ano */}
      <div className="mb-6 flex flex-wrap gap-3 items-end">
        <div>
          <label htmlFor="filtro-mes" className="block text-sm font-medium mb-1">
            Mês
          </label>
          <select
            id="filtro-mes"
            aria-label="Mês"
            value={filtroMes}
            onChange={(e) => setFiltroMes(Number(e.target.value))}
            className="mt-1 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
          >
            {MESES.map((m) => (
              <option key={m.value} value={m.value}>
                {m.label}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label htmlFor="filtro-ano" className="block text-sm font-medium mb-1">
            Ano
          </label>
          <input
            id="filtro-ano"
            aria-label="Ano"
            type="number"
            value={filtroAno}
            onChange={(e) => setFiltroAno(Number(e.target.value))}
            className="mt-1 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 w-28"
          />
        </div>
      </div>

      {/* Resumo de totais */}
      <p className="mb-4 text-sm font-medium">
        Total recebido: {formatarMoeda(totalRecebido)} | Total distribuído: {formatarMoeda(totalDistribuido)} | Total reinvestido: {formatarMoeda(totalReinvestido)}
      </p>

      {/* Formulário inline de criação / edição */}
      {formMode !== 'hidden' && (
        <Card className="mb-6">
          <CardContent className="pt-4">
            <form onSubmit={handleSubmit(onSubmit)} noValidate>
              <div className="mb-4">
                <Label htmlFor="descricao">Descrição</Label>
                <Input
                  id="descricao"
                  {...register('descricao')}
                  placeholder="Ex: Dividendos, Juros sobre capital..."
                  className="mt-1"
                />
                {errors.descricao && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.descricao.message}</p>
                )}
              </div>

              <div className="mb-4">
                <Label htmlFor="membro_id">Membro</Label>
                <select
                  id="membro_id"
                  {...register('membro_id')}
                  className={SELECT_CLASS}
                >
                  <option value="">Selecione um membro</option>
                  {membros.filter((m) => m.ativo).map((membro) => (
                    <option key={membro.id} value={membro.id}>
                      {membro.nome}
                    </option>
                  ))}
                </select>
                {errors.membro_id && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.membro_id.message}</p>
                )}
              </div>

              <div className="mb-4">
                <Label htmlFor="data">Data</Label>
                <Input
                  id="data"
                  type="date"
                  {...register('data')}
                  className="mt-1"
                />
                {errors.data && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.data.message}</p>
                )}
              </div>

              <div className="mb-4">
                <Label htmlFor="valor">Valor total recebido (R$)</Label>
                <Input
                  id="valor"
                  type="number"
                  min={0}
                  step="0.01"
                  {...register('valor')}
                  placeholder="Ex: 500.00"
                  className="mt-1"
                />
                {errors.valor && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.valor.message}</p>
                )}
              </div>

              <div className="mb-4">
                <Label htmlFor="valor_distribuido">Valor distribuído (R$)</Label>
                <Input
                  id="valor_distribuido"
                  type="number"
                  min={0}
                  step="0.01"
                  {...register('valor_distribuido')}
                  placeholder="Ex: 200.00"
                  className="mt-1"
                />
                {errors.valor_distribuido && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.valor_distribuido.message}</p>
                )}
              </div>

              <div className="flex gap-2">
                <Button type="submit" disabled={isSubmitting}>
                  {isSubmitting ? 'Salvando...' : 'Salvar'}
                </Button>
                <Button type="button" variant="outline" onClick={fecharForm}>
                  Cancelar
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {/* Lista de rendimentos */}
      {loading ? (
        <p className="text-muted-foreground">Carregando...</p>
      ) : rendimentos.length === 0 ? (
        <p className="text-muted-foreground">Nenhum rendimento cadastrado</p>
      ) : (
        <ul className="space-y-3">
          {rendimentos.map((rendimento) => (
            <li key={rendimento.id}>
              <Card>
                <CardContent className="pt-4 flex items-center justify-between">
                  <div>
                    <p className="font-medium">{rendimento.descricao}</p>
                    <p className="text-sm text-muted-foreground">{nomeMembroById(rendimento.membro_id)}</p>
                    <p className="text-sm text-muted-foreground">
                      {formatarMoeda(rendimento.valor)}
                      {' · '}
                      Distribuído: {formatarMoeda(rendimento.valor_distribuido)}
                      {' · '}
                      {rendimento.data}
                    </p>
                    <span className="inline-block mt-1 text-xs bg-muted px-2 py-0.5 rounded">
                      Reinvestido: {formatarMoeda(rendimento.valor - rendimento.valor_distribuido)}
                    </span>
                  </div>
                  <div className="flex gap-2 flex-wrap justify-end">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => abrirFormEdicao(rendimento)}
                    >
                      Editar
                    </Button>
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => handleExcluir(rendimento.id)}
                    >
                      Excluir
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
