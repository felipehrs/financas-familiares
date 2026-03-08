import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { listarContasFixas, criarContaFixa, atualizarContaFixa, alterarAtivoContaFixa } from '@/offline/contas_fixas'
import { listarMembros } from '@/api/membros'
import { listarCategorias } from '@/api/categorias'
import type { ContaFixa, AtualizarContaFixaRequest } from '@/types/conta_fixa'
import type { Membro } from '@/types/membro'
import type { Categoria } from '@/types/categoria'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

// ─── Constantes ───────────────────────────────────────────────────────────────

const FORMAS_PAGAMENTO = ['Cartão de crédito', 'Cartão de débito', 'PIX', 'Boleto', 'Débito automático']

// ─── Schema de validação ──────────────────────────────────────────────────────

const contaFixaSchema = z.object({
  descricao: z.string().min(2, 'Descrição é obrigatória e deve ter pelo menos 2 caracteres'),
  membro_id: z.string().min(1, 'Membro é obrigatório'),
  categoria_id: z.string().optional(),
  valor: z.coerce.number({ invalid_type_error: 'Informe o valor' }).positive('Valor deve ser maior que zero'),
  dia_vencimento: z.coerce.number({ invalid_type_error: 'Informe o dia' }).int().min(1, 'Mínimo 1').max(28, 'Máximo 28'),
  forma_pagamento: z.string().min(1, 'Forma de pagamento é obrigatória'),
})

type ContaFixaFormValues = z.infer<typeof contaFixaSchema>

// ─── Componente ───────────────────────────────────────────────────────────────

type FormMode = 'hidden' | 'criar' | 'editar'

export function ContasFixasPage() {
  const { accessToken } = useAuth()

  const [contasFixas, setContasFixas] = useState<ContaFixa[]>([])
  const [membros, setMembros] = useState<Membro[]>([])
  const [categorias, setCategorias] = useState<Categoria[]>([])
  const [loading, setLoading] = useState(true)
  const [apiError, setApiError] = useState<string | null>(null)
  const [formMode, setFormMode] = useState<FormMode>('hidden')
  const [editandoId, setEditandoId] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<ContaFixaFormValues>({
    resolver: zodResolver(contaFixaSchema),
    defaultValues: {
      descricao: '',
      membro_id: '',
      categoria_id: '',
      valor: '' as unknown as number,
      dia_vencimento: '' as unknown as number,
      forma_pagamento: '',
    },
  })

  // ─── Carrega dados ──────────────────────────────────────────────────────────

  async function carregarDados() {
    if (!accessToken) return
    try {
      const [listaContasFixas, listaMembros, listaCategorias] = await Promise.all([
        listarContasFixas(accessToken),
        listarMembros(accessToken),
        listarCategorias(accessToken),
      ])
      setContasFixas(listaContasFixas)
      setMembros(listaMembros)
      setCategorias(listaCategorias)
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

  // ─── Handlers de UI ─────────────────────────────────────────────────────────

  function abrirFormCriacao() {
    setFormMode('criar')
    setEditandoId(null)
    setApiError(null)
    reset({
      descricao: '',
      membro_id: '',
      categoria_id: '',
      valor: undefined,
      dia_vencimento: undefined,
      forma_pagamento: '',
    })
  }

  function abrirFormEdicao(conta: ContaFixa) {
    setFormMode('editar')
    setEditandoId(conta.id)
    setApiError(null)
    reset({
      descricao: conta.descricao,
      membro_id: conta.membro_id,
      categoria_id: conta.categoria_id ?? '',
      valor: conta.valor,
      dia_vencimento: conta.dia_vencimento,
      forma_pagamento: conta.forma_pagamento,
    })
  }

  function fecharForm() {
    setFormMode('hidden')
    setEditandoId(null)
    setApiError(null)
    reset({
      descricao: '',
      membro_id: '',
      categoria_id: '',
      valor: undefined,
      dia_vencimento: undefined,
      forma_pagamento: '',
    })
  }

  // ─── Submit ─────────────────────────────────────────────────────────────────

  async function onSubmit(values: ContaFixaFormValues) {
    if (!accessToken) return
    setApiError(null)

    try {
      if (formMode === 'criar') {
        await criarContaFixa(accessToken, {
          descricao: values.descricao,
          membro_id: values.membro_id,
          categoria_id: values.categoria_id || undefined,
          valor: values.valor,
          dia_vencimento: values.dia_vencimento,
          forma_pagamento: values.forma_pagamento,
        })
      } else if (formMode === 'editar' && editandoId) {
        const conta = contasFixas.find((c) => c.id === editandoId)!
        const payload: AtualizarContaFixaRequest = {
          descricao: values.descricao,
          membro_id: values.membro_id,
          categoria_id: values.categoria_id || undefined,
          valor: values.valor,
          dia_vencimento: values.dia_vencimento,
          forma_pagamento: values.forma_pagamento,
          ativa: conta.ativa,
        }
        await atualizarContaFixa(accessToken, editandoId, payload)
      }

      fecharForm()
      setLoading(true)
      await carregarDados()
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro desconhecido')
    }
  }

  async function handleAlterarAtivo(conta: ContaFixa) {
    if (!accessToken) return
    setApiError(null)
    try {
      await alterarAtivoContaFixa(accessToken, conta.id, !conta.ativa)
      setLoading(true)
      await carregarDados()
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro ao alterar status')
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
        <h1 className="text-2xl font-bold">Contas Fixas</h1>
        {formMode === 'hidden' && (
          <Button onClick={abrirFormCriacao}>Nova conta fixa</Button>
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
                  placeholder="Ex: Condomínio, Aluguel..."
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

              <div className="mb-4">
                <Label htmlFor="categoria_id">Categoria (opcional)</Label>
                <select
                  id="categoria_id"
                  {...register('categoria_id')}
                  className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                >
                  <option value="">Sem categoria</option>
                  {categorias.map((cat) => (
                    <option key={cat.id} value={cat.id}>
                      {cat.nome}
                    </option>
                  ))}
                </select>
              </div>

              <div className="mb-4 grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="valor">Valor (R$)</Label>
                  <Input
                    id="valor"
                    type="number"
                    min={0}
                    step="0.01"
                    {...register('valor')}
                    placeholder="Ex: 800.00"
                    className="mt-1"
                  />
                  {errors.valor && (
                    <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.valor.message}</p>
                  )}
                </div>
                <div>
                  <Label htmlFor="dia_vencimento">Dia de vencimento</Label>
                  <Input
                    id="dia_vencimento"
                    type="number"
                    min={1}
                    max={28}
                    {...register('dia_vencimento')}
                    placeholder="Ex: 10"
                    className="mt-1"
                  />
                  {errors.dia_vencimento && (
                    <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.dia_vencimento.message}</p>
                  )}
                </div>
              </div>

              <div className="mb-4">
                <Label htmlFor="forma_pagamento">Forma de pagamento</Label>
                <select
                  id="forma_pagamento"
                  {...register('forma_pagamento')}
                  className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                >
                  <option value="">Selecione uma forma de pagamento</option>
                  {FORMAS_PAGAMENTO.map((forma) => (
                    <option key={forma} value={forma}>
                      {forma}
                    </option>
                  ))}
                </select>
                {errors.forma_pagamento && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.forma_pagamento.message}</p>
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

      {/* Lista de contas fixas */}
      {loading ? (
        <p className="text-muted-foreground">Carregando...</p>
      ) : contasFixas.length === 0 ? (
        <p className="text-muted-foreground">Nenhuma conta fixa cadastrada</p>
      ) : (
        <ul className="space-y-3">
          {contasFixas.map((conta) => (
            <li key={conta.id}>
              <Card>
                <CardContent className="pt-4 flex items-center justify-between">
                  <div>
                    <p className="font-medium">{conta.descricao}</p>
                    <p className="text-sm text-muted-foreground">{nomeMembro(conta.membro_id)}</p>
                    <p className="text-sm text-muted-foreground">
                      {conta.valor.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })}
                      {' · '}Dia {conta.dia_vencimento}
                      {' · '}{conta.forma_pagamento}
                    </p>
                    <span
                      className={`inline-block mt-1 text-xs font-semibold px-2 py-0.5 rounded-full ${conta.ativa ? 'bg-green-100 dark:bg-green-950/30 text-green-700 dark:text-green-400' : 'bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400'}`}
                    >
                      {conta.ativa ? 'Ativa' : 'Inativa'}
                    </span>
                  </div>
                  <div className="flex gap-2 flex-wrap justify-end">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => abrirFormEdicao(conta)}
                    >
                      Editar
                    </Button>
                    {conta.ativa ? (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleAlterarAtivo(conta)}
                      >
                        Inativar
                      </Button>
                    ) : (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleAlterarAtivo(conta)}
                      >
                        Reativar
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
