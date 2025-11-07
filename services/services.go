package services

import (
	"context"
	"jinovatka/assert"
	"jinovatka/queue"
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

type ServiceSettings struct {
	ServerHost          string
	SeedDetailPath      string
	WaybackRedirectPath string
}

func NewServices(ctx context.Context, log *slog.Logger, repository *storage.Repository, queue queue.Queue, settings *ServiceSettings) *Services {
	assert.Must(log != nil, "NewServices: log can't be nil")
	assert.Must(repository != nil, "NewServices: repository can't be nil")
	seedService := NewSeedService(
		ctx,
		log,
		repository.SeedRepository,
		NewIsUniqueService(repository.IdListRepository),
		MaxUrlAdressLength,
		MaxInputedUrlAddresses,
	)
	exporterService := NewExporterService(settings)
	captureService := NewCaptureService(log, queue, seedService)
	return &Services{
		SeedService:     seedService,
		ExporterService: exporterService,
		CaptureService:  captureService,
	}
}

type Services struct {
	SeedService     *SeedService
	ExporterService *ExporterService
	CaptureService  *CaptureService
}
