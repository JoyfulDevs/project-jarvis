package app

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/joyfuldevs/project-jarvis/pkg/dataportal"
	"github.com/joyfuldevs/project-jarvis/pkg/kst"
	"github.com/joyfuldevs/project-jarvis/service/dataportal/server"
)

var _ server.ServiceV1 = (*DataPortalService)(nil)

func (s *DataPortalService) ListHolidays(ctx context.Context, year int, month int) ([]*server.HolidayV1, error) {
	req := &dataportal.ListHolidaysRequest{
		Year:  year,
		Month: month,
	}
	resp, err := s.client.ListHolidays(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		slog.Warn("received empty response body for holidays")
		return nil, errors.New("empty response body")
	}

	result := make([]*server.HolidayV1, 0, len(resp.Body.Data.Items))
	for _, item := range resp.Body.Data.Items {
		t, err := time.Parse("20060102", item.Date)
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
	baseDate := now.Format("20060102")
	// 초단기 예보는 매 시간 30분 단위로 제공된다.
	baseTime := func() string {
		if now.Minute() > 30 {
			return now.Format("15")
		} else {
			return now.Add(-time.Hour).Format("15")
		}
	}() + "30"

	var (
		page  = 1
		count = 10
	)

	newRequest := func(page int) *dataportal.UltraShortTermForecastRequest {
		return &dataportal.UltraShortTermForecastRequest{
			Page:     page,
			Count:    count,
			BaseDate: baseDate,
			BaseTime: baseTime,
			NX:       int(nx),
			NY:       int(ny),
		}
	}

	req := newRequest(page)
	resp, err := s.client.GetUltraShortTermForecast(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.Body == nil {
		slog.Warn("received empty response body for forecast")
		return nil, errors.New("empty response body")
	}

	// 초단기 예보는 6시간 까지만 제공된다.
	var (
		result      = make([]*server.ForecastV1, 0, 6)
		forecastMap = make(map[string]*server.ForecastV1, 6)
		itemCh      = make(chan dataportal.UltraShortTermForecastItem)
	)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, item := range resp.Body.Data.Items {
			itemCh <- item
		}
	}()

	for page*count < resp.Body.Total {
		page += 1
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			req := newRequest(page)
			resp, err := s.client.GetUltraShortTermForecast(ctx, req)
			if err != nil {
				slog.Warn("failed to get ultra short term forecast", slog.Any("error", err))
				return
			}
			if resp.Body == nil {
				slog.Warn("received empty response body for forecast",
					slog.Int("page", req.Page),
					slog.Int("count", req.Count),
				)
				return
			}
			for _, item := range resp.Body.Data.Items {
				itemCh <- item
			}
		}(page)
	}
	go func() {
		wg.Wait()
		close(itemCh)
	}()

	for item := range itemCh {
		key := item.ForecastDate + item.ForecastTime
		if _, ok := forecastMap[key]; !ok {
			forecastMap[key] = &server.ForecastV1{
				Time: item.ForecastDate + item.ForecastTime,
			}
		}
		switch item.Category {
		case dataportal.CategoryTemperature:
			t, err := strconv.Atoi(item.ForecastValue)
			if err != nil {
				slog.Warn("failed to parse temperature", "value", item.ForecastValue, "error", err)
				continue
			}
			forecastMap[key].Temperature = int32(t)
		case dataportal.CategoryPrecipitation:
			forecastMap[key].Precipitation = dataportal.PrecipitationCode(item.ForecastValue).String()
		case dataportal.CategorySky:
			forecastMap[key].Sky = dataportal.SkyCode(item.ForecastValue).String()
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
