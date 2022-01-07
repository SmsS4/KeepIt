package db

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/SmsS4/KeepIt/cache/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDBConfig(data map[string]string) DBConfig {
	var config DBConfig
	config.Host = data["db_host"]
	config.Username = data["db_username"]
	config.Password = data["db_password"]
	config.Name = data["db_name"]
	db_port, err := strconv.Atoi(data["db_port"])
	utils.CheckError(err)
	config.Port = db_port
	log.Printf(
		"DB Config: %s:%d user:%s name:%s",
		config.Host,
		config.Port,
		config.Username,
		config.Name,
	)
	return config
}

var dbConnection *gorm.DB

func createTables() {
	log.Print("Creating tables")
	dbConnection.AutoMigrate(keyValue)
}

func GetValue(key string) (KeyValue, error) {
	log.Printf("Get key: %s", key)
	var result KeyValue
	report := dbConnection.First(&result, "key = ?", key)
	if errors.Is(report.Error, gorm.ErrRecordNotFound) {
		log.Printf("Get key: %s -> error: value not found", key)
		return result, gorm.ErrRecordNotFound
	}
	return result, report.Error
}

func SetValue(key string, value string) {
	dbConnection.Create(&KeyValue{Key: key, Value: value})
}

func CreateConnection(dbConfig DBConfig) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		dbConfig.Host,
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Name,
		dbConfig.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	utils.CheckError(err)
	dbConnection = db

	createTables()
}
