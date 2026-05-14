package auth

import "github.com/rafaeldepontes/voting-go/internal/auth/model"

type Service interface {
	Register(userReq model.UserReq) error
	Login(userReq model.UserReq) (int64, error)
}
