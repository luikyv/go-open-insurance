package main

import (
	"context"

	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/customer"
	"github.com/luikyv/go-open-insurance/internal/user"
)

func loadMocks(
	userService user.Service,
	customerService customer.Service,
) error {
	ctx := context.Background()

	dateNow := api.DateNow()
	dateTimeNow := api.DateTimeNow()

	var companyA = user.Company{
		Name: "A Business",
		CNPJ: "27737785000136",
	}
	var userBob = user.User{
		UserName:     "bob@mail.com",
		Email:        "bob@mail.com",
		CPF:          "78628584099",
		Name:         "Mr. Bob",
		CompanyCNPJs: []string{companyA.CNPJ},
	}
	if err := userService.Create(ctx, userBob); err != nil {
		return err
	}

	customerService.SetPersonalIdentifications(
		userBob.UserName,
		api.PersonalIdentificationData{
			CpfNumber:       userBob.CPF,
			BirthDate:       &dateNow,
			CivilName:       userBob.Name,
			CivilStatusCode: pointerOf(api.CivilStatusCodeSOLTEIRO),
			CompanyInfo: api.CompanyInfo{
				CnpjNumber: companyA.CNPJ,
				Name:       companyA.Name,
			},
			Contact: api.PersonalContact{
				Emails: pointerOf([]api.CustomerEmail{
					{Email: pointerOf(userBob.Email)},
				}),
				PostalAddresses: []api.PersonalPostalAddress{
					{
						Address:            "street x, number 1",
						Country:            "BR",
						CountrySubDivision: "SP",
						PostCode:           "00000000",
						TownName:           "São Paulo",
					},
				},
			},
			HasBrazilianNationality: pointerOf(true),
			SocialName:              pointerOf(userBob.Name),
			UpdateDateTime:          dateTimeNow,
		},
	)
	customerService.SetPersonalQualifications(
		userBob.UserName,
		api.PersonalQualificationData{
			LifePensionPlans:  api.LifePensionPlanApplicabilityNAOSEAPLICA,
			PepIdentification: api.PoliticalExposureNAOEXPOSTO,
			UpdateDateTime:    dateTimeNow,
		},
	)
	customerService.SetPersonalComplimentaryInfos(
		userBob.UserName,
		api.PersonalComplimentaryInfoData{
			ProductsServices: []api.ProductService{
				{
					Contract: "1234",
					Type:     api.ProductServiceTypeSEGUROSDEPESSOAS,
				},
			},
			StartDate:      api.NewDate(dateNow.AddDate(0, 0, -1)),
			UpdateDateTime: dateTimeNow,
		},
	)

	return nil
}

func pointerOf[T any](t T) *T {
	return &t
}
