package server

import (
	"context"

	channelconfigv2 "github.com/genians/endpoint-lab-slack-bot/gen/go/channelconfig/v2"
)

type ChannelConfigV2 = channelconfigv2.ChannelConfig
type DailyScrumConfigV2 = channelconfigv2.DailyScrumConfig
type WeeklyReportConfigV2 = channelconfigv2.WeeklyReportConfig

type ServiceV2 interface {
	GetChannelConfig(ctx context.Context, channel string) (*ChannelConfigV2, error)
	SetDailyScrumConfig(ctx context.Context, channel string, config *DailyScrumConfigV2) error
	GetDailyScrumConfig(ctx context.Context, channel string) (*DailyScrumConfigV2, error)
	SetWeeklyReportConfig(ctx context.Context, channel string, config *WeeklyReportConfigV2) error
	GetWeeklyReportConfig(ctx context.Context, channel string) (*WeeklyReportConfigV2, error)

	AddScrumMessageHistory(ctx context.Context, channelID string, messageID float64) error
	GetScrumMessageHistory(ctx context.Context, channelID string) ([]float64, error)
}

type serverV2 struct {
	channelconfigv2.UnimplementedChannelConfigServiceServer

	serviceV2 ServiceV2
}

func (s *serverV2) GetChannelConfig(
	ctx context.Context,
	req *channelconfigv2.GetChannelConfigRequest,
) (*channelconfigv2.GetChannelConfigResponse, error) {
	config, err := s.serviceV2.GetChannelConfig(ctx, req.ChannelId)
	if err != nil {
		return nil, err
	}

	return &channelconfigv2.GetChannelConfigResponse{
		Config: config,
	}, nil
}

func (s *serverV2) GetDailyScrumConfig(
	ctx context.Context,
	req *channelconfigv2.GetDailyScrumConfigRequest,
) (*channelconfigv2.GetDailyScrumConfigResponse, error) {
	config, err := s.serviceV2.GetDailyScrumConfig(ctx, req.ChannelId)
	if err != nil {
		return nil, err
	}

	return &channelconfigv2.GetDailyScrumConfigResponse{
		Config: config,
	}, nil
}

func (s *serverV2) SetDailyScrumConfig(
	ctx context.Context,
	req *channelconfigv2.SetDailyScrumConfigRequest,
) (*channelconfigv2.SetDailyScrumConfigResponse, error) {
	if err := s.serviceV2.SetDailyScrumConfig(ctx, req.ChannelId, req.Config); err != nil {
		return nil, err
	}

	return &channelconfigv2.SetDailyScrumConfigResponse{
		Config: req.Config,
	}, nil
}

func (s *serverV2) GetWeeklyReportConfig(
	ctx context.Context,
	req *channelconfigv2.GetWeeklyReportConfigRequest,
) (*channelconfigv2.GetWeeklyReportConfigResponse, error) {
	config, err := s.serviceV2.GetWeeklyReportConfig(ctx, req.ChannelId)
	if err != nil {
		return nil, err
	}

	return &channelconfigv2.GetWeeklyReportConfigResponse{
		Config: config,
	}, nil
}

func (s *serverV2) SetWeeklyReportConfig(
	ctx context.Context,
	req *channelconfigv2.SetWeeklyReportConfigRequest,
) (*channelconfigv2.SetWeeklyReportConfigResponse, error) {
	if err := s.serviceV2.SetWeeklyReportConfig(ctx, req.ChannelId, req.Config); err != nil {
		return nil, err
	}

	return &channelconfigv2.SetWeeklyReportConfigResponse{
		Config: req.Config,
	}, nil
}

func (s *serverV2) AddScrumMessageHistory(
	ctx context.Context,
	req *channelconfigv2.AddScrumMessageHistoryRequest,
) (*channelconfigv2.AddScrumMessageHistoryResponse, error) {
	if err := s.serviceV2.AddScrumMessageHistory(ctx, req.ChannelId, req.MessageId); err != nil {
		return nil, err
	}

	return &channelconfigv2.AddScrumMessageHistoryResponse{}, nil
}

func (s *serverV2) GetScrumMessageHistory(
	ctx context.Context,
	req *channelconfigv2.GetScrumMessageHistoryRequest,
) (*channelconfigv2.GetScrumMessageHistoryResponse, error) {
	messageIDs, err := s.serviceV2.GetScrumMessageHistory(ctx, req.ChannelId)
	if err != nil {
		return nil, err
	}

	return &channelconfigv2.GetScrumMessageHistoryResponse{
		MessageIds: messageIDs,
	}, nil
}
