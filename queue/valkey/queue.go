package valkeyq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"jinovatka/assert"
	"jinovatka/entities"
	"log/slog"
	"time"

	"github.com/valkey-io/valkey-go"
)

// Valkey (fork of Redis) implementation of Queue for capture requests and responses
type Queue struct {
	Log    *slog.Logger
	Client valkey.Client
}

func NewQueue(log *slog.Logger, client valkey.Client) *Queue {
	assert.Must(log != nil, "valkeyq/NewQueue: log can't be nil")
	assert.Must(client != nil, "valkeyq/NewQueue: client can't be nil")
	return &Queue{
		Log:    log,
		Client: client,
	}
}

const (
	RequestListKey = "queue:requests"
	ResultListKey  = "queue:results"
)

func (queue *Queue) Enqueue(ctx context.Context, request *entities.CaptureRequest) error {
	if request == nil {
		return errors.New("Queue.Enqueue recived nil request")
	}
	if request.SeedShadowID == "" {
		return errors.New("Queue.Enqueue recived request with no SeedShadowID")
	}
	if request.SeedURL == "" {
		return errors.New("Queue.Enqueue recived request with no SeedURL")
	}
	// TODO: Maybe create model struct for this?
	requestData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("Queue.Enqueue failed to marshal request to json: %w", err)
	}
	err = queue.Client.Do(ctx, queue.Client.B().Rpush().Key(RequestListKey).Element(string(requestData)).Build()).Error()
	if err != nil {
		return fmt.Errorf("Queue.Enqueue valkey client returned error: %w", err)
	}
	queue.Log.Info("Enqueued request", "URL", request.SeedURL, "ID", request.SeedShadowID)
	return nil
}

func (queue *Queue) AwaitResult(ctx context.Context, timeout time.Duration) (*entities.CaptureResult, error) {
	// This call blocks
	resultData, err := queue.Client.Do(ctx, queue.Client.B().Blpop().Key(ResultListKey).Timeout(timeout.Seconds()).Build()).AsBytes()
	if err != nil {
		return nil, fmt.Errorf("Queue.AwaitResult valkey client returned error: %w", err)
	}
	result := new(entities.CaptureResult)
	err = json.Unmarshal(resultData, result)
	if err != nil {
		return nil, fmt.Errorf("Queue.AwaitResult failed to unmarshal result from json: %w", err)
	}
	return result, nil
}
