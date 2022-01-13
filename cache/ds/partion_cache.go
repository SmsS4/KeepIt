package ds

import (
	"hash/fnv"
	"log"
	"strconv"
	"sync"

	"github.com/SmsS4/KeepIt/cache/db"
	"github.com/SmsS4/KeepIt/cache/utils"
)

type CacheConfig struct {
	MaxSize       int
	PartionsCount int
}

func GetCacheConfig(configMap map[string]string) CacheConfig {
	log.Print("Getting cache config")
	maxSize, err := strconv.Atoi(configMap["max_size"])
	utils.CheckError(err)
	partionsCount, err := strconv.Atoi(configMap["partions_count"])
	utils.CheckError(err)
	log.Printf("Got cache config MaxSize: %d, PartionsCount: %d", maxSize, partionsCount)
	return CacheConfig{
		MaxSize:       maxSize,
		PartionsCount: partionsCount,
	}
}

type PartionCache struct {
	MaxSize       int
	PartionsCount int
	caches        []Cache
	db            *db.DbConnection
	lock          *sync.Mutex
}

func NewPartionCache(Config CacheConfig, db *db.DbConnection) PartionCache {
	caches := make([]Cache, Config.PartionsCount)
	for i := 0; i < Config.PartionsCount; i++ {
		caches[i] = NewCache(
			Config.MaxSize,
			db,
		)
	}
	return PartionCache{
		MaxSize:       Config.MaxSize,
		PartionsCount: Config.PartionsCount,
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
	go partionCache.db.Clear()
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
