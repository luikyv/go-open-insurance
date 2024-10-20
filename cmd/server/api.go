package main

import (
	"context"

	"github.com/luikyv/go-open-insurance/internal/api"
	capitalizationtitlev1 "github.com/luikyv/go-open-insurance/internal/capitalizationtitle/v1"
	consentv2 "github.com/luikyv/go-open-insurance/internal/consent/v2"
	customersv1 "github.com/luikyv/go-open-insurance/internal/customer/v1"
	endorsementv1 "github.com/luikyv/go-open-insurance/internal/endorsement/v1"
	resourcev2 "github.com/luikyv/go-open-insurance/internal/resource/v2"
)

type opinServer struct {
	consentV2Server             consentv2.Server
	customerV1Server            customersv1.Server
	resouceV2Server             resourcev2.Server
	capitalizationTitleV1Server capitalizationtitlev1.Server
	endorsementV1Server         endorsementv1.Server
}

func (s opinServer) CreateConsentV2(
	ctx context.Context,
	request api.CreateConsentV2RequestObject,
) (
	api.CreateConsentV2ResponseObject,
	error,
) {
	return s.consentV2Server.CreateConsentV2(ctx, request)
}

func (s opinServer) DeleteConsentV2(
	ctx context.Context,
	request api.DeleteConsentV2RequestObject,
) (
	api.DeleteConsentV2ResponseObject,
	error,
) {
	return s.consentV2Server.DeleteConsentV2(ctx, request)
}

func (s opinServer) ConsentV2(
	ctx context.Context,
	request api.ConsentV2RequestObject,
) (
	api.ConsentV2ResponseObject,
	error,
) {
	return s.consentV2Server.ConsentV2(ctx, request)
}

func (s opinServer) PersonalIdentificationsV1(
	ctx context.Context,
	request api.PersonalIdentificationsV1RequestObject,
) (
	api.PersonalIdentificationsV1ResponseObject,
	error,
) {
	return s.customerV1Server.PersonalIdentificationsV1(ctx, request)
}

func (s opinServer) PersonalQualificationsV1(
	ctx context.Context,
	request api.PersonalQualificationsV1RequestObject,
) (
	api.PersonalQualificationsV1ResponseObject,
	error,
) {
	return s.customerV1Server.PersonalQualificationsV1(ctx, request)
}

func (s opinServer) PersonalComplimentaryInfoV1(
	ctx context.Context,
	request api.PersonalComplimentaryInfoV1RequestObject,
) (
	api.PersonalComplimentaryInfoV1ResponseObject,
	error,
) {
	return s.customerV1Server.PersonalComplimentaryInfoV1(ctx, request)
}

func (s opinServer) ResourcesV2(
	ctx context.Context,
	request api.ResourcesV2RequestObject,
) (
	api.ResourcesV2ResponseObject,
	error,
) {
	return s.resouceV2Server.ResourcesV2(ctx, request)
}

func (s opinServer) CapitalizationTitlePlansV1(
	ctx context.Context,
	request api.CapitalizationTitlePlansV1RequestObject,
) (
	api.CapitalizationTitlePlansV1ResponseObject,
	error,
) {
	return s.capitalizationTitleV1Server.CapitalizationTitlePlans(ctx, request)
}

func (s opinServer) CapitalizationTitleEventsV1(
	ctx context.Context,
	request api.CapitalizationTitleEventsV1RequestObject,
) (
	api.CapitalizationTitleEventsV1ResponseObject,
	error,
) {
	return s.capitalizationTitleV1Server.CapitalizationTitleEvents(ctx, request)
}

func (s opinServer) CapitalizationTitlePlanInfoV1(
	ctx context.Context,
	request api.CapitalizationTitlePlanInfoV1RequestObject,
) (
	api.CapitalizationTitlePlanInfoV1ResponseObject,
	error,
) {
	return s.capitalizationTitleV1Server.CapitalizationTitlePlanInfo(ctx, request)
}

func (s opinServer) CapitalizationTitleSettlementsV1(
	ctx context.Context,
	request api.CapitalizationTitleSettlementsV1RequestObject,
) (
	api.CapitalizationTitleSettlementsV1ResponseObject,
	error,
) {
	return s.capitalizationTitleV1Server.CapitalizationTitleSettlements(ctx, request)
}

func (s opinServer) CreateEndorsementV1(
	ctx context.Context,
	request api.CreateEndorsementV1RequestObject,
) (
	api.CreateEndorsementV1ResponseObject,
	error,
) {
	return s.endorsementV1Server.CreateEndorsementV1(ctx, request)
}
