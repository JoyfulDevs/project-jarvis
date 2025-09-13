package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	channelconfigv1 "github.com/joyfuldevs/project-jarvis/gen/go/channelconfig/v1"
	channelconfigv2 "github.com/joyfuldevs/project-jarvis/gen/go/channelconfig/v2"
)

type serverOptions struct {
	serviceV1 ServiceV1
	serviceV2 ServiceV2
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

func WithServiceV2(s ServiceV2) Option {
	return func(opt *serverOptions) {
		opt.serviceV2 = s
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
		channelconfigv1.RegisterChannelConfigServiceServer(grpcServer, &serverV1{
			serviceV1: s.options.serviceV1,
		})
	} else {
		slog.Warn("serviceV1 is not set, skipping registration of v1 service")
	}
	if s.options.serviceV2 != nil {
		channelconfigv2.RegisterChannelConfigServiceServer(grpcServer, &serverV2{
			serviceV2: s.options.serviceV2,
		})
	} else {
		slog.Warn("serviceV2 is not set, skipping registration of v2 service")
	}

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	return grpcServer.Serve(listener)
}
