package main

import (
	"context"

	"github.com/luikyv/go-open-insurance/internal/user"
)

var companyA = user.Company{
	Name: "A Business",
	CNPJ: "27737785000136",
}

var userBob = user.User{
	UserName:  "bob@mail.com",
	CPF:       "78628584099",
	Name:      "Mr. Bob",
	Companies: []string{companyA.CNPJ},
}

func loadUsers(userService user.Service) error {
	ctx := context.Background()
	if err := userService.Create(ctx, userBob); err != nil {
		return err
	}

	return nil
}
