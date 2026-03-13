package repository

import (
	"fmt"
	"log"

	"github.com/felipehrs/financas-familiares/backend/config"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// SeedUsuarios cria os 2 usuários iniciais se não existirem.
// Após criar cada usuário, cria (ou reutiliza) a família compartilhada e chama SeedCategorias.
// Os dois usuários de dev compartilham a mesma família: user1 como admin, user2 como membro.
func SeedUsuarios(db *sqlx.DB, cfg *config.SeedConfig) error {
	usuarios := []struct {
		Email string
		Nome  string
		Senha string
	}{
		{cfg.User1Email, cfg.User1Nome, cfg.User1Senha},
		{cfg.User2Email, cfg.User2Nome, cfg.User2Senha},
	}

	var familiaIDCompartilhada string

	for i, u := range usuarios {
		if u.Email == "" || u.Senha == "" {
			log.Printf("seed: usuário com email ou senha vazia ignorado, verifique as variáveis SEED_USER*")
			continue
		}

		// Verificar se já existe
		var usuarioID string
		err := db.QueryRowx(`SELECT id FROM usuarios WHERE email = $1`, u.Email).Scan(&usuarioID)
		if err == nil {
			log.Printf("seed: usuário %s já existe, ignorando criação", u.Email)

			// Garantir que este usuário tenha família
			var famID string
			err2 := db.QueryRowx(`SELECT familia_id FROM familia_usuarios WHERE usuario_id = $1 LIMIT 1`, usuarioID).Scan(&famID)
			if err2 == nil {
				if familiaIDCompartilhada == "" {
					familiaIDCompartilhada = famID
				}
			}
			continue
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.Senha), 12)
		if err != nil {
			return fmt.Errorf("seed: erro ao hashear senha de %s: %w", u.Email, err)
		}

		err = db.QueryRowx(`
			INSERT INTO usuarios (nome, email, senha_hash)
			VALUES ($1, $2, $3)
			RETURNING id
		`, u.Nome, u.Email, string(hash)).Scan(&usuarioID)
		if err != nil {
			return fmt.Errorf("seed: erro ao inserir usuário %s: %w", u.Email, err)
		}

		log.Printf("seed: usuário %s criado com sucesso (id=%s)", u.Email, usuarioID)

		if i == 0 {
			// Primeiro usuário: criar a família compartilhada
			err = db.QueryRowx(`
				INSERT INTO familias (nome, owner_id)
				VALUES ('Família Principal', $1)
				RETURNING id
			`, usuarioID).Scan(&familiaIDCompartilhada)
			if err != nil {
				return fmt.Errorf("seed: erro ao criar família para %s: %w", u.Email, err)
			}
			log.Printf("seed: família criada (id=%s)", familiaIDCompartilhada)
		}

		if familiaIDCompartilhada == "" {
			// Família não criada ainda (user1 já existia mas sem família) — criar agora
			err = db.QueryRowx(`
				INSERT INTO familias (nome, owner_id)
				VALUES ('Família Principal', $1)
				RETURNING id
			`, usuarioID).Scan(&familiaIDCompartilhada)
			if err != nil {
				return fmt.Errorf("seed: erro ao criar família para %s: %w", u.Email, err)
			}
		}

		role := "membro"
		if i == 0 {
			role = "admin"
		}

		_, err = db.Exec(`
			INSERT INTO familia_usuarios (familia_id, usuario_id, role)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING
		`, familiaIDCompartilhada, usuarioID, role)
		if err != nil {
			return fmt.Errorf("seed: erro ao vincular usuário %s à família: %w", u.Email, err)
		}
	}

	if familiaIDCompartilhada != "" {
		if err := SeedCategorias(db, familiaIDCompartilhada); err != nil {
			return err
		}
	}

	return nil
}

// SeedCategorias insere as categorias padrão da família se não existirem.
func SeedCategorias(db *sqlx.DB, familiaID string) error {
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
			INSERT INTO categorias (nome, familia_id)
			VALUES ($1, $2)
			ON CONFLICT (familia_id, nome) DO NOTHING
		`, nome, familiaID)
		if err != nil {
			return fmt.Errorf("seed: erro ao inserir categoria %s: %w", nome, err)
		}
	}

	log.Printf("seed: categorias verificadas/inseridas para família %s", familiaID)
	return nil
}
