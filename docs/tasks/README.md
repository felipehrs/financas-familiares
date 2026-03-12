# Tasks — Tarefas Executáveis por Sprint

Esta pasta contém tarefas detalhadas e executáveis, organizadas por sprint e item técnico (TT) ou user story (US).

Cada tarefa segue a abordagem **TDD (Test-Driven Development)** conforme especificado em [docs/tecnico/01-tech-spec.md](../tecnico/01-tech-spec.md).

---

## Organização

```
docs/tasks/
├── README.md           (este arquivo)
├── tt-08/              (Sprint 9 - TT-08: CORS e Rate Limiting)
│   ├── 01-cors-restrito.md
│   ├── 02-rate-limiting.md
│   └── 03-testes-e-docs.md
└── tt-09/              (Sprint 9 - TT-09: Índices — a criar)
```

---

## Tarefas Disponíveis

### Sprint 9 — Segurança e Isolamento de Dados

#### TT-08: CORS e Rate Limiting

| # | Tarefa | Estimativa | Abordagem |
|---|--------|------------|-----------|
| 1 | [CORS Restrito](./tt-08/01-cors-restrito.md) | 45 min | TDD (Red → Green → Refactor) |
| 2 | [Rate Limiting](./tt-08/02-rate-limiting.md) | 60 min | TDD (Red → Green → Refactor) |
| 3 | [Validação e Docs](./tt-08/03-testes-e-docs.md) | 30 min | Validação E2E + Documentação |

**Total estimado:** ~2h15min

**Plano técnico:** [docs/tecnico/05-plano-tt08-cors-rate-limiting.md](../tecnico/05-plano-tt08-cors-rate-limiting.md)

---

## Abordagem TDD

Todas as tarefas de implementação seguem o ciclo TDD:

### Red Phase (Teste Falhando)
1. Escrever o teste ANTES do código
2. Executar e verificar que FALHA
3. Commit: "test: add failing test for X"

### Green Phase (Teste Passando)
1. Implementar o código mínimo para passar
2. Executar e verificar que PASSA
3. Commit: "feat: implement X"

### Refactor Phase (Melhoria)
1. Melhorar o código mantendo testes verdes
2. Executar testes para garantir que continuam passando
3. Commit: "refactor: improve X"

---

## Como Usar

### 1. Escolher uma Tarefa
Navegue para a pasta (ex: `tt-08/`) e abra a tarefa.

### 2. Seguir o Checklist
Marque os checkboxes `[ ]` como `[x]` conforme avança.

### 3. Validar Critério de Conclusão
Verifique a seção final de cada tarefa.

### 4. Fazer Commit
Use a mensagem sugerida ou adapte.

### 5. Prosseguir
Cada tarefa indica o próximo passo.

---

## Convenções

Cada tarefa documenta:

- **Metadados:** Sprint, item, fase, estimativa, dependências, abordagem
- **Objetivo:** O que será feito e por quê
- **Checklist TDD:** Red → Green → Refactor
- **Critério de Conclusão:** Validação de testes + build
- **Commit Sugerido:** Mensagem pré-formatada
- **Próximos Passos:** Link para próxima tarefa
- **Referências:** Links para docs técnicos

---

## Referências

- [Sprints](../produto/sprints.md) — Visão geral de todas as sprints
- [Stories](../produto/stories.md) — User stories e critérios
- [Tech Spec](../tecnico/01-tech-spec.md) — Estratégia de TDD

---

**Última atualização:** 12/03/2026 (TT-08 criado com abordagem TDD)
