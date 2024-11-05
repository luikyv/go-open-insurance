package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/capitalizationtitle"
	"github.com/luikyv/go-open-insurance/internal/customer"
	"github.com/luikyv/go-open-insurance/internal/resource"
	"github.com/luikyv/go-open-insurance/internal/user"
)

func loadMocks(
	userService user.Service,
	customerService customer.Service,
	resourceService resource.Service,
	capitalizationTitleService capitalizationtitle.Service,
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

	customerService.AddPersonalIdentification(
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
	customerService.AddPersonalQualification(
		userBob.UserName,
		api.PersonalQualificationData{
			LifePensionPlans:  api.LifePensionPlanApplicabilityNAOSEAPLICA,
			PepIdentification: api.PoliticalExposureNAOEXPOSTO,
			UpdateDateTime:    dateTimeNow,
		},
	)
	customerService.AddPersonalComplimentaryInfo(
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

	capitalizationTitlePlanID1 := "cbad06ae-5f44-483a-bded-e61593ea195c"
	capitalizationTitleService.AddPlan(
		userBob.UserName,
		api.CapitalizationTitlePlanData{
			Brand: api.CapitalizationTitleBrand{
				Name: "Mock Insurance",
				Companies: []api.CapitalizationTitleCompany{
					{
						CnpjNumber:  "90990354000113",
						CompanyName: "Mock Insurance",
						Products: []api.CapitalizationTitleProduct{
							{
								PlanId:      capitalizationTitlePlanID1,
								ProductName: "Random Capitalization Title",
							},
						},
					},
				},
			},
		},
	)
	capitalizationTitleSeriesID1 := "eb71e4d5-ff97-41ca-923f-efa08536793e"
	capitalizationTitleService.AddPlanInfo(
		userBob.UserName,
		capitalizationTitlePlanID1,
		api.CapitalizationTitlePlanInfo{
			Series: []api.CapitalizationTitleSeries{
				{
					PlanId:            &capitalizationTitlePlanID1,
					SeriesId:          capitalizationTitleSeriesID1,
					Modality:          api.CapitalizationTitleSeriesModalityPOPULAR,
					UpdateIndex:       api.CapitalizationTitleSeriesUpdateIndexIGPM,
					ReadjustmentIndex: api.CapitalizationTitleSeriesReadjustmentIndexIPCA,
				},
			},
		},
	)
	capitalizationTitleService.AddPlanEvent(
		userBob.UserName,
		capitalizationTitlePlanID1,
		api.CapitalizationTitleEvent{
			TitleId: pointerOf("random_title"),
		},
	)
	capitalizationTitleService.AddPlanSettlement(
		userBob.UserName,
		capitalizationTitlePlanID1,
		api.CapitalizationTitleSettlement{
			SettlementId:          "random_settlement",
			SettlementDueDate:     api.DateNow(),
			SettlementPaymentDate: api.DateNow(),
			SettlementFinancialAmount: api.AmountNumberDetails{
				Amount:   100.0,
				Currency: "BRL",
			},
		},
	)

	resourceService.Add(
		userBob.UserName,
		api.ResourceData{
			ResourceId: uuid.NewString(),
			Status:     api.ResourceStatusUNAVAILABLE,
			Type:       api.ResourceTypeCAPITALIZATIONTITLES,
		},
	)
	resourceService.Add(
		userBob.UserName,
		api.ResourceData{
			ResourceId: uuid.NewString(),
			Status:     api.ResourceStatusUNAVAILABLE,
			Type:       api.ResourceTypeCAPITALIZATIONTITLES,
		},
	)
	resourceService.Add(
		userBob.UserName,
		api.ResourceData{
			ResourceId: uuid.NewString(),
			Status:     api.ResourceStatusUNAVAILABLE,
			Type:       api.ResourceTypeCAPITALIZATIONTITLES,
		},
	)
	resourceService.Add(
		userBob.UserName,
		api.ResourceData{
			ResourceId: uuid.NewString(),
			Status:     api.ResourceStatusUNAVAILABLE,
			Type:       api.ResourceTypeCAPITALIZATIONTITLES,
		},
	)
	resourceService.Add(
		userBob.UserName,
		api.ResourceData{
			ResourceId: uuid.NewString(),
			Status:     api.ResourceStatusUNAVAILABLE,
			Type:       api.ResourceTypeCAPITALIZATIONTITLES,
		},
	)
	resourceService.Add(
		userBob.UserName,
		api.ResourceData{
			ResourceId: uuid.NewString(),
			Status:     api.ResourceStatusUNAVAILABLE,
			Type:       api.ResourceTypeCAPITALIZATIONTITLES,
		},
	)
	resourceService.Add(
		userBob.UserName,
		api.ResourceData{
			ResourceId: uuid.NewString(),
			Status:     api.ResourceStatusUNAVAILABLE,
			Type:       api.ResourceTypeCAPITALIZATIONTITLES,
		},
	)
	resourceService.Add(
		userBob.UserName,
		api.ResourceData{
			ResourceId: uuid.NewString(),
			Status:     api.ResourceStatusUNAVAILABLE,
			Type:       api.ResourceTypeCAPITALIZATIONTITLES,
		},
	)
	resourceService.Add(
		userBob.UserName,
		api.ResourceData{
			ResourceId: uuid.NewString(),
			Status:     api.ResourceStatusUNAVAILABLE,
			Type:       api.ResourceTypeCAPITALIZATIONTITLES,
		},
	)

	return nil
}

func pointerOf[T any](t T) *T {
	return &t
}
