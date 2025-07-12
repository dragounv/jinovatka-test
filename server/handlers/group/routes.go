package group

import (
	"jinovatka/assert"
	"jinovatka/services"
	"log/slog"
	"net/http"
)

func Routes(router *http.ServeMux, args *RoutesArgs) {
	router.Handle("GET /seeds/{id}", NewGroupHandler(args.Log, args.SeedService))
}

func NewRoutesArgs(log *slog.Logger, seedService *services.SeedService) *RoutesArgs {
	assert.Must(log != nil, "group/NewRoutesArgs: log can't be nil")
	assert.Must(seedService != nil, "group/NewRoutesArgs: log can't be nil")
	return &RoutesArgs{
		Log:         log,
		SeedService: seedService,
	}
}

type RoutesArgs struct {
	Log         *slog.Logger
	SeedService *services.SeedService
}
