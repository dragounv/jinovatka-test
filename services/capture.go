package services

import (
	"context"
	"fmt"
	"jinovatka/assert"
	"jinovatka/entities"
	"jinovatka/queue"
)

type CaptureService struct {
	Queue queue.Queue
}

func NewCaptureService(queue queue.Queue) *CaptureService {
	assert.Must(queue != nil, "NewCaptureService: queue can't be nil")
	return &CaptureService{
		Queue: queue,
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
