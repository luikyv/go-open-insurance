package user

import (
	"context"
	"net/http"

	"github.com/luikyv/go-open-insurance/internal/api"
)

type Service struct {
	storage *Storage
}

func NewService(storage *Storage) Service {
	return Service{
		storage: storage,
	}
}

func (s Service) Create(ctx context.Context, user User) error {
	return s.storage.create(ctx, user)
}

func (s Service) User(username string) (User, error) {
	user, err := s.storage.user(username)
	if err != nil {
		return User{}, api.NewError("USER_NOT_FOUND", http.StatusNotFound, "user not found")
	}

	return user, nil
}

func (s Service) UserByCPF(cpf string) (User, error) {
	user, err := s.storage.userByCPF(cpf)
	if err != nil {
		return User{}, api.NewError("USER_NOT_FOUND", http.StatusNotFound, "cpf not found")
	}

	return user, nil
}
