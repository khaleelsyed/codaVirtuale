package main

import (
	"os"

	"github.com/khaleelsyed/codaVirtuale/internal/api"
	"github.com/khaleelsyed/codaVirtuale/internal/storage"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	storage, err := storage.NewMockStorage()
	if err != nil {
		sugar.Panic(err)
	}

	listenAddress := os.Getenv("LISTEN_ADDRESS")

	server := api.NewAPIServer(listenAddress, storage, sugar)

	server.Run()
}
