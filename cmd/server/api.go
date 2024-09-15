package main

import (
	"context"

	"github.com/luikyv/go-open-insurance/internal/api"
	consentv2 "github.com/luikyv/go-open-insurance/internal/consent/v2"
	customersv1 "github.com/luikyv/go-open-insurance/internal/customer/v1"
)

type opinServer struct {
	consentV2Server  consentv2.Server
	customerV1Server customersv1.Server
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
