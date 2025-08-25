package client

import (
	"context"

	aigatewayv1 "github.com/devafterdark/project-jarvis/gen/go/aigateway/v1"
)

func (c *Client) GenerateText(
	ctx context.Context,
	texts []string,
) (string, error) {
	req := &aigatewayv1.GenerateTextRequest{
		Texts: texts,
	}
	resp, err := c.serviceV1.GenerateText(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}
