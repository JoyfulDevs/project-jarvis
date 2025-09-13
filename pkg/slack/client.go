package slack

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/joyfuldevs/project-jarvis/pkg/rest"
)

const domain = "https://slack.com/api"

type Client struct {
	AppToken string
	BotToken string
}

func (c *Client) GetWebSocketURL(ctx context.Context) (*GetWebSocketURLResponse, error) {
	path := "/apps.connections.open"
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + c.AppToken,
	}

	data, err := rest.NewClient(domain).RequestAPI(
		ctx, "POST", path,
		rest.WithHeaders(header),
	)
	if err != nil {
		return nil, err
	}

	resp := &GetWebSocketURLResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if !resp.OK {
		return nil, errors.New(resp.Error)
	}

	return resp, nil
}

func (c *Client) ListChannels(ctx context.Context, req *ListChannelsRequest) (*ListChannelsResponse, error) {
	path := "/conversations.list"
	if req != nil {
		path = path + req.URLParams()
	}
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + c.BotToken,
	}

	data, err := rest.NewClient(domain).RequestAPI(
		ctx, "GET", path,
		rest.WithHeaders(header),
	)
	if err != nil {
		return nil, err
	}

	resp := &ListChannelsResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) PostMessage(ctx context.Context, req *PostMessageRequest) (*PostMessageResponse, error) {
	path := "/chat.postMessage"
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + c.BotToken,
	}
	body, err := json.Marshal(*req)
	if err != nil {
		return nil, err
	}

	data, err := rest.NewClient(domain).RequestAPI(
		ctx, "POST", path,
		rest.WithHeaders(header),
		rest.WithBody(body),
	)
	if err != nil {
		return nil, err
	}

	resp := &PostMessageResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) ListMessages(ctx context.Context, req *ListMessagesRequest) (*ListMessagesResponse, error) {
	path := "/conversations.history"
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + c.BotToken,
	}
	body, err := json.Marshal(*req)
	if err != nil {
		return nil, err
	}

	data, err := rest.NewClient(domain).RequestAPI(
		ctx, "POST", path,
		rest.WithHeaders(header),
		rest.WithBody(body),
	)
	if err != nil {
		return nil, err
	}

	resp := &ListMessagesResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) ListReplies(ctx context.Context, req *ListRepliesRequest) (*ListRepliesResponse, error) {
	path := "/conversations.replies"
	param := map[string]string{
		"channel": req.Channel,
		"ts":      strconv.FormatFloat(req.Timestamp, 'f', -1, 64),
	}
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + c.BotToken,
	}

	data, err := rest.NewClient(domain).RequestAPI(
		ctx, "GET", path,
		rest.WithHeaders(header),
		rest.WithParams(param),
	)
	if err != nil {
		return nil, err
	}

	resp := &ListRepliesResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
	path := "/users.profile.get"

	param := map[string]string{
		"user": userID,
	}
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + c.BotToken,
	}

	data, err := rest.NewClient(domain).RequestAPI(
		ctx, "GET", path,
		rest.WithHeaders(header),
		rest.WithParams(param),
	)
	if err != nil {
		return nil, err
	}
	profile := &UserProfile{}
	if err := json.Unmarshal(data, profile); err != nil {
		return nil, err
	}

	return profile, nil
}
