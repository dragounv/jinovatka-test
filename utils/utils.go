package utils

import (
	"context"
	"log/slog"
	"net/http"
)

var ShutdownFunc context.CancelFunc

func LogRequestInfo(r *http.Request) slog.Attr {
	return slog.Group(
		"request",
		slog.String("path", r.URL.Path),
		slog.String("pattern", r.Pattern),
		slog.String("method", r.Method),
	)
}
