package app

import (
	"context"
	"log/slog"
	"slices"
	"time"

	dataportal "github.com/devafterdark/project-jarvis/service/dataportal/client"
)

func IsHoliday(t time.Time) bool {
	client, err := dataportal.NewClient()
	if err != nil {
		return false
	}
	defer func() {
		if err := client.Close(); err != nil {
			slog.Error("failed to close client", slog.Any("error", err))
		}
	}()

	holidays, err := client.ListHolidays(context.Background(), t.Year(), int(t.Month()))
	if err != nil {
		slog.Error("failed to list holidays", slog.Any("error", err))
		return false
	}
	day := int32(t.Day())
	return slices.ContainsFunc(holidays, func(holiday *dataportal.HolidayV1) bool {
		return holiday.Day == day
	})
}
