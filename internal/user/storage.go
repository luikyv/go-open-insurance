package user

import "github.com/luikyv/go-opf/internal/slice"

type Storage struct {
	users []User
}

func NewStorage() Storage {
	return Storage{
		users: []User{
			userBob,
		},
	}
}

func (st Storage) user(username string) (User, error) {
	user, ok := slice.FindFirst(st.users, func(user User) bool {
		return user.UserName == username
	})
	if !ok {
		return User{}, errorUserNotFound
	}
	return user, nil
}

func (st Storage) userByCPF(cpf string) (User, error) {
	user, ok := slice.FindFirst(st.users, func(user User) bool {
		return user.CPF == cpf
	})
	if !ok {
		return User{}, errorUserNotFound
	}

	return user, nil
}
