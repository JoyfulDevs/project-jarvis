package app

import (
	"context"
	"log/slog"

	channelconfig "github.com/devafterdark/project-jarvis/service/channelconfig/client"
)

func enableDailyScrumConfig(channel string, enabled bool) (*channelconfig.DailyScrumConfigV2, error) {
	client, err := channelconfig.NewClient()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Close(); err != nil {
			slog.Error("failed to close client", slog.Any("error", err))
		}
	}()
	config, err := client.GetDailyScrumConfig(context.Background(), channel)
	if err != nil {
		return nil, err
	}
	config.Enabled = enabled

	return client.SetDailyScrumConfig(context.Background(), channel, config)
}

func enableWeeklyReportConfig(channel string, enabled bool) (*channelconfig.WeeklyReportConfigV2, error) {
	client, err := channelconfig.NewClient()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Close(); err != nil {
			slog.Error("failed to close client", slog.Any("error", err))
		}
	}()
	config, err := client.GetWeeklyReportConfig(context.Background(), channel)
	if err != nil {
		return nil, err
	}
	config.Enabled = enabled

	return client.SetWeeklyReportConfig(context.Background(), channel, config)
}
