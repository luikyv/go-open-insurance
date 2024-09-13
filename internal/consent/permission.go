package consent

import "slices"

const (
	PermissionResourcesRead Permission = "RESOURCES_READ"

	PermissionCustomersPersonalIdentificationRead Permission = "CUSTOMERS_PERSONAL_IDENTIFICATIONS_READ"
	PermissionCustomersPersonalQualificationRead  Permission = "CUSTOMERS_PERSONAL_QUALIFICATION_READ"
	PermissionCustomersPersonalAdditionalInfoRead Permission = "CUSTOMERS_PERSONAL_ADDITIONALINFO_READ"

	PermissionCustomersBusinessIdentificationRead Permission = "CUSTOMERS_BUSINESS_IDENTIFICATIONS_READ"
	PermissionCustomersBusinessQualificationRead  Permission = "CUSTOMERS_BUSINESS_QUALIFICATION_READ"
	PermissionCustomersBusinessAdditionalInfoRead Permission = "CUSTOMERS_BUSINESS_ADDITIONALINFO_READ"

	PermissionCapitalizationTitleRead            Permission = "CAPITALIZATION_TITLE_READ"
	PermissionCapitalizationTitlePlanInfoRead    Permission = "CAPITALIZATION_TITLE_PLANINFO_READ"
	PermissionCapitalizationTitleEventsRead      Permission = "CAPITALIZATION_TITLE_EVENTS_READ"
	PermissionCapitalizationTitleSettlementsRead Permission = "CAPITALIZATION_TITLE_SETTLEMENTS_READ"

	PermissionPensionPlanRead              Permission = "PENSION_PLAN_READ"
	PermissionPensionPlanContractInfoRead  Permission = "PENSION_PLAN_CONTRACTINFO_READ"
	PermissionPensionPlanMovementsRead     Permission = "PENSION_PLAN_MOVEMENTS_READ"
	PermissionPensionPlanPortabilitiesRead Permission = "PENSION_PLAN_PORTABILITIES_READ"
	PermissionPensionPlanWithdrawalsRead   Permission = "PENSION_PLAN_WITHDRAWALS_READ"
	PermissionPensionPlanClaim             Permission = "PENSION_PLAN_CLAIM"

	PermissionLifePensionRead              Permission = "LIFE_PENSION_READ"
	PermissionLifePensionContractInfoRead  Permission = "LIFE_PENSION_CONTRACTINFO_READ"
	PermissionLifePensionMovementsRead     Permission = "LIFE_PENSION_MOVEMENTS_READ"
	PermissionLifePensionPortabilitiesRead Permission = "LIFE_PENSION_PORTABILITIES_READ"
	PermissionLifePensionWithdrawalsRead   Permission = "LIFE_PENSION_WITHDRAWALS_READ"
	PermissionLifePensionClaim             Permission = "LIFE_PENSION_CLAIM"

	PermissionFinancialAssistanceRead             Permission = "FINANCIAL_ASSISTANCE_READ"
	PermissionFinancialAssistanceContractInfoRead Permission = "FINANCIAL_ASSISTANCE_CONTRACTINFO_READ"
	PermissionFinancialAssistanceMovementsRead    Permission = "FINANCIAL_ASSISTANCE_MOVEMENTS_READ"

	PermissionDamagesAndPeoplePatrimonialRead           Permission = "DAMAGES_AND_PEOPLE_PATRIMONIAL_READ"
	PermissionDamagesAndPeoplePatrimonialPolicyInfoRead Permission = "DAMAGES_AND_PEOPLE_PATRIMONIAL_POLICYINFO_READ"
	PermissionDamagesAndPeoplePatrimonialPremiumRead    Permission = "DAMAGES_AND_PEOPLE_PATRIMONIAL_PREMIUM_READ"
	PermissionDamagesAndPeoplePatrimonialClaimRead      Permission = "DAMAGES_AND_PEOPLE_PATRIMONIAL_CLAIM_READ"

	PermissionDamagesAndPeopleResponsibilityRead           Permission = "DAMAGES_AND_PEOPLE_RESPONSIBILITY_READ"
	PermissionDamagesAndPeopleResponsibilityPolicyInfoRead Permission = "DAMAGES_AND_PEOPLE_RESPONSIBILITY_POLICYINFO_READ"
	PermissionDamagesAndPeopleResponsibilityPremiumRead    Permission = "DAMAGES_AND_PEOPLE_RESPONSIBILITY_PREMIUM_READ"
	PermissionDamagesAndPeopleResponsibilityClaimRead      Permission = "DAMAGES_AND_PEOPLE_RESPONSIBILITY_CLAIM_READ"

	PermissionDamagesAndPeopleTransportRead           Permission = "DAMAGES_AND_PEOPLE_TRANSPORT_READ"
	PermissionDamagesAndPeopleTransportPolicyInfoRead Permission = "DAMAGES_AND_PEOPLE_TRANSPORT_POLICYINFO_READ"
	PermissionDamagesAndPeopleTransportPremiumRead    Permission = "DAMAGES_AND_PEOPLE_TRANSPORT_PREMIUM_READ"
	PermissionDamagesAndPeopleTransportClaimRead      Permission = "DAMAGES_AND_PEOPLE_TRANSPORT_CLAIM_READ"

	PermissionDamagesAndPeopleFinancialRisksRead           Permission = "DAMAGES_AND_PEOPLE_FINANCIAL_RISKS_READ"
	PermissionDamagesAndPeopleFinancialRisksPolicyInfoRead Permission = "DAMAGES_AND_PEOPLE_FINANCIAL_RISKS_POLICYINFO_READ"
	PermissionDamagesAndPeopleFinancialRisksPremiumRead    Permission = "DAMAGES_AND_PEOPLE_FINANCIAL_RISKS_PREMIUM_READ"
	PermissionDamagesAndPeopleFinancialRisksClaimRead      Permission = "DAMAGES_AND_PEOPLE_FINANCIAL_RISKS_CLAIM_READ"

	PermissionDamagesAndPeopleRuralRead           Permission = "DAMAGES_AND_PEOPLE_RURAL_READ"
	PermissionDamagesAndPeopleRuralPolicyInfoRead Permission = "DAMAGES_AND_PEOPLE_RURAL_POLICYINFO_READ"
	PermissionDamagesAndPeopleRuralPremiumRead    Permission = "DAMAGES_AND_PEOPLE_RURAL_PREMIUM_READ"
	PermissionDamagesAndPeopleRuralClaimRead      Permission = "DAMAGES_AND_PEOPLE_RURAL_CLAIM_READ"

	PermissionDamagesAndPeopleAutoRead           Permission = "DAMAGES_AND_PEOPLE_AUTO_READ"
	PermissionDamagesAndPeopleAutoPolicyInfoRead Permission = "DAMAGES_AND_PEOPLE_AUTO_POLICYINFO_READ"
	PermissionDamagesAndPeopleAutoPremiumRead    Permission = "DAMAGES_AND_PEOPLE_AUTO_PREMIUM_READ"
	PermissionDamagesAndPeopleAutoClaimRead      Permission = "DAMAGES_AND_PEOPLE_AUTO_CLAIM_READ"

	PermissionDamagesAndPeopleHousingRead           Permission = "DAMAGES_AND_PEOPLE_HOUSING_READ"
	PermissionDamagesAndPeopleHousingPolicyInfoRead Permission = "DAMAGES_AND_PEOPLE_HOUSING_POLICYINFO_READ"
	PermissionDamagesAndPeopleHousingPremiumRead    Permission = "DAMAGES_AND_PEOPLE_HOUSING_PREMIUM_READ"
	PermissionDamagesAndPeopleHousingClaimRead      Permission = "DAMAGES_AND_PEOPLE_HOUSING_CLAIM_READ"

	PermissionDamagesAndPeopleAcceptanceAndBranchesAbroadRead           Permission = "DAMAGES_AND_PEOPLE_ACCEPTANCE_AND_BRANCHES_ABROAD_READ"
	PermissionDamagesAndPeopleAcceptanceAndBranchesAbroadPolicyInfoRead Permission = "DAMAGES_AND_PEOPLE_ACCEPTANCE_AND_BRANCHES_ABROAD_POLICYINFO_READ"
	PermissionDamagesAndPeopleAcceptanceAndBranchesAbroadPremiumRead    Permission = "DAMAGES_AND_PEOPLE_ACCEPTANCE_AND_BRANCHES_ABROAD_PREMIUM_READ"
	PermissionDamagesAndPeopleAcceptanceAndBranchesAbroadClaimRead      Permission = "DAMAGES_AND_PEOPLE_ACCEPTANCE_AND_BRANCHES_ABROAD_CLAIM_READ"

	PermissionDamagesAndPeoplePersonRead           Permission = "DAMAGES_AND_PEOPLE_PERSON_READ"
	PermissionDamagesAndPeoplePersonPolicyInfoRead Permission = "DAMAGES_AND_PEOPLE_PERSON_POLICYINFO_READ"
	PermissionDamagesAndPeoplePersonPremiumRead    Permission = "DAMAGES_AND_PEOPLE_PERSON_PREMIUM_READ"
	PermissionDamagesAndPeoplePersonClaimRead      Permission = "DAMAGES_AND_PEOPLE_PERSON_CLAIM_READ"

	PermissionClaimNotificationRequestDamageCreate Permission = "CLAIM_NOTIFICATION_REQUEST_DAMAGE_CREATE"

	PermissionClaimNotificationRequestPersonCreate Permission = "CLAIM_NOTIFICATION_REQUEST_PERSON_CREATE"

	PermissionEndorsementRequestCreate Permission = "ENDORSEMENT_REQUEST_CREATE"

	PermissionQuotePatrimonialLeadCreate Permission = "QUOTE_PATRIMONIAL_LEAD_CREATE"
	PermissionQuotePatrimonialLeadUpdate Permission = "QUOTE_PATRIMONIAL_LEAD_UPDATE"

	PermissionQuotePatrimonialHomeRead   Permission = "QUOTE_PATRIMONIAL_HOME_READ"
	PermissionQuotePatrimonialHomeCreate Permission = "QUOTE_PATRIMONIAL_HOME_CREATE"
	PermissionQuotePatrimonialHomeUpdate Permission = "QUOTE_PATRIMONIAL_HOME_UPDATE"

	PermissionQuotePatrimonialCondominiumRead   Permission = "QUOTE_PATRIMONIAL_CONDOMINIUM_READ"
	PermissionQuotePatrimonialCondominiumCreate Permission = "QUOTE_PATRIMONIAL_CONDOMINIUM_CREATE"
	PermissionQuotePatrimonialCondominiumUpdate Permission = "QUOTE_PATRIMONIAL_CONDOMINIUM_UPDATE"

	PermissionQuotePatrimonialBusinessRead   Permission = "QUOTE_PATRIMONIAL_BUSINESS_READ"
	PermissionQuotePatrimonialBusinessCreate Permission = "QUOTE_PATRIMONIAL_BUSINESS_CREATE"
	PermissionQuotePatrimonialBusinessUpdate Permission = "QUOTE_PATRIMONIAL_BUSINESS_UPDATE"

	PermissionQuotePatrimonialDiverseRisksRead   Permission = "QUOTE_PATRIMONIAL_DIVERSE_RISKS_READ"
	PermissionQuotePatrimonialDiverseRisksCreate Permission = "QUOTE_PATRIMONIAL_DIVERSE_RISKS_CREATE"
	PermissionQuotePatrimonialDiverseRisksUpdate Permission = "QUOTE_PATRIMONIAL_DIVERSE_RISKS_UPDATE"

	PermissionQuoteAcceptanceAndBranchesAbroadLeadCreate Permission = "QUOTE_ACCEPTANCE_AND_BRANCHES_ABROAD_LEAD_CREATE"
	PermissionQuoteAcceptanceAndBranchesAbroadLeadUpdate Permission = "QUOTE_ACCEPTANCE_AND_BRANCHES_ABROAD_LEAD_UPDATE"

	PermissionQuoteAutoLeadCreate Permission = "QUOTE_AUTO_LEAD_CREATE"
	PermissionQuoteAutoLeadUpdate Permission = "QUOTE_AUTO_LEAD_UPDATE"

	PermissionQuoteAutoRead   Permission = "QUOTE_AUTO_READ"
	PermissionQuoteAutoCreate Permission = "QUOTE_AUTO_CREATE"
	PermissionQuoteAutoUpdate Permission = "QUOTE_AUTO_UPDATE"

	PermissionQuoteFinancialRiskLeadCreate Permission = "QUOTE_FINANCIAL_RISK_LEAD_CREATE"
	PermissionQuoteFinancialRiskLeadUpdate Permission = "QUOTE_FINANCIAL_RISK_LEAD_UPDATE"

	PermissionQuoteHousingLeadCreate Permission = "QUOTE_HOUSING_LEAD_CREATE"
	PermissionQuoteHousingLeadUpdate Permission = "QUOTE_HOUSING_LEAD_UPDATE"

	PermissionQuoteResponsibilityLeadCreate Permission = "QUOTE_RESPONSIBILITY_LEAD_CREATE"
	PermissionQuoteResponsibilityLeadUpdate Permission = "QUOTE_RESPONSIBILITY_LEAD_UPDATE"

	PermissionQuoteRuralLeadCreate Permission = "QUOTE_RURAL_LEAD_CREATE"
	PermissionQuoteRuralLeadUpdate Permission = "QUOTE_RURAL_LEAD_UPDATE"

	PermissionQuoteTransportLeadCreate Permission = "QUOTE_TRANSPORT_LEAD_CREATE"
	PermissionQuoteTransportLeadUpdate Permission = "QUOTE_TRANSPORT_LEAD_UPDATE"

	PermissionQuotePersonLeadCreate Permission = "QUOTE_PERSON_LEAD_CREATE"
	PermissionQuotePersonLeadUpdate Permission = "QUOTE_PERSON_LEAD_UPDATE"

	PermissionQuotePersonLifeRead   Permission = "QUOTE_PERSON_LIFE_READ"
	PermissionQuotePersonLifeCreate Permission = "QUOTE_PERSON_LIFE_CREATE"
	PermissionQuotePersonLifeUpdate Permission = "QUOTE_PERSON_LIFE_UPDATE"

	PermissionQuoteTravelRead   Permission = "QUOTE_TRAVEL_READ"
	PermissionQuoteTravelCreate Permission = "QUOTE_TRAVEL_CREATE"
	PermissionQuoteTravelUpdate Permission = "QUOTE_TRAVEL_UPDATE"

	PermissionQuoteCapitalizationTitleLeadCreate Permission = "QUOTE_CAPITALIZATION_TITLE_LEAD_CREATE"
	PermissionQuoteCapitalizationTitleLeadUpdate Permission = "QUOTE_CAPITALIZATION_TITLE_LEAD_UPDATE"

	PermissionQuoteCapitalizationTitleRead   Permission = "QUOTE_CAPITALIZATION_TITLE_READ"
	PermissionQuoteCapitalizationTitleCreate Permission = "QUOTE_CAPITALIZATION_TITLE_CREATE"
	PermissionQuoteCapitalizationTitleUpdate Permission = "QUOTE_CAPITALIZATION_TITLE_UPDATE"

	PermissionQuoteCapitalizationTitleRaffleCreate Permission = "QUOTE_CAPITALIZATION_TITLE_RAFFLE_CREATE"

	PermissionContractPensionPlanLeadCreate = "CONTRACT_PENSION_PLAN_LEAD_CREATE"
	PermissionContractPensionPlanLeadUpdate = "CONTRACT_PENSION_PLAN_LEAD_UPDATE"

	PermissionContractPensionPlanLeadPortabilityCreate = "CONTRACT_PENSION_PLAN_LEAD_PORTABILITY_CREATE"
	PermissionContractPensionPlanLeadPortabilityUpdate = "CONTRACT_PENSION_PLAN_LEAD_PORTABILITY_UPDATE"

	PermissionContractLifePensionPlanLeadCreate = "CONTRACT_LIFE_PENSION_PLAN_LEAD_CREATE"
	PermissionContractLifePensionPlanLeadUpdate = "CONTRACT_LIFE_PENSION_PLAN_LEAD_UPDATE"

	PermissionContractLifePensionPlanLeadPortabilityCreate = "CONTRACT_LIFE_PENSION_PLAN_LEAD_PORTABILITY_CREATE"
	PermissionContractLifePensionPlanLeadPortabilityUpdate = "CONTRACT_LIFE_PENSION_PLAN_LEAD_PORTABILITY_UPDATE"

	PermissionWithdrawalCreate Permission = "PENSION_WITHDRAWAL_CREATE"

	PermissionCapitalizationTitleWithdrawalCreate Permission = "CAPITALIZATION_TITLE_WITHDRAWAL_CREATE"
)

var permissions = append(permissionsPhase2, permissionsPhase3...)

var permissionsPhase2 = []Permission{
	PermissionResourcesRead,
	PermissionCustomersPersonalIdentificationRead,
	PermissionCustomersPersonalQualificationRead,
	PermissionCustomersPersonalAdditionalInfoRead,
	PermissionCustomersBusinessIdentificationRead,
	PermissionCustomersBusinessQualificationRead,
	PermissionCustomersBusinessAdditionalInfoRead,
	PermissionCapitalizationTitleRead,
	PermissionCapitalizationTitlePlanInfoRead,
	PermissionCapitalizationTitleEventsRead,
	PermissionCapitalizationTitleSettlementsRead,
	PermissionPensionPlanRead,
	PermissionPensionPlanContractInfoRead,
	PermissionPensionPlanMovementsRead,
	PermissionPensionPlanPortabilitiesRead,
	PermissionPensionPlanWithdrawalsRead,
	PermissionPensionPlanClaim,
	PermissionLifePensionRead,
	PermissionLifePensionContractInfoRead,
	PermissionLifePensionMovementsRead,
	PermissionLifePensionPortabilitiesRead,
	PermissionLifePensionWithdrawalsRead,
	PermissionLifePensionClaim,
	PermissionFinancialAssistanceRead,
	PermissionFinancialAssistanceContractInfoRead,
	PermissionFinancialAssistanceMovementsRead,
	PermissionDamagesAndPeoplePatrimonialRead,
	PermissionDamagesAndPeoplePatrimonialPolicyInfoRead,
	PermissionDamagesAndPeoplePatrimonialPremiumRead,
	PermissionDamagesAndPeoplePatrimonialClaimRead,
	PermissionDamagesAndPeopleResponsibilityRead,
	PermissionDamagesAndPeopleResponsibilityPolicyInfoRead,
	PermissionDamagesAndPeopleResponsibilityPremiumRead,
	PermissionDamagesAndPeopleResponsibilityClaimRead,
	PermissionDamagesAndPeopleTransportRead,
	PermissionDamagesAndPeopleTransportPolicyInfoRead,
	PermissionDamagesAndPeopleTransportPremiumRead,
	PermissionDamagesAndPeopleTransportClaimRead,
	PermissionDamagesAndPeopleFinancialRisksRead,
	PermissionDamagesAndPeopleFinancialRisksPolicyInfoRead,
	PermissionDamagesAndPeopleFinancialRisksPremiumRead,
	PermissionDamagesAndPeopleFinancialRisksClaimRead,
	PermissionDamagesAndPeopleRuralRead,
	PermissionDamagesAndPeopleRuralPolicyInfoRead,
	PermissionDamagesAndPeopleRuralPremiumRead,
	PermissionDamagesAndPeopleRuralClaimRead,
	PermissionDamagesAndPeopleAutoRead,
	PermissionDamagesAndPeopleAutoPolicyInfoRead,
	PermissionDamagesAndPeopleAutoPremiumRead,
	PermissionDamagesAndPeopleAutoClaimRead,
	PermissionDamagesAndPeopleHousingRead,
	PermissionDamagesAndPeopleHousingPolicyInfoRead,
	PermissionDamagesAndPeopleHousingPremiumRead,
	PermissionDamagesAndPeopleHousingClaimRead,
	PermissionDamagesAndPeopleAcceptanceAndBranchesAbroadRead,
	PermissionDamagesAndPeopleAcceptanceAndBranchesAbroadPolicyInfoRead,
	PermissionDamagesAndPeopleAcceptanceAndBranchesAbroadPremiumRead,
	PermissionDamagesAndPeopleAcceptanceAndBranchesAbroadClaimRead,
	PermissionDamagesAndPeoplePersonRead,
	PermissionDamagesAndPeoplePersonPolicyInfoRead,
	PermissionDamagesAndPeoplePersonPremiumRead,
	PermissionDamagesAndPeoplePersonClaimRead,
}

var permissionsPhase3 = []Permission{
	PermissionClaimNotificationRequestDamageCreate,
	PermissionClaimNotificationRequestPersonCreate,
	PermissionEndorsementRequestCreate,
	PermissionQuotePatrimonialLeadCreate,
	PermissionQuotePatrimonialLeadUpdate,
	PermissionQuotePatrimonialHomeRead,
	PermissionQuotePatrimonialHomeCreate,
	PermissionQuotePatrimonialHomeUpdate,
	PermissionQuotePatrimonialCondominiumRead,
	PermissionQuotePatrimonialCondominiumCreate,
	PermissionQuotePatrimonialCondominiumUpdate,
	PermissionQuotePatrimonialBusinessRead,
	PermissionQuotePatrimonialBusinessCreate,
	PermissionQuotePatrimonialBusinessUpdate,
	PermissionQuotePatrimonialDiverseRisksRead,
	PermissionQuotePatrimonialDiverseRisksCreate,
	PermissionQuotePatrimonialDiverseRisksUpdate,
	PermissionQuoteAcceptanceAndBranchesAbroadLeadCreate,
	PermissionQuoteAcceptanceAndBranchesAbroadLeadUpdate,
	PermissionQuoteAutoLeadCreate,
	PermissionQuoteAutoLeadUpdate,
	PermissionQuoteAutoRead,
	PermissionQuoteAutoCreate,
	PermissionQuoteAutoUpdate,
	PermissionQuoteFinancialRiskLeadCreate,
	PermissionQuoteFinancialRiskLeadUpdate,
	PermissionQuoteHousingLeadCreate,
	PermissionQuoteHousingLeadUpdate,
	PermissionQuoteResponsibilityLeadCreate,
	PermissionQuoteResponsibilityLeadUpdate,
	PermissionQuoteRuralLeadCreate,
	PermissionQuoteRuralLeadUpdate,
	PermissionQuoteTransportLeadCreate,
	PermissionQuoteTransportLeadUpdate,
	PermissionQuotePersonLeadCreate,
	PermissionQuotePersonLeadUpdate,
	PermissionQuotePersonLifeRead,
	PermissionQuotePersonLifeCreate,
	PermissionQuotePersonLifeUpdate,
	PermissionQuoteTravelRead,
	PermissionQuoteTravelCreate,
	PermissionQuoteTravelUpdate,
	PermissionQuoteCapitalizationTitleLeadCreate,
	PermissionQuoteCapitalizationTitleLeadUpdate,
	PermissionQuoteCapitalizationTitleRead,
	PermissionQuoteCapitalizationTitleCreate,
	PermissionQuoteCapitalizationTitleUpdate,
	PermissionQuoteCapitalizationTitleRaffleCreate,
	PermissionContractPensionPlanLeadCreate,
	PermissionContractPensionPlanLeadUpdate,
	PermissionContractPensionPlanLeadPortabilityCreate,
	PermissionContractPensionPlanLeadPortabilityUpdate,
	PermissionContractLifePensionPlanLeadCreate,
	PermissionContractLifePensionPlanLeadUpdate,
	PermissionContractLifePensionPlanLeadPortabilityCreate,
	PermissionContractLifePensionPlanLeadPortabilityUpdate,
	PermissionWithdrawalCreate,
	PermissionCapitalizationTitleWithdrawalCreate,
}

var (
	permissionCategoryPersonalRegistration PermissionCategory = []Permission{
		PermissionResourcesRead,
		PermissionCustomersPersonalIdentificationRead,
		PermissionCustomersPersonalQualificationRead,
		PermissionCustomersPersonalAdditionalInfoRead,
	}
	permissionCategoryBusinessRegistration PermissionCategory = []Permission{
		PermissionResourcesRead,
		PermissionCustomersBusinessIdentificationRead,
		PermissionCustomersBusinessQualificationRead,
		PermissionCustomersBusinessAdditionalInfoRead,
	}
	permissionCategoryCapitalizationTitle PermissionCategory = []Permission{
		PermissionResourcesRead,
		PermissionCapitalizationTitleRead,
		PermissionCapitalizationTitlePlanInfoRead,
		PermissionCapitalizationTitleEventsRead,
		PermissionCapitalizationTitleSettlementsRead,
	}
	permissionCategoryPensionPlan PermissionCategory = []Permission{
		PermissionResourcesRead,
		PermissionPensionPlanRead,
		PermissionPensionPlanContractInfoRead,
		PermissionPensionPlanMovementsRead,
		PermissionPensionPlanPortabilitiesRead,
		PermissionPensionPlanWithdrawalsRead,
		PermissionPensionPlanClaim,
	}
	permissionCategoryLifePensionPlan PermissionCategory = []Permission{
		PermissionResourcesRead,
		PermissionLifePensionRead,
		PermissionLifePensionContractInfoRead,
		PermissionLifePensionMovementsRead,
		PermissionLifePensionPortabilitiesRead,
		PermissionLifePensionWithdrawalsRead,
		PermissionLifePensionClaim,
	}
	permissionCategoryFinancialAssistence PermissionCategory = []Permission{
		PermissionFinancialAssistanceRead,
		PermissionFinancialAssistanceContractInfoRead,
		PermissionFinancialAssistanceMovementsRead,
	}
	permissionCategoryDamagesAndPeoplePatrimonial PermissionCategory = []Permission{
		PermissionResourcesRead,
		PermissionDamagesAndPeoplePatrimonialRead,
		PermissionDamagesAndPeoplePatrimonialPolicyInfoRead,
		PermissionDamagesAndPeoplePatrimonialPremiumRead,
		PermissionDamagesAndPeoplePatrimonialClaimRead,
	}
	permissionCategoryDamagesAndPeopleResponsibility PermissionCategory = []Permission{
		PermissionResourcesRead,
		PermissionDamagesAndPeopleResponsibilityRead,
		PermissionDamagesAndPeopleResponsibilityPolicyInfoRead,
		PermissionDamagesAndPeopleResponsibilityPremiumRead,
		PermissionDamagesAndPeopleResponsibilityClaimRead,
	}
	permissionCategoryDamagesAndPeopleTransport PermissionCategory = []Permission{
		PermissionResourcesRead,
		PermissionDamagesAndPeopleTransportRead,
		PermissionDamagesAndPeopleTransportPolicyInfoRead,
		PermissionDamagesAndPeopleTransportPremiumRead,
		PermissionDamagesAndPeopleTransportClaimRead,
	}
	permissionCategoryDamagesAndPeopleFinancialRisks PermissionCategory = []Permission{
		PermissionResourcesRead,
		PermissionDamagesAndPeopleFinancialRisksRead,
		PermissionDamagesAndPeopleFinancialRisksPolicyInfoRead,
		PermissionDamagesAndPeopleFinancialRisksPremiumRead,
		PermissionDamagesAndPeopleFinancialRisksClaimRead,
	}
	permissionCategoryDamagesAndPeopleRural PermissionCategory = []Permission{
		PermissionResourcesRead,
		PermissionDamagesAndPeopleRuralRead,
		PermissionDamagesAndPeopleRuralPolicyInfoRead,
		PermissionDamagesAndPeopleRuralPremiumRead,
		PermissionDamagesAndPeopleRuralClaimRead,
	}
	permissionCategoryDamagesAndPeopleAuto PermissionCategory = []Permission{
		PermissionResourcesRead,
		PermissionDamagesAndPeopleAutoRead,
		PermissionDamagesAndPeopleAutoPolicyInfoRead,
		PermissionDamagesAndPeopleAutoPremiumRead,
		PermissionDamagesAndPeopleAutoClaimRead,
	}
	permissionCategoryDamagesAndPeopleHousing PermissionCategory = []Permission{
		PermissionResourcesRead,
		PermissionDamagesAndPeopleHousingRead,
		PermissionDamagesAndPeopleHousingPolicyInfoRead,
		PermissionDamagesAndPeopleHousingPremiumRead,
		PermissionDamagesAndPeopleHousingClaimRead,
	}
	permissionCategoryDamagesAndPeopleAcceptanceAndBranchesAbroad PermissionCategory = []Permission{
		PermissionResourcesRead,
		PermissionDamagesAndPeopleAcceptanceAndBranchesAbroadRead,
		PermissionDamagesAndPeopleAcceptanceAndBranchesAbroadPolicyInfoRead,
		PermissionDamagesAndPeopleAcceptanceAndBranchesAbroadPremiumRead,
		PermissionDamagesAndPeopleAcceptanceAndBranchesAbroadClaimRead,
	}
	permissionCategoryDamagesAndPeoplePerson PermissionCategory = []Permission{
		PermissionResourcesRead,
		PermissionDamagesAndPeoplePersonRead,
		PermissionDamagesAndPeoplePersonPolicyInfoRead,
		PermissionDamagesAndPeoplePersonPremiumRead,
		PermissionDamagesAndPeoplePersonClaimRead,
	}
	permissionCategoryClaimNotificationRequestDamage PermissionCategory = []Permission{
		PermissionClaimNotificationRequestDamageCreate,
	}
	permissionCategoryClaimNotificationRequestPerson PermissionCategory = []Permission{
		PermissionClaimNotificationRequestPersonCreate,
	}
	permissionCategoryEndorsementRequest PermissionCategory = []Permission{
		PermissionEndorsementRequestCreate,
	}
	permissionCategoryQuotePatrimonialLead PermissionCategory = []Permission{
		PermissionQuotePatrimonialLeadCreate,
		PermissionQuotePatrimonialLeadUpdate,
	}
	permissionCategoryQuotePatrimonialHome PermissionCategory = []Permission{
		PermissionQuotePatrimonialHomeRead,
		PermissionQuotePatrimonialHomeCreate,
		PermissionQuotePatrimonialHomeUpdate,
	}
	permissionCategoryQuotePatrimonialCondominium PermissionCategory = []Permission{
		PermissionQuotePatrimonialCondominiumRead,
		PermissionQuotePatrimonialCondominiumCreate,
		PermissionQuotePatrimonialCondominiumUpdate,
	}
	permissionCategoryQuotePatrimonialBusiness PermissionCategory = []Permission{
		PermissionQuotePatrimonialBusinessRead,
		PermissionQuotePatrimonialBusinessCreate,
		PermissionQuotePatrimonialBusinessUpdate,
	}
	permissionCategoryQuotePatrimonialDiverseRisks PermissionCategory = []Permission{
		PermissionQuotePatrimonialDiverseRisksRead,
		PermissionQuotePatrimonialDiverseRisksCreate,
		PermissionQuotePatrimonialDiverseRisksUpdate,
	}
	permissionCategoryQuoteAcceptanceAndBranchesAbroadLead PermissionCategory = []Permission{
		PermissionQuoteAcceptanceAndBranchesAbroadLeadCreate,
		PermissionQuoteAcceptanceAndBranchesAbroadLeadUpdate,
	}
	permissionCategoryQuoteAutoLead PermissionCategory = []Permission{
		PermissionQuoteAutoCreate,
		PermissionQuoteAutoUpdate,
	}
	permissionCategoryQuoteAuto PermissionCategory = []Permission{
		PermissionQuoteAutoRead,
		PermissionQuoteAutoCreate,
		PermissionQuoteAutoUpdate,
	}
	permissionCategoryQuoteFinancialRiskLead PermissionCategory = []Permission{
		PermissionQuoteFinancialRiskLeadCreate,
		PermissionQuoteFinancialRiskLeadUpdate,
	}
	permissionCategoryQuoteHousingLead PermissionCategory = []Permission{
		PermissionQuoteHousingLeadCreate,
		PermissionQuoteHousingLeadUpdate,
	}
	permissionCategoryQuoteResponsibilityLead PermissionCategory = []Permission{
		PermissionQuoteResponsibilityLeadCreate,
		PermissionQuoteResponsibilityLeadUpdate,
	}
	permissionCategoryQuoteRuralLead PermissionCategory = []Permission{
		PermissionQuoteRuralLeadCreate,
		PermissionQuoteRuralLeadUpdate,
	}
	permissionCategoryQuoteTransportLead PermissionCategory = []Permission{
		PermissionQuoteTransportLeadCreate,
		PermissionQuoteTransportLeadUpdate,
	}
	permissionCategoryQuotePersonLead PermissionCategory = []Permission{
		PermissionQuotePersonLeadCreate,
		PermissionQuotePersonLeadUpdate,
	}
	permissionCategoryQuotePersonLifeAuto PermissionCategory = []Permission{
		PermissionQuotePersonLifeRead,
		PermissionQuotePersonLifeCreate,
		PermissionQuotePersonLifeUpdate,
	}
	permissionCategoryQuoteTravelAuto PermissionCategory = []Permission{
		PermissionQuoteTravelRead,
		PermissionQuoteTravelCreate,
		PermissionQuoteTravelUpdate,
	}
	permissionCategoryQuoteCapitalizationTitleLead PermissionCategory = []Permission{
		PermissionQuoteCapitalizationTitleLeadCreate,
		PermissionQuoteCapitalizationTitleLeadUpdate,
	}
	permissionCategoryQuoteCapitalizationTitle PermissionCategory = []Permission{
		PermissionQuoteCapitalizationTitleRead,
		PermissionQuoteCapitalizationTitleCreate,
		PermissionQuoteCapitalizationTitleUpdate,
	}
	permissionCategoryQuoteCapitalizationTitleRaffle PermissionCategory = []Permission{
		PermissionQuoteCapitalizationTitleRaffleCreate,
	}
	permissionCategoryContractPensionPlanLead PermissionCategory = []Permission{
		PermissionContractPensionPlanLeadCreate,
		PermissionContractPensionPlanLeadUpdate,
	}
	permissionCategoryContractPensionPlanLeadPortability PermissionCategory = []Permission{
		PermissionContractPensionPlanLeadPortabilityCreate,
		PermissionContractPensionPlanLeadPortabilityUpdate,
	}
	permissionCategoryContractLifePensionPlanLead PermissionCategory = []Permission{
		PermissionContractLifePensionPlanLeadCreate,
		PermissionContractLifePensionPlanLeadUpdate,
	}
	permissionCategoryContractLifePensionPlanLeadPortability PermissionCategory = []Permission{
		PermissionContractLifePensionPlanLeadPortabilityCreate,
		PermissionContractLifePensionPlanLeadPortabilityUpdate,
	}
	permissionCategoryWithdrawal PermissionCategory = []Permission{
		PermissionWithdrawalCreate,
	}
	permissionCategoryCapitalizationTitleWithdrawalWithdrawal PermissionCategory = []Permission{
		PermissionCapitalizationTitleWithdrawalCreate,
	}
)

var permissionCategories = []PermissionCategory{
	permissionCategoryPersonalRegistration,
	permissionCategoryBusinessRegistration,
	permissionCategoryCapitalizationTitle,
	permissionCategoryPensionPlan,
	permissionCategoryLifePensionPlan,
	permissionCategoryFinancialAssistence,
	permissionCategoryDamagesAndPeoplePatrimonial,
	permissionCategoryDamagesAndPeopleResponsibility,
	permissionCategoryDamagesAndPeopleTransport,
	permissionCategoryDamagesAndPeopleFinancialRisks,
	permissionCategoryDamagesAndPeopleRural,
	permissionCategoryDamagesAndPeopleAuto,
	permissionCategoryDamagesAndPeopleHousing,
	permissionCategoryDamagesAndPeopleAcceptanceAndBranchesAbroad,
	permissionCategoryDamagesAndPeoplePerson,
	permissionCategoryClaimNotificationRequestDamage,
	permissionCategoryClaimNotificationRequestPerson,
	permissionCategoryEndorsementRequest,
	permissionCategoryQuotePatrimonialLead,
	permissionCategoryQuotePatrimonialHome,
	permissionCategoryQuotePatrimonialCondominium,
	permissionCategoryQuotePatrimonialBusiness,
	permissionCategoryQuotePatrimonialDiverseRisks,
	permissionCategoryQuoteAcceptanceAndBranchesAbroadLead,
	permissionCategoryQuoteAutoLead,
	permissionCategoryQuoteAuto,
	permissionCategoryQuoteFinancialRiskLead,
	permissionCategoryQuoteHousingLead,
	permissionCategoryQuoteResponsibilityLead,
	permissionCategoryQuoteRuralLead,
	permissionCategoryQuoteTransportLead,
	permissionCategoryQuotePersonLead,
	permissionCategoryQuotePersonLifeAuto,
	permissionCategoryQuoteTravelAuto,
	permissionCategoryQuoteCapitalizationTitleLead,
	permissionCategoryQuoteCapitalizationTitle,
	permissionCategoryQuoteCapitalizationTitleRaffle,
	permissionCategoryContractPensionPlanLead,
	permissionCategoryContractPensionPlanLeadPortability,
	permissionCategoryContractLifePensionPlanLead,
	permissionCategoryContractLifePensionPlanLeadPortability,
	permissionCategoryWithdrawal,
	permissionCategoryCapitalizationTitleWithdrawalWithdrawal,
}

type Permission string

type PermissionCategory []Permission

func (pc PermissionCategory) contains(p Permission) bool {
	return slices.Contains(pc, p)
}
