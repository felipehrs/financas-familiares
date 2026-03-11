package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

// RendaFixaRepository implementa service.RendaFixaRepository usando PostgreSQL via sqlx.
type RendaFixaRepository struct {
	db *sqlx.DB
}

// NewRendaFixaRepository cria uma nova instância do RendaFixaRepository.
func NewRendaFixaRepository(db *sqlx.DB) *RendaFixaRepository {
	return &RendaFixaRepository{db: db}
}

// rendaFixaRow é a estrutura de scan para as colunas do banco.
type rendaFixaRow struct {
	ID             string     `db:"id"`
	Descricao      string     `db:"descricao"`
	MembroID       string     `db:"membro_id"`
	Valor          float64    `db:"valor"`
	DiaRecebimento int        `db:"dia_recebimento"`
	Ativa          bool       `db:"ativa"`
	DataInicio     time.Time  `db:"data_inicio"`
	DataFim        *time.Time `db:"data_fim"`
}

func (r rendaFixaRow) toDomain() *domain.RendaFixa {
	return &domain.RendaFixa{
		ID:             r.ID,
		Descricao:      r.Descricao,
		MembroID:       r.MembroID,
		Valor:          r.Valor,
		DiaRecebimento: r.DiaRecebimento,
		Ativa:          r.Ativa,
		DataInicio:     r.DataInicio,
		DataFim:        r.DataFim,
	}
}

// Criar persiste uma nova renda fixa no banco e retorna o registro criado (com ID gerado).
func (r *RendaFixaRepository) Criar(familiaID string, renda *domain.RendaFixa) (*domain.RendaFixa, error) {
	var row rendaFixaRow
	err := r.db.QueryRowx(`
		INSERT INTO rendas_fixas (familia_id, descricao, membro_id, valor, dia_recebimento, ativa, data_inicio, data_fim)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, descricao, membro_id, valor, dia_recebimento, ativa, data_inicio, data_fim
	`, familiaID, renda.Descricao, renda.MembroID, renda.Valor, renda.DiaRecebimento, renda.Ativa, renda.DataInicio, renda.DataFim).StructScan(&row)
	if err != nil {
		return nil, err
	}
	return row.toDomain(), nil
}

// BuscarPorID busca uma renda fixa pelo seu UUID.
// Retorna domain.ErrRendaFixaNaoEncontrada se não existir ou estiver soft-deleted.
func (r *RendaFixaRepository) BuscarPorID(familiaID, id string) (*domain.RendaFixa, error) {
	var row rendaFixaRow
	err := r.db.QueryRowx(`
		SELECT id, descricao, membro_id, valor, dia_recebimento, ativa, data_inicio, data_fim
		FROM rendas_fixas
		WHERE id = $1 AND familia_id = $2 AND deleted_at IS NULL
	`, id, familiaID).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRendaFixaNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Listar retorna todas as rendas fixas não excluídas (ativas e inativas).
func (r *RendaFixaRepository) Listar(familiaID string) ([]*domain.RendaFixa, error) {
	var rows []rendaFixaRow
	err := r.db.Select(&rows, `
		SELECT id, descricao, membro_id, valor, dia_recebimento, ativa, data_inicio, data_fim
		FROM rendas_fixas
		WHERE familia_id = $1 AND deleted_at IS NULL
		ORDER BY descricao ASC
	`, familiaID)
	if err != nil {
		return nil, err
	}

	rendas := make([]*domain.RendaFixa, 0, len(rows))
	for _, row := range rows {
		rendas = append(rendas, row.toDomain())
	}
	return rendas, nil
}

// ListarAtivas retorna apenas as rendas fixas ativas (para cálculo de saldo e projeções).
func (r *RendaFixaRepository) ListarAtivas(familiaID string) ([]*domain.RendaFixa, error) {
	var rows []rendaFixaRow
	err := r.db.Select(&rows, `
		SELECT id, descricao, membro_id, valor, dia_recebimento, ativa, data_inicio, data_fim
		FROM rendas_fixas
		WHERE familia_id = $1 AND ativa = TRUE AND deleted_at IS NULL
		ORDER BY descricao ASC
	`, familiaID)
	if err != nil {
		return nil, err
	}

	rendas := make([]*domain.RendaFixa, 0, len(rows))
	for _, row := range rows {
		rendas = append(rendas, row.toDomain())
	}
	return rendas, nil
}

// ListarVigentesPorMes retorna as rendas fixas ativas e vigentes para o mês/ano informado.
// Uma renda é vigente se: data_inicio (mês/ano) <= alvo AND (data_fim IS NULL OR data_fim (mês/ano) >= alvo).
// Compara usando: (YEAR * 12 + MONTH) como inteiro para granularidade mensal.
func (r *RendaFixaRepository) ListarVigentesPorMes(familiaID string, mes, ano int) ([]*domain.RendaFixa, error) {
	var rows []rendaFixaRow
	err := r.db.Select(&rows, `
		SELECT id, descricao, membro_id, valor, dia_recebimento, ativa, data_inicio, data_fim
		FROM rendas_fixas
		WHERE familia_id = $1
		  AND ativa = TRUE
		  AND deleted_at IS NULL
		  AND (EXTRACT(YEAR FROM data_inicio) * 12 + EXTRACT(MONTH FROM data_inicio)) <= ($2 * 12 + $3)
		  AND (data_fim IS NULL OR (EXTRACT(YEAR FROM data_fim) * 12 + EXTRACT(MONTH FROM data_fim)) >= ($2 * 12 + $3))
		ORDER BY descricao ASC
	`, familiaID, ano, mes)
	if err != nil {
		return nil, err
	}

	rendas := make([]*domain.RendaFixa, 0, len(rows))
	for _, row := range rows {
		rendas = append(rendas, row.toDomain())
	}
	return rendas, nil
}

// Atualizar atualiza os dados de uma renda fixa existente e retorna o registro atualizado.
func (r *RendaFixaRepository) Atualizar(familiaID string, renda *domain.RendaFixa) (*domain.RendaFixa, error) {
	var row rendaFixaRow
	err := r.db.QueryRowx(`
		UPDATE rendas_fixas
		SET descricao = $3, membro_id = $4, valor = $5, dia_recebimento = $6, ativa = $7,
		    data_inicio = $8, data_fim = $9, updated_at = NOW()
		WHERE id = $2 AND familia_id = $1 AND deleted_at IS NULL
		RETURNING id, descricao, membro_id, valor, dia_recebimento, ativa, data_inicio, data_fim
	`, familiaID, renda.ID, renda.Descricao, renda.MembroID, renda.Valor, renda.DiaRecebimento, renda.Ativa, renda.DataInicio, renda.DataFim).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRendaFixaNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Inativar marca uma renda fixa como inativa sem excluir o registro.
func (r *RendaFixaRepository) Inativar(familiaID, id string) error {
	result, err := r.db.Exec(`
		UPDATE rendas_fixas
		SET ativa = FALSE, updated_at = NOW()
		WHERE id = $2 AND familia_id = $1 AND deleted_at IS NULL
	`, familiaID, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrRendaFixaNaoEncontrada
	}
	return nil
}
