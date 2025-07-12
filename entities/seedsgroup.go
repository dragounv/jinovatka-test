package entities

type SeedsGroup struct {
	// Seeds in this group
	Seeds []*Seed

	// Unique randomly generated base32 encoded string with at least 128 bits of randomness.
	// Exact size is unspecified. This allowes the use of rand.Text to generate it.
	//
	// It serves as unique name for harvest, that is should be unpredictable and resilient to brute force guessing.
	// It will be used for generating URLs, that cannot be easily guessed, to preserve privacy of harvest creators.
	ShadowID string
}
