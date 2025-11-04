package storage

import (
	"jinovatka/assert"
	"jinovatka/entities"
	"time"
)

func NewRepository(seed SeedRepository, idList IdListRepository) *Repository {
	assert.Must(seed != nil, "NewRepository: seed repository can't be nil")
	assert.Must(idList != nil, "NewRepository: idList repository can't be nil")
	return &Repository{
		SeedRepository:   seed,
		IdListRepository: idList,
	}
}

type Repository struct {
	SeedRepository   SeedRepository
	IdListRepository IdListRepository
}

type SeedRepository interface {
	Save([]*entities.Seed) error
	SaveGroup(*entities.SeedsGroup) error
	GetGroup(shadow string) (*entities.SeedsGroup, error)
	GetSeed(shadow string) (*entities.Seed, error)
	UpdateState(shadow string, state entities.CaptureState) error
	UpdateMetadata(shadow, archivalURL string, harvestedAt time.Time) error
}

type IdListRepository interface {
	Add(id string) error
	AlredyExists(id string) (bool, error)
}
