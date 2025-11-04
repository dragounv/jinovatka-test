package entities

import (
	"time"
)

type Seed struct {
	// The original URL of the seed.
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

	// Unique string. Possibly base32 encoded.
	// Exact length is unspecified. Must be usable in URL.
	ShadowID string
}
