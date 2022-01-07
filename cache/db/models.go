package db

import "time"

type DBConfig struct {
	Host     string
	Username string
	Password string
	Name     string
	Port     int
}

type User struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
