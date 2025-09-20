package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joyfuldevs/project-jarvis/service/jarvis/server"
)

func Run() {
	appToken, ok := os.LookupEnv("SLACK_APP_TOKEN")
	if !ok {
		slog.Error("no such SLACK_APP_TOKEN")
		return
	}
	botToken, ok := os.LookupEnv("SLACK_BOT_TOKEN")
	if !ok {
		slog.Error("no such SLACK_BOT_TOKEN")
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	wg := sync.WaitGroup{}

	wg.Go(func() {
		slog.Info("starting jarvis service")
		jarvisService := NewJarvisService(appToken, botToken)
		jarvisServer := server.NewServer(
			server.WithServiceV1(jarvisService),
		)
		if err := jarvisServer.Serve(ctx); err != nil {
			slog.Error("failed to run service", slog.Any("error", err))
		}
		stop()
	})

	wg.Go(func() {
		slog.Info("starting jarvis bot")
		jarvisBot := NewJarvisBot(appToken, botToken)
		if err := jarvisBot.Run(ctx); err != nil {
			slog.Error("failed to run bot", slog.Any("error", err))
		}
		stop()
	})

	wg.Wait()
}
