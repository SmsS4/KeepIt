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
	// config = config
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
		config.cacheConfig,
		db.CreateConnection(config.dbConfig),
	)

	go server.RunServerCache(config.apiConfig, &partionCache)
	time.Sleep(time.Second * 1)
	log.Print("start client")
	server.RunApi()

	partionCache.Put("test1", "hello1")
	partionCache.Put("test2", "hello2")
	partionCache.Put("test3", "hello3")
	partionCache.Put("test4", "hello4")
	partionCache.Put("test5", "hello5")
	partionCache.Put("test6", "hello6")
	partionCache.Put("test7", "hello7")
	partionCache.Put("test8", "hello8")
	partionCache.Put("test9", "hello9")
	partionCache.Get("test1")
	partionCache.ClearAll()
	partionCache.Print()
}
