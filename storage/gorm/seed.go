package gormStorage

import (
	"database/sql"
	"errors"
	"fmt"
	"jinovatka/assert"
	"jinovatka/entities"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type Seed struct {
	gorm.Model

	// The original seed URL
	URL string `gorm:"index"`

	// If the seed is public, hten it may be visible on the main page and can be searched for.
	Public bool

	// Determines if the seed was alredy harvested or not.
	State string

	// URL of the archived resource
	ArchivalURL sql.NullString

	// Time of harvest. Should be eqivalent to the time used to generate ArchivalUrl.
	// If Null, the seed wasn't harvested yet (it may be waiting in a queue, or error happend during harvest)
	HarvestedAt sql.NullTime

	// Unique randomly generated base32 encoded string with at least 128 bits of randomness.
	// Exact size is unspecified. This allowes the use of rand.Text to generate it.
	//
	// It serves as unique name for harvest, that is should be unpredictable and resilient to brute force guessing.
	// It will be used for generating URLs, that cannot be easily guessed, to preserve privacy of harvest creators.
	ShadowID string `gorm:"unique;index"`

	// Foreing key for SeedsGroup.
	SeedsGroupID uint
}

// Function that creates new seed records that can be inserted into DB.
// It asserts that all necessary values are present to prevent incomplete records.
func NewSeedRecord(seed *entities.Seed) *Seed {
	assert.Must(seed != nil, "NewSeedRecord: seed can't be nil")
	assert.Must(seed.URL != "", "NewSeedRecord: seed.URL can't be empty string")
	assert.Must(seed.ShadowID != "", "NewSeedRecord: seed.ShadowID can't be empty string")
	assert.Must(seed.State == entities.NotEnqueued, "NewSeedRecord: seed.State can't be empty and must be entities.NotEnqueued")
	// We are creating seed that is not yet harvested, so it can't have harvest metadata.
	assert.Must(seed.HarvestedAt.IsZero(), "NewSeedRecord: seed.HarvestedAt must be zero time value")
	assert.Must(seed.ArchivalURL == "", "NewSeedRecord: seed.ArchivalURL must be empty string")

	// Ignore ID, Create, Update adn Delete time. We are creating new record. GORM will fill it in.
	seedRecord := &Seed{
		URL:      seed.URL,
		Public:   seed.Public,
		State:    string(seed.State),
		ShadowID: seed.ShadowID,
	}

	// If we have harvest metadata, put them in the record.
	// We want to allow partially filled records so
	if !seed.HarvestedAt.IsZero() {
		seedRecord.HarvestedAt.Time = seed.HarvestedAt
		seedRecord.HarvestedAt.Valid = true
	}
	if seed.ArchivalURL != "" {
		seedRecord.ArchivalURL.String = seed.ArchivalURL
		seedRecord.ArchivalURL.Valid = true
	}

	return seedRecord
}

func (seed *Seed) ToEntity() *entities.Seed {
	entity := &entities.Seed{
		URL:      seed.URL,
		Public:   seed.Public,
		State:    entities.CaptureState(seed.State),
		ShadowID: seed.ShadowID,
	}
	if seed.ArchivalURL.Valid {
		entity.ArchivalURL = seed.ArchivalURL.String
	} else {
		entity.ArchivalURL = ""
	}
	if seed.HarvestedAt.Valid {
		entity.HarvestedAt = seed.HarvestedAt.Time
	} else {
		entity.HarvestedAt = time.Time{}
	}
	return entity
}

type SeedsGroup struct {
	gorm.Model
	// Has many seeds.
	Seeds    []*Seed
	ShadowID string `gorm:"unique;index"`
}

func NewSeedGroup(seedsGroup *entities.SeedsGroup) *SeedsGroup {
	assert.Must(seedsGroup != nil, "NewSeedGroup: seedsGroup can't be nil")
	assert.Must(seedsGroup.Seeds != nil, "NewSeedGroup: seedsGroup.Seeds can't be nil")
	assert.Must(seedsGroup.ShadowID != "", "NewSeedGroup: seedsGroup.ShadowID can't be empty string")
	seedRecords := make([]*Seed, 0, len(seedsGroup.Seeds))
	for _, seed := range seedsGroup.Seeds {
		seedRecords = append(seedRecords, NewSeedRecord(seed))
	}
	return &SeedsGroup{
		Seeds:    seedRecords,
		ShadowID: seedsGroup.ShadowID,
	}
}

func (group *SeedsGroup) ToEntity() *entities.SeedsGroup {
	seeds := make([]*entities.Seed, 0, len(group.Seeds))
	for _, seed := range group.Seeds {
		seeds = append(seeds, seed.ToEntity())
	}
	return &entities.SeedsGroup{
		Seeds:    seeds,
		ShadowID: group.ShadowID,
	}
}

func NewSeedRepository(log *slog.Logger, db *gorm.DB) *SeedRepository {
	assert.Must(log != nil, "NewSeedRepository: log can't be nil")
	assert.Must(db != nil, "NewSeedRepository: db can't be nil")
	var err error
	err = db.AutoMigrate(Seed{})
	assert.Must(err == nil, "NewSeedRepository: db.AutoMigrate failed for Seed with error: "+assert.AddErrorMessage(err))
	err = db.AutoMigrate(SeedsGroup{})
	assert.Must(err == nil, "NewSeedRepository: db.AutoMigrate failed for SeedsGroup with error: "+assert.AddErrorMessage(err))
	return &SeedRepository{
		Log: log,
		DB:  db,
	}
}

type SeedRepository struct {
	Log *slog.Logger
	DB  *gorm.DB
}

// Create new seed and save it to db.
// func (repository *SeedRepository) Save(seed *entities.Seed) error {
// 	result := repository.DB.Create(NewSeedRecord(seed))
// 	if result.Error != nil {
// 		return fmt.Errorf("SeedRepository.Save: failed to create new record %w", result.Error)
// 	}
// 	return nil
// }

// Save list of SeedRecords to the repository.
func (repository *SeedRepository) Save(seeds []*entities.Seed) error {
	if len(seeds) == 0 {
		// TODO: Decide if error is good here. Maybe assert or even doing nothing could make more sense?
		// It dies on nil anyway and just zero seeds could be allowed?
		return errors.New("SeedsRepository.Save recieved slice with zero length")
	}
	seedRecords := make([]*Seed, 0, len(seeds))
	for _, seed := range seeds {
		seedRecords = append(seedRecords, NewSeedRecord(seed))
	}
	result := repository.DB.Create(seedRecords)
	if result.Error != nil {
		return fmt.Errorf("SeedRepository.Save failed to create new seeds: %w", result.Error)
	}
	return nil
}

func (repository *SeedRepository) SaveGroup(seedsGroup *entities.SeedsGroup) error {
	if seedsGroup == nil {
		return errors.New("SeedRepository.SaveGroup recieved nil seedsGroup")
	}
	groupRecord := NewSeedGroup(seedsGroup)
	result := repository.DB.Create(groupRecord)
	if result.Error != nil {
		return fmt.Errorf("SeedRepository.SaveGroup failed to create new group: %w", result.Error)
	}
	return nil
}

func (repoository *SeedRepository) GetGroup(shadow string) (*entities.SeedsGroup, error) {
	groupRecord := new(SeedsGroup)
	err := repoository.DB.Model(&SeedsGroup{}).Preload("Seeds").First(groupRecord, "shadow_id = ?", shadow).Error
	if err != nil {
		return nil, fmt.Errorf("SeedRepository.GetGroup failed to fetch SeedsGroup from db: %w", err)
	}
	group := groupRecord.ToEntity()
	return group, nil
}

func (repository *SeedRepository) GetSeed(shadow string) (*entities.Seed, error) {
	seedRecord := new(Seed)
	err := repository.DB.First(seedRecord, "shadow_id = ?", shadow).Error
	if err != nil {
		return nil, fmt.Errorf("SeedRepository.GetSeed failed to fetch Seed from db: %w", err)
	}
	seed := seedRecord.ToEntity()
	return seed, nil
}

func (repository *SeedRepository) UpdateState(shadow string, state entities.CaptureState) error {
	err := repository.DB.Model(Seed{}).Where("shadow_id = ?", shadow).Select("state").Updates(Seed{State: string(state)}).Error
	if err != nil {
		return fmt.Errorf("SeedRepository.UpdateStatus failed to update Seed with shadow %s :%w", shadow, err)
	}
	return nil
}
