package valkeyq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"jinovatka/assert"
	"jinovatka/entities"
	q "jinovatka/queue"
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
		return errors.New("Queue.Enqueue recieved nil request")
	}
	if request.SeedShadowID == "" {
		return errors.New("Queue.Enqueue recieved request with no SeedShadowID")
	}
	if request.SeedURL == "" {
		return errors.New("Queue.Enqueue recieved request with no SeedURL")
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

// WARNING: This function blocks indefinitely and should be run in separate goroutine.
//
// If timeout is zero, this function blocks until CaptureResult can be dequeued.
// If timeout is nonzero and no CaptureResult is available, this function blocks
// until timeout runs out and then returns QueueTimeoutError and nil CaptureResult.
func (queue *Queue) AwaitResult(ctx context.Context, timeout time.Duration) (*entities.CaptureResult, error) {
	// This call blocks. BLPOP returns array: [Key, Value]
	valkeyResult := queue.Client.Do(ctx, queue.Client.B().Blpop().Key(ResultListKey).Timeout(timeout.Seconds()).Build())

	// We need to get the message and handle errors returned by client.
	valkeyMessage, err := valkeyResult.ToMessage()
	if err != nil {
		return nil, fmt.Errorf("Queue.AwaitResult valkey client returned error: %w", err)
	}

	// If the call timed out then the message will be null.
	if valkeyMessage.IsNil() {
		return nil, fmt.Errorf("%w: BLPOP in Queue.AwaitResult", q.QueueTimeoutError)
	}

	// Call .ToArray to unwrap the the actual messages.
	valkeyMessageArray, err := valkeyMessage.ToArray()
	if err != nil { // This error is likely to be error from Client.Do, but may also be parsig error from .ToArray
		return nil, fmt.Errorf("Queue.AwaitResult failed to unwrap valkeyMessage: %w", err)
	}

	// This assertion will never fail, unless:
	// 1. A bug exists in previous part of this function (most likely)
	// 2. We are using wrong version of Valkey server or Redis server
	// 3. We are using wrong version of Valkey client library
	assert.Must(len(valkeyMessageArray) == 2, "Queue.AwaitResult: Array returned from Valkey client must have exactly two elements")
	// Drop the message containig key, keep only the value.
	valueMessage := valkeyMessageArray[1]

	// Convert to json. This just delegates to json.Unmarshal but saves as a call.
	result := new(entities.CaptureResult)
	err = valueMessage.DecodeJSON(result)
	if err != nil {
		return nil, fmt.Errorf("Queue.AwaitResult failed to unmarshal result from json: %w", err)
	}

	// Now I think I deserve a coffee.
	return result, nil
}
