package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"jinovatka/assert"
	"jinovatka/entities"
	"jinovatka/storage"
	"log/slog"
	"strings"
	"time"
)

var ErrEmptyList = errors.New("list was empty")

func NewSeedService(
	log *slog.Logger,
	repository storage.SeedRepository,
	maxInputListLineLength,
	maxInputListLines int,
) *SeedService {
	assert.Must(log != nil, "NewSeedService: log can't be nil")
	assert.Must(repository != nil, "NewSeedService: repository can't be nil")
	return &SeedService{
		Log:                    log,
		Repository:             repository,
		UrlParser:              new(UrlParserService),
		MaxInputListLineLength: maxInputListLineLength,
		MaxInputListLines:      maxInputListLines,
	}
}

type SeedService struct {
	Log        *slog.Logger
	Repository storage.SeedRepository

	UrlParser *UrlParserService

	// Maximum length of seed URL adress
	MaxInputListLineLength int
	// Maximum number of lines (seeds) that can be parsed in one call
	MaxInputListLines int
}

type SeedState int

const (
	NotHarvested = iota
	Sucess
	Failed
)

type FindSeedsArgs struct {
	URL       string
	StartDate *time.Time
	EndDate   *time.Time
}

func (service *SeedService) FindSeeds(arguments *FindSeedsArgs) ([]*entities.Seed, error) {
	seeds := make([]*entities.Seed, 58)
	return seeds, nil
}

// Save single seed to repository. This function is deprecated.
// TODO: Get rid of this function. SaveList can hadle both signle and multiple seeds.
func (service *SeedService) SaveSeed(seedURL string) error {
	service.Log.Info("recieved seed", slog.String("url", seedURL))
	shadow := rand.Text()

	seed := &entities.Seed{
		URL:      seedURL,
		Public:   true,
		State:    entities.NotEnqueued,
		ShadowID: shadow,
	}

	err := service.Repository.Save([]*entities.Seed{seed})
	if err != nil {
		return fmt.Errorf("SeedService.SaveSeed failed to save seed to repository: %w", err)
	}

	return nil
}

// Takes string consisting of newline delimited list of URL adresses.
// Checks input data size, parses them into slice of strings and delegates to SaveList.
func (service *SeedService) Save(urlsList string, storeGroup bool) (*entities.SeedsGroup, error) {
	if urlsList == "" {
		return nil, errors.New("SeedServie.Save no input data")
	}
	// Check the entire length of the string.
	if len(urlsList) > (service.MaxInputListLineLength * service.MaxInputListLines) {
		return nil, errors.New("SeedService.Save input data is too large")
	}
	lines := strings.Split(urlsList, "\n")
	// Now check just the number of lines.
	if len(lines) > service.MaxInputListLines {
		return nil, errors.New("SeedService.Save too many lines")
	}
	return service.SaveList(lines, storeGroup)
}

// Save list of URL adresses as Seeds. Does format validation but does not check input data size.
// For saving input from untrusted source use SeedService.Save instead.
func (service *SeedService) SaveList(lines []string, storeGroup bool) (*entities.SeedsGroup, error) {
	if len(lines) == 0 {
		return nil, errors.New("SeedService.SaveList no input data")
	}
	seeds := make([]*entities.Seed, 0, len(lines))
	for _, url := range lines {
		url = strings.TrimSpace(url)
		if url == "" { // Skip empty lines
			continue
		}
		url, err := service.UrlParser.ParseAndCleanURL(url, false)
		if err != nil {
			// TODO: Log this in some smart way.
			return nil, fmt.Errorf("SeedService.SaveList failed to parse URL: %w", err)
		}
		shadow := rand.Text()
		seed := &entities.Seed{
			URL:      url.String(),
			Public:   true,
			State:    entities.NotEnqueued,
			ShadowID: shadow,
		}
		seeds = append(seeds, seed)
	}

	if len(seeds) == 0 {
		return nil, ErrEmptyList
	}

	var err error
	group := &entities.SeedsGroup{Seeds: seeds}
	if storeGroup {
		// Only create the shadow if we are gonna store the group.
		group.ShadowID = rand.Text()
		err = service.Repository.SaveGroup(group)
	} else {
		err = service.Repository.Save(seeds)
	}
	if err != nil {
		return group, fmt.Errorf("SeedService.SaveList failed to save seeds to repository: %w", err)
	}

	return group, nil
}

func (service *SeedService) GetGroup(shadow string) (*entities.SeedsGroup, error) {
	return service.Repository.GetGroup(shadow)
}

func (service *SeedService) GetSeed(shadow string) (*entities.Seed, error) {
	return service.Repository.GetSeed(shadow)
}

func (service *SeedService) UpdateState(shadow string, state entities.CaptureState) error {
	if !state.IsCaptureState() {
		return errors.New("SeedService.UpdateState received invalid state argument")
	}
	return service.Repository.UpdateState(shadow, state)
}

func (service *SeedService) UpdateMetadata(shadow string, metadata *entities.CaptureMetadata) error {
	// Check that timestamp is valid. It needs to be 14 (second precision) or 17 (nanosecond precision) chars long.
	timestampLength := len(metadata.Timestamp)
	if timestampLength != 14 && timestampLength != 17 {
		return errors.New("SeedService.UpdateMetadata recieved CaptureMetadata with invalid timestamp length")
	}

	// Create archivalURL
	// TODO: Move this to config
	const waybackURL = "https://wayback.webarchiv.cz/wayback/"
	// Intentinally left unescaped. Escaping the capturedUrl, while it may seem reasonable, will break compatibility with openwayback.
	waybackPagePath := metadata.Timestamp + "/" + metadata.CapturedUrl
	archivalURL := waybackURL + waybackPagePath

	// If timestamp is the long version, then add decimal point before the fractional second part to allow parsing by time.Parse.
	timestamp := metadata.Timestamp
	if timestampLength == 17 {
		timestamp = timestamp[:14] + "." + timestamp[14:]
	}

	service.Log.Info("Timestamp: " + timestamp)

	// Parse timestamp into time.Time. This format covers both short and long variants as specifying fractional seconds part is optinal.
	const format = "20060102150405"
	archivedAt, err := time.Parse(format, timestamp)
	if err != nil {
		return fmt.Errorf("SeedService.UpdateMetadata failed to parse timestamp: %w", err)
	}

	return service.Repository.UpdateMetadata(shadow, archivalURL, archivedAt)
}
