package server

import (
	"context"

	aigatewayv1 "github.com/genians/endpoint-lab-slack-bot/gen/go/aigateway/v1"
)

type ServiceV1 interface {
	GenerateText(ctx context.Context, texts []string) (string, error)
}

type serverV1 struct {
	aigatewayv1.UnimplementedAIGatewayServiceServer

	service ServiceV1
}

func (s *serverV1) GenerateText(
	ctx context.Context,
	req *aigatewayv1.GenerateTextRequest,
) (*aigatewayv1.GenerateTextResponse, error) {
	result, err := s.service.GenerateText(ctx, req.Texts)
	if err != nil {
		return nil, err
	}
	resp := &aigatewayv1.GenerateTextResponse{
		Text: result,
	}
	return resp, nil
}
