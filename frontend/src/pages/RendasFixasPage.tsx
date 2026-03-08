import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { listarRendasFixas, criarRendaFixa, atualizarRendaFixa, inativarRendaFixa } from '@/offline/rendas_fixas'
import { listarMembros } from '@/api/membros'
import type { RendaFixa, AtualizarRendaFixaRequest } from '@/types/renda_fixa'
import type { Membro } from '@/types/membro'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

// ─── Schema de validação ──────────────────────────────────────────────────────

const rendaFixaSchema = z.object({
  descricao: z.string().min(2, 'Descrição é obrigatória e deve ter pelo menos 2 caracteres'),
  membro_id: z.string().min(1, 'Membro responsável é obrigatório'),
  valor: z.coerce.number({ invalid_type_error: 'Informe o valor' }).positive('Valor deve ser positivo'),
  dia_recebimento: z.coerce
    .number({ invalid_type_error: 'Informe o dia de recebimento' })
    .int()
    .min(1, 'Deve ser entre 1 e 31')
    .max(31, 'Deve ser entre 1 e 31'),
  data_inicio: z.string().min(1, 'Data de início é obrigatória'),
  data_fim: z.string().optional(),
})

type RendaFixaFormValues = z.infer<typeof rendaFixaSchema>

// ─── Componente ───────────────────────────────────────────────────────────────

type FormMode = 'hidden' | 'criar' | 'editar'

export function RendasFixasPage() {
  const { accessToken } = useAuth()

  const [rendas, setRendas] = useState<RendaFixa[]>([])
  const [membros, setMembros] = useState<Membro[]>([])
  const [loading, setLoading] = useState(true)
  const [apiError, setApiError] = useState<string | null>(null)
  const [formMode, setFormMode] = useState<FormMode>('hidden')
  const [editandoId, setEditandoId] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<RendaFixaFormValues>({
    resolver: zodResolver(rendaFixaSchema),
    defaultValues: { descricao: '', membro_id: '', valor: '' as unknown as number, dia_recebimento: '' as unknown as number, data_inicio: '', data_fim: '' },
  })

  // ─── Carrega dados ──────────────────────────────────────────────────────────

  async function carregarDados() {
    if (!accessToken) return
    try {
      const [listaRendas, listaMembros] = await Promise.all([
        listarRendasFixas(accessToken),
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
    void carregarDados()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  // ─── Helpers ────────────────────────────────────────────────────────────────

  function nomeMembro(membroId: string): string {
    return membros.find((m) => m.id === membroId)?.nome ?? membroId
  }

  function formatarValor(valor: number): string {
    return valor.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })
  }

  // ─── Handlers de UI ─────────────────────────────────────────────────────────

  function abrirFormCriacao() {
    setFormMode('criar')
    setEditandoId(null)
    setApiError(null)
    reset({ descricao: '', membro_id: '', valor: undefined, dia_recebimento: undefined, data_inicio: '', data_fim: '' })
  }

  function abrirFormEdicao(renda: RendaFixa) {
    setFormMode('editar')
    setEditandoId(renda.id)
    setApiError(null)
    reset({
      descricao: renda.descricao,
      membro_id: renda.membro_id,
      valor: renda.valor,
      dia_recebimento: renda.dia_recebimento,
      data_inicio: renda.data_inicio,
      data_fim: renda.data_fim ?? '',
    })
  }

  function fecharForm() {
    setFormMode('hidden')
    setEditandoId(null)
    setApiError(null)
    reset({ descricao: '', membro_id: '', valor: undefined, dia_recebimento: undefined, data_inicio: '', data_fim: '' })
  }

  // ─── Submit ─────────────────────────────────────────────────────────────────

  async function onSubmit(values: RendaFixaFormValues) {
    if (!accessToken) return
    setApiError(null)

    try {
      if (formMode === 'criar') {
        await criarRendaFixa(accessToken, {
          descricao: values.descricao,
          membro_id: values.membro_id,
          valor: values.valor,
          dia_recebimento: values.dia_recebimento,
          data_inicio: values.data_inicio,
          data_fim: values.data_fim || undefined,
        })
      } else if (formMode === 'editar' && editandoId) {
        const renda = rendas.find((r) => r.id === editandoId)!
        const payload: AtualizarRendaFixaRequest = {
          descricao: values.descricao,
          membro_id: values.membro_id,
          valor: values.valor,
          dia_recebimento: values.dia_recebimento,
          ativa: renda.ativa,
          data_inicio: values.data_inicio,
          data_fim: values.data_fim || undefined,
        }
        await atualizarRendaFixa(accessToken, editandoId, payload)
      }

      fecharForm()
      setLoading(true)
      await carregarDados()
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro desconhecido')
    }
  }

  async function handleInativar(id: string) {
    if (!accessToken) return
    setApiError(null)
    try {
      await inativarRendaFixa(accessToken, id)
      setLoading(true)
      await carregarDados()
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro ao inativar renda fixa')
    }
  }

  // ─── Render ──────────────────────────────────────────────────────────────────

  return (
    <div className="p-6 max-w-2xl mx-auto">
      <div className="mb-4">
        <Link to="/" className="text-blue-600 dark:text-blue-400 hover:underline text-sm">
          ← Voltar para o dashboard
        </Link>
      </div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Rendas Fixas</h1>
        {formMode === 'hidden' && (
          <Button onClick={abrirFormCriacao}>Adicionar renda fixa</Button>
        )}
      </div>

      {/* Erro global de API */}
      {apiError && (
        <p className="mb-4 text-sm text-red-600 dark:text-red-400" role="alert">
          {apiError}
        </p>
      )}

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
                  placeholder="Ex: Salário, Aluguel recebido..."
                  className="mt-1"
                />
                {errors.descricao && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.descricao.message}</p>
                )}
              </div>

              <div className="mb-4">
                <Label htmlFor="membro_id">Membro responsável</Label>
                <select
                  id="membro_id"
                  {...register('membro_id')}
                  className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
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
                  <Label htmlFor="valor">Valor</Label>
                  <Input
                    id="valor"
                    type="number"
                    min={0}
                    step="0.01"
                    {...register('valor')}
                    placeholder="Ex: 5000.00"
                    className="mt-1"
                  />
                  {errors.valor && (
                    <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.valor.message}</p>
                  )}
                </div>
                <div>
                  <Label htmlFor="dia_recebimento">Dia de recebimento</Label>
                  <Input
                    id="dia_recebimento"
                    type="number"
                    min={1}
                    max={31}
                    {...register('dia_recebimento')}
                    placeholder="Ex: 5"
                    className="mt-1"
                  />
                  {errors.dia_recebimento && (
                    <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.dia_recebimento.message}</p>
                  )}
                </div>
              </div>

              <div className="mb-4 grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="data_inicio">Data de início</Label>
                  <Input
                    id="data_inicio"
                    type="date"
                    {...register('data_inicio')}
                    className="mt-1"
                  />
                  {errors.data_inicio && (
                    <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.data_inicio.message}</p>
                  )}
                </div>
                <div>
                  <Label htmlFor="data_fim">Data de fim (opcional)</Label>
                  <Input
                    id="data_fim"
                    type="date"
                    {...register('data_fim')}
                    className="mt-1"
                  />
                </div>
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

      {/* Lista de rendas fixas */}
      {loading ? (
        <p className="text-muted-foreground">Carregando...</p>
      ) : rendas.length === 0 ? (
        <p className="text-muted-foreground">Nenhuma renda fixa cadastrada</p>
      ) : (
        <ul className="space-y-3">
          {rendas.map((renda) => (
            <li key={renda.id}>
              <Card>
                <CardContent className="pt-4 flex items-center justify-between">
                  <div>
                    <p className="font-medium">{renda.descricao}</p>
                    <p className="text-sm text-muted-foreground">{nomeMembro(renda.membro_id)}</p>
                    <p className="text-sm text-muted-foreground">
                      {formatarValor(renda.valor)} · Dia {renda.dia_recebimento}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      Início: {renda.data_inicio}
                      {renda.data_fim ? ` · Fim: ${renda.data_fim}` : ''}
                    </p>
                    <span
                      className={`inline-block mt-1 text-xs font-semibold px-2 py-0.5 rounded-full ${
                        renda.ativa
                          ? 'bg-green-100 dark:bg-green-950/30 text-green-700 dark:text-green-400'
                          : 'bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400'
                      }`}
                    >
                      {renda.ativa ? 'Ativa' : 'Inativa'}
                    </span>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => abrirFormEdicao(renda)}
                    >
                      Editar
                    </Button>
                    {renda.ativa && (
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleInativar(renda.id)}
                      >
                        Inativar
                      </Button>
                    )}
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
