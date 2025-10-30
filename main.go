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
	"net"
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

	const defaultDBPath = "storage.db"
	sqlitePath, ok := os.LookupEnv("DB_PATH")
	if !ok {
		log.Warn("the database path is not set, using default " + defaultDBPath)
		sqlitePath = defaultDBPath
	}

	// Prepare db conection.
	db, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
	if err != nil {
		log.Error("could not open database connection", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Prepare queue client
	valkeyOptions := valkeyq.NewValkeyOptionsFromEnv()
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{net.JoinHostPort(valkeyOptions.Addr, valkeyOptions.Port)}})
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

	const defaultServerAdderss = "localhost:8080"
	serverAddress, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok {
		log.Warn("the server adress is not set, using default " + defaultServerAdderss)
		serverAddress = defaultServerAdderss
	}
	server := server.NewServer(
		stopSignal,
		log,
		serverAddress,
		initiatedServices,
	)

	// Start the server in new goroutine
	go server.ListenAndServe()
	log.Info("Server is listening at http://" + serverAddress)

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
