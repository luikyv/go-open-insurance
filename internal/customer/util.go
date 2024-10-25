package customer

import (
	"github.com/luikyv/go-open-insurance/internal/api"
)

type Service struct {
	storage *Storage
}

func NewService(storage *Storage) Service {
	return Service{
		storage: storage,
	}
}

func (s Service) AddPersonalIdentification(
	sub string,
	identification api.PersonalIdentificationData,
) {
	s.storage.addPersonalIdentification(sub, identification)
}

func (s Service) PersonalIdentifications(
	meta api.RequestMeta,
) []api.PersonalIdentificationData {
	return s.storage.personalIdentifications(meta.Subject)
}

func (s Service) AddPersonalQualification(
	sub string,
	qualification api.PersonalQualificationData,
) {
	s.storage.addPersonalQualification(sub, qualification)
}

func (s Service) PersonalQualifications(
	meta api.RequestMeta,
) []api.PersonalQualificationData {
	return s.storage.personalQualifications(meta.Subject)
}

func (s Service) AddPersonalComplimentaryInfo(
	sub string,
	info api.PersonalComplimentaryInfoData,
) {
	s.storage.addPersonalComplimentaryInfo(sub, info)
}

func (s Service) PersonalComplimentaryInfos(
	meta api.RequestMeta,
) []api.PersonalComplimentaryInfoData {
	return s.storage.personalComplimentaryInfos(meta.Subject)
}

type Storage struct {
	personalIdentificationsMap   map[string][]api.PersonalIdentificationData
	personalQualificationsMap    map[string][]api.PersonalQualificationData
	personalComplimentaryInfoMap map[string][]api.PersonalComplimentaryInfoData
}

func NewStorage() *Storage {
	return &Storage{
		personalIdentificationsMap:   make(map[string][]api.PersonalIdentificationData),
		personalQualificationsMap:    make(map[string][]api.PersonalQualificationData),
		personalComplimentaryInfoMap: make(map[string][]api.PersonalComplimentaryInfoData),
	}
}

func (s *Storage) addPersonalIdentification(
	sub string,
	identification api.PersonalIdentificationData,
) {
	s.personalIdentificationsMap[sub] = append(
		s.personalIdentificationsMap[sub],
		identification,
	)
}

func (s *Storage) personalIdentifications(
	sub string,
) []api.PersonalIdentificationData {
	return s.personalIdentificationsMap[sub]
}

func (s *Storage) addPersonalQualification(
	sub string,
	qualification api.PersonalQualificationData,
) {
	s.personalQualificationsMap[sub] = append(
		s.personalQualificationsMap[sub],
		qualification,
	)
}

func (s *Storage) personalQualifications(sub string) []api.PersonalQualificationData {
	return s.personalQualificationsMap[sub]
}

func (s *Storage) addPersonalComplimentaryInfo(
	sub string,
	info api.PersonalComplimentaryInfoData,
) {
	s.personalComplimentaryInfoMap[sub] = append(
		s.personalComplimentaryInfoMap[sub],
		info,
	)
}

func (s *Storage) personalComplimentaryInfos(sub string) []api.PersonalComplimentaryInfoData {
	return s.personalComplimentaryInfoMap[sub]
}
