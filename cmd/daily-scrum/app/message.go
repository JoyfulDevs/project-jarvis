package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/devafterdark/project-jarvis/pkg/kst"
	channelconfig "github.com/devafterdark/project-jarvis/service/channelconfig/client"
	jarvis "github.com/devafterdark/project-jarvis/service/jarvis/client"
)

func DailyScrumMessage(t time.Time) string {
	weekday := kst.Weekday(t)
	format := fmt.Sprintf("2006-01-02 (%s) 스크럼", weekday)
	message := t.Format(format)
	return message
}

func SendScrumMessage(channels []string, message string) {
	jarvisClient, err := jarvis.NewClient()
	if err != nil {
		slog.Error("failed to create client", slog.Any("error", err))
		return
	}
	defer func() {
		if err := jarvisClient.Close(); err != nil {
			slog.Warn("failed to close client", slog.Any("error", err))
		}
	}()
	configClient, err := channelconfig.NewClient()
	if err != nil {
		slog.Error("failed to create channel config client", slog.Any("error", err))
		return
	}
	defer func() {
		if err := configClient.Close(); err != nil {
			slog.Warn("failed to close channel config client", slog.Any("error", err))
		}
	}()

	wg := sync.WaitGroup{}
	for _, channel := range channels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ts, err := jarvisClient.SendMessage(context.Background(), channel, message, nil, false)
			if err != nil {
				slog.Error("failed to send message", slog.Any("error", err))
				return
			}
			err = configClient.AddScrumMessageHistory(context.Background(), channel, ts)
			if err != nil {
				slog.Error("failed to add scrum message history", slog.Any("error", err))
			}

			slog.Info("send scrum message", slog.String("channel", channel), slog.Float64("timestamp", ts))
		}()
	}
	wg.Wait()
}
