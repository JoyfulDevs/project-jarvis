package server

import (
	"context"

	channelconfigv1 "github.com/devafterdark/project-jarvis/gen/go/channelconfig/v1"
)

type FeatureV1 = channelconfigv1.Feature

type ServiceV1 interface {
	Subscribe(ctx context.Context, channel string, feature FeatureV1) ([]FeatureV1, error)
	Unsubscribe(ctx context.Context, channel string, feature FeatureV1) ([]FeatureV1, error)
	UnsubscribeAll(ctx context.Context, channel string) error
	ListSubscriptions(ctx context.Context, channel string) ([]FeatureV1, error)
	ListChannelsByFeature(ctx context.Context, feature FeatureV1) ([]string, error)
}

type serverV1 struct {
	channelconfigv1.UnimplementedChannelConfigServiceServer

	serviceV1 ServiceV1
}

func (s *serverV1) Subscribe(
	ctx context.Context,
	req *channelconfigv1.SubscribeRequest,
) (*channelconfigv1.SubscribeResponse, error) {
	features, err := s.serviceV1.Subscribe(ctx, req.ChannelId, req.Feature)
	if err != nil {
		return nil, err
	}
	return &channelconfigv1.SubscribeResponse{
		Features: features,
	}, nil
}

func (s *serverV1) Unsubscribe(
	ctx context.Context,
	req *channelconfigv1.UnsubscribeRequest,
) (*channelconfigv1.UnsubscribeResponse, error) {
	features, err := s.serviceV1.Unsubscribe(ctx, req.ChannelId, req.Feature)
	if err != nil {
		return nil, err
	}

	return &channelconfigv1.UnsubscribeResponse{
		Features: features,
	}, nil
}

func (s *serverV1) UnsubscribeAll(
	ctx context.Context,
	req *channelconfigv1.UnsubscribeAllRequest,
) (*channelconfigv1.UnsubscribeAllResponse, error) {
	if err := s.serviceV1.UnsubscribeAll(ctx, req.ChannelId); err != nil {
		return nil, err
	}
	return &channelconfigv1.UnsubscribeAllResponse{}, nil
}

func (s *serverV1) ListSubscriptions(
	ctx context.Context,
	req *channelconfigv1.ListSubscriptionsRequest,
) (*channelconfigv1.ListSubscriptionsResponse, error) {
	features, err := s.serviceV1.ListSubscriptions(ctx, req.ChannelId)
	if err != nil {
		return nil, err
	}

	return &channelconfigv1.ListSubscriptionsResponse{
		Features: features,
	}, nil
}

func (s *serverV1) ListChannelsByFeature(
	ctx context.Context,
	req *channelconfigv1.ListChannelsByFeatureRequest,
) (*channelconfigv1.ListChannelsByFeatureResponse, error) {
	channelIds, err := s.serviceV1.ListChannelsByFeature(ctx, req.Feature)
	if err != nil {
		return nil, err
	}

	return &channelconfigv1.ListChannelsByFeatureResponse{
		ChannelIds: channelIds,
	}, nil
}
