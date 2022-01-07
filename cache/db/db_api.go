package db

import (
	"fmt"
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
	return config
}

var dbConnection *gorm.DB

func createTables() {
	dbConnection.AutoMigrate(&User{})
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
