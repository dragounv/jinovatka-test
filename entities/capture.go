package entities

// Information necessary for capturing/crawling a page.
type CaptureRequest struct {
	// The URL adress we want to capture.
	SeedURL string
	// The ShadowID of Seed we want to capture. This will be used in CaptureResult.
	SeedShadowID string
	// The status of the request.
	Status CaptureState
}

func NewRequestFromSeed(seed *Seed) *CaptureRequest {
	return &CaptureRequest{
		SeedURL:      seed.URL,
		SeedShadowID: seed.ShadowID,
		Status:       NewRequest,
	}
}

type CaptureState string

const (
	// The request is newly created and is not enqueued.
	NewRequest CaptureState = "NewRequest"
	// The request was enqueued for processing.
	Pending CaptureState = "Pending"
	// The capture was successful.
	DoneSuccess CaptureState = "DoneSuccess"
	// The capture failed.
	DoneFailure CaptureState = "DoneFailure"
)

type CaptureResult struct {
	// ShadowID of the seed to which the result belongs to.
	SeedShadowID string
	// Was the capture completed
	Done bool
	// Recieved errors
	ErrorMessages []string
}
