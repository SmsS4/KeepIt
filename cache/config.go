package main

import (
	"io/ioutil"
	"log"
	"strconv"

	"github.com/SmsS4/KeepIt/cache/db"
	"github.com/SmsS4/KeepIt/cache/utils"

	"gopkg.in/yaml.v2"
)

type CacheConfig struct {
	MaxSize       int
	PartionsCount int
}

func GetCacheConfig(configMap map[string]string) CacheConfig {
	maxSize, err := strconv.Atoi(configMap["cache_max_size"])
	utils.CheckError(err)
	partionsCount, err := strconv.Atoi(configMap["cache_partions_count"])
	utils.CheckError(err)
	return CacheConfig{
		MaxSize:       maxSize,
		PartionsCount: partionsCount,
	}
}

type ApiConfig struct {
	Port int
}

func GetApiConfig(configMap map[string]string) ApiConfig {
	port, err := strconv.Atoi(configMap["api_port"])
	utils.CheckError(err)
	return ApiConfig{
		Port: port,
	}
}

type Config struct {
	db          db.DBConfig
	cacheConfig CacheConfig
	apiConfig   ApiConfig
}

func getConfig(configPath string) Config {
	log.Print("Get config")
	configFile, err := ioutil.ReadFile("config.yml")
	utils.CheckError(err)
	configMap := make(map[string]string)
	utils.CheckError(yaml.Unmarshal(configFile, &configMap))
	return Config{
		db:          db.GetDBConfig(configMap),
		cacheConfig: GetCacheConfig(configMap),
		apiConfig:   GetApiConfig(configMap),
	}
}
