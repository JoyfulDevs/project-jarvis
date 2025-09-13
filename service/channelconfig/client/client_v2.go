package client

import (
	"context"

	channelconfigv2 "github.com/joyfuldevs/project-jarvis/gen/go/channelconfig/v2"
)

type ChannelConfigV2 = channelconfigv2.ChannelConfig
type DailyScrumConfigV2 = channelconfigv2.DailyScrumConfig
type WeeklyReportConfigV2 = channelconfigv2.WeeklyReportConfig

// 채널의 전체 설정을 조회한다.
func (c *Client) GetChannelConfig(
	ctx context.Context,
	channel string,
) (*ChannelConfigV2, error) {
	req := &channelconfigv2.GetChannelConfigRequest{
		ChannelId: channel,
	}
	resp, err := c.serviceV2.GetChannelConfig(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Config, nil
}

// 주어진 채널에 대한 스크럼 알림 설정을 조회 한다.
func (c *Client) GetDailyScrumConfig(
	ctx context.Context,
	channel string,
) (*DailyScrumConfigV2, error) {
	req := &channelconfigv2.GetDailyScrumConfigRequest{
		ChannelId: channel,
	}
	resp, err := c.serviceV2.GetDailyScrumConfig(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Config, nil
}

// 주어진 채널에 대한 스크럼 알림 설정을 저장한다.
func (c *Client) SetDailyScrumConfig(
	ctx context.Context,
	channel string,
	config *DailyScrumConfigV2,
) (*DailyScrumConfigV2, error) {
	req := &channelconfigv2.SetDailyScrumConfigRequest{
		ChannelId: channel,
		Config:    config,
	}
	resp, err := c.serviceV2.SetDailyScrumConfig(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Config, nil
}

// 주어진 채널에 대한 주간 보고 알림 설정을 조회 한다.
func (c *Client) GetWeeklyReportConfig(
	ctx context.Context,
	channel string,
) (*WeeklyReportConfigV2, error) {
	req := &channelconfigv2.GetWeeklyReportConfigRequest{
		ChannelId: channel,
	}
	resp, err := c.serviceV2.GetWeeklyReportConfig(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Config, nil
}

// 주어진 채널에 대한 주간 보고 알림 설정을 저장한다.
func (c *Client) SetWeeklyReportConfig(
	ctx context.Context,
	channel string,
	config *WeeklyReportConfigV2,
) (*WeeklyReportConfigV2, error) {
	req := &channelconfigv2.SetWeeklyReportConfigRequest{
		ChannelId: channel,
		Config:    config,
	}
	resp, err := c.serviceV2.SetWeeklyReportConfig(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Config, nil
}

func (c *Client) AddScrumMessageHistory(
	ctx context.Context,
	channelID string,
	messageID float64,
) error {
	req := &channelconfigv2.AddScrumMessageHistoryRequest{
		ChannelId: channelID,
		MessageId: messageID,
	}
	_, err := c.serviceV2.AddScrumMessageHistory(ctx, req)
	return err
}

func (c *Client) GetScrumMessageHistory(
	ctx context.Context,
	channelID string,
) ([]float64, error) {
	req := &channelconfigv2.GetScrumMessageHistoryRequest{
		ChannelId: channelID,
	}
	resp, err := c.serviceV2.GetScrumMessageHistory(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.MessageIds, nil
}
