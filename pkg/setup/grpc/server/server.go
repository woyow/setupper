package server

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

const (
	tcpProto          = "tcp"
	setupLoggingKey   = "setup"
	setupLoggingValue = "grpc-server"
)

type Server struct {
	Server   *grpc.Server
	listener net.Listener
	log      *logrus.Logger
}

func NewServer(cfg *Config, log *logrus.Logger) *Server {
	address := os.Getenv(cfg.AddressEnvKey)
	if address == "" {
		panic(ErrEmptyAddress)
	}

	log.Info("setup grpc-server: NewServer listen on ", address)

	listener, err := net.Listen(tcpProto, address)
	if err != nil {
		grpclog.Fatalf("setup grpc-server: failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	reflection.Register(grpcServer)

	return &Server{
		Server:   grpcServer,
		listener: listener,
		log:      log,
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.log.WithContext(ctx).
		WithField(setupLoggingKey, setupLoggingValue).
		Info("Run - running grpc server")

	return s.Server.Serve(s.listener)
}

func (s *Server) Shutdown() error {
	s.log.WithField(setupLoggingKey, setupLoggingValue).
		Info("Shutdown - graceful shutdown grpc server")
	s.Server.GracefulStop()
	return nil
}
