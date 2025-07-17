package storage

import (
	"jinovatka/assert"
	"jinovatka/entities"
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
}
