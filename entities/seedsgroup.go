package entities

type SeedsGroup struct {
	// Seeds in this group
	Seeds []*Seed

	// Unique string. Possibly base32 encoded.
	// Exact length is unspecified. Must be usable in URL.
	ShadowID string
}
