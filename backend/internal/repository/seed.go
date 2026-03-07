package repository

import (
	"fmt"
	"log"

	"github.com/felipehrs/financas-familiares/backend/config"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// SeedUsuarios cria os 2 usuários iniciais se não existirem.
// As senhas são hasheadas com bcrypt (custo 12) antes de inserir.
func SeedUsuarios(db *sqlx.DB, cfg *config.SeedConfig) error {
	usuarios := []struct {
		Email string
		Nome  string
		Senha string
	}{
		{cfg.User1Email, cfg.User1Nome, cfg.User1Senha},
		{cfg.User2Email, cfg.User2Nome, cfg.User2Senha},
	}

	for _, u := range usuarios {
		if u.Email == "" || u.Senha == "" {
			log.Printf("seed: usuário com email ou senha vazia ignorado, verifique as variáveis SEED_USER*")
			continue
		}

		// Verificar se já existe
		var count int
		err := db.Get(&count, `SELECT COUNT(*) FROM usuarios WHERE email = $1`, u.Email)
		if err != nil {
			return fmt.Errorf("seed: erro ao verificar usuário %s: %w", u.Email, err)
		}
		if count > 0 {
			log.Printf("seed: usuário %s já existe, ignorando", u.Email)
			continue
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.Senha), 12)
		if err != nil {
			return fmt.Errorf("seed: erro ao hashear senha de %s: %w", u.Email, err)
		}

		_, err = db.Exec(`
			INSERT INTO usuarios (nome, email, senha_hash)
			VALUES ($1, $2, $3)
		`, u.Nome, u.Email, string(hash))
		if err != nil {
			return fmt.Errorf("seed: erro ao inserir usuário %s: %w", u.Email, err)
		}

		log.Printf("seed: usuário %s criado com sucesso", u.Email)
	}

	return nil
}

// SeedCategorias insere as categorias padrão se não existirem.
func SeedCategorias(db *sqlx.DB) error {
	categorias := []string{
		"Alimentação",
		"Transporte",
		"Lazer",
		"Saúde",
		"Educação",
		"Moradia",
		"Vestuário",
		"Outros",
	}

	for _, nome := range categorias {
		_, err := db.Exec(`
			INSERT INTO categorias (nome)
			VALUES ($1)
			ON CONFLICT DO NOTHING
		`, nome)
		if err != nil {
			return fmt.Errorf("seed: erro ao inserir categoria %s: %w", nome, err)
		}
	}

	log.Printf("seed: categorias verificadas/inseridas com sucesso")
	return nil
}
