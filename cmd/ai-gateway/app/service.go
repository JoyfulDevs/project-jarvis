package app

import (
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type GatewayService struct {
	client anthropic.Client
}

func NewGatewayService(claudeKey string) *GatewayService {
	return &GatewayService{
		client: anthropic.NewClient(
			option.WithAPIKey(claudeKey),
		),
	}
}
