package server

import (
	"context"
	"jinovatka/server/handlers"
	"jinovatka/server/handlers/admin"
	"jinovatka/server/handlers/group"
	"jinovatka/server/handlers/index"
	"jinovatka/server/handlers/static"
	"jinovatka/services"
	"log/slog"
	"net"
	"net/http"
	"time"
)

func NewServer(ctx context.Context, log *slog.Logger, addr string, services *services.Services) *http.Server {
	// Create router
	mux := http.NewServeMux()
	router := handlers.NewRouterHandler(mux)

	// Add all handlers to the router
	router.AddHandlers(
		index.NewIndexHandler(log, services.SeedService),
		static.NewStaticHandler(log, staticFiles /* from embed.go */),
		group.NewGroupHandler(log, services.SeedService),
		admin.NewAdminHandler(log),
	)

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		BaseContext:  func(l net.Listener) context.Context { return ctx },
	}

	return server
}
