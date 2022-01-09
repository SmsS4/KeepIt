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

func (cp *PartionCache) hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32() % uint32(cp.PartionsCount)
}

func (cp *PartionCache) getCache(key string) *Cache {
	return &cp.caches[cp.hash(key)]
}

func (cp *PartionCache) ClearAll() {
	cp.db.Clear()
	for _, cache := range cp.caches {
		cache.Clear()
	}
}

func (cp *PartionCache) Get(key string) (string, bool, error) {
	return cp.getCache(key).Get(key)
}

func (cp *PartionCache) Put(key string, value string) {
	cp.getCache(key).Put(key, value)
}

func (cp *PartionCache) Print() {
	log.Printf("Number of caches: %d MaxSize: %d", cp.PartionsCount, cp.MaxSize)
	for i, cache := range cp.caches {
		log.Printf("Cache-%d:", i)
		cache.linkList.Print()
		log.Print(cache.keyToNode)
	}
}
