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

type DbConnection struct {
	conn *gorm.DB
}

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

func (db *DbConnection) createTables() {
	log.Print("Creating tables")
	db.conn.AutoMigrate(keyValue)
}

func (db *DbConnection) GetValue(key string) (KeyValue, error) {
	log.Printf("Get key: %s", key)
	var result KeyValue
	report := db.conn.First(&result, "key = ?", key)
	if errors.Is(report.Error, gorm.ErrRecordNotFound) {
		log.Printf("Get key: %s -> error: value not found", key)
		return result, gorm.ErrRecordNotFound
	}
	return result, report.Error
}

func (db *DbConnection) SetValue(key string, value string) {
	data := KeyValue{Key: key, Value: value}
	if db.conn.Model(&data).Where("Key = ?", key).Updates(&data).RowsAffected == 0 {
		log.Printf("Set key: %s to %s", key, value)
		db.conn.Create(&data)
	} else {
		log.Printf("Update key: %s to %s", key, value)
	}
}

func CreateConnection(dbConfig DBConfig) *DbConnection {
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
	conn := DbConnection{conn: db}
	conn.createTables()
	return &conn
}

func (db *DbConnection) Clear() {
	db.conn.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&KeyValue{})
}
