package services

import (
	"context"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Generates random base32 ids. Can be used from multiple goroutines.
type IdGeneratorService struct {
	Closed          bool
	IsUniqueChecker IsUniqueChecker
	MaxRetries      int

	requests  chan IdRequest
	closeOnce sync.Once
}

// Creates new IdGeneratorService and starts backgroud goroutine.
// Please call IdGeneratorService.Close to gracefuly stop.
func NewIdGeneratorService(ctx context.Context, isUniqueChecker IsUniqueChecker, maxRetries int) *IdGeneratorService {
	service := &IdGeneratorService{
		requests:        make(chan IdRequest),
		IsUniqueChecker: isUniqueChecker,
		MaxRetries:      maxRetries,
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
	maxRetries := service.MaxRetries
	if maxRetries < 0 {
		maxRetries = 0
	}

	for ; maxRetries >= 0; maxRetries-- {
		id, err := service.getId()
		if err != nil {
			return "", err
		}

		isUnique, err := service.IsUniqueChecker.IsUnique(id)
		if err != nil {
			return "", fmt.Errorf("IdGeneratorService.GetId IsUniqueChecker returned error: %w", err)
		}

		if isUnique {
			err := service.IsUniqueChecker.Add(id)
			if err != nil {
				return "", fmt.Errorf("IdGeneratorService.GetId IsUniqueChecker returned error: %w", err)
			}
			return id, nil
		}
	}

	return "", errors.New("IdGeneratorService.GetId could not generate unique id and run out of retries")
}

func (service *IdGeneratorService) getId() (string, error) {
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
			id := string(outputBuffer[:desiredIdLength])

			request.Reply <- id
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

type IsUniqueChecker interface {
	// Return true if the id is unique, return false otherwise
	IsUnique(id string) (bool, error)
	Add(id string) error
}

// IsUniqueChecker that always returns true
type AlwaysUnique struct{}

func (*AlwaysUnique) IsUnique(id string) (bool, error) {
	return true, nil
}

func (*AlwaysUnique) Add(id string) error {
	return nil
}
