package server

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

type ApiConfig struct {
	Ip        string
	Port      int
	Instances []Instance
}

func (instance *Instance) GetUrl() string {
	return fmt.Sprintf("%s:%d", instance.Ip, instance.Port)
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

func GetApiConfig(configMap map[string]string) ApiConfig {
	log.Print("Getting api config")
	port, err := strconv.Atoi(configMap["port"])
	utils.CheckError(err)
	log.Printf("Got api config, port: %d", port)
	instances := ParseInsances(configMap["instances"])
	for i := 1; i < len(instances); i++ {
		if instances[i].Priority == instances[i-1].Priority {
			log.Fatalf(
				"Two instances have same priority: %s:%d:%d and %s:%d:%d",
				instances[i].Ip,
				instances[i].Port,
				instances[i].Priority,
				instances[i-1].Ip,
				instances[i-1].Port,
				instances[i-1].Priority,
			)
		}
	}
	return ApiConfig{
		Ip:        configMap["ip"],
		Port:      port,
		Instances: instances,
	}
}
