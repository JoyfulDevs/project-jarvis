package main

import (
	"log/slog"

	"github.com/joyfuldevs/project-jarvis/cmd/data-portal/app"
)

func main() {
	slog.Info("data portal service starting")
	app.Run()
	slog.Info("data portal service finished")
}
