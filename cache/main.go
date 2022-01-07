package main

import (
	"log"

	"github.com/SmsS4/KeepIt/cache/db"
)

func main() {
	log.Print("Starting server...")
	config := getConfig("config.yml")
	db.CreateConnection(config.db)

	// db.SetValue("test_key", "test_value2")
	// kv, e := db.GetValue("test_key")
	// fmt.Print(kv.Value)
	// db.SetValue("test_key", "test_value1")
	// kv, e = db.GetValue("test_key")
	// fmt.Print(kv.Value)
	// e = e
}
