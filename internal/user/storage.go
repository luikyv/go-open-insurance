package user

import (
	"context"
	"errors"
)

type Storage struct {
	users []User
}

func NewStorage() *Storage {
	return &Storage{
		users: []User{},
	}
}

func (st *Storage) create(_ context.Context, user User) error {
	_, ok := findFirst(st.users, func(u User) bool {
		return u.UserName == user.UserName
	})

	if ok {
		return errors.New("user already exists")
	}

	st.users = append(st.users, user)
	return nil
}

func (st *Storage) user(username string) (User, error) {
	user, ok := findFirst(st.users, func(user User) bool {
		return user.UserName == username
	})
	if !ok {
		return User{}, errorUserNotFound
	}
	return user, nil
}

func (st *Storage) userByCPF(cpf string) (User, error) {
	user, ok := findFirst(st.users, func(user User) bool {
		return user.CPF == cpf
	})
	if !ok {
		return User{}, errorUserNotFound
	}

	return user, nil
}

func findFirst[T any](elements []T, condition func(t T) bool) (T, bool) {
	for _, e := range elements {
		if condition(e) {
			return e, true
		}
	}

	return *new(T), false
}
