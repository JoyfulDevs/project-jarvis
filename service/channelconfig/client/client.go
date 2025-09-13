package client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	channelconfigv1 "github.com/joyfuldevs/project-jarvis/gen/go/channelconfig/v1"
	channelconfigv2 "github.com/joyfuldevs/project-jarvis/gen/go/channelconfig/v2"
)

type clientOptions struct {
	address string
	port    int
}

type Option func(*clientOptions)

func WithAddress(address string) Option {
	return func(opt *clientOptions) {
		opt.address = address
	}
}

func WithPort(port int) Option {
	return func(opt *clientOptions) {
		opt.port = port
	}
}

type Client struct {
	options    *clientOptions
	grpcClient *grpc.ClientConn
	serviceV1  channelconfigv1.ChannelConfigServiceClient
	serviceV2  channelconfigv2.ChannelConfigServiceClient
}

func NewClient(opts ...Option) (*Client, error) {
	options := &clientOptions{
		address: "channel-config-service",
		port:    50051,
	}
	for _, opt := range opts {
		opt(options)
	}

	grpcClient, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", options.address, options.port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		options:    options,
		grpcClient: grpcClient,
		serviceV1:  channelconfigv1.NewChannelConfigServiceClient(grpcClient),
		serviceV2:  channelconfigv2.NewChannelConfigServiceClient(grpcClient),
	}, nil
}

func (c *Client) Close() error {
	return c.grpcClient.Close()
}
