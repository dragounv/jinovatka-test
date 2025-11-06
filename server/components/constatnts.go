package components

import (
	"fmt"
	"html"
	"jinovatka/assert"
	"net/url"
)

// Using global variable is not the best solution, but it will do for now.
// Passing the constatnts as parameter from from handlers would be a lot better,
// but it would require a larger refactor.

var (
	g_constants             *ComponentConstants
	g_setConstantsWasCalled = false
)

// This function can only be called once.
func SetComponentConstants(constants *ComponentConstants) {
	assert.Must(constants != nil, "components:SetComponentConstants argument constants can't be nil")
	assert.Must(!g_setConstantsWasCalled, "components:SetComponentConstants can only be called once")
	g_setConstantsWasCalled = true
	g_constants = constants
}

func Constants() *ComponentConstants {
	return g_constants
}

func NewComponentConstants(
	serverHost,
	staticPath,
	seedDetailPath,
	groupDetailPath,
	waybackRedirectPath string,
) *ComponentConstants {
	return &ComponentConstants{
		serverHost:          serverHost,
		staticPath:          staticPath,
		seedDetailPath:      seedDetailPath,
		groupDetailPath:     groupDetailPath,
		waybackRedirectPath: waybackRedirectPath,
	}
}

// Values that can be used by templ components, that are initialized at server
// start and don't change during runtime.
type ComponentConstants struct {
	// Public address of the application as reachable from internet.
	// This value does not have to (and in production should not be) the same
	// as the actual server address.
	//
	// If the application is behind proxy that uses URL paths to route to different applications,
	// then the path segment must be included in this adress.
	//
	// Also remember to add a correct protocol (http:// or https://)
	//
	// So in dev it could be ip like: http://127.0.0.1:8080
	// In prod: https://this-app.example
	// Or when sharing domain with other apps: https://some-apps.example/this-app
	serverHost string

	// Constants for paths
	staticPath          string
	seedDetailPath      string
	groupDetailPath     string
	waybackRedirectPath string
}

func (constants *ComponentConstants) GetServerHost() string {
	return constants.serverHost
}

// Appends path to serverHost.
func (constants *ComponentConstants) GetFullURL(path string) string {
	url, err := url.JoinPath(constants.serverHost, path)
	if err != nil {
		return fmt.Sprintf("<!-- Error: %s -->", html.EscapeString(err.Error()))
	}
	return url
}

func (constants *ComponentConstants) GetStaticPath() string {
	return constants.staticPath
}

func (constants *ComponentConstants) GetSeedPath() string {
	return constants.seedDetailPath
}

func (constants *ComponentConstants) GetGroupPath() string {
	return constants.groupDetailPath
}

func (constants *ComponentConstants) GetWaybackRedirectPath() string {
	return constants.waybackRedirectPath
}
