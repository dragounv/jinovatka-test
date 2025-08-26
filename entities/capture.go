package entities

// Information necessary for capturing/crawling a page.
type CaptureRequest struct {
	// The URL adress we want to capture.
	SeedURL string `json:"seedURL"`
	// The ShadowID of Seed we want to capture. This will be used in CaptureResult.
	SeedShadowID string `json:"seedShadowID"`
	// The status of the request.
	State CaptureState `json:"state"`
}

func NewRequestFromSeed(seed *Seed) *CaptureRequest {
	return &CaptureRequest{
		SeedURL:      seed.URL,
		SeedShadowID: seed.ShadowID,
		State:        NotEnqueued,
	}
}

type CaptureState string

const (
	// The request is newly created and is not enqueued.
	NotEnqueued CaptureState = "NotEnqueued"
	// The request was enqueued for processing.
	Pending CaptureState = "Pending"
	// The capture was successful.
	DoneSuccess CaptureState = "DoneSuccess"
	// The capture failed.
	DoneFailure CaptureState = "DoneFailure"
)

func (state CaptureState) IsCaptureState() bool {
	return state == NotEnqueued ||
		state == Pending ||
		state == DoneSuccess ||
		state == DoneFailure
}

type CaptureResult struct {
	// ShadowID of the seed to which the result belongs to.
	SeedShadowID string `json:"seedShadowID"`
	// Was the capture completed
	Done bool `json:"done"`
	// Recieved errors
	ErrorMessages []string `json:"errorMessages"`

	CaptureMetadata *CaptureMetadata `json:"captureMetadata"`
}

type CaptureMetadata struct {
	// CDXJ timestamp of the instatnt the capture was taken as recorded in index/WARC https://specs.webrecorder.net/cdxj/0.1.0/#timestamp
	Timestamp string `json:"timestamp"`
	// URL from CDXJ JSON block.
	CapturedUrl string `json:"capturedUrl"`
}
