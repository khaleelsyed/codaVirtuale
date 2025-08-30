package main

import (
	"log"
	"os"

	"github.com/khaleelsyed/codaVirtuale/internal/api"
	"github.com/khaleelsyed/codaVirtuale/internal/storage"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	storage, err := storage.NewMockStorage()
	if err != nil {
		log.Panic(err)
	}

	listenAddress := os.Getenv("LISTEN_ADDRESS")

	server := api.NewAPIServer(listenAddress, storage, logger)

	server.Run()
}
