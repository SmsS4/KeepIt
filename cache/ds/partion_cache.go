package ds

import (
	"hash/fnv"
	"log"
	"sync"

	"github.com/SmsS4/KeepIt/cache/db"
)

type PartionCache struct {
	MaxSize       int
	PartionsCount int
	caches        []Cache
	db            *db.DbConnection
	lock          *sync.Mutex
}

func NewPartionCache(MaxSize int, PartionsCount int, db *db.DbConnection) PartionCache {
	caches := make([]Cache, PartionsCount)
	for i := 0; i < PartionsCount; i++ {
		caches[i] = NewCache(
			MaxSize,
			db,
		)
	}
	return PartionCache{
		MaxSize:       MaxSize,
		PartionsCount: PartionsCount,
		caches:        caches,
		db:            db,
		lock:          new(sync.Mutex),
	}
}

func (partionCache *PartionCache) hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32() % uint32(partionCache.PartionsCount)
}

func (partionCache *PartionCache) getCache(key string) *Cache {
	return &partionCache.caches[partionCache.hash(key)]
}

func (partionCache *PartionCache) ClearAll() {
	partionCache.db.Clear()
	for _, cache := range partionCache.caches {
		cache.Clear()
	}
}

func (partionCache *PartionCache) Get(key string) (string, bool, error) {
	return partionCache.getCache(key).Get(key)
}

func (partionCache *PartionCache) Put(key string, value string) {
	partionCache.getCache(key).Put(key, value)
}

func (partionCache *PartionCache) Print() {
	log.Printf("Number of caches: %d MaxSize: %d", partionCache.PartionsCount, partionCache.MaxSize)
	for i, cache := range partionCache.caches {
		log.Printf("Cache-%d:", i)
		cache.linkList.Print()
		log.Print(cache.keyToNode)
	}
}
