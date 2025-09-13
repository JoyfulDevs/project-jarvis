package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joyfuldevs/project-jarvis/service/channelconfig/server"
)

func Run() {
	addr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		slog.Error("no such REDIS_ADDR")
		return
	}

	pw, _ := os.LookupEnv("REDIS_PW")

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	redis := NewRedisService(addr, pw)
	s := server.NewServer(
		server.WithServiceV1(redis),
		server.WithServiceV2(redis),
	)

	if err := s.Serve(ctx); err != nil {
		slog.Error("failed to start channel config service", slog.Any("error", err))
	}
}
