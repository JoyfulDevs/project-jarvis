package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/devafterdark/project-jarvis/service/dataportal/server"
)

func Run() {
	key, isExists := os.LookupEnv("DATA_PORTAL_AUTH_KEY")
	if !isExists {
		slog.Error("no such DATA_PORTAL_AUTH_KEY")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	server := server.NewServer(
		server.WithServiceV1(NewDataPortalService(key)),
	)

	if err := server.Serve(ctx); err != nil {
		slog.Error("failed to start data portal service", slog.Any("error", err))
	}
}
