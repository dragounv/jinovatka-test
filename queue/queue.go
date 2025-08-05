package queue

import (
	"context"
	"errors"
	"jinovatka/entities"
	"time"
)

// Package queue defines an interface that is used to enqueue capture requests and recieve capture results.
// It is expected, that the interface implementations will use external inmemory database as is Redis or Valkey
// and thus it is possible that the requests/results may get lost.
// Always store the request data to persistent db before enquing them.
// This way it should be possible to repeat the request.
//
// This is the producer side of the queue. We are enquing requests and recieving responses.
// The consumer side will be implemented by capture software running as separate process. Possibly on separete server.

// Queue that enqueues CaptureRequests and recieves CpatureResults.
type Queue interface {
	// Enqueue the request. This method should never block the caller.
	Enqueue(context.Context, *entities.CaptureRequest) error
	// Fetch result from queue.
	// If there are no reusts waiting in the queue then this method should block until result is recieved.
	// If timeout is zero, no timeout will be used.
	AwaitResult(ctx context.Context, timeout time.Duration) (*entities.CaptureResult, error)
}

// Use to cath potential timeouts that are not supposed to propagate.
// Only methods with timeout can return this error.
var QueueTimeoutError = errors.New("operation timed out")
