package capitalizationtitle

import (
	"net/http"

	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/resource"
)

type Service struct {
	storage         *Storage
	resourceService resource.Service
}

func NewService(
	storage *Storage,
	resourceService resource.Service,
) Service {
	return Service{
		storage:         storage,
		resourceService: resourceService,
	}
}

func (s Service) AddPlan(
	sub string,
	plan api.CapitalizationTitlePlanData,
) {
	s.storage.addPlan(sub, plan)
	for _, company := range plan.Brand.Companies {
		for _, product := range company.Products {
			s.resourceService.Add(sub, api.ResourceData{
				ResourceId: product.PlanId,
				Status:     api.ResourceStatusAVAILABLE,
				Type:       api.ResourceTypeCAPITALIZATIONTITLES,
			})
		}
	}
}

func (s Service) plans(
	meta api.RequestMeta,
	page api.Pagination,
) api.GetCapitalizationTitlePlansResponse {
	plans := s.storage.plans(meta.Subject, page)
	return newPlansResponse(meta, plans)
}

func (s Service) AddPlanInfo(
	sub string,
	planID string,
	info api.CapitalizationTitlePlanInfo,
) {
	s.storage.addPlanInfo(sub, planID, info)
}

func (s Service) planInfo(
	meta api.RequestMeta,
	planID string,
) (
	api.GetCapitalizationTitlePlanInfoResponse,
	error,
) {
	info, err := s.storage.planInfo(meta.Subject, planID)
	if err != nil {
		return api.GetCapitalizationTitlePlanInfoResponse{},
			api.NewError("NOT_FOUND", http.StatusNotFound, err.Error())
	}
	return newPlanInfoResponse(meta, info), nil
}

func (s Service) AddPlanEvent(
	sub string,
	planID string,
	event api.CapitalizationTitleEvent,
) {
	s.storage.addPlanEvent(sub, planID, event)
}

func (s Service) planEvents(
	meta api.RequestMeta,
	planID string,
	page api.Pagination,
) (
	api.GetCapitalizationTitleEventsResponse,
	error,
) {
	events, err := s.storage.planEvents(meta.Subject, planID, page)
	if err != nil {
		return api.GetCapitalizationTitleEventsResponse{},
			api.NewError("NOT_FOUND", http.StatusNotFound, err.Error())
	}
	return newPlanEventsResponse(meta, events), nil
}

func (s Service) AddPlanSettlement(
	sub string,
	planID string,
	settlement api.CapitalizationTitleSettlement,
) {
	s.storage.addPlanSettlement(sub, planID, settlement)
}

func (s Service) planSettlements(
	meta api.RequestMeta,
	planID string,
	page api.Pagination,
) (
	api.GetCapitalizationTitleSettlementsResponse,
	error,
) {
	settlements, err := s.storage.planSettlements(meta.Subject, planID, page)
	if err != nil {
		return api.GetCapitalizationTitleSettlementsResponse{},
			api.NewError("NOT_FOUND", http.StatusNotFound, err.Error())
	}
	resp := newPlanSettlementsResponse(meta, settlements)
	return resp, nil
}
