package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/lusiker/clipper/internal/service"
)

type Hub struct {
	clients    map[string]map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *BroadcastMessage
	mu         sync.RWMutex

	clipService   *service.ClipService
	deviceService *service.DeviceService
}

type BroadcastMessage struct {
	UserID  string
	Message []byte
	Exclude string
}

func NewHub(clipService *service.ClipService, deviceService *service.DeviceService) *Hub {
	return &Hub{
		clients:      make(map[string]map[string]*Client),
		register:     make(chan *Client, 10),
		unregister:   make(chan *Client, 10),
		broadcast:    make(chan *BroadcastMessage, 100),
		clipService:  clipService,
		deviceService: deviceService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[string]*Client)
			}
			h.clients[client.UserID][client.DeviceID] = client
			h.mu.Unlock()

			h.notifyDeviceStatus(client.UserID, client.DeviceID, client.DeviceName, true)

		case client := <-h.unregister:
			h.mu.Lock()
			if devices, ok := h.clients[client.UserID]; ok {
				if _, ok := devices[client.DeviceID]; ok {
					delete(devices, client.DeviceID)
					close(client.Send)
				}
				if len(devices) == 0 {
					delete(h.clients, client.UserID)
				}
			}
			h.mu.Unlock()

			h.deviceService.SetOffline(client.DeviceID)
			h.notifyDeviceStatus(client.UserID, client.DeviceID, client.DeviceName, false)

		case msg := <-h.broadcast:
			h.mu.RLock()
			devices := h.clients[msg.UserID]
			h.mu.RUnlock()

			for deviceID, client := range devices {
				if deviceID != msg.Exclude {
					select {
					case client.Send <- msg.Message:
					default:
						log.Printf("Client %s channel full", deviceID)
					}
				}
			}
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Broadcast(userID string, message []byte, excludeDevice string) {
	h.broadcast <- &BroadcastMessage{
		UserID:  userID,
		Message: message,
		Exclude: excludeDevice,
	}
}

func (h *Hub) notifyDeviceStatus(userID, deviceID, deviceName string, isOnline bool) {
	msgType := MessageTypeDeviceOnline
	if !isOnline {
		msgType = MessageTypeDeviceOffline
	}

	msg := Message{
		Type: msgType,
		Data: DeviceStatusData{
			DeviceID:   deviceID,
			DeviceName: deviceName,
			IsOnline:   isOnline,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.Broadcast(userID, data, "")
}

func (h *Hub) NotifyClipCreated(userID string, clip interface{}, excludeDevice string) {
	msg := Message{
		Type: MessageTypeClipCreated,
		Data: ClipCreatedData{Clip: clip},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.Broadcast(userID, data, excludeDevice)
}

func (h *Hub) NotifyClipDeleted(userID, clipID string, excludeDevice string) {
	msg := Message{
		Type: MessageTypeClipDeleted,
		Data: ClipDeletedData{ClipID: clipID},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.Broadcast(userID, data, excludeDevice)
}

func (h *Hub) GetOnlineDevices(userID string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if devices, ok := h.clients[userID]; ok {
		result := make([]string, 0, len(devices))
		for deviceID := range devices {
			result = append(result, deviceID)
		}
		return result
	}
	return nil
}