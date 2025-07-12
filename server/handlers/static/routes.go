package static

import (
	"io/fs"
	"jinovatka/assert"
	"log/slog"
	"net/http"
)

func Routes(router *http.ServeMux, args *RoutesArgs) {
	router.Handle("GET /static/", NewStaticHandler(args.Log, args.FS))
}

func NewRoutesArgs(log *slog.Logger, fs fs.FS) *RoutesArgs {
	assert.Must(log != nil, "static/NewRoutesArgs: log can't be nil")
	assert.Must(fs != nil, "static/NewRoutesArgs: fs can't be nil")
	return &RoutesArgs{
		Log: log,
		FS:  fs,
	}
}

type RoutesArgs struct {
	Log *slog.Logger
	FS  fs.FS
}
