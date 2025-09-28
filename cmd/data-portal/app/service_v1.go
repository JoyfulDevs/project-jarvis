package app

import (
	"context"
	"log/slog"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/jh1104/publicapi/forecast"
	"github.com/jh1104/publicapi/specialday"
	"github.com/joyfuldevs/project-jarvis/pkg/kst"
	"github.com/joyfuldevs/project-jarvis/service/dataportal/server"
)

var _ server.ServiceV1 = (*DataPortalService)(nil)

func (s *DataPortalService) ListHolidays(ctx context.Context, year int, month int) ([]*server.HolidayV1, error) {
	resp, err := specialday.ListHolidays(ctx, specialday.NewParameters(year, month))
	if err != nil {
		return nil, err
	}

	result := make([]*server.HolidayV1, 0, len(resp.Body.Data.Items))
	for _, item := range resp.Body.Data.Items {
		t, err := time.Parse("20060102", strconv.Itoa(item.Date))
		if err != nil {
			slog.Warn("failed to parse holiday date", "date", item.Date, "error", err)
			continue
		}
		year, month, day := t.Date()
		result = append(result, &server.HolidayV1{
			Year:  int32(year),
			Month: int32(month),
			Day:   int32(day),
			Name:  item.Name,
		})
	}
	return result, nil
}

func (s *DataPortalService) GetUltraShortTermForecast(
	ctx context.Context,
	nx int32,
	ny int32,
) ([]*server.ForecastV1, error) {
	now := kst.Now()

	baseDate, baseTime := forecast.BaseForUltraShortTermForecast(now)
	params := forecast.Parameters{
		BaseDate:     baseDate,
		BaseTime:     baseTime,
		NX:           int(nx),
		NY:           int(ny),
		NumberOfRows: 10,
		PageNo:       1,
	}

	resp, err := forecast.GetUltraShortTermForecast(ctx, params)
	if err != nil {
		return nil, err
	}

	// 초단기 예보는 6시간 까지만 제공된다.
	var (
		result      = make([]*server.ForecastV1, 0, 6)
		forecastMap = make(map[string]*server.ForecastV1, 6)
		itemCh      = make(chan forecast.Item)
	)

	wg := sync.WaitGroup{}
	wg.Go(func() {
		for _, item := range resp.Body.Data.Items {
			itemCh <- item
		}
	})

	remains := (resp.Body.Total-params.NumberOfRows)/params.NumberOfRows + 1
	for range remains {
		wg.Add(1)
		params = params.NextPage()
		go func(params forecast.Parameters) {
			defer wg.Done()
			resp, err := forecast.GetUltraShortTermForecast(ctx, params)
			if err != nil {
				slog.Warn("failed to get ultra short term forecast", slog.Any("error", err))
				return
			}
			for _, item := range resp.Body.Data.Items {
				itemCh <- item
			}
		}(params)
	}

	go func() {
		wg.Wait()
		close(itemCh)
	}()

	for item := range itemCh {
		key := item.Date + item.Time
		if _, ok := forecastMap[key]; !ok {
			forecastMap[key] = &server.ForecastV1{
				Time: item.Date + item.Time,
			}
		}
		switch item.Category {
		case forecast.CategoryTemperature:
			t, err := strconv.Atoi(item.Value)
			if err != nil {
				slog.Warn("failed to parse temperature", "value", item.Value, "error", err)
				continue
			}
			forecastMap[key].Temperature = int32(t)
		case forecast.CategoryPrecipitation:
			forecastMap[key].Precipitation = forecast.PrecipitationCode(item.Value).String()
		case forecast.CategorySky:
			forecastMap[key].Sky = forecast.SkyCode(item.Value).String()
		}
	}

	for _, v := range forecastMap {
		result = append(result, v)
	}
	slices.SortFunc(result, func(a, b *server.ForecastV1) int {
		at, _ := time.Parse("200601021504", a.Time)
		bt, _ := time.Parse("200601021504", b.Time)
		return at.Compare(bt)
	})

	return result, nil
}
