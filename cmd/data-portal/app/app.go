package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joyfuldevs/project-jarvis/service/dataportal/server"
)

func Run() {
	key, isExists := os.LookupEnv("DATA_PORTAL_AUTH_KEY")
	if !isExists {
		slog.Error("no such DATA_PORTAL_AUTH_KEY")
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := server.NewServer(
		server.WithServiceV1(NewDataPortalService(key)),
	)

	if err := server.Serve(ctx); err != nil {
		slog.Error("failed to start data portal service", slog.Any("error", err))
	}
}
