package main

import (
	"log"
	"os"

	"github.com/khaleelsyed/codaVirtuale/internal/api"
	"github.com/khaleelsyed/codaVirtuale/internal/storage"
	"github.com/khaleelsyed/codaVirtuale/internal/types"
)

func main() {
	logger, err := types.NewLogger()
	if err != nil {
		log.Fatal("failed to initialise zap.logger")
	}
	defer logger.Sync()

	storage, err := storage.NewPostgresStorage(logger)
	if err != nil {
		return
	}

	if err = storage.Init(); err != nil {
		return
	}
	logger.Info("database connection is stable")

	listenAddress := os.Getenv("LISTEN_ADDRESS")

	server := api.NewAPIServer(listenAddress, storage, logger)
	server.Run()
}
