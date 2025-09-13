package server

import (
	"context"

	jarvisv1 "github.com/joyfuldevs/project-jarvis/gen/go/jarvis/v1"
)

type UserProfileV1 = jarvisv1.UserProfile

type ServiceV1 interface {
	ListInvitedChannels(ctx context.Context) ([]string, error)
	SendSlackMessage(ctx context.Context, channel string, message string, blocksData []byte, markdown bool) (float64, error)
	GetUserProfile(ctx context.Context, userID string) (*UserProfileV1, error)
}

type serverV1 struct {
	jarvisv1.UnimplementedJarvisServiceServer

	service ServiceV1
}

func (s *serverV1) ListInvitedChannels(
	ctx context.Context,
	req *jarvisv1.ListInvitedChannelsRequest,
) (*jarvisv1.ListInvitedChannelsResponse, error) {
	channelIds, err := s.service.ListInvitedChannels(ctx)
	if err != nil {
		return nil, err
	}
	return &jarvisv1.ListInvitedChannelsResponse{
		ChannelIds: channelIds,
	}, nil
}

func (s *serverV1) SendSlackMessage(
	ctx context.Context,
	req *jarvisv1.SendSlackMessageRequest,
) (*jarvisv1.SendSlackMessageResponse, error) {
	timestamp, err := s.service.SendSlackMessage(
		ctx,
		req.ChannelId,
		req.Message,
		req.Blocks,
		req.Markdown,
	)
	if err != nil {
		return nil, err
	}
	return &jarvisv1.SendSlackMessageResponse{
		Timestamp: timestamp,
	}, nil
}

func (s *serverV1) GetUserProfile(
	ctx context.Context,
	req *jarvisv1.GetUserProfileRequest,
) (*jarvisv1.GetUserProfileResponse, error) {
	profile, err := s.service.GetUserProfile(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &jarvisv1.GetUserProfileResponse{
		UserId:  req.UserId,
		Profile: profile,
	}, nil
}
