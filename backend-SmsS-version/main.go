package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/SmsS4/KeepIt/backend-SmsS-version/cache_api"
	"github.com/SmsS4/KeepIt/cache/utils"
	"gopkg.in/yaml.v2"
)

type Config struct {
	CacheApi cache_api.CacheConfig
}

func getConfig(configPath string) Config {
	log.Print("Getting config")
	configFile, err := ioutil.ReadFile(configPath)
	utils.CheckError(err)
	configMap := make(map[string]map[string]string)
	utils.CheckError(yaml.Unmarshal(configFile, &configMap))
	return Config{
		CacheApi: cache_api.GetCacheConfig(configMap["cache"]),
	}
}

func main() {
	config := getConfig(os.Args[1])
	api := cache_api.CreateApi(config.CacheApi)
	for {
		api.Put("key", "value")
		response := api.Get("key")
		log.Printf("Response put %s", response)
		response = api.Get("key")
		log.Printf("Response put %s", response)
		api.Clear()
	}

}
