package main

import (
	"log"
	"os"

	server "github.com/SmsS4/KeepIt/cache/cache_server"
	"github.com/SmsS4/KeepIt/cache/db"
	"github.com/SmsS4/KeepIt/cache/ds"
)

func main() {
	server.RunApi()

	log.Print("Starting server...")
	log.Printf("Path is %s", os.Args[1])
	config := getConfig(os.Args[1])
	partionCache := ds.NewPartionCache(
		config.cacheConfig,
		db.CreateConnection(config.dbConfig),
	)
	server.RunServer(config.apiConfig, &partionCache)
}
