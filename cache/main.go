package main

import (
	"fmt"
	"log"

	"github.com/SmsS4/KeepIt/cache/db"
)

func main() {
	log.Print("Starting server...")
	config := getConfig("config.yml")
	db.CreateConnection(config.db)

	// db.SetValue("test_key", "test_value")
	kv, e := db.GetValue("test_key2")
	fmt.Print(kv.Value)
	fmt.Print(e)
}
