package main

import (
	"fmt"
	"log"
	"strconv"

	"io/ioutil"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	host     string
	username string
	password string
	name     string
	port     int
}

type Config struct {
	db DBConfig
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func getDBConfig(data map[string]string) DBConfig {
	var config DBConfig
	config.host = data["db_host"]
	config.username = data["db_username"]
	config.password = data["db_password"]
	config.name = data["db_name"]
	db_port, err := strconv.Atoi(data["db_port"])
	checkError(err)
	config.port = db_port
	return config
}

func getConfig(configPath string) Config {
	log.Print("Get config")
	configFile, err := ioutil.ReadFile("config.yml")
	checkError(err)
	configMap := make(map[string]string)
	checkError(yaml.Unmarshal(configFile, &configMap))
	var config Config
	config.db = getDBConfig(configMap)
	log.Printf(
		"DB Config: %s:%d user:%s name:%s",
		config.db.host,
		config.db.port,
		config.db.username,
		config.db.name,
	)
	return config
}

func main() {
	log.Print("Starting server...")
	config := getConfig("config.yml")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		config.db.host,
		config.db.username,
		config.db.password,
		config.db.name,
		config.db.port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	checkError(err)
}
