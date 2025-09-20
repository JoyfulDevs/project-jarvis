package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joyfuldevs/project-jarvis/service/aigateway/server"
)

func Run() {
	claudeKey, ok := os.LookupEnv("CLAUDE_API_KEY")
	if !ok {
		slog.Error("no such CLAUDE_API_KEY")
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	service := NewGatewayService(claudeKey)
	s := server.NewServer(
		server.WithServiceV1(service),
	)

	if err := s.Serve(ctx); err != nil {
		slog.Error("failed to start AI gateway service", slog.Any("error", err))
		return
	}
}
