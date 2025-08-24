package server

import (
	"context"

	dataportalv1 "github.com/genians/endpoint-lab-slack-bot/gen/go/dataportal/v1"
)

type HolidayV1 = dataportalv1.Holiday
type ForecastV1 = dataportalv1.Forecast

type ServiceV1 interface {
	ListHolidays(ctx context.Context, year int, month int) ([]*HolidayV1, error)
	GetUltraShortTermForecast(ctx context.Context, nx int32, ny int32) ([]*ForecastV1, error)
}

type serverV1 struct {
	dataportalv1.UnimplementedDataPortalServiceServer

	serviceV1 ServiceV1
}

func (s *serverV1) ListHolidays(
	ctx context.Context,
	req *dataportalv1.ListHolidaysRequest,
) (*dataportalv1.ListHolidaysResponse, error) {
	holidays, err := s.serviceV1.ListHolidays(ctx, int(req.Year), int(req.Month))
	if err != nil {
		return nil, err
	}
	resHolidays := make([]*dataportalv1.Holiday, 0, len(holidays))
	for _, h := range holidays {
		resHolidays = append(resHolidays, &dataportalv1.Holiday{
			Year:  req.Year,
			Month: req.Month,
			Day:   h.Day,
			Name:  h.Name,
		})
	}
	resp := &dataportalv1.ListHolidaysResponse{
		Holidays: resHolidays,
	}
	return resp, nil
}

func (s *serverV1) GetUltraShortTermForecast(
	ctx context.Context,
	req *dataportalv1.GetUltraShortTermForecastRequest,
) (*dataportalv1.GetUltraShortTermForecastResponse, error) {
	forecasts, err := s.serviceV1.GetUltraShortTermForecast(ctx, req.Nx, req.Ny)
	if err != nil {
		return nil, err
	}
	return &dataportalv1.GetUltraShortTermForecastResponse{
		Forecasts: forecasts,
	}, nil
}
