package services

import (
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

func NewServices(log *slog.Logger, repository *storage.Repository, queue queue.Queue) *Services {
	assert.Must(log != nil, "NewServices: log can't be nil")
	assert.Must(repository != nil, "NewServices: repository can't be nil")
	seedService := NewSeedService(log, repository.SeedRepository, MaxUrlAdressLength, MaxInputedUrlAddresses)
	exporterService := NewExporterService()
	captureService := NewCaptureService(queue)
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
