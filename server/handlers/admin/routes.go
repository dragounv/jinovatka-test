package admin

import (
	"jinovatka/assert"
	"log/slog"
	"net/http"
)

func Routes(router *http.ServeMux, args *RoutesArgs) {
	router.Handle("/admin/", NewAdminHandler(args.Log))
}

func NewRoutesArgs(log *slog.Logger) *RoutesArgs {
	assert.Must(log != nil, "admin/NewRoutesArgs: log can't be nil")
	return &RoutesArgs{
		Log: log,
	}
}

type RoutesArgs struct {
	Log *slog.Logger
}
