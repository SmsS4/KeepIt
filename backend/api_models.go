package main

import (
	"github.com/SmsS4/KeepIt/backend/cache_api"
)

type GatewayConfig struct {
	Port               string
	RateLimitPerMinute int
}

type Config struct {
	CacheApi      cache_api.CacheConfig
	GatewayConfig GatewayConfig
}

type UserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type NewNoteInput struct {
	Note string `json:"note" binding:"required"`
}

type NotesListInput struct {
	Username string `json:"username"`
}

type UpdateNoteInput struct {
	Note_id string `json:"note_id" binding:"required"`
	Note    string `json:"note" binding:"required"`
}
