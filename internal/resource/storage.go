package resource

import (
	"fmt"
	"slices"

	"github.com/luikyv/go-open-insurance/internal/api"
)

type Storage struct {
	resourcesMap map[string][]api.ResourceData
}

func NewStorage() *Storage {
	return &Storage{
		resourcesMap: make(map[string][]api.ResourceData),
	}
}

func (s *Storage) add(sub string, resource api.ResourceData) {
	s.resourcesMap[sub] = append(s.resourcesMap[sub], resource)
}

func (s *Storage) resources(
	sub string,
	types []api.ResourceType,
	page api.Pagination,
) api.Page[api.ResourceData] {
	var rs []api.ResourceData
	for _, r := range s.resourcesMap[sub] {
		if slices.Contains(types, r.Type) {
			rs = append(rs, r)
		}
	}
	return api.Paginate(rs, page)
}

func (s *Storage) get(sub string, id string) (api.ResourceData, error) {
	for _, r := range s.resourcesMap[sub] {
		if r.ResourceId == id {
			return r, nil
		}
	}

	return api.ResourceData{}, fmt.Errorf("resource %s not found", id)
}
