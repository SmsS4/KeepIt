package main

import (
	"log"
	"time"

	server "github.com/SmsS4/KeepIt/cache/cache_server"
	"github.com/SmsS4/KeepIt/cache/db"
	"github.com/SmsS4/KeepIt/cache/ds"
)

func main() {
	log.Print("Starting server...")
	config := getConfig("config.yml")
	// dbConn := db.CreateConnection(config.db)
	// dbConn.SetValue("test_key", "test_value2")
	// kv, e := dbConn.GetValue("test_key")
	// fmt.Print(kv.Value)
	// dbConn.SetValue("test_key", "test_value1")
	// kv, e = dbConn.GetValue("test_key")
	// fmt.Print(kv.Value)
	// e = e

	// ll := ds.NewLinkList()
	// node1 := ll.AppendValue("test")
	// ll.AppendValue("demol")
	// node3 := ll.AppendValue("hello")
	// node1 = node1
	// node3 = node3
	// ll.MoveToTail(node1)
	// ll.MoveToTail(node3)
	// ll.PopHead()

	partionCache := ds.NewPartionCache(
		config.cacheConfig.MaxSize,
		config.cacheConfig.PartionsCount,
		db.CreateConnection(config.db),
	)

	go server.RunServer(config.apiConfig, &partionCache)
	time.Sleep(time.Second * 1)
	log.Print("start client")
	server.RunApi()

	// cp.Put("test1", "hello1")
	// cp.Put("test2", "hello2")
	// cp.Put("test3", "hello3")
	// cp.Put("test4", "hello4")
	// cp.Put("test5", "hello5")
	// cp.Put("test6", "hello6")
	// cp.Put("test7", "hello7")
	// cp.Put("test8", "hello8")
	// cp.Put("test9", "hello9")
	// cp.Get("test1")
	// cp.Print()
}
