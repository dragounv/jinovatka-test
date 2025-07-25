package entities

// Information necessary for capturing/crawling a page.
type CaptureRequest struct {
	// The URL adress we want to capture.
	SeedURL string
	// The ShadowID of Seed we want to capture. This will be used in CaptureResult.
	SeedShadowID string
	// The status of the request.
	Status CaptureStatus
}

func NewRequestFromSeed(seed *Seed) *CaptureRequest {
	return &CaptureRequest{
		SeedURL:      seed.URL,
		SeedShadowID: seed.ShadowID,
		Status:       NewRequest,
	}
}

type CaptureResult struct {
	// ShadowID of the seed to which the result belongs to.
	SeedShadowID string
	// Status of the result.
	Status CaptureStatus
}

type CaptureStatus string

const (
	// The request is newly created and is not enqueued.
	NewRequest CaptureStatus = "NewRequest"
	// The request was enqueued for processing.
	Pending CaptureStatus = "Pending"
	// The capture was successful.
	DoneSuccess CaptureStatus = "DoneSuccess"
	// The capture failed.
	DoneFailure CaptureStatus = "DoneFailure"
)
