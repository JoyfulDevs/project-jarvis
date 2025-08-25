package app

import (
	"context"
	"log/slog"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"

	"github.com/devafterdark/project-jarvis/service/aigateway/server"
)

var _ server.ServiceV1 = (*GatewayService)(nil)

func (g *GatewayService) GenerateText(ctx context.Context, texts []string) (string, error) {
	contents := make([]anthropic.ContentBlockParamUnion, 0, len(texts))
	for _, text := range texts {
		contents = append(contents, anthropic.ContentBlockParamUnion{
			OfText: &anthropic.TextBlockParam{
				Text: text,
			},
		})
	}

	param := anthropic.MessageNewParams{
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			{
				Role:    anthropic.MessageParamRoleUser,
				Content: contents,
			},
		},
		Model: anthropic.ModelClaude4Sonnet20250514,
	}
	resp, err := g.client.Messages.New(
		ctx,
		param,
	)
	if err != nil {
		return "", err
	}

	slog.Info("token usage",
		slog.Int("input_tokens", int(resp.Usage.InputTokens)),
		slog.Int("output_tokens", int(resp.Usage.OutputTokens)),
	)
	builder := strings.Builder{}
	for _, content := range resp.Content {
		if len(content.Text) == 0 {
			continue
		}
		builder.WriteString(content.Text)
	}

	return builder.String(), nil
}
