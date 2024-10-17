package customer

import "github.com/luikyv/go-open-insurance/internal/api"

type Service struct {
	personalIdentificationsMap   *map[string][]api.PersonalIdentificationData
	personalQualificationsMap    *map[string][]api.PersonalQualificationData
	personalComplimentaryInfoMap *map[string][]api.PersonalComplimentaryInfoData
}

func NewService() Service {
	return Service{
		personalIdentificationsMap:   &map[string][]api.PersonalIdentificationData{},
		personalQualificationsMap:    &map[string][]api.PersonalQualificationData{},
		personalComplimentaryInfoMap: &map[string][]api.PersonalComplimentaryInfoData{},
	}
}

func (s *Service) SetPersonalIdentifications(
	sub string,
	identifications ...api.PersonalIdentificationData,
) {
	(*s.personalIdentificationsMap)[sub] = identifications
}

func (s *Service) PersonalIdentifications(sub string) []api.PersonalIdentificationData {
	return (*s.personalIdentificationsMap)[sub]
}

func (s *Service) SetPersonalQualifications(
	sub string,
	qualifications ...api.PersonalQualificationData,
) {
	(*s.personalQualificationsMap)[sub] = qualifications
}

func (s *Service) PersonalQualifications(sub string) []api.PersonalQualificationData {
	return (*s.personalQualificationsMap)[sub]
}

func (s *Service) SetPersonalComplimentaryInfos(
	sub string,
	infos ...api.PersonalComplimentaryInfoData,
) {
	(*s.personalComplimentaryInfoMap)[sub] = infos
}

func (s *Service) PersonalComplimentaryInfos(sub string) []api.PersonalComplimentaryInfoData {
	return (*s.personalComplimentaryInfoMap)[sub]
}
