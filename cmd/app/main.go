package main

import (
	"os"

	"github.com/khaleelsyed/codaVirtuale/internal/api"
	"github.com/khaleelsyed/codaVirtuale/internal/storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const TraceLevel zapcore.Level = -2 // TRACE is lower than DEBUG (-1)

type SugarWithTrace struct {
	*zap.SugaredLogger
}

func (l *SugarWithTrace) Trace(msg string, keysAndValues ...interface{}) {
	if ce := l.Desugar().Check(TraceLevel, msg); ce != nil {
		ce.Write(zap.Any("extra", keysAndValues))
	}
}

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	storage, err := storage.NewPostgresStorage(sugar)
	if err != nil {
		return
	}

	if err = storage.Init(); err != nil {
		return
	}

	listenAddress := os.Getenv("LISTEN_ADDRESS")

	server := api.NewAPIServer(listenAddress, storage, sugar)

	server.Run()
}
