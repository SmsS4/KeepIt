package cache_api

import (
	"context"
	"log"

	server "github.com/SmsS4/KeepIt/cache/cache_server"
	"github.com/SmsS4/KeepIt/cache/utils"
	"google.golang.org/grpc"
)

type CacheApi struct {
	config          CacheConfig
	currentInstance string
}

func CreateApi(config CacheConfig) *CacheApi {
	api := &CacheApi{
		config: config,
	}
	api.currentInstance = config.instances[0].GetUrl()
	return api
}

func (api *CacheApi) createConnection(url string) server.CacheServiceClient {
	log.Printf("Create connection to %s", url)
	var conn *grpc.ClientConn
	cert, err := server.LoadTLSCredentials()
	utils.CheckError(err)
	conn, err = grpc.Dial(
		url,
		grpc.WithTransportCredentials(cert),
	)
	utils.CheckError(err)
	return server.NewCacheServiceClient(conn)
}

func (api *CacheApi) Get(key string) *server.Result {
	response, err := api.createConnection(
		api.currentInstance,
	).Get(
		context.Background(),
		&server.Key{
			Key: key,
		},
	)
	log.Printf("Get response %d", response)
	if err != nil {
		log.Printf("Error in get %s", err)
		/// todo: check type of error
		/// for now assume it's change ip
		api.currentInstance = response.ActiveIp
		return api.Get(key)
	}
	return response
}

func (api *CacheApi) Put(key string, value string) *server.OprationResult {
	response, err := api.createConnection(
		api.currentInstance,
	).Put(
		context.Background(),
		&server.KeyValue{
			Key:   key,
			Value: value,
		},
	)
	log.Printf("Put response %d", response)
	if err != nil {
		log.Printf("Error in put %s", err)
		/// todo: check type of error
		/// for now assume it's change ip
		api.currentInstance = response.ActiveIp
		return api.Put(key, value)
	}
	return response
}

func (api *CacheApi) Clear() *server.OprationResult {
	response, err := api.createConnection(
		api.currentInstance,
	).Clear(
		context.Background(),
		&server.Nil{},
	)
	log.Printf("Clear response %d", response)
	if err != nil {
		log.Printf("Error in clear %s", err)
		/// todo: check type of error
		/// for now assume it's change ip
		api.currentInstance = response.ActiveIp
		return api.Clear()
	}
	return response
}
