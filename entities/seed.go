package entities

import (
	"time"
)

type InvalidSeedStateError struct {
	InvalidString string
}

func (err InvalidSeedStateError) Error() string {
	if err.InvalidString == "" {
		return "empty string is not valid SeedState"
	}
	return string(err.InvalidString) + " is not valid SeedState"
}

// State of the seed. Determines if the seed was alredy harvested or not.
type SeedState string

// These constatns represent the states a seed can be in.
// The values are used to persist the state. They should be changed carefully!
const (
	NotHarvested         SeedState = "not_harvested"
	HarvestedSucessfully SeedState = "harvested"
	HarvestFailed        SeedState = "failed"
)

// Set used to test if string is valid state. The values of the map are not important,
// only the keys matter. Remember to update the map whenever the constants are updated.
var seedStateSet = map[string]SeedState{
	string(NotHarvested):         NotHarvested,
	string(HarvestedSucessfully): HarvestedSucessfully,
	string(HarvestFailed):        HarvestFailed,
}

// Returns SeedState or InvalidSeedStateError if stateString is invalid.
func StateFromString(stateString string) (SeedState, error) {
	if !IsValidSeedState(stateString) {
		return "", InvalidSeedStateError{stateString}
	}
	return SeedState(stateString), nil
}

func IsValidSeedState(stateString string) bool {
	_, ok := seedStateSet[stateString]
	return ok
}

type Seed struct {
	// ID uint
	// CreatedAt,
	// UpdatedAt time.Time
	// // May be nil.
	// DeletedAt *time.Time

	// The original URL od seed.
	URL string

	// If the seed is public, hten it may be visible on the main page and can be searched for.
	Public bool

	// Determines if the seed was alredy harvested or not.
	State SeedState

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

	// Name of WACZ file that stores the resource
	// Filename sql.NullString

	// Harvest to which the resource belongs
	// HarvestID uint
	// Harvest   Harvest
}
