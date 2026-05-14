package repository

import (
	"database/sql"
	"log"

	"github.com/rafaeldepontes/voting-go/internal/auth"
	"github.com/rafaeldepontes/voting-go/internal/auth/model"
	"github.com/rafaeldepontes/voting-go/pkg/database/postgres"
)

type repository struct {
	db *sql.DB
}

func NewRepository() auth.Repository {
	return repository{
		db: postgres.GetDb(),
	}
}

// FindByEmail implements [auth.Repository].
func (r repository) FindByEmail(email string) (model.User, error) {
	query := `
		SELECT id, email, hashed_password, created_at FROM users u where u.email = $1
	`
	var user model.User
	if err := r.db.QueryRow(
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.HashedPassword,
		&user.CreatedAt,
	); err != nil {
		if err != sql.ErrNoRows {
			log.Println("[ERROR] database error finding user, ", err)
		}
		return model.User{}, err
	}

	return user, nil
}

// Save implements [auth.Repository].
func (r repository) Save(user model.User) error {
	query := `
		INSERT INTO users (email, hashed_password, created_at) VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, &user.Email, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		log.Printf("[ERROR] couldn't insert user: %s because %v\n", user.Email, err)
		return err
	}
	return nil
}
