package server

import (
	context "context"
	"fmt"
	"log"
	"net"

	"github.com/SmsS4/KeepIt/cache/ds"
	"github.com/SmsS4/KeepIt/cache/utils"
	"google.golang.org/grpc"
)

type Server struct {
	activeIp     string
	myIp         string
	config       ApiConfig
	checkAgain   bool
	partionCache *ds.PartionCache
}

func (s *Server) CheckActiveIp() bool {
	if s.activeIp == s.myIp && !s.checkAgain {
		return false
	}
	for _, instance := range s.config.Instances {
		conn := s.makeConnection(instance)
		_, err := conn.Check(context.Background(), &Nil{})
		if s.config.Ip == instance.Ip && s.config.Port == instance.Port {
			s.activeIp = s.myIp
			log.Printf("New active is %s (me)", s.activeIp)
			return false
		}
		if err == nil {
			s.activeIp = fmt.Sprintf("%s:%d", instance.Ip, instance.Port)
			log.Printf("New active is %s", s.activeIp)
			return true
		}
	}
	log.Fatal("what ?")
	return false
}

func (s *Server) Get(ctx context.Context, key *Key) (*Result, error) {
	if s.CheckActiveIp() {
		return &Result{
			Value:     "",
			MissCache: false,
			ActiveIp:  s.activeIp,
		}, nil
	}
	value, miss, err := s.partionCache.Get(key.Key)
	return &Result{
		Value:     value,
		MissCache: miss,
		ActiveIp:  "",
	}, err
}

func (s *Server) Put(ctx context.Context, keyValue *KeyValue) (*OprationResult, error) {
	if s.CheckActiveIp() {
		return &OprationResult{
			ActiveIp: s.activeIp,
		}, nil
	}
	s.partionCache.Put(keyValue.Key, keyValue.Value)
	return &OprationResult{
		ActiveIp: "",
	}, nil
}

func (s *Server) Clear(ctx context.Context, _ *Nil) (*OprationResult, error) {
	if s.CheckActiveIp() {
		return &OprationResult{
			ActiveIp: s.activeIp,
		}, nil
	}
	s.partionCache.ClearAll(true)
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
	gatewayServer := Server{
		activeIp:     apiConfig.Instances[0].Ip,
		partionCache: partionCache,
		myIp:         fmt.Sprintf("%s:%d", apiConfig.Ip, apiConfig.Port),
		config:       apiConfig,
	}
	// distributionServer := Distribution{
	// 	server: &gatewayServer,
	// 	config: apiConfig,
	// }
	grpcServer := grpc.NewServer()
	RegisterCacheServiceServer(grpcServer, &gatewayServer)
	// RegisterDistributionServiceServer(grpcServer, &distributionServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port %d: %v", apiConfig.Port, err)
	}
}

func (s *Server) makeConnection(instance Instance) CacheServiceClient {
	var conn *grpc.ClientConn
	addr := fmt.Sprintf("%s:%d", instance.Ip, instance.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not makeConnection: %s", err)
	}
	defer conn.Close()

	return NewCacheServiceClient(conn)
}

func (s *Server) Use(ctx context.Context, keyValue *DistKeyValue) (*Nil, error) {
	log.Printf(
		"Use %s %s %s:%d %v",
		keyValue.Key,
		keyValue.Value,
		keyValue.Data.FromIp,
		keyValue.Data.FromPort,
		keyValue.Data.Distribute,
	)
	s.partionCache.Touch(keyValue.Key, keyValue.Value)
	if keyValue.Data.Distribute {
		log.Print("Distribute is true")
		keyValue.Data.Distribute = false
		for _, instance := range s.config.Instances {
			if instance.Ip != keyValue.Data.FromIp && instance.Port != int(keyValue.Data.FromPort) {
				log.Printf(
					"Distribute to %s:%d",
					instance.Ip,
					instance.Port,
				)
				conn := s.makeConnection(instance)
				_, err := conn.Use(context.Background(), keyValue)
				utils.CheckError(err)
			}
		}
	}
	return &Nil{}, nil
}
func (s *Server) ClearDist(ctx context.Context, data *Data) (*Nil, error) {
	log.Printf(
		"Clear %s:%d %v",
		data.FromIp,
		data.FromPort,
		data.Distribute,
	)
	s.partionCache.ClearAll(false)
	if data.Distribute {
		log.Print("Clear distribute is true")
		data.Distribute = false
		for _, instance := range s.config.Instances {
			if instance.Ip != data.FromIp && instance.Port != int(data.FromPort) {
				log.Printf(
					"Clear distribute to %s:%d",
					instance.Ip,
					instance.Port,
				)
				conn := s.makeConnection(instance)
				_, err := conn.ClearDist(context.Background(), data)
				utils.CheckError(err)
			}
		}
	}
	return &Nil{}, nil
}
func (s *Server) Check(ctx context.Context, _ *Nil) (*Nil, error) {
	return &Nil{}, nil
}
func (s *Server) ImAlive(ctx context.Context, data *Data) (*Nil, error) {
	log.Printf(
		"ImAlive %s:%d %v",
		data.FromIp,
		data.FromPort,
		data.Distribute,
	)
	s.checkAgain = true
	if data.Distribute {
		log.Print("ImAlive distribute is true")
		data.Distribute = false
		for _, instance := range s.config.Instances {
			if instance.Ip != data.FromIp && instance.Port != int(data.FromPort) {
				log.Printf(
					"ImAlive distribute to %s:%d",
					instance.Ip,
					instance.Port,
				)
				conn := s.makeConnection(instance)
				_, err := conn.ImAlive(context.Background(), data)
				utils.CheckError(err)
			}
		}
	}
	return &Nil{}, nil
}

// func RunServerDistribution() {
// 	log.Printf("Starting server on port %d", apiConfig.Port)
// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", apiConfig.Port))
// 	if err != nil {
// 		log.Fatalf("Failed to listen on port %d: %v", apiConfig.Port, err)
// 	}
// 	s := Server{
// 		partionCache: partionCache,
// 	}
// 	grpcServer := grpc.NewServer()
// 	RegisterCacheServiceServer(grpcServer, &s)
// 	if err := grpcServer.Serve(lis); err != nil {
// 		log.Fatalf("Failed to serve gRPC server over port %d: %v", apiConfig.Port, err)
// 	}
// }
