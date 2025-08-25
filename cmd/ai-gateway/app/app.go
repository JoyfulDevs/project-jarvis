package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/devafterdark/project-jarvis/service/aigateway/server"
)

func Run() {
	claudeKey, ok := os.LookupEnv("CLAUDE_API_KEY")
	if !ok {
		slog.Error("no such CLAUDE_API_KEY")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	service := NewGatewayService(claudeKey)
	s := server.NewServer(
		server.WithServiceV1(service),
	)

	if err := s.Serve(ctx); err != nil {
		slog.Error("failed to start AI gateway service", slog.Any("error", err))
		return
	}
}
