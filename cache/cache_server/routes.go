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
	partionCache *ds.PartionCache
}

func (s *Server) Get(ctx context.Context, key *Key) (*Result, error) {
	value, miss, err := s.partionCache.Get(key.Key)
	return &Result{
		Value:     value,
		MissCache: miss,
		ActiveIp:  "",
	}, err
}

func (s *Server) Put(ctx context.Context, keyValue *KeyValue) (*OprationResult, error) {
	s.partionCache.Put(keyValue.Key, keyValue.Value)
	return &OprationResult{
		ActiveIp: "",
	}, nil
}

func (s *Server) Clear(ctx context.Context, _ *Nil) (*OprationResult, error) {
	s.partionCache.ClearAll()
	return &OprationResult{
		ActiveIp: "",
	}, nil
}

func RunServerCache(apiConfig ApiConfig, partionCache *ds.PartionCache) {
	log.Printf("Starting server on port %d", apiConfig.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", apiConfig.Port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", apiConfig.Port, err)
	}
	s := Server{
		partionCache: partionCache,
	}
	grpcServer := grpc.NewServer()
	RegisterCacheServiceServer(grpcServer, &s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port %d: %v", apiConfig.Port, err)
	}
}

type Distribution struct {
}

func RunServerDistribution() {

}
