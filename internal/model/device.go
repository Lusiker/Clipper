package model

import "time"

type Device struct {
	ID       string    `json:"id"`
	UserID   string    `json:"user_id"`
	Name     string    `json:"name"`
	IP       string    `json:"ip"`
	LastSeen time.Time `json:"last_seen"`
	IsOnline bool      `json:"is_online"`
}