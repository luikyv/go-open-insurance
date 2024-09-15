package main

import (
	"context"

	"github.com/luikyv/go-open-insurance/internal/api"
	customersv1 "github.com/luikyv/go-open-insurance/internal/customer/v1"
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

var userBobPersonalIdentificationsV1 = []api.PersonalIdentificationDataV1{}

func loadMocks(
	userService user.Service,
	customerServiceV1 customersv1.Service,
) error {
	ctx := context.Background()
	if err := userService.Create(ctx, userBob); err != nil {
		return err
	}

	customerServiceV1.AddPersonalIdentifications(
		userBob.UserName,
		userBobPersonalIdentificationsV1,
	)

	return nil
}
