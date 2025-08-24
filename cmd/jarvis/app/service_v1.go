package app

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/genians/endpoint-lab-slack-bot/pkg/slack"
	"github.com/genians/endpoint-lab-slack-bot/pkg/slack/blockkit"
	"github.com/genians/endpoint-lab-slack-bot/service/jarvis/server"
)

var _ server.ServiceV1 = (*JarvisService)(nil)

func (j *JarvisService) ListInvitedChannels(ctx context.Context) ([]string, error) {
	client := &slack.Client{
		AppToken: j.AppToken,
		BotToken: j.BotToken,
	}
	req := &slack.ListChannelsRequest{
		Types: slack.PrivateChannel,
	}
	resp, err := client.ListChannels(ctx, req)
	if err != nil {
		return nil, err
	}

	if !resp.OK {
		return nil, errors.New(resp.Error)
	}

	channels := make([]string, 0, len(resp.Channels))
	for _, channel := range resp.Channels {
		if channel.IsMember {
			channels = append(channels, channel.ID)
		}
	}

	return channels, nil
}

func (j *JarvisService) SendSlackMessage(
	ctx context.Context,
	channel string,
	message string,
	blocksData []byte,
	markdown bool,
) (float64, error) {
	req := &slack.PostMessageRequest{
		Channel:  channel,
		Text:     message,
		Markdown: markdown,
	}
	if len(blocksData) > 0 {
		blocks := make([]blockkit.SlackBlock, 0)
		if err := json.Unmarshal(blocksData, &blocks); err != nil {
			return 0, err
		}
		req.Blocks = blocks
	}

	client := &slack.Client{
		AppToken: j.AppToken,
		BotToken: j.BotToken,
	}
	resp, err := client.PostMessage(ctx, req)
	if err != nil {
		return 0, err
	}
	if !resp.OK {
		return 0, errors.New(resp.Error)
	}

	return resp.Timestamp, nil
}

func (j *JarvisService) GetUserProfile(
	ctx context.Context,
	userID string,
) (*server.UserProfileV1, error) {
	client := &slack.Client{
		AppToken: j.AppToken,
		BotToken: j.BotToken,
	}
	resp, err := client.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	profile := &server.UserProfileV1{
		Title:       resp.Title,
		DisplayName: resp.DisplayName,
		RealName:    resp.RealName,
		FirstName:   resp.FirstName,
		LastName:    resp.LastName,
		Email:       resp.Email,
		Phone:       resp.Phone,
	}

	return profile, nil
}
