package user

import (
	"context"
	"slices"
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
	return s.storage.user(username)
}

func (s Service) UserByCPF(cpf string) (User, error) {
	return s.storage.userByCPF(cpf)
}

func (s Service) UserBelongsToCompany(user User, businessCNPJ string) bool {
	return slices.Contains(user.CompanyCNPJs, businessCNPJ)
}
