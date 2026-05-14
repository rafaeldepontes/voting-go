package auth

import "github.com/rafaeldepontes/voting-go/internal/auth/model"

type Repository interface {
	FindByEmail(email string) (model.User, error)
	Save(user model.User) error
}
