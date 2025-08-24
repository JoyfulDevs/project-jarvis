package client

import (
	"context"

	dataportalv1 "github.com/genians/endpoint-lab-slack-bot/gen/go/dataportal/v1"
)

type HolidayV1 = dataportalv1.Holiday
type ForecastV1 = dataportalv1.Forecast

// 주어진 연도와 월에 해당하는 공휴일 목록을 조회한다.
func (c *Client) ListHolidays(ctx context.Context, year int, month int) ([]*HolidayV1, error) {
	req := &dataportalv1.ListHolidaysRequest{Year: int32(year), Month: int32(month)}
	resp, err := c.serviceClient.ListHolidays(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Holidays, nil
}

// 주어진 좌표(nx, ny)에 대한 현재 초단기 예보를 확인한다.
func (c *Client) GetUltraShortTermForecast(ctx context.Context, nx int32, ny int32) ([]*ForecastV1, error) {
	req := &dataportalv1.GetUltraShortTermForecastRequest{Nx: nx, Ny: ny}
	resp, err := c.serviceClient.GetUltraShortTermForecast(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Forecasts, nil
}
