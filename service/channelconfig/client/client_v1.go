package client

import (
	"context"

	channelconfigv1 "github.com/genians/endpoint-lab-slack-bot/gen/go/channelconfig/v1"
)

type FeatureV1 = channelconfigv1.Feature

const (
	FeatureV1Unspecified  = channelconfigv1.Feature_FEATURE_UNSPECIFIED
	FeatureV1DailyScrum   = channelconfigv1.Feature_FEATURE_DAILY_SCRUM
	FeatureV1WeeklyReport = channelconfigv1.Feature_FEATURE_WEEKLY_REPORT
)

// 채널이 구독하는 이벤트를 추가한다.
func (c *Client) Subscribe(
	ctx context.Context,
	channel string,
	feature FeatureV1,
) ([]FeatureV1, error) {
	req := &channelconfigv1.SubscribeRequest{
		ChannelId: channel,
		Feature:   feature,
	}
	resp, err := c.serviceV1.Subscribe(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Features, nil
}

// 채널이 구독하는 이벤트를 제거한다.
func (c *Client) Unsubscribe(
	ctx context.Context,
	channel string,
	feature FeatureV1,
) ([]FeatureV1, error) {
	req := &channelconfigv1.UnsubscribeRequest{
		ChannelId: channel,
		Feature:   feature,
	}
	resp, err := c.serviceV1.Unsubscribe(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Features, nil
}

// 채널이 구독하는 모든 이벤트를 제거한다.
func (c *Client) UnsubscribeAll(
	ctx context.Context,
	channel string,
) error {
	req := &channelconfigv1.UnsubscribeAllRequest{
		ChannelId: channel,
	}
	_, err := c.serviceV1.UnsubscribeAll(ctx, req)
	return err
}

// 채널이 구독하는 이벤트 목록을 가져온다.
func (c *Client) ListSubscriptions(
	ctx context.Context,
	channel string,
) ([]FeatureV1, error) {
	req := &channelconfigv1.ListSubscriptionsRequest{
		ChannelId: channel,
	}
	resp, err := c.serviceV1.ListSubscriptions(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Features, nil
}

// 특정 이벤트를 구독하는 채널 목록을 가져온다.
func (c *Client) ListChannelsByFeature(
	ctx context.Context,
	feature FeatureV1,
) ([]string, error) {
	req := &channelconfigv1.ListChannelsByFeatureRequest{
		Feature: feature,
	}
	resp, err := c.serviceV1.ListChannelsByFeature(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.ChannelIds, nil
}
