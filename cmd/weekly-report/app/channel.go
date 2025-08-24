package app

import (
	"context"
	"log/slog"

	channelconfig "github.com/genians/endpoint-lab-slack-bot/service/channelconfig/client"
	jarvis "github.com/genians/endpoint-lab-slack-bot/service/jarvis/client"
)

func ListInvitedChannels() []string {
	client, err := jarvis.NewClient()
	if err != nil {
		slog.Error("failed to create client", slog.Any("error", err))
		return nil
	}
	defer func() {
		if err := client.Close(); err != nil {
			slog.Error("failed to close client", slog.Any("error", err))
		}
	}()
	channels, err := client.ListInvitedChannels(context.Background())
	if err != nil {
		slog.Error("failed to get invited channels", slog.Any("error", err))
		return nil
	}
	return channels
}

func WeeklyReportConfig(channel string) (*channelconfig.WeeklyReportConfigV2, error) {
	client, err := channelconfig.NewClient()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Close(); err != nil {
			slog.Error("failed to close client", slog.Any("error", err))
		}
	}()
	return client.GetWeeklyReportConfig(context.Background(), channel)
}
