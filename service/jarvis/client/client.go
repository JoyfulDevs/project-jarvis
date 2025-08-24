package client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	jarvisv1 "github.com/genians/endpoint-lab-slack-bot/gen/go/jarvis/v1"
)

type clientOptions struct {
	address string
	port    int
}

type Option func(*clientOptions)

func WithAddress(address string) Option {
	return func(opts *clientOptions) {
		opts.address = address
	}
}

func WithPort(port int) Option {
	return func(opts *clientOptions) {
		opts.port = port
	}
}

type Client struct {
	options       *clientOptions
	grpcClient    *grpc.ClientConn
	serviceClient jarvisv1.JarvisServiceClient
}

func NewClient(opts ...Option) (*Client, error) {
	options := &clientOptions{
		address: "jarvis-service",
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
		options:       options,
		grpcClient:    grpcClient,
		serviceClient: jarvisv1.NewJarvisServiceClient(grpcClient),
	}, nil
}

func (c *Client) Close() error {
	return c.grpcClient.Close()
}
