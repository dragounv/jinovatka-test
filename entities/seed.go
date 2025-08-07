package entities

import (
	"time"
)

type Seed struct {
	// The original URL od seed.
	URL string

	// If the seed is public, hten it may be visible on the main page and can be searched for.
	Public bool

	// State of capture of the seed.
	State CaptureState

	// URL of the archived resource. Must be empty unless seed was sucessfully harvested ( state is HarvestedSucessfully).
	ArchivalURL string

	// Time of harvest. Should be eqivalent to the time used to generate ArchivalUrl.
	// Must be zero value if seed wasn't harvested yet (state is NotHarvested).
	HarvestedAt time.Time

	// Unique randomly generated base32 encoded string with at least 128 bits of randomness.
	// Exact size is unspecified. This allowes the use of rand.Text to generate it.
	//
	// It serves as unique name for harvest, that is should be unpredictable and resilient to brute force guessing.
	// It will be used for generating URLs, that cannot be easily guessed, to preserve privacy of harvest creators.
	ShadowID string
}
