import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { listarAssinaturas, criarAssinatura, atualizarAssinatura, alterarStatusAssinatura } from '@/offline/assinaturas'
import { listarMembros } from '@/api/membros'
import { listarCategorias } from '@/api/categorias'
import type { Assinatura, AtualizarAssinaturaRequest, StatusAssinatura } from '@/types/assinatura'
import type { Membro } from '@/types/membro'
import type { Categoria } from '@/types/categoria'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

// ─── Constantes ───────────────────────────────────────────────────────────────

const FORMAS_PAGAMENTO = ['Cartão de crédito', 'Cartão de débito', 'PIX', 'Boleto', 'Débito automático']

// ─── Schema de validação ──────────────────────────────────────────────────────

const assinaturaSchema = z.object({
  nome: z.string().min(2, 'Nome é obrigatório'),
  membro_id: z.string().min(1, 'Membro é obrigatório'),
  categoria_id: z.string().optional(),
  valor: z.coerce.number().positive('Valor deve ser maior que zero'),
  dia_cobranca: z.coerce.number().int().min(1).max(31, 'Dia deve ser entre 1 e 31'),
  forma_pagamento: z.string().min(1, 'Forma de pagamento é obrigatória'),
  status: z.enum(['ativa', 'pausada', 'cancelada']).optional(),
})

type AssinaturaFormValues = z.infer<typeof assinaturaSchema>

// ─── Helpers de badge de status ───────────────────────────────────────────────

function badgeClasses(status: StatusAssinatura): string {
  switch (status) {
    case 'ativa':
      return 'bg-green-100 dark:bg-green-950/30 text-green-700 dark:text-green-400'
    case 'pausada':
      return 'bg-amber-100 dark:bg-amber-950/30 text-amber-700 dark:text-amber-400'
    case 'cancelada':
      return 'bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400'
  }
}

function labelStatus(status: StatusAssinatura): string {
  switch (status) {
    case 'ativa':
      return 'Ativa'
    case 'pausada':
      return 'Pausada'
    case 'cancelada':
      return 'Cancelada'
  }
}

// ─── Componente ───────────────────────────────────────────────────────────────

type FormMode = 'hidden' | 'criar' | 'editar'

export function AssinaturasPage() {
  const { accessToken } = useAuth()

  const [assinaturas, setAssinaturas] = useState<Assinatura[]>([])
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
  } = useForm<AssinaturaFormValues>({
    resolver: zodResolver(assinaturaSchema),
    defaultValues: {
      nome: '',
      membro_id: '',
      categoria_id: '',
      valor: '' as unknown as number,
      dia_cobranca: '' as unknown as number,
      forma_pagamento: '',
    },
  })

  // ─── Carrega dados ──────────────────────────────────────────────────────────

  async function carregarDados() {
    if (!accessToken) return
    try {
      const [listaAssinaturas, listaMembros, listaCategorias] = await Promise.all([
        listarAssinaturas(accessToken),
        listarMembros(accessToken),
        listarCategorias(accessToken),
      ])
      setAssinaturas(listaAssinaturas)
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
      nome: '',
      membro_id: '',
      categoria_id: '',
      valor: undefined,
      dia_cobranca: undefined,
      forma_pagamento: '',
    })
  }

  function abrirFormEdicao(assinatura: Assinatura) {
    setFormMode('editar')
    setEditandoId(assinatura.id)
    setApiError(null)
    reset({
      nome: assinatura.nome,
      membro_id: assinatura.membro_id,
      categoria_id: assinatura.categoria_id ?? '',
      valor: assinatura.valor,
      dia_cobranca: assinatura.dia_cobranca,
      forma_pagamento: assinatura.forma_pagamento,
      status: assinatura.status,
    })
  }

  function fecharForm() {
    setFormMode('hidden')
    setEditandoId(null)
    setApiError(null)
    reset({
      nome: '',
      membro_id: '',
      categoria_id: '',
      valor: undefined,
      dia_cobranca: undefined,
      forma_pagamento: '',
    })
  }

  // ─── Submit ─────────────────────────────────────────────────────────────────

  async function onSubmit(values: AssinaturaFormValues) {
    if (!accessToken) return
    setApiError(null)

    try {
      if (formMode === 'criar') {
        await criarAssinatura(accessToken, {
          nome: values.nome,
          membro_id: values.membro_id,
          categoria_id: values.categoria_id || undefined,
          valor: values.valor,
          dia_cobranca: values.dia_cobranca,
          forma_pagamento: values.forma_pagamento,
        })
      } else if (formMode === 'editar' && editandoId) {
        const assinatura = assinaturas.find((a) => a.id === editandoId)!
        const payload: AtualizarAssinaturaRequest = {
          nome: values.nome,
          membro_id: values.membro_id,
          categoria_id: values.categoria_id || undefined,
          valor: values.valor,
          dia_cobranca: values.dia_cobranca,
          forma_pagamento: values.forma_pagamento,
          status: values.status ?? assinatura.status,
        }
        await atualizarAssinatura(accessToken, editandoId, payload)
      }

      fecharForm()
      setLoading(true)
      await carregarDados()
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro desconhecido')
    }
  }

  async function handleAlterarStatus(id: string, novoStatus: StatusAssinatura) {
    if (!accessToken) return
    setApiError(null)
    try {
      await alterarStatusAssinatura(accessToken, id, novoStatus)
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
        <h1 className="text-2xl font-bold">Assinaturas</h1>
        {formMode === 'hidden' && (
          <Button onClick={abrirFormCriacao}>Nova assinatura</Button>
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
                <Label htmlFor="nome">Nome da assinatura</Label>
                <Input
                  id="nome"
                  {...register('nome')}
                  placeholder="Ex: Netflix, Spotify..."
                  className="mt-1"
                />
                {errors.nome && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.nome.message}</p>
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
                    placeholder="Ex: 49.90"
                    className="mt-1"
                  />
                  {errors.valor && (
                    <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.valor.message}</p>
                  )}
                </div>
                <div>
                  <Label htmlFor="dia_cobranca">Dia de cobrança</Label>
                  <Input
                    id="dia_cobranca"
                    type="number"
                    min={1}
                    max={31}
                    {...register('dia_cobranca')}
                    placeholder="Ex: 15"
                    className="mt-1"
                  />
                  {errors.dia_cobranca && (
                    <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.dia_cobranca.message}</p>
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

              {formMode === 'editar' && (
                <div className="mb-4">
                  <Label htmlFor="status">Status</Label>
                  <select
                    id="status"
                    {...register('status')}
                    className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                  >
                    <option value="ativa">Ativa</option>
                    <option value="pausada">Pausada</option>
                    <option value="cancelada">Cancelada</option>
                  </select>
                </div>
              )}

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

      {/* Lista de assinaturas */}
      {loading ? (
        <p className="text-muted-foreground">Carregando...</p>
      ) : assinaturas.length === 0 ? (
        <p className="text-muted-foreground">Nenhuma assinatura cadastrada</p>
      ) : (
        <ul className="space-y-3">
          {assinaturas.map((assinatura) => (
            <li key={assinatura.id}>
              <Card>
                <CardContent className="pt-4 flex items-center justify-between">
                  <div>
                    <p className="font-medium">{assinatura.nome}</p>
                    <p className="text-sm text-muted-foreground">{nomeMembro(assinatura.membro_id)}</p>
                    <p className="text-sm text-muted-foreground">
                      {assinatura.valor.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })}
                      {' · '}Dia {assinatura.dia_cobranca}
                      {' · '}{assinatura.forma_pagamento}
                    </p>
                    <span
                      className={`inline-block mt-1 text-xs font-semibold px-2 py-0.5 rounded-full ${badgeClasses(assinatura.status)}`}
                    >
                      {labelStatus(assinatura.status)}
                    </span>
                  </div>
                  <div className="flex gap-2 flex-wrap justify-end">
                    {assinatura.status !== 'cancelada' && (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => abrirFormEdicao(assinatura)}
                      >
                        Editar
                      </Button>
                    )}
                    {assinatura.status === 'ativa' && (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleAlterarStatus(assinatura.id, 'pausada')}
                      >
                        Pausar
                      </Button>
                    )}
                    {assinatura.status === 'pausada' && (
                      <>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleAlterarStatus(assinatura.id, 'ativa')}
                        >
                          Reativar
                        </Button>
                        <Button
                          variant="destructive"
                          size="sm"
                          onClick={() => handleAlterarStatus(assinatura.id, 'cancelada')}
                        >
                          Cancelar
                        </Button>
                      </>
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
