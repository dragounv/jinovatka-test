package components

import (
	"jinovatka/entities"
	"path"
)

func prettyPrintCaptureState(state entities.CaptureState) string {
	return entities.PrettyPrintCaptureState(state)
}

// NOTE
// Be careful with using url.JoinPath. It does not work unless the base is valid URL with protocol.
// path.Join is better when joining just path segments. Especially when we don't know if they start and end with slashes.

// URL that will redirect to the archived page in wayback when the seed is captured
func shortWaybackLink(seed *entities.Seed) string {
	return Constants().GetFullURL(path.Join(Constants().GetWaybackRedirectPath(), seed.ShadowID))
}

// URL for the seeds detail view
func seedViewLink(seed *entities.Seed) string {
	return Constants().GetFullURL(path.Join(Constants().GetSeedPath(), seed.ShadowID))
}

// URL for the group view
func groupViewLink(group *entities.SeedsGroup) string {
	return Constants().GetFullURL(path.Join(Constants().GetGroupPath(), group.ShadowID))
}

// Get full path to static file "filename"
func fullStaticPath(filename string) string {
	return Constants().GetFullURL(path.Join(Constants().GetStaticPath(), filename))
}
