package service

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/rafaeldepontes/voting-go/internal/auth"
	"github.com/rafaeldepontes/voting-go/internal/auth/model"
	"github.com/rafaeldepontes/voting-go/internal/utils"
)

type service struct {
	r auth.Repository
}

func NewService(r auth.Repository) auth.Service {
	return &service{
		r,
	}
}

// Login implements [auth.Service].
func (s *service) Login(userReq model.UserReq) (int64, error) {
	if userReq.Email == "" {
		return 0, utils.ErrInvalidEmail
	}

	if userReq.Password == "" || len(userReq.Password) < 5 {
		return 0, utils.ErrInvalidPassword
	}

	log.Printf("[INFO] log in attempt to email: %s\n", userReq.Email)

	user, err := s.r.FindByEmail(userReq.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, utils.ErrInvalidLogin
		}
		return 0, utils.ErrGenericError
	}

	match, err := VerifyPassword(userReq.Password, user.HashedPassword)
	if err != nil {
		log.Printf("[ERROR] password verification failed: %v\n", err)
		return 0, utils.ErrGenericError
	}

	if !match {
		return 0, utils.ErrInvalidCredentials
	}

	return user.ID, nil
}

// Register implements [auth.Service].
func (s *service) Register(userReq model.UserReq) error {
	if userReq.Email == "" {
		return utils.ErrInvalidEmail
	}

	if userReq.Password == "" || len(userReq.Password) < 5 {
		return utils.ErrInvalidPassword
	}

	log.Printf("[INFO] registering a new user, email: %s\n", userReq.Email)

	_, err := s.r.FindByEmail(userReq.Email)
	if err == nil {
		return utils.ErrInvalidLogin
	}

	hp, err := HashPassword(userReq.Password)
	if err != nil {
		log.Printf("[ERROR] failed to hash password: %v\n", err)
		return utils.ErrGenericError
	}

	user := model.User{
		Email:          userReq.Email,
		HashedPassword: hp,
		CreatedAt:      time.Now(),
	}
	if err := s.r.Save(user); err != nil {
		return err
	}

	return nil
}
