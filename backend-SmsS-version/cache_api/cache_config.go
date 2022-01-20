package cache_api

import (
	"log"

	server "github.com/SmsS4/KeepIt/cache/cache_server"
)

type CacheConfig struct {
	instances []server.Instance
}

func GetCacheConfig(data map[string]string) CacheConfig {
	log.Printf("Getting cache config %s", data)
	return CacheConfig{
		instances: server.ParseInsances(data["instances"]),
	}
}
