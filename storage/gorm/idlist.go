package gormStorage

import (
	"errors"
	"fmt"
	"jinovatka/assert"

	"gorm.io/gorm"
)

type Id struct {
	gorm.Model
	Value string `gorm:"unique;index"`
}

func NewIdListRepository(db *gorm.DB) *IdListRepository {
	assert.Must(db != nil, "NewIdListRepository: db can't be nil")
	err := db.AutoMigrate(Id{})
	assert.Must(err == nil, "NewIdListRepository: db.AutoMigrate failed for model Id with error: "+assert.AddErrorMessage(err))
	return &IdListRepository{
		DB: db,
	}
}

type IdListRepository struct {
	DB *gorm.DB
}

func (repository *IdListRepository) Add(id string) error {
	result := repository.DB.Create(&Id{Value: id})
	if result.Error != nil {
		return fmt.Errorf("IdListRepository.Add could not create database record: %w", result.Error)
	}
	return nil
}

func (repository *IdListRepository) AlredyExists(id string) (bool, error) {
	idRecord := &Id{}
	result := repository.DB.Where("value = ?", id).First(idRecord)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if result.Error != nil {
		return false, fmt.Errorf("IdListRepository.AlredyExists database returned error: %w", result.Error)
	}
	return true, nil
}
