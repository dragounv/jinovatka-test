package handlers

import (
	"io/fs"
	"jinovatka/assert"
	"jinovatka/server/handlers/admin"
	"jinovatka/server/handlers/group"
	"jinovatka/server/handlers/index"
	"jinovatka/server/handlers/static"
	"jinovatka/services"
	"log/slog"
	"net/http"
)

func Routes(router *http.ServeMux, args *RoutesArgs) {
	index.Routes(router, index.NewRoutesArgs(args.Log, args.Services.SeedService))
	static.Routes(router, static.NewRoutesArgs(args.Log, args.StaticFiles))
	admin.Routes(router, admin.NewRoutesArgs(args.Log))
	group.Routes(router, group.NewRoutesArgs(args.Log, args.Services.SeedService))
}

func NewRoutesArgs(log *slog.Logger, services *services.Services, staticFiles fs.FS) *RoutesArgs {
	assert.Must(log != nil, "handlers/NewRoutesArgs: log can't be nil")
	assert.Must(services != nil, "handlers/NewRoutesArgs: services can't be nil")
	assert.Must(staticFiles != nil, "handlers/NewRoutesArgs: staticFiles can't be nil")
	return &RoutesArgs{
		Log:         log,
		Services:    services,
		StaticFiles: staticFiles,
	}
}

type RoutesArgs struct {
	Log         *slog.Logger
	Services    *services.Services
	StaticFiles fs.FS
}
