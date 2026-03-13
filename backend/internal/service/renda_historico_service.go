package service

import (
	"github.com/felipehrs/financas-familiares/backend/internal/domain"
)

// RendaFixaRepositoryForHistorico define os métodos de renda fixa usados pelo RendaHistoricoService.
type RendaFixaRepositoryForHistorico interface {
	Listar(familiaID string) ([]*domain.RendaFixa, error)
	ListarVigentesPorMes(familiaID string, mes, ano int) ([]*domain.RendaFixa, error)
}

// RendaVariavelRepositoryForHistorico define os métodos de renda variável usados pelo RendaHistoricoService.
type RendaVariavelRepositoryForHistorico interface {
	Listar(familiaID string) ([]*domain.RendaVariavel, error)
	ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendaVariavel, error)
}

// RendaExtraRepositoryForHistorico define os métodos de renda extra usados pelo RendaHistoricoService.
type RendaExtraRepositoryForHistorico interface {
	Listar(familiaID string) ([]*domain.RendaExtra, error)
	ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendaExtra, error)
}

// RendimentoRepositoryForHistorico define os métodos de rendimento de investimento usados pelo RendaHistoricoService.
type RendimentoRepositoryForHistorico interface {
	Listar(familiaID string) ([]*domain.RendimentoInvestimento, error)
	ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendimentoInvestimento, error)
}

// RendaHistoricoService implementa a lógica de busca e filtro do histórico de rendas.
type RendaHistoricoService struct {
	rendaFixaRepo     RendaFixaRepositoryForHistorico
	rendaVariavelRepo RendaVariavelRepositoryForHistorico
	rendaExtraRepo    RendaExtraRepositoryForHistorico
	rendimentoRepo    RendimentoRepositoryForHistorico
}

// NewRendaHistoricoService cria uma nova instância do RendaHistoricoService.
func NewRendaHistoricoService(
	rendaFixaRepo RendaFixaRepositoryForHistorico,
	rendaVariavelRepo RendaVariavelRepositoryForHistorico,
	rendaExtraRepo RendaExtraRepositoryForHistorico,
	rendimentoRepo RendimentoRepositoryForHistorico,
) *RendaHistoricoService {
	return &RendaHistoricoService{
		rendaFixaRepo:     rendaFixaRepo,
		rendaVariavelRepo: rendaVariavelRepo,
		rendaExtraRepo:    rendaExtraRepo,
		rendimentoRepo:    rendimentoRepo,
	}
}

// BuscarHistorico retorna o histórico de rendas aplicando os filtros informados.
// Regras de filtragem por data:
//   - mes+ano informados: usa ListarPorMes (variavel/extra/investimento) ou ListarVigentesPorMes (fixa)
//   - apenas ano informado: usa Listar() e filtra em memória
//   - sem filtro de data: usa Listar()
//
// O filtro membro_id e tipo são sempre aplicados em memória após a busca.
func (s *RendaHistoricoService) BuscarHistorico(familiaID string, filtro domain.FiltroHistoricoRendas) (*domain.HistoricoRendas, error) {
	itens := make([]domain.ItemRendaHistorico, 0)

	comMesEAno := filtro.Mes != 0 && filtro.Ano != 0
	comApenasAno := filtro.Mes == 0 && filtro.Ano != 0

	incluirFixa := filtro.Tipo == "" || filtro.Tipo == domain.TipoRendaFixa
	incluirVariavel := filtro.Tipo == "" || filtro.Tipo == domain.TipoRendaVariavel
	incluirExtra := filtro.Tipo == "" || filtro.Tipo == domain.TipoRendaExtra
	incluirInvestimento := filtro.Tipo == "" || filtro.Tipo == domain.TipoRendaInvestimento

	// --- Renda Fixa ---
	if incluirFixa {
		var rendas []*domain.RendaFixa
		var err error
		if comMesEAno {
			rendas, err = s.rendaFixaRepo.ListarVigentesPorMes(familiaID, filtro.Mes, filtro.Ano)
		} else {
			rendas, err = s.rendaFixaRepo.Listar(familiaID)
		}
		if err != nil {
			return nil, err
		}

		for _, r := range rendas {
			if filtro.MembroID != "" && r.MembroID != filtro.MembroID {
				continue
			}
			if comApenasAno {
				// para fixa, vigente no ano: inicio.Year <= ano && (fim == nil || fim.Year >= ano)
				if r.DataInicio.Year() > filtro.Ano {
					continue
				}
				if r.DataFim != nil && r.DataFim.Year() < filtro.Ano {
					continue
				}
			}

			var valor float64
			if comMesEAno {
				valor = ValorProporcionado(r, filtro.Mes, filtro.Ano)
			} else {
				valor = r.Valor
			}

			ativa := r.Ativa
			item := domain.ItemRendaHistorico{
				ID:        r.ID,
				Tipo:      domain.TipoRendaFixa,
				Descricao: r.Descricao,
				MembroID:  r.MembroID,
				Valor:     valor,
				Ativa:     &ativa,
			}
			if comMesEAno {
				mes := filtro.Mes
				ano := filtro.Ano
				item.Mes = &mes
				item.Ano = &ano
			}
			itens = append(itens, item)
		}
	}

	// --- Renda Variável ---
	if incluirVariavel {
		var rendas []*domain.RendaVariavel
		var err error
		if comMesEAno {
			rendas, err = s.rendaVariavelRepo.ListarPorMes(familiaID, filtro.Mes, filtro.Ano)
		} else {
			rendas, err = s.rendaVariavelRepo.Listar(familiaID)
		}
		if err != nil {
			return nil, err
		}

		for _, r := range rendas {
			if filtro.MembroID != "" && r.MembroID != filtro.MembroID {
				continue
			}
			if comApenasAno && r.AnoReferencia != filtro.Ano {
				continue
			}
			mes := r.MesReferencia
			ano := r.AnoReferencia
			item := domain.ItemRendaHistorico{
				ID:        r.ID,
				Tipo:      domain.TipoRendaVariavel,
				Descricao: r.Descricao,
				MembroID:  r.MembroID,
				Valor:     r.Valor,
				Mes:       &mes,
				Ano:       &ano,
			}
			itens = append(itens, item)
		}
	}

	// --- Renda Extra ---
	if incluirExtra {
		var rendas []*domain.RendaExtra
		var err error
		if comMesEAno {
			rendas, err = s.rendaExtraRepo.ListarPorMes(familiaID, filtro.Mes, filtro.Ano)
		} else {
			rendas, err = s.rendaExtraRepo.Listar(familiaID)
		}
		if err != nil {
			return nil, err
		}

		for _, r := range rendas {
			if filtro.MembroID != "" && r.MembroID != filtro.MembroID {
				continue
			}
			if comApenasAno && r.DataRecebimento.Year() != filtro.Ano {
				continue
			}
			dataStr := r.DataRecebimento.Format("2006-01-02")
			item := domain.ItemRendaHistorico{
				ID:        r.ID,
				Tipo:      domain.TipoRendaExtra,
				Descricao: r.Descricao,
				MembroID:  r.MembroID,
				Valor:     r.Valor,
				Data:      &dataStr,
			}
			itens = append(itens, item)
		}
	}

	// --- Rendimento de Investimento ---
	if incluirInvestimento {
		var rendimentos []*domain.RendimentoInvestimento
		var err error
		if comMesEAno {
			rendimentos, err = s.rendimentoRepo.ListarPorMes(familiaID, filtro.Mes, filtro.Ano)
		} else {
			rendimentos, err = s.rendimentoRepo.Listar(familiaID)
		}
		if err != nil {
			return nil, err
		}

		for _, r := range rendimentos {
			if filtro.MembroID != "" && r.MembroID != filtro.MembroID {
				continue
			}
			if comApenasAno && r.Data.Year() != filtro.Ano {
				continue
			}
			dataStr := r.Data.Format("2006-01-02")
			vd := r.ValorDistribuido
			item := domain.ItemRendaHistorico{
				ID:               r.ID,
				Tipo:             domain.TipoRendaInvestimento,
				Descricao:        r.Descricao,
				MembroID:         r.MembroID,
				Valor:            r.Valor, // informativo
				Data:             &dataStr,
				ValorDistribuido: &vd,
			}
			itens = append(itens, item)
		}
	}

	// --- Resumo ---
	resumo := calcularResumo(itens)

	return &domain.HistoricoRendas{
		Itens:  itens,
		Resumo: resumo,
	}, nil
}

// calcularResumo agrega os totais por tipo de renda a partir da lista de itens.
func calcularResumo(itens []domain.ItemRendaHistorico) domain.ResumoRendaHistorico {
	var resumo domain.ResumoRendaHistorico
	for _, item := range itens {
		switch item.Tipo {
		case domain.TipoRendaFixa:
			resumo.TotalFixas += item.Valor
		case domain.TipoRendaVariavel:
			resumo.TotalVariaveis += item.Valor
		case domain.TipoRendaExtra:
			resumo.TotalExtras += item.Valor
		case domain.TipoRendaInvestimento:
			if item.ValorDistribuido != nil {
				resumo.TotalInvestimentosDistribuidos += *item.ValorDistribuido
			}
		}
	}
	resumo.TotalGeral = resumo.TotalFixas + resumo.TotalVariaveis + resumo.TotalExtras + resumo.TotalInvestimentosDistribuidos
	return resumo
}
