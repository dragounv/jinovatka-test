package main

import (
	"context"
	valkeyq "jinovatka/queue/valkey"
	"jinovatka/server"
	"jinovatka/services"
	"jinovatka/storage"
	gormStorage "jinovatka/storage/gorm"
	"jinovatka/utils"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/valkey-io/valkey-go"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	log := slog.New(slog.Default().Handler())

	// Prepare db conection.
	db, err := gorm.Open(sqlite.Open("storage.db"), &gorm.Config{})
	if err != nil {
		log.Error("could not open database connection", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Prepare queue client
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
	if err != nil {
		log.Error("failed to create valkey client", "error", err.Error())
	}

	// Catch SIGINT and SIGHUP. Prepare gentle shutdown.
	// TODO: There are more signals that need catching
	signals := []os.Signal{os.Interrupt}
	if runtime.GOOS == "linux" {
		signals = append(signals, syscall.SIGHUP)
	}
	stopSignal, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()
	utils.ShutdownFunc = stop // Setup function, that can be used in cases, where shutdown of the server is necessary.

	seedRepository := gormStorage.NewSeedRepository(log, db)
	repository := storage.NewRepository(seedRepository)

	queue := valkeyq.NewQueue(log, client)

	initiatedServices := services.NewServices(log, repository, queue)

	const addr = "localhost:8080"
	server := server.NewServer(
		stopSignal,
		log,
		addr,
		initiatedServices,
	)

	// Start the server in new goroutine
	go server.ListenAndServe()
	log.Info("Server is listening at http://" + addr)

	// Start listening for results from queue
	initiatedServices.CaptureService.ListenForResults(stopSignal)
	log.Info("CaptureService is listening for CaptureResults")

	// Wait for interupt
	<-stopSignal.Done()
	// Wait for shutdown (or timeout and go eat dirt)
	shutdownTimeout, stop := context.WithTimeout(context.Background(), 120*time.Second)
	defer stop()
	err = server.Shutdown(shutdownTimeout)
	if err != nil {
		log.Error("shutdown timeout run out", slog.String("error", err.Error()))
	}
	log.Info("Server shutdown")
}
