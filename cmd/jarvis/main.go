package main

import (
	"io"
	"log/slog"
	"os"

	"github.com/joyfuldevs/project-jarvis/cmd/jarvis/app"
)

func main() {
	filebeat, ok := os.LookupEnv("USE_FILEBEAT")
	if ok && filebeat == "yes" {
		logPath := "/var/log/filebeat/jarvis.log"
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o777)
		if err != nil {
			slog.Error("failed to open log file", slog.Any("error", err))
		} else {
			writer := io.MultiWriter(os.Stdout, file)
			handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{})
			slog.SetDefault(slog.New(handler))
		}
		defer func() {
			if err := file.Close(); err != nil {
				slog.Error("failed to close log file", slog.Any("error", err))
			}
		}()
	}

	slog.Info("jarvis starting")
	app.Run()
	slog.Info("jarvis finished")
}
