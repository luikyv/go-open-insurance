package resource

import (
	"context"
	"net/http"
	"strings"

	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/consent"
)

type Service struct {
	storage        *Storage
	consentService consent.Service
}

func NewService(storage *Storage, consentService consent.Service) Service {
	return Service{
		storage:        storage,
		consentService: consentService,
	}
}

func (s Service) Add(sub string, resource api.ResourceData) {
	s.storage.add(sub, resource)
}

func (s Service) Resource(
	ctx context.Context,
	meta api.RequestMeta,
	id string,
) (
	api.ResourceData,
	error,
) {
	r, err := s.storage.get(meta.Subject, id)
	if err != nil {
		return api.ResourceData{},
			api.NewError("NAO_FOUND", http.StatusNotFound, err.Error())
	}
	return r, nil
}

func (s Service) resources(
	ctx context.Context,
	meta api.RequestMeta,
	page api.Pagination,
) (
	api.GetResourcesResponse,
	error,
) {
	consent, err := s.consentService.Fetch(ctx, meta, meta.ConsentID)
	if err != nil {
		return api.GetResourcesResponse{}, err
	}

	consentedTypes := consentedResourceTypes(consent.Permissions)
	rs := s.storage.resources(meta.Subject, consentedTypes, page)
	return newResourcesResponse(meta, rs), nil
}

func consentedResourceTypes(permissions []api.ConsentPermission) []api.ResourceType {
	// TODO: Should I enumerate this?
	var consentedTypes []api.ResourceType
	for _, typ := range resourceTypes {
		// Some resource types have a trailing 'S' whereas the corresponding
		// permisions don't. E.g., the resource type 'CAPITALIZATION_TITLES' and
		// the permission 'CAPITALIZATION_TITLE_EVENTS_READ'.
		typStr := strings.TrimRight(string(typ), "S")
		for _, p := range permissions {

			if strings.HasPrefix(string(p), typStr) {
				consentedTypes = append(consentedTypes, typ)
				break
			}
		}
	}
	return consentedTypes
}
