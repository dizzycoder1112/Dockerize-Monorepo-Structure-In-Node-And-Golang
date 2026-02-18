package grpc

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Logger interface for optional logging
type Logger interface {
	Info(msg string, context ...any)
}

// defaultLogger uses standard library log
type defaultLogger struct{}

func (d *defaultLogger) Info(msg string, context ...any) {
	log.Println(append([]any{msg}, context...)...)
}

type Server struct {
	grpcServer *grpc.Server
	log        Logger
}

type ServerDeps struct {
	Log Logger // optional, uses default if nil
}

func NewServer(deps *ServerDeps) *Server {
	var logger Logger = &defaultLogger{}
	if deps != nil && deps.Log != nil {
		logger = deps.Log
	}

	return &Server{
		grpcServer: grpc.NewServer(),
		log:        logger,
	}
}

// GrpcServer exposes the underlying grpc.Server for service registration
func (s *Server) GrpcServer() *grpc.Server {
	return s.grpcServer
}

// EnableReflection enables gRPC reflection for debugging tools like grpcurl
func (s *Server) EnableReflection() {
	reflection.Register(s.grpcServer)
}

func (s *Server) start(port string, ready chan<- struct{}) error {
	address := fmt.Sprintf(":%s", port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	if ready != nil {
		close(ready)
	}

	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *Server) GracefulStop() {
	s.grpcServer.GracefulStop()
	s.log.Info("gRPC server stopped")
}

// Run starts the server in a goroutine and blocks until shutdown signal is received.
// Returns control to the caller for cleanup. Caller should call GracefulStop().
func (s *Server) Run(port string) error {
	ready := make(chan struct{})
	errChan := make(chan error, 1)

	go func() {
		if err := s.start(port, ready); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ready:
		s.log.Info("gRPC server listening", map[string]any{"port": port})
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	return nil
}
