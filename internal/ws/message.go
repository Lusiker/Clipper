package ws

const (
	MessageTypeClipCreated   = "clip_created"
	MessageTypeClipDeleted   = "clip_deleted"
	MessageTypeDeviceOnline  = "device_online"
	MessageTypeDeviceOffline = "device_offline"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type ClipCreatedData struct {
	Clip interface{} `json:"clip"`
}

type ClipDeletedData struct {
	ClipID string `json:"clip_id"`
}

type DeviceStatusData struct {
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	IsOnline   bool   `json:"is_online"`
}