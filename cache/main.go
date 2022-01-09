package main

import (
	"fmt"
	"log"

	"github.com/SmsS4/KeepIt/cache/db"
)

func main() {
	log.Print("Starting server...")
	config := getConfig("config.yml")
	dbConn := db.CreateConnection(config.db)
	dbConn.SetValue("test_key", "test_value2")
	kv, e := dbConn.GetValue("test_key")
	fmt.Print(kv.Value)
	dbConn.SetValue("test_key", "test_value1")
	kv, e = dbConn.GetValue("test_key")
	fmt.Print(kv.Value)
	e = e

	// ll := ds.NewLinkList()
	// node1 := ll.AppendValue("test")
	// ll.AppendValue("demol")
	// node3 := ll.AppendValue("hello")
	// node1 = node1
	// node3 = node3
	// ll.MoveToTail(node1)
	// ll.MoveToTail(node3)
	// ll.PopHead()
}
