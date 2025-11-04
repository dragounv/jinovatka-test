package services

import (
	"jinovatka/assert"
	"jinovatka/storage"
)

func NewIsUniqueService(idListRepository storage.IdListRepository) *IsUniqueService {
	assert.Must(idListRepository != nil, "NewIsUniqueService: idListRepository can't be nil")
	return &IsUniqueService{
		IdListRepository: idListRepository,
	}
}

// Service implementing IsUniqiueChecker
type IsUniqueService struct {
	IdListRepository storage.IdListRepository
}

func (service *IsUniqueService) Add(id string) error {
	return service.IdListRepository.Add(id)
}

func (service *IsUniqueService) IsUnique(id string) (bool, error) {
	alredyExists, err := service.IdListRepository.AlredyExists(id)
	if err != nil {
		return false, err
	}
	return !alredyExists, nil
}
