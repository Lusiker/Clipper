package model

import "time"

type ClipType string

const (
	ClipTypeText  ClipType = "text"
	ClipTypeImage ClipType = "image"
)

type Clip struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	DeviceID  string    `json:"device_id"`
	Type      ClipType  `json:"type"`
	Content   string    `json:"content"`
	Meta      string    `json:"meta,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ClipCreate struct {
	Type    ClipType `json:"type" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Meta    string   `json:"meta,omitempty"`
}

type ClipMeta struct {
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
	Size   int64  `json:"size,omitempty"`
}

type ClipMetaImage struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Size      int64  `json:"size"`
	Format    string `json:"format"`
	ThumbPath string `json:"thumb_path"`
}