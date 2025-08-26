package services

import (
	"context"
	"fmt"
	"jinovatka/assert"
	"jinovatka/entities"
	"jinovatka/queue"
	"log/slog"
	"time"
)

type CaptureService struct {
	Log         *slog.Logger
	Queue       queue.Queue
	SeedService *SeedService
}

func NewCaptureService(log *slog.Logger, queue queue.Queue, seedService *SeedService) *CaptureService {
	assert.Must(log != nil, "NewCaptureService: log can't be nil")
	assert.Must(queue != nil, "NewCaptureService: queue can't be nil")
	assert.Must(seedService != nil, "NewCaptureService: seedService can't be nil")
	return &CaptureService{
		Log:         log,
		Queue:       queue,
		SeedService: seedService,
	}
}

// Capture all seeds in group. This will create CaptureRequests for all seeds and enqueue them for capturing.
func (service *CaptureService) CaptureGroup(ctx context.Context, group *entities.SeedsGroup) error {
	for _, seed := range group.Seeds {
		err := service.CaptureSeed(ctx, seed)
		if err != nil {
			// TODO: Maybe wrap this error
			return err
		}
		err = service.SeedService.UpdateState(seed.ShadowID, entities.Pending)
		if err != nil {
			return err
		}
	}
	return nil
}

// Capture single seed. This will create CaptureRequest and enqueue it.
func (service *CaptureService) CaptureSeed(ctx context.Context, seed *entities.Seed) error {
	request := entities.NewRequestFromSeed(seed)
	err := service.Queue.Enqueue(ctx, request)
	if err != nil {
		return fmt.Errorf("CaptureService.CaptureSeed Queue.Enqueue returned error: %w", err)
	}
	return nil
}

// WARNING: This function blocks indefinitely and should be run in separate goroutine.
//
// If timeout is zero, this function blocks until CaptureResult can be dequeued.
// If timeout is nonzero and no CaptureResult is available, this function blocks
// until timeout runs out and then returns QueueTimeoutError and nil CaptureResult.
func (service *CaptureService) AwaitResult(ctx context.Context, timeout time.Duration) (*entities.CaptureResult, error) {
	return service.Queue.AwaitResult(ctx, timeout)
}

// Starts a new goroutine that listens for and handles CaptureResults that are being enququed from workers.
func (service *CaptureService) ListenForResults(ctx context.Context) {
	go service.listenForResults(ctx)
}

func (service *CaptureService) listenForResults(ctx context.Context) {
	// While context is not done, try fetchig next result.
	for ctx.Err() == nil {
		const noTimeout = 0
		result, err := service.AwaitResult(ctx, noTimeout)
		if err != nil {
			service.Log.Error("CaptureService.listenForResults failed to AwaitResult", "error", err.Error())
			break
		}
		service.Log.Info("Got result", "shadowID", result.SeedShadowID, "done", result.Done, "errors", result.ErrorMessages)

		// Update state of seed
		var state entities.CaptureState
		if result.Done && len(result.ErrorMessages) == 0 {
			state = entities.DoneSuccess
		} else if result.Done {
			state = entities.DoneFailure
		} else {
			// TODO add another option
			state = entities.NotEnqueued
		}
		err = service.SeedService.UpdateState(result.SeedShadowID, state)
		if err != nil {
			service.Log.Error("CaptureService.listenForResults failed to update SeedState", "error", err.Error())
			break
		}
		// Update seed metadata
		if result.CaptureMetadata != nil {
			err = service.SeedService.UpdateMetadata(result.SeedShadowID, result.CaptureMetadata)
			if err != nil {
				service.Log.Error("CaptureService.listenForResults failed to update seed metadata", "error", err.Error())
			}
		}
	}
	service.Log.Info("CaptureService.listenForResults context is done", "error", ctx.Err().Error())
}
