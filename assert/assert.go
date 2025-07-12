package assert

import (
	"os"
)

// Function used to write error message to stderr and exit immediatly and ungracefully.
func Die(message string) {
	// Ignore errors, we are going down anyways.
	_, _ = os.Stderr.WriteString(message)
	_ = os.Stderr.Sync()
	os.Exit(1)
}

// Expression must be true or the assertion crashes the app.
// Message will be printed on fail.
func Must(expresion bool, message string) {
	if !expresion {
		Die("assertion failed: " + message + "\n")
	}
}

// Safely get error message.
func AddErrorMessage(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
