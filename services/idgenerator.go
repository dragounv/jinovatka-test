package services

import (
	"context"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"math/rand"
	"sync"
	"time"
)

type IdGeneratorService struct {
	requests  chan IdRequest
	closeOnce sync.Once
	Closed    bool
}

// Creates new IdGeneratorService and starts backgroud goroutine.
// Please call IdGeneratorService.Close to gracefuly stop.
func NewIdGeneratorService(ctx context.Context) *IdGeneratorService {
	service := &IdGeneratorService{
		requests: make(chan IdRequest),
	}
	go service.run(ctx)
	return service
}

type IdRequest struct {
	Reply chan<- string
}

// Generate and return a base32 encoded string that should be used as id for Seeds and Groups.
// Can be called from multiple goroutines.
func (service *IdGeneratorService) GetId() (string, error) {
	if service.Closed {
		return "", errors.New("IdGeneratorService.GetId the service is closed")
	}

	replyChan := make(chan string)
	var err error

	// Recover from panic caused by reading from closed chan
	defer func() {
		if r := recover(); r != nil {
			replyChan = nil
			err = errors.New("IdGeneratorService.GetId stoped panic, the service is probably closed")
		}
	}()
	service.requests <- IdRequest{Reply: replyChan}
	id := <-replyChan
	return id, err
}

func (service *IdGeneratorService) run(ctx context.Context) {
	source := rand.NewSource(time.Now().UnixMicro())
	randomGenerator := rand.New(source)

	// The ids will be used in URLs and padding would look bad
	encoding := base32.StdEncoding.WithPadding(base32.NoPadding)

	outputBuffer := make([]byte, 8 /* need to be large enough for encoded value */)
	randNumBuffer := make([]byte, 4 /* bytes in uint32 */)

	for {
		select {
		case <-ctx.Done():
			service.Close()
			return
		case request, ok := <-service.requests:
			if !ok {
				return
			}

			randomValue := randomGenerator.Uint32()
			binary.NativeEndian.PutUint32(randNumBuffer, randomValue)
			encoding.Encode(outputBuffer, randNumBuffer)

			const desiredIdLength = 6 // If more than 7 chars is required then bigger random number is necessary
			request.Reply <- string(outputBuffer[:desiredIdLength])
		}
	}
}

// This will stop the goroutine asociated with this object
func (service *IdGeneratorService) Close() {
	service.closeOnce.Do(func() {
		close(service.requests)
		service.Closed = true
	})
}
