# Avaliação do Codebase e Roadmap Multi-Tenant

> Gerado em: 10/03/2026

---

## Parte 1: Avaliação Geral do Código

### 1. Problemas de Segurança

#### CORS Muito Aberto (Alto Risco)

Em `backend/cmd/server/main.go`:

```go
r.Use(cors.New(cors.Config{
    AllowAllOrigins:  true,  // Aceita requisições de qualquer origem
    AllowCredentials: false,
}))
```

Qualquer site malicioso pode fazer requisições HTTP para a API. Embora `AllowCredentials: false` mitigue cookies, JWTs no header ainda estão expostos.

**Fix**: Whitelist apenas `http://localhost:5173` em dev e os domínios do Vercel em prod.

#### Falta de Rate Limiting

Não há middleware de rate limiting. Um atacante pode fazer força bruta em `POST /api/v1/auth/login` sem limitação de tentativas.

**Fix**: Adicionar rate limiting por IP via `github.com/ulule/limiter` ou similar.

---

### 2. Problema Crítico: Ausência de Isolamento de Dados

**Este é o problema mais grave do sistema atualmente.**

O banco de dados **não filtra dados por usuário**. Nenhuma tabela (`membros`, `categorias`, `cartoes`, `despesas`, etc.) possui uma coluna `usuario_id` ou `familia_id`. Isso significa:

- **Usuário A consegue ver TODOS os dados de TODOS os usuários do banco**
- Toda requisição `GET /api/v1/membros` retorna os membros de todos os usuários sem exceção

Exemplo em `backend/internal/repository/membro_repository.go`:

```go
func (r *MembroRepository) Listar() ([]*domain.Membro, error) {
    err := r.db.Select(&rows, `
        SELECT id, nome, relacionamento, ativo
        FROM membros
        WHERE deleted_at IS NULL  // Sem filtro de usuario_id!
        ORDER BY nome ASC
    `)
}
```

O mesmo padrão se repete em todos os 15+ repositórios. O dashboard também está comprometido: `dashboard_service.go` busca despesas globalmente via `ListarPorFaturaGlobal()` sem filtragem.

**Severidade**: CRÍTICA — todos os usuários do sistema veem os dados uns dos outros.

---

### 3. Seed de Categorias: Compartilhamento Global

Em `backend/internal/repository/seed.go`, as categorias padrão são inseridas **uma única vez no banco inteiro**, sem vínculo com nenhum usuário:

```go
INSERT INTO categorias (nome) VALUES ($1) ON CONFLICT DO NOTHING
```

Todas as categorias são compartilhadas globalmente — não há como um usuário ter sua própria taxonomia.

---

### 4. Regras de Negócio

| Regra | Status | Observação |
|-------|--------|------------|
| RN01 (fatura do cartão) | ✅ Correto | Calcula corretamente `fatura_mes`/`fatura_ano` via `dia_fechamento` |
| RN03 (parcelamento) | ⚠️ Parcial | Estrutura existe, mas sem validação que as parcelas estão nas faturas certas |
| RN06 (saldo mensal) | ✅ Correto | Fórmula implementada corretamente, mas sem filtro de usuário fica sem sentido |
| RN07 (projeção) | ⚠️ Global | `ProjecaoProximosMeses()` retorna projeções globais (todos os usuários) |

**Validação ausente**: `dia_fechamento` pode ser 31 mesmo em meses de 28 dias, sem tratamento de edge case no cálculo de fatura.

---

### 5. Performance

#### Falta de Paginação

`Listar()` em todos os repositórios retorna **todas as linhas** sem paginação. Em uso real com centenas de despesas isso se tornará um gargalo.

#### Índices Ausentes

As migrations criam as tabelas sem índices em colunas de filtro:

- `membros.deleted_at` — sem índice
- `categorias.deleted_at` — sem índice
- `refresh_tokens.usuario_id` — sem índice (busca de tokens por usuário)
- Sem índice composto em `despesas_cartao(fatura_mes, fatura_ano)`

#### Dashboard: Múltiplas Queries Isoladas

`dashboard_service.go` executa 8+ queries isoladas sem joins ou agregações SQL. Em produção com volume real, isso pode ser lento. Uma query consolidada ou uso de `MATERIALIZED VIEW` seria mais eficiente.

---

### 6. Arquitetura e Qualidade

**Pontos positivos:**
- Separação clara Handler → Service → Repository
- JWT com expiração de 15min (adequado)
- Refresh token rotation implementada
- Soft delete com `deleted_at`
- Domain errors bem definidos

**Pontos a melhorar:**
- Sem logs estruturados — apenas `log.Printf` no main; sem rastreabilidade por request ID
- Sem transações em operações que envolvem múltiplos writes (ex: criar parcelas de um cartão)
- Sem validação de domínio robusta (ex: `dia_fechamento` deve ser 1-28)

---

### 7. Offline Sync (Frontend)

Em `frontend/src/lib/syncQueue.ts`, a fila armazena `endpoint`, `method` e `body`, mas:

- Sem rastreamento de `resourceId` ou `version` para detecção de conflito
- Sem implementação de **last-write-wins** real — a última requisição vence silenciosamente
- Se dois dispositivos modificarem o mesmo registro offline, o conflito não é detectado nem reportado ao usuário

---

## Parte 2: Roadmap para Multi-Tenant

Para suportar múltiplas famílias/tenants isolados, estas são as mudanças necessárias:

### 1. Modelo de Dados

Criar tabela `familias`:

```sql
CREATE TABLE familias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nome VARCHAR(100) NOT NULL,
    owner_usuario_id UUID NOT NULL REFERENCES usuarios(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE familia_usuarios (
    familia_id UUID NOT NULL REFERENCES familias(id),
    usuario_id UUID NOT NULL REFERENCES usuarios(id),
    role VARCHAR(20) NOT NULL DEFAULT 'membro', -- admin | membro
    PRIMARY KEY (familia_id, usuario_id)
);
```

Adicionar `familia_id UUID NOT NULL REFERENCES familias(id)` a todas as tabelas:
- `membros`, `categorias`, `cartoes_credito`, `despesas_cartao`
- `assinaturas`, `contas_fixas`, `despesas_gerais`
- `rendas_fixas`, `rendas_variaveis`, `rendas_extras`, `rendimentos_investimento`

**~15 migrations necessárias** | Impacto: **ALTO**

---

### 2. Autenticação e JWT

**Atual**: JWT contém apenas `sub` (usuarioID).

**Necessário**: JWT deve incluir `familia_id` e `role`:

```json
{
  "sub": "usuario-uuid",
  "familia_id": "familia-uuid",
  "role": "admin",
  "exp": 1234567890
}
```

Novos endpoints:

| Endpoint | Descrição |
|----------|-----------|
| `POST /api/v1/familias` | Criar nova família (owner) |
| `POST /api/v1/familias/:id/convidar` | Convidar usuário para a família |
| `GET /api/v1/familias` | Listar famílias do usuário autenticado |

O `AuthMiddleware` deve validar o `familia_id` do JWT e disponibilizá-lo no contexto Gin.

**Arquivos impactados**: `auth_service.go`, `middleware/auth.go`, todos os handlers | Impacto: **ALTO**

---

### 3. Isolamento de Queries

Todos os repositórios precisam receber `familiaID` como parâmetro:

```go
// Antes
func (r *MembroRepository) Listar() ([]*domain.Membro, error) {
    SELECT * FROM membros WHERE deleted_at IS NULL
}

// Depois
func (r *MembroRepository) Listar(familiaID string) ([]*domain.Membro, error) {
    SELECT * FROM membros WHERE familia_id = $1 AND deleted_at IS NULL
}
```

**15+ repositórios a refatorar** + todos os handlers extraindo `familiaID` do JWT no contexto Gin.

Impacto: **ALTO** — risco de vazamento se algum repositório for esquecido

---

### 4. Dashboard e Projeção

`dashboard_service.go` deve receber `familiaID` e passá-lo a todas as chamadas de repositório:

```go
func (s *DashboardService) ResumoMensal(familiaID string, mes, ano int) (*ResumoMensal, error) {
    rendas, _ := s.rendaRepo.ListarVigentesPorMes(familiaID, mes, ano)
    despesas, _ := s.despesaRepo.ListarPorFatura(familiaID, mes, ano)
    ...
}
```

Impacto: **MÉDIO** — 1 dia de trabalho

---

### 5. Seed de Dados

Modificar o seed de categorias para criar categorias por família:

```go
func SeedCategorias(db *sqlx.DB, familiaID string) error {
    // Inserir categorias padrão vinculadas à familia
    INSERT INTO categorias (nome, familia_id) VALUES ($1, $2)
}
```

Ao criar uma nova família, disparar o seed automaticamente.

Impacto: **BAIXO**

---

### 6. Frontend: Estado Global de Tenant

O `authStore` precisa armazenar e expor `familiaId`:

```typescript
interface AuthContextValue {
  accessToken: string | null
  familiaId: string | null  // NOVO
  role: 'admin' | 'membro' | null  // NOVO
  isAuthenticated: boolean
}

// Ao fazer login, decodificar JWT para extrair familia_id
const decoded = jwtDecode(response.access_token)
setFamiliaId(decoded.familia_id)
```

Para usuários que pertencem a múltiplas famílias (ex: administrador de mais de uma família):
- Novo `FamiliaContext` com lista de famílias disponíveis
- Dropdown para trocar de família
- Família selecionada persistida em `localStorage`

Impacto: **MÉDIO** — 2 dias de trabalho

---

### 7. Offline Sync com Família

A fila de sync precisa rastrear a família para evitar sincronizar dados da família errada se o usuário trocar:

```typescript
interface SyncQueueItem {
  endpoint: string
  method: string
  body: string
  familiaId: string   // NOVO
  resourceType: string
  resourceId: string
  version: number     // Para last-write-wins real
}

// Ao processar sync, ignorar itens de outras familias
if (item.familiaId !== currentFamiliaId) continue
```

Impacto: **BAIXO**

---

### 8. Row-Level Security no PostgreSQL (Defesa em Profundidade)

Camada extra de proteção além da lógica da aplicação:

```sql
-- Habilitar RLS em todas as tabelas
ALTER TABLE membros ENABLE ROW LEVEL SECURITY;

-- Policy: acesso apenas aos dados da família do contexto atual
CREATE POLICY "isolamento_por_familia" ON membros
  FOR ALL USING (familia_id = current_setting('app.current_familia_id')::uuid);
```

No backend, antes de cada query:

```go
db.Exec("SET LOCAL app.current_familia_id = $1", familiaID)
```

Isso garante que mesmo um bug na lógica da aplicação não vaze dados de outras famílias.

Impacto: **BAIXO** (defesa extra) — 1-2 dias de implementação

---

### Resumo e Estimativas

| Componente | Estimativa | Impacto |
|------------|-----------|---------|
| Modelo de dados (migrations + familia_id) | 2-3 dias | **ALTO** |
| Autenticação + JWT com familia_id | 3-4 dias | **ALTO** |
| Isolamento de queries (15+ repos) | 4-5 dias | **ALTO** |
| Dashboard e projeção | 1 dia | **MÉDIO** |
| Seed por família | 1 dia | **BAIXO** |
| Frontend (estado + UI de troca de família) | 2 dias | **MÉDIO** |
| Offline sync com familia_id | 1 dia | **BAIXO** |
| PostgreSQL RLS | 1-2 dias | **BAIXO** |
| Testes de regressão e isolamento | 2-3 dias | **ALTO** |
| **TOTAL** | **17-25 dias** | |

---

### Recomendações por Prioridade

| Prioridade | Ação |
|-----------|------|
| 🔴 Imediato | Adicionar `usuario_id`/`familia_id` em todas as tabelas e filtros — risco de vazamento de dados é crítico |
| 🔴 Antes de prod multi-usuário | Implementar isolamento completo de queries em todos os repositórios |
| 🟡 Curto prazo | Restringir CORS a domínios específicos, adicionar rate limiting em `/auth/login` |
| 🟡 Curto prazo | Criar índices em `deleted_at`, colunas de filtro e `familia_id` |
| 🟢 Médio prazo | Implementar PostgreSQL RLS como camada defensiva |
| 🟢 Médio prazo | Adicionar paginação nas listagens |
| 🟢 Médio prazo | Implementar versioning na fila offline para last-write-wins real |

> O código está bem estruturado arquiteturalmente, mas **inseguro para multi-usuário no estado atual**. Uma mudança fundamental é necessária antes de qualquer deploy com múltiplos usuários reais não confiáveis.
