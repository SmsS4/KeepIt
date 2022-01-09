package server

import (
	"strconv"

	"github.com/SmsS4/KeepIt/cache/utils"
)

type ApiConfig struct {
	Port int
}

func GetApiConfig(configMap map[string]string) ApiConfig {
	port, err := strconv.Atoi(configMap["api_port"])
	utils.CheckError(err)
	return ApiConfig{
		Port: port,
	}
}
