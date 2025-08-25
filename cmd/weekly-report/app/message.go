package app

import (
	"context"
	"log/slog"

	jarvis "github.com/devafterdark/project-jarvis/service/jarvis/client"
)

func WeeklyReportMessage() string {
	return `오늘은 이번 주 마지막 업무일입니다.
한 주 동안 고생 많으셨습니다.
주간업무보고 작성 부탁드립니다.`
}

func SendReportMessage(channels []string, message string) {
	client, err := jarvis.NewClient()
	if err != nil {
		slog.Error("failed to create client", slog.Any("error", err))
		return
	}
	defer func() {
		if err := client.Close(); err != nil {
			slog.Error("failed to close client", slog.Any("error", err))
		}
	}()

	for _, channel := range channels {
		ts, err := client.SendMessage(context.Background(), channel, message, nil, false)
		if err != nil {
			slog.Error("failed to send message", slog.Any("error", err))
		} else {
			slog.Info("send report message", slog.String("channel", channel), slog.Float64("timestamp", ts))
		}
	}
}
