package db

import "time"

type DBConfig struct {
	Host     string
	Username string
	Password string
	Name     string
	Port     int
}

type KeyValue struct {
	Key       string `gorm:"primaryKey;type:varchar(64)"`
	Value     string `gorm:"type:varchar(2048)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var keyValue KeyValue
