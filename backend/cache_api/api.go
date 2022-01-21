package cache_api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	server "github.com/SmsS4/KeepIt/cache/cache_server"
	"github.com/SmsS4/KeepIt/cache/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
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

func LoadTLSCredentials() (credentials.TransportCredentials, error) {
	pemServerCA, err := ioutil.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}
	config := &tls.Config{
		RootCAs: certPool,
	}
	return credentials.NewTLS(config), nil
}

func (api *CacheApi) createConnection(url string) (server.CacheServiceClient, error) {
	log.Printf("Create connection to %s", url)
	var conn *grpc.ClientConn
	cert, err := LoadTLSCredentials()
	utils.CheckError(err)
	conn, err = grpc.Dial(
		url,
		grpc.WithTransportCredentials(cert),
	)
	if err != nil {
		return nil, err
	}
	return server.NewCacheServiceClient(conn), err
}

func (api *CacheApi) findConnection() string {
	for _, instance := range api.config.instances {
		log.Printf("Connecting to %s", instance.GetUrl())
		connection, err := api.createConnection(instance.GetUrl())
		if err == nil {
			response, err := connection.Check(
				context.Background(),
				&server.Nil{},
			)
			log.Printf("Response is %s", response)
			if err == nil {
				log.Print("Connected")
				return instance.GetUrl()
			}
		}
		log.Print("Failed")
	}
	log.Fatal("Could not find any connection")
	return ""
}

func (api *CacheApi) FindNewActive(err error) {
	st := status.Convert(err)
	if st.Code() == codes.Unavailable {
		api.currentInstance = api.findConnection()
	} else {
		log.Fatalf("Error %s", err)
	}

}

func (api *CacheApi) Get(key string) string {
	conn, err := api.createConnection(
		api.currentInstance,
	)
	utils.CheckError(err)
	response, err := conn.Get(
		context.Background(),
		&server.Key{
			Key: key,
		},
	)
	log.Printf("Get response %s", response)
	st := status.Convert(err)
	if st.Code() == codes.Unavailable {
		log.Printf("network error")
		api.FindNewActive(err)
		return api.Get(key)
	}

	if err != nil {
		return ""
	}
	if response.ActiveIp != "" {
		api.currentInstance = response.ActiveIp
		return api.Get(key)
	}
	return response.Value
}

func (api *CacheApi) Put(key string, value string) *server.OprationResult {
	conn, err := api.createConnection(
		api.currentInstance,
	)
	utils.CheckError(err)
	response, err := conn.Put(
		context.Background(),
		&server.KeyValue{
			Key:   key,
			Value: value,
		},
	)
	log.Printf("Put response %s", response)
	if err != nil {
		api.FindNewActive(err)
		return api.Put(key, value)
	}
	if response.ActiveIp != "" {
		api.currentInstance = response.ActiveIp
		return api.Put(key, value)
	}
	return response
}

func (api *CacheApi) Clear() *server.OprationResult {
	conn, err := api.createConnection(
		api.currentInstance,
	)
	utils.CheckError(err)
	response, err := conn.Clear(
		context.Background(),
		&server.Nil{},
	)
	log.Printf("Clear response %s", response)
	if err != nil {
		api.FindNewActive(err)
		return api.Clear()
	}
	if response.ActiveIp != "" {
		api.currentInstance = response.ActiveIp
		return api.Clear()
	}

	return response
}
