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
	log.Printf(
		"CheckActiveIp active: %s, my:%s, checkAgain: %v",
		s.activeIp,
		s.myIp,
		s.checkAgain,
	)
	if s.activeIp == s.myIp && !s.checkAgain {
		log.Print("False")
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

func (s *Server) GetDistrubter() CacheServiceClient {
	return s.makeConnection(s.config.Instances[len(s.config.Instances)-1])
}

func (s *Server) getData() *Data {
	return &Data{
		FromIp:     s.config.Ip,
		FromPort:   int32(s.config.Port),
		Distribute: true,
	}
}

func (s *Server) UseOthers(key string, value string) {
	log.Print("Distribute")
	_, err := s.GetDistrubter().Use(
		context.Background(),
		&DistKeyValue{
			Key:   key,
			Value: value,
			Data:  s.getData(),
		},
	)
	log.Print(err)
}

func (s *Server) Get(ctx context.Context, key *Key) (*Result, error) {
	log.Printf("Get %s", key.Key)
	if s.CheckActiveIp() {
		log.Printf("Not active call: %s", s.activeIp)
		return &Result{
			Value:     "",
			MissCache: false,
			ActiveIp:  s.activeIp,
		}, nil
	}
	value, miss, err := s.partionCache.Get(key.Key)
	if err == nil {
		s.UseOthers(key.Key, value)
	}
	log.Printf("Value is: %s, miss: %v", value, miss)
	return &Result{
		Value:     value,
		MissCache: miss,
		ActiveIp:  "",
	}, err
}

func (s *Server) Put(ctx context.Context, keyValue *KeyValue) (*OprationResult, error) {
	log.Printf("Put %s:%s", keyValue.Key, keyValue.Value)
	if s.CheckActiveIp() {
		log.Printf("Not active call %s", s.activeIp)
		return &OprationResult{
			ActiveIp: s.activeIp,
		}, nil
	}
	s.partionCache.Put(keyValue.Key, keyValue.Value)
	s.UseOthers(keyValue.Key, keyValue.Value)
	log.Print("Puted")
	return &OprationResult{
		ActiveIp: "",
	}, nil
}

func (s *Server) Clear(ctx context.Context, _ *Nil) (*OprationResult, error) {
	log.Print("Clearing")
	if s.CheckActiveIp() {
		log.Printf("Not active call %s", s.activeIp)
		return &OprationResult{
			ActiveIp: s.activeIp,
		}, nil
	}
	log.Print("Clear all partions")
	s.partionCache.ClearAll(true)
	log.Print("Distribute")
	_, err := s.GetDistrubter().ClearDist(
		context.Background(),
		s.getData(),
	)
	log.Print(err)
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
	grpcServer := grpc.NewServer()
	RegisterCacheServiceServer(grpcServer, &gatewayServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port %d: %v", apiConfig.Port, err)
	}
}

func (s *Server) makeConnection(instance Instance) CacheServiceClient {
	log.Printf("makeConnection to instance %s:%d", instance.Ip, instance.Port)
	var conn *grpc.ClientConn
	addr := fmt.Sprintf("%s:%d", instance.Ip, instance.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not makeConnection: %s", err)
	}
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
			if instance.Ip != keyValue.Data.FromIp || instance.Port != int(keyValue.Data.FromPort) {
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
			if instance.Ip != data.FromIp || instance.Port != int(data.FromPort) {
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
	log.Printf("Check on %s", s.myIp)
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
			if instance.Ip != data.FromIp || instance.Port != int(data.FromPort) {
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
