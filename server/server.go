package server

import (
	"context"
	"jinovatka/server/handlers"
	"jinovatka/services"
	"log/slog"
	"net"
	"net/http"
	"time"
)

func NewServer(ctx context.Context, log *slog.Logger, addr string, services *services.Services) *http.Server {
	// routerArgs := &RouterArgs{
	// 	Log:      log,
	// 	Services: services.NewServices(log),
	// }
	// router := NewRouter(routerArgs)

	router := http.NewServeMux()
	handlers.Routes(router, handlers.NewRoutesArgs(
		log,
		services,
		staticFiles,
	))

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		BaseContext:  func(l net.Listener) context.Context { return ctx },
	}

	return server
}
