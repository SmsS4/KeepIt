package cache_api

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/SmsS4/KeepIt/cache/utils"
)

type Instance struct {
	Ip       string
	Port     int
	Priority int
}

func ParseInsances(data string) []Instance {
	instancesStrings := strings.Split(data, ",")
	instances := make([]Instance, len(instancesStrings))
	for i := 0; i < len(instancesStrings); i++ {
		ipAndPort := strings.Split(instancesStrings[i], ":")
		ip := ipAndPort[0]
		port, err := strconv.Atoi(ipAndPort[1])
		utils.CheckError(err)
		priority, err := strconv.Atoi(ipAndPort[2])
		utils.CheckError(err)
		log.Printf("Instance %d is on %s:%d with priority %d", i, ip, port, priority)
		instances[i] = Instance{
			Ip:       ip,
			Port:     port,
			Priority: priority,
		}
	}
	sort.Slice(instances, func(i, j int) bool {
		return instances[i].Priority < instances[j].Priority
	})
	return instances
}
func (instance *Instance) GetUrl() string {
	return fmt.Sprintf("%s:%d", instance.Ip, instance.Port)
}

type CacheConfig struct {
	instances []Instance
}

func GetCacheConfig(data map[string]string) CacheConfig {
	log.Printf("Getting cache config %s", data)
	return CacheConfig{
		instances: ParseInsances(data["instances"]),
	}
}
