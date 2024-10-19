package capitalizationtitle

import (
	"fmt"

	"github.com/luikyv/go-open-insurance/internal/api"
)

type Storage struct {
	plansMap           map[string][]api.CapitalizationTitlePlanData
	planInfoMap        map[string]api.CapitalizationTitlePlanInfo
	planEventsMap      map[string][]api.CapitalizationTitleEvent
	planSettlementsMap map[string][]api.CapitalizationTitleSettlement
}

func NewStorage() *Storage {
	return &Storage{
		plansMap:           make(map[string][]api.CapitalizationTitlePlanData),
		planInfoMap:        make(map[string]api.CapitalizationTitlePlanInfo),
		planEventsMap:      make(map[string][]api.CapitalizationTitleEvent),
		planSettlementsMap: make(map[string][]api.CapitalizationTitleSettlement),
	}
}

func (s *Storage) addPlan(
	sub string,
	title api.CapitalizationTitlePlanData,
) {
	s.plansMap[sub] = append(s.plansMap[sub], title)
}

func (s *Storage) plans(
	sub string,
	page api.Pagination,
) api.Page[api.CapitalizationTitlePlanData] {
	return api.Paginate(s.plansMap[sub], page)
}

func (s *Storage) addPlanInfo(
	sub string,
	planID string,
	info api.CapitalizationTitlePlanInfo,
) {
	s.planInfoMap[sub+"_"+planID] = info
}

func (s *Storage) planInfo(
	sub string,
	planID string,
) (
	api.CapitalizationTitlePlanInfo,
	error,
) {
	info, ok := s.planInfoMap[sub+"_"+planID]
	if !ok {
		return api.CapitalizationTitlePlanInfo{}, fmt.Errorf("plan %s not found", planID)
	}

	return info, nil
}

func (s *Storage) addPlanEvent(
	sub string,
	planID string,
	event api.CapitalizationTitleEvent,
) {
	s.planEventsMap[sub+"_"+planID] = append(s.planEventsMap[sub+"_"+planID], event)
}

func (s *Storage) planEvents(
	sub string,
	planID string,
	page api.Pagination,
) (
	api.Page[api.CapitalizationTitleEvent],
	error,
) {
	events, ok := s.planEventsMap[sub+"_"+planID]
	if !ok {
		return api.Page[api.CapitalizationTitleEvent]{}, fmt.Errorf("plan %s not found", planID)
	}

	return api.Paginate(events, page), nil
}

func (s *Storage) addPlanSettlement(
	sub string,
	planID string,
	settlement api.CapitalizationTitleSettlement,
) {
	s.planSettlementsMap[sub+"_"+planID] = append(
		s.planSettlementsMap[sub+"_"+planID],
		settlement,
	)
}

func (s *Storage) planSettlements(
	sub string,
	planID string,
	page api.Pagination,
) (
	api.Page[api.CapitalizationTitleSettlement],
	error,
) {
	settlements, ok := s.planSettlementsMap[sub+"_"+planID]
	if !ok {
		return api.Page[api.CapitalizationTitleSettlement]{},
			fmt.Errorf("plan %s not found", planID)
	}

	return api.Paginate(settlements, page), nil
}
