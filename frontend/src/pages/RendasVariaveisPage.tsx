import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { listarRendasVariaveis, criarRendaVariavel, atualizarRendaVariavel, excluirRendaVariavel } from '@/offline/rendas_variaveis'
import { listarMembros } from '@/api/membros'
import type { RendaVariavel } from '@/types/renda_variavel'
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

const rendaVariavelSchema = z.object({
  descricao: z.string().min(2, 'Descrição é obrigatória e deve ter pelo menos 2 caracteres'),
  membro_id: z.string().min(1, 'Membro é obrigatório'),
  mes_referencia: z.coerce.number({ invalid_type_error: 'Informe o mês' }).int().min(1).max(12),
  ano_referencia: z.coerce.number({ invalid_type_error: 'Informe o ano' }).int().positive('Ano deve ser maior que zero'),
  valor: z.coerce.number({ invalid_type_error: 'Informe o valor' }).positive('Valor deve ser maior que zero'),
  data_recebimento: z.string().min(1, 'Data de recebimento é obrigatória'),
})

type RendaVariavelFormValues = z.infer<typeof rendaVariavelSchema>

// ─── Helpers ─────────────────────────────────────────────────────────────────

function formatarMoeda(valor: number): string {
  return valor.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })
}

function nomeMes(mesNumero: number): string {
  return MESES.find((m) => m.value === mesNumero)?.label ?? String(mesNumero)
}

// ─── Componente ───────────────────────────────────────────────────────────────

type FormMode = 'hidden' | 'criar' | 'editar'

export function RendasVariaveisPage() {
  const { accessToken } = useAuth()

  const hoje = new Date()

  const [rendas, setRendas] = useState<RendaVariavel[]>([])
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
  } = useForm<RendaVariavelFormValues>({
    resolver: zodResolver(rendaVariavelSchema),
    defaultValues: {
      descricao: '',
      membro_id: '',
      mes_referencia: '' as unknown as number,
      ano_referencia: '' as unknown as number,
      valor: '' as unknown as number,
      data_recebimento: '',
    },
  })

  // ─── Carrega dados ──────────────────────────────────────────────────────────

  async function carregarDados(mes?: number, ano?: number) {
    if (!accessToken) return
    try {
      const [listaRendas, listaMembros] = await Promise.all([
        listarRendasVariaveis(accessToken, mes, ano),
        listarMembros(accessToken),
      ])
      setRendas(listaRendas)
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
      mes_referencia: undefined,
      ano_referencia: undefined,
      valor: undefined,
      data_recebimento: '',
    })
  }

  function abrirFormEdicao(renda: RendaVariavel) {
    setFormMode('editar')
    setEditandoId(renda.id)
    setApiError(null)
    reset({
      descricao: renda.descricao,
      membro_id: renda.membro_id,
      mes_referencia: renda.mes_referencia,
      ano_referencia: renda.ano_referencia,
      valor: renda.valor,
      data_recebimento: renda.data_recebimento,
    })
  }

  function fecharForm() {
    setFormMode('hidden')
    setEditandoId(null)
    setApiError(null)
    reset({
      descricao: '',
      membro_id: '',
      mes_referencia: undefined,
      ano_referencia: undefined,
      valor: undefined,
      data_recebimento: '',
    })
  }

  // ─── Submit ─────────────────────────────────────────────────────────────────

  async function onSubmit(values: RendaVariavelFormValues) {
    if (!accessToken) return
    setApiError(null)

    try {
      if (formMode === 'criar') {
        await criarRendaVariavel(accessToken, {
          descricao: values.descricao,
          membro_id: values.membro_id,
          mes_referencia: values.mes_referencia,
          ano_referencia: values.ano_referencia,
          valor: values.valor,
          data_recebimento: values.data_recebimento,
        })
      } else if (formMode === 'editar' && editandoId) {
        await atualizarRendaVariavel(accessToken, editandoId, {
          descricao: values.descricao,
          membro_id: values.membro_id,
          mes_referencia: values.mes_referencia,
          ano_referencia: values.ano_referencia,
          valor: values.valor,
          data_recebimento: values.data_recebimento,
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
      await excluirRendaVariavel(accessToken, id)
      setLoading(true)
      await carregarDados(filtroMes, filtroAno)
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro ao excluir renda variável')
    }
  }

  // ─── Totais ───────────────────────────────────────────────────────────────

  const total = rendas.reduce((acc, r) => acc + r.valor, 0)

  // ─── Render ──────────────────────────────────────────────────────────────────

  return (
    <div className="p-6 max-w-2xl mx-auto">
      <div className="mb-4">
        <Link to="/" className="text-blue-600 dark:text-blue-400 hover:underline text-sm">
          ← Voltar para o dashboard
        </Link>
      </div>

      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Rendas Variáveis</h1>
        {formMode === 'hidden' && (
          <Button onClick={abrirFormCriacao}>Nova renda variável</Button>
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

      {/* Total das rendas filtradas */}
      <p className="mb-4 text-sm font-medium">
        Total: {formatarMoeda(total)} ({rendas.length} {rendas.length === 1 ? 'renda' : 'rendas'})
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
                  placeholder="Ex: Freelance, Comissão..."
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

              <div className="mb-4 grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="mes_referencia">Mês de referência</Label>
                  <select
                    id="mes_referencia"
                    {...register('mes_referencia')}
                    className={SELECT_CLASS}
                  >
                    <option value="">Selecione o mês</option>
                    {MESES.map((m) => (
                      <option key={m.value} value={m.value}>
                        {m.label}
                      </option>
                    ))}
                  </select>
                  {errors.mes_referencia && (
                    <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.mes_referencia.message}</p>
                  )}
                </div>
                <div>
                  <Label htmlFor="ano_referencia">Ano de referência</Label>
                  <Input
                    id="ano_referencia"
                    type="number"
                    min={1}
                    {...register('ano_referencia')}
                    placeholder="Ex: 2026"
                    className="mt-1"
                  />
                  {errors.ano_referencia && (
                    <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.ano_referencia.message}</p>
                  )}
                </div>
              </div>

              <div className="mb-4">
                <Label htmlFor="valor">Valor (R$)</Label>
                <Input
                  id="valor"
                  type="number"
                  min={0}
                  step="0.01"
                  {...register('valor')}
                  placeholder="Ex: 1500.00"
                  className="mt-1"
                />
                {errors.valor && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.valor.message}</p>
                )}
              </div>

              <div className="mb-4">
                <Label htmlFor="data_recebimento">Data de recebimento</Label>
                <Input
                  id="data_recebimento"
                  type="date"
                  {...register('data_recebimento')}
                  className="mt-1"
                />
                {errors.data_recebimento && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.data_recebimento.message}</p>
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

      {/* Lista de rendas variáveis */}
      {loading ? (
        <p className="text-muted-foreground">Carregando...</p>
      ) : rendas.length === 0 ? (
        <p className="text-muted-foreground">Nenhuma renda variável cadastrada</p>
      ) : (
        <ul className="space-y-3">
          {rendas.map((renda) => (
            <li key={renda.id}>
              <Card>
                <CardContent className="pt-4 flex items-center justify-between">
                  <div>
                    <p className="font-medium">{renda.descricao}</p>
                    <p className="text-sm text-muted-foreground">{nomeMembroById(renda.membro_id)}</p>
                    <p className="text-sm text-muted-foreground">
                      {formatarMoeda(renda.valor)}
                      {' · '}
                      {nomeMes(renda.mes_referencia)}/{renda.ano_referencia}
                      {' · '}
                      {renda.data_recebimento}
                    </p>
                  </div>
                  <div className="flex gap-2 flex-wrap justify-end">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => abrirFormEdicao(renda)}
                    >
                      Editar
                    </Button>
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => handleExcluir(renda.id)}
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
