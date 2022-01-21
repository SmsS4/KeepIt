package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/SmsS4/KeepIt/backend/cache_api"
	"github.com/SmsS4/KeepIt/cache/utils"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"gopkg.in/yaml.v2"
)

func GenerateId() string {
	res := strconv.Itoa(rand.Intn(99999999-10000000) + 10000000)
	return res
}

func ParseToken(tokenStr string) string {

	//creating the token
	token, _ := jwt_lib.Parse(tokenStr, func(token *jwt_lib.Token) (interface{}, error) {
		return mysupersecretpassword, nil
	})
	claims, _ := token.Claims.(jwt_lib.MapClaims)
	username := claims["username"].(string)
	return username
}

func GetGatewayConfig(data map[string]string) GatewayConfig {
	port, err := strconv.Atoi(data["ratelimit-per-minute"])
	utils.CheckError(err)
	return GatewayConfig{
		Port:               data["port"],
		RateLimitPerMinute: port,
	}
}

func GetConfig(configPath string) Config {
	log.Print("Getting config")
	configFile, err := ioutil.ReadFile(configPath)
	utils.CheckError(err)
	configMap := make(map[string]map[string]string)
	utils.CheckError(yaml.Unmarshal(configFile, &configMap))
	return Config{
		CacheApi:      cache_api.GetCacheConfig(configMap["cache"]),
		GatewayConfig: GetGatewayConfig(configMap["gateway"]),
	}
}

func RelaxRatelimit(requests []int64) {
	timestamp := time.Now().Unix()
	for len(requests) > 0 {
		if (timestamp - requests[0]) > 60 {
			requests = requests[1:]
		}
	}
}

func CheckRatelimit(username string) bool {
	if _, ok := rateLimit[username]; !ok {
		rateLimit[username] = make([]int64, 0)
	}
	RelaxRatelimit(rateLimit[username])
	if len(rateLimit[username]) > config.GatewayConfig.RateLimitPerMinute {
		return false
	}
	rateLimit[username] = append(rateLimit[username], time.Now().Unix())
	return true
}
