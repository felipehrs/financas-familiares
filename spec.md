# Especificação Funcional - Sistema de Gestão de Finanças Familiares

**Versão:** 1.1
**Data:** 06 de Março de 2026
**Status:** Especificação Inicial

---

## 1. Visão Geral

### 1.1 Objetivo
Desenvolver um sistema para gestão completa das finanças familiares, permitindo controle detalhado de rendas e gastos de todos os membros da família, incluindo despesas em cartões de crédito (com parcelamento), assinaturas recorrentes, contas fixas, despesas gerais e rendimentos variados.

### 1.2 Problema a Resolver
Necessidade de ter visibilidade clara e organizada das finanças da família, com recursos específicos para:
- Consolidar rendas de múltiplos membros da família
- Controlar múltiplos cartões de crédito com diferentes datas de fechamento
- Gerenciar compras parceladas
- Acompanhar assinaturas e serviços recorrentes
- Controlar contas fixas mensais
- Projetar despesas futuras considerando parcelas a vencer
- Visualizar o saldo familiar (rendas menos despesas) por mês

### 1.3 Público-Alvo
Famílias que desejam ter controle detalhado e consolidado de suas finanças, com visibilidade individual por membro e visão geral da família.

---

## 2. Requisitos Funcionais

### 2.1 Gestão de Membros da Família

#### RF01 - Cadastro de Membros
- O usuário deve poder cadastrar os membros da família (ex: Pai, Mãe, Filho(a))
- Informações por membro:
  - Nome (obrigatório)
  - Relacionamento (ex: cônjuge, filho, etc.) — opcional
  - Status (ativo/inativo)
- Cartões, despesas e rendas podem ser vinculados a um membro específico
- O sistema deve oferecer visão consolidada da família e visão individual por membro

### 2.2 Gestão de Categorias

#### RF02 - Cadastro de Categorias
- O usuário deve poder criar categorias personalizadas de despesas e rendas
- Sugestões iniciais: Alimentação, Transporte, Lazer, Saúde, Educação, Moradia, Vestuário, Outros
- Deve permitir adicionar, editar e excluir categorias
- Cada despesa e renda deve estar vinculada a uma categoria

### 2.3 Gestão de Cartões de Crédito

#### RF03 - Cadastro de Cartões
- O sistema deve permitir cadastro de múltiplos cartões de crédito
- Informações obrigatórias por cartão:
  - Nome/Identificação do cartão
  - Membro responsável (FK para Membro)
  - Dia de fechamento da fatura
  - Dia de vencimento da fatura
  - Limite do cartão (opcional)

#### RF04 - Lançamento de Despesas no Cartão
- Campos obrigatórios:
  - Data da compra
  - Cartão utilizado
  - Descrição da despesa
  - Categoria
  - Valor total
  - Número de parcelas (padrão: 1 para compra à vista)
  - Parcela atual (para controle de parcelamento)

#### RF05 - Cálculo Automático de Parcelas
- O sistema deve calcular automaticamente o valor de cada parcela (valor total / número de parcelas)
- Deve identificar automaticamente a qual fatura cada parcela pertence

#### RF06 - Identificação de Fatura
- Baseado na data da compra e no dia de fechamento do cartão:
  - Se a compra foi feita ANTES do dia de fechamento → vai para a fatura do mês corrente
  - Se a compra foi feita APÓS o dia de fechamento → vai para a fatura do mês seguinte
- A fatura deve ser identificada como "MÊS/ANO" (ex: "JAN/26", "FEV/26")
- O vencimento da fatura deve ser calculado automaticamente baseado no dia configurado

#### RF07 - Controle de Limite do Cartão
- Exibir saldo disponível considerando compras lançadas
- Alertar quando o limite estiver próximo de ser atingido (ex: >80%)

### 2.4 Gestão de Assinaturas

#### RF08 - Cadastro de Assinaturas
- Campos obrigatórios:
  - Nome do serviço
  - Membro responsável (FK para Membro)
  - Categoria
  - Valor mensal
  - Dia da cobrança
  - Forma de pagamento (cartão de crédito, débito, etc)
  - Status (ativa/cancelada)
  - Observações (opcional)

#### RF09 - Recorrência Automática
- Assinaturas ativas devem ser automaticamente consideradas nas projeções mensais
- Deve permitir pausar ou cancelar uma assinatura sem excluí-la do histórico

### 2.5 Gestão de Contas Fixas

#### RF10 - Cadastro de Contas Fixas
- Campos obrigatórios:
  - Descrição (ex: Escola, Financiamento, Condomínio)
  - Membro responsável (FK para Membro)
  - Categoria
  - Valor
  - Dia do vencimento
  - Forma de pagamento
  - Status (ativa/inativa)
  - Observações (opcional)

#### RF11 - Recorrência de Contas Fixas
- Contas fixas ativas devem ser consideradas automaticamente nas projeções mensais
- Deve permitir alterar o valor quando houver reajuste

### 2.6 Gestão de Despesas Gerais

#### RF12 - Lançamento de Despesas à Vista
- Campos obrigatórios:
  - Data
  - Membro responsável (FK para Membro)
  - Descrição
  - Categoria
  - Valor
  - Forma de pagamento (dinheiro, débito, PIX, etc)
  - Observações (opcional)

#### RF13 - Histórico de Despesas
- Visualizar todas as despesas gerais lançadas
- Filtrar por período, categoria, forma de pagamento, membro
- Exportar relatórios

### 2.7 Gestão de Rendas

#### RF14 - Cadastro de Renda Fixa (Salário)
- Para rendas recorrentes com valor fixo mensal
- Campos obrigatórios:
  - Descrição (ex: "Salário - Empresa X")
  - Membro (FK para Membro)
  - Valor mensal
  - Dia de recebimento
  - Status (ativa/inativa)
  - Observações (opcional)
- Rendas fixas ativas são consideradas automaticamente nas projeções mensais

#### RF15 - Registro de Renda Variável
- Para rendas recorrentes cujo valor muda a cada mês (freelance, comissões, bônus)
- Campos obrigatórios:
  - Descrição (ex: "Freelance Design")
  - Membro (FK para Membro)
  - Mês/Ano de referência
  - Valor recebido no mês
  - Data de recebimento
  - Observações (opcional)
- Cada lançamento é um registro individual por mês

#### RF16 - Registro de Renda Extra / Pontual
- Para recebimentos não recorrentes (venda de item, restituição, décimo terceiro, etc.)
- Campos obrigatórios:
  - Descrição
  - Membro (FK para Membro)
  - Data de recebimento
  - Valor
  - Observações (opcional)

#### RF17 - Registro de Rendimentos de Investimentos
- Para dividendos, rendimentos de CDB, poupança, FIIs, etc.
- Campos obrigatórios:
  - Descrição (ex: "Dividendos ITSA4", "Rendimento CDB")
  - Membro (FK para Membro)
  - Data de recebimento
  - Valor
  - Valor distribuído (opcional, padrão: R$ 0,00)
  - Observações (opcional)
- O campo **Valor Distribuído** representa a parcela do rendimento que o usuário decide injetar no orçamento familiar do mês. O restante é considerado reinvestido e não impacta o saldo do mês.
- O valor distribuído não pode ser maior que o valor do rendimento.

#### RF17a - Distribuição de Rendimentos no Dashboard
- Rendimentos de investimentos são exibidos em seção própria no dashboard, separados das demais rendas.
- Apenas o **valor distribuído** de cada rendimento é somado ao total de rendas operacionais para cálculo do saldo.
- O dashboard exibe:
  - Total de rendimentos recebidos no mês
  - Total distribuído (valor que entra no orçamento)
  - Total reinvestido (diferença: recebido − distribuído)
- O usuário pode editar o valor distribuído de qualquer rendimento a qualquer momento dentro do mês.

#### RF18 - Histórico e Resumo de Rendas
- Visualizar todas as rendas por tipo, membro e período
- Filtrar por tipo de renda, membro e mês/ano
- Exibir total de rendas consolidado por membro e da família

---

## 3. Dashboard e Visualizações

### 3.1 Dashboard Principal

#### RF19 - Resumo do Mês Atual
Exibir consolidado do mês corrente:

**Rendas Operacionais:**
- Total de rendas fixas
- Total de rendas variáveis e extras lançadas no mês
- **TOTAL DE RENDAS OPERACIONAIS**

**Rendimentos de Investimentos (seção separada):**
- Total recebido no mês (informativo)
- Total distribuído para o orçamento (soma dos valores_distribuídos)
- Total reinvestido (recebido − distribuído) — informativo

**Despesas:**
- Total de despesas gerais
- Total de assinaturas ativas
- Total de contas fixas
- Total da fatura atual dos cartões
- **TOTAL DE DESPESAS**

**Saldo:**
- **SALDO DO MÊS = (TOTAL RENDAS OPERACIONAIS + TOTAL DISTRIBUÍDO) − TOTAL DESPESAS**
- Indicação visual positivo (verde) / negativo (vermelho)
- Visão por membro e visão consolidada da família

#### RF20 - Despesas por Categoria
- Gráfico ou tabela mostrando distribuição de gastos por categoria
- Exibir valor e percentual de cada categoria sobre o total de despesas
- Considerar todas as fontes: despesas gerais, cartões, assinaturas e contas fixas
- Filtrar por membro ou exibir consolidado familiar

#### RF21 - Evolução Mensal
- Gráfico de linha mostrando evolução mês a mês de: rendas, despesas e saldo
- Comparativo entre meses (ex: saldo do mês atual vs mês anterior)
- Identificar tendências de aumento ou redução

#### RF22 - Projeção de Despesas Futuras
- Projetar despesas dos próximos 3 meses considerando:
  - Parcelas de cartão a vencer
  - Assinaturas ativas
  - Contas fixas
- Projetar rendas dos próximos 3 meses considerando rendas fixas ativas
- Exibir por mês:
  - Total estimado de cartões (parcelas)
  - Total de contas fixas + assinaturas
  - **Total geral estimado de despesas**
  - **Total estimado de rendas**
  - **Saldo estimado**

---

## 4. Regras de Negócio

### RN01 - Cálculo de Fatura de Cartão
```
Se DATA_COMPRA <= DIA_FECHAMENTO do mês corrente:
    FATURA = MÊS_CORRENTE
Senão:
    FATURA = MÊS_SEGUINTE

VENCIMENTO_FATURA = DIA_VENCIMENTO do mês da fatura
```

**Exemplo:**
- Cartão com fechamento dia 12 e vencimento dia 18
- Compra em 10/03/2026 → Fatura MAR/26, vence 18/03/2026
- Compra em 15/03/2026 → Fatura ABR/26, vence 18/04/2026

### RN02 - Cálculo de Parcelas
```
VALOR_PARCELA = VALOR_TOTAL / NUMERO_PARCELAS
```

### RN03 - Distribuição de Parcelas nas Faturas
Para compra parcelada:
- Parcela 1 vai para a fatura calculada pela data da compra
- Parcela 2 vai para a fatura do mês seguinte
- Parcela 3 vai para a fatura do mês seguinte à parcela 2
- E assim sucessivamente

### RN04 - Consolidação de Despesas por Categoria
O total de uma categoria deve somar:
- Despesas gerais da categoria no mês
- Parcelas de cartão da categoria na fatura do mês
- Assinaturas ativas da categoria
- Contas fixas ativas da categoria

### RN05 - Status de Assinaturas e Contas
- Status "Ativa": incluída em todos os cálculos e projeções
- Status "Cancelada/Inativa": mantida no histórico mas não incluída em cálculos futuros

### RN06 - Cálculo de Saldo Mensal
```
TOTAL_RENDAS_OPERACIONAIS = renda_fixa + renda_variável_lançada + renda_extra

TOTAL_RENDIMENTOS_RECEBIDOS = soma(valor) de todos os rendimentos do mês

TOTAL_DISTRIBUÍDO = soma(valor_distribuido) de todos os rendimentos do mês
                    (padrão 0 se não distribuído)

TOTAL_DESPESAS = despesas_gerais + assinaturas_ativas
               + contas_fixas_ativas + fatura_cartões_do_mês

SALDO_MÊS = (TOTAL_RENDAS_OPERACIONAIS + TOTAL_DISTRIBUÍDO) − TOTAL_DESPESAS
```

> Rendimentos de investimentos **não distribuídos** são informativos e não impactam o saldo do mês — representam reinvestimento.

### RN07 - Projeção de Rendas Futuras
- Apenas rendas fixas ativas são incluídas automaticamente nas projeções
- Rendas variáveis e extras não são projetadas (valor desconhecido)
- A projeção de saldo futuro usa: rendas fixas − despesas projetadas

### RN09 - Tratamento de Rendimentos de Investimentos
- Rendimentos de investimentos são sempre registrados com seu valor real recebido.
- `valor_distribuido` representa a parcela que entra no orçamento do mês (padrão: 0).
- Somente `valor_distribuido` é incluído no `SALDO_MÊS`.
- O total recebido não distribuído é exibido como informação de reinvestimento, sem impacto no saldo.
- Projeções futuras **não incluem** rendimentos de investimentos (nem distribuídos), pois são imprevisíveis.

### RN08 - Consolidação Familiar vs. Individual
- A visão familiar soma todas as rendas e despesas de todos os membros ativos
- A visão individual filtra rendas e despesas por membro
- Cartões, assinaturas e contas fixas herdam o membro do seu cadastro

---

## 5. Requisitos Não-Funcionais

### RNF01 - Usabilidade
- Interface intuitiva e fácil de usar, com foco em entrada rápida de dados (poucos cliques/toques)
- Design responsivo (mobile-first), compatível com telas de 320px a 1920px
- Navegação fluida em dispositivos touch (celular e tablet)
- Feedback visual imediato após cada ação (salvar, excluir, calcular)

### RNF02 - Performance
- Carregamento do dashboard em menos de 2 segundos em conexão normal
- Cálculos de saldo, parcelas e projeções em tempo real, sem recarregar a página
- Interface responsiva mesmo em conexões lentas (modo offline-first)

### RNF03 - Armazenamento e Sincronização
- Dados armazenados localmente no dispositivo (offline-first): o sistema deve funcionar sem conexão à internet
- Sincronização automática com servidor cloud quando conexão disponível
- Sincronização entre os dispositivos dos dois usuários, com resolução de conflitos simples (última escrita vence, por padrão)
- Backup automático periódico dos dados no servidor cloud
- Exportação de dados em CSV e Excel para uso externo

### RNF04 - Segurança e Autenticação
- Sistema de autenticação obrigatório com suporte a 2 usuários (casal), cada um com credenciais próprias
- Dados compartilhados entre os dois usuários autenticados
- Comunicação cliente-servidor exclusivamente via HTTPS
- Tokens de sessão com expiração configurável
- Dados financeiros criptografados em trânsito e em repouso no servidor

### RNF05 - Disponibilidade e Offline
- O sistema deve operar em modo offline com todas as funcionalidades de lançamento e consulta
- Ao recuperar conexão, a sincronização deve ocorrer automaticamente em segundo plano
- Em caso de conflito de dados entre dispositivos, o usuário deve ser notificado e ter a opção de resolver manualmente

### RNF06 - Compatibilidade
- Web: últimas 2 versões dos navegadores Chrome, Firefox e Safari
- Mobile: Safari no iOS e Chrome no Android (via PWA ou app responsivo)
- Sem necessidade de instalação de aplicativo nativo na versão inicial (PWA suficiente)

### RNF07 - Manutenibilidade
- Código-fonte versionado em repositório Git
- Migrações de banco de dados versionadas e reversíveis
- Separação clara entre frontend, backend e camada de dados

---

## 6. Casos de Uso Principais

### UC01 - Lançar Compra Parcelada no Cartão
**Ator:** Usuário  
**Fluxo Principal:**
1. Usuário acessa módulo de Cartões de Crédito
2. Clica em "Nova Despesa"
3. Preenche:
   - Data: 10/03/2026
   - Cartão: Cartão 1
   - Descrição: "Notebook Dell"
   - Categoria: "Outros"
   - Valor Total: R$ 3.000,00
   - Parcelas: 10x
   - Parcela Atual: 1
4. Sistema calcula:
   - Valor da Parcela: R$ 300,00
   - Fatura: MAR/26 (pois 10/03 < 12)
5. Sistema salva e atualiza dashboard

### UC02 - Cadastrar Nova Assinatura
**Ator:** Usuário  
**Fluxo Principal:**
1. Usuário acessa módulo de Assinaturas
2. Clica em "Nova Assinatura"
3. Preenche:
   - Serviço: "Netflix"
   - Categoria: "Lazer"
   - Valor: R$ 49,90
   - Dia Cobrança: 15
   - Forma Pagamento: "Cartão 1"
   - Status: "Ativa"
4. Sistema salva e inclui nos cálculos mensais

### UC03 - Visualizar Projeção Futura
**Ator:** Usuário  
**Fluxo Principal:**
1. Usuário acessa Dashboard
2. Visualiza seção "Projeção Próximos 3 Meses"
3. Sistema exibe:
   - Abril/26: R$ 4.500,00 (9 parcelas do notebook + fixas + assinaturas)
   - Maio/26: R$ 4.500,00 (8 parcelas + fixas + assinaturas)
   - Junho/26: R$ 4.500,00 (7 parcelas + fixas + assinaturas)

### UC04 - Cancelar Assinatura
**Ator:** Usuário  
**Fluxo Principal:**
1. Usuário acessa lista de Assinaturas
2. Seleciona assinatura "Spotify"
3. Altera status para "Cancelada"
4. Sistema remove dos cálculos futuros mas mantém histórico

---

## 7. Modelo de Dados (Conceitual)

### Entidades Principais

#### Membro
- id
- nome
- relacionamento (opcional)
- ativo

#### Categoria
- id
- nome
- ativo

#### CartaoCredito
- id
- nome
- membro_id (FK)
- dia_fechamento
- dia_vencimento
- limite
- ativo

#### DespesaCartao
- id
- data
- cartao_id (FK)
- descricao
- categoria_id (FK)
- valor_total
- numero_parcelas
- parcela_atual
- valor_parcela (calculado)
- fatura (calculado: "MMM/AA")

#### Assinatura
- id
- servico
- membro_id (FK)
- categoria_id (FK)
- valor_mensal
- dia_cobranca
- forma_pagamento
- status (ativa/cancelada)
- observacoes

#### ContaFixa
- id
- descricao
- membro_id (FK)
- categoria_id (FK)
- valor
- dia_vencimento
- forma_pagamento
- status (ativa/inativa)
- observacoes

#### DespesaGeral
- id
- data
- membro_id (FK)
- descricao
- categoria_id (FK)
- valor
- forma_pagamento
- observacoes

#### RendaFixa
- id
- descricao
- membro_id (FK)
- valor_mensal
- dia_recebimento
- status (ativa/inativa)
- observacoes

#### RendaVariavel
- id
- descricao
- membro_id (FK)
- mes_ano (referência: "MMM/AA")
- valor
- data_recebimento
- observacoes

#### RendaExtra
- id
- descricao
- membro_id (FK)
- data_recebimento
- valor
- observacoes

#### RendimentoInvestimento
- id
- descricao
- membro_id (FK)
- data_recebimento
- valor
- valor_distribuido (default: 0) — parcela incluída no orçamento do mês
- observacoes

---

## 8. Wireframes (Descrição Textual)

### Tela: Dashboard
```
+------------------------------------------------------------+
| DASHBOARD FINANCEIRO     [Mês: MAR/26] [Membro: Todos ▼]  |
+------------------------------------------------------------+
| RENDAS OPERACIONAIS                                        |
| Rendas Fixas (Salários):      R$ 12.000,00                 |
| Rendas Variáveis:             R$  1.500,00                 |
| Rendas Extras:                R$    500,00                 |
| ========================================================== |
| TOTAL RENDAS OPERACIONAIS:    R$ 14.000,00                 |
+------------------------------------------------------------+
| RENDIMENTOS DE INVESTIMENTOS                               |
| Total Recebido:               R$    500,00                 |
| Total Distribuído:            R$    200,00  [editar]       |
| Total Reinvestido:            R$    300,00  (informativo)  |
+------------------------------------------------------------+
| DESPESAS DO MÊS                                            |
| Total Despesas Gerais:        R$  2.500,00                 |
| Total Assinaturas:            R$    150,00                 |
| Total Contas Fixas:           R$  4.000,00                 |
| Total Cartões (Fatura Atual): R$  1.800,00                 |
| ========================================================== |
| TOTAL DE DESPESAS:            R$  8.450,00                 |
+------------------------------------------------------------+
| SALDO DO MÊS:                 R$  5.750,00  ✅             |
| (R$14.000 + R$200 distribuído − R$8.450)                   |
+------------------------------------------------------------+
|                                                            |
| DESPESAS POR CATEGORIA                                     |
| [Gráfico Pizza ou Barras]                                  |
| Alimentação:    R$ 1.500 (18%)                             |
| Moradia:        R$ 4.000 (47%)                             |
| Transporte:     R$   800  (9%)                             |
| ...                                                        |
+------------------------------------------------------------+
|                                                            |
| EVOLUÇÃO MENSAL                                            |
| [Gráfico Linha: JAN, FEV, MAR — Rendas / Despesas / Saldo]|
+------------------------------------------------------------+
|                                                            |
| PROJEÇÃO PRÓXIMOS 3 MESES                                  |
|          | Rendas   | Despesas | Saldo                     |
| ABR/26:  | R$12.000 | R$ 8.200 | R$ 3.800                  |
| MAI/26:  | R$12.000 | R$ 8.100 | R$ 3.900                  |
| JUN/26:  | R$12.000 | R$ 8.000 | R$ 4.000                  |
+------------------------------------------------------------+
```

### Tela: Cartões de Crédito
```
+--------------------------------------------------+
| CARTÕES DE CRÉDITO              [+ Nova Despesa] |
+--------------------------------------------------+
| Filtros: [Cartão: Todos ▼] [Fatura: MAR/26 ▼]   |
+--------------------------------------------------+
| Data       | Cartão   | Descrição      | Valor  | Parc. |
|------------|----------|----------------|--------|-------|
| 10/03/2026 | Cartão 1 | Notebook Dell  | 300,00 | 1/10  |
| 15/03/2026 | Cartão 2 | Supermercado   | 450,00 | 1/1   |
| 20/03/2026 | Cartão 1 | Restaurante    | 280,00 | 1/1   |
+--------------------------------------------------+
| TOTAL FATURA MAR/26 - Cartão 1:    R$ 580,00    |
| TOTAL FATURA MAR/26 - Cartão 2:    R$ 450,00    |
+--------------------------------------------------+
```

---

## 9. Tecnologias Sugeridas

### Frontend
- **Web:** React.js ou Vue.js
- **Mobile:** React Native ou Flutter
- **Desktop:** Electron (se necessário)

### Backend (se cloud)
- **API:** Node.js (Express) ou Python (FastAPI/Django)
- **Banco de Dados:** PostgreSQL ou MongoDB

### Local (se offline-first)
- **Armazenamento:** SQLite ou IndexedDB
- **Sincronização:** Opcional com cloud

---

## 10. Roadmap Sugerido

### Fase 1 - MVP (Mínimo Produto Viável)
- Cadastro de membros da família
- Cadastro de categorias
- Cadastro de cartões
- Lançamento de despesas em cartão (sem parcelamento)
- Cadastro de rendas fixas
- Dashboard básico com totais de rendas, despesas e saldo

### Fase 2 - Funcionalidades Core
- Parcelamento de compras no cartão
- Assinaturas
- Contas fixas
- Despesas gerais
- Rendas variáveis e extras
- Rendimentos de investimentos
- Dashboard completo com saldo e visão por membro

### Fase 3 - Análises Avançadas
- Gráficos e relatórios detalhados
- Projeções futuras (despesas, rendas e saldo estimado)
- Comparativos mensais
- Metas de gastos por categoria

### Fase 4 - Extras
- Integração bancária (Open Banking)
- Alertas e notificações
- Importação de extratos
- Relatórios personalizados
- App mobile

---

## 11. Considerações Finais

Este sistema visa oferecer controle completo das finanças familiares com foco especial em:
- Consolidação de rendas de múltiplos membros da família
- Gestão inteligente de múltiplos cartões de crédito
- Controle detalhado de parcelamentos
- Visibilidade de despesas recorrentes
- Balanço mensal real (rendas − despesas = saldo)
- Projeções futuras confiáveis incluindo rendas e despesas

A implementação pode começar com uma versão web simples e evoluir para aplicativos mobile e integrações bancárias conforme a necessidade e feedback dos usuários.

---

**Documento sujeito a revisões e atualizações conforme evolução do projeto.**
