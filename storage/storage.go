package storage

import (
	"jinovatka/assert"
	"jinovatka/entities"
	"time"
)

func NewRepository(seed SeedRepository) *Repository {
	assert.Must(seed != nil, "NewRepository: seed repository can't be nil")
	return &Repository{
		SeedRepository: seed,
	}
}

type Repository struct {
	SeedRepository SeedRepository
}

type SeedRepository interface {
	Save([]*entities.Seed) error
	SaveGroup(*entities.SeedsGroup) error
	GetGroup(shadow string) (*entities.SeedsGroup, error)
	GetSeed(shadow string) (*entities.Seed, error)
	UpdateState(shadow string, state entities.CaptureState) error
	UpdateMetadata(shadow, archivalURL string, harvestedAt time.Time) error
}
