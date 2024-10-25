package resource

import (
	"context"
	"net/http"
	"strings"

	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/consent"
	"github.com/luikyv/go-open-insurance/internal/opinerr"
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

func (s Service) Resources(
	ctx context.Context,
	meta api.RequestMeta,
	page api.Pagination,
) (
	api.Page[api.ResourceData],
	error,
) {
	consent, err := s.consentService.Get(ctx, meta, meta.ConsentID)
	if err != nil {
		return api.Page[api.ResourceData]{}, err
	}

	consentedTypes := consentedResourceTypes(consent.Permissions)
	return s.storage.resources(meta.Subject, consentedTypes, page), nil
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

func (s Service) Get(
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
			opinerr.New("NAO_FOUND", http.StatusNotFound, err.Error())
	}
	return r, nil
}
