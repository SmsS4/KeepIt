package main

import (
	"io/ioutil"
	"log"

	server "github.com/SmsS4/KeepIt/cache/cache_server"
	"github.com/SmsS4/KeepIt/cache/db"
	"github.com/SmsS4/KeepIt/cache/ds"
	"github.com/SmsS4/KeepIt/cache/utils"

	"gopkg.in/yaml.v2"
)

type Config struct {
	dbConfig    db.DBConfig
	cacheConfig ds.CacheConfig
	apiConfig   server.ApiConfig
}

func getConfig(configPath string) Config {
	log.Print("Get config")
	configFile, err := ioutil.ReadFile("config.yml")
	utils.CheckError(err)
	configMap := make(map[string]map[string]string)
	utils.CheckError(yaml.Unmarshal(configFile, &configMap))
	return Config{
		dbConfig:    db.GetDBConfig(configMap["db"]),
		cacheConfig: ds.GetCacheConfig(configMap["cache"]),
		apiConfig:   server.GetApiConfig(configMap["api"]),
	}
}
