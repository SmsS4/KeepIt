package server

import (
	context "context"
	"fmt"
	"log"
	"net"

	"github.com/SmsS4/KeepIt/cache/ds"
	"google.golang.org/grpc"
)

type Server struct {
}

func (s *Server) Get(ctx context.Context, key *Key) (*Result, error) {
	return &Result{
		Value:     "test",
		MissCache: false,
		ActiveIp:  "",
	}, nil
}

func (s *Server) Put(ctx context.Context, keyValue *KeyValue) (*OprationResult, error) {
	return &OprationResult{
		ActiveIp: "",
	}, nil
}

func (s *Server) Clear(ctx context.Context, _ *Nil) (*OprationResult, error) {
	return &OprationResult{
		ActiveIp: "",
	}, nil
}

func RunServer(apiConfig ApiConfig, partionCache *ds.PartionCache) {
	log.Printf("Starting server on port %d", apiConfig.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", apiConfig.Port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", apiConfig.Port, err)
	}
	s := Server{}
	grpcServer := grpc.NewServer()
	RegisterCacheServiceServer(grpcServer, &s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port %d: %v", apiConfig.Port, err)
	}
}
