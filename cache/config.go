package main

import (
	"io/ioutil"
	"log"

	"github.com/SmsS4/KeepIt/cache/db"
	"github.com/SmsS4/KeepIt/cache/utils"

	"gopkg.in/yaml.v2"
)

type Config struct {
	db db.DBConfig
}

func getConfig(configPath string) Config {
	log.Print("Get config")
	configFile, err := ioutil.ReadFile("config.yml")
	utils.CheckError(err)
	configMap := make(map[string]string)
	utils.CheckError(yaml.Unmarshal(configFile, &configMap))
	var config Config
	config.db = db.GetDBConfig(configMap)
	return config
}
