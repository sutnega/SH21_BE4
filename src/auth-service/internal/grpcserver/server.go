package grpcserver

import (
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	srv *grpc.Server
}

func New() *Server {
	return &Server{
		srv: grpc.NewServer(),
	}
}

func (s *Server) GRPC() *grpc.Server {
	return s.srv
}

func (s *Server) Serve(lis net.Listener) error {
	return s.srv.Serve(lis)
}

func (s *Server) GracefulStop() {
	s.srv.GracefulStop()
}
