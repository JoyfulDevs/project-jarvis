package main

import (
	"log/slog"

	"github.com/joyfuldevs/project-jarvis/cmd/ai-gateway/app"
)

func main() {
	slog.Info("AI gateway service starting")
	app.Run()
	slog.Info("AI gateway service finished")
}
