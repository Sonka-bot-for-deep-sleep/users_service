package grpc

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	i "github.com/Sonka-bot-for-deep-sleep/common/pkg/interceptors"
	pb "github.com/Sonka-bot-for-deep-sleep/proto_files/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	handler pb.UsersServiceServer
	logger  *zap.Logger
}

func New(handler pb.UsersServiceServer, logger *zap.Logger) *Server {
	return &Server{
		handler: handler,
		logger:  logger,
	}
}

func (s *Server) StartServer(port string) error {
	addr := fmt.Sprintf("localhost:%s", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("StartServer: failed create tcp net listener: %w", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(i.LoggingInterceptor))
	pb.RegisterUsersServiceServer(grpcServer, s.handler)

	errCH := make(chan error, 1)
	s.logger.Info("Server start work", zap.String("addr", addr),
		zap.Time("server_start_time", time.Now()))

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			errCH <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGBUS, syscall.SIGTERM)

	select {
	case <-quit:
		s.logger.Info("Server shutting down", zap.Time("time_shutting_down", time.Now()))
	case <-errCH:
		err := <-errCH
		s.logger.Error("Failed start server", zap.Error(err))
		return fmt.Errorf("StartServer: failed start server work: %w", err)
	}

	grpcServer.GracefulStop()
	s.logger.Info("Server Stop work", zap.Time("time_stopped_server", time.Now()))
	return nil
}
