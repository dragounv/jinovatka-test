package services

import (
	"jinovatka/assert"
	"jinovatka/storage"
	"log/slog"
)

// TODO: Make some better way for dealing with settings/constants
const (
	// This sets the maximum allowed size of URL adresses.
	// 64kB. Some quick reaserch seems to show that larger URLs could cause issues during crawls.
	MaxUrlAdressLength     = 64 << 10
	MaxInputedUrlAddresses = 20
)

func NewServices(log *slog.Logger, repository *storage.Repository) *Services {
	assert.Must(log != nil, "NewServices: log can't be nil")
	assert.Must(repository != nil, "NewServices: repository can't be nil")
	seedService := NewSeedService(log, repository.SeedRepository, MaxUrlAdressLength, MaxInputedUrlAddresses)
	exporterService := NewExporterService()
	return &Services{
		SeedService:     seedService,
		ExporterService: exporterService,
	}
}

type Services struct {
	SeedService     *SeedService
	ExporterService *ExporterService
}
