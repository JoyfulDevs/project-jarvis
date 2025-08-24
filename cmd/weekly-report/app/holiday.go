package app

import (
	"context"
	"log/slog"
	"slices"
	"time"

	dataportal "github.com/genians/endpoint-lab-slack-bot/service/dataportal/client"
)

func IsLastWorkday(t time.Time) bool {
	client, err := dataportal.NewClient()
	if err != nil {
		slog.Error("failed to create open api client", slog.Any("error", err))
		return false
	}

	year, month, day := t.Date()
	holidays, err := client.ListHolidays(context.Background(), year, int(month))
	if err != nil {
		slog.Error("failed to connect open api service", slog.Any("error", err))
	}

	// 오늘이 휴일인 경우 마지막 업무일이 아니다.
	isHoliday := slices.ContainsFunc(holidays, func(holiday *dataportal.HolidayV1) bool {
		return holiday.Day == int32(day)
	})

	if isHoliday {
		return false
	}

	// 오늘이 금요일인 경우 마지막 업무일이다.
	if t.Weekday() == time.Friday {
		return true
	}

	// 토요일까지 남은 요일 중 업무일이 있는지 확인한다.
	for next := t.AddDate(0, 0, 1); next.Weekday() != time.Saturday; next = next.AddDate(0, 0, 1) {
		y, m, d := next.Date()
		// 년, 월이 변경될 때 휴일 정보를 갱신한다.
		if year != y || month != m {
			year, month = y, m
			holidays, err = client.ListHolidays(context.Background(), year, int(month))
			if err != nil {
				slog.Error("failed to connect open api service", slog.Any("error", err))
			}
		}

		isHoliday := slices.ContainsFunc(holidays, func(holiday *dataportal.HolidayV1) bool {
			return holiday.Day == int32(d)
		})

		if !isHoliday {
			slog.Info("remaining workday", slog.Int("day", d))
			return false
		}
	}
	// 남은 모든 요일이 휴일이면 메시지를 보낸다.
	return true
}
