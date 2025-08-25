package client

import (
	"context"

	jarvisv1 "github.com/devafterdark/project-jarvis/gen/go/jarvis/v1"
)

// 슬랙봇이 초대된 채널 목록을 가져온다.
func (c *Client) ListInvitedChannels(ctx context.Context) ([]string, error) {
	req := &jarvisv1.ListInvitedChannelsRequest{}
	resp, err := c.serviceClient.ListInvitedChannels(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.ChannelIds, nil
}

// 슬랙봇으로 메시지를 보낸다.
// 전송이 성공하면 메시지의 타임스탬프를 리턴한다.
func (c *Client) SendMessage(
	ctx context.Context,
	channel string,
	message string,
	blocks []byte,
	markdown bool,
) (float64, error) {
	req := &jarvisv1.SendSlackMessageRequest{
		ChannelId: channel,
		Message:   message,
		Blocks:    blocks,
		Markdown:  markdown,
	}
	resp, err := c.serviceClient.SendSlackMessage(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.Timestamp, nil
}

// 사용자 프로필을 가져온다.
func (c *Client) GetUserProfile(
	ctx context.Context,
	userID string,
) (*jarvisv1.UserProfile, error) {
	req := &jarvisv1.GetUserProfileRequest{
		UserId: userID,
	}
	resp, err := c.serviceClient.GetUserProfile(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Profile, nil
}
