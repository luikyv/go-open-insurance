package capitalizationtitle

import (
	"net/http"

	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/opinerr"
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

func (s Service) Plans(
	meta api.RequestMeta,
	page api.Pagination,
) api.Page[api.CapitalizationTitlePlanData] {
	return s.storage.plans(meta.Subject, page)
}

func (s Service) AddPlanInfo(
	sub string,
	planID string,
	info api.CapitalizationTitlePlanInfo,
) {
	s.storage.addPlanInfo(sub, planID, info)
}

func (s Service) PlanInfo(
	meta api.RequestMeta,
	planID string,
) (
	api.CapitalizationTitlePlanInfo,
	error,
) {
	info, err := s.storage.planInfo(meta.Subject, planID)
	if err != nil {
		return api.CapitalizationTitlePlanInfo{},
			opinerr.New("NOT_FOUND", http.StatusNotFound, err.Error())
	}
	return info, nil
}

func (s Service) AddPlanEvent(
	sub string,
	planID string,
	event api.CapitalizationTitleEvent,
) {
	s.storage.addPlanEvent(sub, planID, event)
}

func (s Service) PlanEvents(
	meta api.RequestMeta,
	planID string,
	page api.Pagination,
) (
	api.Page[api.CapitalizationTitleEvent],
	error,
) {
	events, err := s.storage.planEvents(meta.Subject, planID, page)
	if err != nil {
		return api.Page[api.CapitalizationTitleEvent]{},
			opinerr.New("NOT_FOUND", http.StatusNotFound, err.Error())
	}
	return events, nil
}

func (s Service) AddPlanSettlement(
	sub string,
	planID string,
	settlement api.CapitalizationTitleSettlement,
) {
	s.storage.addPlanSettlement(sub, planID, settlement)
}

func (s Service) PlanSettlements(
	meta api.RequestMeta,
	planID string,
	page api.Pagination,
) (
	api.Page[api.CapitalizationTitleSettlement],
	error,
) {
	settlements, err := s.storage.planSettlements(meta.Subject, planID, page)
	if err != nil {
		return api.Page[api.CapitalizationTitleSettlement]{},
			opinerr.New("NOT_FOUND", http.StatusNotFound, err.Error())
	}
	return settlements, nil
}
