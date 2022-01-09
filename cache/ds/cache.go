package ds

import (
	"log"
	"sync"

	"github.com/SmsS4/KeepIt/cache/db"
)

type Cache struct {
	MaxSize   int
	linkList  LinkList
	keyToNode map[string]*Node
	db        *db.DbConnection
	lock      *sync.Mutex
}

func NewCache(MaxSize int, db *db.DbConnection) Cache {
	return Cache{
		MaxSize:   MaxSize,
		linkList:  NewLinkList(),
		keyToNode: make(map[string]*Node),
		db:        db,
		lock:      new(sync.Mutex),
	}
}

func (cache *Cache) relax() {
	if cache.linkList.Size >= cache.MaxSize {
		log.Print("Relaxing cache")
		/// > cache.MaxSize/2 to reduce locks
		for cache.linkList.Size > cache.MaxSize {
			node := cache.linkList.PopHead()
			delete(cache.keyToNode, node.Value.Key)
		}
	}
}

func (cache *Cache) addToMap(key string, value string) {
	log.Printf("Add %s:%s to cache map", key, value)
	node := cache.linkList.AppendValue(Pair{key, value})
	cache.keyToNode[node.Value.Key] = node
	cache.relax()
}

func (cache *Cache) Put(key string, value string) {
	log.Printf("Put %s:%s to db", key, value)
	go cache.db.SetValue(key, value)
	if node, ok := cache.keyToNode[key]; ok {
		cache.linkList.MoveToTail(node)
		node.Value = Pair{key, value}
	} else {
		cache.addToMap(key, value)
	}
}

func (cache *Cache) Get(key string) (string, bool, error) {
	if node, ok := cache.keyToNode[key]; ok {
		log.Printf("Get %s from cache", key)
		cache.linkList.MoveToTail(node)
		return node.Value.Value, true, nil
	} else {
		log.Printf("Get %s from db", key)
		value, err := cache.db.GetValue(key)
		cache.addToMap(value.Key, value.Value)
		return value.Value, false, err
	}
}

func (cache *Cache) Clear() {
	log.Print("Clear cache")
	cache.linkList = NewLinkList()
	cache.keyToNode = make(map[string]*Node)
}
