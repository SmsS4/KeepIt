package server

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func RunApi() {

	var conn *grpc.ClientConn
	cert, err := LoadTLSCredentials()
	conn, err = grpc.Dial(":7000", grpc.WithTransportCredentials(cert))
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}
	defer conn.Close()

	c := NewCacheServiceClient(conn)

	// key := Key{
	// 	Key: "notExists",
	// }

	// response, err := c.Get(context.Background(), &key)
	// log.Print(err)
	// if err == nil {
	// 	log.Printf(
	// 		"Response get %s, value:%s active:%s miss:%v",
	// 		key.Key,
	// 		response.Value,
	// 		response.ActiveIp,
	// 		response.MissCache,
	// 	)
	// 	log.Fatal("x")
	// }

	// keyValue := KeyValue{
	// 	Key:   "Exists",
	// 	Value: "Valve",
	// }

	// responsePut, errPut := c.Put(context.Background(), &keyValue)
	// if errPut != nil {
	// 	log.Fatalf("Error when calling SayHello: %s", errPut)
	// }

	// log.Printf("Response put %s %s, %s", keyValue.Key, keyValue.Value, responsePut.ActiveIp)

	key := Key{
		Key: "Exists",
	}
	// go c.Get(context.Background(), &key)
	// go c.Get(context.Background(), &key)
	// go c.Get(context.Background(), &key)
	// go c.Get(context.Background(), &key)
	response, err := c.Get(context.Background(), &key)
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}

	log.Printf("Response get %s, value:%s active:%s miss:%v", key.Key, response.Value, response.ActiveIp, response.MissCache)

	// key = Key{
	// 	Key: "Exists",
	// }

	// response, err = c.Get(context.Background(), &key)
	// if err != nil {
	// 	log.Fatalf("Error when calling SayHello: %s", err)
	// }

	// log.Printf("Response get %s, value:%s active:%s miss:%v", key.Key, response.Value, response.ActiveIp, response.MissCache)

}
