package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	aigatewayv1 "github.com/joyfuldevs/project-jarvis/gen/go/aigateway/v1"
)

type serverOptions struct {
	serviceV1 ServiceV1
	port      int
}

type Option func(*serverOptions)

func WithPort(port int) Option {
	return func(opt *serverOptions) {
		opt.port = port
	}
}

func WithServiceV1(s ServiceV1) Option {
	return func(opt *serverOptions) {
		opt.serviceV1 = s
	}
}

type Server struct {
	options *serverOptions
}

func NewServer(opts ...Option) *Server {
	options := &serverOptions{
		port: 50051,
	}
	for _, opt := range opts {
		opt(options)
	}
	return &Server{
		options: options,
	}
}

func (s *Server) Serve(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.options.port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	if s.options.serviceV1 != nil {
		aigatewayv1.RegisterAIGatewayServiceServer(grpcServer, &serverV1{
			service: s.options.serviceV1,
		})
	} else {
		slog.Warn("serviceV1 is not set, skipping registration of v1 service")
	}

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	return grpcServer.Serve(listener)
}
